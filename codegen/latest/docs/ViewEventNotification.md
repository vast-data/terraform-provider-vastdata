# ViewEventNotification

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | A unique name of an event among all other events related to the view holding this event | [optional] [default to null]
**Topic** | **string** | The name of the kafka topic to send alert to | [optional] [default to null]
**BrokerId** | **int** | The id of the external kafka broker | [optional] [default to null]
**Triggers** | **[]string** | List of S3 triggers which will trigger event notification, The following events are supported: - S3_OBJECT_CREATED_ALL - S3_OBJECT_CREATED_PUT - S3_OBJECT_CREATED_POST - S3_OBJECT_CREATED_COPY - S3_OBJECT_CREATED_COMPLETE_MULTIPART_UPLOAD - S3_OBJECT_REMOVED_ALL - S3_OBJECT_REMOVED_DELETE - S3_OBJECT_REMOVED_DELETE_MARKER_CREATED | [optional] [default to null]
**PrefixFilter** | **string** | Event prefix filter | [optional] [default to null]
**SuffixFilter** | **string** | Event suffix filter | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

