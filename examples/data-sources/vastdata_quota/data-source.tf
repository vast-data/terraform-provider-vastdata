data "vastdata_quota" "vastdb_quota_by_id" {
  id = 1
}

data "vastdata_quota" "vastdb_quota_by_guid" {
  guid = "00000000-0000-0000-0000-000000000001"
}

data "vastdata_quota" "vastdb_quota_by_name" {
  name = "vastdb_quota"
}
