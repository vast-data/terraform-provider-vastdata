# ignore:e2e
# Undeployed pipeline with per-function config, secrets, configmaps,
# event batching, and resource limits. Requires existing function / trigger / secret.

resource "vastdata_de_pipeline" "advanced_undeployed" {
  name                   = "de-pipe-advanced"
  tenant_name            = "de-lab"
  kubernetes_cluster_vrn = vastdata_de_kubernetes_cluster.example.vrn
  namespace              = "vast-pipelines"
  deploy                 = false
  description            = "Undeployed pipeline with advanced function deployment options"
  tags                   = ["example", "advanced"]

  depends_on = [vastdata_de_kubernetes_secret.example]

  manifest = {
    triggers = [
      {
        name = vastdata_de_schedule_trigger.hello.name
        vrn  = vastdata_de_schedule_trigger.hello.vrn
      },
    ]
    function_deployments = [
      {
        name         = "hello-fn"
        function_vrn = vastdata_de_function.hello.vrn
        config = {
          log_level = "DEBUG"
          secrets   = [vastdata_de_kubernetes_secret.example.secret_name]
          environment_variables = [
            { name = "PIPE_KIND", value = "advanced" },
          ]
          configmaps = [
            {
              name = "de-pipe-cm"
              data = { kind = "advanced" }
            },
          ]
          event_batching = {
            batch_max_size   = 5
            batch_timeout_ms = 1000
          }
        }
        resources = {
          min_cpu         = "100m"
          max_cpu         = "500m"
          min_memory      = "128Mi"
          max_memory      = "256Mi"
          timeout         = 60
          num_workers     = 4
          max_concurrency = 10
        }
      },
    ]
    links = [
      {
        source      = [vastdata_de_schedule_trigger.hello.name]
        destination = ["hello-fn"]
      },
    ]
    config = {
      secrets = [vastdata_de_kubernetes_secret.example.secret_name]
      environment_variables = [
        { name = "PIPE_SCOPE", value = "pipeline" },
      ]
    }
  }
}
