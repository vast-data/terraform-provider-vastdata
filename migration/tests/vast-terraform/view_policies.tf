# Copyright (c) HashiCorp, Inc.


## holds shared view policies that will be applicable to a range of views

resource "vastdata_view_policy" "default" {
  name = "default"

  lifecycle {
    ignore_changes = [
      count_views,
    ]
  }

}

resource "vastdata_view_policy" "s3_default_policy" {
  name = "s3_default_policy"

  lifecycle {
    ignore_changes = [
      count_views,
    ]
  }

}

resource "vastdata_view_policy" "root" {
  name          = "root"

  // only allow root admin access (no squash) from cdplatform01
  vip_pools        = [vastdata_vip_pool.prod.id]
  flavor           = "NFS"
  nfs_no_squash    = ["10.66.4.67"]
  nfs_read_write   = ["10.66.4.67"]

  lifecycle {
    ignore_changes = [
      count_views,
    ]
  }

}

resource "vastdata_view_policy" "k8s-csi-cd-prod" {
  name          = "k8s-csi-cd-prod"

  // do not allow access over the dev VIP pool
  vip_pools     = [vastdata_vip_pool.prod.id]
  flavor        = "NFS"

  // nfs_no_squash is needed to allow the k8s csi delete volumes
  nfs_no_squash       = ["10.72.34.0/24"]

  // read write from k8s CD prod
  nfs_read_write   = ["10.72.34.0/24"]

  lifecycle {
    ignore_changes = [
      count_views,
    ]
  }

}

resource "vastdata_view_policy" "k8s-csi-cd-dev" {
  name          = "k8s-csi-cd-dev"

  // nfs_no_squash is needed to allow the k8s csi delete volumes
  nfs_no_squash       = ["10.72.33.0/24"]

  // do not allow access over the prod VIP pool
  vip_pools     = [vastdata_vip_pool.dev.id]
  flavor        = "NFS"

  // read write from k8s CD dev
  nfs_read_write   = ["10.72.33.0/24"]

  lifecycle {
    ignore_changes = [
      count_views,
    ]
  }

}

resource "vastdata_view_policy" "k8s-csi" {
  name          = "k8s-csi"
  nfs_no_squash       = ["10.72.33.0/24"]

  nfs_root_squash       = ["*"]

  // do not allow access over the prod VIP pool
  vip_pools     = [vastdata_vip_pool.dev.id,vastdata_vip_pool.prod.id]
  flavor        = "NFS"

  // read write from k8s CD prod and dev only
  // this should prevent us from mounting a volume into a UK k8s cluster inadvertently
  nfs_read_write   = ["10.72.33.0/24","10.72.34.0/24"]

  lifecycle {
    ignore_changes = [
      count_views,
    ]
  }

}

resource "vastdata_view_policy" "data-worldread" {
  name          = "data-worldread"
  nfs_root_squash       = ["*"]
  s3_read_write   = ["*"]


  // allow over both vip pools
  vip_pools     = [vastdata_vip_pool.dev.id,vastdata_vip_pool.prod.id]
  flavor        = "NFS"
  use_auth_provider = true

  lifecycle {
    ignore_changes = [
      count_views,
    ]
  }

}

resource "vastdata_view_policy" "logs" {
  // allow write access from within k8s, read from everywhere else
  name          = "logs"
  //nfs_root_squash       = ["*"]

  // allow over both vip pools
  vip_pools     = [vastdata_vip_pool.dev.id,vastdata_vip_pool.prod.id]
  flavor        = "NFS"
  use_auth_provider = true
 nfs_read_write   = ["10.72.33.0/24","10.72.34.0/24"]

  lifecycle {
    ignore_changes = [
      count_views,
    ]
  }

}

resource "vastdata_view_policy" "services" {
  name          = "services"
  nfs_root_squash       = ["*"]
  nfs_read_write   = ["*"]
  s3_read_write   = ["*"]
  smb_read_write   = ["*"]

  // allow over both vip pools
  vip_pools     = [vastdata_vip_pool.dev.id,vastdata_vip_pool.prod.id]
  flavor        = "NFS"
  use_auth_provider = true

  lifecycle {
    ignore_changes = [
      count_views,
    ]
  }

}
