data "vastdata_cluster" "vastdb_cluster" {
  name = "loop101-kfs"
}

resource "vastdata_cluster_ekm" "vastdb_cluster_ekm" {
  cluster_id      = data.vastdata_cluster.vastdb_cluster.id
  ekm_servers     = "10.27.14.107:5696,10.27.14.108:5696"
  encryption_type = "CIPHER_TRUST_KMIP"

  # Certificate Configuration
  ekm_certificate = <<-EOT
    -----BEGIN CERTIFICATE-----
    .
    .  <content>
    .
-----END CERTIFICATE-----
  EOT

  ekm_private_key = <<-EOT
    -----BEGIN PRIVATE KEY-----
    .
    .  <content>
    .
-----END PRIVATE KEY-----
  EOT

  ekm_ca_certificate = <<-EOT
    -----BEGIN CERTIFICATE-----
    .
    .  <content>
    .
-----END CERTIFICATE-----
  EOT

  # Thales-specific Configuration
  ekm_auth_domain   = "vastdata-auth"
  ekm_domain        = "vastdata-domain"
  ekm_proxy_address = "https://10.27.14.109:8443"

  # Security Configuration
  ekm_bypass_validation = false
}

