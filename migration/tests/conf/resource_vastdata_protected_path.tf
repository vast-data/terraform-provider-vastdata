# Copyright (c) HashiCorp, Inc.

variable protected_path_name {
    type = string
}

variable protected_path_tenant_name {
    type = string
}

variable protected_path_policy_name {
    type = string
}

variable enabled {
    type = bool
}

resource vastdata_protection_policy ppolicy1 {
        name = var.protected_path_policy_name
        indestructible = "false"
        prefix = "policy-1"
        clone_type = "LOCAL"
        frames {
                every = "1D"
                keep_local = "2D"
                start_at = "2023-06-04 09:00:00"
        }
}

resource vastdata_tenant tenant2 {
        name = var.protected_path_tenant_name
        client_ip_ranges {
                start_ip = "192.168.0.50"
                end_ip = "192.168.0.51"
        }
}

resource vastdata_protected_path protected_path1 {
        name = var.protected_path_name
        source_dir = "/"
        tenant_id = vastdata_tenant.tenant2.id
        target_exported_dir = "/view1"
        protection_policy_id = vastdata_protection_policy.ppolicy1.id
        enabled = var.enabled
        capabilities = "ASYNC_REPLICATION"
}

output tf_protected_path {
  value = vastdata_protected_path.protected_path1
}

output tf_protected_path_policy {
  value = vastdata_protection_policy.ppolicy1
}

output tf_protected_path_tenant {
  value = vastdata_tenant.tenant2
}