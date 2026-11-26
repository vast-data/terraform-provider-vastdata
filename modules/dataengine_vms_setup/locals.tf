locals {
  tenant_name      = coalesce(var.tenant_name, var.name_prefix)
  vip_pool_name    = coalesce(var.vip_pool_name, "${var.name_prefix}-vippool")
  view_policy_name = coalesce(var.view_policy_name, "${var.name_prefix}-s3-policy")
  broker_bucket    = coalesce(var.broker_bucket, "${var.name_prefix}-broker")
  api_token_name   = coalesce(var.api_token_name, "${var.name_prefix}-vastde-admin")

  is_local_identity = var.identity_provider_type == "local"
  is_ldap           = var.identity_provider_type == "ldap"
  is_ad             = var.identity_provider_type == "active_directory"

  local_provider_name = coalesce(var.local_provider_name, "${var.name_prefix}-provider")
  local_provider_id = local.is_local_identity ? (
    var.create_local_provider ? vastdata_local_provider.this[0].id : var.local_provider_id
  ) : null

  # Created local names vs existing LDAP/AD names
  app_group_name = local.is_local_identity ? coalesce(var.app_group_name, "${var.name_prefix}-app-group") : var.app_group_name
  broker_owner   = local.is_local_identity ? coalesce(var.broker_owner_name, "${var.name_prefix}-broker-owner") : var.broker_owner_name

  # External Kubernetes cluster (after DE enable)
  k8s_cluster_name        = coalesce(var.kubernetes_cluster_name, "${var.name_prefix}-k8s")
  k8s_mtls_name           = coalesce(var.kubernetes_mtls_name, "${var.name_prefix}-mtls")
  container_registry_name = coalesce(var.container_registry_name, "${var.name_prefix}-registry")

  # Prefer inline kubeconfig, else path (e.g. rke2/out/kubeconfig).
  k8s_kubeconfig_data = var.kubernetes_kubeconfig != null ? var.kubernetes_kubeconfig : (
    var.kubernetes_kubeconfig_path != null ? file(var.kubernetes_kubeconfig_path) : null
  )
  k8s_use_kubeconfig = try(length(trimspace(local.k8s_kubeconfig_data)) > 0, false)

  k8s_artifacts = var.kubernetes_mtls_artifacts_dir
  # PEM path only when not using kubeconfig (conditional avoids file() on unused branch).
  k8s_ca_b64 = local.k8s_use_kubeconfig ? null : (
    var.kubernetes_certificate_authority_b64 != null ? var.kubernetes_certificate_authority_b64 : (
      local.k8s_artifacts != null ? trimspace(file("${local.k8s_artifacts}/ca_cert")) : null
    )
  )
  k8s_client_cert_b64 = local.k8s_use_kubeconfig ? null : (
    var.kubernetes_client_certificate_b64 != null ? var.kubernetes_client_certificate_b64 : (
      local.k8s_artifacts != null ? trimspace(file("${local.k8s_artifacts}/client_cert")) : null
    )
  )
  k8s_client_key_b64 = local.k8s_use_kubeconfig ? null : (
    var.kubernetes_client_key_b64 != null ? var.kubernetes_client_key_b64 : (
      local.k8s_artifacts != null ? trimspace(file("${local.k8s_artifacts}/client_key")) : null
    )
  )
  k8s_has_pems = (
    try(length(trimspace(local.k8s_ca_b64)) > 0, false) &&
    try(length(trimspace(local.k8s_client_cert_b64)) > 0, false) &&
    try(length(trimspace(local.k8s_client_key_b64)) > 0, false)
  )

  # DataEngine tenant-admin
  tenant_admin_create_role       = var.create_tenant_admin && var.create_tenant_admin_role
  tenant_admin_use_existing_role = var.create_tenant_admin && !var.create_tenant_admin_role
  tenant_admin_role_name         = var.create_tenant_admin_role ? coalesce(var.tenant_admin_role_name, "${var.name_prefix}-tenant-admin") : var.tenant_admin_role_name
  tenant_admin_username          = coalesce(var.tenant_admin_username, "${var.name_prefix}-admin")
  tenant_admin_perms_list = setunion(
    var.tenant_admin_permissions_list,
    var.tenant_admin_extra_permissions_list
  )
  tenant_admin_use_list = length(local.tenant_admin_perms_list) > 0
  tenant_admin_role_id = var.create_tenant_admin ? (
    var.create_tenant_admin_role ? vastdata_administrator_role.tenant_admin[0].id : data.vastdata_administrator_role.tenant_admin[0].id
  ) : null

  # CSI naming
  csi_tenant_name               = coalesce(var.csi_tenant_name, "${var.name_prefix}-csi-tenant")
  csi_vip_pool_name             = coalesce(var.csi_vip_pool_name, "${var.name_prefix}-csi-vippool")
  csi_view_policy_name          = coalesce(var.csi_view_policy_name, "${var.name_prefix}-csi-policy")
  docker_repo_view_name         = coalesce(var.docker_repo_view_name, "${var.name_prefix}-docker-repo")
  csi_tenant_admin_role_name    = coalesce(var.csi_tenant_admin_role_name, "${var.name_prefix}-csi-tenant-admin")
  csi_service_account_role_name = coalesce(var.csi_service_account_role_name, "${var.name_prefix}-csi-service-account")
  csi_tenant_admin_username     = coalesce(var.csi_tenant_admin_username, "${var.name_prefix}-csi-admin")
  csi_driver_username           = coalesce(var.csi_driver_username, "${var.name_prefix}-csi-driver")

  csi_tenant_id = var.enable_csi ? (
    var.create_csi_tenant ? vastdata_tenant.csi[0].id : vastdata_tenant.dataengine.id
  ) : null

  csi_view_policy_auth_source     = try(var.csi_view_policy.auth_source, "RPC")
  csi_view_policy_bind_vip_pool   = try(var.csi_view_policy.bind_vip_pool, true)
  csi_view_policy_nfs_no_squash   = try(var.csi_view_policy.nfs_no_squash, toset(["*"]))
  csi_view_policy_nfs_root_squash = try(var.csi_view_policy.nfs_root_squash, toset([]))

  csi_tenant_admin_perms_list = setunion(
    var.csi_tenant_admin_permissions_list,
    var.csi_tenant_admin_extra_permissions_list
  )
  csi_tenant_admin_use_list = length(local.csi_tenant_admin_perms_list) > 0

  csi_service_account_perms_list = setunion(
    var.csi_service_account_permissions_list,
    var.csi_service_account_extra_permissions_list
  )
  csi_service_account_use_list = length(local.csi_service_account_perms_list) > 0
}
