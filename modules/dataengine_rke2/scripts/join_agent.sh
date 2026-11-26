#!/usr/bin/env bash
# Install RKE2 agents one node at a time. Parallel GitHub release downloads
# from every worker at once tend to stall (curl exit 56).
#
# Required env:
#   SSH_* / RKE2_*
#   JOIN_ADDRESS   registration host (storage IP preferred)
#   AGENT_HOSTS    space-separated worker management IPs
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
JOIN="${JOIN_ADDRESS:?JOIN_ADDRESS is required}"
AGENT_HOSTS="${AGENT_HOSTS:?AGENT_HOSTS is required}"

for HOST in $AGENT_HOSTS; do
  echo "Installing RKE2 agent on ${HOST} via ${JOIN}"
  "$ROOT/scripts/ssh_copy.sh" "$HOST" "$ROOT/scripts/bootstrap_rke2.sh" /tmp/bootstrap_rke2.sh
  "$ROOT/scripts/ssh_run.sh" "$HOST" \
    "chmod 755 /tmp/bootstrap_rke2.sh && sudo env RKE2_CHANNEL='$RKE2_CHANNEL' RKE2_TOKEN='$RKE2_TOKEN' RKE2_WRITE_KUBECONFIG='$RKE2_WRITE_KUBECONFIG' RKE2_EXTRA_TLS_SANS='$RKE2_EXTRA_TLS_SANS' /tmp/bootstrap_rke2.sh agent $JOIN"
done
