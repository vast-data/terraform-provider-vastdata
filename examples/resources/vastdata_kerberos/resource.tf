# Basic Kerberos configuration
resource "vastdata_kerberos" "basic_kerberos" {
  realm              = "EXAMPLE.COM"
  service_principals = ["nfs/vastcluster.example.com"]
  kdc                = ["kdc.example.com"]
  kadmin_servers     = ["kadmin.example.com"]
  enabled            = true
}

# Advanced Kerberos configuration with multiple service principals
resource "vastdata_kerberos" "advanced_kerberos" {
  realm = "CORP.EXAMPLE.COM"
  service_principals = [
    "nfs/vastcluster.corp.example.com",
    "cifs/vastcluster.corp.example.com",
    "http/vastcluster.corp.example.com"
  ]
  kdc = [
    "kdc1.corp.example.com",
    "kdc2.corp.example.com"
  ]
  kadmin_servers = [
    "kadmin1.corp.example.com",
    "kadmin2.corp.example.com"
  ]
  enabled = true
}

# Kerberos configuration with IP addresses
resource "vastdata_kerberos" "ip_based_kerberos" {
  realm              = "VASTDB.LOCAL"
  service_principals = ["nfs/vastdb"]
  kdc                = ["10.27.14.107"]
  kadmin_servers     = ["10.27.14.107"]
  enabled            = true
}

