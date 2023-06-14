# UserQuotaCreateParams

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | A name for the quota | [optional] [default to null]
**Identifier** | **string** |  | [optional] [default to null]
**IdentifierType** | **string** |  | [optional] [default to null]
**IsGroup** | **bool** |  | [optional] [default to null]
**QuotaId** | **int32** |  | [optional] [default to null]
**GracePeriod** | **string** | Quota enforcement grace period. An alarm is triggered and write operations are blocked if storage usage continues to exceed the soft limit for the grace period. Format: [DD] [HH:[MM:]]ss | [optional] [default to null]
**SoftLimit** | **int32** | Storage usage limit at which warnings of exceeding the quota are issued. | [optional] [default to null]
**HardLimit** | **int32** | Storage usage limit beyond which no writes will be allowed. | [optional] [default to null]
**HardLimitInodes** | **int32** | Number of directories and unique files under the path beyond which no writes will be allowed. A file with multiple hardlinks is counted only once. | [optional] [default to null]
**SoftLimitInodes** | **int32** | Number of directories and unique files under the path at which warnings of exceeding the quota will be issued. A file with multiple hardlinks is counted only once. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


