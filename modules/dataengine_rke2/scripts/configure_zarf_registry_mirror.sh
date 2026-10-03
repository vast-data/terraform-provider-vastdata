#!/usr/bin/env bash
# Merge Zarf in-cluster registry (HTTP NodePort) into RKE2 registries.yaml on
# every node, then optionally restart rke2-server/agent.
#
# Discovers NodePort + hostIP from the cluster unless ZARF_REGISTRY_HOST /
# ZARF_REGISTRY_PORT are set. Merges with any existing registries.yaml (e.g.
# Podman registry CA trust) on the Terraform runner, then pushes the result.
#
# Required env (plus SSH_*):
#   KUBECONFIG MASTER_IP RKE2_HOSTS CONTROL_HOSTS
#   WORKER_HOSTS (may be empty)
#   RESTART_RKE2 true|false
#   REQUIRE_ZARF true|false   — if false and zarf svc missing, exit 0 (no-op)
# Optional:
#   ZARF_REGISTRY_HOST ZARF_REGISTRY_PORT ZARF_REGISTRY_EXTRA_HOSTS
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

: "${KUBECONFIG:?}"
: "${MASTER_IP:?}"
: "${RKE2_HOSTS:?}"
: "${CONTROL_HOSTS:?}"
WORKER_HOSTS="${WORKER_HOSTS:-}"
RESTART_RKE2="${RESTART_RKE2:-true}"
REQUIRE_ZARF="${REQUIRE_ZARF:-true}"

if ! python3 -c 'import yaml' 2>/dev/null; then
  echo "ERROR: python3-yaml required on the Terraform runner (pip/apt install PyYAML)" >&2
  exit 1
fi

if ! kubectl -n zarf get svc zarf-docker-registry >/dev/null 2>&1; then
  if [[ "${REQUIRE_ZARF}" == "true" ]]; then
    echo "ERROR: zarf/zarf-docker-registry not found. Run zarf init, then re-apply." >&2
    exit 1
  fi
  echo "Zarf registry not present; skipping zarf HTTP mirror (REQUIRE_ZARF=false)."
  exit 0
fi

PORT="${ZARF_REGISTRY_PORT:-$(kubectl -n zarf get svc zarf-docker-registry -o jsonpath='{.spec.ports[0].nodePort}')}"
HOST="${ZARF_REGISTRY_HOST:-$(kubectl -n zarf get pod -l app=docker-registry -o jsonpath='{.items[0].status.hostIP}')}"

if [[ -z "${PORT}" || -z "${HOST}" ]]; then
  echo "ERROR: could not resolve Zarf registry host/NodePort" >&2
  exit 1
fi

# 127.0.0.1 (Zarf default) + registry node IP + optional extras (e.g. all node IPs).
HOSTS_CSV="${HOST},127.0.0.1"
if [[ -n "${ZARF_REGISTRY_EXTRA_HOSTS:-}" ]]; then
  HOSTS_CSV="${HOSTS_CSV},${ZARF_REGISTRY_EXTRA_HOSTS}"
fi

WORKDIR="$(mktemp -d)"
trap 'rm -rf "${WORKDIR}"' EXIT

python3 - "${HOSTS_CSV}" "${PORT}" "${WORKDIR}/fragment.yaml" <<'PY'
import sys, yaml
hosts_csv, port, out = sys.argv[1], sys.argv[2], sys.argv[3]
seen = []
for h in hosts_csv.replace(" ", ",").split(","):
    h = h.strip()
    if h and h not in seen:
        seen.append(h)
doc = {"mirrors": {}, "configs": {}}
for h in seen:
    key = f"{h}:{port}"
    doc["mirrors"][key] = {"endpoint": [f"http://{key}"]}
    doc["configs"][key] = {"tls": {"insecure_skip_verify": True}}
with open(out, "w") as f:
    yaml.safe_dump(doc, f, default_flow_style=False, sort_keys=False)
print("Zarf HTTP mirrors:", ", ".join(f"{h}:{port}" for h in seen))
PY

restart_group() {
  local svc=$1
  shift
  for host in "$@"; do
    echo "Restarting ${svc} on ${host}"
    "$ROOT/scripts/ssh_run.sh" "$host" "sudo systemctl restart ${svc}"
  done
}

for host in $RKE2_HOSTS; do
  echo "Merging zarf HTTP mirrors on ${host}"
  EXISTING="${WORKDIR}/existing-${host}.yaml"
  MERGED="${WORKDIR}/merged-${host}.yaml"
  # Pull current registries.yaml if any (empty file if missing).
  if "$ROOT/scripts/ssh_run.sh" "$host" "sudo test -f /etc/rancher/rke2/registries.yaml"; then
    "$ROOT/scripts/ssh_run.sh" "$host" "sudo cat /etc/rancher/rke2/registries.yaml" >"${EXISTING}"
  else
    : >"${EXISTING}"
  fi
  python3 - "${EXISTING}" "${WORKDIR}/fragment.yaml" "${MERGED}" <<'PY'
import sys, yaml, os
existing_path, frag_path, out_path = sys.argv[1], sys.argv[2], sys.argv[3]
base = {}
if os.path.getsize(existing_path) > 0:
    with open(existing_path) as f:
        base = yaml.safe_load(f) or {}
with open(frag_path) as f:
    frag = yaml.safe_load(f) or {}
base.setdefault("mirrors", {})
base.setdefault("configs", {})
base["mirrors"].update(frag.get("mirrors") or {})
base["configs"].update(frag.get("configs") or {})
with open(out_path, "w") as f:
    yaml.safe_dump(base, f, default_flow_style=False, sort_keys=False)
PY
  "$ROOT/scripts/ssh_copy.sh" "$host" "${MERGED}" /tmp/registries.yaml
  "$ROOT/scripts/ssh_run.sh" "$host" "
    set -euo pipefail
    sudo mkdir -p /etc/rancher/rke2
    if [ -f /etc/rancher/rke2/registries.yaml ]; then
      sudo cp /etc/rancher/rke2/registries.yaml /etc/rancher/rke2/registries.yaml.pre-zarf-http
    fi
    sudo install -m 0644 /tmp/registries.yaml /etc/rancher/rke2/registries.yaml
    rm -f /tmp/registries.yaml
    echo \"\$(hostname): registries.yaml:\"
    sudo awk '1' /etc/rancher/rke2/registries.yaml
  "
done

if [ "${RESTART_RKE2}" = "true" ]; then
  # shellcheck disable=SC2086
  restart_group rke2-server ${CONTROL_HOSTS}
  if [ -n "${WORKER_HOSTS}" ]; then
    # shellcheck disable=SC2086
    restart_group rke2-agent ${WORKER_HOSTS}
  fi
  echo "Waiting for API after restart..."
  for _ in $(seq 1 90); do
    if "$ROOT/scripts/ssh_run.sh" "$MASTER_IP" \
      "sudo /var/lib/rancher/rke2/bin/kubectl --kubeconfig /etc/rancher/rke2/rke2.yaml get nodes >/dev/null 2>&1"; then
      break
    fi
    sleep 10
  done
fi

echo "Zarf registry HTTP mirror configuration complete"
