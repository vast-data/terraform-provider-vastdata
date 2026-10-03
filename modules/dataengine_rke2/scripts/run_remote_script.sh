#!/usr/bin/env bash
# Copy a local script to a remote host and run it with selected env vars forwarded.
#
# Required env (plus SSH_*):
#   REMOTE_HOST LOCAL_SCRIPT REMOTE_SCRIPT
# Optional passthrough: any VAR listed in FORWARD_VARS (default common registry vars).
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

REMOTE_HOST="${REMOTE_HOST:?}"
LOCAL_SCRIPT="${LOCAL_SCRIPT:?}"
REMOTE_SCRIPT="${REMOTE_SCRIPT:?}"
FORWARD_VARS="${FORWARD_VARS:-REGISTRY_HOST REGISTRY_PORT NFS_EXPORT MOUNTPOINT CERT_DAYS FORCE_RECERT RUN_SMOKE_TEST CONTAINER_NAME}"

"$ROOT/scripts/ssh_copy.sh" "$REMOTE_HOST" "$LOCAL_SCRIPT" "$REMOTE_SCRIPT"

ENV_PREFIX=""
for var in $FORWARD_VARS; do
  # Always forward listed vars (including empty NFS_EXPORT) so the remote script sees them.
  if [ -n "${var+x}" ]; then
    ENV_PREFIX="${ENV_PREFIX}${var}=$(printf '%q' "${!var-}") "
  fi
done

"$ROOT/scripts/ssh_run.sh" "$REMOTE_HOST" "chmod 755 ${REMOTE_SCRIPT} && ${ENV_PREFIX}${REMOTE_SCRIPT}"