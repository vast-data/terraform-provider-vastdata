# VipPool

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **int32** | The unique ID of the virtual IP pool. | [optional] [default to null]
**Guid** | **string** | The unique GUID of the virtual IP pool. | [optional] [default to null]
**Name** | **string** | The unique name of the virtual IP pool. | [optional] [default to null]
**SubnetCidr** | **int32** | IPv4 subnet CIDR prefix (number of bits). | [optional] [default to null]
**SubnetCidrIpv6** | **int32** | IPv6 subnet CIDR prefix (number of bits). | [optional] [default to null]
**GwIp** | **string** | Gateway IPv4 address. | [optional] [default to null]
**GwIpv6** | **string** | Gateway IPv6 address. | [optional] [default to null]
**Vlan** | **int** | The VLAN of the virtual IP pool. | [optional] [default to null]
**State** | **string** |  | [optional] [default to null]
**CnodeIds** | **[]int32** | IDs of CNodes comprising the CNode group. | [optional] [default to null]
**Cluster** | **string** | Parent cluster. | [optional] [default to null]
**Url** | **string** |  | [optional] [default to null]
**DomainName** | **string** |  | [optional] [default to null]
**Role** | **string** | Role. | [optional] [default to null]
**IpRanges** | [**[][]string**](array.md) | IP ranges. | [optional] [default to null]
**VmsPreferred** | **bool** | If &#x27;true&#x27;, the CNodes included in this virtual IP pool are handled as preferred CNodes during VMS host election. | [optional] [default to null]
**Enabled** | **bool** | Enables or disables the virtual IP pool. | [optional] [default to null]
**PortMembership** | **string** | The port(s) on the CNode that this pool will use: &#x27;Right&#x27;, &#x27;Left&#x27; or &#x27;All&#x27;. | [optional] [default to null]
**ActiveInterfaces** | **int32** | The number of active interfaces. | [optional] [default to null]
**EnableL3** | **bool** | Enables or disables L3 CNode access. | [optional] [default to null]
**VastAsn** | **int32** | VAST ASN. | [optional] [default to null]
**PeerAsn** | **int32** | Peer ASN. | [optional] [default to null]
**TenantId** | **int64** | The ID of the tenant associated with the virtual IP pool. An ID of &#x27;0&#x27; (zero) means the virtual IP pool is available for all tenants. | [optional] [default to null]
**ActiveCnodeIds** | **[]int32** | IDs of active CNodes | [optional] [default to null]
**ClusterId** | **int32** | Cluster ID | [optional] [default to null]
**Cnodes** | **[]string** |  | [optional] [default to null]
**EnableWeightedBalancing** | **bool** | Weighted Balancing Enabled | [optional] [default to null]
**RangesSummary** | **string** | IP ranges | [optional] [default to null]
**SyncTime** | **string** | Synchronization time with leader | [optional] [default to null]
**Sync** | **string** | Synchronization state with leader | [optional] [default to null]
**TenantName** | **string** | Tenant Name | [optional] [default to null]
**Title** | **string** | IP range of the VIP pool | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

