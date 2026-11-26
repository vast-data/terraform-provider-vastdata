# Count functions active within a time window.
# Both `start_time` and `end_time` are required and must be epoch microseconds
# (e.g., 1 hour = 3600000000 microseconds).
data "vastdata_de_count_functions" "example" {
  tenant_name = "de-lab"
  start_time  = 1700000000000000
  end_time    = 1700003600000000
}

output "function_count" {
  value = data.vastdata_de_count_functions.example.function_count
}
