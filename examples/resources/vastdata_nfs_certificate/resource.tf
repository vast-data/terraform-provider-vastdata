
resource "vastdata_nfs_certificate" "nfs" {
  cluster_id  = 1
  certificate = file("path/to/tls.crt")
  private_key = file("path/to/tls.key")
}

# ---------------------
# Complete examples
# ---------------------

# e2e: install NFSv4 RPC-with-TLS certificate (write-only PEMs + certificate_md5).
# Not importable; destroy is a no-op.
resource "vastdata_nfs_certificate" "vastdb_nfs_certificate" {
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

