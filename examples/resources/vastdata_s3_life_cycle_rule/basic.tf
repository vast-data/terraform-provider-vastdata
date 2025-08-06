
resource "vastdata_s3_life_cycle_rule" "vastdb_s3_lifecycle_rule" {
  name                      = "vastdb_s3_lifecycle_rule"
  max_size                  = 10000000
  min_size                  = 100000
  newer_noncurrent_versions = 3
  prefix                    = "/s3/"
  view_id                   = 2
  expiration_days           = 30
  enabled                   = true
  noncurrent_days           = 20
}
