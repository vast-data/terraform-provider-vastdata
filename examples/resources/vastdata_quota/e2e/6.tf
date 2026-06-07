# ignore:example
resource "vastdata_tenant" "vastdb_quota_zero_limits" {
  name             = "vastdb-zero-limits-tenant"
  client_ip_ranges = [["192.168.5.1", "192.168.5.2"]]
  force_delete     = true
}

resource "vastdata_view_policy" "vastdb_quota_zero_limits" {
  name      = "vastdb-zero-limits-policy"
  flavor    = "NFS"
  tenant_id = vastdata_tenant.vastdb_quota_zero_limits.id
}

resource "vastdata_view" "vastdb_quota_zero_limits" {
  path       = "/vastdb-zero-limits-view"
  policy_id  = vastdata_view_policy.vastdb_quota_zero_limits.id
  create_dir = true
  protocols  = ["NFS"]
  tenant_id  = vastdata_tenant.vastdb_quota_zero_limits.id
}

resource "vastdata_quota" "vastdb_quota_zero_limits" {
  name              = "vastdb-zero-limits-quota"
  path              = vastdata_view.vastdb_quota_zero_limits.path
  tenant_id         = vastdata_tenant.vastdb_quota_zero_limits.id
  is_user_quota     = false
  hard_limit        = 0
  hard_limit_inodes = 0
  soft_limit        = 0
  soft_limit_inodes = 0
}
