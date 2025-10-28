resource "vastdata_user" "vastdb_user" {
  name              = "runner"
  local_provider_id = 1
}

resource "vastdata_user_tenant_data" "vastdb_user_tenant_data" {
  user_id = vastdata_user.vastdb_user.id
}