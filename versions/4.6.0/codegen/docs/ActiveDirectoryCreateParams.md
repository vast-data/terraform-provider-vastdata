# ActiveDirectoryCreateParams

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**LdapId** | **string** | ID of the LDAP configuration for binding to the LDAP domain of the AD server. | [optional] [default to null]
**DomainName** | **string** | The fully qualified domain name (FQDN) of the Active Directory domain to join. | [optional] [default to null]
**MachineAccountName** | **string** | The name for the machine object representing the VAST Cluster to be created within the OU | [optional] [default to null]
**OrganizationalUnit** | **string** | A non default organizational unit (OU) in the AD domain in which to create the machine object. If left empty, the machine object will be created in the default Computers OU. | [optional] [default to null]
**PreferredDcList** | [**[]ErrorUnknown**](.md) | Specify multiple DCs using &#39;urls&#39; parameter in LDAP configuration. | [optional] [default to null]
**SmbAllowed** | **bool** | Indicates if AD is allowed for SMB. There may only be 1 such AD. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


