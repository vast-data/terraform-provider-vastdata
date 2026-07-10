# ignore:e2e

# Start a stopped compute cluster.
# Prerequisites: cluster state is STOPPED.

data "vastdata_compute_cluster" "vastdb_compute_cluster" {
  name = "existing-compute-cluster"
}

resource "vastdata_compute_cluster_control" "vastdb_compute_cluster_control_start" {
  compute_cluster_id = data.vastdata_compute_cluster.vastdb_compute_cluster.id
  action             = "start"
}
