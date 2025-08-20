resource "vastdata_tenant" "vastdb_tenant" {
  name             = "vastdbtenant-for-metrics"
  client_ip_ranges = [["192.168.0.11", "192.168.0.12"]]
  force_delete     = true
}

resource "vastdata_user" "vastdb_user_for_metrics" {
  name = "vastdb_user_for_metrics"
  uid  = 997160
}

resource "vastdata_view_policy" "vastdb_view_policy_for_metrics" {
  name                 = "vastdb_view_policy_for_metrics"
  tenant_id            = vastdata_tenant.vastdb_tenant.id
  flavor               = "S3_NATIVE"
  is_s3_default_policy = true
}

resource "vastdata_tenant_client_metrics" "client_metrics1" {
  depends_on = [vastdata_view_policy.vastdb_view_policy_for_metrics]
  tenant_id  = vastdata_tenant.vastdb_tenant.id
  config = {
    enabled            = true
    bucket_name        = "vastdb-metrics-dist"
    bucket_owner       = vastdata_user.vastdb_user_for_metrics.name
    max_capacity_mb    = 1024
    retention_time_sec = 86400
  }
  user_defined_columns = [
    {
      name = "ENV_USER_ID"
      field = {
        column_type = "string"
      }
    }
  ]
}
