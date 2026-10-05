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

module "dataengine" {
  source = "git::ssh://git@git.vastdata.com/dev/terraform-provider-vastdata.git//modules/dataengine_vms_setup?ref=85ba8e8f53fafcae17b4fafa7422348e06069c57"

  name_prefix = var.name_prefix

  identity_provider_type = var.identity_provider_type
  create_local_provider  = var.create_local_provider
  local_provider_name    = var.local_provider_name
  local_provider_id      = var.local_provider_id
  ldap_provider_id       = var.ldap_provider_id
  ad_provider_id         = var.ad_provider_id

  app_group_name        = var.app_group_name
  app_group_gid         = var.app_group_gid
  broker_owner_name     = var.broker_owner_name
  broker_owner_uid      = var.broker_owner_uid
  broker_owner_password = var.broker_owner_password

  vip_pool         = var.vip_pool
  broker_view_path = var.broker_view_path
  delete_dir       = var.delete_dir

  force_delete_tenant = var.force_delete_tenant

  enable_data_engine = var.enable_data_engine

  enable_kubernetes_cluster            = var.enable_kubernetes_cluster
  kubernetes_cluster_name              = var.kubernetes_cluster_name
  kubernetes_kube_api_url              = var.kubernetes_kube_api_url
  kubernetes_kubeconfig                = var.kubernetes_kubeconfig
  kubernetes_kubeconfig_path           = var.kubernetes_kubeconfig_path
  kubernetes_kubeconfig_context        = var.kubernetes_kubeconfig_context
  kubernetes_mtls_artifacts_dir        = var.kubernetes_mtls_artifacts_dir
  kubernetes_certificate_authority_b64 = var.kubernetes_certificate_authority_b64
  kubernetes_client_certificate_b64    = var.kubernetes_client_certificate_b64
  kubernetes_client_key_b64            = var.kubernetes_client_key_b64
  kubernetes_mtls_name                 = var.kubernetes_mtls_name
  kubernetes_namespaces                = var.kubernetes_namespaces
  kubernetes_is_internal               = var.kubernetes_is_internal

  enable_container_registry                       = var.enable_container_registry
  container_registry_name                         = var.container_registry_name
  container_registry_url                          = var.container_registry_url
  container_registry_auth_type                    = var.container_registry_auth_type
  container_registry_username                     = var.container_registry_username
  container_registry_password                     = var.container_registry_password
  container_registry_primary_namespace            = var.container_registry_primary_namespace
  container_registry_is_accessible_to_data_engine = var.container_registry_is_accessible_to_data_engine

  create_tenant_admin                 = var.create_tenant_admin
  tenant_admin_username               = var.tenant_admin_username
  tenant_admin_password               = var.tenant_admin_password
  create_tenant_admin_role            = var.create_tenant_admin_role
  tenant_admin_role_name              = var.tenant_admin_role_name
  tenant_admin_permissions_list       = var.tenant_admin_permissions_list
  tenant_admin_extra_permissions_list = var.tenant_admin_extra_permissions_list
  tenant_admin_permission_types       = var.tenant_admin_permission_types

  s3_view_policy_auth_source = var.s3_view_policy_auth_source

  main_topic_name       = var.main_topic_name
  deadletter_topic_name = var.deadletter_topic_name
  topic_partitions      = var.topic_partitions

  create_api_token = var.create_api_token
  api_token_name   = var.api_token_name
  api_token_owner  = var.api_token_owner
  api_token_expiry = var.api_token_expiry

  enable_csi           = var.enable_csi
  create_csi_tenant    = var.create_csi_tenant
  csi_vip_pool         = var.csi_vip_pool
  csi_view_policy      = var.csi_view_policy
  docker_repo_path     = var.docker_repo_path
  csi_manager_password = var.csi_manager_password
  trash_gid            = var.trash_gid

  csi_tenant_admin_permissions_list          = var.csi_tenant_admin_permissions_list
  csi_tenant_admin_extra_permissions_list    = var.csi_tenant_admin_extra_permissions_list
  csi_service_account_permissions_list       = var.csi_service_account_permissions_list
  csi_service_account_extra_permissions_list = var.csi_service_account_extra_permissions_list
}
