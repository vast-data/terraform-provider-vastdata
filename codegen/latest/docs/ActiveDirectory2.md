# ActiveDirectory2

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **int32** |  | [optional] [default to null]
**Guid** | **string** |  | [optional] [default to null]
**MachineAccountName** | **string** |  | [optional] [default to null]
**OrganizationalUnit** | **string** |  | [optional] [default to null]
**SmbAllowed** | **bool** |  | [optional] [default to true]
**NtlmEnabled** | **bool** |  | [optional] [default to true]
**UseAutoDiscovery** | **bool** |  | [optional] [default to false]
**UseLdaps** | **bool** |  | [optional] [default to false]
**Port** | **int** |  | [optional] [default to null]
**Binddn** | **string** |  | [optional] [default to null]
**Searchbase** | **string** |  | [optional] [default to null]
**DomainName** | **string** |  | [optional] [default to null]
**Method** | **string** |  | [optional] [default to METHOD.SIMPLE]
**QueryGroupsMode** | **string** |  | [optional] [default to QUERY_GROUPS_MODE.COMPATIBLE]
**PosixAttributesSource** | **string** |  | [optional] [default to JOINED_DOMAIN]
**UseTls** | **bool** |  | [optional] [default to false]
**ReverseLookup** | **bool** |  | [optional] [default to false]
**GidNumber** | **string** |  | [optional] [default to null]
**UseMultiForest** | **bool** |  | [optional] [default to false]
**Uid** | **string** |  | [optional] [default to null]
**UidNumber** | **string** |  | [optional] [default to null]
**MatchUser** | **string** |  | [optional] [default to null]
**UidMemberValuePropertyName** | **string** |  | [optional] [default to null]
**UidMember** | **string** |  | [optional] [default to null]
**PosixAccount** | **string** |  | [optional] [default to null]
**PosixGroup** | **string** |  | [optional] [default to null]
**UsernamePropertyName** | **string** |  | [optional] [default to null]
**UserLoginName** | **string** |  | [optional] [default to null]
**GroupLoginName** | **string** |  | [optional] [default to null]
**MailPropertyName** | **string** |  | [optional] [default to null]
**IsVmsAuthProvider** | **bool** |  | [optional] [default to false]
**Bindpw** | **string** |  | [optional] [default to null]
**Urls** | **[]string** | List of LDAP servers urls , starting with ldap:// | [optional] [default to null]
**SkipLdap** | **string** | When creating an active directory using this resource an ldap configuration is also created if set to true than whed deleting this resource the ldap configuration will be kept otherwise it will also be deleted | [optional] [default to false]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

