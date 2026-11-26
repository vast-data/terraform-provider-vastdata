# Query recent execution traces for a Data Engine pipeline.
# Timestamps (`start_time`, `end_time`, `start_time_after`, `start_time_before`) are in epoch microseconds.
# Use `pagination.next_cursor` / `pagination.previous_cursor` with the `cursor` argument to page through results.
data "vastdata_de_traces" "example" {
  tenant_name   = "de-lab"
  pipeline_name = "my-pipeline"
  limit         = 50
  order_by      = "desc"
}

output "traces" {
  value = data.vastdata_de_traces.example.data
}

output "traces_pagination" {
  value = data.vastdata_de_traces.example.pagination
}
