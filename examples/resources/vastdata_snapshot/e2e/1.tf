# Copyright (c) HashiCorp, Inc.

data "vastdata_view_policy" "vastdb_view_policy_default" {
  name = "default"
}

resource "vastdata_view" "vastdb_view" {
  path       = "/vastdb_view/snap-example"
  policy_id  = data.vastdata_view_policy.vastdb_view_policy_default.id
  create_dir = true
}

resource "vastdata_snapshot" "vastdb_snapshot" {
  name            = "vastdb_snapshot"
  path            = "${vastdata_view.vastdb_view.path}/"
  indestructible  = false
  expiration_time = "2030-06-02T12:22:32Z"
}
