# DnsModifyParams

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** |  | [optional] [default to null]
**Vip** | **string** | A virtual IP to assign to the DNS service. DNS requests from your external DNS server must be delegated to this IP. | [optional] [default to null]
**DomainSuffix** | **string** | A suffix for domain names. Requests for domain names with this suffix are resolved to the VIPs configured on the cluster. | [optional] [default to null]
**VipGateway** | **string** | If the external DNS server doesn&#39;t reside on the same subnet as the DNS VIP, enter the IP of a gateway through which to connect to the DNS server.  | [optional] [default to null]
**Enabled** | **bool** | Set to true to enable the DNS service | [optional] [default to null]
**NetType** | **string** |  | [optional] [default to null]
**InvalidNameResponse** | **string** |  | [optional] [default to null]
**InvalidTypeResponse** | **string** |  | [optional] [default to null]
**VipSubnetCidr** | **int32** | The subnet, in CIDR format, on which the DNS VIP resides. | [optional] [default to null]
**CnodeIds** | **string** | To dedicate a specific group of CNodes to the DNS, list the IDs of the CNodes. | [optional] [default to null]
**VipVlan** | **int32** | If your external DNS server is only exposed to a specific VLAN, you can enter the VLAN here to enable communication with the DNS server. | [optional] [default to null]
**VipIpv6** | **string** | Assigns an IPv6 to the DNS service. | [optional] [default to null]
**VipIpv6SubnetCidr** | **int32** | Specifies the subnet, as a CIDR index, on which the DNS resides. [1..128] | [optional] [default to null]
**VipIpv6Gateway** | **string** | Specifies a gateway IPv6 to external DNS server if on different subnet. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


