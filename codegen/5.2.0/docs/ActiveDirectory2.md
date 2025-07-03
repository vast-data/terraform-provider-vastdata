# ActiveDirectory2

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **int32** | The unique ID of the resource. | [optional] [default to null]
**Guid** | **string** | The unique GUID of the resource. | [optional] [default to null]
**MachineAccountName** | **string** | The name of the computer object/machine account to add. Recommended to use the name of the cluster. | [optional] [default to null]
**OrganizationalUnit** | **string** | Organizational unit within Active Directory where the cluster&#x27;s machine account will be created. If left empty, defaults to Computers OU. | [optional] [default to null]
**SmbAllowed** | **bool** | Indicates if Active Directory is allowed for SMB. | [optional] [default to true]
**NtlmEnabled** | **bool** | Enables or disables support of NTLM authentication for SMB. | [optional] [default to true]
**UseAutoDiscovery** | **bool** | If &#x27;true&#x27;, Active Directory Domain Controllers (DCs) and Active Directory domains are automatically discovered. Queries extend beyond the joined domain to all domains in the forest. If &#x27;false&#x27;, queries are restricted to the joined domain and DCs must be provided in the URLs field. | [optional] [default to false]
**UseLdaps** | **bool** | Specifies whether to use LDAPS for auto-discovery. To enable use of LDAPS, also set &#x27;use_auto_discovery&#x27; to &#x27;true&#x27;. | [optional] [default to false]
**Port** | **int** | Which port to use. | [optional] [default to null]
**Binddn** | **string** | Distinguished name of the Active Directory superuser. | [optional] [default to null]
**Searchbase** | **string** | The base DN is the starting point that the Active Directory provider uses when searching for users and groups. If a group base DN is configured, it will be used instead of the base DN (for groups only). | [optional] [default to null]
**DomainName** | **string** | FQDN of the domain. | [optional] [default to null]
**Method** | **string** | Bind authentication method. | [optional] [default to METHOD.SIMPLE]
**QueryGroupsMode** | **string** | Query group mode. | [optional] [default to QUERY_GROUPS_MODE.COMPATIBLE]
**PosixAttributesSource** | **string** | Defines which domains POSIX attributes will be supported from. | [optional] [default to JOINED_DOMAIN]
**UseTls** | **bool** | Set to &#x27;true&#x27; to enable use of TLS to secure communication between the VAST cluster and the Active Directory server. | [optional] [default to false]
**TlsCertificate** | **string** | TLS certificate to use for verifying the remote Active Directory server’s TLS certificate. | [optional] [default to null]
**ReverseLookup** | **bool** | Specifies whether to resolve Active Directory netgroups into hostnames. | [optional] [default to false]
**GidNumber** | **string** | The attribute of a user entry on the Active Directory server that contains the UID number, if different from &#x27;uidNumber&#x27;. Often, when binding the VAST cluster to Active Directory, this does not need to be set. | [optional] [default to null]
**UseMultiForest** | **bool** | Allows or prohibits access for users from trusted domains on other forests. | [optional] [default to false]
**Uid** | **string** | The attribute of a user entry on the Active Directory server that contains the user name, if different from &#x27;uid&#x27;. When binding the VAST cluster to Active Directory, you may need to set this to &#x27;sAMAccountname&#x27;. | [optional] [default to null]
**UidNumber** | **string** | The attribute of a user entry on the Active Directory server that contains the UID number, if different from &#x27;uidNumber&#x27;. Often when binding the VAST cluster to Active Directory, this does not need to be set. | [optional] [default to null]
**MatchUser** | **string** | The attribute to use when querying a provider for a user that matches a user that has already been retrieved from another provider. A user entry that contains a matching value in this attribute will be considered the same user as the user previously retrieved. | [optional] [default to null]
**UidMemberValuePropertyName** | **string** | Specifies the attribute which represents the value of the Active Directory group’s member property. | [optional] [default to null]
**UidMember** | **string** | The attribute of a group entry on the Active Directory server that contains names of group members, if different from &#x27;memberUid&#x27;. When binding the VAST cluster to Active Directory, you may need to set this to &#x27;memberUID&#x27;. | [optional] [default to null]
**PosixAccount** | **string** | The object class that defines a user entry on the Active Directory server, if different from &#x27;posixAccount&#x27;. When binding the VAST cluster to Active Directory, set this parameter to &#x27;user&#x27; to ensure that authorization works properly. | [optional] [default to null]
**PosixGroup** | **string** | The object class that defines a group entry on the Active Directory server, if different from &#x27;posixGroup&#x27;. When binding the VAST cluster to Active Directory, set this parameter to &#x27;group&#x27; to ensure that authorization works properly. | [optional] [default to null]
**UsernamePropertyName** | **string** | The attribute to use for querying users in VMS user-initated user queries. Default is &#x27;name&#x27;. Sometimes it can be set to &#x27;cn&#x27;. | [optional] [default to null]
**UserLoginName** | **string** | Specifies the attribute used to query Active Directory for the user login name in NFS ID mapping. | [optional] [default to null]
**GroupLoginName** | **string** | Specifies the attribute used to query Active Directory for the group login name in NFS ID mapping. | [optional] [default to null]
**MailPropertyName** | **string** | Specifies the attribute to use for the user’s email address. | [optional] [default to null]
**IsVmsAuthProvider** | **bool** | Enables or disables use of the Active Directory for VMS authentication. Two Active Directory configurations per cluster can be used for VMS authentication: one with Active Directory and the other without Active Directory. | [optional] [default to false]
**Bindpw** | **string** | The password used with the Bind DN to authenticate to the Active Directory server. | [optional] [default to null]
**Urls** | **[]string** | A comma-separated list of URIs of Active Directory servers in the format &#x27;SCHEME://ADDRESS&#x27;. The order of listing defines the priority order. The URI with the highest priority that has a good health status is used. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

