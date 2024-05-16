#We have defined 2 clusters called clusterA, clusterB.
#We define all needed componants in order to define the proteced path
# this will create a replication of view with path /view1 from clusterA to remote clusterB path /view1
resource "vastdata_vip_pool" "pool" {
  provider    = vastdata.clusterA
  name        = "protocols-pool"
  role        = "PROTOCOLS"
  subnet_cidr = "24"
  ip_ranges {
    end_ip   = "11.0.0.80"
    start_ip = "11.0.0.70"
  }

}

resource "vastdata_view_policy" "view-policy" {
  provider      = vastdata.clusterA
  name          = "view-policy1"
  vip_pools     = [vastdata_vip_pool.pool.id]
  tenant_id     = vastdata_tenant.tenant.id
  flavor        = "NFS"
  nfs_no_squash = ["10.0.0.1", "10.0.0.2"]
}

resource "vastdata_view" "view" {
  path       = "/view1"
  policy_id  = vastdata_view_policy.view-policy.id
  tenant_id  = vastdata_tenant.tenant.id
  create_dir = "true"
  provider   = vastdata.clusterA
}

resource "vastdata_protection_policy" "protection-policy" {
  provider         = vastdata.clusterA
  name             = "protection-policy-1"
  clone_type       = "NATIVE_REPLICATION"
  indestructible   = "false"
  prefix           = "policy-1"
  target_object_id = vastdata_replication_peers.clusterA-clusterB-peer.id
  frames {
    every       = "1D"
    keep_local  = "2D"
    keep_remote = "3D"
    start_at    = "2023-06-04 09:00:00"
  }


}

resource "vastdata_tenant" "tenant-clusterB" {
  name     = "tenant2"
  provider = vastdata.clusterB
}

resource "vastdata_protected_path" "protected-path-view" {
  name                 = "protected-path-view"
  source_dir           = vastdata_view.view.path
  tenant_id            = vastdata_view.view.tenant_id
  target_exported_dir  = "/view1"
  protection_policy_id = vastdata_protection_policy.protection-policy.id
  remote_tenant_guid   = vastdata_tenant.tenant-clusterB.guid
  target_id            = vastdata_protection_policy.protection-policy.target_object_id

}
