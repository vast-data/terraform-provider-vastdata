data "vastdata_compute_cluster_events" "calico_controller_events" {
  compute_cluster_name = "tf-test-cluster"
  resource_namespace   = "calico-system"
  resource_name        = "calico-kube-controllers"
  resource_kind        = "Deployment"
}
