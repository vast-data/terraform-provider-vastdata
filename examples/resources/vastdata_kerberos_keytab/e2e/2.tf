# ignore:e2e

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

