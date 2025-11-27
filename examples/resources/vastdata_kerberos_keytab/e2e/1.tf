# ignore:e2e

# Basic e2e test: Generate keytab only
resource "vastdata_kerberos" "vastdb_kerberos" {
  realm              = "VASTDB.LOCAL"
  service_principals = ["nfs/vastdb.local"]
  kdc                = ["10.27.14.107"]
  kadmin_servers     = ["10.27.14.107"]
  enabled            = true
}

resource "vastdata_kerberos_keytab" "vastdb_keytab" {
  kerberos_id    = vastdata_kerberos.vastdb_kerberos.id
  admin_username = "admin@VASTDB.LOCAL"
  admin_password = "vastdata123"
}

