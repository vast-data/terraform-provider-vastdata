#Create a view with S3 protocol + user , and attach an s3 lifecycle rul
resource "vastdata_user" "s3user" {
  name = "s3user"
  uid  = 2000

}

resource "vastdata_view" "s3-bucket-view" {
  policy_id    = data.vastdata_view_policy.s3_default_policy.id
  path         = "/s3view"
  bucket       = "s3view"
  protocols    = ["S3"]
  bucket_owner = vastdata_user.s3user.name
  create_dir   = true
}

resource "vastdata_s3_life_cycle_rule" "s3-bucket-view-lifecycle-rule" {
  name                      = "rule1"
  max_size                  = 10000000
  min_size                  = 100000
  newer_noncurrent_versions = 3
  prefix                    = "prefix"
  view_id                   = vastdata_view.s3-bucket-view.id
  expiration_days           = 30
  enabled                   = true
  noncurrent_days           = 20


}
