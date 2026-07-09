data "vastdata_compute_cluster" "vastdb_compute_cluster" {
  name = "existing-compute-cluster"
}

resource "vastdata_compute_cluster_control" "vastdb_compute_cluster_control" {
  compute_cluster_id = data.vastdata_compute_cluster.vastdb_compute_cluster.id
  action             = "start"
}

# ---------------------
# Complete examples
# ---------------------

# ignore:e2e - requires an existing compute cluster and triggers lifecycle actions

data "vastdata_compute_cluster" "vastdb_compute_cluster" {
  name = "existing-compute-cluster"
}

resource "vastdata_compute_cluster_control" "vastdb_compute_cluster_control" {
  compute_cluster_id = data.vastdata_compute_cluster.vastdb_compute_cluster.id
  action             = "start"
}

# --------------------

