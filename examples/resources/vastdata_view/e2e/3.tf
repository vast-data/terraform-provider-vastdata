
data "vastdata_user" "vastdb_user" {
  name = "runner"
}

data "vastdata_view_policy" "vastdb_view_policy_s3_default" {
  name = "s3_default_policy"
}


resource "vastdata_view" "vastdb_view" {
  path                 = "/vastdb_view/s3"
  bucket               = "vastdb-s3-bucket"
  create_dir           = true
  bucket_owner         = data.vastdata_user.vastdb_user.name
  policy_id            = data.vastdata_view_policy.vastdb_view_policy_s3_default.id
  s3_unverified_lookup = true
  protocols            = ["S3"]
}