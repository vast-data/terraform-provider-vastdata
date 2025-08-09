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
