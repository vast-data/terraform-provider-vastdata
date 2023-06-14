# UserModifyParams

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**LeadingGid** | **int32** | Leading GID | [optional] [default to null]
**Gids** | **[]int32** | GID list | [optional] [default to null]
**Uid** | **int32** | NFS UID | [optional] [default to null]
**Local** | **bool** |  | [optional] [default to null]
**AllowCreateBucket** | **bool** | Set to true to give the user permission to create S3 buckets | [optional] [default to null]
**AllowDeleteBucket** | **bool** | Set to true to give the user permission to delete S3 buckets | [optional] [default to null]
**S3Superuser** | **bool** | Set to true to give the user S3 superuser permission | [optional] [default to null]
**S3PoliciesIds** | [***interface{}**](interface{}.md) | list of s3 policy ids | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


