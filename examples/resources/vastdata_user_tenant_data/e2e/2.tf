resource "vastdata_user" "vastdb_user" {
  name = "runner"
}

resource "vastdata_user_tenant_data" "vastdb_user_tenant_data" {
  user_id = vastdata_user.vastdb_user.id
}