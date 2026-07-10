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


# Start a stopped compute cluster.
# Prerequisites: cluster state is STOPPED.

data "vastdata_compute_cluster" "vastdb_compute_cluster" {
  name = "existing-compute-cluster"
}

resource "vastdata_compute_cluster_control" "vastdb_compute_cluster_control_start" {
  compute_cluster_id = data.vastdata_compute_cluster.vastdb_compute_cluster.id
  action             = "start"
}

# --------------------


# Stop a running compute cluster.
# Prerequisites: cluster state is RUNNING.

data "vastdata_compute_cluster" "vastdb_compute_cluster" {
  name = "existing-compute-cluster"
}

resource "vastdata_compute_cluster_control" "vastdb_compute_cluster_control_stop" {
  compute_cluster_id = data.vastdata_compute_cluster.vastdb_compute_cluster.id
  action             = "stop"
}

# --------------------


# Rotate leaf certificates on a running compute cluster.
# Prerequisites: cluster state is RUNNING.

data "vastdata_compute_cluster" "vastdb_compute_cluster" {
  name = "existing-compute-cluster"
}

resource "vastdata_compute_cluster_control" "vastdb_compute_cluster_control_rotate_leaf" {
  compute_cluster_id = data.vastdata_compute_cluster.vastdb_compute_cluster.id
  action             = "rotate_leaf_certificates"
}

# --------------------


# Rotate the compute cluster service key.
# Prerequisites: cluster state is RUNNING.

data "vastdata_compute_cluster" "vastdb_compute_cluster" {
  name = "existing-compute-cluster"
}

resource "vastdata_compute_cluster_control" "vastdb_compute_cluster_control_rotate_service_key" {
  compute_cluster_id = data.vastdata_compute_cluster.vastdb_compute_cluster.id
  action             = "rotate_service_key"
}

# --------------------


# Reconcile compute cluster create state with the VAST backend.
# Prerequisites: cluster state is RUNNING.

data "vastdata_compute_cluster" "vastdb_compute_cluster" {
  name = "existing-compute-cluster"
}

resource "vastdata_compute_cluster_control" "vastdb_compute_cluster_control_reconcile_create" {
  compute_cluster_id = data.vastdata_compute_cluster.vastdb_compute_cluster.id
  action             = "reconcile_create"
}

# --------------------


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

# --------------------


# Generate metric viewer certificates for the compute cluster.
# Prerequisites: cluster state is RUNNING.

data "vastdata_compute_cluster" "vastdb_compute_cluster" {
  name = "existing-compute-cluster"
}

resource "vastdata_compute_cluster_control" "vastdb_compute_cluster_control_metric_viewer" {
  compute_cluster_id = data.vastdata_compute_cluster.vastdb_compute_cluster.id
  action             = "metric_viewer_certificates"
}

# --------------------


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

# --------------------

