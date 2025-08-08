data "vastdata_view_policy" "vastdb_policy" {
  name = "default"
}

data "vastdata_tenant" "vastdb_tenant" {
  name = "default"
}


resource "vastdata_view" "vastdb_view" {
  path       = "/vastdb_view/snapclone"
  policy_id  = data.vastdata_view_policy.vastdb_policy.id
  tenant_id  = data.vastdata_tenant.vastdb_tenant.id
  create_dir = "true"
}

resource "vastdata_snapshot" "vastdb_snapshot" {
  name = "vastdb_snapshot"
  path = "${vastdata_view.vastdb_view.path}/"
}


resource "vastdata_global_local_snapshot" "vastdb_local_snapshot" {
  name               = "vastdb_local_snapshot"
  loanee_root_path   = "/vastdb_local_snapshot"
  enabled            = true
  loanee_snapshot_id = vastdata_snapshot.vastdb_snapshot.id
  loanee_tenant_id   = data.vastdata_tenant.vastdb_tenant.id
}
