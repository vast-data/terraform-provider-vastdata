# NOTE: All metrics in a comma-separated `metric_names` query must share the same
# `metric_type` (e.g., all Counter or all Histogram). Mixing types returns HTTP 422.
# Use `data.vastdata_de_metric_names` to discover available metric names and their types.
data "vastdata_de_metrics" "example" {
  tenant_name  = "de-lab"
  metric_names = "dataengine.function.invocations,dataengine.event.total"
  pipeline     = "my-pipeline"
  aggregation  = "avg"
  timeframe    = "1h"
  max_points   = 60
}

output "metrics_query" {
  value = data.vastdata_de_metrics.example.query
}

output "metrics_results" {
  value = data.vastdata_de_metrics.example.results
}
