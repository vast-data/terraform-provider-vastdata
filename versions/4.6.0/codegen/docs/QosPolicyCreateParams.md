# QosPolicyCreateParams

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** |  | [default to null]
**Mode** | **string** | Allocation of performance resources mode | [optional] [default to null]
**IoSizeBytes** | **int32** | Size of a single IO, default is 64K | [optional] [default to null]
**StaticLimits** | [***RequestQosStaticLimits**](RequestQOSStaticLimits.md) | Static mode limits | [optional] [default to null]
**CapacityLimits** | [***RequestQosDynamicLimits**](RequestQOSDynamicLimits.md) | Capacity mode limits | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


