# ignore:e2e - requires live CipherTrust/KMIP EKM servers and matching PEMs
resource "vastdata_cluster_ekm" "vastdb_cluster_ekm" {
  cluster_id         = 1
  ekm_servers        = "127.0.0.1:5696"
  encryption_type    = "CIPHER_TRUST_KMIP"
  ekm_certificate    = "PLACEHOLDER"
  ekm_private_key    = "PLACEHOLDER"
  ekm_ca_certificate = "PLACEHOLDER"
}
