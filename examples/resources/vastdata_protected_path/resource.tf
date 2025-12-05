
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


# ---------------------
# Complete examples
# ---------------------

resource "vastdata_tenant" "vastdb_tenant" {
  name             = "vastdbtenant"
  client_ip_ranges = [["192.168.0.50", "192.168.0.51"]]
}


resource "vastdata_protection_policy" "vastdb_ppolicy" {
  name           = "vastdb_ppolicy"
  indestructible = "false"
  prefix         = "vastdb-policy"
  clone_type     = "LOCAL"
  tenant_id      = vastdata_tenant.vastdb_tenant.id
  frames = [{
    every      = "1D"
    keep_local = "14D"
    start_at   = "2025-08-01 09:00:00"
    }, {
    every      = "1D"
    keep_local = "8D"
    start_at   = "2025-08-01 09:00:00"
  }]
}


resource "vastdata_protected_path" "protected_path" {
  name                 = "vastdbppath"
  source_dir           = "/"
  tenant_id            = vastdata_tenant.vastdb_tenant.id
  target_exported_dir  = "/view1"
  protection_policy_id = vastdata_protection_policy.vastdb_ppolicy.id
  enabled              = false
  capabilities         = "ASYNC_REPLICATION"
}

# --------------------


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

# --------------------

