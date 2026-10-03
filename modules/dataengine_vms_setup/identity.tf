check "local_gid_uid_required" {
  assert {
    condition     = !local.is_local_identity || (var.app_group_gid != null && var.broker_owner_uid != null)
    error_message = "app_group_gid and broker_owner_uid are required when identity_provider_type=local."
  }
}

check "ldap_provider_required" {
  assert {
    condition     = !local.is_ldap || var.ldap_provider_id != null
    error_message = "ldap_provider_id is required when identity_provider_type=ldap."
  }
}

check "ad_provider_required" {
  assert {
    condition     = !local.is_ad || var.ad_provider_id != null
    error_message = "ad_provider_id is required when identity_provider_type=active_directory."
  }
}

check "external_identity_names_required" {
  assert {
    condition = local.is_local_identity || (
      var.app_group_name != null && var.app_group_name != "" &&
      var.broker_owner_name != null && var.broker_owner_name != ""
    )
    error_message = "app_group_name and broker_owner_name (existing LDAP/AD names) are required when identity_provider_type is ldap or active_directory."
  }
}

resource "vastdata_local_provider" "this" {
  count = local.is_local_identity && var.create_local_provider ? 1 : 0

  name = local.local_provider_name
}

resource "vastdata_group" "app" {
  count = local.is_local_identity ? 1 : 0

  name              = local.app_group_name
  gid               = var.app_group_gid
  local_provider_id = local.local_provider_id
}

resource "vastdata_user" "broker_owner" {
  count = local.is_local_identity ? 1 : 0

  name                = local.broker_owner
  uid                 = var.broker_owner_uid
  local_provider_id   = local.local_provider_id
  leading_gid         = vastdata_group.app[0].gid
  gids                = [vastdata_group.app[0].gid]
  allow_create_bucket = true
  allow_delete_bucket = true

  password = var.broker_owner_password == null ? null : var.broker_owner_password
}
