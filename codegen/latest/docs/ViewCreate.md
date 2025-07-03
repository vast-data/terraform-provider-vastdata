# ViewCreate

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | A unique name for the view. | [optional] [default to null]
**Path** | **string** | File system path. Begin with &#x27;/&#x27;. Do not include a trailing slash. | [default to null]
**CreateDir** | **bool** | If &#x27;true&#x27;, creates the directory specified by the path. | [optional] [default to null]
**Alias** | **string** | Alias for NFS export. The alias must start with &#x27;/&#x27; and can include ASCII characters only. If configured, this setting supersedes the exposed NFS export path. | [optional] [default to null]
**Bucket** | **string** | S3 bucket name. | [optional] [default to null]
**PolicyId** | **int32** | The ID of the view policy to be associated with the view. | [optional] [default to null]
**Cluster** | **string** | Parent cluster. | [optional] [default to null]
**ClusterId** | **int32** | Parent cluster ID. | [optional] [default to null]
**TenantId** | **int32** | The ID of the tenant associated with this view. | [optional] [default to null]
**Directory** | **bool** | If &#x27;true&#x27;, creates the directory if it does not exist. | [optional] [default to null]
**S3Versioning** | **bool** | Enables or disables S3 versioning. | [optional] [default to null]
**S3UnverifiedLookup** | **bool** | Allows or prohibits S3 Unverified Lookup. | [optional] [default to null]
**AllowAnonymousAccess** | **bool** | Allows or prohibits S3 anonymous access. | [optional] [default to null]
**AllowS3AnonymousAccess** | **bool** | Allows or prohibits S3 anonymous access. | [optional] [default to null]
**Protocols** | **[]string** | Protocols exposed by this view. | [optional] [default to null]
**Share** | **string** | Name of the SMB share. The name cannot not include the following characters: \&quot; \\ / [ ] : | &lt; &gt; + &#x3D; ; , * ? | [optional] [default to null]
**BucketOwner** | **string** | S3 bucket owner. | [optional] [default to null]
**BucketCreators** | **[]string** | A list of bucket creator users. | [optional] [default to null]
**BucketCreatorsGroups** | **[]string** | A list of bucket creator groups. | [optional] [default to null]
**S3Locks** | **bool** | Enables or disables S3 object locks. | [optional] [default to null]
**S3LocksRetentionMode** | **string** | S3 locks retention mode. | [optional] [default to null]
**S3LocksRetentionPeriod** | **string** | Retention period for S3 locks. The period is specified as a positive integer suffixed by a time unit of measure, for example: &#x27;0d&#x27;|&#x27;2d&#x27;|&#x27;1y&#x27;|&#x27;2y&#x27; | [optional] [default to null]
**PhysicalCapacity** | **int64** | Physical capacity. | [optional] [default to null]
**LogicalCapacity** | **int64** | Logical capacity. | [optional] [default to null]
**NfsInteropFlags** | **string** | Indicates whether the view supports simultaneous access using NFSv3/NFSv4/SMB protocols. | [optional] [default to null]
**IsRemote** | **bool** |  | [optional] [default to null]
**ShareAcl** | [***ViewShareAcl**](View_share_acl.md) |  | [optional] [default to null]
**QosPolicyId** | **int32** | The ID of the QoS policy associated with the view. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

