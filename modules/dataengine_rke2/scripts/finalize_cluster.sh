#!/usr/bin/env bash
# Post-install: kubeconfig on manager + local copy, wait Ready, label workers.
# Invoked by Terraform local-exec; configuration via environment variables.
#
# Required env:
#   SSH_* / RKE2_*     (from module ssh_env)
#   FIRST_HOST         first control-plane management IP
#   FIRST_USER         SSH user (for chown on remote kubeconfig)
#   KUBE_OUT           absolute path for local kubeconfig
#   EXPECTED           expected Ready node count
#   READY_TIMEOUT      seconds
#   LABEL_WORKERS      true|false
#   WORKER_NAMES       space-separated node names (may be empty)
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
FIRST="${FIRST_HOST:?FIRST_HOST is required}"
USER_NAME="${FIRST_USER:?FIRST_USER is required}"
KUBE_OUT="${KUBE_OUT:?KUBE_OUT is required}"
EXPECTED="${EXPECTED:?EXPECTED is required}"
READY_TIMEOUT="${READY_TIMEOUT:?READY_TIMEOUT is required}"
LABEL_WORKERS="${LABEL_WORKERS:-true}"
WORKER_NAMES="${WORKER_NAMES:-}"

kubectl_remote() {
  "$ROOT/scripts/ssh_run.sh" "$FIRST" \
    "sudo env KUBECONFIG=/etc/rancher/rke2/rke2.yaml /var/lib/rancher/rke2/bin/kubectl $*"
}

# Manager-node kubeconfig layout (matches the DataEngine guide).
"$ROOT/scripts/ssh_run.sh" "$FIRST" \
  "mkdir -p \$HOME/.kubeconfig && sudo cp /etc/rancher/rke2/rke2.yaml \$HOME/.kubeconfig/rke2.yaml && sudo chown ${USER_NAME}:${USER_NAME} \$HOME/.kubeconfig/rke2.yaml"
"$ROOT/scripts/ssh_run.sh" "$FIRST" \
  "grep -q 'KUBECONFIG' \$HOME/.profile || echo 'export KUBECONFIG=\$HOME/.kubeconfig/rke2.yaml' >> \$HOME/.profile"
"$ROOT/scripts/ssh_run.sh" "$FIRST" \
  "grep -q \"alias k=\" \$HOME/.profile || echo \"alias k='kubectl'\" >> \$HOME/.profile"

OUT_DIR="$(dirname "$KUBE_OUT")"
mkdir -p "${OUT_DIR}"
OUT_DIR="$(cd "${OUT_DIR}" && pwd)"
"$ROOT/scripts/ssh_run.sh" "$FIRST" "sudo cat /etc/rancher/rke2/rke2.yaml" >"$KUBE_OUT"
# Point local kubeconfig at the first control-plane management IP.
sed -i.bak -E \
  "s#server: https://127\\.0\\.0\\.1:6443#server: https://${FIRST}:6443#; s#server: https://localhost:6443#server: https://${FIRST}:6443#" \
  "$KUBE_OUT"
rm -f "${KUBE_OUT}.bak"
chmod 600 "$KUBE_OUT"

# VMS “Link Kubernetes cluster” wizard needs these base64 fields from kubeconfig
# (paste as-is into CA Certificate / Client Certificate / Private Key).
awk '/certificate-authority-data:/ {printf "%s", $2; exit}' "$KUBE_OUT" >"${OUT_DIR}/ca_cert"
awk '/client-certificate-data:/ {printf "%s", $2; exit}' "$KUBE_OUT" >"${OUT_DIR}/client_cert"
awk '/client-key-data:/ {printf "%s", $2; exit}' "$KUBE_OUT" >"${OUT_DIR}/client_key"
chmod 600 "${OUT_DIR}/client_key"
chmod 644 "${OUT_DIR}/ca_cert" "${OUT_DIR}/client_cert"
cat >"${OUT_DIR}/README.txt" <<EOF
VMS — Link Kubernetes cluster
  Name:               (e.g. de-lab-rke2)
  URL API server:     https://${FIRST}:6443
  CA Certificate:     contents of ca_cert
  Client Certificate: contents of client_cert
  Private Key:        contents of client_key

Files in this directory:
  kubeconfig     full kubectl config
  ca_cert        base64 certificate-authority-data
  client_cert    base64 client-certificate-data
  client_key     base64 client-key-data (secret)
EOF
echo "Wrote VMS link artifacts under ${OUT_DIR}/ (kubeconfig ca_cert client_cert client_key)"

DEADLINE=$((SECONDS + READY_TIMEOUT))
while true; do
  READY="$(kubectl_remote get nodes --no-headers 2>/dev/null | grep -c ' Ready' || true)"
  echo "cluster Ready count: ${READY}/${EXPECTED}"
  if [ "${READY}" -ge "${EXPECTED}" ]; then
    break
  fi
  if [ "$SECONDS" -ge "$DEADLINE" ]; then
    echo "Timed out waiting for all nodes Ready" >&2
    kubectl_remote get nodes -o wide || true
    exit 1
  fi
  sleep 10
done

if [ "$LABEL_WORKERS" = "true" ]; then
  for NODE in $WORKER_NAMES; do
    echo "Labeling worker ${NODE}"
    kubectl_remote label node "${NODE}" "node-role.kubernetes.io/worker=" --overwrite
  done
fi

kubectl_remote get nodes -o wide
