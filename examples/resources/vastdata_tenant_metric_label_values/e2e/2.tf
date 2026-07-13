resource "vastdata_tenant" "vastdb_tenant" {
  name         = "vastdb-test-mlv"
  force_delete = true
}

resource "vastdata_tenant_metric_labels" "vastdb_metric_label" {
  key           = "vastdb_environment"
  default_value = "production"
  description   = "Deployment environment tag for tenant metrics."
}

resource "vastdata_tenant_metric_label_values" "vastdb_metric_label_val" {
  tenant_id = vastdata_tenant.vastdb_tenant.id
  label_id  = vastdata_tenant_metric_labels.vastdb_metric_label.id
  value     = "production"
}
