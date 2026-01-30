
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
