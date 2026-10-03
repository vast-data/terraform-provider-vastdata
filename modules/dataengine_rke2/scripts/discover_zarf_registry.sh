#!/usr/bin/env bash
# terraform external data source: discover Zarf registry NodePort endpoint.
# stdin: JSON query with kubeconfig, require ("true"|"false")
# stdout: JSON {present, host, port}
set -euo pipefail

QUERY="$(cat)"
KUBECONFIG="$(python3 -c 'import json,sys; print(json.loads(sys.argv[1]).get("kubeconfig") or "")' "${QUERY}")"
REQUIRE="$(python3 -c 'import json,sys; print(json.loads(sys.argv[1]).get("require") or "false")' "${QUERY}")"
export KUBECONFIG

if [[ -z "${KUBECONFIG}" || ! -f "${KUBECONFIG}" ]]; then
  python3 -c 'import json; print(json.dumps({"present":"false","host":"","port":""}))'
  exit 0
fi

if ! kubectl -n zarf get svc zarf-docker-registry >/dev/null 2>&1; then
  if [[ "${REQUIRE}" == "true" ]]; then
    echo "zarf-docker-registry missing" >&2
    exit 1
  fi
  python3 -c 'import json; print(json.dumps({"present":"false","host":"","port":""}))'
  exit 0
fi

PORT="$(kubectl -n zarf get svc zarf-docker-registry -o jsonpath='{.spec.ports[0].nodePort}')"
HOST="$(kubectl -n zarf get pod -l app=docker-registry -o jsonpath='{.items[0].status.hostIP}')"
python3 -c 'import json,sys; print(json.dumps({"present":"true","host":sys.argv[1],"port":sys.argv[2]}))' "${HOST}" "${PORT}"
