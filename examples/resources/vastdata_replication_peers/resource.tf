#Suppose that two providers are defined, one for each VAST cluster:
#1. A provider with the alias `clusterA`
#2. A provider with the alias `clusterB`

#Define replication virtual IP for cluster A:
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

#Define replication virtual IP pool for cluster B:
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
#Define a replication peer on cluster A using virtual IP pool settings from cluster B:
resource "vastdata_replication_peers" "clusterA-clusterB-peer" {
  name        = "peer-loop-b"
  leading_vip = vastdata_vip_pool.pool1-clusterB.ip_ranges[0].start_ip
  pool_id     = vastdata_vip_pool.pool1-clusterA.id
  secure_mode = "NONE"
  provider    = vastdata.clusterA
}
