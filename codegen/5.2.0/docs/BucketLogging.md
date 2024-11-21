# BucketLogging

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**DestinationId** | **int64** | The Logging bucket ID | [optional] [default to null]
**Prefix** | **string** | Log line prefix to add | [optional] [default to null]
**KeyFormat** | **string** | The format for log object keys. SIMPLE_PREFIX&#x3D;[DestinationPrefix][YYYY]-[MM]-[DD]-[hh]-[mm]-[ss]-[UniqueString], PARTITIONED_PREFIX_EVENT_TIME&#x3D;[DestinationPrefix][SourceUsername]/[SourceBucket]/[YYYY]/[MM]/[DD]/[YYYY]-[MM]-[DD]-[hh]-[mm]-[ss]-[UniqueString] where the partitioning is done based on the time when the logged events occurred, PARTITIONED_PREFIX_DELIVERY_TIME&#x3D;[DestinationPrefix][SourceUsername]/[SourceBucket]/[YYYY]/[MM]/[DD]/[YYYY]-[MM]-[DD]-[hh]-[mm]-[ss]-[UniqueString] where the partitioning is done based on the time when the log object has been delivered to the destination bucket. Default: SIMPLE_PREFIX | [optional] [default to KEY_FORMAT.SIMPLE_PREFIX]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

