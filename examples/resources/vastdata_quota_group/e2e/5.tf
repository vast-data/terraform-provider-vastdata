# Quota group with quota assignment and user quota refresh.

data "vastdata_view_policy" "vastdb_view_policy_default" {
  name = "default"
}

resource "vastdata_view" "vastdb_qgrp_full_view" {
  path       = "/vastdb_qgrp_full_view"
  create_dir = true
  policy_id  = data.vastdata_view_policy.vastdb_view_policy_default.id
  protocols  = ["NFS", "NFS4"]
}

resource "vastdata_quota" "vastdb_qgrp_full_quota" {
  name          = "vastdb_qgrp_full_quota"
  path          = vastdata_view.vastdb_qgrp_full_view.path
  soft_limit    = 107374182400
  hard_limit    = 214748364800
  is_user_quota = true
  create_dir    = true
}

resource "vastdata_quota_group" "vastdb_qgrp_full" {
  name          = "vastdb_qgrp_full"
  grace_period  = "7d"
  soft_limit    = 107374182400
  hard_limit    = 214748364800
  enable_alarms = true

  quotas_ids          = [vastdata_quota.vastdb_qgrp_full_quota.id]
  refresh_user_quotas = true
}
