# ignore:e2e
# Secret in the pipelines namespace (typical app secret mount target).

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
