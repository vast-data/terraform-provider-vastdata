# ignore:e2e

# =============================
# clusterA
provider "vastdata" {
  host            = "v183"
  username        = "admin"
  password        = "123456"
  skip_ssl_verify = true
}

# clusterB
provider "vastdata" {
  alias           = "clusterB"
  host            = "v184"
  username        = "admin"
  password        = "123456"
  skip_ssl_verify = true
}

############################################
# Cluster A: VIP pool (PROTOCOLS), Tenant, View Policy, View
############################################

data "vastdata_tenant" "tenant_clusterA" {
  name = "default"
}


data "vastdata_tenant" "tenant_clusterB" {
  provider = vastdata.clusterB
  name     = "default"
}


data "vastdata_view_policy" "view_policy_clusterB" {
  provider = vastdata.clusterB
  name     = "default"
}

data "vastdata_view_policy" "view_policy_clusterA" {
  name = "default"
}

resource "vastdata_view" "view_custerA" {
  path       = "/source-1"
  policy_id  = data.vastdata_view_policy.view_policy_clusterA.id
  tenant_id  = data.vastdata_tenant.tenant_clusterA.id
  create_dir = true
}

resource "vastdata_view" "view_custerB" {
  provider   = vastdata.clusterB
  path       = "/destination-1"
  policy_id  = data.vastdata_view_policy.view_policy_clusterB.id
  tenant_id  = data.vastdata_tenant.tenant_clusterB.id
  create_dir = false
}

############################################
# Replication setup (A <-> B)
############################################
data "vastdata_vip_pool" "replication_poolA" {
  name = "gateway-1"
}

data "vastdata_vip_pool" "replication_poolB" {
  provider = vastdata.clusterB
  name     = "gateway-1"
}

resource "vastdata_replication_peer" "clusterA_clusterB_peer" {
  # Peer is created on clusterA
  name        = "clusterA-clusterB-peer"
  leading_vip = data.vastdata_vip_pool.replication_poolB.start_ip
  pool_id     = data.vastdata_vip_pool.replication_poolA.id
}

resource "vastdata_protection_policy" "protection_policy" {
  name             = "protection-policy-1"
  clone_type       = "NATIVE_REPLICATION"
  indestructible   = false
  prefix           = "policy-1"
  target_object_id = vastdata_replication_peer.clusterA_clusterB_peer.id

  frames = [{
    every       = "30m"
    keep_local  = "2D"
    keep_remote = "3D"
  }]
}


############################################
# Protected path (cross-cluster)
############################################
resource "vastdata_protected_path" "protected_path_view" {
  name                 = "protected-path-view"
  source_dir           = vastdata_view.view_custerA.path
  tenant_id            = vastdata_view.view_custerA.tenant_id
  target_exported_dir  = "/destination-uu"
  protection_policy_id = vastdata_protection_policy.protection_policy.id
  remote_tenant_guid   = data.vastdata_tenant.tenant_clusterB.guid
  sync_interval        = 900
}
