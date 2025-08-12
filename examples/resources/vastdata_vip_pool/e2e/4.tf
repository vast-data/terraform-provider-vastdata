# create IPv6 VIP pool
resource "vastdata_vip_pool" "vastdata_vip_pool_ipv6" {
  name             = "vastdata_vip_pool_ipv6"
  role             = "PROTOCOLS"
  subnet_cidr_ipv6 = 64
  ip_ranges        = [["fec0:10::11", "fec0:10::18"]]
}