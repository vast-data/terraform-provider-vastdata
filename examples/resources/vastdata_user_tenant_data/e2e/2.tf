data "vastdata_tenant" "vastdb_tenant" {
  name = "default"
}

resource "vastdata_user" "vastdb_user" {
  name              = "vastdb_user"
  local_provider_id = 1
}

resource "vastdata_user_tenant_data" "vastdb_user_tenant_data" {
  user_id   = vastdata_user.vastdb_user.id
  tenant_id = data.vastdata_tenant.vastdb_tenant.id
}