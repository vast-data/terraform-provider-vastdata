output "identity_provider_type" {
  value = module.dataengine.identity_provider_type
}

output "local_provider_id" {
  value = module.dataengine.local_provider_id
}

output "local_provider_created" {
  value = module.dataengine.local_provider_created
}

output "local_provider_name" {
  value = module.dataengine.local_provider_name
}

output "ldap_provider_id" {
  value = module.dataengine.ldap_provider_id
}

output "ad_provider_id" {
  value = module.dataengine.ad_provider_id
}

output "tenant_id" {
  value = module.dataengine.tenant_id
}

output "tenant_guid" {
  value = module.dataengine.tenant_guid
}

output "tenant_name" {
  value = module.dataengine.tenant_name
}

output "tenant_admin_username" {
  value = module.dataengine.tenant_admin_username
}

output "tenant_admin_manager_id" {
  value = module.dataengine.tenant_admin_manager_id
}

output "data_engine_enabled" {
  description = "Tenant flag from VMS (refresh after vastdata_de_engine completes)."
  value       = module.dataengine.data_engine_enabled
}

output "de_engine_setup_status" {
  value = module.dataengine.de_engine_setup_status
}

output "de_engine_deletion_status" {
  value = module.dataengine.de_engine_deletion_status
}

output "kubernetes_cluster_enabled" {
  value = module.dataengine.kubernetes_cluster_enabled
}

output "kubernetes_mtls_guid" {
  value = module.dataengine.kubernetes_mtls_guid
}

output "kubernetes_cluster_guid" {
  value = module.dataengine.kubernetes_cluster_guid
}

output "kubernetes_cluster_vrn" {
  value = module.dataengine.kubernetes_cluster_vrn
}

output "kubernetes_cluster_name" {
  value = module.dataengine.kubernetes_cluster_name
}

output "container_registry_enabled" {
  value = module.dataengine.container_registry_enabled
}

output "container_registry_guid" {
  value = module.dataengine.container_registry_guid
}

output "container_registry_vrn" {
  value = module.dataengine.container_registry_vrn
}

output "container_registry_name" {
  value = module.dataengine.container_registry_name
}

output "app_group_name" {
  value = module.dataengine.app_group_name
}

output "app_group_gid" {
  value = module.dataengine.app_group_gid
}

output "broker_owner" {
  value = module.dataengine.broker_owner
}

output "vip_pool_id" {
  value = module.dataengine.vip_pool_id
}

output "broker_bucket" {
  value = module.dataengine.broker_bucket
}

output "main_topic_name" {
  value = module.dataengine.main_topic_name
}

output "deadletter_topic_name" {
  value = module.dataengine.deadletter_topic_name
}

output "api_token_id" {
  value = module.dataengine.api_token_id
}

output "csi_enabled" {
  value = module.dataengine.csi_enabled
}

output "csi_tenant_id" {
  value = module.dataengine.csi_tenant_id
}
