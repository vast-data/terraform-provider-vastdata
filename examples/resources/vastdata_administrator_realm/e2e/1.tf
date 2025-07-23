# Copyright (c) HashiCorp, Inc.

resource "vastdata_administrator_realm" "vastdb_realm" {
  name         = "vastdb_realm"
  object_types = ["nic", "viewpolicy"]
}