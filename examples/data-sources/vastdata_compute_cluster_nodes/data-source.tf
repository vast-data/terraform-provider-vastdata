data "vastdata_compute_cluster_nodes" "tf_test_cluster_nodes_by_name" {
  compute_cluster_name = "tf-test-cluster"
}

data "vastdata_compute_cluster_nodes" "tf_test_cluster_nodes_by_id" {
  compute_cluster_id = 2
}
