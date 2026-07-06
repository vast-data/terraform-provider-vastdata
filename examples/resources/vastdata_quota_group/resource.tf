resource "vastdata_quota_group" "example" {
  name         = "my-quota-group"
  tenant_id    = vastdata_tenant.example.id
  grace_period = "7d"
  soft_limit   = 1099511627776 # 1 TiB
  hard_limit   = 2199023255552 # 2 TiB

  enable_alarms = true
  default_email = "storage-alerts@example.com"

  # Assign existing quotas to this group.
  quotas_ids = [
    vastdata_quota.quota1.id,
    vastdata_quota.quota2.id,
  ]

  # Uncomment to refresh user quota counters after every apply.
  # refresh_user_quotas = true

  # Uncomment to reset the grace period countdown after every apply.
  # reset_grace_period = true
}
