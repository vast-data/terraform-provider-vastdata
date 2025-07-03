# QosPolicy

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **int32** |  | [optional] [default to null]
**Guid** | **string** | QoS policy GUID. | [optional] [default to null]
**Name** | **string** |  | [optional] [default to null]
**Mode** | **string** | QoS provisioning mode. | [optional] [default to null]
**PolicyType** | **string** | The QoS policy type. | [optional] [default to null]
**LimitBy** | **string** | Specifies which attributes are setting the limitations. | [optional] [default to LIMIT_BY.BW_IOPS]
**TenantId** | **int64** | When setting &#x27;is_default&#x27;, this is the tenant for which the policy will be used as the default user QoS policy. | [optional] [default to null]
**AttachedUsersIdentifiers** | **[]string** | A list of local user IDs to which this QoS policy applies. | [optional] [default to null]
**IsDefault** | **bool** | Specifies whether this QoS policy is to be used as the default QoS policy per user for this tenant. Setting this attribute requires that &#x27;tenant_id&#x27; is also supplied. | [optional] [default to null]
**IoSizeBytes** | **int64** | Sets the size of IO for static and capacity limit definitions. The number of IOs per request is obtained by dividing the request size by IO size. Default: 64K. Recommended range: 4K - 1M. | [optional] [default to null]
**StaticLimits** | [***QosStaticLimits**](QosStaticLimits.md) |  | [optional] [default to null]
**CapacityLimits** | [***QosDynamicLimits**](QosDynamicLimits.md) |  | [optional] [default to null]
**StaticTotalLimits** | [***QoSStaticTotalLimits**](QoSStaticTotalLimits.md) |  | [optional] [default to null]
**CapacityTotalLimits** | [***QoSDynamicTotalLimits**](QoSDynamicTotalLimits.md) |  | [optional] [default to null]
**AttachedUsers** | [**[]QosUser**](QosUser.md) |  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

