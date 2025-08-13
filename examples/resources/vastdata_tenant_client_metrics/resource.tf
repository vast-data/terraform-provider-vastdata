resource "vastdata_tenant_client_metrics" "vastdb_tenant_client_metrics" {
  tenant_id = 1
  config = {
    enabled            = true
    max_capacity_mb    = 1024
    retention_time_sec = 86400
    bucket_owner       = "admin"
    bucket_name        = "client-metrics-bucket"
  }
}

# ---------------------
# Complete examples
# ---------------------

data "vastdata_tenant" "vastdb_tenant" {
  name = "default"
}

resource "vastdata_tenant_client_metrics" "vastdb_tenant_client_metrics" {
  tenant_id = data.vastdata_tenant.vastdb_tenant.id

  config = {
    enabled            = true
    max_capacity_mb    = 2048
    retention_time_sec = 172800
    bucket_owner       = "metrics-user"
    bucket_name        = "vastdb-metrics"
  }

  user_defined_columns = [
    {
      name = "ENV_USER_ID"
      field = {
        column_type = "string"
      }
    },
    {
      name = "ENV_ACCESS_COUNT"
      field = {
        column_type = "int16"
      }
    }
  ]
}

# --------------------

data "vastdata_user" "vastdb_user" {
  name = "runner"
}

data "vastdata_tenant" "vastdb_tenant" {
  name = "default"
}

resource "vastdata_tenant_client_metrics" "vastdb_tenant_client_metrics" {
  tenant_id = data.vastdata_tenant.vastdb_tenant.id
  config = {
    enabled            = true
    bucket_owner       = data.vastdata_user.vastdb_user.name
    bucket_name        = "vastdb-metrics"
    max_capacity_mb    = 2048
    retention_time_sec = 172800
  }
}

# --------------------

