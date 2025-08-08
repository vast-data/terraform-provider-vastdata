resource "vastdata_administrator_role" "vastdb_role" {
  name        = "vastdb_role"
  permissions = "view"
  realm       = 4
}
