# Copyright (c) HashiCorp, Inc.

# To refer to a specific system, need to add `provider = vastdata.system<IDX>` where
# IDX is the system's index (in the order of appearance in Comet's commandline)

variable view_policy_name {
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

variable tenant_name {
    type = string
}

variable tenant_client_start_ip {
    type = string
}

variable tenant_client_end_ip {
    type = string
}

resource vastdata_vip_pool pool1 {
    name = var.pool_name
    role = "PROTOCOLS"
    subnet_cidr = "24"
    ip_ranges {
        start_ip = var.pool_start_ip
        end_ip = var.pool_end_ip
    }
}

resource vastdata_tenant tenant1 {
    name = var.tenant_name
    client_ip_ranges {
        start_ip = var.tenant_client_start_ip
        end_ip = var.tenant_client_end_ip
    }
}

resource vastdata_view_policy vpolicy1 {
    name = var.view_policy_name
    vip_pools = [vastdata_vip_pool.pool1.id]
    tenant_id = vastdata_tenant.tenant1.id
    flavor = "NFS"
    nfs_no_squash = ["10.0.0.1","10.0.0.2"]
}

output tf_policy {
  value = vastdata_view_policy.vpolicy1
}

output tf_policy_tenant {
  value = vastdata_tenant.tenant1
}

output tf_policy_vippool {
  value = vastdata_vip_pool.pool1
}