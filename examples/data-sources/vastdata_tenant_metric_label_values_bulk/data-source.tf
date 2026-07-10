data "vastdata_tenant" "example" {
  name = "default"
}

data "vastdata_tenant_metric_label_values_bulk" "example" {
  tenant_id = data.vastdata_tenant.example.id
}
