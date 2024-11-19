/*
 * VastData API
 *
 * A API document representing VastData API
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger

type VipPool struct {
	// A uniq id given to the vippool
	Id int32 `json:"id,omitempty"`
	// A uniq guid given to the vippool
	Guid string `json:"guid,omitempty"`
	// A uniq name given to the vippool
	Name string `json:"name,omitempty"`
	// IPv4 Subnet CIDR prefix (bits number)
	SubnetCidr int32 `json:"subnet_cidr,omitempty"`
	// IPv6 Subnet CIDR prefix (bits number)
	SubnetCidrIpv6 int32 `json:"subnet_cidr_ipv6,omitempty"`
	// Gateway IP Address
	GwIp string `json:"gw_ip,omitempty"`
	// GW IPv6 Address
	GwIpv6 string `json:"gw_ipv6,omitempty"`
	// VIPPool VLAN
	Vlan int `json:"vlan,omitempty"`
	State string `json:"state,omitempty"`
	// IDs of cnodes comprising cnode group
	CnodeIds []int32 `json:"cnode_ids,omitempty"`
	// Parent Cluster
	Cluster string `json:"cluster,omitempty"`
	Url string `json:"url,omitempty"`
	DomainName string `json:"domain_name,omitempty"`
	// Role
	Role string `json:"role,omitempty"`
	// IP ranges
	IpRanges [][]string `json:"ip_ranges,omitempty"`
	// If true, CNodes participating in the vip pool are preferred in VMS host election
	VmsPreferred bool `json:"vms_preferred,omitempty"`
	// True for enable, False for disable
	Enabled bool `json:"enabled,omitempty"`
	// The port on the CNode this pool will use. Right, left or all
	PortMembership string `json:"port_membership,omitempty"`
	// Numver of active interfaces
	ActiveInterfaces int32 `json:"active_interfaces,omitempty"`
	// Enables L3 CNode access
	EnableL3 bool `json:"enable_l3,omitempty"`
	// VAST ASN
	VastAsn int32 `json:"vast_asn,omitempty"`
	// Peer ASN
	PeerAsn int32 `json:"peer_asn,omitempty"`
	// The Tenant id to which this Vip Pool is assigned to , if set to 0 it means all tenants 
	TenantId int64 `json:"tenant_id,omitempty"`
}
