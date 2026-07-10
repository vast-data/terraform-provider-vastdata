resource "vastdata_tenant_metric_label_values_bulk" "vastdb_tenant_metric_values_bulk" {
  # Tenant whose metric label values are managed in bulk.
  tenant_id = vastdata_tenant.example.id

  # Label keys must match vastdata_tenant_metric_labels.key values.
  values = {
    environment = "production"
    region      = "us-east-1"
  }
}
