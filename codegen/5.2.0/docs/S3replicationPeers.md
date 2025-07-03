# S3replicationPeers

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **int64** | The unique ID of the S3 replication peer configuration. | [optional] [default to null]
**Guid** | **string** | The unique GUID of the S3 replication peer configuration. | [optional] [default to null]
**Name** | **string** | The name of the S3 replication peer configuration. | [optional] [default to null]
**Url** | **string** | Direct URL of the S3 replication peer configuration. | [optional] [default to null]
**BucketName** | **string** | The name of the peer bucket to replicate to. | [optional] [default to null]
**HttpProtocol** | **string** | The HTTP protocol used (HTTP or HTTPS). | [optional] [default to null]
**Type_** | **string** | The type of the peer bucket: CUSTOM_S3/AWS_S3 | [default to null]
**Proxies** | **[]string** | A list of HTTP proxies. | [optional] [default to null]
**AwsRegion** | **string** | The AWS region of the bucket. Valid only when type is AWS_S3. | [optional] [default to null]
**AccessKey** | **string** | The S3 access key. | [optional] [default to null]
**SecretKey** | **string** | The S3 secret key. | [optional] [default to null]
**CustomBucketUrl** | **string** | The S3 URL of the bucket (DNS name/IP), used only when using CUSTOM_S3. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

