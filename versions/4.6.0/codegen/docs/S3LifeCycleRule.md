# S3LifeCycleRule

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **int32** |  | [optional] [default to null]
**Name** | **string** | A unique name | [optional] [default to null]
**Guid** | **string** |  | [optional] [default to null]
**Enabled** | **bool** |  | [optional] [default to null]
**Prefix** | **string** | The prefix to filter object names | [optional] [default to null]
**MinSize** | **int32** | The minimum size of the object | [optional] [default to null]
**MaxSize** | **int32** | The maximum size of the object | [optional] [default to null]
**ExpirationDays** | **int32** | The number of days from creation until an object expires | [optional] [default to null]
**ExpirationDate** | **string** | The expiration date of the object | [optional] [default to null]
**ExpiredObjDeleteMarker** | **bool** | Remove expired objects delete markers | [optional] [default to null]
**NoncurrentDays** | **int32** | Number of days after objects become noncurrent | [optional] [default to null]
**NewerNoncurrentVersions** | **int32** | The number of newer versions to retain | [optional] [default to null]
**AbortMpuDaysAfterInitiation** | **int32** | The number of days until expiration after an incomplete multipart upload | [optional] [default to null]
**ViewPath** | **string** | The path of the related View | [optional] [default to null]
**ViewId** | **int32** | The ID of the related View | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


