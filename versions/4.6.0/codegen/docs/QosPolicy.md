# QosPolicy

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **int32** |  | [default to null]
**Guid** | **string** | QoS Policy guid | [default to null]
**Name** | **string** |  | [default to null]
**Mode** | **string** | Allocation of performance resources mode | [default to null]
**IoSizeBytes** | **int32** | Size of a single IO, default is 64K | [optional] [default to null]
**StaticLimits** | [***QosStaticLimits**](QOSStaticLimits.md) | Static mode limits | [default to null]
**CapacityLimits** | [***QosDynamicLimits**](QOSDynamicLimits.md) | Capacity mode limits | [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


