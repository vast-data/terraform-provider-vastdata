# UserQuota

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **int32** |  | [optional] [default to null]
**Guid** | **string** | Quota guid | [optional] [default to null]
**Name** | **string** | The name | [optional] [default to null]
**State** | **string** |  | [optional] [default to null]
**GracePeriod** | **string** | Quota enforcement grace period in seconds, minutes, hours or days. Example: 90m | [optional] [default to null]
**TimeToBlock** | **string** | Grace period expiration time | [optional] [default to null]
**SoftLimit** | **int32** | Soft quota limit | [optional] [default to null]
**HardLimit** | **int32** | Hard quota limit | [optional] [default to null]
**HardLimitInodes** | **int32** | Hard inodes quota limit | [optional] [default to null]
**SoftLimitInodes** | **int32** | Soft inodes quota limit | [optional] [default to null]
**UsedInodes** | **int32** | Used inodes | [optional] [default to null]
**UsedCapacity** | **int64** | Used capacity in bytes | [optional] [default to null]
**IsAccountable** | **bool** |  | [optional] [default to null]
**QuotaSystemId** | **int32** |  | [optional] [default to null]
**Entity** | [***interface{}**](interface{}.md) |  | [optional] [default to null]
**EntityIdentifier** | **string** |  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


