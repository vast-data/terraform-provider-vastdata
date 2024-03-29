/*
 * VAST API Swagger Schema
 *
 * VAST Management API definition
 *
 * API version: 1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type DnsModifyParams struct {
	// 
	Name string `json:"name,omitempty"`
	// A virtual IP to assign to the DNS service. DNS requests from your external DNS server must be delegated to this IP.
	Vip string `json:"vip,omitempty"`
	// A suffix for domain names. Requests for domain names with this suffix are resolved to the VIPs configured on the cluster.
	DomainSuffix string `json:"domain_suffix,omitempty"`
	// If the external DNS server doesn't reside on the same subnet as the DNS VIP, enter the IP of a gateway through which to connect to the DNS server. 
	VipGateway string `json:"vip_gateway,omitempty"`
	// Set to true to enable the DNS service
	Enabled bool `json:"enabled,omitempty"`
	// 
	NetType string `json:"net_type,omitempty"`
	// 
	InvalidNameResponse string `json:"invalid_name_response,omitempty"`
	// 
	InvalidTypeResponse string `json:"invalid_type_response,omitempty"`
	// The subnet, in CIDR format, on which the DNS VIP resides.
	VipSubnetCidr int32 `json:"vip_subnet_cidr,omitempty"`
	// To dedicate a specific group of CNodes to the DNS, list the IDs of the CNodes.
	CnodeIds string `json:"cnode_ids,omitempty"`
	// If your external DNS server is only exposed to a specific VLAN, you can enter the VLAN here to enable communication with the DNS server.
	VipVlan int32 `json:"vip_vlan,omitempty"`
	// Assigns an IPv6 to the DNS service.
	VipIpv6 string `json:"vip_ipv6,omitempty"`
	// Specifies the subnet, as a CIDR index, on which the DNS resides. [1..128]
	VipIpv6SubnetCidr int32 `json:"vip_ipv6_subnet_cidr,omitempty"`
	// Specifies a gateway IPv6 to external DNS server if on different subnet.
	VipIpv6Gateway string `json:"vip_ipv6_gateway,omitempty"`
}
