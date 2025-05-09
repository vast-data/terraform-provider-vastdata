# Role

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **int64** | A unique id given to the role | [optional] [default to null]
**Name** | **string** | A uniqe name of the role | [optional] [default to null]
**PermissionsList** | **[]string** | List of allowed permissions | [optional] [default to null]
**Permissions** | **[]string** | List of allowed permissions returned from the VMS | [optional] [default to null]
**Tenants** | **[]int64** | List of tenants to which this role is associated with | [optional] [default to null]
**IsAdmin** | **bool** | Is the role is an admin role | [optional] [default to null]
**IsDefault** | **bool** | Is the role is a default role | [optional] [default to null]
**LdapGroups** | **[]string** | LDAP group(s) associated with the role. Members of the specified groups on a connected LDAP/Active Directory provider can access VMS and are granted whichever permissions are included in the role. A group can be associated with multiple roles. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

