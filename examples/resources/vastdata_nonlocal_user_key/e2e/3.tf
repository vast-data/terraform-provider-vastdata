# Copyright (c) HashiCorp, Inc.

data "vastdata_tenant" "vastdb_tenant" {
  name = "default"
}


resource "vastdata_user" "vastdb_user" {
  name = "vastdb_user"
  uid  = 30109
}


resource "vastdata_nonlocal_user_key" "vastdb_nonlocal_user_key" {
  username = vastdata_user.vastdb_user.name
}
