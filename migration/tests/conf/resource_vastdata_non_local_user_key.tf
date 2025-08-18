# Copyright (c) HashiCorp, Inc.

variable user_uid {
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

resource vastdata_ldap ldap_for_non_local_user2 {
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

resource vastdata_tenant tenant_for_non_local_user2 {
  name = var.tenant_name
  ldap_provider_id = vastdata_ldap.ldap_for_non_local_user2.id

  dynamic "client_ip_ranges" {
    for_each = var.tenant_client_ip_ranges
    content {
      start_ip = client_ip_ranges.value["start_ip"]
      end_ip = client_ip_ranges.value["end_ip"]
    }
  }
}

resource "vastdata_non_local_user" "non_local_user2" {
    uid                 = var.user_uid
    tenant_id           = vastdata_tenant.tenant_for_non_local_user2.id
    allow_create_bucket = true
    allow_delete_bucket = false
    s3_policies_ids     = []
}

resource "vastdata_non_local_user_key" "external_user_key2" {
    depends_on          = [vastdata_non_local_user.non_local_user2]
    uid                 = var.user_uid
    tenant_id           = vastdata_tenant.tenant_for_non_local_user2.id
    enabled             = true
}

resource "vastdata_non_local_user_key" "external_user_key_encrypted2" {
    depends_on          = [vastdata_non_local_user.non_local_user2]
    uid                 = var.user_uid
    tenant_id           = vastdata_tenant.tenant_for_non_local_user2.id
    enabled             = true
}

/* resource "vastdata_non_local_user_key" "external_user_key_encrypted2" {
    depends_on          = [vastdata_non_local_user.non_local_user2]
    uid                 = var.user_uid
    tenant_id           = vastdata_tenant.tenant_for_non_local_user2.id
    enabled             = true
    pgp_public_key = <<-EOT
    -----BEGIN PGP PUBLIC KEY BLOCK-----
Version: Keybase OpenPGP v1.0.0
Comment: https://keybase.io/crypto

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
qyjqu/s3g8Fn6+sNO0X5aRM=
=Wtzq
-----END PGP PUBLIC KEY BLOCK-----
    EOT
} */

output tf_user_key {
  value = vastdata_non_local_user_key.external_user_key2
  sensitive = true
}

output tf_user_key_encrypted {
  value = vastdata_non_local_user_key.external_user_key_encrypted2
  sensitive = true
}

output tf_user {
  value = vastdata_non_local_user.non_local_user2
}

output tf_tenant {
  value = vastdata_tenant.tenant_for_non_local_user2
}