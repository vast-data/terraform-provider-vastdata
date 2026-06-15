
resource "vastdata_view" "vastdb_view" {
  path       = "/vastdb_view/example"
  policy_id  = 2
  create_dir = true
  protocols  = ["NFS", "NFS4"]
}

# ---------------------
# Complete examples
# ---------------------


data "vastdata_view_policy" "vastdb_view_policy_default" {
  name = "default"
}

resource "vastdata_view" "vastdb_view" {
  path       = "/vastdb_view/example"
  policy_id  = data.vastdata_view_policy.vastdb_view_policy_default.id
  create_dir = true
  protocols  = ["NFS", "NFS4"]
}

# --------------------


data "vastdata_tenant" "vastdb_default_tenant" {
  name = "default"
}

data "vastdata_view_policy" "vastdb_view_policy_default" {
  name = "default"
}

resource "vastdata_view" "vastdb_view" {
  path                       = "/vastdb_view/example"
  alias                      = "/vastdb_view-aliased"
  tenant_id                  = data.vastdata_tenant.vastdb_default_tenant.id
  policy_id                  = data.vastdata_view_policy.vastdb_view_policy_default.id
  create_dir                 = true
  select_for_live_monitoring = true
  protocols                  = ["NFS"]
}

# --------------------

resource "vastdata_user" "vastdb_user" {
  name              = "vastdb_user"
  local_provider_id = 1
}

data "vastdata_view_policy" "vastdb_view_policy_s3_default" {
  name = "s3_default_policy"
}

resource "vastdata_view" "vastdb_view" {
  path                 = "/vastdb_view/s3"
  bucket               = "vastdb-s3-bucket"
  create_dir           = true
  bucket_owner         = vastdata_user.vastdb_user.name
  policy_id            = data.vastdata_view_policy.vastdb_view_policy_s3_default.id
  s3_unverified_lookup = true
  protocols            = ["S3"]
}

# --------------------

resource "vastdata_user" "vastdb_user" {
  name              = "vastdb_user"
  local_provider_id = 1
}

data "vastdata_view_policy" "vastdb_view_policy_s3_default" {
  name = "s3_default_policy"
}

resource "vastdata_view" "vastdb_view1" {
  path         = "/vastdb_view/s3-1"
  bucket       = "vastdb-s3-bucket-1"
  create_dir   = true
  bucket_owner = vastdata_user.vastdb_user.name
  policy_id    = data.vastdata_view_policy.vastdb_view_policy_s3_default.id
  protocols    = ["S3"]
}

resource "vastdata_view" "vastdb_view2" {
  path         = "/vastdb_view/s3-2"
  bucket       = "vastdb-s3-bucket-2"
  create_dir   = true
  bucket_owner = vastdata_user.vastdb_user.name
  policy_id    = data.vastdata_view_policy.vastdb_view_policy_s3_default.id
  bucket_logging = {
    prefix         = "/logs"
    destination_id = vastdata_view.vastdb_view1.id
    key_format     = "PARTITIONED_PREFIX_DELIVERY_TIME"
  }
  protocols = ["S3"]
}

# --------------------


data "vastdata_view_policy" "vastdb_view_policy_default" {
  name = "default"
}

resource "vastdata_view" "vastdb_view" {
  path                 = "/vastdb_view/subsystem"
  name                 = "vastdb-subsystem"
  create_dir           = true
  is_default_subsystem = true
  policy_id            = data.vastdata_view_policy.vastdb_view_policy_default.id
  protocols            = ["BLOCK"]
}

# --------------------

resource "vastdata_user" "vastdb_user" {
  name              = "vastdb_user"
  local_provider_id = 1
}

data "vastdata_view_policy" "vastdb_view_policy_s3_default" {
  name = "s3_default_policy"
}

resource "vastdata_view" "vastdb_view" {
  path                      = "/vastdb_view-bucket"
  bucket                    = "vastdb-bucket"
  create_dir                = true
  bucket_owner              = vastdata_user.vastdb_user.name
  policy_id                 = data.vastdata_view_policy.vastdb_view_policy_s3_default.id
  allow_s3_anonymous_access = true
  s3_versioning             = true
  create_dir_mode           = 777
  protocols                 = ["S3"]
}

# --------------------

# View with NFSv4 triggers enabled
data "vastdata_view_policy" "nfs4_policy" {
  name = "default"
}

resource "vastdata_view" "nfs4_triggers_view" {
  path                 = "/vastdb_view/nfs4"
  create_dir           = true
  policy_id            = data.vastdata_view_policy.nfs4_policy.id
  protocols            = ["NFS4"]
  enable_nfs4_triggers = true
}

# --------------------

# S3 view with CORS configuration
data "vastdata_view_policy" "s3_cors_policy" {
  name = "s3_default_policy"
}

resource "vastdata_user" "cors_bucket_owner" {
  name              = "cors_owner"
  local_provider_id = 1
}

resource "vastdata_view" "s3_cors_view" {
  path         = "/vastdb_view/s3-cors"
  bucket       = "vastdb-s3-cors-bucket"
  create_dir   = true
  bucket_owner = vastdata_user.cors_bucket_owner.name
  policy_id    = data.vastdata_view_policy.s3_cors_policy.id
  protocols    = ["S3"]

  s3cors_configuration = {
    cors_rules = [
      {
        allowed_methods = ["GET", "PUT", "POST"]
        allowed_origins = ["https://example.com", "https://app.example.com"]
        allowed_headers = ["Authorization", "Content-Type"]
        expose_headers  = ["ETag", "x-amz-request-id"]
        max_age_seconds = 3600
      },
      {
        allowed_methods = ["GET"]
        allowed_origins = ["*"]
      },
    ]
  }
}

# --------------------

