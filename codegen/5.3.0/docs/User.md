# User

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **int32** | A uniq id given to user | [optional] [default to null]
**Guid** | **string** | A uniq guid given to the user | [optional] [default to null]
**Name** | **string** | A uniq name given to the user | [optional] [default to null]
**Uid** | **int32** | The user unix UID | [optional] [default to null]
**LeadingGid** | **int32** | The user leading unix GID | [optional] [default to null]
**Gids** | **[]int32** | List of supplementary GID list | [optional] [default to null]
**Groups** | **[]string** | List of supplementary Group list | [optional] [default to null]
**GroupCount** | **int32** | Group Count | [optional] [default to null]
**LeadingGroupName** | **string** | Leading Group Name | [optional] [default to null]
**LeadingGroupGid** | **int** | Leading Group GID | [optional] [default to null]
**Sid** | **string** | The user SID | [optional] [default to null]
**PrimaryGroupSid** | **string** | The user primary group SID | [optional] [default to null]
**Sids** | **[]string** | supplementary SID list | [optional] [default to null]
**Local** | **bool** | IS this a local user | [optional] [default to null]
**AllowCreateBucket** | **bool** | Allow create bucket | [optional] [default to null]
**AllowDeleteBucket** | **bool** | Allow delete bucket | [optional] [default to null]
**S3Superuser** | **bool** | Is S3 superuser | [optional] [default to null]
**S3PoliciesIds** | **[]int32** | List S3 policies IDs | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

