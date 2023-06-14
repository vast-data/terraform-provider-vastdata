# Vms

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **int32** |  | [optional] [default to null]
**Ip** | **string** | Data Bond IP | [optional] [default to null]
**Ip1** | **string** | Data IP1 | [optional] [default to null]
**Ip2** | **string** | Data IP2 | [optional] [default to null]
**MgmtIp** | **string** | Management IP | [optional] [default to null]
**MgmtVip** | **string** | Management VIP | [optional] [default to null]
**MgmtVipIpv6** | **string** | Management IPv6 VIP | [optional] [default to null]
**MgmtCnode** | **string** | Management CNode | [optional] [default to null]
**MgmtInnerVip** | **string** | VMS inner (NFS) vip | [optional] [default to null]
**MgmtInnerVipCnode** | **string** | VMS inner VIP CNode | [optional] [default to null]
**Name** | **string** |  | [default to null]
**Build** | **string** |  | [optional] [default to null]
**SwVersion** | **string** |  | [optional] [default to null]
**AutoLogoutTimeout** | **string** |  | [optional] [default to null]
**Created** | [**time.Time**](time.Time.md) |  | [optional] [default to null]
**State** | **string** |  | [optional] [default to null]
**Url** | **string** |  | [optional] [default to null]
**DisableMgmtHa** | **bool** | Is Management HA disabled? | [optional] [default to null]
**DisableVmsMetrics** | **bool** | Disable vms metrics collection | [optional] [default to null]
**CapacityBase10** | **bool** | Format capacity properties to base 10 | [optional] [default to null]
**PerformanceBase10** | **bool** | Format performance properties to base 10 | [optional] [default to null]
**MinTlsVersion** | **string** | Minimum supported TLS version (e.g.: 1.2) | [optional] [default to null]
**DegradedReason** | **string** | The reason for VMS degraded state | [optional] [default to null]
**AccessTokenLifetime** | **string** | Validity duration for JWT access token, specify as [DD [HH:[MM:]]]ss | [optional] [default to null]
**RefreshTokenLifetime** | **string** | Validity duration for JWT refresh token, specify as [DD [HH:[MM:]]]ss | [optional] [default to null]
**LoginBanner** | **string** | Customize login banner for VMS Web UI and CLI | [optional] [default to null]
**TotalUsageCapacityPercentage** | **string** | the remaining capacity out of the active capacity | [optional] [default to null]
**TotalActiveCapacity** | **string** | sum of all active licenses capacity | [optional] [default to null]
**TotalRemainingCapacity** | **string** | the total active capacity minus the system capacity | [optional] [default to null]
**TabularSupport** | **string** | Parameter that controls everything related to database | [optional] [default to null]
**Ipv6Support** | **bool** | Parameter that controls visibility of ipv6 fields for VIP Pools | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


