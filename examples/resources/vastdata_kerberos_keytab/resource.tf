# Basic keytab generation
resource "vastdata_kerberos_keytab" "generate_keytab" {
  kerberos_id    = 1
  admin_username = "admin@EXAMPLE.COM"
  admin_password = "admin_password"
}

# ---------------------
# Complete examples
# ---------------------


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


# --------------------


# Generate keytab and upload custom keytab file
resource "vastdata_kerberos" "vastdb_kerberos_with_upload" {
  realm = "VASTENG.LAB"
  service_principals = [
    "nfs/vast-cluster.vasteng.lab",
    "cifs/vast-cluster.vasteng.lab"
  ]
  kdc            = ["kdc.vasteng.lab"]
  kadmin_servers = ["kadmin.vasteng.lab"]
  enabled        = true
}

resource "vastdata_kerberos_keytab" "vastdb_keytab_with_upload" {
  kerberos_id    = vastdata_kerberos.vastdb_kerberos_with_upload.id
  admin_username = "Administrator@VASTENG.LAB"
  admin_password = "VastData123!"

  # Upload a keytab file using base64 encoded test data
  keytab_file = base64decode("BQIAAABUAAIAC0VYQU1QTEUuQ09NAAR2YXN0ABB0ZXN0LmV4YW1wbGUuY29tAAAAAWHPmYABABIAIAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
  filename    = "vast-test.keytab"
}


# --------------------


# Multiple environments with different keytab strategies
resource "vastdata_kerberos" "production_kerberos" {
  realm = "PROD.COMPANY.COM"
  service_principals = [
    "nfs/vast-prod.prod.company.com",
    "cifs/vast-prod.prod.company.com"
  ]
  kdc            = ["kdc1.prod.company.com", "kdc2.prod.company.com"]
  kadmin_servers = ["kadmin.prod.company.com"]
  enabled        = true
}

resource "vastdata_kerberos" "staging_kerberos" {
  realm = "STAGING.COMPANY.COM"
  service_principals = [
    "nfs/vast-staging.staging.company.com"
  ]
  kdc            = ["kdc.staging.company.com"]
  kadmin_servers = ["kadmin.staging.company.com"]
  enabled        = true
}

# Production environment - generate keytab only
resource "vastdata_kerberos_keytab" "production_keytab" {
  kerberos_id    = vastdata_kerberos.production_kerberos.id
  admin_username = "svc-vast@PROD.COMPANY.COM"
  admin_password = "ProdPassword123!"
}

# Staging environment - upload pre-generated keytab
resource "vastdata_kerberos_keytab" "staging_keytab" {
  kerberos_id    = vastdata_kerberos.staging_kerberos.id
  admin_username = "admin@STAGING.COMPANY.COM"
  admin_password = "StagingPassword123!"

  # Upload staging keytab using base64 test data
  keytab_file = base64decode("BQIAAABUAAIAC0VYQU1QTEUuQ09NAAR2YXN0ABB0ZXN0LmV4YW1wbGUuY29tAAAAAWHPmYABABIAIAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
  filename    = "vast-staging.keytab"
}


# --------------------

