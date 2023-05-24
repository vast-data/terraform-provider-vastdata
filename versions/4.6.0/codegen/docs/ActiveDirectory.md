# ActiveDirectory

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **int32** |  | [optional] [default to null]
**Guid** | **string** | GUID | [optional] [default to null]
**Name** | **string** |  | [optional] [default to null]
**MachineAccountName** | **string** | Name of the computer object/machine account to add. Recommended to be the name of the cluster. | [optional] [default to null]
**OrganizationalUnit** | **string** | Organizational Unit within AD where the Cluster Machine account will be created. If left empty, it will go into default Computers OU | [optional] [default to null]
**PreferredDcList** | [***interface{}**](interface{}.md) | List of Domain Controllers to prefer for authentication. DCs listed here will be queried exclusively unless they fail or do not respond. In such a case, other DCs will be consulted. Specify as a comma-separated list. Each entry can be a fully-qualified hostname or an IP address. | [optional] [default to null]
**DomainName** | **string** | FQDN of the domain. | [optional] [default to null]
**LdapId** | **string** | the id of the attached LDAP object | [optional] [default to null]
**Enabled** | **bool** | enabled/disabled | [optional] [default to null]
**State** | **string** | Active Directory state | [optional] [default to null]
**SmbAllowed** | **bool** | Indicates if AD is allowed for SMB. There may only be 1 such AD. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


