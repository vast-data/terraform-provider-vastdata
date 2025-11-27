resource "vastdata_user" "vastdb_user" {
  name              = "vastdbuser"
  local_provider_id = 1
}

resource "vastdata_qos_policy" "vastdb_qos_policy" {
  name        = "vastdb_qos_policy"
  policy_type = "USER"
  attached_users = [
    {
      name             = vastdata_user.vastdb_user.name
      fqdn             = "${vastdata_user.vastdb_user.name}.vastdb.local"
      identifier_type  = "username"
      identifier_value = vastdata_user.vastdb_user.name
    }
  ]

  static_limits = {
    max_writes_bw_mbps = 1000
    max_reads_iops     = 2000
    max_writes_iops    = 3000
  }
}

