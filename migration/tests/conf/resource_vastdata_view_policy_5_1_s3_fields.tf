# Copyright (c) HashiCorp, Inc.

# To refer to a specific system, need to add `provider = vastdata.system<IDX>` where
# IDX is the system's index (in the order of appearance in Comet's commandline)

variable view_policy_name {
    type = string
}

variable tenant_name {
    type = string
}

variable pool_name {
    type = string
}

variable s3_special_chars_support {
    type = bool
}

variable smb_is_ca {
    type = bool
}

variable nfs_case_insensitive {
    type = bool
}

variable enable_access_to_snapshot_dir_in_subdirs {
    type = bool
}

variable enable_visibility_of_snapshot_dir {
    type = bool
}

variable nfs_enforce_tls {
    type = bool
}

variable path_length {
    type = string
}

variable nfs_minimal_protection_level {
    type = string
}

resource vastdata_vip_pool s3pool1 {
    name = var.pool_name
    role = "PROTOCOLS"
    subnet_cidr = "24"
    ip_ranges {
        end_ip = "11.0.0.80"
        start_ip = "11.0.0.70"
    }
}

resource vastdata_tenant s3tenant1 {
    name = var.tenant_name
    client_ip_ranges {
        start_ip = "192.168.0.100"
        end_ip = "192.168.0.200"
    }
}

resource vastdata_view_policy s3vpolicy1 {
    name = var.view_policy_name
    vip_pools = [vastdata_vip_pool.s3pool1.id]
    tenant_id = vastdata_tenant.s3tenant1.id
    flavor = "S3_NATIVE"
    nfs_no_squash = ["10.0.0.1","10.0.0.2"]

    s3_special_chars_support = var.s3_special_chars_support
    smb_is_ca = var.smb_is_ca
    nfs_case_insensitive = var.nfs_case_insensitive
    enable_access_to_snapshot_dir_in_subdirs = var.enable_access_to_snapshot_dir_in_subdirs
    enable_visibility_of_snapshot_dir = var.enable_visibility_of_snapshot_dir
    nfs_enforce_tls = var.nfs_enforce_tls
    path_length = var.path_length
    nfs_minimal_protection_level = var.nfs_minimal_protection_level
}

output tf_policy {
  value = vastdata_view_policy.s3vpolicy1
}

output tf_policy_tenant {
  value = vastdata_tenant.s3tenant1
}

output tf_policy_vippool {
  value = vastdata_vip_pool.s3pool1
}