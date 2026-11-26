#!/usr/bin/env bash
# Run a remote command over SSH. Args: <host> <remote-command...>
set -euo pipefail

HOST=${1:?host required}
shift

SSH_USER=${SSH_USER:?}
SSH_PORT=${SSH_PORT:-22}
SSH_STRICT=${SSH_STRICT_HOST_KEY_CHECKING:-accept-new}
SSH_KEY=${SSH_IDENTITY_FILE:?}

exec ssh \
  -i "$SSH_KEY" \
  -p "$SSH_PORT" \
  -o "StrictHostKeyChecking=${SSH_STRICT}" \
  -o "UserKnownHostsFile=/dev/null" \
  -o "LogLevel=ERROR" \
  -o "BatchMode=yes" \
  -o "ConnectTimeout=30" \
  "${SSH_USER}@${HOST}" \
  "$@"
