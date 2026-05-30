# Allocate IPs automatically with additional configuration options (VAST 5.5+).
resource "vastdata_vip_pool" "vastdb_vippool_allocated" {
  name                      = "vastdb_vippool_allocated"
  role                      = "PROTOCOLS"
  ips_count                 = 5
  enable_weighted_balancing = true
  domain_name               = "vastdb.example.com"
  vms_preferred             = true
}
