# Quota group with refresh_user_quotas = true.
# After every apply the provider calls PATCH /quotagroups/{id}/refresh_user_quotas/.

resource "vastdata_quota_group" "vastdb_qgrp_refresh" {
  name         = "vastdb_qgrp_refresh"
  grace_period = "1d"
  soft_limit   = 107374182400
  hard_limit   = 214748364800

  refresh_user_quotas = true
}
