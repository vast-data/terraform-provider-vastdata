resource "vastdata_folder_read_only" "vastdb_folder_readonly" {
  path      = "/vastdb/folder_readonly"
  tenant_id = 1
}
