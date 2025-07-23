# Copyright (c) HashiCorp, Inc.

resource "vastdata_group" "vastdb_group" {
  name = "vastdb_group"
  gid  = 30097
}

resource "vastdata_user" "vastdb_user" {
  name                = "vastdb_user"
  uid                 = 30109
  local               = true
  allow_create_bucket = true
  allow_delete_bucket = true
  s3_superuser        = false
  leading_gid         = vastdata_group.vastdb_group.gid
  gids = [
    1001,
    vastdata_group.vastdb_group.gid
  ]
}
