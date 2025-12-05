resource "vastdata_tenant" "vastdb_tenant" {
  name             = "vastdbtenant"
  client_ip_ranges = [["192.168.0.50", "192.168.0.51"]]
}


resource "vastdata_protection_policy" "vastdb_ppolicy" {
  name           = "vastdb_ppolicy"
  indestructible = "false"
  prefix         = "vastdb-policy"
  clone_type     = "LOCAL"
  tenant_id      = vastdata_tenant.vastdb_tenant.id
  frames = [{
    every      = "1D"
    keep_local = "14D"
    start_at   = "2025-08-01 09:00:00"
    }, {
    every      = "1D"
    keep_local = "8D"
    start_at   = "2025-08-01 09:00:00"
  }]
}


resource "vastdata_protected_path" "protected_path" {
  name                 = "vastdbppath"
  source_dir           = "/"
  tenant_id            = vastdata_tenant.vastdb_tenant.id
  target_exported_dir  = "/view1"
  protection_policy_id = vastdata_protection_policy.vastdb_ppolicy.id
  enabled              = false
  capabilities         = "ASYNC_REPLICATION"
}
