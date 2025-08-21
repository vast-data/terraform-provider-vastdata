# Copyright (c) HashiCorp, Inc.

variable tenant_name {
    type = string
}

variable tenant_id {
    type = number
}

variable tenant_client_ip_ranges {
    type = list(object({
        start_ip = string
        end_ip = string
    }))
}

variable vippool_name {
    type = string
}

variable vippool_range1_start {
    type = string
}

variable vippool_range1_end {
    type = string
}

variable vippool_range2_start {
    type = string
}

variable vippool_range2_end {
    type = string
}

variable enable_weighted_balancing {
    type = bool
}

variable active_cnode_ids {
    type = list(number)
}

resource "vastdata_tenant" "tenant_for_adv_pool1" {
    name = var.tenant_name
    client_ip_ranges {
        start_ip = var.tenant_client_ip_ranges[0].start_ip
        end_ip = var.tenant_client_ip_ranges[0].end_ip
    }
}

resource "vastdata_vip_pool" "adv_pool1" {
    depends_on = [vastdata_tenant.tenant_for_adv_pool1]
    name = var.vippool_name
    role = "PROTOCOLS"
    subnet_cidr = 24
    enable_weighted_balancing = var.enable_weighted_balancing
    tenant_id = var.tenant_id
    active_cnode_ids = var.active_cnode_ids
    ip_ranges {
        start_ip = var.vippool_range1_start
        end_ip = var.vippool_range1_end
    }
    ip_ranges {
        start_ip = var.vippool_range2_start
        end_ip = var.vippool_range2_end
    }
}

output tf_vippool {
    value = vastdata_vip_pool.adv_pool1
}

output tf_tenant {
    value = vastdata_tenant.tenant_for_adv_pool1
}
