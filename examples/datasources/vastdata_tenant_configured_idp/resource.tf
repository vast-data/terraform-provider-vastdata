data "vastdata_tenant_configured_idp" "vastdb_tenant_configured_idp" {
  name = "default"
}

# ---------------------
# Complete examples
# ---------------------

data "vastdata_tenant" "vastdb_tenant" {
  name = "default"
}

data "vastdata_tenant_configured_idp" "vastdb_tenant_configured_idp" {
  name = data.vastdata_tenant.vastdb_tenant.name
}

# --------------------

data "vastdata_tenant_configured_idp" "vastdb_tenant_configured_idp" {
  name = "my-tenant"
}

# -------------------- 