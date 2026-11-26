output "function_guid" {
  value = vastdata_de_function.hello.guid
}

output "function_vrn" {
  value = vastdata_de_function.hello.vrn
}

output "schedule_trigger_guid" {
  value = vastdata_de_schedule_trigger.hello.guid
}

output "schedule_trigger_vrn" {
  value = vastdata_de_schedule_trigger.hello.vrn
}

output "element_trigger_guid" {
  value = vastdata_de_element_trigger.hello.guid
}

output "element_trigger_vrn" {
  value = vastdata_de_element_trigger.hello.vrn
}

output "kubernetes_secret_name" {
  value = vastdata_de_kubernetes_secret.hello.secret_name
}

output "kubernetes_secret_namespace" {
  value = vastdata_de_kubernetes_secret.hello.namespace
}

output "pipeline_guid" {
  value = vastdata_de_pipeline.hello.guid
}

output "pipeline_vrn" {
  value = vastdata_de_pipeline.hello.vrn
}

output "pipeline_status" {
  value = vastdata_de_pipeline.hello.status
}

# Back-compat aliases
output "trigger_guid" {
  value = vastdata_de_schedule_trigger.hello.guid
}

output "trigger_vrn" {
  value = vastdata_de_schedule_trigger.hello.vrn
}
