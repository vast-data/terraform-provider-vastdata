#!/usr/bin/env bash
# Install TLS Podman registry on the registry VM (DataEngine guide).
# Runs on the registry host (copied via SSH and executed with sudo where needed).
#
# Required env:
#   REGISTRY_HOST REGISTRY_PORT MOUNTPOINT CERT_DAYS
# Optional:
#   NFS_EXPORT          if set, mount this NFS export at MOUNTPOINT (VAST lab)
#                       if empty, use MOUNTPOINT as a local directory (created here)
#   FORCE_RECERT=true|false  RUN_SMOKE_TEST=true|false
set -euo pipefail

REGISTRY_HOST="${REGISTRY_HOST:?}"
REGISTRY_PORT="${REGISTRY_PORT:?}"
MOUNTPOINT="${MOUNTPOINT:?}"
NFS_EXPORT="${NFS_EXPORT:-}"
CERT_DAYS="${CERT_DAYS:-3650}"
FORCE_RECERT="${FORCE_RECERT:-false}"
RUN_SMOKE_TEST="${RUN_SMOKE_TEST:-false}"
CONTAINER_NAME="${CONTAINER_NAME:-local-registry}"
UNIT_NAME="container-${CONTAINER_NAME}.service"

export DEBIAN_FRONTEND=noninteractive

sudo apt-get update -y
sudo apt-get install -y podman openssl curl
if [ -n "${NFS_EXPORT}" ]; then
  sudo apt-get install -y nfs-common
fi

sudo mkdir -p "${MOUNTPOINT}"
if [ -n "${NFS_EXPORT}" ]; then
  if ! mountpoint -q "${MOUNTPOINT}"; then
    sudo mount -t nfs "${NFS_EXPORT}" "${MOUNTPOINT}"
  fi
  if ! grep -qsE "[[:space:]]${MOUNTPOINT}[[:space:]]" /etc/fstab; then
    echo "${NFS_EXPORT} ${MOUNTPOINT} nfs defaults,_netdev 0 0" | sudo tee -a /etc/fstab >/dev/null
  fi
else
  echo "No NFS_EXPORT set — using local directory ${MOUNTPOINT}"
fi

sudo mkdir -p "${MOUNTPOINT}/certs"

CERT="${MOUNTPOINT}/certs/domain.crt"
KEY="${MOUNTPOINT}/certs/domain.key"
CNF="${MOUNTPOINT}/certs/openssl.cnf"

need_cert=false
if [ "${FORCE_RECERT}" = "true" ] || [ ! -f "${CERT}" ] || [ ! -f "${KEY}" ]; then
  need_cert=true
fi

if [ "${need_cert}" = "true" ]; then
  cat <<EOF | sudo tee "${CNF}" >/dev/null
[req]
default_bits       = 4096
prompt             = no
default_md         = sha256
distinguished_name = dn
x509_extensions    = v3_req
[dn]
C  = US
ST = Lab
L = Lab
O  = VAST Lab
CN = ${REGISTRY_HOST}
[v3_req]
subjectAltName = @alt_names
keyUsage = digitalSignature, keyEncipherment
extendedKeyUsage = serverAuth
[alt_names]
IP.1 = ${REGISTRY_HOST}
EOF
  sudo openssl req -x509 -nodes -days "${CERT_DAYS}" \
    -newkey rsa:4096 \
    -keyout "${KEY}" \
    -out "${CERT}" \
    -config "${CNF}"
fi

sudo chmod 600 "${KEY}"
sudo chmod 644 "${CERT}"
# registry:2 runs as uid 1000 — storage + certs must be writable/readable by that uid
sudo chown -R 1000:1000 "${MOUNTPOINT}"

sudo mkdir -p "/etc/containers/certs.d/${REGISTRY_HOST}:${REGISTRY_PORT}"
sudo cp "${CERT}" "/etc/containers/certs.d/${REGISTRY_HOST}:${REGISTRY_PORT}/ca.crt"

# Tear down any previous rootless or rootful instance.
systemctl --user disable --now "${UNIT_NAME}" 2>/dev/null || true
sudo systemctl disable --now "${UNIT_NAME}" 2>/dev/null || true
if podman container exists "${CONTAINER_NAME}" 2>/dev/null; then
  podman stop "${CONTAINER_NAME}" || true
  podman rm "${CONTAINER_NAME}" || true
fi
if sudo podman container exists "${CONTAINER_NAME}" 2>/dev/null; then
  sudo podman stop "${CONTAINER_NAME}" || true
  sudo podman rm "${CONTAINER_NAME}" || true
fi

# Rootful + host network: avoids Docker/Podman FORWARD DROP that breaks
# published ports between Multipass VMs (DNAT works locally, fails from peers).
sudo podman run -d \
  --name "${CONTAINER_NAME}" \
  --network host \
  -v "${MOUNTPOINT}:/var/lib/registry:Z" \
  -v "${MOUNTPOINT}/certs:/certs:Z" \
  -e REGISTRY_HTTP_ADDR=0.0.0.0:${REGISTRY_PORT} \
  -e REGISTRY_HTTP_TLS_CERTIFICATE=/certs/domain.crt \
  -e REGISTRY_HTTP_TLS_KEY=/certs/domain.key \
  docker.io/library/registry:2

TMP_UNIT="$(mktemp -d)"
(
  cd "${TMP_UNIT}"
  sudo podman generate systemd --name "${CONTAINER_NAME}" --files --new
)
sudo cp "${TMP_UNIT}/${UNIT_NAME}" /etc/systemd/system/
rm -rf "${TMP_UNIT}"
sudo systemctl daemon-reload
sudo systemctl enable --now "${UNIT_NAME}"

echo "Waiting for registry HTTPS..."
for i in $(seq 1 30); do
  if curl -sS --cacert "${CERT}" "https://${REGISTRY_HOST}:${REGISTRY_PORT}/v2/_catalog" >/dev/null 2>&1; then
    break
  fi
  sleep 2
done

echo "TLS without CA (expect failure):"
curl -sS "https://${REGISTRY_HOST}:${REGISTRY_PORT}/v2/_catalog" && exit 1 || true

echo "TLS with CA:"
curl -sS --cacert "${CERT}" "https://${REGISTRY_HOST}:${REGISTRY_PORT}/v2/_catalog"
echo

sudo openssl x509 -in "${CERT}" -noout -subject -dates -fingerprint -sha256 -ext subjectAltName

if [ "${RUN_SMOKE_TEST}" = "true" ]; then
  # Client can be rootless; registry is rootful and writable.
  podman pull docker.io/library/alpine:latest
  podman tag alpine:latest "${REGISTRY_HOST}:${REGISTRY_PORT}/test-alpine:v1"
  podman push "${REGISTRY_HOST}:${REGISTRY_PORT}/test-alpine:v1"
  curl -sS --cacert "/etc/containers/certs.d/${REGISTRY_HOST}:${REGISTRY_PORT}/ca.crt" \
    "https://${REGISTRY_HOST}:${REGISTRY_PORT}/v2/_catalog"
  echo
fi

echo "Podman registry ready at https://${REGISTRY_HOST}:${REGISTRY_PORT}"
