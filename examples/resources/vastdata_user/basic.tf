# Create a user with a specific UID.
resource "vastdata_user" "example-user" {
  name              = "example"
  uid               = 5001
  local_provider_id = 1
}

resource "vastdata_local_provider" "vastdb_local_provider1" {
  name = "vastdb_local_provider1"
}

// Create a user with local provider.
resource "vastdata_user" "vastdb_user1" {
  name              = "vastdb_user1"
  uid               = 5002
  local_provider_id = vastdata_local_provider.vastdb_local_provider1.id
}
