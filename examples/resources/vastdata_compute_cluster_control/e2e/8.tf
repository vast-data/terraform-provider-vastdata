# ignore:e2e

# Rotate base certificates including root (keep_root = false).
# Prerequisites: cluster state is RUNNING.
# VMS regenerates the full certificate chain when custom PEMs are not supplied.

data "vastdata_compute_cluster" "vastdb_compute_cluster" {
  name = "existing-compute-cluster"
}

resource "vastdata_compute_cluster_control" "vastdb_compute_cluster_control_rotate_base_full" {
  compute_cluster_id = data.vastdata_compute_cluster.vastdb_compute_cluster.id
  action             = "rotate_base_certificates"
  keep_root          = false
}
