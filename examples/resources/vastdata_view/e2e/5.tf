
data "vastdata_view_policy" "vastdb_view_policy_default" {
  name = "default"
}

resource "vastdata_view" "vastdb_view" {
  path                 = "/vastdb_view/subsystem"
  name                 = "vastdb-subsystem"
  create_dir           = true
  is_default_subsystem = true
  policy_id            = data.vastdata_view_policy.vastdb_view_policy_default.id
  protocols            = ["BLOCK"]
}
