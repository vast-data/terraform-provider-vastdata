data "vastdata_view_policy" "default_policy_groups" {
  name = "default"
}

resource "vastdata_view" "group_defaults_only" {
  path       = "/vastdb-group-defaults"
  policy_id  = data.vastdata_view_policy.default_policy_groups.id
  create_dir = true
  protocols  = ["NFS", "NFS4"]
}

# Quota that applies default limits to ALL groups without exceptions
resource "vastdata_quota" "group_defaults_only" {
  name          = "group-vastdb-defaults-only"
  path          = vastdata_view.group_defaults_only.path
  default_email = "group-admin@example.com"

  # Enable user/group quotas
  is_user_quota = true
  enable_alarms = true

  # Create directory with specific permissions
  create_dir      = true
  create_dir_mode = 775
  inherit_acl     = false

  # Overall directory limits (capacity and inodes)
  soft_limit        = 5000000000  # 5GB total
  hard_limit        = 10000000000 # 10GB total
  soft_limit_inodes = 50000
  hard_limit_inodes = 100000
  grace_period      = "12:00:00" # 12 hours

  # Every group gets the same limits
  # Useful for enforcing uniform quota policies across all groups
  default_group_quota = {
    soft_limit        = 500000000  # 500MB per group
    hard_limit        = 1000000000 # 1GB per group
    soft_limit_inodes = 5000       # 5000 files/dirs per group
    hard_limit_inodes = 10000      # 10000 files/dirs per group
    grace_period      = "6:00:00"  # 6 hours
  }
}
