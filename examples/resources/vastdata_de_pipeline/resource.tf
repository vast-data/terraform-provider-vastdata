# Create a Data Engine pipeline linking a schedule trigger to a function,
# then deploy it (default) and wait until status is Ready.
# Exactly one of kubernetes_cluster_vrn or kubernetes_cluster_name.

resource "vastdata_de_pipeline" "hello" {
  name                    = "de-test-pipeline"
  tenant_name             = "de-lab"
  kubernetes_cluster_name = vastdata_de_kubernetes_cluster.example.name
  namespace               = "vast-pipelines"
  deploy                  = true

  manifest = {
    triggers = [
      {
        name = vastdata_de_schedule_trigger.hello.name
        vrn  = vastdata_de_schedule_trigger.hello.vrn
      }
    ]
    function_deployments = [
      {
        name         = "hello-fn"
        function_vrn = vastdata_de_function.hello.vrn
      }
    ]
    links = [
      {
        source      = [vastdata_de_schedule_trigger.hello.name]
        destination = ["hello-fn"]
      }
    ]
  }
}

# ---------------------
# Complete examples
# ---------------------

# Undeployed pipeline with per-function config, secrets, configmaps,
# event batching, and resource limits.
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

# --------------------

# Element-trigger-only graph (no schedule).
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
