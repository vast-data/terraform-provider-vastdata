
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
  name    = "tf-test-cluster"
  netmask = "255.255.255.0"

  static_ip_ranges = [
    ["172.21.87.50", "172.21.87.54"],
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

  get_nodes       = true
  get_pods        = true
  get_namespaces  = true
  get_services    = true
  get_deployments = true
  get_tenants     = true
  get_dashboard   = true
}
