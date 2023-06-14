# MonitorAddParams

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | A name for the monitor. | [optional] [default to null]
**ObjectType** | **string** | The type of object to monitor. | [default to null]
**FromTime** | **string** |  | [optional] [default to null]
**ToTime** | **string** |  | [optional] [default to null]
**TimeFrame** | **string** | Default time frame to report over. Examples: 2h (2 hours), 1D (1 Day), 10m (10 minutes), 1M (1 month) | [optional] [default to null]
**ObjectIds** | **string** | Specific objects to include in the report, specified as a comma separated list of object IDs. | [optional] [default to null]
**PropList** | **string** | A list of metrics to query. To get the full list of metrics, use GET /metrics/. | [optional] [default to null]
**Granularity** | **string** | Data granularity: seconds (raw), minutes (five minute aggregated samples), hours (hourly aggregated samples), or days (daily aggregated samples) | [optional] [default to null]
**Aggregation** | **string** | If data granularity is minutes, hours or days, the data is aggregated. This parameter selects which aggregation function to use. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


