# Optional CSI helpers (lab guide). Enable with enable_csi=true.

check "csi_vip_pool_required" {
  assert {
    condition     = !var.enable_csi || var.csi_vip_pool != null
    error_message = "csi_vip_pool must be set when enable_csi=true."
  }
}

check "csi_password_required" {
  assert {
    condition     = !var.enable_csi || (var.csi_manager_password != null && var.csi_manager_password != "")
    error_message = "csi_manager_password is required when enable_csi=true."
  }
}

resource "vastdata_tenant" "csi" {
  count = var.enable_csi && var.create_csi_tenant ? 1 : 0

  name         = local.csi_tenant_name
  trash_gid    = var.trash_gid
  force_delete = var.force_delete_tenant

  local_provider_id = local.is_local_identity ? local.local_provider_id : null
  ldap_provider_id  = local.is_ldap ? var.ldap_provider_id : null
  ad_provider_id    = local.is_ad ? var.ad_provider_id : null
}

resource "vastdata_vip_pool" "csi" {
  count = var.enable_csi ? 1 : 0

  name        = local.csi_vip_pool_name
  role        = "PROTOCOLS"
  tenant_id   = local.csi_tenant_id
  subnet_cidr = var.csi_vip_pool.subnet_cidr
  ip_ranges   = var.csi_vip_pool.ip_ranges
  gw_ip       = try(var.csi_vip_pool.gw_ip, null)
  vlan        = try(var.csi_vip_pool.vlan, null)
  domain_name = try(var.csi_vip_pool.domain_name, null)
}

resource "vastdata_view_policy" "csi_nfs" {
  count = var.enable_csi ? 1 : 0

  name            = local.csi_view_policy_name
  tenant_id       = local.csi_tenant_id
  flavor          = "NFS"
  auth_source     = local.csi_view_policy_auth_source
  nfs_no_squash   = local.csi_view_policy_nfs_no_squash
  nfs_root_squash = local.csi_view_policy_nfs_root_squash
  vip_pools       = local.csi_view_policy_bind_vip_pool ? [vastdata_vip_pool.csi[0].id] : null
}

resource "vastdata_view" "docker_repo" {
  count = var.enable_csi ? 1 : 0

  name       = local.docker_repo_view_name
  path       = var.docker_repo_path
  create_dir = true
  delete_dir = var.delete_dir
  policy_id  = vastdata_view_policy.csi_nfs[0].id
  tenant_id  = local.csi_tenant_id
  protocols  = ["NFS", "NFS4"]
}

resource "vastdata_administrator_role" "csi_tenant_admin" {
  count = var.enable_csi ? 1 : 0

  name      = local.csi_tenant_admin_role_name
  tenant_id = local.csi_tenant_id

  # Fine-grained list takes precedence; otherwise permission types (current default).
  permissions_list = local.csi_tenant_admin_use_list ? local.csi_tenant_admin_perms_list : null
  permissions      = local.csi_tenant_admin_use_list ? null : var.csi_tenant_admin_permission_types
}

resource "vastdata_administrator_role" "csi_service_account" {
  count = var.enable_csi ? 1 : 0

  name      = local.csi_service_account_role_name
  tenant_id = local.csi_tenant_id

  permissions_list = local.csi_service_account_use_list ? local.csi_service_account_perms_list : null
  permissions      = local.csi_service_account_use_list ? null : var.csi_service_account_permission_types
}

resource "vastdata_administrator_manager" "csi_tenant_admin" {
  count = var.enable_csi ? 1 : 0

  username                     = local.csi_tenant_admin_username
  password                     = var.csi_manager_password
  password_expiration_disabled = true
  is_temporary_password        = false
  user_type                    = "TENANT_ADMIN"
  tenant_id                    = local.csi_tenant_id
  roles                        = [vastdata_administrator_role.csi_tenant_admin[0].id]
}

resource "vastdata_administrator_manager" "csi_driver" {
  count = var.enable_csi ? 1 : 0

  username                     = local.csi_driver_username
  password                     = var.csi_manager_password
  password_expiration_disabled = true
  is_temporary_password        = false
  user_type                    = "TENANT_ADMIN"
  tenant_id                    = local.csi_tenant_id
  roles                        = [vastdata_administrator_role.csi_service_account[0].id]
}
