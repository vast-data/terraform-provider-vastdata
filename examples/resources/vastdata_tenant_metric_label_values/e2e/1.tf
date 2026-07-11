resource "vastdata_tenant_metric_labels" "env" {
  key           = "vastdb_environment"
  default_value = "staging"
  description   = "Deployment environment"
}

resource "vastdata_tenant" "vastdb_tenant" {
  name         = "vastdb_test_mlv"
  force_delete = true
}

resource "vastdata_tenant_metric_label_values" "env_val" {
  tenant_id = vastdata_tenant.vastdb_tenant.id
  label_id  = vastdata_tenant_metric_labels.env.id
  value     = "production"
}
