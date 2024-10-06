# Dns

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | Specifies a name for the VAST DNS server configuration | [optional] [default to null]
**Id** | **int32** | A uniqe id given to the VAST DNS server configurations | [optional] [default to null]
**Vip** | **string** | Assigns a IP to the DNS service. DNS requests from your external DNS server must be delegated to this IP. | [optional] [default to null]
**DomainSuffix** | **string** | Specifies a suffix to append to domain names of each VIP pool. The suffix should complete each domain name to form a valid FQDN for DNS requests to target. | [optional] [default to null]
**VipGateway** | **string** | Specifies a gateway IP to external DNS server if on different subnet. Must be on same subnet as the IP and reachable from the relevant nework interface. | [optional] [default to null]
**Enabled** | **bool** | Enable the VAST DNS server configurations | [optional] [default to null]
**Guid** | **string** | A uniqe guid assigned to the VAST DNS server configurations | [optional] [default to null]
**VipSubnetCidr** | **int32** | Specifies the subnet, as a CIDR index, on which the DNS resides. | [optional] [default to null]
**VipVlan** | **int32** | Specifies a VLAN if needed to enable communication with external DNS server(s). | [optional] [default to null]
**CnodeIds** | **[]int32** |  | [optional] [default to null]
**VipIpv6** | **string** | Assigns an IPv6 to the DNS service. | [optional] [default to null]
**VipIpv6SubnetCidr** | **int32** | Specifies the subnet, as a CIDR index, on which the DNS resides. [1..128] | [optional] [default to null]
**VipIpv6Gateway** | **string** | Specifies a gateway IPv6 to external DNS server if on different subnet. | [optional] [default to null]
**NetType** | **string** | Select the interface, that listens for DNS service delegation requests | [optional] [default to null]
**InvalidNameResponse** | **string** | The response DNS type for invalid dns name | [optional] [default to null]
**InvalidTypeResponse** | **string** | The response DNS type for invalid dns type | [optional] [default to null]
**Ttl** | **int32** | The reposne TTL in seconds | [optional] [default to null]
**DnsPort** | **int32** | The DNS listenning port | [optional] [default to 53]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

