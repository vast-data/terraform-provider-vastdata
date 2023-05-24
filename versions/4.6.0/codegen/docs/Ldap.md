# Ldap

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **int32** |  | [optional] [default to null]
**Name** | **string** |  | [optional] [default to null]
**Guid** | **string** |  | [optional] [default to null]
**Url** | **string** | Comma-separated list of URIs of LDAP servers (Domain Controllers (DCs) in Active Directory), in priority order. The URI with highest priority that has a good health status is used. Specify each URI in the format &lt;scheme&gt;://&lt;address&gt;. &lt;address&gt; can be either a DNS name or an IP address. e.g. ldap://ldap.company.com, ldaps://ldaps.company.com, ldap://192.0.2.2 | [optional] [default to null]
**Urls** | **[]string** | Comma-separated list of URIs of LDAP servers (Domain Controllers (DCs) in Active Directory), in priority order. The URI with highest priority that has a good health status is used. Specify each URI in the format &lt;scheme&gt;://&lt;address&gt;. &lt;address&gt; can be either a DNS name or an IP address. e.g. ldap://ldap.company.com, ldaps://ldaps.company.com, ldap://192.0.2.2 | [optional] [default to null]
**Port** | **int32** | LDAP server port. 389 (LDAP)  636 (LDAPS) | [optional] [default to null]
**Binddn** | **string** | Distinguished name of LDAP superuser | [optional] [default to null]
**Bindpw** | **string** | Password for the LDAP superuser | [optional] [default to null]
**Searchbase** | **string** | The Base DN is the starting point the LDAP provider uses when searching for users and groups. If the Group Base DN is configured it will be used instead of the Base DN, for groups only | [optional] [default to null]
**GroupSearchbase** | **string** | Base DN for group queries within the joined domain only. When auto discovery is enabled, group queries outside the joined domain use auto-discovered Base DNs. | [optional] [default to null]
**Method** | **string** | Bind Authentication Method | [optional] [default to null]
**State** | **string** |  | [optional] [default to null]
**TenantId** | **int32** | Tenant ID | [optional] [default to null]
**GidNumber** | **string** |  | [optional] [default to null]
**Uid** | **string** |  | [optional] [default to null]
**UidNumber** | **string** |  | [optional] [default to null]
**MatchUser** | **string** |  | [optional] [default to null]
**UidMember** | **string** |  | [optional] [default to null]
**PosixAccount** | **string** |  | [optional] [default to null]
**PosixGroup** | **string** |  | [optional] [default to null]
**UseTls** | **bool** | configure LDAP with TLS | [optional] [default to null]
**UsePosix** | **bool** | POSIX support | [optional] [default to null]
**PosixPrimaryProvider** | **bool** | POSIX primary provider | [optional] [default to null]
**TlsCertificate** | **string** |  | [optional] [default to null]
**ActiveDirectory** | **string** |  | [optional] [default to null]
**QueryGroupsMode** | **string** | Query group mode | [optional] [default to null]
**UsernamePropertyName** | **string** | Username property name | [optional] [default to null]
**DomainName** | **string** | FQDN of the domain. | [optional] [default to null]
**UserLoginName** | **string** | The attribute used to query AD for the user login name in NFS ID mapping. Applicable only with AD and NFSv4.1. | [optional] [default to null]
**GroupLoginName** | **string** | The attribute used to query AD for the group login name in NFS ID mapping. Applicable only with AD and NFSv4.1. | [optional] [default to null]
**MailPropertyName** | **string** |  | [optional] [default to null]
**UseAutoDiscovery** | **bool** | When enabled, Active Directory Domain Controllers (DCs) and Active Directory domains are auto discovered. Queries extend beyond the joined domain to all domains in the forest. When disabled, queries are restricted to the joined domain and DCs must be provided in the URLs field. | [optional] [default to null]
**UseLdaps** | **bool** | Use LDAPS for Auto-Discovery | [optional] [default to null]
**IsVmsAuthProvider** | **bool** | Whether the LDAP should be used for VMS auth. There is only two LDAPs allowed for VMS auth: one with AD and one w/o. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


