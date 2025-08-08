resource "vastdata_quota" "vastdb_quota" {
  name            = "vastdb_quota_example"
  default_email   = "user@example.com"
  path            = "/vastdb_view/quota-example"
  soft_limit      = 100000
  hard_limit      = 200000
  create_dir_mode = 755
  create_dir      = true
  is_user_quota   = true
  grace_period    = "01:00:00"
  enable_alarms   = true
}