#!/usr/bin/env bash
# Install Rancher local-path provisioner (optional dataengine_rke2 storage backend).
set -euo pipefail

: "${KUBECONFIG:?KUBECONFIG is required}"
: "${LOCAL_PATH_VERSION:?LOCAL_PATH_VERSION is required}"
SET_DEFAULT="${SET_DEFAULT:-true}"

MANIFEST_URL="https://raw.githubusercontent.com/rancher/local-path-provisioner/${LOCAL_PATH_VERSION}/deploy/local-path-storage.yaml"

kubectl apply -f "${MANIFEST_URL}"
kubectl wait --for=condition=available --timeout=180s \
  deployment/local-path-provisioner -n local-path-storage

if [[ "${SET_DEFAULT}" == "true" ]]; then
  kubectl annotate storageclass local-path \
    storageclass.kubernetes.io/is-default-class=true --overwrite
else
  kubectl annotate storageclass local-path \
    storageclass.kubernetes.io/is-default-class- --overwrite 2>/dev/null || true
fi

kubectl get storageclass local-path
echo "local-path StorageClass is ready (default=${SET_DEFAULT})."
