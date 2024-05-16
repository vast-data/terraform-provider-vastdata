#We defined 2 porivders , each one represents a cluster , clusterA & clusterB.
#We define vip pools for replication for each cluster & make them replication peers.
#Than we define a protection policy
resource "vastdata_vip_pool" "pool1-clusterB" {
  name        = "pool1"
  role        = "REPLICATION"
  subnet_cidr = "24"
  provider    = vastdata.clusterB
  ip_ranges {
    end_ip   = "11.0.0.10"
    start_ip = "11.0.0.10"
  }

}

resource "vastdata_vip_pool" "pool1-clusterA" {
  name        = "pool1"
  role        = "REPLICATION"
  subnet_cidr = "24"
  provider    = vastdata.clusterA
  ip_ranges {
    end_ip   = "12.0.0.10"
    start_ip = "12.0.0.10"
  }

}

resource "vastdata_replication_peers" "clusterA-clusterB-peer" {
  name        = "peer-loop-b"
  leading_vip = vastdata_vip_pool.pool1-clusterB.ip_ranges[0].start_ip
  pool_id     = vastdata_vip_pool.pool1-clusterA.id
  secure_mode = "NONE"
  provider    = vastdata.clusterA
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
