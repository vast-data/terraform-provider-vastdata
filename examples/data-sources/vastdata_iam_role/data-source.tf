data "vastdata_iam_role" "vastdb_iam_role_by_id" {
  id = 1
}

data "vastdata_iam_role" "vastdb_iam_role_by_guid" {
  guid = "00000000-0000-0000-0000-000000000001"
}

data "vastdata_iam_role" "role_by_name" {
  name = "vastdb_iam_role"
}
