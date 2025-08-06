
data "vastdata_user" "vastdb_user" {
  name = "runner"
}

data "vastdata_view_policy" "vastdb_view_policy_s3_default" {
  name = "s3_default_policy"
}

resource "vastdata_view" "vastdb_view1" {
  path         = "/vastdb_view/s3-1"
  bucket       = "vastdb-s3-bucket-1"
  create_dir   = true
  bucket_owner = data.vastdata_user.vastdb_user.name
  policy_id    = data.vastdata_view_policy.vastdb_view_policy_s3_default.id
  protocols    = ["S3"]
}

resource "vastdata_view" "vastdb_view2" {
  path         = "/vastdb_view/s3-2"
  bucket       = "vastdb-s3-bucket-2"
  create_dir   = true
  bucket_owner = data.vastdata_user.vastdb_user.name
  policy_id    = data.vastdata_view_policy.vastdb_view_policy_s3_default.id
  bucket_logging = {
    prefix         = "/logs"
    destination_id = vastdata_view.vastdb_view1.id
    key_format     = "PARTITIONED_PREFIX_DELIVERY_TIME"
  }
  protocols = ["S3"]
}