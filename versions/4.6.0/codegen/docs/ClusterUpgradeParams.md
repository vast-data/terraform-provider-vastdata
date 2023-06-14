# ClusterUpgradeParams

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Build** | **string** | Specifies the build for upgrade | [default to null]
**OsUpgrade** | **bool** | OS upgrade | [optional] [default to null]
**BmcUpgrade** | **bool** | BMC upgrade | [optional] [default to null]
**EnableDr** | **bool** | Enables data reduction (DR) for a cluster without DR enabled prior to upgrade | [optional] [default to null]
**Force** | **bool** | Forces upgrade regardless of version or upgrade state | [optional] [default to null]
**SkipHwCheck** | **bool** | Skips validation of hardware component health. Use with caution since component redundancy is important in NDU. Do not use with OS upgrade. | [optional] [default to null]
**Prepare** | **bool** | Pull docker images only | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


