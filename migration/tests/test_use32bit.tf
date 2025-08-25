# Copyright (c) HashiCorp, Inc.

resource "vastdata_view_policy" "data-pcaps" {
    name                                 = "data-pcaps"
    use32bit_fileid                      = false
    nfs_read_only                        = ["user1", "user2"]
    smb_read_write                       = ["admin"]
    port_membership                      = "ALL"
}

resource "vastdata_view_policy" "other_policy" {
    name                                 = "other"
    use32bit_fileid                      = true
    permissions_list                     = ["read", "write"]
}
