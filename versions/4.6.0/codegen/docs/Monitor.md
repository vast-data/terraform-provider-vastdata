# Monitor

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **int32** |  | [optional] [default to null]
**Name** | **string** | Monitor name | [optional] [default to null]
**ObjectType** | **string** | Object type | [optional] [default to null]
**MonitorType** | **string** | Monitor type (predefined, custom, etc) | [optional] [default to null]
**FromTime** | [**time.Time**](time.Time.md) |  | [optional] [default to null]
**ToTime** | [**time.Time**](time.Time.md) |  | [optional] [default to null]
**TimeFrame** | **string** | Time frame to use | [optional] [default to null]
**ObjectIds** | [***interface{}**](interface{}.md) | Only query metrics on these objects (optional) | [optional] [default to null]
**PropList** | [***interface{}**](interface{}.md) | Only query these metrics (optional) | [optional] [default to null]
**Granularity** | **string** | Data granularity | [optional] [default to null]
**Aggregation** | **string** | Aggregation function. avg, min, max etc. | [optional] [default to null]
**QueryAggregation** | **string** | Special aggregations to apply on query, e.g intersampling | [optional] [default to null]
**MetricsExposure** | **string** | Monitor&#39;s metrics exposure | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


