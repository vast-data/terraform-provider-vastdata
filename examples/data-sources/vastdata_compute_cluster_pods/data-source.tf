data "vastdata_compute_cluster_pods" "tf_test_cluster_pods" {
  compute_cluster_name = "tf-test-cluster"
}

data "vastdata_compute_cluster_pods" "calico_system_pods" {
  compute_cluster_name = "tf-test-cluster"
  namespace            = "calico-system"
}
