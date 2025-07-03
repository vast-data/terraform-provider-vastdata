# Dns

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | The name of the VAST DNS server configuration. | [optional] [default to null]
**Id** | **int32** | The unique ID of the VAST DNS server configuration. | [optional] [default to null]
**Vip** | **string** | The virtual IP for the DNS service. DNS requests from your external DNS server must be delegated to this IP. | [optional] [default to null]
**DomainSuffix** | **string** | A suffix to append to domain names of each virtual IP pool. The suffix should append each domain name to form a valid FQDN for DNS requests to target. | [optional] [default to null]
**VipGateway** | **string** | The IPv4 address of the gateway to the external DNS server if it is on a different subnet. Must be on the same subnet as the IP and reachable from the relevant network interface. | [optional] [default to null]
**Enabled** | **bool** | Enables or disables the VAST DNS server configuration. | [optional] [default to null]
**Guid** | **string** | The unique GUID of the VAST DNS server configuration. | [optional] [default to null]
**VipSubnetCidr** | **int32** | The subnet, in the CIDR format, on which the DNS resides. | [optional] [default to null]
**VipVlan** | **int32** | The VLAN (optional) to enable communication with the external DNS server(s). | [optional] [default to null]
**CnodeIds** | **[]int32** |  | [optional] [default to null]
**VipIpv6** | **string** | The IPv6 address of the DNS service. | [optional] [default to null]
**VipIpv6SubnetCidr** | **int32** | The subnet, in the CIDR format, on which the DNS resides. Valid values: [1..128] | [optional] [default to null]
**VipIpv6Gateway** | **string** | The IPv6 address of the gateway to the external DNS server if it is on a different subnet. | [optional] [default to null]
**NetType** | **string** | The interface that listens for DNS service delegation requests. | [optional] [default to NET_TYPE.EXTERNAL_PORT]
**InvalidNameResponse** | **string** | The response DNS type for invalid DNS name. | [optional] [default to null]
**InvalidTypeResponse** | **string** | The response DNS type for invalid DNS type. | [optional] [default to null]
**Ttl** | **int32** | The response TTL in seconds. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

