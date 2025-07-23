resource "vastdata_administrator_manager" "vastdb_manager" {
  password_expiration_disabled = true
  username                     = "vastdb_manager"
  password                     = "Www##12345678"
  roles = [
    1
  ]
}
