# ignore:e2e

# Rotate base certificates while keeping the existing root CA (keep_root = true).
# Prerequisites: cluster state is RUNNING.

data "vastdata_compute_cluster" "vastdb_compute_cluster" {
  name = "existing-compute-cluster"
}

resource "vastdata_compute_cluster_control" "vastdb_compute_cluster_control_rotate_base_keep_root" {
  compute_cluster_id = data.vastdata_compute_cluster.vastdb_compute_cluster.id
  action             = "rotate_base_certificates"
  keep_root          = true
}
