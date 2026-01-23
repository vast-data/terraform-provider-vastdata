data "vastdata_view_policy" "default_policy" {
  name = "default"
}

resource "vastdata_view" "user_defaults_only" {
  path       = "/vastdb-user-defaults"
  policy_id  = data.vastdata_view_policy.default_policy.id
  create_dir = true
  protocols  = ["NFS"]
}

# Quota that applies default limits to ALL users without exceptions
resource "vastdata_quota" "user_defaults_only" {
  name          = "user-vastdb-defaults-only"
  path          = vastdata_view.user_defaults_only.path
  default_email = "quota-admin@example.com"

  # Enable user quotas
  is_user_quota          = true
  enable_alarms          = true
  enable_email_providers = true

  # Overall directory limits
  soft_limit   = 10000000000 # 10GB total
  hard_limit   = 20000000000 # 20GB total
  grace_period = "1:00:00"   # 1 hour

  # Every user gets the same limits
  default_user_quota = {
    soft_limit   = 100000000 # 100MB per user
    hard_limit   = 200000000 # 200MB per user
    grace_period = "30m"     # 30 minutes
  }
}
