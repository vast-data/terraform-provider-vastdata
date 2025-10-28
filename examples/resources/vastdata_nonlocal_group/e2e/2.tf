
resource "vastdata_group" "vastdb_group" {
  name              = "vastdb_group"
  gid               = 1001
  local_provider_id = 1
}

resource "vastdata_nonlocal_group" "vastdb_nonlocal_group" {
  gid = vastdata_group.vastdb_group.gid
}
