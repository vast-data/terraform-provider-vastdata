# QosPolicy

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **int32** |  | [optional] [default to null]
**Guid** | **string** | QoS Policy guid | [optional] [default to null]
**Name** | **string** |  | [optional] [default to null]
**Mode** | **string** | QoS provisioning mode | [optional] [default to null]
**IoSizeBytes** | **int64** | Sets the size of IO for static and capacity limit definitions. The number of IOs per request is obtained by dividing request size by IO size. Default: 64K, Recommended range: 4K - 1M | [optional] [default to null]
**StaticLimits** | [***QosStaticLimits**](QosStaticLimits.md) |  | [optional] [default to null]
**CapacityLimits** | [***QosDynamicLimits**](QosDynamicLimits.md) |  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

