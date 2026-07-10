# ignore:e2e

# Reconcile compute cluster create state with the VAST backend.
# Prerequisites: cluster state is RUNNING.

data "vastdata_compute_cluster" "vastdb_compute_cluster" {
  name = "existing-compute-cluster"
}

resource "vastdata_compute_cluster_control" "vastdb_compute_cluster_control_reconcile_create" {
  compute_cluster_id = data.vastdata_compute_cluster.vastdb_compute_cluster.id
  action             = "reconcile_create"
}
