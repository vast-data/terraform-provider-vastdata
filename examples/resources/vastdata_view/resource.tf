
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

# --------------------


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

# --------------------

