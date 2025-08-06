resource "vastdata_volume" "vastdb_volume" {
  name    = "vastdb-volume"
  size    = 10737418240
  view_id = 1
}

# ---------------------
# Complete examples
# ---------------------

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

resource "vastdata_volume" "vastdb_volume" {
  name    = "vastdb-volume"
  size    = 10737418240
  view_id = vastdata_view.vastdb_view.id
}

# --------------------

