
resource "vastdata_vip_pool" "vastdb_vippool" {
  name        = "vastdb_vippool"
  role        = "PROTOCOLS"
  subnet_cidr = "24"

  ip_ranges = [
    ["11.0.0.6", "11.0.0.10"],
  ]
}

# ---------------------
# Complete examples
# ---------------------


resource "vastdata_vip_pool" "vastdb_vippool" {
  name        = "vastdb_vippool"
  role        = "PROTOCOLS"
  subnet_cidr = "24"

  ip_ranges = [
    ["11.0.0.6", "11.0.0.10"],
    ["11.0.0.20", "11.0.0.40"]
  ]
}

# --------------------


resource "vastdata_vip_pool" "vastdb_vippool" {
  name                      = "vastdb_vippool"
  role                      = "PROTOCOLS"
  subnet_cidr               = "24"
  enable_weighted_balancing = true
  ip_ranges = [
    ["11.0.0.50", "11.0.0.80"],
  ]
}

# --------------------


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
    ["11.0.0.50", "11.0.0.80"],
  ]
}

# --------------------

