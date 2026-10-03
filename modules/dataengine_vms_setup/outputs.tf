output "identity_provider_type" {
  value = var.identity_provider_type
}

output "local_provider_id" {
  description = "Local provider ID when identity_provider_type=local; otherwise null."
  value       = local.local_provider_id
}

output "local_provider_created" {
  description = "True if this module created the local provider."
  value       = local.is_local_identity && var.create_local_provider
}

output "local_provider_name" {
  description = "Created provider name, or null when not creating one."
  value       = try(vastdata_local_provider.this[0].name, null)
}

output "ldap_provider_id" {
  value = local.is_ldap ? var.ldap_provider_id : null
}

output "ad_provider_id" {
  value = local.is_ad ? var.ad_provider_id : null
}

output "tenant_id" {
  description = "DataEngine tenant numeric ID."
  value       = vastdata_tenant.dataengine.id
}

output "tenant_guid" {
  description = "DataEngine tenant GUID (setup-provisioning)."
  value       = vastdata_tenant.dataengine.guid
}

output "tenant_name" {
  value = vastdata_tenant.dataengine.name
}

output "tenant_admin_username" {
  description = "Created DataEngine tenant-admin username when create_tenant_admin=true; otherwise null."
  value       = try(vastdata_administrator_manager.tenant_admin[0].username, null)
}

output "tenant_admin_role_name" {
  description = "Role name assigned to the tenant-admin manager (created or looked up); null when create_tenant_admin=false."
  value       = var.create_tenant_admin ? local.tenant_admin_role_name : null
}

output "tenant_admin_role_id" {
  description = "Role ID assigned to the tenant-admin manager (created or looked up); null when create_tenant_admin=false."
  value       = local.tenant_admin_role_id
}

output "tenant_admin_role_created" {
  description = "True if this module created the tenant-admin role."
  value       = local.tenant_admin_create_role
}

output "tenant_admin_manager_id" {
  description = "Created tenant-admin manager ID when create_tenant_admin=true; otherwise null."
  value       = try(vastdata_administrator_manager.tenant_admin[0].id, null)
}

output "data_engine_enabled" {
  description = "Read-only VMS flag (false until POST /data-engine / vastdata_de_engine succeeds)."
  value       = vastdata_tenant.dataengine.data_engine_enabled
}

output "de_engine_setup_status" {
  description = "Setup status from vastdata_de_engine when enable_data_engine=true; otherwise null."
  value       = try(vastdata_de_engine.this[0].setup.status, null)
}

output "de_engine_deletion_status" {
  description = "Deletion status from vastdata_de_engine when enable_data_engine=true; otherwise null."
  value       = try(vastdata_de_engine.this[0].deletion.status, null)
}

output "kubernetes_cluster_enabled" {
  value = var.enable_kubernetes_cluster
}

output "kubernetes_mtls_guid" {
  description = "GUID of vastdata_de_mtls_authentication_credential when enable_kubernetes_cluster=true."
  value       = try(vastdata_de_mtls_authentication_credential.cluster[0].guid, null)
}

output "kubernetes_cluster_guid" {
  description = "GUID of vastdata_de_kubernetes_cluster when enable_kubernetes_cluster=true."
  value       = try(vastdata_de_kubernetes_cluster.this[0].guid, null)
}

output "kubernetes_cluster_vrn" {
  description = "VRN of vastdata_de_kubernetes_cluster (use for container registry binding)."
  value       = try(vastdata_de_kubernetes_cluster.this[0].vrn, null)
}

output "kubernetes_cluster_name" {
  value = var.enable_kubernetes_cluster ? local.k8s_cluster_name : null
}

output "container_registry_enabled" {
  value = var.enable_container_registry
}

output "container_registry_guid" {
  description = "GUID of vastdata_de_container_registry when enable_container_registry=true."
  value       = try(vastdata_de_container_registry.this[0].guid, null)
}

output "container_registry_vrn" {
  value = try(vastdata_de_container_registry.this[0].vrn, null)
}

output "container_registry_name" {
  value = var.enable_container_registry ? local.container_registry_name : null
}

output "app_group_id" {
  value = try(vastdata_group.app[0].id, null)
}

output "app_group_name" {
  value = local.app_group_name
}

output "app_group_gid" {
  value = try(vastdata_group.app[0].gid, var.app_group_gid)
}

output "broker_owner" {
  value = local.broker_owner
}

output "vip_pool_id" {
  value = vastdata_vip_pool.dataengine.id
}

output "vip_pool_name" {
  value = vastdata_vip_pool.dataengine.name
}

output "view_policy_id" {
  value = vastdata_view_policy.broker.id
}

output "broker_bucket" {
  description = "Kafka database_name / kafka_broker.name for setup-provisioning."
  value       = vastdata_view.broker.bucket
}

output "broker_view_id" {
  value = vastdata_view.broker.id
}

output "main_topic_name" {
  value = vastdata_topic.main.name
}

output "deadletter_topic_name" {
  value = vastdata_topic.deadletter.name
}

output "api_token_id" {
  description = "Created API token ID (token value is not exported; retrieve from VMS if needed)."
  value       = try(vastdata_api_token.vastde[0].id, null)
}

output "csi_enabled" {
  value = var.enable_csi
}

output "csi_tenant_id" {
  value = local.csi_tenant_id
}

output "csi_vip_pool_id" {
  value = try(vastdata_vip_pool.csi[0].id, null)
}

output "docker_repo_view_id" {
  value = try(vastdata_view.docker_repo[0].id, null)
}
