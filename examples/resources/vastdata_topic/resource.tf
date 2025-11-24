resource "vastdata_topic" "vastdb_topic" {
  database_name    = "vastdb_bucket"
  name             = "vastdb_topic"
  topic_partitions = 3
  retention_ms     = 86400000 # 1 day retention
}
