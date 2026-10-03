# ignore:e2e
# Element-trigger-only graph (no schedule). Deploy left false for docs/lab use.

resource "vastdata_de_pipeline" "element_only" {
  name                   = "de-pipe-element"
  tenant_name            = "de-lab"
  kubernetes_cluster_vrn = vastdata_de_kubernetes_cluster.example.vrn
  namespace              = "vast-pipelines"
  deploy                 = false
  description            = "Pipeline driven only by an element (S3) trigger"
  tags                   = ["example", "element"]

  manifest = {
    triggers = [
      {
        name = vastdata_de_element_trigger.example.name
        vrn  = vastdata_de_element_trigger.example.vrn
      },
    ]
    function_deployments = [
      {
        name         = "hello-fn"
        function_vrn = vastdata_de_function.hello.vrn
      },
    ]
    links = [
      {
        source      = [vastdata_de_element_trigger.example.name]
        destination = ["hello-fn"]
      },
    ]
  }
}
