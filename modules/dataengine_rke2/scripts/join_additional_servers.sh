#!/usr/bin/env bash
# Join additional RKE2 control-plane nodes sequentially (etcd-safe).
#
# Required env:
#   SSH_* / RKE2_*
#   JOIN_ADDRESS       registration host (storage IP preferred)
#   FIRST_HOST         first control-plane management IP
#   ADDITIONAL_HOSTS   space-separated management IPs to join
#   EXPECTED_CPS       expected Ready control-plane count after join
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
JOIN="${JOIN_ADDRESS:?JOIN_ADDRESS is required}"
FIRST="${FIRST_HOST:?FIRST_HOST is required}"
EXPECTED="${EXPECTED_CPS:?EXPECTED_CPS is required}"
ADDITIONAL_HOSTS="${ADDITIONAL_HOSTS:?ADDITIONAL_HOSTS is required}"

for HOST in $ADDITIONAL_HOSTS; do
  echo "Joining additional control-plane ${HOST} via ${JOIN}"
  "$ROOT/scripts/ssh_copy.sh" "$HOST" "$ROOT/scripts/bootstrap_rke2.sh" /tmp/bootstrap_rke2.sh
  "$ROOT/scripts/ssh_run.sh" "$HOST" \
    "chmod 755 /tmp/bootstrap_rke2.sh && sudo env RKE2_CHANNEL='$RKE2_CHANNEL' RKE2_TOKEN='$RKE2_TOKEN' RKE2_WRITE_KUBECONFIG='$RKE2_WRITE_KUBECONFIG' RKE2_EXTRA_TLS_SANS='$RKE2_EXTRA_TLS_SANS' /tmp/bootstrap_rke2.sh next $JOIN"
  # Give etcd a moment before the next member joins.
  sleep 30
done

for i in $(seq 1 90); do
  READY="$("$ROOT/scripts/ssh_run.sh" "$FIRST" \
    "sudo env KUBECONFIG=/etc/rancher/rke2/rke2.yaml /var/lib/rancher/rke2/bin/kubectl get nodes --no-headers 2>/dev/null | grep -c ' Ready' || true")"
  echo "control-plane Ready count: ${READY}/${EXPECTED}"
  if [ "${READY}" -ge "${EXPECTED}" ]; then
    exit 0
  fi
  sleep 10
done

echo "Timed out waiting for control-plane nodes" >&2
exit 1
