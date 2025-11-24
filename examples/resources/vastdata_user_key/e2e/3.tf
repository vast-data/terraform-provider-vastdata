resource "vastdata_local_provider" "vastdb_local_provider" {
  name       = "vastdb_local_provider"
  managed_by = ["TENANT_ADMIN", "SUPER_ADMIN"]
}

resource "vastdata_tenant" "vastdb_tenant" {
  name              = "vastdb-tenant"
  client_ip_ranges  = [["192.168.100.1", "192.168.100.10"]]
  local_provider_id = vastdata_local_provider.vastdb_local_provider.id
}

resource "vastdata_user" "vastdb_user" {
  name              = "vastdb_user"
  uid               = 5001
  local_provider_id = vastdata_local_provider.vastdb_local_provider.id
}

resource "vastdata_user_key" "vastdb_user_key" {
  user_id   = vastdata_user.vastdb_user.id
  tenant_id = vastdata_tenant.vastdb_tenant.id
  enabled   = false
}
