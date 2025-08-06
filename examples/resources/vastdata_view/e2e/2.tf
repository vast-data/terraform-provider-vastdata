
data "vastdata_tenant" "vastdb_default_tenant" {
  name = "default"
}

data "vastdata_view_policy" "vastdb_view_policy_default" {
  name = "default"
}

resource "vastdata_view" "vastdb_view" {
  path                       = "/vastdb_view/example"
  alias                      = "/vastdb_view-aliased"
  tenant_id                  = data.vastdata_tenant.vastdb_default_tenant.id
  policy_id                  = data.vastdata_view_policy.vastdb_view_policy_default.id
  create_dir                 = true
  select_for_live_monitoring = true
  protocols                  = ["NFS"]
}