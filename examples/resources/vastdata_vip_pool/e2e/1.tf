
resource "vastdata_vip_pool" "vastdb_vippool" {
  name        = "vastdb_vippool"
  role        = "PROTOCOLS"
  subnet_cidr = "24"

  ip_ranges = [
    ["16.0.0.6", "16.0.0.10"],
    ["16.0.0.20", "16.0.0.40"]
  ]
}