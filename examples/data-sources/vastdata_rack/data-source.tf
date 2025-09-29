data "vastdata_rack" "vastdb_rack_by_id" {
  id = 1
}

data "vastdata_rack" "vastdb_rack_by_guid" {
  guid = "00000000-0000-0000-0000-000000000001"
}

data "vastdata_rack" "vastdb_rack_by_name" {
  name = "rack-1"
}
