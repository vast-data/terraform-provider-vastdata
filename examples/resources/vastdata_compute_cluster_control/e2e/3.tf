# ignore:e2e

# Rotate leaf certificates on a running compute cluster.
# Prerequisites: cluster state is RUNNING.

data "vastdata_compute_cluster" "vastdb_compute_cluster" {
  name = "existing-compute-cluster"
}

resource "vastdata_compute_cluster_control" "vastdb_compute_cluster_control_rotate_leaf" {
  compute_cluster_id = data.vastdata_compute_cluster.vastdb_compute_cluster.id
  action             = "rotate_leaf_certificates"
}
