# Role

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **int64** | The unique ID of the role. | [optional] [default to null]
**Guid** | **string** | The unique GUID  of the role. | [optional] [default to null]
**Name** | **string** | The unique name of the role. | [optional] [default to null]
**PermissionsList** | **[]string** | A list of permissions granted. | [optional] [default to null]
**Permissions** | **[]string** | A list of granted permissions returned from the VMS. | [optional] [default to null]
**Tenants** | **[]int64** | A list of tenants the role is associated with. | [optional] [default to null]
**IsAdmin** | **bool** | If true, the role is an admin role. | [optional] [default to null]
**IsDefault** | **bool** | If true, the role is a default role. | [optional] [default to null]
**LdapGroups** | **[]string** | LDAP group(s) associated with the role. Members of the specified groups on the connected LDAP/Active Directory provider can access VMS and are granted whichever permissions are included in the role. A group can be associated with multiple roles. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

