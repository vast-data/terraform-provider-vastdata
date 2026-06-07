# ignore:example
data "vastdata_view_policy" "vastdb_view_policy_s3_default_dup" {
  name = "s3_default_policy"
}

resource "vastdata_user" "vastdb_user_dup" {
  name              = "vastdb_user_dup"
  local_provider_id = 1
}

resource "vastdata_view" "vastdb_view_dup_1" {
  path         = "/vastdb_view_dup/s3/1"
  bucket       = "vastdb-s3-bucket-dup-1"
  create_dir   = true
  bucket_owner = vastdata_user.vastdb_user_dup.name
  policy_id    = data.vastdata_view_policy.vastdb_view_policy_s3_default_dup.id
  protocols    = ["S3"]
}

resource "vastdata_view" "vastdb_view_dup_2" {
  path         = "/vastdb_view_dup/s3/2"
  bucket       = "vastdb-s3-bucket-dup-2"
  create_dir   = true
  bucket_owner = vastdata_user.vastdb_user_dup.name
  policy_id    = data.vastdata_view_policy.vastdb_view_policy_s3_default_dup.id
  protocols    = ["S3"]
}

# Same name, different views — both must be created as separate API objects.
resource "vastdata_s3_life_cycle_rule" "vastdb_s3_lifecycle_rule_dup_1" {
  name            = "expiration-30-days"
  prefix          = ""
  view_id         = vastdata_view.vastdb_view_dup_1.id
  expiration_days = 30
  enabled         = true
}

resource "vastdata_s3_life_cycle_rule" "vastdb_s3_lifecycle_rule_dup_2" {
  name            = "expiration-30-days"
  prefix          = ""
  view_id         = vastdata_view.vastdb_view_dup_2.id
  expiration_days = 30
  enabled         = true
}
