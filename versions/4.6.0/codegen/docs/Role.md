# Role

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **int32** |  | [optional] [default to null]
**Name** | **string** | The friendly name of the role | [optional] [default to null]
**Managers** | [***interface{}**](interface{}.md) |  | [optional] [default to null]
**Permissions** | [***interface{}**](interface{}.md) |  | [optional] [default to null]
**IsDefault** | **bool** | Is the role is a default role | [optional] [default to null]
**LdapGroups** | **[]string** | Ldap groups which will be granted this role in login | [optional] [default to null]
**Tenants** | **[]string** | Tenants for that role | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


