resource "vastdata_block_host" "vastdb_block_host" {
  name      = "vastdb-block-host"
  tenant_id = 1
  nqn       = "nqn.2014-08.org.nvmexpress:uuid:12345678-1234-1234-1234-123456789012"
}