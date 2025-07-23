resource "vastdata_vip_pool" "vastdb_replication_poolA" {
  name        = "vastdb_replication_poolA"
  role        = "REPLICATION"
  subnet_cidr = "24"
  provider    = vastdata.clusterA
  ip_ranges = [
    ["12.0.0.10", "12.0.0.10"]
  ]
}

resource "vastdata_vip_pool" "vastdb_replication_poolB" {
  name        = "vastdb_replication_poolB"
  role        = "REPLICATION"
  subnet_cidr = "24"
  provider    = vastdata.clusterB
  ip_ranges = [
    ["11.0.0.10", "11.0.0.10"]
  ]
}

resource "vastdata_replication_peer" "vastdb_replication_peer" {
  name        = "vastdb_replication_peer"
  password    = "####Wwww11111"
  leading_vip = vastdata_vip_pool.vastdb_replication_poolB.start_ip
  pool_id     = vastdata_vip_pool.vastdb_replication_poolA.id
  provider    = vastdata.clusterA
}
