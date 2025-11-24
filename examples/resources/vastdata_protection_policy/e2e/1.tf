resource "vastdata_tenant" "vastdb_tenant" {
  name             = "vastdb-tenant"
  client_ip_ranges = [["192.168.0.52", "192.168.0.53"]]
}

resource "vastdata_protection_policy" "vastdb_ppolicy" {
  name           = "vastdb-ppolicy"
  indestructible = "false"
  prefix         = "ppolicy"
  clone_type     = "LOCAL"
  tenant_id      = vastdata_tenant.vastdb_tenant.id
  frames = [{
    every      = "1D"
    keep_local = "2W"
    start_at   = "2025-08-01 09:00:00"
    }, {
    every      = "1D"
    keep_local = "8D"
    start_at   = "2025-08-01 09:00:00"
  }]
}
