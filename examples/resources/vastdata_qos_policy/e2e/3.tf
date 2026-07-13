locals {
  ips = ["172.18.1.123", "172.18.1.130", "172.18.1.134", "172.18.1.135", "172.18.1.137", "172.18.1.148", "172.18.1.232", "172.18.1.237", "172.18.1.238", "172.18.1.239", "172.18.31.105", "172.18.31.106", "172.18.31.107", "172.18.31.108",
  ]
}

resource "vastdata_view_policy" "vastdb_policy" {
  access_flavor                            = "ALL"
  allowed_characters                       = "NPL"
  apple_sid                                = true
  auth_source                              = "RPC"
  disable_handle_lease                     = false
  disable_read_lease                       = false
  disable_write_lease                      = false
  enable_access_to_snapshot_dir_in_subdirs = true
  enable_visibility_of_snapshot_dir        = false
  expose_id_in_fsid                        = false
  flavor                                   = "NFS"
  gid_inheritance                          = "LINUX"
  inherit_parent_mode_bits                 = false
  is_s3_default_policy                     = false
  name                                     = "vastdb_view_policy"
  nfs_all_squash                           = []
  nfs_case_insensitive                     = false
  nfs_enforce_tls                          = false
  nfs_enforce_tls_relaxed                  = false
  nfs_minimal_protection_level             = "NONE"
  nfs_no_squash                            = local.ips
  nfs_posix_acl                            = false
  nfs_read_only                            = []
  nfs_read_write                           = local.ips
  nfs_return_open_permissions              = false
  nfs_root_squash                          = []
  path_length                              = "NPL"
  permission_per_vip_pool = {
    "${vastdata_vip_pool.vastdb_vippool.id}" = "RW"
  }
  protocols            = []
  s3_read_only         = []
  s3_visibility        = []
  s3_visibility_groups = []
  smb_directory_mode   = 0
  smb_file_mode        = 0
  smb_is_ca            = false
  smb_read_only        = []
  smb_read_write = [
    "172.18.1.11",
    "172.18.1.12",
    "172.18.1.123",
    "172.18.1.148",
    "172.18.1.18",
    "172.18.1.25",
    "172.18.1.75",
    "172.18.1.76",
    "172.18.32.105",
  ]
  trash_access      = []
  use_32bit_fileid  = false
  use_auth_provider = false
}

resource "vastdata_vip_pool" "vastdb_vippool" {
  cnode_ids                 = []
  enable_l3                 = false
  enable_weighted_balancing = false
  enabled                   = true
  gw_ip                     = "172.18.31.1"
  gw_ipv6                   = ""
  ip_ranges = [
    [
      "172.18.31.11",
      "172.18.31.42",
    ],
  ]
  name            = "vastdb_vippool"
  port_membership = "ALL"
  role            = "PROTOCOLS"
  subnet_cidr     = 24
  tenant_id       = 1
}