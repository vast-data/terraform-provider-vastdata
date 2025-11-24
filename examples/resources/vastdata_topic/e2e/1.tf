data "vastdata_view_policy" "vastdb_view_policy_s3" {
  name = "s3_default_policy"
}

resource "vastdata_vip_pool" "vastdb_vippool" {
  name        = "vastdb_vippool"
  role        = "PROTOCOLS"
  subnet_cidr = "24"

  ip_ranges = [
    ["15.0.0.6", "15.0.0.10"],
    ["15.0.0.20", "15.0.0.40"]
  ]
}

resource "vastdata_user" "vastdb_user" {
  name              = "vastdb_user"
  local_provider_id = 1
}

resource "vastdata_view" "vastdb_view" {
  path            = "/vastdb_kafkastore"
  bucket          = "vastdb-kafkastore"
  bucket_owner    = vastdata_user.vastdb_user.name
  create_dir      = true
  policy_id       = data.vastdata_view_policy.vastdb_view_policy_s3.id
  protocols       = ["KAFKA", "S3", "DATABASE"]
  kafka_vip_pools = [vastdata_vip_pool.vastdb_vippool.id]
}


resource "vastdata_topic" "vastdb_topic" {
  database_name    = vastdata_view.vastdb_view.bucket
  name             = "vastdb_topic"
  topic_partitions = 3
  retention_ms     = 86400000 # 1 day retention
}
