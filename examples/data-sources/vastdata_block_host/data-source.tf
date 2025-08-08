data "vastdata_block_host" "vastdb_block_host_by_id" {
  id = 1
}

data "vastdata_block_host" "vastdb_block_host_by_name" {
  name = "vastdb-block-host"
}

data "vastdata_block_host" "vastdb_block_host_by_nqn" {
  nqn = "nqn.2014-08.org.nvmexpress:uuid:12345678-1234-1234-1234-123456789012"
}
