resource "vastdata_bgp_config" "vastdb_bgp_config" {
  name         = "bgp-config-1"
  self_asn     = 65001
  external_asn = 65002
  method       = "numbered"
}

