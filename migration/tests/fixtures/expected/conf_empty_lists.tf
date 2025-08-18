# Copyright (c) HashiCorp, Inc.

variable user_uid {
    type = number
}

variable local_user_name {
    type = string
}

variable local_user_uid {
    type = number
}

variable local_group_name {
    type = string
}

variable local_group_gid {
    type = number
}

variable tenant_name {
    type = string
}

variable tenant_client_ip_ranges {
    type = list(object({
      start_ip = string
      end_ip = string
    }))
}

resource vastdata_ldap ldap_for_empty_lists1 {
    domain_name = "VastEng.lab"
    urls = ["ldap://10.27.252.30"]
    binddn = "cn=admin,dc=qa,dc=vastdata,dc=com"
    searchbase = "dc=qa,dc=vastdata,dc=com"
    bindpw = "vastdata"
    use_auto_discovery = "false"
    use_ldaps = "false"
    port = "389"
    method = "simple"
    query_groups_mode = "COMPATIBLE"
    use_tls = "false"
}

resource vastdata_tenant tenant_for_empty_lists1 {
  name = var.tenant_name
  ldap_provider_id = vastdata_ldap.ldap_for_empty_lists1.id

  dynamic "client_ip_ranges" {
    for_each = var.tenant_client_ip_ranges
    content {
      start_ip = client_ip_ranges.value["start_ip"]
      end_ip = client_ip_ranges.value["end_ip"]
    }
  }

  # vippool_ids = [] check
}

resource "vastdata_nonlocal_user" "non_local_user_for_empty_lists1" {
    uid                 = var.user_uid
    tenant_id           = vastdata_tenant.tenant_for_empty_lists1.id
    allow_create_bucket = true
    allow_delete_bucket = false
    context = "ldap"
    # s3_policies_ids     = [] check
}

resource vastdata_group user_group_for_empty_lists1 {
  name = var.local_group_name
  gid = var.local_group_gid
  # s3_policies_ids     = [] check
}

resource "vastdata_user" "local_user_for_empty_lists1" {
  name        = var.local_user_name
  uid         = var.local_user_uid
  leading_gid = resource.vastdata_group.user_group_for_empty_lists1.gid
  # s3_policies_ids     = [] check
}

output tf_tenant {
  value = vastdata_tenant.tenant_for_empty_lists1
}

output tf_non_local_user {
  value = vastdata_non_local_user.non_local_user_for_empty_lists1
}

output tf_local_group {
  value = vastdata_group.user_group_for_empty_lists1
}

output tf_local_user {
  value = vastdata_user.local_user_for_empty_lists1
}
