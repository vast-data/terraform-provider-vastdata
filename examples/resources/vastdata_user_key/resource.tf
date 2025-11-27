
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
  name              = "vastdb_user"
  uid               = 30109
  local_provider_id = 1
}

resource "vastdata_user_key" "vastdb_user_key" {
  username  = vastdata_user.vastdb_user.name
  tenant_id = data.vastdata_tenant.vastdb_tenant.id
  enabled   = false
}

# --------------------

data "vastdata_tenant" "vastdb_tenant" {
  name = "default"
}

resource "vastdata_user" "vastdb_user" {
  name              = "vastdb_user"
  uid               = 5001
  local_provider_id = 1
}

resource "vastdata_user_key" "vastdb_user_key" {
  user_id   = vastdata_user.vastdb_user.id
  tenant_id = data.vastdata_tenant.vastdb_tenant.id

  pgp_public_key = <<-EOT
    -----BEGIN PGP PUBLIC KEY BLOCK-----
    .
    .  <content>
    .
-----END PGP PUBLIC KEY BLOCK-----
  EOT

}


# --------------------

resource "vastdata_local_provider" "vastdb_local_provider" {
  name       = "vastdb_local_provider"
  managed_by = ["TENANT_ADMIN", "SUPER_ADMIN"]
}

resource "vastdata_tenant" "vastdb_tenant" {
  name              = "vastdb-tenant"
  client_ip_ranges  = [["192.168.100.1", "192.168.100.10"]]
  local_provider_id = vastdata_local_provider.vastdb_local_provider.id
}

resource "vastdata_user" "vastdb_user" {
  name              = "vastdb_user"
  uid               = 5001
  local_provider_id = vastdata_local_provider.vastdb_local_provider.id
}

resource "vastdata_user_key" "vastdb_user_key" {
  user_id   = vastdata_user.vastdb_user.id
  tenant_id = vastdata_tenant.vastdb_tenant.id
  enabled   = false
}

# --------------------

