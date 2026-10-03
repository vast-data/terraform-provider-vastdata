# Register a container registry with Data Engine and bind it to a Kubernetes cluster.
# Requires vastdata_de_engine, mTLS credentials, and a kubernetes cluster resource.

provider "vastdata" {
  host     = "vast-cluster"
  username = "admin"
  password = "******"
  # tenant = "de-lab"
}

resource "vastdata_de_mtls_authentication_credential" "example" {
  name            = "de-mtls"
  tenant_name     = "de-lab" # required for cluster-admin auth
  kubeconfig_data = file("${path.module}/kubeconfig")
  context_name    = "my-cluster-context"
}

resource "vastdata_de_kubernetes_cluster" "example" {
  name                  = "de-test"
  tenant_name           = "de-lab"
  kube_api_url          = "https://k8s-api.example.com:6443"
  mtls_credentials_guid = vastdata_de_mtls_authentication_credential.example.guid
  namespaces = [
    "default",
    "istio-system",
    "vast-dataengine",
    "vast-pipelines",
  ]
}

resource "vastdata_de_container_registry" "example" {
  name        = "de-registry"
  tenant_name = "de-lab"
  url         = "https://registry.example.com"
  auth_type   = "none"

  primary_kubernetes_cluster = {
    kubernetes_cluster_vrn = vastdata_de_kubernetes_cluster.example.vrn
    namespace              = "default"
  }

  additional_kubernetes_clusters = [
    for ns in setsubtract(toset(vastdata_de_kubernetes_cluster.example.namespaces), toset(["default"])) : {
      kubernetes_cluster_vrn = vastdata_de_kubernetes_cluster.example.vrn
      namespace              = ns
    }
  ]
}
