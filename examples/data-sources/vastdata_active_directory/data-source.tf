data "vastdata_active_directory" "vastdb_active_directory_by_id" {
  id = 1
}

data "vastdata_active_directory" "vastdb_active_directory_by_guid" {
  guid = "00000000-0000-0000-0000-000000000001"
}

data "vastdata_active_directory" "vastdb_active_directory1" {
  machine_account_name = "machine_acc"
}
