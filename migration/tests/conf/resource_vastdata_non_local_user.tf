# Copyright (c) HashiCorp, Inc.

variable context {
    type = string
}

variable user_uid {
    type = number
}

variable user_name {
    type = string
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

variable use_s3_policies {
    type = string
}

variable s3_policy_id1 {
    type = number
}

variable s3_policy_id2 {
  type = number
}

resource vastdata_ldap ldap_for_non_local_user1 {
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

resource vastdata_tenant tenant_for_non_local_user1 {
  name = var.tenant_name
  ldap_provider_id = vastdata_ldap.ldap_for_non_local_user1.id

  dynamic "client_ip_ranges" {
    for_each = var.tenant_client_ip_ranges
    content {
      start_ip = client_ip_ranges.value["start_ip"]
      end_ip = client_ip_ranges.value["end_ip"]
    }
  }
}

resource "vastdata_non_local_user" "non_local_user1" {
    uid                 = var.user_uid
    context   = var.context
    tenant_id           = vastdata_tenant.tenant_for_non_local_user1.id
    allow_create_bucket = true
    allow_delete_bucket = false
    s3_policies_ids     = var.use_s3_policies == "none" ? [] : (
        var.use_s3_policies == "all" ? [
            var.s3_policy_id1,
            var.s3_policy_id2
        ] : [
            var.s3_policy_id1
        ]
    )
}

output tf_user {
  value = vastdata_non_local_user.non_local_user1
}

output tf_tenant {
  value = vastdata_tenant.tenant_for_non_local_user1
}

output tf_ldap {
    value = vastdata_ldap.ldap_for_non_local_user1
    sensitive = true
}

# Data Source:
data "vastdata_non_local_user" "non_local_user_ds1" {
    username  = var.user_name
    context   = var.context
    tenant_id = vastdata_tenant.tenant_for_non_local_user1.id
}

# Output for data source validation
output tf_user_ds {
    value = data.vastdata_non_local_user.non_local_user_ds1
}