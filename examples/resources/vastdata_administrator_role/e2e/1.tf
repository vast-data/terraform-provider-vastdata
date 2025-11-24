resource "vastdata_administrator_realm" "vastdb_realm" {
  name         = "vastdb_realm"
  object_types = ["nic", "viewpolicy"]
}

resource "vastdata_administrator_role" "vastdb_role" {
  name             = "vastdb_role"
  permissions_list = ["create_events"]
  realm            = vastdata_administrator_realm.vastdb_realm.id
}
