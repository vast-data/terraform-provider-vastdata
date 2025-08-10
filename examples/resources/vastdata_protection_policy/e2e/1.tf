# ignore:example

############################################
# Providers
############################################

# clusterB
provider "vastdata" {
  alias           = "clusterB"
  host            = "vast_secondary_host"
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
