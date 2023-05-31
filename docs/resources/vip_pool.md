---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "vastdata_vip_pool Resource - terraform-provider-vastdata"
subcategory: ""
description: |-
  
---

# vastdata_vip_pool (Resource)



## Example Usage

```terraform
resource vastdata_vip_pool pool1{
 name = "pool1"
 role = "PROTOCOLS"
 subnet_cidr = "24"
 ip_ranges {
        end_ip = "11.0.0.40"
        start_ip = "11.0.0.20"
  }

 ip_ranges {
        start_ip = "11.0.0.5"
        end_ip = "11.0.0.10"
  }
}

resource vastdata_tenant tenant1 {
 name = "tenant01"
 client_ip_ranges {
         start_ip = "192.168.0.100"
         end_ip = "192.168.0.200"
    }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `ip_ranges` (Block List, Min: 1) (see [below for nested schema](#nestedblock--ip_ranges))
- `name` (String)
- `role` (String)
- `subnet_cidr` (Number)

### Optional

- `active_interfaces` (Number) Numver of active interfaces
- `cluster` (String) Parent Cluster
- `cnode_ids` (List of Number) IDs of cnodes comprising cnode group
- `domain_name` (String)
- `enable_l3` (Boolean) Enables L3 CNode access
- `enabled` (Boolean) True for enable, False for disable
- `gw_ip` (String) Gateway IP Address
- `gw_ipv6` (String) GW IPv6 Address
- `peer_asn` (Number) Peer ASN
- `port_membership` (String) The port on the CNode this pool will use. Right, left or all
- `state` (String)
- `subnet_cidr_ipv6` (Number) IPv6 Subnet CIDR prefix (bits number)
- `url` (String)
- `vast_asn` (Number) VAST ASN
- `vlan` (Number) VIPPool VLAN
- `vms_preferred` (Boolean) If true, CNodes participating in the vip pool are preferred in VMS host election

### Read-Only

- `guid` (String) A uniq guid given to the vippool
- `id` (String) The ID of this resource.

<a id="nestedblock--ip_ranges"></a>
### Nested Schema for `ip_ranges`

Optional:

- `end_ip` (String)
- `start_ip` (String)