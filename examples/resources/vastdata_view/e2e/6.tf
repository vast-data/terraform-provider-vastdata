
data "vastdata_user" "vastdb_user" {
  name = "runner"
}

data "vastdata_view_policy" "vastdb_view_policy_s3_default" {
  name = "s3_default_policy"
}

resource "vastdata_view" "vastdb_view" {
  path                      = "/vastdb_view-bucket"
  bucket                    = "vastdb-bucket"
  create_dir                = true
  bucket_owner              = data.vastdata_user.vastdb_user.name
  policy_id                 = data.vastdata_view_policy.vastdb_view_policy_s3_default.id
  allow_s3_anonymous_access = true
  s3_versioning             = true
  create_dir_mode           = 777
  protocols                 = ["S3"]
}
