# UserKey

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | The unique ID of the access key. | [optional] [default to null]
**UserId** | **int32** | The ID of the user to create the key for. | [optional] [default to null]
**AccessKey** | **string** | The access key of the user. | [optional] [default to null]
**SecretKey** | **string** | The secret key of the user. This secret key is not encrypted and should be kept in an highly secure backend. This field will only be returned if &#x27;pgp_public_key&#x27; is not provided. | [optional] [default to null]
**PgpPublicKey** | **string** | The PGP public key in the ASCII armor format to encrypt the secret key returned by the VAST cluster. If this option is set, the &#x27;encrypted_secret_key&#x27; will be returned, while &#x27;secret_key&#x27; will be empty. Changing this after apply will have no effect. | [optional] [default to null]
**EncryptedSecretKey** | **string** | The secret key returned from the VAST cluster. This key is encrypted with the public key that was supplied in &#x27;pgp_public_key&#x27;. | [optional] [default to null]
**Enabled** | **bool** | Sets the key to be enabled or disabled. | [optional] [default to true]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

