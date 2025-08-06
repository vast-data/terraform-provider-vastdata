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