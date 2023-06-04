# S3replicationPeers

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **int64** | A unique id given to the s3 replication peer configuration | [optional] [default to null]
**Guid** | **string** | A unique guid given to the s3 replication peer configuration | [optional] [default to null]
**Name** | **string** | The name of the s3 replication peer configuration | [optional] [default to null]
**Url** | **string** | Direct link to the s3 replication peer configurations | [optional] [default to null]
**BucketName** | **string** | The name of the peer bucket to replicate to | [optional] [default to null]
**HttpProtocol** | **string** | The http protocol user http/https | [optional] [default to null]
**Type_** | **string** | The type of the peer bucket CUSTOM_S3/AWS_S3 | [default to null]
**Proxies** | **[]string** | List of http procies | [optional] [default to null]
**AwsRegion** | **string** | The Bucket AWS region, Valid only when type is AWS_S3 | [optional] [default to null]
**AccessKey** | **string** | The S3 access key | [optional] [default to null]
**SecretKey** | **string** | The S3 secret key | [optional] [default to null]
**CustomBucketUrl** | **string** | The S3 url of the bucket (dns name/ip) used only when using CUSTOM_S3 | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

