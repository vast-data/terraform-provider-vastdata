data "vastdata_blob_expansion" "example" {
  database_name      = "tfmild-bucket"
  table_name         = "tfcautious-kowari"
  source_column_name = "value"
}
