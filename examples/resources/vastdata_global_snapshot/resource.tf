#we have 2 providers defined for each cluster with the aliases clusterA & clusterB
#we create the following resources in order to create a global snapshot
#1. Replication vippool for ceach cluster
#2. Define them as replication  peers
#3. A tenant 
#4. A view with the prefix /view1 whcih belongs to the previous tenant just created
#5  A snapshot to the view called snapshot1 on clusterB
#6  A Global Snapshot of snapshot1 from clusterB to clusterA


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
resource "vastdata_replication_peers" "clusterA-clusterB-peer" {
  name        = "peer-loop-b"
  leading_vip = vastdata_vip_pool.pool1-clusterB.ip_ranges[0].start_ip
  pool_id     = vastdata_vip_pool.pool1-clusterA.id
  secure_mode = "NONE"
  provider    = vastdata.clusterA
}

resource "vastdata_view" "view" {
  provider   = vastdata.clusterB
  path       = "/view1"
  policy_id  = vastdata_view_policy.view-policy.id
  tenant_id  = vastdata_tenant.tenant.id
  create_dir = "true"
}

resource "vastdata_tenant" "tenant" {
  provider = vastdata.clusterB
  name     = "tenant1"
  client_ip_ranges {
    start_ip = "192.168.0.100"
    end_ip   = "192.168.0.200"
  }
}

resource "vastdata_snapshot" "snapshot" {
  provider        = vastdata.clusterB
  name            = "snapshot1"
  path            = vastdata_view.view.path
  tenant_id       = vastdata_tenant.tenant.id
  indestructible  = false
  expiration_time = "2023-06-20T12:22:32Z"
  lifecycle {
    ignore_changes = [path]
  }

}


resource "vastdata_global_snapshot" "gsnap1" {
  name               = "gsnap1"
  enabled            = true
  provider           = vastdata.clusterA
  loanee_root_path   = "/gsnap1"
  remote_target_path = "/view1/"
  loanee_tenant_id   = 1
  remote_target_guid = vastdata_replication_peers.clusterA-clusterB-peer.guid
  remote_target_id   = vastdata_replication_peers.clusterA-clusterB-peer.id
  owner_root_snapshot {
    name = vastdata_snapshot.snapshot.name

  }
  owner_tenant {
    guid = vastdata_tenant.tenant.guid
    name = vastdata_tenant.tenant.name
  }
  lifecycle {
    ignore_changes = [remote_target_guid, remote_target_path, loanee_root_path]
  }

}
