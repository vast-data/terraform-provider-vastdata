# ClusterUpgrade

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Build** | **string** | Specifies build for upgrade | [default to null]
**SkipHwCheck** | **bool** | Skips validation of hardware component health. Use with caution since component redundancy is important in NDU. Do not use with OS upgrade. | [optional] [default to null]
**EnableDr** | **bool** | Change system settings from nondr to dr | [optional] [default to null]
**Prepare** | **bool** | Pull docker images only | [optional] [default to null]
**Force** | **bool** | Force upgrade | [optional] [default to null]
**Os** | **bool** | Execute OS upgrade if possible | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


