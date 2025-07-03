# UserQuota

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**GracePeriod** | **string** | Quota enforcement grace period in the format of HH:MM:SS. | [optional] [default to null]
**TimeToBlock** | **string** | Grace period expiration time. | [optional] [default to null]
**SoftLimit** | **int64** | Soft quota limit. | [optional] [default to null]
**HardLimit** | **int64** | Hard quota limit. | [optional] [default to null]
**HardLimitInodes** | **int64** | Hard quota limit for the number of inodes. | [optional] [default to null]
**SoftLimitInodes** | **int64** | Soft quota limit for the number of inodes. | [optional] [default to null]
**UsedInodes** | **int64** | Used inodes. | [optional] [default to null]
**UsedCapacity** | **int64** | Used capacity in bytes. | [optional] [default to null]
**IsAccountable** | **bool** |  | [optional] [default to null]
**QuotaSystemId** | **int32** |  | [optional] [default to null]
**Entity** | [***QuotaEntityInfo**](QuotaEntityInfo.md) |  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

