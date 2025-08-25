# Copyright (c) HashiCorp, Inc.

resource "vastdata_view_policy" "data-pcaps" {
    name                                 = "data-pcaps"
    use32bit_fileid                      = false
    port_membership                      = "ALL"  # This should remain uppercase
    
    # Attributes that should be removed (unsupported)
    s3_bucket_full_control               = null
    
    # Attributes that should be removed (read-only)
    tenant_name                          = "default"
    smb_directory_mode_padded            = "775"
    smb_file_mode_padded                 = "664"
    log_username                         = false
    log_hostname                         = false
    log_full_path                        = false
    log_deleted                          = false
    enable_snapshot_lookup               = true
    enable_listing_of_snapshot_dir       = false
    data_modify                          = false
    data_create_delete                   = false
    data_read                            = false
    cluster                              = "VAST-MAVEN-2"
    count_views                          = 0
    
    # Attributes that should remain
    nfs_read_only                        = ["user1", "user2"]
}

resource "vastdata_administrator_manager" "prometheus_reader" {
    username                            = "prometheus_reader"
    # This should be renamed to permissions_list for administrator_manager
    permissions                         = ["create_monitoring"]
}

resource "vastdata_administators_managers" "old_name" {
    username                            = "old_admin"
    # This starts as permissions_list, should become permissions_list in new resource type
    permissions_list                    = ["read", "write"]
}
