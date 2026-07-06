data "vastdata_tenant_metric_label_values" "env_label" {
  tenant_id = vastdata_tenant.example.id
  label_id  = vastdata_tenant_metric_labels.environment.id
}
