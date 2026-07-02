resource "vastdata_tenant_metric_labels" "env" {
  key           = "environment"
  default_value = "staging"
  description   = "Deployment environment"
}

resource "vastdata_tenant" "test" {
  name = "test-mlv"
}

resource "vastdata_tenant_metric_label_values" "env_val" {
  tenant_id = vastdata_tenant.test.id
  label_id  = vastdata_tenant_metric_labels.env.id
  value     = "production"
}
