#!/usr/bin/env bash
# Copy a local file to a remote path. Args: <host> <local> <remote>
set -euo pipefail

HOST=${1:?host required}
LOCAL=${2:?local path required}
REMOTE=${3:?remote path required}

SSH_USER=${SSH_USER:?}
SSH_PORT=${SSH_PORT:-22}
SSH_STRICT=${SSH_STRICT_HOST_KEY_CHECKING:-accept-new}
SSH_KEY=${SSH_IDENTITY_FILE:?}

exec scp \
  -i "$SSH_KEY" \
  -P "$SSH_PORT" \
  -o "StrictHostKeyChecking=${SSH_STRICT}" \
  -o "UserKnownHostsFile=/dev/null" \
  -o "LogLevel=ERROR" \
  -o "BatchMode=yes" \
  "$LOCAL" \
  "${SSH_USER}@${HOST}:${REMOTE}"
