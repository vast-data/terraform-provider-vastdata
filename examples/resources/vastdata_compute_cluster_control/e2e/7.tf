# ignore:e2e

# Generate metric viewer certificates for the compute cluster.
# Prerequisites: cluster state is RUNNING.

data "vastdata_compute_cluster" "vastdb_compute_cluster" {
  name = "existing-compute-cluster"
}

resource "vastdata_compute_cluster_control" "vastdb_compute_cluster_control_metric_viewer" {
  compute_cluster_id = data.vastdata_compute_cluster.vastdb_compute_cluster.id
  action             = "metric_viewer_certificates"
}
