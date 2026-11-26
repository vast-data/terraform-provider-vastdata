# Inspect the root trace and function spans for a specific `trace_id`.
# Can be chained directly from `data.vastdata_de_traces` to inspect the latest execution.
data "vastdata_de_traces" "recent" {
  tenant_name   = "de-lab"
  pipeline_name = "my-pipeline"
  limit         = 1
  order_by      = "desc"
}

data "vastdata_de_trace_tree" "example" {
  tenant_name = "de-lab"
  trace_id    = data.vastdata_de_traces.recent.data[0].trace_id
}

output "root_trace" {
  value = data.vastdata_de_trace_tree.example.root_trace
}

output "trace_functions" {
  value = data.vastdata_de_trace_tree.example.functions
}
