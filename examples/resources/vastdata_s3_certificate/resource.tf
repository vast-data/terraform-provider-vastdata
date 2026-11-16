
resource "vastdata_s3_certificate" "s3" {
  cluster_id  = 1
  certificate = file("path/to/tls.crt")
  private_key = file("path/to/tls.key")
}

# ---------------------
# Complete examples
# ---------------------

# e2e: install S3 HTTPS certificate (write-only PEMs + certificate_md5).
# Not importable; destroy is a no-op.
resource "vastdata_s3_certificate" "vastdb_s3_certificate" {
  cluster_id  = 1
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

