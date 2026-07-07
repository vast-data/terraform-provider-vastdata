# ignore:e2e

resource "vastdata_quota_group" "vastdb_qgrp_grace" {
  name         = "vastdb_qgrp_grace"
  grace_period = "7d"
  soft_limit   = 107374182400
  hard_limit   = 214748364800

  reset_grace_period = true
}
