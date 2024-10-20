# Manager

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **int64** | A unique id given to the manager | [optional] [default to null]
**Username** | **string** | The username of the manager | [optional] [default to null]
**Password** | **string** | The username of the manager | [optional] [default to null]
**FirstName** | **string** | The user firstname | [optional] [default to null]
**LastName** | **string** | The user last name | [optional] [default to null]
**PermissionsList** | **[]string** | List of allowed permissions | [optional] [default to null]
**Roles** | **[]int** | List of roles ids | [optional] [default to null]
**PasswordExpirationDisabled** | **bool** | Disable apssword expiration | [optional] [default to true]
**IsTemporaryPassword** | **bool** | Disable apssword expiration | [optional] [default to true]
**PasswordHash** | **string** | password sha256 to be used to check for password updates | [optional] [default to null]
**RealmsPermissions** | [**[]RealmPermission**](RealmPermission.md) | List of realms related permissions | [optional] [default to null]
**Permissions** | **[]string** | List of allowed permissions returned from the VMS | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

