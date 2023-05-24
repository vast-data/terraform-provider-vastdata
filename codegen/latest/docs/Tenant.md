# Tenant

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **int32** | A uniq id given to the tenant | [optional] [default to null]
**Guid** | **string** | A uniq guid given to the tenant | [optional] [default to null]
**Name** | **string** | A uniq name given to the tenant | [optional] [default to null]
**Sync** | **string** | Synchronization state with leader | [optional] [default to null]
**SyncTime** | **string** | Synchronization time with leader | [optional] [default to null]
**SmbPrivilegedUserName** | **string** | Optional custom username for the SMB privileged user. If not set, the SMB privileged user name is &#x27;vastadmin&#x27; | [optional] [default to null]
**SmbPrivilegedGroupSid** | **string** | Optional custom SID to specify a non default SMB privileged group. If not set, SMB privileged group is the Backup Operators domain group. | [optional] [default to null]
**SmbAdministratorsGroupName** | **string** | Optional custom name to specify a non default privileged group. If not set, privileged group is the Backup Operators domain group. | [optional] [default to null]
**DefaultOthersShareLevelPerm** | **string** | Default Share-level permissions for Others | [optional] [default to null]
**TrashGid** | **int32** | GID with permissions to the trash folder | [optional] [default to null]
**ClientIpRanges** | [**[][]string**](array.md) | Array of source IP ranges to allow for the tenant. | [optional] [default to null]
**PosixPrimaryProvider** | **string** | POSIX primary provider type | [optional] [default to null]
**AdProviderId** | **int32** | AD provider ID | [optional] [default to null]
**LdapProviderId** | **int32** | Open-LDAP provider ID specified separately by the user | [optional] [default to null]
**NisProviderId** | **int32** | NIS provider ID | [optional] [default to null]
**EncryptionCrn** | **string** | Tenant&#x27;s encryption group unique identifier | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

