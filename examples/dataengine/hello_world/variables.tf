variable "vms_host" {
  type = string
}

variable "vms_port" {
  type    = number
  default = 443
}

variable "vms_username" {
  type    = string
  default = null
}

variable "vms_password" {
  type      = string
  default   = null
  sensitive = true
}

variable "vms_api_token" {
  type      = string
  default   = null
  sensitive = true
}

variable "skip_ssl_verify" {
  type    = bool
  default = true
}

variable "tenant_name" {
  type        = string
  description = "Data Engine tenant name (query param for cluster-admin auth)."
}

variable "container_registry_vrn" {
  type        = string
  description = "VRN of vastdata_de_container_registry (from dataengine_vms_setup output)."
}

variable "kubernetes_cluster_vrn" {
  type        = string
  description = "VRN of vastdata_de_kubernetes_cluster."
}

variable "broker_name" {
  type        = string
  description = "Internal broker / view bucket name."
  default     = "de-lab-broker"
}

variable "topic_name" {
  type        = string
  description = "Topic name on the broker (engine main topic)."
  default     = "engine-main"
}

variable "artifact_source" {
  type        = string
  description = "Image path in the registry without host/tag (e.g. fnn45/hello-world)."
  default     = "de-lab-test/hello-world"
}

variable "image_tag" {
  type    = string
  default = "latest"
}

variable "function_name" {
  type    = string
  default = "de-test-function"
}

variable "schedule_trigger_name" {
  type        = string
  description = "Schedule trigger name."
  default     = "de-test-trigger"
}

variable "element_trigger_name" {
  type    = string
  default = "de-test-element-trigger"
}

variable "element_source_bucket" {
  type        = string
  description = "S3 view/bucket that emits object events for the element trigger."
  default     = null
  nullable    = true
}

variable "element_events" {
  type        = list(string)
  description = "TriggerEvent enum values (ObjectCreated:*, ObjectRemoved:*, …)."
  default     = ["ObjectCreated:*", "ObjectRemoved:*"]
}

variable "element_name_prefixes" {
  type    = list(string)
  default = ["hello/"]
}

variable "element_name_suffixes" {
  type    = list(string)
  default = [".json"]
}

variable "pipeline_name" {
  type    = string
  default = "de-test-pipeline"
}

variable "pipeline_namespace" {
  type    = string
  default = "vast-pipelines"
}

variable "pipeline_deploy" {
  type        = bool
  description = "POST /pipeline/{guid}/deploy after create/update and wait for Ready."
  default     = true
}

variable "cron_schedule" {
  type        = string
  description = "Quartz cron for the schedule trigger (every 1 minute by default)."
  default     = "0 * * ? * * *"
}

variable "skip_image_validation" {
  type        = bool
  description = "Skip DE image reachability check (needed for Docker Hub 401 on create)."
  default     = true
}

variable "secret_namespace" {
  type    = string
  default = "vast-pipelines"
}

variable "secret_name" {
  type    = string
  default = "de-hello-app-secret"
}

variable "secret_entries" {
  type = list(object({
    key   = string
    value = string
  }))
  description = "Key/value pairs written into the Kubernetes secret via DE API."
  sensitive   = true
  default = [
    { key = "DEMO_USER", value = "hello" },
    { key = "DEMO_TOKEN", value = "change-me" },
  ]
}
