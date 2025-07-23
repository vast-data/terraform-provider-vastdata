# Copyright (c) HashiCorp, Inc.

resource "vastdata_qos_policy" "vastdb_qos_policy" {
  name        = "vastdb_qos_policy"
  policy_type = "USER"
  attached_users = [
    {
      name             = "runner"
      fqdn             = "runner.vastdb.local"
      identifier_type  = "username"
      identifier_value = "runner"
    }
  ]

  static_limits = {
    max_writes_bw_mbps = 1000
    max_reads_iops     = 2000
    max_writes_iops    = 3000
  }
}

