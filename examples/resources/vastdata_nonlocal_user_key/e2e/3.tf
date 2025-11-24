
data "vastdata_tenant" "vastdb_tenant" {
  name = "default"
}

resource "vastdata_user" "vastdb_user" {
  name              = "vastdb_user"
  uid               = 5001
  local_provider_id = 1
}

resource "vastdata_nonlocal_user_key" "vastdb_nonlocal_user_key" {
  username = vastdata_user.vastdb_user.name
}
