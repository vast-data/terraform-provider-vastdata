
resource "vastdata_tenant" "vastdb_tenant" {
  name = "vastdbtenant"
  client_ip_ranges = [
    ["12.0.0.6", "12.0.0.10"],
    ["12.0.0.20", "12.0.0.40"],
    ["192.168.0.100", "192.168.0.201"],
  ]
}
