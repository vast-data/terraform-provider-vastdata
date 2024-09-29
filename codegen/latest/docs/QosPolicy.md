# QosPolicy

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **int32** |  | [optional] [default to null]
**Guid** | **string** | QoS Policy guid | [optional] [default to null]
**Name** | **string** |  | [optional] [default to null]
**Mode** | **string** | QoS provisioning mode | [optional] [default to null]
**PolicyType** | **string** | The QoS type | [optional] [default to null]
**LimitBy** | **string** | What attributes are setting the limitations. | [optional] [default to LIMIT_BY.BW_IOPS]
**TenantId** | **int64** | When setting is_default this is the tenant which will take affect | [optional] [default to null]
**AttachedUsersIdentifiers** | **[]string** | List of local user IDs to which this QoS Policy is affective. | [optional] [default to null]
**IsDefault** | **bool** | Should this QoS Policy be the default QoS per user for this tenant ?, tnenat_id should be also provided when settingthis attribute | [optional] [default to null]
**IoSizeBytes** | **int64** | Sets the size of IO for static and capacity limit definitions. The number of IOs per request is obtained by dividing request size by IO size. Default: 64K, Recommended range: 4K - 1M | [optional] [default to null]
**StaticLimits** | [***QosStaticLimits**](QosStaticLimits.md) |  | [optional] [default to null]
**CapacityLimits** | [***QosDynamicLimits**](QosDynamicLimits.md) |  | [optional] [default to null]
**StaticTotalLimits** | [***QoSStaticTotalLimits**](QoSStaticTotalLimits.md) |  | [optional] [default to null]
**CapacityTotalLimits** | [***QoSDynamicTotalLimits**](QoSDynamicTotalLimits.md) |  | [optional] [default to null]
**AttachedUsers** | [**[]QosUser**](QosUser.md) |  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

