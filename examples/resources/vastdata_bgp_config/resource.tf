resource "vastdata_bgp_config" "vastdb_bgp_config" {
  name         = "bgp-config-1"
  self_asn     = 65001
  external_asn = 65002
  method       = "numbered"
}


# ---------------------
# Complete examples
# ---------------------

resource "vastdata_bgp_config" "vastdb_bgp_config" {
  name                           = "bgp-config-1"
  self_asn                       = 65001
  external_asn                   = 65002
  any_external_asn               = false
  bfd_enabled                    = true
  method                         = "numbered"
  subnet_bits                    = 32
  vip_migration_grace_period_sec = 30
}


# --------------------

