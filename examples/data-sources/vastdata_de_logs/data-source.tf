# Query function/pipeline logs (includes both `vast-runtime` and `user` scopes).
# Optional `message` filtering supports postfixes: `__iexact`, `__icontains`, `__istartswith`.
data "vastdata_de_logs" "example" {
  tenant_name   = "de-lab"
  function_name = "my-function"
  limit         = 100
  order_by      = "desc"
}

output "logs" {
  value = data.vastdata_de_logs.example.data
}

output "logs_pagination" {
  value = data.vastdata_de_logs.example.pagination
}
