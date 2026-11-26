# Register a container registry with Data Engine and bind it to the Kubernetes cluster.

check "container_registry_requires_kubernetes_cluster" {
  assert {
    condition     = !var.enable_container_registry || var.enable_kubernetes_cluster
    error_message = "enable_container_registry requires enable_kubernetes_cluster=true."
  }
}

check "container_registry_requires_url" {
  assert {
    condition     = !var.enable_container_registry || try(length(trimspace(var.container_registry_url)) > 0, false)
    error_message = "enable_container_registry requires container_registry_url (e.g. https://<registry-host>:5000)."
  }
}

resource "vastdata_de_container_registry" "this" {
  count = var.enable_container_registry ? 1 : 0

  name                         = local.container_registry_name
  tenant_name                  = vastdata_tenant.dataengine.name
  url                          = var.container_registry_url
  auth_type                    = var.container_registry_auth_type
  username                     = var.container_registry_username
  password                     = var.container_registry_password
  is_accessible_to_data_engine = var.container_registry_is_accessible_to_data_engine

  primary_kubernetes_cluster = {
    kubernetes_cluster_vrn = vastdata_de_kubernetes_cluster.this[0].vrn
    namespace              = var.container_registry_primary_namespace
  }

  additional_kubernetes_clusters = [
    for ns in setsubtract(toset(vastdata_de_kubernetes_cluster.this[0].namespaces), toset([var.container_registry_primary_namespace])) : {
      kubernetes_cluster_vrn = vastdata_de_kubernetes_cluster.this[0].vrn
      namespace              = ns
    }
  ]

  depends_on = [vastdata_de_kubernetes_cluster.this]
}
