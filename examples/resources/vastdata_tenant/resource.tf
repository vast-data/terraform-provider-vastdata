
resource "vastdata_tenant" "vastdb_tenant" {
  name = "vastdbtenant"
}

# ---------------------
# Complete examples
# ---------------------


resource "vastdata_tenant" "vastdb_tenant" {
  name = "vastdbtenant"
  client_ip_ranges = [
    ["12.0.0.6", "12.0.0.10"],
    ["12.0.0.20", "12.0.0.40"],
    ["192.168.0.100", "192.168.0.201"],
  ]
}

# --------------------


resource "vastdata_vip_pool" "vastdb_vippool" {
  name        = "vastdb_vippool"
  role        = "PROTOCOLS"
  subnet_cidr = "24"

  ip_ranges = [
    ["14.0.0.6", "14.0.0.10"],
  ]
}

resource "vastdata_tenant" "vastdb_tenant" {
  name                 = "vastdbtenant"
  allow_locked_users   = true
  allow_disabled_users = true
  access_ip_ranges     = ["14.0.0.6", "14.0.0.7"]
}

# --------------------

