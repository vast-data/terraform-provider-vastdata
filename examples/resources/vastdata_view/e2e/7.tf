# ignore:e2e - nfs4_triggers endpoint requires cluster-level feature enablement
data "vastdata_view_policy" "vastdb_view_policy_default" {
  name = "default"
}

resource "vastdata_view" "vastdb_view" {
  path                 = "/vastdb_view/nfs4-triggers"
  policy_id            = data.vastdata_view_policy.vastdb_view_policy_default.id
  create_dir           = true
  protocols            = ["NFS4"]
  enable_nfs4_triggers = true
}
