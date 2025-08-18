# Copyright (c) HashiCorp, Inc.

variable snapshot_name {
    type = string
}

variable snapshot_tenant_name {
    type = string
}

resource vastdata_tenant tenant2 {
 name = var.snapshot_tenant_name
 client_ip_ranges {
         start_ip = "192.168.0.100"
         end_ip = "192.168.0.102"
    }
}

resource vastdata_snapshot snapshot {
    name = var.snapshot_name
    tenant_id = vastdata_tenant.tenant2.id
    indestructible = false
    expiration_time = "2025-06-02T12:22:32Z"
}

output tf_snapshot {
  value = vastdata_snapshot.snapshot
}

output tf_snapshot_tenant {
  value = vastdata_tenant.tenant2
}