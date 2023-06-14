# SystemShardExpandParams

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**EstoreShardCount** | **int32** | EStore shard count | [optional] [default to null]
**DrShardCount** | **int32** | DR shard count | [optional] [default to null]
**DrWbShardCount** | **int32** | DR WB shard count | [optional] [default to null]
**Force** | **bool** | Force shard expansion. Use if you want to run shard expansion even though shards are denylisted (indicated by MAINTENANCE_DENYLIST_EXISTS in error code when running without &#39;force&#39;). | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


