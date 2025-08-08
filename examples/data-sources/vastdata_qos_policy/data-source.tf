data "vastdata_qos_policy" "vastdb_qos_policy_by_id" {
  id = 1
}

data "vastdata_qos_policy" "vastdb_qos_policy_by_guid" {
  guid = "00000000-0000-0000-0000-000000000001"
}

data "vastdata_qos_policy" "vastdb_qos_policy_by_name" {
  name = "vastdb_qos_policy"
}

