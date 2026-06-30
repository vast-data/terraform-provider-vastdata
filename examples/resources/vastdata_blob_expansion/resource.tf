resource "vastdata_blob_expansion" "example" {
  database_name    = "my-kafka-bucket"
  table_name       = "my-kafka-topic"
  expansion_format = "json"
}
