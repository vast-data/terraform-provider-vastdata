# Basic Kerberos configuration
resource "vastdata_kerberos" "basic_kerberos" {
  realm              = "EXAMPLE.COM"
  service_principals = ["nfs/vastcluster.example.com"]
  kdc                = ["kdc.example.com"]
  kadmin_servers     = ["kadmin.example.com"]
  enabled            = true
}

