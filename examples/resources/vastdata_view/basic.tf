# Copyright (c) HashiCorp, Inc.

resource "vastdata_view" "vastdb_view" {
  path       = "/vastdb_view/example"
  policy_id  = 2
  create_dir = true
  protocols  = ["NFS", "NFS4"]
}
