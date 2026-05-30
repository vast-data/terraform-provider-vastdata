# Allocate IPs automatically using ips_count (VAST 5.5+).
resource "vastdata_vip_pool" "vastdb_vippool_allocated" {
  name      = "vastdb_vippool_allocated"
  role      = "PROTOCOLS"
  ips_count = 3
}
