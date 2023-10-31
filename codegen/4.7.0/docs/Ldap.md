# Ldap

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **int32** |  | [optional] [default to null]
**Guid** | **string** |  | [optional] [default to null]
**Urls** | **[]string** | List of URIs of LDAP servers (Domain Controllers (DCs) in Active Directory), in priority order. The URI with highest priority that has a good health status is used. Specify each URI in the format &lt;scheme&gt;://&lt;address&gt;. &lt;address&gt; can be either a DNS name or an IP address. e.g. ldap://ldap.company.com, ldaps://ldaps.company.com, ldap://192.0.2.2 | [optional] [default to null]
**Port** | **int32** | LDAP server port. 389 (LDAP)  636 (LDAPS) | [optional] [default to 389]
**Binddn** | **string** | Distinguished name of LDAP superuser | [optional] [default to null]
**Bindpw** | **string** | Password for the LDAP superuser | [optional] [default to null]
**Searchbase** | **string** | The Base DN is the starting point the LDAP provider uses when searching for users and groups. If the Group Base DN is configured it will be used instead of the Base DN, for groups only | [optional] [default to null]
**GroupSearchbase** | **string** | Base DN for group queries within the joined domain only. When auto discovery is enabled, group queries outside the joined domain use auto-discovered Base DNs. | [optional] [default to null]
**Method** | **string** | Bind Authentication Method | [optional] [default to null]
**GidNumber** | **string** | Attrirbute mapping for gid number | [optional] [default to gidNumber]
**Uid** | **string** | Attrirbute mapping for uid | [optional] [default to uid]
**UidNumber** | **string** | Attrirbute mapping for uid number | [optional] [default to uidNumber]
**MatchUser** | **string** | Attribute mapping for user matching | [optional] [default to uid]
**UidMember** | **string** | Attrirbute mapping for uid member | [optional] [default to memberUID]
**PosixAccount** | **string** | Attrirbute mapping for posix account | [optional] [default to posixAccount]
**PosixGroup** | **string** | Attrirbute mapping for posix account | [optional] [default to posixGroup]
**UseTls** | **bool** | configure LDAP with TLS | [optional] [default to false]
**PosixPrimaryProvider** | **bool** | POSIX primary provider | [optional] [default to null]
**PosixAttributesSource** | **string** |  | [optional] [default to JOINED_DOMAIN]
**ReverseLookup** | **bool** |  | [optional] [default to false]
**TlsCertificate** | **string** |  | [optional] [default to null]
**ActiveDirectory** | **string** |  | [optional] [default to null]
**QueryGroupsMode** | **string** | Query group mode | [optional] [default to null]
**UsernamePropertyName** | **string** | Username property name | [optional] [default to cn]
**DomainName** | **string** | FQDN of the domain. | [optional] [default to null]
**UserLoginName** | **string** | The attribute used to query AD for the user login name in NFS ID mapping. Applicable only with AD and NFSv4.1. | [optional] [default to uid]
**GroupLoginName** | **string** | The attribute used to query AD for the group login name in NFS ID mapping. Applicable only with AD and NFSv4.1. | [optional] [default to cn]
**MailPropertyName** | **string** |  | [optional] [default to mail]
**UidMemberValuePropertyName** | **string** |  | [optional] [default to uid]
**UseAutoDiscovery** | **bool** | When enabled, Active Directory Domain Controllers (DCs) and Active Directory domains are auto discovered. Queries extend beyond the joined domain to all domains in the forest. When disabled, queries are restricted to the joined domain and DCs must be provided in the URLs field. | [optional] [default to null]
**UseLdaps** | **bool** | Use LDAPS for Auto-Discovery | [optional] [default to null]
**IsVmsAuthProvider** | **bool** | Whether the LDAP should be used for VMS auth. There is only two LDAPs allowed for VMS auth: one with AD and one w/o. | [optional] [default to false]
**QueryPosixAttributesFromGc** | **bool** | When set to True - users/groups from non-joined domain POSIX attributes are supported, when set to False - Posix attributes of users/groups from non-joined domain are not supported. As a condition Global catalog needs to be configured to support Posix attributes.  | [optional] [default to false]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

