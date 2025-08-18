# Copyright (c) HashiCorp, Inc.

variable view_policy_name {
    type = string
}

variable tenant_name {
    type = string
}

variable pool_name {
    type = string
}

variable pool_start_ip {
    type = string
}

variable pool_end_ip {
    type = string
}

variable tenant_client_start_ip {
    type = string
}

variable tenant_client_end_ip {
    type = string
}

variable nfs_all_squash {
    type = string
}

variable nfs_root_squash {
    type = string
}

variable nfs_no_squash {
    type = string
}

variable nfs_read_only {
    type = string
}

variable nfs_read_write {
    type = string
}

variable s3_read_only {
    type = string
}

variable s3_read_write {
    type = string
}

variable smb_read_only {
    type = string
}

variable smb_read_write {
    type = string
}

resource vastdata_vip_pool pool52 {
    name = var.pool_name
    role = "PROTOCOLS"
    subnet_cidr = "24"
    ip_ranges {
        start_ip = var.pool_start_ip
        end_ip = var.pool_end_ip
    }
}

resource vastdata_tenant tenant52 {
    name = var.tenant_name
    client_ip_ranges {
        start_ip = var.tenant_client_start_ip
        end_ip = var.tenant_client_end_ip
    }
}

resource vastdata_view_policy vpolicy52 {
    name = var.view_policy_name
    vip_pools = [vastdata_vip_pool.pool51.id]
    tenant_id = vastdata_tenant.tenant51.id
    flavor = "NFS"

    nfs_all_squash= var.nfs_all_squash
    nfs_root_squash= var.nfs_root_squash
    nfs_no_squash= var.nfs_no_squash
    nfs_read_only= var.nfs_read_only
    nfs_read_write= var.nfs_read_write
    s3_read_only= var.s3_read_only
    s3_read_write= var.s3_read_write
    smb_read_only= var.smb_read_only
    smb_read_write= var.smb_read_write
}

output tf_policy {
  value = vastdata_view_policy.vpolicy52
}

output tf_policy_tenant {
  value = vastdata_tenant.tenant52
}

output tf_policy_vippool {
  value = vastdata_vip_pool.pool52
}