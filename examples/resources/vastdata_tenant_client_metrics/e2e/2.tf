data "vastdata_tenant" "vastdb_tenant" {
  name = "default"
}

resource "vastdata_tenant_client_metrics" "vastdb_tenant_client_metrics" {
  tenant_id = data.vastdata_tenant.vastdb_tenant.id
  config = {
    enabled = false
  }
} 