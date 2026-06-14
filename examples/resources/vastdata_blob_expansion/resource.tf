resource "vastdata_blob_expansion" "example" {
  database_name     = "my-kafka-bucket"
  table_name        = "my-kafka-topic"
  expansion_format  = "json"
  target_table_name = "my-expanded-table"
  target_table_schema = "kafka_topics"

  arrow_schema = [
    {
      name = "user_id"
      field = {
        column_type = "int64"
      }
    },
    {
      name = "event_type"
      field = {
        column_type = "string"
      }
    },
    {
      name = "timestamp_ms"
      field = {
        column_type = "int64"
      }
    },
  ]

  copy_source_column          = false
  add_missing_values_output   = false
  add_excessive_values_output = false
}
