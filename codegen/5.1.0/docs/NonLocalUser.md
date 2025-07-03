# NonLocalUser

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | The ID of the non-local user. | [optional] [default to null]
**Uid** | **int32** | The user&#x27;&#x27;s Unix UID. | [optional] [default to null]
**Username** | **string** | The  username of the non-local user. | [optional] [default to null]
**AllowCreateBucket** | **bool** | Allows or prohibits bucket creation by the user. | [optional] [default to null]
**AllowDeleteBucket** | **bool** | Allows or prohibits bucket deletion by the user. | [optional] [default to null]
**TenantId** | **int32** | Tenant ID. | [optional] [default to null]
**S3PoliciesIds** | **[]int32** | A list of identity policy IDs. | [optional] [default to null]
**Context** | **string** | Context from which the user originates. Valid values: &#x27;ad&#x27;, &#x27;nis&#x27; and &#x27;ldap&#x27;. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

