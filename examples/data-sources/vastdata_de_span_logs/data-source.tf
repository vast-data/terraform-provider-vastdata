# Fetch ordered log records belonging to a specific trace (and optionally a specific span).
# `functions` on `vastdata_de_trace_tree` is a set, so convert with `tolist(...)` to index a span.
data "vastdata_de_traces" "recent" {
  tenant_name   = "de-lab"
  pipeline_name = "my-pipeline"
  limit         = 1
  order_by      = "desc"
}

data "vastdata_de_trace_tree" "tree" {
  tenant_name = "de-lab"
  trace_id    = data.vastdata_de_traces.recent.data[0].trace_id
}

data "vastdata_de_span_logs" "example" {
  tenant_name = "de-lab"
  trace_id    = data.vastdata_de_traces.recent.data[0].trace_id
  span_id     = tolist(data.vastdata_de_trace_tree.tree.functions)[0].value.span_id
}

output "span_logs" {
  value = data.vastdata_de_span_logs.example.data
}

output "span_logs_pagination" {
  value = data.vastdata_de_span_logs.example.pagination
}
