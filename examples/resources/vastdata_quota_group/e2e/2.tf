# Quota group with a quota assigned via quotas_ids.

data "vastdata_view_policy" "vastdb_view_policy_default" {
  name = "default"
}

resource "vastdata_view" "vastdb_qgrp_view" {
  path       = "/vastdb_qgrp_view"
  create_dir = true
  policy_id  = data.vastdata_view_policy.vastdb_view_policy_default.id
  protocols  = ["NFS", "NFS4"]
}

resource "vastdata_quota" "vastdb_qgrp_quota" {
  name       = "vastdb_qgrp_quota"
  path       = vastdata_view.vastdb_qgrp_view.path
  soft_limit = 107374182400
  hard_limit = 214748364800
  create_dir = true
}

resource "vastdata_quota_group" "vastdb_qgrp_assigned" {
  name         = "vastdb_qgrp_assigned"
  grace_period = "90m"
  soft_limit   = 107374182400
  hard_limit   = 214748364800

  quotas_ids = [vastdata_quota.vastdb_qgrp_quota.id]
}
