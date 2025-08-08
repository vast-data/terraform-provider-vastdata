resource "vastdata_administrator_realm" "vastdb_realm" {
  name         = "vastdb_realm"
  object_types = ["nic", "viewpolicy"]
}

resource "vastdata_administrator_role" "vastdb_role" {
  name        = "vastdb_role"
  permissions = "view"
  realm       = vastdata_administrator_realm.vastdb_realm.id
}
