
resource "vastdata_nfs_certificate" "nfs" {
  cluster_id  = 1
  certificate = file("path/to/tls.crt")
  private_key = file("path/to/tls.key")
}
