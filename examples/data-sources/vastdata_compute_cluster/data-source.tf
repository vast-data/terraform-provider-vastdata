data "vastdata_compute_cluster" "vastdb_compute_cluster_by_id" {
  id        = 1
  get_nodes = true
}

data "vastdata_compute_cluster" "vastdb_compute_cluster_by_name" {
  name           = "existing-compute-cluster"
  get_pods       = true
  get_namespaces = true
  get_services   = true
}
