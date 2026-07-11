resource "vastdata_tenant_metric_labels" "env_label" {
  key           = "vastdb_environment"
  default_value = "production"
  description   = "Deployment environment tag for tenant metrics."
}
