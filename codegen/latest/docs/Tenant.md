# Tenant

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **int32** | A uniq id given to the tenant | [optional] [default to null]
**Guid** | **string** | A uniq guid given to the tenant | [optional] [default to null]
**Name** | **string** | A uniq name given to the tenant | [optional] [default to null]
**UseSmbPrivilegedUser** | **bool** | Enables SMB privileged user | [optional] [default to null]
**SmbPrivilegedUserName** | **string** | Optional custom username for the SMB privileged user. If not set, the SMB privileged user name is &#x27;vastadmin&#x27; | [optional] [default to null]
**UseSmbPrivilegedGroup** | **bool** | Enables SMB privileged user group | [optional] [default to null]
**SmbPrivilegedGroupSid** | **string** | Optional custom SID to specify a non default SMB privileged group. If not set, SMB privileged group is the Backup Operators domain group. | [optional] [default to null]
**SmbPrivilegedGroupFullAccess** | **bool** | True&#x3D;The SMB privileged user group has read and write control access. Members of the group can perform backup and restore operations on all files and directories, without requiring read or write access to the specific files and directories. False&#x3D;the privileged group has read only access. | [optional] [default to null]
**SmbAdministratorsGroupName** | **string** | Optional custom name to specify a non default privileged group. If not set, privileged group is the Backup Operators domain group. | [optional] [default to null]
**DefaultOthersShareLevelPerm** | **string** | Default Share-level permissions for Others | [optional] [default to null]
**TrashGid** | **int32** | GID with permissions to the trash folder | [optional] [default to null]
**ClientIpRanges** | [**[][]string**](array.md) | Array of source IP ranges to allow for the tenant. | [optional] [default to null]
**PosixPrimaryProvider** | **string** | POSIX primary provider type | [optional] [default to null]
**AdProviderId** | **int32** | AD provider ID | [optional] [default to null]
**LdapProviderId** | **int32** | Open-LDAP provider ID specified separately by the user | [optional] [default to null]
**NisProviderId** | **int32** | NIS provider ID | [optional] [default to null]
**EncryptionCrn** | **string** | Tenant&#x27;s encryption group unique identifier | [optional] [default to null]
**IsNfsv42Supported** | **bool** | Enable NFSv4.2 | [optional] [default to null]
**AllowLockedUsers** | **bool** | Allow IO from users whose Active Directory accounts are locked out by lockout policies due to unsuccessful login attempts. | [optional] [default to false]
**AllowDisabledUsers** | **bool** | Allow IO from users whose Active Directory accounts are explicitly disabled. | [optional] [default to false]
**UseSmbNative** | **bool** | Use native SMB authentication | [optional] [default to null]
**VippoolNames** | **[]string** | An array of VIP Pool names attached to this tenant. | [optional] [default to null]
**VippoolIds** | **[]int64** | An array of VIP Pool ids to attach to tenant. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

