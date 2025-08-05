resource "vastdata_api_token" "vastdb_api_token" {
  name        = "vastdb_api_token"
  expiry_date = "1Y"
}


# ---------------------
# Complete examples
# ---------------------

resource "vastdata_administrator_realm" "vastdb_realm" {
  name         = "vastdb_realm"
  object_types = ["nic", "viewpolicy"]
}

resource "vastdata_administrator_role" "vastdb_role" {
  name        = "vastdb_role"
  permissions = "view"
  realm       = vastdata_administrator_realm.vastdb_realm.id
}

resource "vastdata_administrator_manager" "vastdb_manager" {
  password_expiration_disabled = true
  username                     = "vastdb_manager"
  password                     = "Www##12345678"
  first_name                   = "me"
  last_name                    = "myself"
  roles = [
    vastdata_administrator_role.vastdb_role.id
  ]
}

resource "vastdata_api_token" "vastdb_api_token" {
  name        = "vastdb_api_token"
  expiry_date = "1Y"
  owner       = vastdata_administrator_manager.vastdb_manager.username

}

# --------------------

