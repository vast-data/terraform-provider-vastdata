# ignore:e2e - compute cluster provisioning requires dedicated CNodes and static IP ranges

# Requires at least 3 unassigned CNodes. Adjust IDs to match your cluster.

data "vastdata_cnode" "vastdb_cnode_1" {
  id = 1
}

data "vastdata_cnode" "vastdb_cnode_2" {
  id = 2
}

data "vastdata_cnode" "vastdb_cnode_3" {
  id = 3
}

resource "vastdata_compute_cluster" "vastdb_compute_cluster" {
  name    = "vastdb-compute-cluster"
  netmask = "24"

  static_ip_ranges = [
    ["10.100.0.10", "10.100.0.20"],
  ]

  cnodes = [
    {
      id              = data.vastdata_cnode.vastdb_cnode_1.id
      resource_preset = "BALANCED"
    },
    {
      id              = data.vastdata_cnode.vastdb_cnode_2.id
      resource_preset = "BALANCED"
    },
    {
      id              = data.vastdata_cnode.vastdb_cnode_3.id
      resource_preset = "BALANCED"
    },
  ]

  description = "E2E compute cluster example"
  get_nodes   = true
  get_pods    = true
}
