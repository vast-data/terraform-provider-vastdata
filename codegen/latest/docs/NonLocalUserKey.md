# NonLocalUserKey

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | The Access key unique identifier | [optional] [default to null]
**Uid** | **int32** | The user unix UID | [optional] [default to null]
**AccessKey** | **string** | The access id of the user key | [optional] [default to null]
**SecretKey** | **string** | The secret id of the user key | [optional] [default to null]
**PgpPublicKey** | **string** | The PGP public key at ascii armor format to encrypt the secret id returned from vast cluster, if this option is set than the encrypted_secret_key will be returned and secret_key will be empty, changing it after apply will have no affect | [optional] [default to null]
**EncryptedSecretKey** | **string** | The secret id returned from the vast cluster encrypted with the public key provided at pgp_public_key | [optional] [default to null]
**TenantId** | **int32** | Tenant ID | [optional] [default to null]
**Enabled** | **bool** | Should the key be enabled or disabled | [optional] [default to true]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

