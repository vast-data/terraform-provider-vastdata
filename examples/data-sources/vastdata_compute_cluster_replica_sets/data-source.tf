data "vastdata_compute_cluster_replica_sets" "calico_controller_replica_sets" {
  compute_cluster_name = "tf-test-cluster"
  resource_namespace   = "calico-system"
  resource_name        = "calico-kube-controllers"
}
