
resource "vastdata_user" "vastdb_user" {
  name              = "vastdb_user"
  local_provider_id = 1
}

data "vastdata_view_policy" "vastdb_view_policy_s3_default" {
  name = "s3_default_policy"
}

resource "vastdata_view" "vastdb_view" {
  path         = "/vastdb_view/s3-cors-full"
  bucket       = "vastdb-s3-cors-full-bucket"
  create_dir   = true
  bucket_owner = vastdata_user.vastdb_user.name
  policy_id    = data.vastdata_view_policy.vastdb_view_policy_s3_default.id
  protocols    = ["S3"]

  s3cors_configuration = {
    cors_rules = [
      {
        allowed_methods = ["GET", "PUT", "POST", "DELETE"]
        allowed_origins = ["https://app.example.com", "https://admin.example.com"]
        allowed_headers = ["Authorization", "Content-Type", "X-Requested-With"]
        expose_headers  = ["ETag", "x-amz-request-id", "x-amz-id-2"]
        max_age_seconds = 3600
      },
      {
        allowed_methods = ["GET"]
        allowed_origins = ["*"]
        max_age_seconds = 86400
      },
    ]
  }
}
