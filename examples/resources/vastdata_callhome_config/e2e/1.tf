# e2e: adopt Call Home config id=1 and set a non-secret field.
resource "vastdata_callhome_config" "vastdb_callhome_config" {
  id         = 1
  obfuscated = true
}
