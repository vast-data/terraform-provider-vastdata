# Tenant

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **int32** | The unique ID of the tenant. | [optional] [default to null]
**Guid** | **string** | The unique GUID of the tenant. | [optional] [default to null]
**Name** | **string** | The unique name of the tenant. | [optional] [default to null]
**UseSmbPrivilegedUser** | **bool** | Enables or disables SMB privileged user. | [optional] [default to null]
**SmbPrivilegedUserName** | **string** | An optional custom username for the SMB privileged user. If not set, the SMB privileged user name is &#x27;vastadmin&#x27;. | [optional] [default to null]
**UseSmbPrivilegedGroup** | **bool** | Enables or disables SMB privileged user group. | [optional] [default to null]
**SmbPrivilegedGroupSid** | **string** | An optional custom SID to specify a non-default SMB privileged group. If not set, the SMB privileged group is the Backup Operators domain group. | [optional] [default to null]
**SmbPrivilegedGroupFullAccess** | **bool** | If &#x27;true&#x27;, the SMB privileged user group has read and write control access. Members of the group can perform backup and restore operations on all files and directories, without requiring read or write access to the specific files and directories. If &#x27;false&#x27;, the privileged group has read-only access. | [optional] [default to null]
**SmbAdministratorsGroupName** | **string** | An optional custom name to specify a non-default privileged group. If not set, the privileged group is the Backup Operators domain group. | [optional] [default to null]
**DefaultOthersShareLevelPerm** | **string** | Default share-level permissions for others. | [optional] [default to null]
**TrashGid** | **int32** | A GID with permissions to the trash folder. | [optional] [default to null]
**ClientIpRanges** | [**[][]string**](array.md) | An array of source IP ranges to allow for the tenant. | [optional] [default to null]
**PosixPrimaryProvider** | **string** | The POSIX primary provider type. | [optional] [default to null]
**AdProviderId** | **int32** | The ID of the Active Directory provider. | [optional] [default to null]
**LdapProviderId** | **int32** | The ID of the OpenLDAP provider specified separately by the user. | [optional] [default to null]
**NisProviderId** | **int32** | The NIS provider ID. | [optional] [default to null]
**EncryptionCrn** | **string** | The unique ID of the tenant&#x27;s encryption group. | [optional] [default to null]
**IsNfsv42Supported** | **bool** | Enables or disables NFSv4.2. | [optional] [default to null]
**AllowLockedUsers** | **bool** | Allows or prohibits IO from users whose Active Directory accounts are locked out by lockout policies due to unsuccessful login attempts. | [optional] [default to false]
**AllowDisabledUsers** | **bool** | Allows or prohibits IO from users whose Active Directory accounts are explicitly disabled. | [optional] [default to false]
**UseSmbNative** | **bool** | Enables or disables use of native SMB authentication. | [optional] [default to null]
**VippoolNames** | **[]string** | An array of names of virtual IP pools attached to the tenant. | [optional] [default to null]
**VippoolIds** | **[]int64** | An array of IDs of virtual IP pools attached to the tenant. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

