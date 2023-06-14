# LdapModifyParams

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Url** | **string** | LDAP server URI in the format &lt;scheme&gt;://&lt;address&gt;. &lt;address&gt; can be either a DNS name or an IP address. Example: ldap://ldap.company.com | [optional] [default to null]
**Urls** | [**[]ErrorUnknown**](.md) | Comma separated list of URIs of LDAP servers in the format &lt;scheme&gt;://&lt;address&gt;. The order of listing defines the priority order. The URI with highest priority that has a good health status is used. | [optional] [default to null]
**Port** | **int32** | The port of the remote LDAP server. Typical values: 389, 636. | [optional] [default to null]
**Binddn** | **string** | The bind DN for authenticating to the LDAP domain. You can specify any user account that has read access to the domain. | [optional] [default to null]
**Bindpw** | **string** | The password used with the Bind DN to authenticate to the LDAP server. | [optional] [default to null]
**Searchbase** | **string** | The entry in the LDAP directory tree to use as a starting point for user queries. | [optional] [default to null]
**GroupSearchbase** | **string** | Base DN for group queries within the joined domain only. When auto discovery is enabled, group queries outside the joined domain use auto-discovered Base DNs. | [optional] [default to null]
**Method** | **string** | The authentication method configured on the LDAP server for authenticating clients. | [optional] [default to null]
**GidNumber** | **string** | Override &#39;gidNumber&#39; as the attribute of a group entry that contains the group&#39;s GID number. When binding VAST Cluster to AD, you may need to set this to &#39;gidnumber&#39; (case sensitive). | [optional] [default to null]
**Uid** | **string** | Override &#39;uid&#39; as the attribute of a user entry on the LDAP server that contains the user name. When binding VAST Cluster to AD, you may need to set this to &#39;sAMAccountname&#39;. | [optional] [default to null]
**UidNumber** | **string** | Override &#39;uidNumber&#39; as the attribute of a user entry on the LDAP server that contains the UID number. Often when binding VAST Cluster to AD this does not need to be set. | [optional] [default to null]
**MatchUser** | **string** | The attribute to use when querying a provider for a user that matches a user that was already retrieved from another provider. A user entry that contains a matching value in this attribute will be considered the same user as the user previously retrieved. | [optional] [default to null]
**UidMember** | **string** | Override &#39;memberUid&#39; as the attribute of a group entry on the LDAP server that contains names of group members. When binding VAST Cluster to AD, you may need to set this to &#39;memberUID&#39; | [optional] [default to null]
**PosixAccount** | **string** | Override &#39;posixAccount&#39;as the object class that defines a user entry on the LDAP server. When binding VAST Cluster to AD, set this parameter to &#39;user&#39; in order for authorization to work properly. | [optional] [default to null]
**PosixGroup** | **string** | Override &#39;posixGroup&#39; as the object class that defines a group entry on the LDAP server. When binding VAST Cluster to AD, set this parameter to &#39;group&#39; in order for authorization to work properly. | [optional] [default to null]
**TlsCertificate** | **string** | TLS certificate to use for verifying the remote LDAP server&#39;s TLS certificate. | [optional] [default to null]
**UseTls** | **bool** | Set to true to enable use of TLS to secure communication between VAST Cluster and the LDAP server. | [optional] [default to null]
**UsePosix** | **bool** | POSIX support | [optional] [default to null]
**QueryGroupsMode** | **string** | A mode setting for how groups are queried: Set to COMPATIBLE to look up user groups using the &#39;memberOf&#39; and &#39;memberUid&#39; attributes. Set to RFC2307BIS_ONLY to look up user groups using only the &#39;memberOf&#39; attribute. Set to RFC2307_ONLY to look up user groups using only the &#39;memberUid&#39; attribute. Set to NONE not to look up user groups other than by leading GID and primary group SID. | [optional] [default to null]
**UsernamePropertyName** | **string** | The attribute to use for querying users in VMS user-initated user queries. Default is &#39;name&#39;. Sometimes set to &#39;cn&#39; | [optional] [default to null]
**DomainName** | **string** | FQDN of AD domain. Must be resolvable in DNS | [optional] [default to null]
**UserLoginName** | **string** | The attribute used to query AD for the user login name in NFS ID mapping. Applicable only with AD and NFSv4.1. | [optional] [default to null]
**GroupLoginName** | **string** | The attribute used to query AD for the group login name in NFS ID mapping. Applicable only with AD and NFSv4.1. | [optional] [default to null]
**MailPropertyName** | **string** |  | [optional] [default to null]
**UseAutoDiscovery** | **bool** | When enabled, Active Directory Domain Controllers (DCs) and Active Directory domains are auto discovered. Queries extend beyond the joined domain to all domains in the forest. When disabled, queries are restricted to the joined domain and DCs must be provided in the URLs field. | [optional] [default to null]
**UseLdaps** | **bool** | Use LDAPS for Auto-Discovery | [optional] [default to null]
**IsVmsAuthProvider** | **bool** | Whether the LDAP should be used for VMS auth. There is only two LDAPs allowed for VMS auth: one with AD and one w/o. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


