# Basic keytab generation
resource "vastdata_kerberos_keytab" "generate_keytab" {
  kerberos_id    = 1
  admin_username = "admin@EXAMPLE.COM"
  admin_password = "admin_password"
}
