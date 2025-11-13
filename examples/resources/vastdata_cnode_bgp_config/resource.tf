resource "vastdata_cnode_bgp_config" "vastdb_cnode_bgp_config" {
  cnode_id           = 1
  enabled            = true
  port1_peer_address = "10.0.1.1"
  port1_self_address = "10.0.1.2"
  port2_peer_address = "10.0.2.1"
  port2_self_address = "10.0.2.2"
  self_asn           = "65001"
  subnet_bits        = 24
}


# ---------------------
# Complete example with data source
# ---------------------

data "vastdata_cnode" "existing_cnode" {
  id = 1
}

resource "vastdata_cnode_bgp_config" "vastdb_cnode_bgp_config_with_datasource" {
  cnode_id           = data.vastdata_cnode.existing_cnode.id
  enabled            = true
  port1_peer_address = "192.168.1.1"
  port1_self_address = "192.168.1.2"
  port2_peer_address = "192.168.2.1"
  port2_self_address = "192.168.2.2"
  self_asn           = "65002"
  subnet_bits        = 30
}
