resource "vastdata_global_local_snapshot" "vastdb_local_snapshot" {
  name               = "vastdb_local_snapshot"
  loanee_root_path   = "/vastdb_local_snapshot"
  enabled            = true
  loanee_snapshot_id = 2
  loanee_tenant_id   = 1
}
