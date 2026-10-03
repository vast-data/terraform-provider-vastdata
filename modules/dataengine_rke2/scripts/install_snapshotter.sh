#!/usr/bin/env bash
# Apply kubernetes-csi external-snapshotter CRDs + controller (VAST CSI guide).
set -euo pipefail

: "${KUBECONFIG:?KUBECONFIG is required}"
: "${SNAPSHOTTER_VERSION:?SNAPSHOTTER_VERSION is required}"

BASE="https://raw.githubusercontent.com/kubernetes-csi/external-snapshotter/${SNAPSHOTTER_VERSION}"

kubectl apply -f "${BASE}/client/config/crd/snapshot.storage.k8s.io_volumesnapshotclasses.yaml"
kubectl apply -f "${BASE}/client/config/crd/snapshot.storage.k8s.io_volumesnapshotcontents.yaml"
kubectl apply -f "${BASE}/client/config/crd/snapshot.storage.k8s.io_volumesnapshots.yaml"
kubectl apply -f "${BASE}/deploy/kubernetes/snapshot-controller/rbac-snapshot-controller.yaml"
kubectl apply -f "${BASE}/deploy/kubernetes/snapshot-controller/setup-snapshot-controller.yaml"

echo "external-snapshotter ${SNAPSHOTTER_VERSION} applied."
