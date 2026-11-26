#!/bin/sh
# Bootstrap RKE2 on this host. Mirrors the DataEngine lab install_server.sh flow
# without ClusterShell. Intended to be copied and run with sudo.
#
# Usage:
#   bootstrap_rke2.sh first
#   bootstrap_rke2.sh next  <join-host>
#   bootstrap_rke2.sh agent <join-host>
#
# Environment:
#   RKE2_CHANNEL          (default: v1.35)
#   RKE2_TOKEN            (default: RKE2CLUSTER)
#   RKE2_WRITE_KUBECONFIG (default: 0644)
#   RKE2_EXTRA_TLS_SANS   (optional, space-separated)

set -eu

TYPE=${1:-}
SERVER=${2:-}
VERSION=${RKE2_CHANNEL:-v1.35}
TOKEN=${RKE2_TOKEN:-RKE2CLUSTER}
WRITE_MODE=${RKE2_WRITE_KUBECONFIG:-0644}
TMP_CONFIG=/tmp/rke2-config.yaml

help() {
  echo "bootstrap_rke2.sh [first|next|agent] [server]"
  echo "  first - Install the first server"
  echo "  next  - Install an additional server"
  echo "  agent - Install an agent"
  echo "  server - registration IP/hostname for next/agent (backend/storage IP preferred)"
}

add_ips() {
  for i in $(ip -br a | awk '{ print $3 }' | grep -vE '^127|^f' | awk -F"/" '{ print $1 }'); do
    echo "  - \"$i\"" >>"$TMP_CONFIG"
  done
  for san in ${RKE2_EXTRA_TLS_SANS:-}; do
    echo "  - \"$san\"" >>"$TMP_CONFIG"
  done
}

create_config() {
  rm -f "$TMP_CONFIG"
  {
    echo "token: ${TOKEN}"
    if [ "$TYPE" != "agent" ]; then
      echo "write-kubeconfig-mode: \"${WRITE_MODE}\""
      echo "tls-san:"
    fi
  } >"$TMP_CONFIG"

  if [ "$TYPE" != "agent" ]; then
    add_ips
  fi

  if [ "$TYPE" != "first" ]; then
    echo "server: https://${SERVER}:9345" >>"$TMP_CONFIG"
  fi
}

# get.rke2.io's download() is a single curl/wget with no retries. GitHub release
# tarballs often die mid-transfer (curl exit 56). Patch retries in, then re-run
# the installer if it still fails.
fetch_installer() {
  attempt=1
  while [ "$attempt" -le 5 ]; do
    if curl -fL --retry 5 --retry-delay 3 --connect-timeout 20 --speed-limit 1024 --speed-time 60 \
      -o /tmp/install-rke2.sh https://get.rke2.io; then
      chmod +x /tmp/install-rke2.sh
      return 0
    fi
    echo "get.rke2.io download failed (attempt ${attempt}/5)" >&2
    attempt=$((attempt + 1))
    sleep 10
  done
  return 1
}

harden_installer_downloads() {
  # Single quotes so "$1" reaches sed as a literal. $ is end-of-line in the pattern.
  if curl --help all 2>/dev/null | grep -q retry-all-errors; then
    sed -i 's|curl -o "\$1" -fsSL "\$2"|curl -o "$1" -fL --retry 8 --retry-all-errors --retry-delay 5 --connect-timeout 20 --speed-limit 1024 --speed-time 60 "$2"|' /tmp/install-rke2.sh
  else
    sed -i 's|curl -o "\$1" -fsSL "\$2"|curl -o "$1" -fL --retry 8 --retry-delay 5 --connect-timeout 20 --speed-limit 1024 --speed-time 60 "$2"|' /tmp/install-rke2.sh
  fi
  sed -i 's|wget -qO "\$1" "\$2"|wget -O "$1" -q --tries=8 --timeout=60 --waitretry=5 --retry-connrefused "$2"|' /tmp/install-rke2.sh

  if ! grep -q 'speed-limit' /tmp/install-rke2.sh; then
    echo "warning: rke2 installer download() was not patched; relying on outer retries" >&2
  fi
}

run_installer() {
  if [ "$TYPE" = "agent" ]; then
    INSTALL_RKE2_CHANNEL="$VERSION" INSTALL_RKE2_TYPE="agent" /tmp/install-rke2.sh
  else
    INSTALL_RKE2_CHANNEL="$VERSION" /tmp/install-rke2.sh
  fi
}

install_rke2() {
  mkdir -p /etc/rancher/rke2/ /var/lib/rancher/rke2/server/manifests/
  cp "$TMP_CONFIG" /etc/rancher/rke2/config.yaml

  fetch_installer
  harden_installer_downloads

  attempt=1
  while [ "$attempt" -le 4 ]; do
    if run_installer; then
      break
    fi
    echo "RKE2 installer failed (attempt ${attempt}/4); retrying in 20s" >&2
    attempt=$((attempt + 1))
    sleep 20
  done
  if [ "$attempt" -gt 4 ]; then
    echo "RKE2 installer failed after 4 attempts" >&2
    exit 1
  fi

  if [ "$TYPE" = "agent" ]; then
    systemctl enable --now rke2-agent.service
  else
    systemctl enable --now rke2-server.service
    ln -sf /var/lib/rancher/rke2/bin/kubectl /usr/local/bin/kubectl
  fi
}

if [ -z "$TYPE" ]; then
  help
  exit 0
fi

case "$TYPE" in
first | next | agent) ;;
*)
  help
  exit 1
  ;;
esac

if [ "$TYPE" != "first" ] && [ -z "$SERVER" ]; then
  echo "Error: server (join address) is required for $TYPE"
  help
  exit 1
fi

echo "Installing RKE2 role=$TYPE channel=$VERSION${SERVER:+ join=$SERVER}"
create_config
install_rke2
echo "RKE2 $TYPE install finished on $(hostname)"
