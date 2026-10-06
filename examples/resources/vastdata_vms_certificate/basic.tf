
resource "vastdata_vms_certificate" "vms" {
  vms_id      = 1
  certificate = file("path/to/tls.crt")
  private_key = file("path/to/tls.key")
}
