resource "vastdata_blob_expansion" "blob_exp" {
  database_name       = "tfmild-bucket"
  table_name          = "tfcautious-kowari"
  target_table_name   = "tfindustrious-goose"
  target_table_schema = "kafka_topics"
  expansion_format    = "json"
  tenant_id           = 1
  arrow_schema = [
    {
      name  = "label"
      field = { column_type = "string" }
    },
    {
      name  = "timestamp"
      field = { column_type = "int64" }
    },
    {
      name  = "value"
      field = { column_type = "float" }
    },
  ]
  copy_source_column          = false
  add_missing_values_output   = false
  add_excessive_values_output = false
}
