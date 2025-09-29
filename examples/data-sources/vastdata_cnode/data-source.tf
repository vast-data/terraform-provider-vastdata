data "vastdata_cnode" "vastdb_cnode_by_id" {
  id = 1
}

data "vastdata_cnode" "vastdb_cnode_by_guid" {
  guid = "00000000-0000-0000-0000-000000000001"
}

data "vastdata_cnode" "vastdb_cnode_by_name" {
  name = "cnode-1"
}

data "vastdata_cnode" "vastdb_cnode_by_hostname" {
  hostname = "cnode01.example.com"
}

data "vastdata_cnode" "vastdb_cnode_by_ip" {
  ip = "10.0.1.100"
}
