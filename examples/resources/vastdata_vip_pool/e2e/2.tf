
resource "vastdata_vip_pool" "vastdb_vippool" {
  name                      = "vastdb_vippool"
  role                      = "PROTOCOLS"
  subnet_cidr               = "24"
  enable_weighted_balancing = true
  ip_ranges = [
    ["16.0.0.50", "16.0.0.80"],
  ]
}