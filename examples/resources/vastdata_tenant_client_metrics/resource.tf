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

# The tenant "metricstenant" should be created manually in the VAST UI or CLI 
# with a default S3 policy assigned to it before running this configuration.
# This is a requirement for client metrics to function properly.
#
# This workaround avoids the issue where terraform destroy fails with:
# "Cannot delete view policy that has referenced views"

data "vastdata_tenant" "vastdb_tenant" {
  name = "metricstenant"
}

data "vastdata_local_provider" "vastdb_local_provider" {
  name = "default"
}

resource "vastdata_user" "vastdb_user" {
  name              = "metricsuser"
  local_provider_id = data.vastdata_local_provider.vastdb_local_provider.id
}

resource "vastdata_tenant_client_metrics" "vastdb_tenant_client_metrics" {
  tenant_id = data.vastdata_tenant.vastdb_tenant.id

  config = {
    enabled            = true
    max_capacity_mb    = 2048
    retention_time_sec = 172800
    bucket_owner       = vastdata_user.vastdb_user.name
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
