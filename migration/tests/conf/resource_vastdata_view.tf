# Copyright (c) HashiCorp, Inc.

# To refer to a specific system, need to add `provider = vastdata.system<IDX>` where
# IDX is the system's index (in the order of appearance in Comet's commandline)

variable view_path {
    type = string
}

variable view_protocols {
    type = list(string)
}

variable view_policy_name {
    type = string
}

variable policy_use_auth_provider {
    type = bool
}

variable policy_auth_source {
    type = string
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

resource vastdata_tenant view_tenant1 {
  name = var.tenant_name

  dynamic "client_ip_ranges" {
    for_each = var.tenant_client_ip_ranges
    content {
      start_ip = client_ip_ranges.value["start_ip"]
      end_ip = client_ip_ranges.value["end_ip"]
    }
  }
}

# Create a view with NFS & NFSv4 protocols
resource vastdata_view_policy viewpolicy1 {
   name = var.view_policy_name
   flavor = "NFS"
   use_auth_provider = var.policy_use_auth_provider
   auth_source = var.policy_auth_source
   nfs_no_squash = ["10.0.0.1", "10.0.0.2"]
   tenant_id = vastdata_tenant.view_tenant1.id
}

resource vastdata_view view1 {
  path = var.view_path
  policy_id = vastdata_view_policy.viewpolicy1.id
  create_dir = "true"
  protocols = var.view_protocols
  nfs_interop_flags = "BOTH_NFS3_AND_NFS4_INTEROP_DISABLED"
  tenant_id = vastdata_tenant.view_tenant1.id
}

output tf_view {
  value = vastdata_view.view1
}

output tf_view_view_policy {
  value = vastdata_view_policy.viewpolicy1
}

output tf_tenant {
  value = vastdata_tenant.view_tenant1
}