# Copyright (c) HashiCorp, Inc.

resource "vastdata_user_key" "vastdb_user_key" {
  username  = "example-user"
  tenant_id = "example-tenant-id"
}
