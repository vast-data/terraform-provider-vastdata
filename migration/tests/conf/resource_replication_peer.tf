# Copyright (c) HashiCorp, Inc.

variable peer_name {
    type = string
}

resource vastdata_vip_pool local-rep-pool1 {
    name = "tfpool1"
    role = "REPLICATION"
    subnet_cidr = "24"
    ip_ranges {
        start_ip = "172.20.0.100"
        end_ip = "172.20.0.101"
    }
}

resource vastdata_replication_peers clusterA-clusterB-peer {
    name=var.peer_name
    leading_vip="172.20.0.2"
    pool_id=vastdata_vip_pool.local-rep-pool1.id
    secure_mode="NONE"
}

output tf_replication_peer {
  value = vastdata_replication_peers.clusterA-clusterB-peer
}