
resource "vastdata_tenant" "vastdb_tenant" {
  name = "vastdbtenant"
  client_ip_ranges = [
    ["192.168.0.100", "192.168.0.201"],
    ["11.0.0.6", "11.0.0.10"],
    ["11.0.0.20", "11.0.0.40"]
  ]
}
