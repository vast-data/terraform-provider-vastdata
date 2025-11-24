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
