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

# ---------------------
# Complete examples
# ---------------------


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

# --------------------

resource "vastdata_group" "vastdb_quota_group" {
  name = "vastdb_quota_group"
  gid  = 57593
}

resource "vastdata_user" "vastdb_quota_user" {
  name = "vastdb-quota-user"
  uid  = 776107
}

resource "vastdata_view_policy" "vastdb_view_policy" {
  name          = "vastdb_view_policy_example"
  flavor        = "NFS"
  nfs_no_squash = ["10.0.0.1", "10.0.0.2"]
}

resource "vastdata_view" "vastdb_view" {
  path       = "/vastdb_view/quota-example"
  policy_id  = vastdata_view_policy.vastdb_view_policy.id
  create_dir = true
  protocols  = ["NFS", "NFS4"]
}

resource "vastdata_quota" "vastdb_quota" {
  name          = "vastdb_quota_example"
  default_email = "user@example.com"
  path          = vastdata_view.vastdb_view.path
  soft_limit    = 100000
  hard_limit    = 100000
  is_user_quota = true
  default_user_quota = {
    grace_period      = "09 01:00:00"
    hard_limit        = 2000
    soft_limit        = 1000
    hard_limit_inodes = 20000000
  }
  user_quotas = [{
    name            = vastdata_user.vastdb_quota_user.name
    identifier      = vastdata_user.vastdb_quota_user.name
    email           = "user1@example.com"
    identifier_type = "username"
    is_group        = false
  }]
  group_quotas = [{
    name            = vastdata_group.vastdb_quota_group.name
    identifier      = vastdata_group.vastdb_quota_group.name
    identifier_type = "group"
    is_group        = true
  }]
}

# --------------------

