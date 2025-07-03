# NonLocalUserKey

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | The unique ID of the access key. | [optional] [default to null]
**Uid** | **int32** | The user&#x27;s Unix UID. | [optional] [default to null]
**AccessKey** | **string** | The user&#x27;s access key. | [optional] [default to null]
**SecretKey** | **string** | The user&#x27;s secret key. | [optional] [default to null]
**PgpPublicKey** | **string** | The PGP public key in the ASCII armor format to encrypt the secret key returned by the VAST cluster. If this option is set, the &#x27;encrypted_secret_key&#x27; will be returned, while &#x27;secret_key&#x27; will be empty. Changing this after apply will have no effect. | [optional] [default to null]
**EncryptedSecretKey** | **string** | The secret key returned from the VAST cluster. This key is encrypted with the public key that was supplied in &#x27;pgp_public_key&#x27;. | [optional] [default to null]
**TenantId** | **int32** | Tenant ID. | [optional] [default to null]
**Enabled** | **bool** | Sets the key to be enabled or disabled. | [optional] [default to true]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

