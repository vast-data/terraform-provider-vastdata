# ActiveDirectory2

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **int32** | A uniqe ID given to this resource | [optional] [default to null]
**Guid** | **string** | A uniqe ID given to this resource | [optional] [default to null]
**MachineAccountName** | **string** | Name of the computer object/machine account to add. Recommended to be the name of the cluster | [optional] [default to null]
**OrganizationalUnit** | **string** | Organizational Unit within AD where the Cluster Machine account will be created. If left empty, it will go into default Computers OU | [optional] [default to null]
**SmbAllowed** | **bool** | Indicates if AD is allowed for SMB. | [optional] [default to true]
**NtlmEnabled** | **bool** | Manages support of NTLM authentication method for SMB protocol. | [optional] [default to true]
**UseAutoDiscovery** | **bool** | When enabled, Active Directory Domain Controllers (DCs) and Active Directory domains are auto discovered. Queries extend beyond the joined domain to all domains in the forest. When disabled, queries are restricted to the joined domain and DCs must be provided in the URLs field. | [optional] [default to false]
**UseLdaps** | **bool** | Use LDAPS for auto-Discovery. To activate, set use_auto_discovery to true also. | [optional] [default to false]
**Port** | **int** | Which port to use | [optional] [default to null]
**Binddn** | **string** | Distinguished name of AD superuser | [optional] [default to null]
**Searchbase** | **string** | The Base DN is the starting point the AD provider uses when searching for users and groups. If the Group Base DN is configured it will be used instead of the Base DN, for groups only | [optional] [default to null]
**DomainName** | **string** | FQDN of the domain. | [optional] [default to null]
**Method** | **string** | Bind Authentication Method | [optional] [default to METHOD.SIMPLE]
**QueryGroupsMode** | **string** | Query group mode | [optional] [default to QUERY_GROUPS_MODE.COMPATIBLE]
**PosixAttributesSource** | **string** | Defines which domains POSIX attributes will be supported from. | [optional] [default to JOINED_DOMAIN]
**UseTls** | **bool** | Set to true to enable use of TLS to secure communication between VAST Cluster and the AD server. | [optional] [default to false]
**TlsCertificate** | **string** | TLS certificate to use for verifying the remote AD server’s TLS certificate. | [optional] [default to null]
**ReverseLookup** | **bool** | Resolve AD netgroups into hostnames | [optional] [default to false]
**GidNumber** | **string** | The attribute of a user entry on the AD server that contains the UID number, if different from ‘uidNumber’. Often when binding VAST Cluster to AD this does not need to be set. | [optional] [default to null]
**UseMultiForest** | **bool** | Allow access for users from trusted domains on other forests. | [optional] [default to false]
**Uid** | **string** | The attribute of a user entry on the AD server that contains the user name, if different from ‘uid’ When binding VAST Cluster to AD, you may need to set this to ‘sAMAccountname’. | [optional] [default to null]
**UidNumber** | **string** | The attribute of a user entry on the AD server that contains the UID number, if different from ‘uidNumber’. Often when binding VAST Cluster to AD this does not need to be set. | [optional] [default to null]
**MatchUser** | **string** | The attribute to use when querying a provider for a user that matches a user that was already retrieved from another provider. A user entry that contains a matching value in this attribute will be considered the same user as the user previously retrieved. | [optional] [default to null]
**UidMemberValuePropertyName** | **string** | Specifies the attribute which represents the value of the AD group’s member property. | [optional] [default to null]
**UidMember** | **string** | The attribute of a group entry on the AD server that contains names of group members, if different from ‘memberUid’. When binding VAST Cluster to AD, you may need to set this to ‘memberUID’. | [optional] [default to null]
**PosixAccount** | **string** | The object class that defines a user entry on the AD server, if different from ‘posixAccount’. When binding VAST Cluster to AD, set this parameter to ‘user’ in order for authorization to work properly. | [optional] [default to null]
**PosixGroup** | **string** |  The object class that defines a group entry on the AD server, if different from ‘posixGroup’. When binding VAST Cluster to AD, set this parameter to ‘group’ in order for authorization to work properly. | [optional] [default to null]
**UsernamePropertyName** | **string** | The attribute to use for querying users in VMS user-initated user queries. Default is ‘name’. Sometimes set to ‘cn’ | [optional] [default to null]
**UserLoginName** | **string** | Specifies the attribute used to query AD for the user login name in NFS ID mapping. Applicable only with AD and NFSv4.1. | [optional] [default to null]
**GroupLoginName** | **string** | Specifies the attribute used to query AD for the group login name in NFS ID mapping. Applicable only with AD and NFSv4.1. | [optional] [default to null]
**MailPropertyName** | **string** | Specifies the attribute to use for the user’s email address. | [optional] [default to null]
**IsVmsAuthProvider** | **bool** | Enables use of the AD for VMS authentication. Two AD configurations per cluster can be used for VMS authentication: one with AD and one without. | [optional] [default to false]
**Bindpw** | **string** | The password used with the Bind DN to authenticate to the AD server. | [optional] [default to null]
**Urls** | **[]string** | Comma separated list of URIs of AD servers in the format SCHEME://ADDRESS. The order of listing defines the priority order. The URI with highest priority that has a good health status is used. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

