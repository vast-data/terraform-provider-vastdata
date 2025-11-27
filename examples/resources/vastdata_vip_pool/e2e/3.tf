# ignore:e2e

data "vastdata_tenant" "vastdb_tenant" {
  name = "default"
}

resource "vastdata_vip_pool" "vastdb_vippool" {
  name          = "vastdb_vippool"
  role          = "PROTOCOLS"
  tenant_id     = data.vastdata_tenant.vastdb_tenant.id
  domain_name   = "vastdb.example.com"
  vms_preferred = true
  subnet_cidr   = "24"
  ip_ranges = [
    ["16.0.0.50", "16.0.0.80"],
  ]
}