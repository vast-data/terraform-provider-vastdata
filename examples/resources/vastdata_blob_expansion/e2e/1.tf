# ignore:example

data "vastdata_view_policy" "vastdb_view_policy_s3" {
  name = "s3_default_policy"
}

resource "vastdata_vip_pool" "vastdb_blob_exp_vippool" {
  name        = "vastdb-blob-exp-vippool"
  role        = "PROTOCOLS"
  subnet_cidr = "24"

  ip_ranges = [
    ["15.0.3.6", "15.0.3.10"],
    ["15.0.3.20", "15.0.3.40"],
  ]
}

resource "vastdata_user" "vastdb_blob_exp_user" {
  name              = "vastdb-blob-exp-user"
  local_provider_id = 1
}

resource "vastdata_view" "vastdb_blob_exp_view" {
  path            = "/vastdb-blob-exp-store"
  bucket          = "vastdb-blob-exp-store"
  bucket_owner    = vastdata_user.vastdb_blob_exp_user.name
  create_dir      = true
  policy_id       = data.vastdata_view_policy.vastdb_view_policy_s3.id
  protocols       = ["KAFKA", "S3", "DATABASE"]
  kafka_vip_pools = [vastdata_vip_pool.vastdb_blob_exp_vippool.id]
}

resource "vastdata_topic" "vastdb_blob_exp_topic" {
  database_name    = vastdata_view.vastdb_blob_exp_view.bucket
  name             = "vastdb-blob-exp-topic"
  topic_partitions = 3
  retention_ms     = 86400000 # 1 day
}

resource "vastdata_blob_expansion" "vastdb_blob_exp" {
  database_name       = vastdata_view.vastdb_blob_exp_view.bucket
  table_name          = vastdata_topic.vastdb_blob_exp_topic.name
  expansion_format    = "json"
  target_table_name   = "vastdb-blob-exp-target"
  target_table_schema = "kafka_topics"

  arrow_schema = [
    {
      name = "message_id"
      field = {
        column_type = "string"
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
