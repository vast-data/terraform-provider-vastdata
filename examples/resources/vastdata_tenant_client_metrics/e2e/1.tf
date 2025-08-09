# data "vastdata_tenant" "vastdb_tenant" {
#   name = "default"
# }
#
# resource "vastdata_tenant_client_metrics" "vastdb_tenant_client_metrics" {
#   tenant_id = data.vastdata_tenant.vastdb_tenant.id
#   config = {
#     enabled            = true
#     max_capacity_mb    = 2048
#     retention_time_sec = 172800
#     bucket_owner       = "metrics-user"
#     bucket_name        = "tenant-metrics-bucket"
#   }
#   user_defined_columns = [
#     {
#       name  = "ENV_USER_ID"
#       field = "string"
#     },
#     {
#       name  = "ENV_ACCESS_COUNT"
#       field = "integer"
#     }
#   ]
# }