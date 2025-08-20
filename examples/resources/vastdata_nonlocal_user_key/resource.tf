resource "vastdata_nonlocal_user_key" "vastdb_nonlocal_user_key" {
  uid       = 1000
  tenant_id = 1
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


resource "vastdata_nonlocal_user_key" "vastdb_nonlocal_user_key" {
  uid       = vastdata_user.vastdb_user.uid
  tenant_id = data.vastdata_tenant.vastdb_tenant.id
  enabled   = false
}

# --------------------


data "vastdata_tenant" "vastdb_tenant" {
  name = "default"
}


resource "vastdata_user" "vastdb_user" {
  name = "vastdb_user"
  uid  = 30109
}


resource "vastdata_nonlocal_user_key" "vastdb_nonlocal_user_key" {
  username = vastdata_user.vastdb_user.name

  pgp_public_key = <<-EOT
    -----BEGIN PGP PUBLIC KEY BLOCK-----
    .
    .  <content>
    .
-----END PGP PUBLIC KEY BLOCK-----
  EOT

}

# --------------------


data "vastdata_tenant" "vastdb_tenant" {
  name = "default"
}


resource "vastdata_user" "vastdb_user" {
  name = "vastdb_user"
  uid  = 30109
}


resource "vastdata_nonlocal_user_key" "vastdb_nonlocal_user_key" {
  username = vastdata_user.vastdb_user.name
}

# --------------------

