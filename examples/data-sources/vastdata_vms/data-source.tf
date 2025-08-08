data "vastdata_vms" "vastdb_vms_by_id" {
  id = 1
}

data "vastdata_vms" "vastdb_vms_by_guid" {
  guid = "00000000-0000-0000-0000-000000000001"
}

data "vastdata_vms" "vastdb_vms_by_name" {
  name = "vastdb-vms"
}

