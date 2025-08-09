
# Create a user with a specific UID.
resource "vastdata_user" "example-user" {
  name = "example"
  uid  = 9000
}


resource "vastdata_local_provider" "vastdb_local_provider1" {
  name = "vastdb_local_provider1"
}

// Create a user with local provider.
resource "vastdata_user" "vastdb_user1" {
  name              = "vastdb_user1"
  uid               = 30117
  local_provider_id = vastdata_local_provider.vastdb_local_provider1.id
}
