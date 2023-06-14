# AddBoxesParams

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Cboxes** | [***AddBoxesParamsCboxes**](AddBoxesParams_cboxes.md) |  | [default to null]
**Dboxes** | [***AddBoxesParamsDboxes**](AddBoxesParams_dboxes.md) |  | [default to null]
**DnodeManagementIpPool** | **[]string** |  | [default to null]
**CnodeManagementIpPool** | **[]string** |  | [default to null]
**DnodeIpmiPool** | **[]string** |  | [default to null]
**CnodeIpmiPool** | **[]string** |  | [default to null]
**ManagementCidr** | **int32** |  | [default to null]
**Ipv6Prefix** | **int32** |  | [optional] [default to null]
**ExternalGateway** | **[]string** |  | [default to null]
**EmptyDbox** | **bool** |  | [default to null]
**CnodeStartIndex** | **int32** | CNode start index valid range: [1, 99], for null value index will be selected automatically, as max_existed_index + 1 | [optional] [default to null]
**DnodeStartIndex** | **int32** | CNode start index valid range: [100, 254], for null value index will be selected automatically, as max_existed_index + 1 | [optional] [default to null]
**HostnamePrefix** | **string** |  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


