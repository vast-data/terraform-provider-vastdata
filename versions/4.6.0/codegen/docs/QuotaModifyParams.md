# QuotaModifyParams

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | Quota name | [optional] [default to null]
**GracePeriod** | **string** | Quota enforcement grace period. An alarm is triggered and write operations are blocked if storage usage continues to exceed the soft limit for the grace period. Format: [DD] [HH:[MM:]]s | [optional] [default to null]
**SoftLimit** | **string** | Storage usage limit at which warnings of exceeding the quota are issued. | [optional] [default to null]
**HardLimit** | **string** | Storage usage limit beyond which no writes will be allowed. | [optional] [default to null]
**HardLimitInodes** | **int32** | Number of directories and unique files under the path beyond which no writes will be allowed. A file with multiple hardlinks is counted only once. | [optional] [default to null]
**SoftLimitInodes** | **int32** | Number of directories and unique files under the path at which warnings of exceeding the quota will be issued. A file with multiple hardlinks is counted only once. | [optional] [default to null]
**DefaultUserQuota** | [***interface{}**](interface{}.md) |  | [optional] [default to null]
**DefaultGroupQuota** | [***interface{}**](interface{}.md) |  | [optional] [default to null]
**UserQuotas** | [**[]interface{}**](interface{}.md) |  | [optional] [default to null]
**GroupQuotas** | [**[]interface{}**](interface{}.md) |  | [optional] [default to null]
**EnableAlarms** | **bool** |  | [optional] [default to null]
**DefaultEmail** | **string** |  | [optional] [default to null]
**EnableEmailProviders** | **bool** |  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


