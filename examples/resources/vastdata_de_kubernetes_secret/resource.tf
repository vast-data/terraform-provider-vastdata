# Create (or upsert) a Kubernetes secret via Data Engine.
# There is no GET API — Terraform tracks identity from config and refreshes from state.
# Exactly one of kubernetes_cluster_vrn or kubernetes_cluster_name.

resource "vastdata_de_kubernetes_secret" "example" {
  tenant_name             = "de-lab" # required for cluster-admin auth
  kubernetes_cluster_name = vastdata_de_kubernetes_cluster.example.name
  namespace               = "vast-dataengine"
  secret_name             = "db-credentials"

  entries = [
    {
      key   = "username"
      value = "app-user"
    },
    {
      key   = "password"
      value = "******"
    },
  ]
}

# ---------------------
# Complete examples
# ---------------------

resource "vastdata_de_kubernetes_secret" "pipelines" {
  tenant_name            = "de-lab"
  kubernetes_cluster_vrn = vastdata_de_kubernetes_cluster.example.vrn
  namespace              = "vast-pipelines"
  secret_name            = "de-app-secret"

  entries = [
    {
      key   = "DEMO_USER"
      value = "hello"
    },
    {
      key   = "DEMO_TOKEN"
      value = "******"
    },
  ]
}
