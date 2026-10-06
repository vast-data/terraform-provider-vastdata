
resource "vastdata_s3_certificate" "s3" {
  cluster_id  = 1
  certificate = file("path/to/tls.crt")
  private_key = file("path/to/tls.key")
}
