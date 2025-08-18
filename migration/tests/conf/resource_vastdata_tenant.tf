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

#Create a tenant with the name and client_ip_ranges
resource vastdata_tenant tenant1 {
  name = var.tenant_name

  dynamic "client_ip_ranges" {
    for_each = var.tenant_client_ip_ranges
    content {
      start_ip = client_ip_ranges.value["start_ip"]
      end_ip = client_ip_ranges.value["end_ip"]
    }
  }
}

output tf_tenant {
  value = vastdata_tenant.tenant1
}