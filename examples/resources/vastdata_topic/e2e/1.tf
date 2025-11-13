data "vastdata_view_policy" "vastdb_view_policy_s3" {
  name = "s3_default_policy"
}

data "vastdata_vip_pool" "vastdb_vip_pool" {
  name = "vippool-comet"
}

resource "vastdata_view" "vastdb_view" {
  path            = "/kafkastore"
  bucket          = "kafkastore"
  bucket_owner    = "runner"
  create_dir      = true
  policy_id       = data.vastdata_view_policy.vastdb_view_policy_s3.id
  protocols       = ["KAFKA", "S3", "DATABASE"]
  kafka_vip_pools = [data.vastdata_vip_pool.vastdb_vip_pool.id]
}


resource "vastdata_topic" "vastdb_topic" {
  database_name    = vastdata_view.vastdb_view.bucket
  name             = "vastdb_topic"
  topic_partitions = 3
  retention_ms     = 86400000 # 1 day retention
}
