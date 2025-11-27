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
# VIP pools
############################################
resource "vastdata_vip_pool" "pool1_clusterB" {
  provider    = vastdata.clusterB
  name        = "gateway-1"
  role        = "REPLICATION"
  subnet_cidr = "24"
  ip_ranges   = [["11.0.0.10", "11.0.0.10"]]
}

resource "vastdata_vip_pool" "pool1_clusterA" {
  name        = "gateway-1"
  role        = "REPLICATION"
  subnet_cidr = "24"
  ip_ranges   = [["12.0.0.10", "12.0.0.10"]]
}

############################################
# Replication peer
############################################
resource "vastdata_replication_peer" "clusterA_clusterB_peer" {
  name        = "peer-loop-b"
  leading_vip = vastdata_vip_pool.pool1_clusterB.start_ip
  pool_id     = vastdata_vip_pool.pool1_clusterA.id
  secure_mode = "NONE"
}

############################################
# Protection policy
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
    start_at    = "2027-06-04 09:00:00"
  }]
}

# ---------------------
# Complete examples
# ---------------------

resource "vastdata_tenant" "vastdb_tenant" {
  name             = "vastdb-tenant"
  client_ip_ranges = [["192.168.0.52", "192.168.0.53"]]
}

resource "vastdata_protection_policy" "vastdb_ppolicy" {
  name           = "vastdb-ppolicy"
  indestructible = "false"
  prefix         = "ppolicy"
  clone_type     = "LOCAL"
  tenant_id      = vastdata_tenant.vastdb_tenant.id
  frames = [{
    every      = "1D"
    keep_local = "2W"
    start_at   = "2025-08-01 09:00:00"
    }, {
    every      = "1D"
    keep_local = "8D"
    start_at   = "2025-08-01 09:00:00"
  }]
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

resource "vastdata_protection_policy" "vastdb_ppolicy" {
  name             = "protection-policy"
  clone_type       = "NATIVE_REPLICATION"
  indestructible   = false
  prefix           = "vastdb_ppolicy"
  target_object_id = vastdata_replication_peer.clusterA_clusterB_peer.id

  frames = [{
    every       = "30m"
    keep_local  = "2D"
    keep_remote = "3D"
  }]
}

# --------------------

