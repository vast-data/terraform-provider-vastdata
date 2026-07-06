data "vastdata_tenant_metric_labels" "by_id" {
  id = 1
}

data "vastdata_tenant_metric_labels" "by_key" {
  key = "environment"
}
