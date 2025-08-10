
############################################
# Providers
############################################
# clusterA
provider "vastdata" {
  host            = "v95"
  username        = "admin"
  password        = "123456"
  skip_ssl_verify = true
}

# clusterB
provider "vastdata" {
  alias           = "clusterB"
  host            = "v117"
  username        = "admin"
  password        = "123456"
  skip_ssl_verify = true
}

############################################
# Cluster A: VIP pool (PROTOCOLS), Tenant, View Policy, View
############################################
resource "vastdata_vip_pool" "protocols_pool" {
  name        = "protocols-pool"
  role        = "PROTOCOLS"
  subnet_cidr = "24"
  # Framework expects list-of-tuples form
  ip_ranges = [
    ["11.0.0.70", "11.0.0.80"]
  ]
}

resource "vastdata_tenant" "tenant" {
  name = "tenant1"
}

resource "vastdata_view_policy" "view_policy" {
  name          = "view-policy1"
  vip_pools     = [vastdata_vip_pool.protocols_pool.id]
  tenant_id     = vastdata_tenant.tenant.id
  flavor        = "NFS"
  nfs_no_squash = ["10.0.0.1", "10.0.0.2"]
}

resource "vastdata_view" "view" {
  path       = "/view1"
  policy_id  = vastdata_view_policy.view_policy.id
  tenant_id  = vastdata_tenant.tenant.id
  create_dir = true
}

############################################
# Replication setup (A <-> B)
############################################
resource "vastdata_vip_pool" "replication_poolA" {
  name        = "gateway-1"
  role        = "REPLICATION"
  subnet_cidr = "24"
  ip_ranges   = [["12.0.0.10", "12.0.0.10"]]
}

resource "vastdata_vip_pool" "replication_poolB" {
  provider    = vastdata.clusterB
  name        = "gateway-1"
  role        = "REPLICATION"
  subnet_cidr = "24"
  ip_ranges   = [["11.0.0.10", "11.0.0.10"]]
}

resource "vastdata_replication_peer" "clusterA_clusterB_peer" {
  # Peer is created on clusterA
  name        = "clusterA-clusterB-peer"
  password    = "####Wwww11111"
  leading_vip = vastdata_vip_pool.replication_poolB.start_ip
  pool_id     = vastdata_vip_pool.replication_poolA.id
}

############################################
# Protection policy (Framework-style frames map)
############################################
resource "vastdata_protection_policy" "protection_policy" {
  name             = "protection-policy-1"
  clone_type       = "NATIVE_REPLICATION"
  indestructible   = false
  prefix           = "policy-1"
  target_object_id = vastdata_replication_peer.clusterA_clusterB_peer.id

  frames = [{
    every       = "1D"
    keep_local  = "2D"
    keep_remote = "3D"
    start_at    = "2026-06-04 09:00:00"
  }]
}

############################################
# Cluster B: Tenant
############################################
resource "vastdata_tenant" "tenant_clusterB" {
  provider = vastdata.clusterB
  name     = "tenant2"
}

############################################
# Protected path (cross-cluster)
############################################
resource "vastdata_protected_path" "protected_path_view" {
  name                 = "protected-path-view"
  source_dir           = vastdata_view.view.path
  tenant_id            = vastdata_view.view.tenant_id
  target_exported_dir  = "/view1"
  protection_policy_id = vastdata_protection_policy.protection_policy.id
  remote_tenant_guid   = vastdata_tenant.tenant_clusterB.guid
}

