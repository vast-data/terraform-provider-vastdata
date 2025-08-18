# Copyright (c) HashiCorp, Inc.


## holds all resources specific to the requirements of the pcap view

resource "vastdata_view" "data-pcaps" {
  path       = "/data/pcaps"
  policy_id  = vastdata_view_policy.data-pcaps.id
  create_dir = "true"
  protocols  = ["NFS"]
}

resource "vastdata_view" "pcaps-raw" {
    path       = "/data/pcaps/raw"
    bucket     = "raw"
    create_dir = "true"
    policy_id  = vastdata_view_policy.data-pcaps.id
    protocols  = ["NFS","S3"]
    bucket_owner = "svcpcaps@mavensecurities.com"
}

resource "vastdata_view" "pcaps-filtered" {
    path       = "/data/pcaps/filtered" ## group_owner:sgdatapcapsfiltered
    create_dir = "true"
    policy_id  = vastdata_view_policy.data-pcaps.id
    protocols  = ["NFS","S3"]
    bucket     = "filtered"
    bucket_owner = "svcpcapsfiltered@mavensecurities.com"
}


resource "vastdata_view_policy" "data-pcaps" {
    access_flavor                            = "ALL"
    allowed_characters                       = "LCD"
    apple_sid                                = true
    atime_frequency                          = null
    auth_source                              = "RPC_AND_PROVIDERS"
    cluster                                  = "VAST-MAVEN-2"
    cluster_id                               = 1
    count_views                              = 0
    data_create_delete                       = false
    data_modify                              = false
    data_read                                = false
    enable_access_to_snapshot_dir_in_subdirs = true
    enable_listing_of_snapshot_dir           = false
    enable_snapshot_lookup                   = true
    enable_visibility_of_snapshot_dir        = false
    expose_id_in_fsid                        = false
    flavor                                   = "NFS"
    gid_inheritance                          = "LINUX"
    log_deleted                              = false
    log_full_path                            = false
    log_hostname                             = false
    log_username                             = false
    name                                     = "data-pcaps"
    nfs_all_squash                           = []
    nfs_case_insensitive                     = false
    nfs_enforce_tls                          = false
    nfs_minimal_protection_level             = "SYSTEM"
    nfs_no_squash                            = []
    nfs_posix_acl                            = true
    nfs_read_only                            = []
    nfs_read_write                           = [
        "*",
    ]
    nfs_return_open_permissions              = false
    nfs_root_squash                          = [
        "*",
    ]
    path_length                              = "LCD"
    protocols                                = []
    read_only                                = []
    read_write                               = [
        "*",
    ]
    s3_bucket_full_control                   = null
    s3_bucket_listing                        = null
    s3_bucket_read                           = null
    s3_bucket_read_acp                       = null
    s3_bucket_write                          = null
    s3_bucket_write_acp                      = null
    s3_object_full_control                   = null
    s3_object_read                           = null
    s3_object_read_acp                       = null
    s3_object_write                          = null
    s3_object_write_acp                      = null
    s3_read_only                             = []
    s3_read_write                            = [
        "*",
    ]
    s3_special_chars_support                 = false
    s3_visibility                            = []
    s3_visibility_groups                     = []
    smb_directory_mode                       = 775
    smb_directory_mode_padded                = "775"
    smb_file_mode                            = 664
    smb_file_mode_padded                     = "664"
    smb_is_ca                                = false
    smb_read_only                            = []
    smb_read_write                           = [
        "*",
    ]
    tenant_id                                = 1
    tenant_name                              = "default"
    trash_access                             = []
    use32bit_fileid                          = false
    use_auth_provider                        = true
    vip_pools     = [
            vastdata_vip_pool.prod.id,
            vastdata_vip_pool.dev.id
    ]


    protocols_audit {
        create_delete_files_dirs_objects = false
        log_deleted_files_dirs           = false
        log_full_path                    = false
        log_hostname                     = false
        log_username                     = false
        modify_data                      = false
        modify_data_md                   = false
        read_data                        = false
        read_data_md                     = false
    }

//    vippool_permissions {
//        vippool_id          = vastdata_vip_pool.prod.id
//        vippool_permissions = "RW"
//    }

  lifecycle {
    ignore_changes = [
      count_views,
    ]
  }

}

resource "vastdata_quota" "pcaps" {
  name          = "pcaps"
  default_email = "michael.moyles@mavensecurities.com"
  path          = vastdata_view.data-pcaps.path
    // limits are in bytes == 300TB
  soft_limit    = 290000000000000
  hard_limit    = 300000000000000
  is_user_quota = false
}


