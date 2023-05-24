# AsyncSupportBundle

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **int32** |  | [optional] [default to null]
**Guid** | **string** |  | [optional] [default to null]
**Name** | **string** |  | [optional] [default to null]
**State** | **string** |  | [optional] [default to null]
**Prefix** | **string** | The name of the support bundle | [optional] [default to null]
**Path** | **string** | The path | [optional] [default to null]
**BundleFile** | **string** |  | [optional] [default to null]
**BundleSize** | **string** |  | [optional] [default to null]
**BundleUrl** | **string** | URL to access/download support bundle file | [optional] [default to null]
**Level** | **string** |  | [optional] [default to null]
**StartTime** | **string** | Start time of logs | [optional] [default to null]
**EndTime** | **string** | End time of logs | [optional] [default to null]
**Text** | **bool** | Include only textual logs in bundle | [optional] [default to null]
**CnodeIds** | **string** | Comma separated IDs to fetch (Fetch all if not specified) | [optional] [default to null]
**DnodeIds** | **string** | Comma separated IDs to fetch (Fetch all if not specified) | [optional] [default to null]
**Created** | [**time.Time**](time.Time.md) |  | [optional] [default to null]
**CreateDatetime** | **string** |  | [optional] [default to null]
**Aggregated** | **bool** | Aggregate harvests into final bundle | [optional] [default to null]
**Preset** | **string** | A predefined preset bundle | [optional] [default to null]
**Obfuscated** | **bool** | Obfuscate text files | [optional] [default to null]
**Cluster** | **string** | Parent Cluster | [optional] [default to null]
**Url** | **string** |  | [optional] [default to null]
**MaxSize** | **float32** | Trace files bundle size limit in GB for each node | [optional] [default to null]
**DeleteAfterSend** | **bool** | Delete the bundle immediately after successfully uploading | [optional] [default to null]
**AsyncTask** | [***interface{}**](interface{}.md) | Creation Async task properties | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


