resource "vastdata_local_provider" "vastdb_local_provider1" {
  name = "vastdb_local_provider1"
}

resource "vastdata_user" "vastdb_user1" {
  name              = "vastdb_user1"
  local_provider_id = vastdata_local_provider.vastdb_local_provider1.id
}
