# ReplicationTargetCreateParams

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** |  | [optional] [default to null]
**Proxies** | [**[]ErrorUnknown**](.md) | If configured, replication traffic is routed via proxies. Separate with commas. Format: http://&lt;username&gt;:&lt;password&gt;@&lt;IP&gt;:&lt;port&gt; | [optional] [default to null]
**AccessKey** | **string** | Access key of a valid key pair for accessing the named S3 bucket | [optional] [default to null]
**SecretKey** | **string** | The secret key of a valid key pair for accessing the destination S3 bucket | [optional] [default to null]
**BucketName** | **string** | The S3 bucket name of an existing S3 bucket that you want to configure as the replication target | [optional] [default to null]
**HttpProtocol** | **string** | For custom S3 buckets (not AWS), specifies which protocol to use to connect to the bucket | [optional] [default to null]
**CustomBucketUrl** | **string** | If the target is a custom S3 bucket, use this parameter to specify the URL of the bucket | [optional] [default to null]
**AwsRegion** | **string** | If the target is an AWS S3 bucket, use this parameter to specify the AWS region of the bucket | [optional] [default to null]
**AwsAccountId** | **string** | Not yet implemented | [optional] [default to null]
**AwsRole** | **string** | Not yet implemented | [optional] [default to null]
**Type_** | **string** | Specify AWS_S3 for an AWS S3 bucket. Specify CUSTOM_S3 for a custom S3 bucket. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


