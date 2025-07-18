# Quota

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **int32** |  | [optional] [default to null]
**Guid** | **string** | Quota GUID. | [optional] [default to null]
**Name** | **string** | Quota name. | [optional] [default to null]
**State** | **string** |  | [optional] [default to null]
**PrettyState** | **string** |  | [optional] [default to null]
**Path** | **string** | Directory path. | [optional] [default to null]
**PrettyGracePeriod** | **string** | The quota enforcement pretty grace period in seconds, minutes, hours or days. Example: 90m | [optional] [default to null]
**GracePeriod** | **string** | QQuota enforcement grace period in seconds, minutes, hours or days. Example: 90m | [optional] [default to null]
**TimeToBlock** | **string** | Grace period expiration time. | [optional] [default to null]
**SoftLimit** | **int64** | Soft quota limit. | [optional] [default to null]
**HardLimit** | **int64** | Hard quota limit. | [optional] [default to null]
**HardLimitInodes** | **int32** | Hard inodes quota limit. | [optional] [default to null]
**SoftLimitInodes** | **int32** | Soft inodes quota limit. | [optional] [default to null]
**UsedInodes** | **int32** | Used inodes. | [optional] [default to null]
**UsedCapacity** | **int64** | Used capacity in bytes. | [optional] [default to null]
**UsedCapacityTb** | **float32** | Used capacity in TB. | [optional] [default to null]
**UsedEffectiveCapacity** | **int64** | Used effective capacity in bytes. | [optional] [default to null]
**UsedEffectiveCapacityTb** | **float32** | Used effective capacity in TB. | [optional] [default to null]
**TenantId** | **int32** | Tenant ID. | [optional] [default to null]
**TenantName** | **string** | Tenant name. | [optional] [default to null]
**Cluster** | **string** | Parent cluster. | [optional] [default to null]
**ClusterId** | **int32** | Parent cluster ID. | [optional] [default to null]
**SystemId** | **int32** |  | [optional] [default to null]
**IsUserQuota** | **bool** |  | [optional] [default to null]
**EnableEmailProviders** | **bool** |  | [optional] [default to null]
**NumExceededUsers** | **int32** |  | [optional] [default to null]
**NumBlockedUsers** | **int32** |  | [optional] [default to null]
**EnableAlarms** | **bool** | Enables or disables alarms when users or groups are exceeding their limits. | [optional] [default to null]
**DefaultEmail** | **string** | The default email address used if there is no suffix specified and no address can be found in the providers. | [optional] [default to null]
**PercentInodes** | **int32** | Percentage of used inodes out of the hard limit. | [optional] [default to null]
**PercentCapacity** | **int32** | Percentage of used capacity out of the hard limit. | [optional] [default to null]
**DefaultUserQuota** | [***DefaultQuota**](DefaultQuota.md) |  | [optional] [default to null]
**DefaultGroupQuota** | [***DefaultQuota**](DefaultQuota.md) |  | [optional] [default to null]
**UserQuotas** | [**[]UserQuota**](UserQuota.md) |  | [optional] [default to null]
**GroupQuotas** | [**[]UserQuota**](UserQuota.md) |  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

