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
  name              = "vastdb_quota_group"
  gid               = 5001
  local_provider_id = 1
}

resource "vastdata_user" "vastdb_quota_user" {
  name              = "vastdb-quota-user"
  uid               = 5002
  local_provider_id = 1
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

  user_quotas = [{
    name            = vastdata_user.vastdb_quota_user.name
    identifier      = vastdata_user.vastdb_quota_user.name
    email           = "user1@example.com"
    identifier_type = "username"
    is_group        = false
    hard_limit      = 100000
    soft_limit      = 50000
    grace_period    = "90m"
  }]
  group_quotas = [{
    name            = vastdata_group.vastdb_quota_group.name
    identifier      = vastdata_group.vastdb_quota_group.name
    identifier_type = "group"
    is_group        = true
    hard_limit      = 100000
    soft_limit      = 50000
    grace_period    = "90m"
  }]
}

# --------------------

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

# --------------------

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

# --------------------


resource "vastdata_local_provider" "vastdb_provider" {
  name        = "vastdb-provider"
  description = "Test local provider"
  managed_by  = ["TENANT_ADMIN", "SUPER_ADMIN"]
}

resource "vastdata_tenant" "vastdb_tenant" {
  name              = "vastdb-tenant"
  client_ip_ranges  = [["192.168.100.50", "192.168.100.54"]]
  local_provider_id = vastdata_local_provider.vastdb_provider.id
  force_delete      = true
}

resource "vastdata_user" "vastdb_user" {
  name              = "vastdb_user"
  local_provider_id = vastdata_local_provider.vastdb_provider.id
}

resource "vastdata_group" "vastdb_group" {
  name              = "vastdb_group"
  gid               = 50001
  local_provider_id = vastdata_local_provider.vastdb_provider.id
}

resource "vastdata_view_policy" "vastdb_view_policy" {
  name          = "vastdb-view-policy"
  flavor        = "NFS"
  nfs_no_squash = ["10.0.0.1", "10.0.0.2"]
  tenant_id     = vastdata_tenant.vastdb_tenant.id
}

resource "vastdata_view" "vastdb_view" {
  path       = "/vastdb-path"
  policy_id  = vastdata_view_policy.vastdb_view_policy.id
  create_dir = true
  protocols  = ["NFS", "NFS4"]
  tenant_id  = vastdata_tenant.vastdb_tenant.id
}

resource "vastdata_quota" "quota1" {
  name              = "vastdb-quota"
  path              = vastdata_view.vastdb_view.path
  default_email     = "user@example.com"
  soft_limit        = 100000
  hard_limit        = 100000
  soft_limit_inodes = 1000000
  hard_limit_inodes = 1000000
  is_user_quota     = true
  enable_alarms     = true
  tenant_id         = vastdata_tenant.vastdb_tenant.id
  default_user_quota = {
    soft_limit        = 50000
    hard_limit        = 100000
    hard_limit_inodes = 10000
  }
  default_group_quota = {
    soft_limit        = 75000
    hard_limit        = 150000
    hard_limit_inodes = 15000
  }
  user_quotas = [{
    grace_period    = "02:00:00"
    hard_limit      = 15000
    soft_limit      = 15000
    identifier      = "vastdb-quota-user"
    identifier_type = "username"
    name            = vastdata_user.vastdb_user.name
    email           = "user1@example.com"
  }]
  group_quotas = [{
    grace_period    = "4 03:00:00"
    hard_limit      = 15000
    soft_limit      = 15000
    identifier      = "vastdb-quota-group"
    identifier_type = "group"
    name            = "vastdb-quota-group"
  }]
}

# --------------------

