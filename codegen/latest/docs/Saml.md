# Saml

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**IdpName** | **string** | SAML IDP name | [optional] [default to null]
**EncryptAssertion** | **bool** | Set to true if the IdP encrypts the assertion. If true, an encryption certificate and key must be uploaded. Use encryption_saml_crt and encryption_saml_key to provide the required certificate and key. Default: false. Set to false to disable encryption. | [optional] [default to null]
**EncryptionSamlCrt** | **string** | Specifies the encryption certificate file content to upload. Required if encrypt_assertion is true. | [optional] [default to null]
**EncryptionSamlKey** | **string** | Specifies the encryption key file content to upload. Required if encrypt_assertion is true. | [optional] [default to null]
**ForceAuthn** | **bool** | Set to true to force authentication with the IDP even if there is an active session with the IdP for the user. Default: false. | [optional] [default to null]
**IdpEntityid** | **string** | A unique identifier for the IdP instance | [optional] [default to null]
**IdpMetadata** | **string** | Use local metadata. Supply local metadata XML. | [optional] [default to null]
**IdpMetadataUrl** | **string** | Use metadata located at specified remote URL. For example: &#x27;https://dev-12914105.okta.com/app/exke7ia133bKXWP2g5d7/sso/saml/metadata&#x27; | [optional] [default to null]
**SigningCert** | **string** | Specifies the certificate file content to use for requiring signed responses from the IdP. Required if want_assertions_or_response_signed is true. | [optional] [default to null]
**SigningKey** | **string** | Specifies the key file content to use for requiring signed responses from the IdP. Required if want_assertions_or_response_signed is true. | [optional] [default to null]
**WantAssertionsOrResponseSigned** | **bool** | Set to true to require a signed response or assertion from the IdP. VMS then fails the user authentication if an unsigned response is received. If true, upload a certificate and key. Use signing_cert and signing_key to provide certificate and key. Default: false. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

