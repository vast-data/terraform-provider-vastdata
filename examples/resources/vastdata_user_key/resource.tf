
resource "vastdata_user_key" "vastdb_user_key" {
  username  = "example-user"
  tenant_id = "example-tenant-id"
}

# ---------------------
# Complete examples
# ---------------------


data "vastdata_tenant" "vastdb_tenant" {
  name = "default"
}

resource "vastdata_user" "vastdb_user" {
  name = "vastdb_user"
  uid  = 30109
}

resource "vastdata_user_key" "vastdb_user_key" {
  username  = vastdata_user.vastdb_user.name
  tenant_id = data.vastdata_tenant.vastdb_tenant.id
  enabled   = false
}

# --------------------



resource "vastdata_user" "vastdb_user" {
  name = "vastdb_user"
  uid  = 30109
}

resource "vastdata_user_key" "vastdb_user_key" {
  user_id = vastdata_user.vastdb_user.id

  pgp_public_key = <<-EOT
    -----BEGIN PGP PUBLIC KEY BLOCK-----
    .
    .  <content>
    .
-----END PGP PUBLIC KEY BLOCK-----
  EOT

}


# --------------------

