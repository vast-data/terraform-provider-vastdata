# ignore:e2e
# Auto-allocate IPs scoped to a specific tenant (VAST 5.5+).
data "vastdata_tenant" "vastdb_tenant" {
  name = "default"
}

resource "vastdata_vip_pool" "vastdb_vippool_allocated" {
  name      = "vastdb_vippool_allocated"
  role      = "PROTOCOLS"
  ips_count = 4
  tenant_id = data.vastdata_tenant.vastdb_tenant.id
}
