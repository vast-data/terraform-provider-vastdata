resource "vastdata_rack" "vastdb_rack" {
  name = "rack-1"
}

resource "vastdata_rack_bgp_config" "vastdb_rack_bgp_config" {
  rack_id       = vastdata_rack.vastdb_rack.id
  ip_ranges     = ["10.0.0.1", "10.0.0.10"]
  ips_represent = "odd"
  self_asn      = "65001"
  subnet_bits   = 24
}
