# Copyright (c) HashiCorp, Inc.

variable tenant_name {
    type = string
}

variable tenant_client_ip_ranges {
    type = list(object({
      start_ip = string
      end_ip = string
    }))
}

resource "vastdata_tenant" "fallback_tenant1" {
    name = var.tenant_name
    client_ip_ranges {
        start_ip = var.tenant_client_ip_ranges[0].start_ip
        end_ip = var.tenant_client_ip_ranges[0].end_ip
    }
}

output tf_tenant {
    value = vastdata_tenant.fallback_tenant1
}
