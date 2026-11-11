
resource "vastdata_vms_certificate" "vms" {
  vms_id      = 1
  certificate = file("path/to/tls.crt")
  private_key = file("path/to/tls.key")
}

# ---------------------
# Complete examples
# ---------------------

# e2e: install VMS web/API TLS certificate (write-only PEMs + certificate_md5).
# Not importable; destroy is a no-op.
resource "vastdata_vms_certificate" "vastdb_vms_certificate" {
  vms_id      = 1
  certificate = <<-EOT
    -----BEGIN CERTIFICATE-----
    .
    .  <content>
    .
-----END CERTIFICATE-----
  EOT
  private_key = <<-EOT
    -----BEGIN PRIVATE KEY-----
    .
    .  <content>
    .
-----END PRIVATE KEY-----
  EOT
}

# --------------------

