resource "vastdata_administrator_role" "vastdb_role" {
  name             = "vastdb_role"
  permissions_list = ["create_events"]
  realm            = 4
}
