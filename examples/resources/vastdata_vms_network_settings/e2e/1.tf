data "vastdata_vms" "vastdb_vms" {
  name = "VMS"
}

resource "vastdata_vms_network_settings" "vastdb_vms_network_settings_full" {
  vms_id               = data.vastdata_vms.vastdb_vms.id
  auto_ports_ext_iface = "outband"
  b2b_ipmi             = false
  dns                  = ["8.8.8.8", "8.8.4.4"]
  ntp                  = ["pool.ntp.org", "time.google.com"]
  eth_mtu              = 9000
  ext_netmask          = "255.255.0.0"
  external_gateways    = ["10.0.0.1"]
  ib_mtu               = 65520
  management_vips      = ["10.0.0.10"]

  hosts = [
    {
      id       = 1
      hostname = "cnode-1"
      mgmt_ip  = "10.0.0.11"
      ipmi_ip  = "10.0.1.11"
    },
    {
      id       = 2
      hostname = "cnode-2"
      mgmt_ip  = "10.0.0.12"
      ipmi_ip  = "10.0.1.12"
    },
  ]
}
