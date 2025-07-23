
# Create a user with a specific UID.
resource "vastdata_user" "example-user" {
  name = "example"
  uid  = 9000
}

# ---------------------
# Complete examples
# ---------------------


resource "vastdata_user" "vastdb_user" {
  name = "vastdb_user"
  uid  = 30109
  gids = [
    1001
  ]
}

# --------------------


resource "vastdata_group" "vastdb_group" {
  name = "vastdb_group"
  gid  = 30097
}

resource "vastdata_user" "vastdb_user" {
  name                = "vastdb_user"
  uid                 = 30109
  local               = true
  allow_create_bucket = true
  allow_delete_bucket = true
  s3_superuser        = false
  leading_gid         = vastdata_group.vastdb_group.gid
  gids = [
    1001,
    vastdata_group.vastdb_group.gid
  ]
}

# --------------------

