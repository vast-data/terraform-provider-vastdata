data "vastdata_s3_policy" "vastdb_s3_policy_by_id" {
  id = 1
}

data "vastdata_s3_policy" "vastdb_s3_policy_by_guid" {
  guid = "00000000-0000-0000-0000-000000000001"
}

