# ignore:example

########################################
# Providers
########################################

# Cluster B
provider "vastdata" {
  alias           = "clusterB"
  host            = "vast_secondary_host"
  port            = 443
  username        = "admin"
  password        = "123456"
  skip_ssl_verify = true
}

########################################
# VIP Pools
########################################

resource "vastdata_vip_pool" "pool1_clusterA" {
  provider    = vastdata.clusterA
  name        = "gateway-1"
  role        = "REPLICATION"
  subnet_cidr = 24
  ip_ranges   = [["12.0.0.10", "12.0.0.10"]]
}

resource "vastdata_vip_pool" "pool1_clusterB" {
  provider    = vastdata.clusterB
  name        = "gateway-1"
  role        = "REPLICATION"
  subnet_cidr = 24
  ip_ranges   = [["11.0.0.10", "11.0.0.10"]]
}

locals {
  pool1_clusterB_leading_vip = vastdata_vip_pool.pool1_clusterB.ip_ranges[0][0]
}

########################################
# Replication peer (Cluster A -> B)
########################################

resource "vastdata_replication_peer" "clusterA_clusterB_peer" {
  provider    = vastdata.clusterA
  name        = "peer-loop-b"
  leading_vip = local.pool1_clusterB_leading_vip
  pool_id     = vastdata_vip_pool.pool1_clusterA.id
  secure_mode = "NONE"
}

########################################
# Tenant and View on Cluster B
########################################

data "vastdata_view_policy" "view_policyB" {
  provider = vastdata.clusterB
  name     = "default"
}

resource "vastdata_tenant" "tenant1" {
  provider         = vastdata.clusterB
  name             = "tenant1"
  client_ip_ranges = [["192.168.0.100", "192.168.0.200"]]
}

resource "vastdata_view" "view1" {
  provider   = vastdata.clusterB
  path       = "/view1"
  policy_id  = data.vastdata_view_policy.view_policyB.id
  tenant_id  = vastdata_tenant.tenant1.id
  create_dir = true
}

########################################
# Snapshot on Cluster B
########################################

resource "vastdata_snapshot" "snapshot1" {
  provider        = vastdata.clusterB
  name            = "snapshot1"
  path            = "${vastdata_view.view1.path}/"
  tenant_id       = vastdata_tenant.tenant1.id
  indestructible  = false
  expiration_time = "2023-06-20T12:22:32Z"

  lifecycle {
    ignore_changes = [path]
  }
}

########################################
# Global Snapshot from Cluster A
########################################

resource "vastdata_global_snapshot" "gsnap1" {
  provider         = vastdata.clusterA
  name             = "gsnap1"
  enabled          = true
  loanee_root_path = "/gsnap1"

  loanee_tenant_id = 1
  remote_target_id = vastdata_replication_peer.clusterA_clusterB_peer.id

  owner_root_snapshot = {
    name = vastdata_snapshot.snapshot1.name
  }
}