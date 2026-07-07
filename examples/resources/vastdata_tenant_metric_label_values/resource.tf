resource "vastdata_tenant_metric_label_values" "example" {
  # The tenant to associate this label value with.
  tenant_id = vastdata_tenant.example.id

  # The metric label (from vastdata_tenant_metric_labels) whose value is set.
  label_id = vastdata_tenant_metric_labels.environment.id

  # The value for this metric label on this tenant.
  value = "production"
}
