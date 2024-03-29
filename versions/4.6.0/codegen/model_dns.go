/*
 * VAST API Swagger Schema
 *
 * VAST Management API definition
 *
 * API version: 1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type Dns struct {
	// Specifies a name for the VAST DNS server configuration
	Name string `json:"name,omitempty"`
	//
	Id int32 `json:"id,omitempty"`
	// Assigns a IP to the DNS service. DNS requests from your external DNS server must be delegated to this IP.
	Vip string `json:"vip,omitempty"`
	// Specifies a suffix to append to domain names of each VIP pool. The suffix should complete each domain name to form a valid FQDN for DNS requests to target.
	DomainSuffix string `json:"domain_suffix,omitempty"`
	// Specifies a gateway IP to external DNS server if on different subnet. Must be on same subnet as the IP and reachable from the relevant nework interface.
	VipGateway string `json:"vip_gateway,omitempty"`
	//
	Enabled bool `json:"enabled,omitempty"`
	//
	Guid string `json:"guid,omitempty"`
	//
	NetType string `json:"net_type,omitempty"`
	//
	InvalidNameResponse string `json:"invalid_name_response,omitempty"`
	//
	InvalidTypeResponse string `json:"invalid_type_response,omitempty"`
	// Specifies the subnet, as a CIDR index, on which the DNS resides.
	VipSubnetCidr int32 `json:"vip_subnet_cidr,omitempty"`
	// Specifies a VLAN if needed to enable communication with external DNS server(s).
	VipVlan  int32   `json:"vip_vlan,omitempty"`
	CnodeIds []int32 `json:"cnode_ids,omitempty"`
	// Synchronization state with leader
	Sync string `json:"sync,omitempty"`
	// Synchronization time with leader
	SyncTime string `json:"sync_time,omitempty"`
	// Assigns an IPv6 to the DNS service.
	VipIpv6 string `json:"vip_ipv6,omitempty"`
	// Specifies the subnet, as a CIDR index, on which the DNS resides. [1..128]
	VipIpv6SubnetCidr int32 `json:"vip_ipv6_subnet_cidr,omitempty"`
	// Specifies a gateway IPv6 to external DNS server if on different subnet.
	VipIpv6Gateway string `json:"vip_ipv6_gateway,omitempty"`
}
