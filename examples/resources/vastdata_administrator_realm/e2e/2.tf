
data "vastdata_tenant" "vastdb_tenant" {
  name = "default"
}

resource "vastdata_administrator_realm" "vastdb_realm" {
  name         = "vastdb_realm"
  object_types = ["nic", "viewpolicy"]
  tenant_id    = data.vastdata_tenant.vastdb_tenant.id
}