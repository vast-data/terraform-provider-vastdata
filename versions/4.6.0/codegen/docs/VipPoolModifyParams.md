# VipPoolModifyParams

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | The VIP pool name | [optional] [default to null]
**StartIp** | **string** | The IP address at the start of a continuous range. | [optional] [default to null]
**EndIp** | **string** | The IP address at the end of the range. | [optional] [default to null]
**SubnetCidr** | **int32** | The subnet expressed as a CIDR index (number of bits in each IP that belong to the subnet) | [optional] [default to null]
**SubnetCidrIpv6** | **int32** | The subnet expressed as a CIDR index (number of bits in each IP that belong to the subnet) | [optional] [default to null]
**GwIp** | **string** | The IP address of a local gateway device if client traffic is routed through one | [optional] [default to null]
**GwIpv6** | **string** | The IP address of a local gateway device if client traffic is routed through one | [optional] [default to null]
**Vlan** | **int32** | To tag the VIP pool with a specific VLAN on the data network, specify the VLAN (0-4096). | [optional] [default to null]
**CnodeIds** | **string** | To dedicate a specific group of CNodes to the VIP pool, list the IDs of the CNodes. | [optional] [default to null]
**CnodeNames** | **string** | list of cnode names | [optional] [default to null]
**DomainName** | **string** | Domain name for the VAST DNS server. The domain suffix defined in the DNS server configuration is appended to this domain name to form a FQDN which the DNS server resolves to this VIP pool. | [optional] [default to null]
**Role** | **string** | &#39;Protocol&#39; dedicates the VIP pool for client access. &#39;Replication&#39; dedicates the VIP pool for native replication | [optional] [default to null]
**IpRanges** | [**[][]string**](array.md) | IP ranges | [optional] [default to null]
**TenantId** | **int32** | Tenant ID | [optional] [default to null]
**VmsPreferred** | **bool** | If true, CNodes participating in the vip pool to be preferred in VMS host election. | [optional] [default to null]
**Enabled** | **bool** | True for enable, False for disable | [optional] [default to null]
**PortMembership** | **string** | The port on the CNode this pool will use. Right, left or all | [optional] [default to null]
**VastAsn** | **int32** |  | [optional] [default to null]
**PeerAsn** | **int32** |  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


