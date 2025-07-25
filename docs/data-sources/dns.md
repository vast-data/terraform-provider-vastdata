---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "vastdata_dns Data Source - terraform-provider-vastdata"
subcategory: ""
description: |-
  
---

# vastdata_dns (Data Source)



## Example Usage

```terraform
data "vastdata_dns" "dns1" {
  name = "dns1"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0) The name of the VAST DNS server configuration.

### Read-Only

- `cnode_ids` (List of Number) (Valid for versions: 5.0.0,5.1.0,5.2.0)
- `domain_suffix` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0) A suffix to append to domain names of each virtual IP pool. The suffix should append each domain name to form a valid FQDN for DNS requests to target.
- `enabled` (Boolean) (Valid for versions: 5.0.0,5.1.0,5.2.0) Enables or disables the VAST DNS server configuration.
- `guid` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0) The unique GUID of the VAST DNS server configuration.
- `id` (Number) (Valid for versions: 5.0.0,5.1.0,5.2.0) The unique ID of the VAST DNS server configuration.
- `invalid_name_response` (String) (Valid for versions: 5.1.0,5.2.0) The response DNS type for invalid DNS name. Allowed Values are [NXDOMAIN REFUSED SERVFAIL NOERROR]
- `invalid_type_response` (String) (Valid for versions: 5.1.0,5.2.0) The response DNS type for invalid DNS type. Allowed Values are [NXDOMAIN REFUSED SERVFAIL NOERROR]
- `net_type` (String) (Valid for versions: 5.1.0,5.2.0) The interface that listens for DNS service delegation requests. Allowed Values are [NORTH_PORT SOUTH_PORT EXTERNAL_PORT]
- `ttl` (Number) (Valid for versions: 5.1.0,5.2.0) The response TTL in seconds.
- `vip` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0) The virtual IP for the DNS service. DNS requests from your external DNS server must be delegated to this IP.
- `vip_gateway` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0) The IPv4 address of the gateway to the external DNS server if it is on a different subnet. Must be on the same subnet as the IP and reachable from the relevant network interface.
- `vip_ipv6` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0) The IPv6 address of the DNS service.
- `vip_ipv6_gateway` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0) The IPv6 address of the gateway to the external DNS server if it is on a different subnet.
- `vip_ipv6_subnet_cidr` (Number) (Valid for versions: 5.0.0,5.1.0,5.2.0) The subnet, in the CIDR format, on which the DNS resides. Valid values: [1..128]
- `vip_subnet_cidr` (Number) (Valid for versions: 5.0.0,5.1.0,5.2.0) The subnet, in the CIDR format, on which the DNS resides.
- `vip_vlan` (Number) (Valid for versions: 5.0.0,5.1.0,5.2.0) The VLAN (optional) to enable communication with the external DNS server(s).
