# S3LifeCycleRule

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **int32** |  | [optional] [default to null]
**Name** | **string** | A unique name. | [optional] [default to null]
**Guid** | **string** |  | [optional] [default to null]
**Enabled** | **bool** |  | [optional] [default to null]
**Prefix** | **string** | Defines a scope of elements (objects, files or directories) by prefix. All objects with keys that begin with the specified prefix are included in the scope. In file and directory nomenclature, a prefix is a file and/or directory path within the view that can include part of the file or directory name. For example, sales/jan would include the file sales/january and the directory sales/jan/week1/. No characters are handled as wildcards. | [optional] [default to null]
**MinSize** | **int64** | The minimum size of the object. | [optional] [default to null]
**MaxSize** | **int64** | The maximum size of the object. | [optional] [default to null]
**ExpirationDays** | **int32** | The number of days from creation until an object expires. | [optional] [default to null]
**ExpirationDate** | **string** | The expiration date of the object. | [optional] [default to null]
**ExpiredObjDeleteMarker** | **bool** | If &#x27;true&#x27;, removes expired object delete markers. | [optional] [default to null]
**NoncurrentDays** | **int32** | Number of days after which objects become non-current | [optional] [default to null]
**NewerNoncurrentVersions** | **int32** | The number of newer versions to retain. | [optional] [default to null]
**AbortMpuDaysAfterInitiation** | **int32** | The number of days until expiration after an incomplete multipart upload. | [optional] [default to null]
**ViewPath** | **string** | The path of the view to which the lifecycle rule applies. | [optional] [default to null]
**ViewId** | **int32** | The ID of the view to which the lifecycle rule applies. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

