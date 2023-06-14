# View

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **int32** |  | [optional] [default to null]
**Guid** | **string** |  | [optional] [default to null]
**Name** | **string** |  | [optional] [default to null]
**Created** | [**time.Time**](time.Time.md) |  | [optional] [default to null]
**Path** | **string** | File system path. Begin with &#39;/&#39;. Do not include a trailing slash | [default to null]
**CreateDir** | **string** | Creates the directory specified by the path | [optional] [default to null]
**Alias** | **string** | Alias for NFS export, must start with &#39;/&#39; and only ASCII characters are allowed. If configured, this supersedes the exposed NFS export path | [optional] [default to null]
**Bucket** | **string** | S3 Bucket name | [optional] [default to null]
**Policy** | **string** | The view policy that applies to this view | [optional] [default to null]
**PolicyId** | **int32** | Associated view policy ID | [optional] [default to null]
**Cluster** | **string** | Parent Cluster | [optional] [default to null]
**ClusterId** | **int32** | Parent Cluster ID | [optional] [default to null]
**TenantId** | **int32** | Tenant ID | [optional] [default to null]
**Url** | **string** |  | [optional] [default to null]
**Directory** | **bool** | Create the directory if it does not exist | [optional] [default to null]
**S3Versioning** | **bool** | S3 Versioning | [optional] [default to null]
**S3UnverifiedLookup** | **bool** | S3 Unverified Lookup | [optional] [default to null]
**AllowAnonymousAccess** | **bool** | Allow S3 anonymous access | [optional] [default to null]
**AllowS3AnonymousAccess** | **bool** | Allow S3 anonymous access | [optional] [default to null]
**Protocols** | **[]string** | Protocols exposed by this view | [optional] [default to null]
**Share** | **string** | Name of the SMB Share. Must not include the following characters: \&quot; \\ / [ ] : | &lt; &gt; + &#x3D; ; , * ? | [optional] [default to null]
**Sync** | **string** | Synchronization state with leader | [optional] [default to null]
**SyncTime** | **string** | Synchronization time with leader | [optional] [default to null]
**BucketOwner** | **string** | S3 Bucket owner | [optional] [default to null]
**BucketCreators** | **[]string** | List of bucket creators users | [optional] [default to null]
**BucketCreatorsGroups** | **[]string** | List of bucket creators groups | [optional] [default to null]
**S3Locks** | **bool** | S3 Object Lock | [optional] [default to null]
**S3LocksRetentionMode** | **string** | S3 Locks retention mode | [optional] [default to null]
**S3LocksRetentionPeriod** | **string** | Period should be positive in format like 0d|2d|1y|2y | [optional] [default to null]
**PhysicalCapacity** | **int32** | Physical Capacity | [optional] [default to null]
**LogicalCapacity** | **int32** | Logical Capacity | [optional] [default to null]
**NfsInteropFlags** | **string** | Indicates whether the view should support simultaneous access to NFS3/NFS4/SMB protocols. | [optional] [default to null]
**IsRemote** | **bool** |  | [optional] [default to null]
**ShareAcl** | [***ViewShareAcl**](View_share_acl.md) |  | [optional] [default to null]
**SelectForLiveMonitoring** | **bool** |  | [optional] [default to null]
**QosPolicyId** | **int32** | QoS Policy ID | [optional] [default to null]
**QosPolicy** | **string** | QoS Policy | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


