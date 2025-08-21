# Copyright (c) HashiCorp, Inc.


resource "vastdata_view_policy" "strike-pnl-cd" {
  name          = "strike-pnl-cd"

  // only allow access via productrion networks
  vip_pools     = [vastdata_vip_pool.prod.id]
  flavor        = "NFS"
  nfs_root_squash       = ["*"]
  // access from US sources only, second CIDR is k8s prod, third is k8s dev
  nfs_read_write   = ["10.66.0.0/16","10.72.34.0/24","10.72.33.0/24"]

  lifecycle {
    ignore_changes = [
      count_views,
    ]
  }

}

resource "vastdata_view" "data-strikepnl" {
  path       = "/data/strikepnl" ## group_owner:sgdatastrikepnl
  policy_id  = vastdata_view_policy.strike-pnl-cd.id
  create_dir = "true"
  protocols  = ["NFS"]
}

// add a bucket
resource "vastdata_view" "strike-pnl-1min-tick-data" {
    path       = "/data/strikepnl/strike-pnl-1min-tick-data"  ## group_owner:sgdatastrikepnl
    bucket     = "strike-pnl-1min-tick-data"
    create_dir = "true"
    policy_id  = vastdata_view_policy.s3_default_policy.id
    protocols  = ["NFS","S3"]
    bucket_owner = "svcPTQ@mavensecurities.com"
}


resource "vastdata_view" "strike-pnl-1min-tick-data-tmp" {
    path       = "/data/strikepnl/strike-pnl-1min-tick-data-tmp"  ## group_owner:sgdatastrikepnl
    bucket     = "strike-pnl-1min-tick-data-tmp"
    create_dir = "true"
    policy_id  = vastdata_view_policy.s3_default_policy.id
    protocols  = ["NFS","S3"]
    bucket_owner = "svcPTQ@mavensecurities.com"
}

resource "vastdata_view" "strike-pnl-backfill" {
    path       = "/data/strikepnl/strike-pnl-backfill"  ## group_owner:sgdatastrikepnl
    bucket     = "strike-pnl-backfill"
    create_dir = "true"
    policy_id  = vastdata_view_policy.s3_default_policy.id
    protocols  = ["NFS","S3"]
    bucket_owner = "svcmmdata@mavensecurities.com"
}

resource "vastdata_view" "strike-pnl-backfill-dev" {
    path       = "/data/strikepnl/strike-pnl-backfill-dev"  ## group_owner:sgdatastrikepnl
    bucket     = "strike-pnl-backfill-dev"
    create_dir = "true"
    policy_id  = vastdata_view_policy.s3_default_policy.id
    protocols  = ["NFS","S3"]
    bucket_owner = "svcmmdata@mavensecurities.com"
}

resource "vastdata_view" "strike-pnl-backfill-prod" {
    path       = "/data/strikepnl/strike-pnl-backfill-prod"  ## group_owner:sgdatastrikepnl
    bucket     = "strike-pnl-backfill-prod"
    create_dir = "true"
    policy_id  = vastdata_view_policy.s3_default_policy.id
    protocols  = ["NFS","S3"]
    bucket_owner = "svcmmdata@mavensecurities.com"
}

resource "vastdata_view" "strike-pnl-backfill-stage" {
    path       = "/data/strikepnl/strike-pnl-backfill-stage"  ## group_owner:sgdatastrikepnl
    bucket     = "strike-pnl-backfill-stage"
    create_dir = "true"
    policy_id  = vastdata_view_policy.s3_default_policy.id
    protocols  = ["NFS","S3"]
    bucket_owner = "svcmmdata@mavensecurities.com"
}

resource "vastdata_view" "strike-pnl-backfill-tmp" {
    path       = "/data/strikepnl/strike-pnl-backfill-tmp"  ## group_owner:sgdatastrikepnl
    bucket     = "strike-pnl-backfill-tmp"
    create_dir = "true"
    policy_id  = vastdata_view_policy.s3_default_policy.id
    protocols  = ["NFS","S3"]
    bucket_owner = "svcmmdata@mavensecurities.com"
}

resource "vastdata_quota" "strikepnl" {
  name          = "strikepnl"
  default_email = "michael.moyles@mavensecurities.com"
  path          = vastdata_view.data-strikepnl.path
    // limits are in bytes: 500GB
  soft_limit    = 500000000000
  hard_limit    = 500000000000
  is_user_quota = false
}
