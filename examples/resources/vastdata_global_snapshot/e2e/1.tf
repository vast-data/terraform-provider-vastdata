resource "vastdata_tenant" "vastdb_tenant" {
  name = "vastdb-drift-tenant"
}

resource "vastdata_view_policy" "vastdb_view_policy" {
  name          = "vastdb-vp"
  flavor        = "NFS"
  nfs_no_squash = ["10.0.0.1"]
}

resource "vastdata_view" "vastdb_view1" {
  path       = "/vastdb-gs-view"
  policy_id  = vastdata_view_policy.vastdb_view_policy.id
  create_dir = true
  protocols  = ["NFS", "NFS4"]
}

resource "vastdata_view" "vastdb_view2" {
  path       = "/vastdb-gs-snapclone"
  policy_id  = vastdata_view_policy.vastdb_view_policy.id
  create_dir = true
  protocols  = ["NFS", "NFS4"]
}


resource "vastdata_snapshot" "snap1" {
  path = "${vastdata_view.vastdb_view1.path}/"
  name = "vastdb-gs-snapshot"
}

resource "vastdata_global_snapshot" "gsnap1" {
  name               = "vastdb-gs-global-snap"
  loanee_root_path   = "${vastdata_view.vastdb_view2.path}/"
  loanee_snapshot_id = vastdata_snapshot.snap1.id
  loanee_tenant_id   = vastdata_tenant.vastdb_tenant.id
  enabled            = true
}
