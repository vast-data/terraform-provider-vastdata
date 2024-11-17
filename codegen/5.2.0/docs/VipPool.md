# VipPool

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **int32** | A uniq id given to the vippool | [optional] [default to null]
**Guid** | **string** | A uniq guid given to the vippool | [optional] [default to null]
**Name** | **string** | A uniq name given to the vippool | [optional] [default to null]
**SubnetCidr** | **int32** | IPv4 Subnet CIDR prefix (bits number) | [optional] [default to null]
**SubnetCidrIpv6** | **int32** | IPv6 Subnet CIDR prefix (bits number) | [optional] [default to null]
**GwIp** | **string** | Gateway IP Address | [optional] [default to null]
**GwIpv6** | **string** | GW IPv6 Address | [optional] [default to null]
**Vlan** | **int** | VIPPool VLAN | [optional] [default to null]
**State** | **string** |  | [optional] [default to null]
**CnodeIds** | **[]int32** | IDs of cnodes comprising cnode group | [optional] [default to null]
**Cluster** | **string** | Parent Cluster | [optional] [default to null]
**Url** | **string** |  | [optional] [default to null]
**DomainName** | **string** |  | [optional] [default to null]
**Role** | **string** | Role | [optional] [default to null]
**IpRanges** | [**[][]string**](array.md) | IP ranges | [optional] [default to null]
**VmsPreferred** | **bool** | If true, CNodes participating in the vip pool are preferred in VMS host election | [optional] [default to null]
**Enabled** | **bool** | True for enable, False for disable | [optional] [default to null]
**PortMembership** | **string** | The port on the CNode this pool will use. Right, left or all | [optional] [default to null]
**ActiveInterfaces** | **int32** | Numver of active interfaces | [optional] [default to null]
**EnableL3** | **bool** | Enables L3 CNode access | [optional] [default to null]
**VastAsn** | **int32** | VAST ASN | [optional] [default to null]
**PeerAsn** | **int32** | Peer ASN | [optional] [default to null]
**TenantId** | **int64** | The Tenant id to which this Vip Pool is assigned to , if set to 0 it means all tenants  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

