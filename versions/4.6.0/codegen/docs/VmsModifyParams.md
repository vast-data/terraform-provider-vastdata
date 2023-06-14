# VmsModifyParams

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**DisableVmsMetrics** | **bool** | Set to True to disable VMS metrics. Set to False to enable VMS metrics. | [optional] [default to null]
**MgmtVip** | **string** | VMS Virtual IP | [optional] [default to null]
**MgmtInnerVip** | **string** | The virtual IP on the internal network used for mounting the VMS database. | [optional] [default to null]
**CapacityBase10** | **bool** | Format capacity properties in base 10 units | [optional] [default to null]
**PerformanceBase10** | **bool** | Format performance properties in base 10 units | [optional] [default to null]
**MinTlsVersion** | **string** | Minimum supported TLS version (e.g.: 1.2) | [optional] [default to null]
**AccessTokenLifetime** | **string** | Validity duration for JWT access token, specify as [DD [HH:[MM:]]]ss | [optional] [default to null]
**RefreshTokenLifetime** | **string** | Validity duration for JWT refresh token, specify as [DD [HH:[MM:]]]ss | [optional] [default to null]
**LoginBanner** | **string** | Custom login banner text for VMS Web UI and CLI | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


