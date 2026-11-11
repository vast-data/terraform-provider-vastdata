data "vastdata_license" "by_id" {
  id = 1
}

data "vastdata_license" "by_name" {
  name = "example-license"
}

data "vastdata_license" "by_guid" {
  guid = "00000000-0000-0000-0000-000000000000"
}
