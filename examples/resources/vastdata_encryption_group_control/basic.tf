data "vastdata_encryption_group" "vastdb_encryption_group" {
  id = 1
}

resource "vastdata_encryption_group_control" "vastdb_encryption_group_control" {
  id     = data.vastdata_encryption_group.vastdb_encryption_group.id
  action = "revoke"
}