resource "vastdata_group" "vastdb_group" {
  name              = "vastdb_group"
  gid               = 5001
  local_provider_id = 1

}

resource "vastdata_user" "vastdb_user" {
  name                = "vastdb_user"
  uid                 = 5003
  allow_create_bucket = true
  allow_delete_bucket = true
  s3_superuser        = false
  leading_gid         = vastdata_group.vastdb_group.gid
  local_provider_id   = 1
  gids = [
    vastdata_group.vastdb_group.gid
  ]
}
