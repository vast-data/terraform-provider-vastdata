# User

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **int32** | The unique ID of the user. | [optional] [default to null]
**Guid** | **string** | The unique GUID of the user. | [optional] [default to null]
**Name** | **string** | The unique name of the user. | [optional] [default to null]
**Uid** | **int32** | The user&#x27;s Unix UID. | [optional] [default to null]
**LeadingGid** | **int32** | The user&#x27;s leading Unix GID. | [optional] [default to null]
**Gids** | **[]int32** | A list of supplementary GIDs. | [optional] [default to null]
**Groups** | **[]string** | A list of supplementary groups. | [optional] [default to null]
**GroupCount** | **int32** | Group count. | [optional] [default to null]
**LeadingGroupName** | **string** | The name of the leading group. | [optional] [default to null]
**LeadingGroupGid** | **int** | The GID of the leading group. | [optional] [default to null]
**Sid** | **string** | The user&#x27;s SID. | [optional] [default to null]
**PrimaryGroupSid** | **string** | The user&#x27;s primary group SID. | [optional] [default to null]
**Sids** | **[]string** | A list of supplementary SIDs. | [optional] [default to null]
**Local** | **bool** | If &#x27;true&#x27;, the user is a local user. | [optional] [default to null]
**AllowCreateBucket** | **bool** | Allows or prohibits bucket creation by the user. | [optional] [default to null]
**AllowDeleteBucket** | **bool** | Allows or prohibits bucket deletion by the user. | [optional] [default to null]
**S3Superuser** | **bool** | If &#x27;true&#x27;, the user is an S3 superuser. | [optional] [default to null]
**S3PoliciesIds** | **[]int32** | A list of identity policy IDs. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

