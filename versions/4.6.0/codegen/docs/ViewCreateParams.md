# ViewCreateParams

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | A name for the view | [optional] [default to null]
**Path** | **string** | The Element Store path to exposed through the view | [default to null]
**Alias** | **string** | For NFS-enabled views, an alias that can be used by NFSv3 clients | [optional] [default to null]
**Bucket** | **string** |  | [optional] [default to null]
**PolicyId** | **int32** | Specify (by ID) which view policy should be applied to the view | [optional] [default to null]
**ClusterId** | **int32** | Cluster ID | [optional] [default to null]
**CreateDir** | **bool** | Create a directory at the specified path | [optional] [default to null]
**Protocols** | **[]string** | Enabled client access protocols | [optional] [default to null]
**Share** | **string** | SMB share name | [optional] [default to null]
**BucketOwner** | **string** | S3 Bucket owner | [optional] [default to null]
**S3Locks** | **bool** | S3 Object Lock | [optional] [default to null]
**S3LocksRetentionMode** | **string** | S3 Locks retention mode | [optional] [default to null]
**S3LocksRetentionPeriod** | **string** | Period should be positive in format like 0d|2d|1y|2y | [optional] [default to null]
**BucketCreators** | **[]string** |  | [optional] [default to null]
**BucketCreatorsGroups** | **[]string** |  | [optional] [default to null]
**S3Versioning** | **bool** | S3 Versioning | [optional] [default to null]
**S3UnverifiedLookup** | **bool** | S3 Unverified Lookup | [optional] [default to null]
**AllowAnonymousAccess** | **bool** | Allow S3 anonymous access | [optional] [default to null]
**AllowS3AnonymousAccess** | **bool** | Allow S3 anonymous access | [optional] [default to null]
**NfsInteropFlags** | **string** | Indicates whether the view should support simultaneous access to NFS3/NFS4/SMB protocols. | [optional] [default to null]
**ShareAcl** | [***interface{}**](interface{}.md) | Share-level ACL details | [optional] [default to null]
**SelectForLiveMonitoring** | **bool** |  | [optional] [default to null]
**QosPolicyId** | **int32** | QoS Policy ID | [optional] [default to null]
**QosPolicy** | **string** | QoS Policy | [optional] [default to null]
**TenantId** | **int32** | Tenant ID | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


