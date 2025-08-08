data "vastdata_volume" "vastdb_volume_by_id" {
  id = 1
}

data "vastdata_volume" "vastdb_volume_by_name" {
  name = "vastdb-volume"
}
