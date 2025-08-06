
resource "vastdata_vip_pool" "vastdb_vippool" {
  name        = "vastdb_vippool"
  role        = "PROTOCOLS"
  subnet_cidr = "24"

  ip_ranges = [
    ["11.0.0.6", "11.0.0.10"],
  ]
}

resource "vastdata_tenant" "vastdb_tenant" {
  name                 = "vastdbtenant"
  allow_locked_users   = true
  allow_disabled_users = true
  access_ip_ranges     = ["11.0.0.6", "11.0.0.7"]
}