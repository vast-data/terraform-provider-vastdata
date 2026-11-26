# Discover available metrics and their metadata (`name`, `type`, `unit`, `scope`, `description`).
# Useful before querying `vastdata_de_metrics`, which requires all metrics in a single query
# to have the same `type` (e.g., all `Counter` or all `Histogram`).
data "vastdata_de_metric_names" "example" {
  tenant_name = "de-lab"
  pipeline    = "my-pipeline"
}

output "metric_names" {
  value = data.vastdata_de_metric_names.example.items
}
