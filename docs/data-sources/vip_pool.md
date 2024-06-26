---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "vastdata_vip_pool Data Source - terraform-provider-vastdata"
subcategory: ""
description: |-
  
---

# vastdata_vip_pool (Data Source)



## Example Usage

```terraform
data "vastdata_vip_pool" "pool1" {
  name = "pool1"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) A uniq name given to the vippool

### Read-Only

- `active_interfaces` (Number) Numver of active interfaces
- `cluster` (String) Parent Cluster
- `cnode_ids` (List of Number) IDs of cnodes comprising cnode group
- `domain_name` (String)
- `enable_l3` (Boolean) Enables L3 CNode access
- `enabled` (Boolean) True for enable, False for disable
- `guid` (String) A uniq guid given to the vippool
- `gw_ip` (String) Gateway IP Address
- `gw_ipv6` (String) GW IPv6 Address
- `id` (Number) A uniq id given to the vippool
- `ip_ranges` (List of Object) IP ranges (see [below for nested schema](#nestedatt--ip_ranges))
- `peer_asn` (Number) Peer ASN
- `port_membership` (String) The port on the CNode this pool will use. Right, left or all
- `role` (String) Role
- `state` (String)
- `subnet_cidr` (Number) IPv4 Subnet CIDR prefix (bits number)
- `subnet_cidr_ipv6` (Number) IPv6 Subnet CIDR prefix (bits number)
- `url` (String)
- `vast_asn` (Number) VAST ASN
- `vlan` (Number) VIPPool VLAN
- `vms_preferred` (Boolean) If true, CNodes participating in the vip pool are preferred in VMS host election

<a id="nestedatt--ip_ranges"></a>
### Nested Schema for `ip_ranges`

Read-Only:

- `end_ip` (String)
- `start_ip` (String)
