resource "vastdata_tenant_metric_label_values_bulk" "example" {
  tenant_id = vastdata_tenant.example.id
  values = {
    environment = "production"
    region      = "us-east-1"
  }
}
