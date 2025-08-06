
data "vastdata_view_policy" "vastdb_view_policy_default" {
  name = "default"
}

resource "vastdata_view" "vastdb_view" {
  path       = "/vastdb_view/quota-example"
  create_dir = true
  policy_id  = data.vastdata_view_policy.vastdb_view_policy_default.id
  protocols  = ["NFS", "NFS4"]
}

resource "vastdata_quota" "vastdb_quota" {
  name            = "vastdb_quota_example"
  default_email   = "user@example.com"
  path            = vastdata_view.vastdb_view.path
  soft_limit      = 100000
  hard_limit      = 200000
  create_dir_mode = 755
  create_dir      = true
  is_user_quota   = true
  grace_period    = "01:00:00"
  enable_alarms   = true
}