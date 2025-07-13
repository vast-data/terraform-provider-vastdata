# Ldap

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **int32** |  | [optional] [default to null]
**Guid** | **string** |  | [optional] [default to null]
**Urls** | **[]string** | A list of URIs of LDAP servers (Domain Controllers in Active Directory), in priority order. The URI with highest priority that has a good health status is used. Specify each URI in the format &#x27;&lt;scheme&gt;://&lt;address&gt;&#x27;. &#x27;&lt;address&gt;&#x27; can be either a DNS name or an IP address, for example: &#x27;ldap://ldap.company.com, ldaps://ldaps.company.com, ldap://192.0.2.2&#x27;. | [optional] [default to null]
**Port** | **int32** | LDAP server port. 389 (LDAP)  636 (LDAPS) | [optional] [default to 389]
**Binddn** | **string** | Distinguished name of the LDAP superuser. | [optional] [default to null]
**Bindpw** | **string** | Password for the LDAP superuser. | [optional] [default to null]
**Searchbase** | **string** | The base DN is the starting point that the LDAP provider uses when searching for users and groups. If a group base DN is configured, it will be used instead of the base DN, for groups only. | [optional] [default to null]
**GroupSearchbase** | **string** | Base DN for group queries within the joined domain only. When auto-discovery is enabled, group queries outside the joined domain use automatically discovered base DNs. | [optional] [default to null]
**Method** | **string** | Bind authentication method. | [optional] [default to null]
**GidNumber** | **string** | Attribute mapping for gid number. | [optional] [default to gidNumber]
**Uid** | **string** | Attribute mapping for uid. | [optional] [default to uid]
**UidNumber** | **string** | Attribute mapping for uid number. | [optional] [default to uidNumber]
**MatchUser** | **string** | Attribute mapping for user matching. | [optional] [default to uid]
**UidMember** | **string** | Attribute mapping for uid member. | [optional] [default to memberUID]
**PosixAccount** | **string** | Attribute mapping for posix account. | [optional] [default to posixAccount]
**PosixGroup** | **string** | Attribute mapping for posix group. | [optional] [default to posixGroup]
**UseTls** | **bool** | Specifies whether to configure LDAP with TLS. | [optional] [default to false]
**PosixPrimaryProvider** | **bool** | POSIX primary provider. | [optional] [default to null]
**PosixAttributesSource** | **string** |  | [optional] [default to JOINED_DOMAIN]
**ReverseLookup** | **bool** |  | [optional] [default to false]
**TlsCertificate** | **string** |  | [optional] [default to null]
**ActiveDirectory** | **string** |  | [optional] [default to null]
**QueryGroupsMode** | **string** | Query group mode. | [optional] [default to null]
**UsernamePropertyName** | **string** | Username property name. | [optional] [default to cn]
**DomainName** | **string** | FQDN of the domain. | [optional] [default to null]
**UserLoginName** | **string** | The attribute used to query the provider for the user login name in NFS ID mapping. | [optional] [default to uid]
**GroupLoginName** | **string** | The attribute used to query the provider for the group login name in NFS ID mapping. | [optional] [default to cn]
**MailPropertyName** | **string** |  | [optional] [default to mail]
**UidMemberValuePropertyName** | **string** |  | [optional] [default to uid]
**UseAutoDiscovery** | **bool** | If &#x27;true&#x27;, Active Directory Domain Controllers (DCs) and Active Directory domains are automatically discovered. Queries extend beyond the joined domain to all domains in the forest. If &#x27;false&#x27;, queries are restricted to the joined domain and URIs must be provided in &#x27;urls&#x27;. | [optional] [default to null]
**UseLdaps** | **bool** | Specifies whether to use LDAPS for auto-discovery. | [optional] [default to null]
**IsVmsAuthProvider** | **bool** | Specifies whether the LDAP is to be used for VMS authentication. There can be only two LDAP configurations that can be used for VMS authentication: one with Active Directory and the other without Active Directory. | [optional] [default to false]
**QueryPosixAttributesFromGc** | **bool** | If &#x27;true&#x27;, users/groups from non-joined domain POSIX attributes are supported. If &#x27;false&#x27;, POSIX attributes of users/groups from non-joined domain are not supported. | [optional] [default to false]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

