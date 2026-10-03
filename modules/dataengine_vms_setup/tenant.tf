resource "vastdata_tenant" "dataengine" {
  name                         = local.tenant_name
  application_users_group_name = local.app_group_name
  force_delete                 = var.force_delete_tenant

  local_provider_id = local.is_local_identity ? local.local_provider_id : null
  ldap_provider_id  = local.is_ldap ? var.ldap_provider_id : null
  ad_provider_id    = local.is_ad ? var.ad_provider_id : null

  # Trash GID on DE tenant only when CSI helpers share this tenant.
  trash_gid = var.enable_csi && !var.create_csi_tenant ? var.trash_gid : null
}
