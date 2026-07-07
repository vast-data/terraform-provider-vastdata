# Basic quota group — limits only, no action triggers.

resource "vastdata_quota_group" "vastdb_qgroup" {
  name         = "vastdb_qgroup"
  grace_period = "90m"
  soft_limit   = 107374182400
  hard_limit   = 214748364800
}
