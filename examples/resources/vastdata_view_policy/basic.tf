# Copyright (c) HashiCorp, Inc.

resource "vastdata_view_policy" "vastdb_view_policy" {
  name          = "vastdb_view_policy"
  vip_pools     = [1, 2, 3]
  tenant_id     = 1
  flavor        = "NFS"
  nfs_no_squash = ["10.0.0.1", "10.0.0.2"]
}