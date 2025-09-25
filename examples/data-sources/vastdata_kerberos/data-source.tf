data "vastdata_kerberos" "vastdb_kerberos_by_id" {
  id = 1
}

data "vastdata_kerberos" "vastdb_kerberos_by_realm" {
  realm = "EXAMPLE.COM"
}
