# Copyright (c) HashiCorp, Inc.

variable tenant_name {
    type = string
}

variable pool_name {
    type = string
}

variable tenant_client_ip_ranges {
    type = list(object({
      start_ip = string
      end_ip = string
    }))
}

variable use_smb_privileged_user {
    type = bool
}

variable use_smb_privileged_group {
    type = bool
}

variable smb_privileged_group_full_access {
    type = bool
}

variable is_nfsv42_supported {
    type = bool
}

variable allow_locked_users {
    type = bool
}

variable allow_disabled_users {
    type = bool
}

variable use_smb_native {
    type = bool
}

variable vippool_ids {
    type = string
}

variable pool1_start_ip {
    type = string
}

variable pool1_end_ip {
    type = string
}

variable pool2_start_ip {
    type = string
}

variable pool2_end_ip {
    type = string
}

variable machine_name {
    type = string
}

resource vastdata_vip_pool ten_pool1 {
    name = "${var.pool_name}_1"
    role = "PROTOCOLS"
    subnet_cidr = "24"
    ip_ranges {
        start_ip = var.pool1_start_ip
        end_ip = var.pool1_end_ip
    }
}

resource vastdata_vip_pool ten_pool2 {
    name = "${var.pool_name}_2"
    role = "PROTOCOLS"
    subnet_cidr = "24"
    ip_ranges {
        start_ip = var.pool2_start_ip
        end_ip = var.pool2_end_ip
    }
}

resource vastdata_active_directory2 tenant_active_dir1 {
    machine_account_name = var.machine_name
    organizational_unit = "OU=VASTs,OU=VastENG,DC=VastENG,DC=lab"
    use_auto_discovery = "false"
    binddn = "cn=admin,dc=qa,dc=vastdata,dc=com"
    searchbase = "dc=qa,dc=vastdata,dc=com"
    bindpw = "vastdata"
    use_ldaps = "false"
    domain_name = "VastEng.lab"
    method = "simple"
    query_groups_mode = "COMPATIBLE"
    use_tls = "false"
    urls = ["ldap://10.27.252.30"]
}

# Create a tenant with the name and client_ip_ranges
resource vastdata_tenant tenant1 {
  name = var.tenant_name
  #ad_provider_id = vastdata_active_directory2.tenant_active_dir1.id

  dynamic "client_ip_ranges" {
    for_each = var.tenant_client_ip_ranges
    content {
      start_ip = client_ip_ranges.value["start_ip"]
      end_ip = client_ip_ranges.value["end_ip"]
    }
  }
  use_smb_privileged_user = var.use_smb_privileged_user
  use_smb_privileged_group = var.use_smb_privileged_group
  smb_privileged_group_full_access = var.smb_privileged_group_full_access
  is_nfsv42_supported = var.is_nfsv42_supported
  allow_locked_users = var.allow_locked_users
  allow_disabled_users = var.allow_disabled_users
  use_smb_native = var.use_smb_native
  vippool_ids = var.vippool_ids == "all" ? [vastdata_vip_pool.ten_pool1.id, vastdata_vip_pool.ten_pool2.id] : []
}

output tf_tenant {
  value = vastdata_tenant.tenant1
}

output tf_pool1{
  value = vastdata_vip_pool.ten_pool1
}

output tf_pool2{
  value = vastdata_vip_pool.ten_pool2
}