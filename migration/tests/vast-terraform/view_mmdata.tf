# Copyright (c) HashiCorp, Inc.


resource "vastdata_view_policy" "mmdata-cd" {
  name          = "mmdata-cd"

  // only allow access via productrion networks
  vip_pools     = [vastdata_vip_pool.prod.id]
  flavor        = "NFS"
  nfs_root_squash       = ["*"]

  // read only from k8s CD dev
  nfs_read_only    = ["*"]
  // read write from k8s CD dev
  nfs_read_write   = ["*"]

  lifecycle {
    ignore_changes = [
      count_views,
    ]
  }

}

resource "vastdata_view_policy" "mmbus-cd" {
  name          = "mmbus-cd"

  ## group_owner:sgdatamarketmaking
  vip_pools     = [vastdata_vip_pool.prod.id,vastdata_vip_pool.dev.id]
  flavor        = "NFS"
  nfs_root_squash       = ["*"]
  ## PE-4200 read/write only allowed from cdmmdata01-03,cdmmdatastg01-02,cdinf03,cdplatform01
  nfs_read_write  = [
    "10.64.116.11",
    "10.66.19.7",
    "10.66.19.8",
    "10.66.19.9",
    "10.66.19.10",
    "10.66.19.11",
    "10.66.18.14",
    "10.66.4.67",
  ]
  nfs_read_only = ["*"]
  // s3 read write from everywhere
  s3_read_write = ["*"]

  // we need to provide write access to two users, svcmmdata (owner) and svccore
  nfs_posix_acl = true

  lifecycle {
    ignore_changes = [
      count_views,
    ]
  }

}


resource "vastdata_view" "data-mmdata" {
  path       = "/data/mmdata" ## group_owner:sgdatamarketmaking
  policy_id  = vastdata_view_policy.mmdata-cd.id
  create_dir = "true"
  protocols  = ["NFS"]
}

// add a bucket
resource "vastdata_view" "test-bucket" {
    path       = "/data/mmdata/test-bucket"
    bucket     = "test-bucket"
    create_dir = "true"
    policy_id  = vastdata_view_policy.s3_default_policy.id
    protocols  = ["NFS","S3"]
    bucket_owner = "svcmmdata@mavensecurities.com"
}

resource "vastdata_view" "mmbus-dev" {
    path       = "/data/mmdata/mmbus-dev"
    bucket     = "mmbus-dev"
    create_dir = "true"
    policy_id  = vastdata_view_policy.mmbus-cd.id
    protocols  = ["NFS","S3"]
    bucket_owner = "svcmmdata@mavensecurities.com" // view policy needs to support ACLs to allow svccore write access
}

resource "vastdata_view" "mmbus-prod" {
    path       = "/data/mmdata/mmbus-prod"
    bucket     = "mmbus-prod"
    create_dir = "true"
    policy_id  = vastdata_view_policy.mmbus-cd.id
    protocols  = ["NFS","S3"]
    bucket_owner = "svcmmdata@mavensecurities.com"
}

resource "vastdata_view" "mmbus-stage" {
    path       = "/data/mmdata/mmbus-stage"
    bucket     = "mmbus-stage"
    create_dir = "true"
    policy_id  = vastdata_view_policy.mmbus-cd.id
    protocols  = ["NFS","S3"]
    bucket_owner = "svcmmdata@mavensecurities.com"
}

resource "vastdata_view" "mmcaptures-raw" {
    path       = "/data/mmdata/mmcaptures-raw"
    bucket     = "mmcaptures-raw"
    create_dir = "true"
    policy_id  = vastdata_view_policy.s3_default_policy.id
    protocols  = ["NFS","S3"]
    bucket_owner = "svcexecution@mavensecurities.com"
}

resource "vastdata_view" "mmdata-importer-archives-prod" {
    path       = "/data/mmdata/importer-archives-prod"
    bucket     = "importer-archives-prod"
    create_dir = "true"
    policy_id  = vastdata_view_policy.mmdata-cd.id
    protocols  = ["NFS","S3"]
    bucket_owner = "svcmmdata@mavensecurities.com"
}

resource "vastdata_view" "mmdata-importer-archives-stage" {
    path       = "/data/mmdata/importer-archives-stage"
    bucket     = "importer-archives-stage"
    create_dir = "true"
    policy_id  = vastdata_view_policy.mmdata-cd.id
    protocols  = ["NFS","S3"]
    bucket_owner = "svcmmdata@mavensecurities.com"
}

resource "vastdata_view" "mmdata-importer-archives-dev" {
    path       = "/data/mmdata/importer-archives-dev"
    bucket     = "importer-archives-dev"
    create_dir = "true"
    policy_id  = vastdata_view_policy.mmdata-cd.id
    protocols  = ["NFS","S3"]
    bucket_owner = "svcmmdata@mavensecurities.com"
}

resource "vastdata_quota" "mmdata" {
  name          = "mmdata"
  default_email = "mmdata@mavensecurities.com"
  path          = vastdata_view.data-mmdata.path
  // Not adding limits here but need a quota so we can track the volume size regardless
  soft_limit    = 0
  hard_limit    = 0
  is_user_quota = false
}

resource "vastdata_quota" "mmdata-importer-archives-dev" {
  name          = "mmdata-importer-archives-dev"
  default_email = "mmdata@mavensecurities.com"
  path          = vastdata_view.mmdata-importer-archives-dev.path
  // limits are in bytes / 30 TB hard 27TB soft
  soft_limit    = 27000000000000
  hard_limit    = 30000000000000
  is_user_quota = false
}

resource "vastdata_quota" "mmdata-importer-archives-stage" {
  name          = "mmdata-importer-archives-stage"
  default_email = "mmdata@mavensecurities.com"
  path          = vastdata_view.mmdata-importer-archives-stage.path
  // limits are in bytes / 30 TB hard 27TB soft
  soft_limit    = 27000000000000
  hard_limit    = 30000000000000
  is_user_quota = false
}

resource "vastdata_quota" "mmdata-importer-archives-prod" {
  name          = "mmdata-importer-archives-prod"
  default_email = "mmdata@mavensecurities.com"
  path          = vastdata_view.mmdata-importer-archives-prod.path
  // limits are in bytes / 30 TB hard 27TB soft
  soft_limit    = 27000000000000
  hard_limit    = 30000000000000
  is_user_quota = false
}
