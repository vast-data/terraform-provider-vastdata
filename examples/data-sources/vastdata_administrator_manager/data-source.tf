data "vastdata_administrator_manager" "vastdb_administrator_manager_by_id" {
  id = 1
}

data "vastdata_administrator_manager" "vastdb_administrator_manager_by_guid" {
  guid = "00000000-0000-0000-0000-000000000001"
}

data "vastdata_administrator_manager" "vastdb_administrator_manager_by_username" {
  username = "vastdb_manager"
}
