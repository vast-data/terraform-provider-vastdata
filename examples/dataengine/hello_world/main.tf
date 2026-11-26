terraform {
  required_version = ">= 1.5.0"

  required_providers {
    vastdata = {
      source  = "vast-data/vastdata"
      version = ">= 2.0.0"
    }
  }
}

provider "vastdata" {
  host            = var.vms_host
  port            = var.vms_port
  username        = var.vms_username
  password        = var.vms_password
  api_token       = var.vms_api_token
  skip_ssl_verify = var.skip_ssl_verify
}

# --- Function -----------------------------------------------------------------

resource "vastdata_de_function" "hello" {
  name                   = var.function_name
  tenant_name            = var.tenant_name
  container_registry_vrn = var.container_registry_vrn
  artifact_source        = var.artifact_source
  image_tag              = var.image_tag
  revision_alias         = "latest"
  description            = "Hello World function"
  skip_image_validation  = var.skip_image_validation
}


resource "vastdata_de_schedule_trigger" "hello" {
  name        = var.schedule_trigger_name
  tenant_name = var.tenant_name
  topic_name  = var.topic_name
  description = "Hello World schedule trigger"

  broker = {
    name = var.broker_name
    type = "Internal"
  }

  config = {
    cron_schedule = var.cron_schedule
  }
}

resource "vastdata_de_element_trigger" "hello" {
  name         = var.element_trigger_name
  tenant_name  = var.tenant_name
  topic_name   = var.topic_name
  description  = "Hello World element (S3) trigger"
  source_types = ["S3"]

  broker = {
    name = var.broker_name
    type = "Internal"
  }

  config = {
    source_bucket_name = coalesce(var.element_source_bucket, var.broker_name)
    events             = var.element_events
    name_filters = {
      prefixes = var.element_name_prefixes
      suffixes = var.element_name_suffixes
    }
  }
}

# --- Kubernetes secret (via Data Engine API) ----------------------------------

resource "vastdata_de_kubernetes_secret" "hello" {
  tenant_name            = var.tenant_name
  kubernetes_cluster_vrn = var.kubernetes_cluster_vrn
  namespace              = var.secret_namespace
  secret_name            = var.secret_name

  entries = [
    for e in var.secret_entries : {
      key   = e.key
      value = e.value
    }
  ]
}

# --- Pipeline (schedule + element → function) ---------------------------------

resource "vastdata_de_pipeline" "hello" {
  name                   = var.pipeline_name
  tenant_name            = var.tenant_name
  kubernetes_cluster_vrn = var.kubernetes_cluster_vrn
  namespace              = var.pipeline_namespace
  deploy                 = var.pipeline_deploy
  description            = "Hello World pipeline (schedule + element triggers)"

  depends_on = [
    vastdata_de_kubernetes_secret.hello,
  ]

  manifest = {
    triggers = [
      {
        name = vastdata_de_schedule_trigger.hello.name
        vrn  = vastdata_de_schedule_trigger.hello.vrn
      },
      {
        name = vastdata_de_element_trigger.hello.name
        vrn  = vastdata_de_element_trigger.hello.vrn
      },
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
      },
      {
        source      = [vastdata_de_element_trigger.hello.name]
        destination = ["hello-fn"]
      },
    ]
  }
}
