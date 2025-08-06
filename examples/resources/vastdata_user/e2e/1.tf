
resource "vastdata_user" "vastdb_user" {
  name = "vastdb_user"
  uid  = 30109
  gids = [
    1001
  ]
}
