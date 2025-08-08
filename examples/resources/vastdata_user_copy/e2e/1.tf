resource "vastdata_local_provider" "vastdb_local_provider1" {
  name = "vastdb_local_provider"
}

resource "vastdata_user" "vastdb_user1" {
  name = "vastdb_user"
  uid  = 30017
}

resource "vastdata_user_copy" "copy_specific_users" {
  destination_provider_id = vastdata_local_provider.vastdb_local_provider1.id
  user_ids = [
    vastdata_user.vastdb_user1.id,
  ]
}
