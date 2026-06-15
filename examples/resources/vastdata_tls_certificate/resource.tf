# TLS certificate with CA only (no CRL)
resource "vastdata_tls_certificate" "vastdb_tls_cert" {
  ca_certificate_file = file("path/to/ca.crt")
  ca_certificate_name = "my-ca"
  protocols           = ["NFS"]
}

# TLS certificate with CA and CRL revocation list
resource "vastdata_tls_certificate" "vastdb_tls_cert_with_crl" {
  ca_certificate_file = file("path/to/ca.crt")
  revocation_file     = file("path/to/ca.crl")
  ca_certificate_name = "my-ca-with-crl"
  revocations_name    = "my-crl"
  protocols           = ["NFS", "KAFKA"]
}
