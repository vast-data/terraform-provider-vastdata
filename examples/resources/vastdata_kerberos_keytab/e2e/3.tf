# ignore:e2e

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

