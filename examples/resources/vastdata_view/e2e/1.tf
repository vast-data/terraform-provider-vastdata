
data "vastdata_view_policy" "vastdb_view_policy_default" {
  name = "default"
}

resource "vastdata_view" "vastdb_view" {
  path       = "/vastdb_view/example"
  policy_id  = data.vastdata_view_policy.vastdb_view_policy_default.id
  create_dir = true
  protocols  = ["NFS", "NFS4"]
}