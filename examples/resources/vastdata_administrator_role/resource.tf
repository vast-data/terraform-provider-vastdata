resource "vastdata_administrator_role" "vastdb_role" {
  name             = "vastdb_role"
  permissions_list = ["create_events"]
  realm            = 4
}

# ---------------------
# Complete examples
# ---------------------

resource "vastdata_administrator_realm" "vastdb_realm" {
  name         = "vastdb_realm"
  object_types = ["nic", "viewpolicy"]
}

resource "vastdata_administrator_role" "vastdb_role" {
  name             = "vastdb_role"
  permissions_list = ["create_events"]
  realm            = vastdata_administrator_realm.vastdb_realm.id
}

# --------------------


resource "vastdata_tenant" "vastdb_tenant1" {
  name = "vastdbtenant1"

}

resource "vastdata_tenant" "vastdb_tenant2" {
  name = "vastdbtenant2"

}

resource "vastdata_administrator_role" "vastdb_role" {
  name = "vastdb_role"
  tenant_ids = [
    vastdata_tenant.vastdb_tenant1.id,
    vastdata_tenant.vastdb_tenant2.id
  ]

  permissions_list = [
    "create_database", "create_hardware", "view_hardware", "edit_database", "delete_hardware",
  ]
}

# --------------------

