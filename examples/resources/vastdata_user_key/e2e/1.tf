
data "vastdata_tenant" "vastdb_tenant" {
  name = "default"
}

resource "vastdata_user" "vastdb_user" {
  name = "vastdb_user"
  uid  = 30109
}

resource "vastdata_user_key" "vastdb_user_key" {
  username  = vastdata_user.vastdb_user.name
  tenant_id = data.vastdata_tenant.vastdb_tenant.id
  enabled   = false
}
