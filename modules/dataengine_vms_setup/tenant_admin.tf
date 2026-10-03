# Optional DataEngine tenant-admin manager.
# Role: create a new one (default) or assign an existing role by name.

check "tenant_admin_password_required" {
  assert {
    condition     = !var.create_tenant_admin || (var.tenant_admin_password != null && var.tenant_admin_password != "")
    error_message = "tenant_admin_password is required when create_tenant_admin=true."
  }
}

check "tenant_admin_existing_role_name_required" {
  assert {
    condition = !var.create_tenant_admin || var.create_tenant_admin_role || (
      var.tenant_admin_role_name != null && var.tenant_admin_role_name != ""
    )
    error_message = "tenant_admin_role_name is required when create_tenant_admin=true and create_tenant_admin_role=false (lookup existing role by name)."
  }
}

data "vastdata_administrator_role" "tenant_admin" {
  count = local.tenant_admin_use_existing_role ? 1 : 0

  name = local.tenant_admin_role_name
}

resource "vastdata_administrator_role" "tenant_admin" {
  count = local.tenant_admin_create_role ? 1 : 0

  name      = local.tenant_admin_role_name
  tenant_id = vastdata_tenant.dataengine.id

  permissions_list = local.tenant_admin_use_list ? local.tenant_admin_perms_list : null
  permissions      = local.tenant_admin_use_list ? null : var.tenant_admin_permission_types
}

resource "vastdata_administrator_manager" "tenant_admin" {
  count = var.create_tenant_admin ? 1 : 0

  username                     = local.tenant_admin_username
  password                     = var.tenant_admin_password
  password_expiration_disabled = true
  is_temporary_password        = false
  user_type                    = "TENANT_ADMIN"
  tenant_id                    = vastdata_tenant.dataengine.id
  roles                        = [local.tenant_admin_role_id]
}
