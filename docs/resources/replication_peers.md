---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "vastdata_replication_peers Resource - terraform-provider-vastdata"
subcategory: ""
description: |-
  
---

# vastdata_replication_peers (Resource)



## Example Usage

```terraform
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0) The name of the replication peer configuration.

### Optional

- `is_local` (Boolean) (Valid for versions: 5.0.0,5.1.0,5.2.0) Specifies whether the source of the replication is local (this host is the source).
- `leading_vip` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0) The virtual IP pool provided for the replication peer configuration.
- `peer_name` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0) The name of the peer cluster.
- `pool_id` (Number) (Valid for versions: 5.0.0,5.1.0,5.2.0) The ID of the replication virtual IP pool.
- `remote_version` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0) The version of the remote peer.
- `remote_vip_range` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0) The range of virtual IPs that were reported by the peer.
- `secure_mode` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0) If true, the connection is secure.
- `url` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0) Direct URL of the replication peer configuration.
- `version` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0) The version of the source.

### Read-Only

- `guid` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0) The unique GUID of the replication peer configuration.
- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import vastdata_replication_peers.example <guid>
terraform import vastdata_replication_peers.example <Name>
```
