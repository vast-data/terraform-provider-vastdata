# ignore:e2e

resource "vastdata_user" "vastdb_user" {
  name              = "vastdb_user"
  uid               = 30109
  local_provider_id = 1
  gids = [
    1001
  ]
}
