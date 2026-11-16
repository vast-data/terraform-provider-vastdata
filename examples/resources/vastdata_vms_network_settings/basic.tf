# Partial config is enough: the provider GETs current settings and fills any
# omitted PATCH fields before calling VMS (VMS rejects true partial updates).
resource "vastdata_vms_network_settings" "vastdb_vms_network_settings" {
  vms_id = 1
  dns    = ["8.8.8.8", "1.1.1.1"]
  ntp    = ["time.google.com"]
}
