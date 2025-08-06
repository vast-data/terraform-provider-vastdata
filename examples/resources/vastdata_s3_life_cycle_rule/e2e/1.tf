
data "vastdata_view_policy" "vastdb_view_policy_s3_default" {
  name = "s3_default_policy"
}

resource "vastdata_view" "vastdb_view" {
  path         = "/vastdb_view/s3"
  bucket       = "vastdb-s3-bucket"
  create_dir   = true
  bucket_owner = "runner"
  policy_id    = data.vastdata_view_policy.vastdb_view_policy_s3_default.id
  protocols    = ["S3"]
}

resource "vastdata_s3_life_cycle_rule" "vastdb_s3_lifecycle_rule" {
  name                      = "vastdb_s3_lifecycle_rule"
  max_size                  = 10000000
  min_size                  = 100000
  newer_noncurrent_versions = 3
  prefix                    = "/s3/"
  view_id                   = vastdata_view.vastdb_view.id
  expiration_days           = 30
  enabled                   = true
  noncurrent_days           = 20
}
