#!/usr/bin/env bash
# Distribute registry CA + registries.yaml to all RKE2 nodes, restart services,
# and ensure containerd hosts.toml exists (DataEngine guide; no ClusterShell).
#
# Required env (plus SSH_*):
#   REGISTRY_HOST REGISTRY_PORT REGISTRY_SSH_HOST
#   CA_FILE MASTER_IP
#   RKE2_HOSTS          space-separated management IPs (control+worker)
#   CONTROL_HOSTS       space-separated control-plane management IPs
#   WORKER_HOSTS        space-separated worker management IPs (may be empty)
#   INSECURE_SKIP_VERIFY true|false
#   RESTART_RKE2        true|false
#   MOUNTPOINT          cert path on registry VM
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

REGISTRY_HOST="${REGISTRY_HOST:?}"
REGISTRY_PORT="${REGISTRY_PORT:?}"
REGISTRY_SSH_HOST="${REGISTRY_SSH_HOST:?}"
CA_FILE="${CA_FILE:?}"
MASTER_IP="${MASTER_IP:?}"
RKE2_HOSTS="${RKE2_HOSTS:?}"
CONTROL_HOSTS="${CONTROL_HOSTS:?}"
WORKER_HOSTS="${WORKER_HOSTS:-}"
INSECURE_SKIP_VERIFY="${INSECURE_SKIP_VERIFY:-false}"
RESTART_RKE2="${RESTART_RKE2:-true}"
MOUNTPOINT="${MOUNTPOINT:?}"

REMOTE_CERT="${MOUNTPOINT}/certs/domain.crt"
LOCAL_CA="/tmp/registry-ca.crt"
LOCAL_REG="/tmp/registries.yaml"
HOSTS_DIR="/var/lib/rancher/rke2/agent/etc/containerd/certs.d/${REGISTRY_HOST}:${REGISTRY_PORT}"
HOSTS_FILE="${HOSTS_DIR}/hosts.toml"

echo "Fetching CA from registry ${REGISTRY_SSH_HOST}:${REMOTE_CERT}"
"$ROOT/scripts/ssh_run.sh" "$REGISTRY_SSH_HOST" \
  "openssl x509 -in '${REMOTE_CERT}' -noout -subject -dates -fingerprint -sha256 -ext subjectAltName"
"$ROOT/scripts/ssh_run.sh" "$REGISTRY_SSH_HOST" "sudo cat '${REMOTE_CERT}'" >"${LOCAL_CA}"

echo "Verifying TLS to registry with staged CA (from manager node; Multipass may block host→VM:5000)"
openssl x509 -in "${LOCAL_CA}" -noout -fingerprint -sha256
"$ROOT/scripts/ssh_copy.sh" "$MASTER_IP" "${LOCAL_CA}" /tmp/registry-ca-verify.crt
"$ROOT/scripts/ssh_run.sh" "$MASTER_IP" "
  set -euo pipefail
  openssl s_client \
    -connect '${REGISTRY_HOST}:${REGISTRY_PORT}' \
    -servername '${REGISTRY_HOST}' \
    -CAfile /tmp/registry-ca-verify.crt \
    -verify_return_error </dev/null >/dev/null
  rm -f /tmp/registry-ca-verify.crt
  echo 'TLS verify OK from manager'
"

EXPECTED_FP="$(openssl x509 -in "${LOCAL_CA}" -noout -fingerprint -sha256)"

install_on_host() {
  local host=$1
  "$ROOT/scripts/ssh_copy.sh" "$host" "${LOCAL_CA}" /tmp/registry-ca.crt
  "$ROOT/scripts/ssh_run.sh" "$host" "
    set -euo pipefail
    sudo install -d -m 0755 /etc/rancher/rke2/certs
    sudo install -m 0644 /tmp/registry-ca.crt '${CA_FILE}'
    FP=\$(sudo openssl x509 -in '${CA_FILE}' -noout -fingerprint -sha256)
    echo \"\$(hostname): \$FP\"
    rm -f /tmp/registry-ca.crt
    test \"\$FP\" = '${EXPECTED_FP}'
  "
}

echo "Installing CA on RKE2 nodes"
for host in $RKE2_HOSTS; do
  install_on_host "$host"
done

SKIP_YAML="false"
if [ "${INSECURE_SKIP_VERIFY}" = "true" ]; then
  SKIP_YAML="true"
fi

cat >"${LOCAL_REG}" <<EOF
mirrors:
  "${REGISTRY_HOST}:${REGISTRY_PORT}":
    endpoint:
      - "https://${REGISTRY_HOST}:${REGISTRY_PORT}"
configs:
  "${REGISTRY_HOST}:${REGISTRY_PORT}":
    tls:
      ca_file: ${CA_FILE}
      insecure_skip_verify: ${SKIP_YAML}
EOF

echo "Writing registries.yaml:"
cat "${LOCAL_REG}"

write_registries() {
  local host=$1
  "$ROOT/scripts/ssh_copy.sh" "$host" "${LOCAL_REG}" /tmp/registries.yaml
  "$ROOT/scripts/ssh_run.sh" "$host" "
    set -euo pipefail
    sudo mkdir -p /etc/rancher/rke2
    if [ -f /etc/rancher/rke2/registries.yaml ]; then
      sudo cp /etc/rancher/rke2/registries.yaml /etc/rancher/rke2/registries.yaml.pre-ca-trust
    fi
    sudo install -m 0644 /tmp/registries.yaml /etc/rancher/rke2/registries.yaml
    rm -f /tmp/registries.yaml
  "
}

for host in $RKE2_HOSTS; do
  write_registries "$host"
done

restart_group() {
  local svc=$1
  shift
  local hosts=$*
  for host in $hosts; do
    echo "Restarting ${svc} on ${host}"
    "$ROOT/scripts/ssh_run.sh" "$host" "sudo systemctl restart ${svc}"
  done
}

if [ "${RESTART_RKE2}" = "true" ]; then
  restart_group rke2-server $CONTROL_HOSTS
  if [ -n "${WORKER_HOSTS}" ]; then
    restart_group rke2-agent $WORKER_HOSTS
  fi
  echo "Waiting for API after restart..."
  for i in $(seq 1 90); do
    if "$ROOT/scripts/ssh_run.sh" "$MASTER_IP" \
      "sudo /var/lib/rancher/rke2/bin/kubectl --kubeconfig /etc/rancher/rke2/rke2.yaml get nodes >/dev/null 2>&1"; then
      break
    fi
    sleep 10
  done
fi

echo "Checking hosts.toml on each node"
MISSING=()
for host in $RKE2_HOSTS; do
  if "$ROOT/scripts/ssh_run.sh" "$host" "sudo test -f '${HOSTS_FILE}'"; then
    echo "${host}: hosts.toml present"
    "$ROOT/scripts/ssh_run.sh" "$host" "sudo awk '1' '${HOSTS_FILE}'" || true
  else
    echo "${host}: hosts.toml MISSING"
    MISSING+=("$host")
  fi
done

if [ "${#MISSING[@]}" -gt 0 ]; then
  echo "Seeding hosts.toml from manager ${MASTER_IP}"
  "$ROOT/scripts/ssh_run.sh" "$MASTER_IP" "
    set -euo pipefail
    sudo test -f '${HOSTS_FILE}'
    sudo install -m 0644 '${HOSTS_FILE}' /tmp/rke2-hosts.toml.seed
  "
  "$ROOT/scripts/ssh_run.sh" "$MASTER_IP" "cat /tmp/rke2-hosts.toml.seed" >/tmp/rke2-hosts.toml.seed
  "$ROOT/scripts/ssh_run.sh" "$MASTER_IP" "rm -f /tmp/rke2-hosts.toml.seed"

  for host in "${MISSING[@]}"; do
    "$ROOT/scripts/ssh_copy.sh" "$host" /tmp/rke2-hosts.toml.seed /tmp/rke2-hosts.toml.seed
    "$ROOT/scripts/ssh_run.sh" "$host" "
      set -euo pipefail
      sudo install -d -m 0700 '${HOSTS_DIR}'
      sudo install -m 0600 /tmp/rke2-hosts.toml.seed '${HOSTS_FILE}'
      rm -f /tmp/rke2-hosts.toml.seed
      echo 'hosts.toml installed'
    "
  done
  rm -f /tmp/rke2-hosts.toml.seed

  if [ "${RESTART_RKE2}" = "true" ]; then
    restart_group rke2-server $CONTROL_HOSTS
    if [ -n "${WORKER_HOSTS}" ]; then
      restart_group rke2-agent $WORKER_HOSTS
    fi
  fi
fi

rm -f "${LOCAL_CA}" "${LOCAL_REG}"
echo "Registry trust configuration complete"
