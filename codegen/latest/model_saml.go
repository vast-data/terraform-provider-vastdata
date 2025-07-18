/*
 * VastData API
 *
 * A API document representing VastData API
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger

type Saml struct {
	// Internal (Terraform) ID of SAML
	Id string `json:"id,omitempty"`
	// VMS ID
	VmsId int32 `json:"vms_id,omitempty"`
	// SAML IDP name
	IdpName string `json:"idp_name,omitempty"`
	// Set to true if the IdP encrypts the assertion. If true, an encryption certificate and key must be uploaded. Use encryption_saml_crt and encryption_saml_key to provide the required certificate and key. Default: false. Set to false to disable encryption.
	EncryptAssertion bool `json:"encrypt_assertion,omitempty"`
	// Specifies the encryption certificate file content to upload. Required if encrypt_assertion is true.
	EncryptionSamlCrt string `json:"encryption_saml_crt,omitempty"`
	// Specifies the encryption key file content to upload. Required if encrypt_assertion is true.
	EncryptionSamlKey string `json:"encryption_saml_key,omitempty"`
	// Set to true to force authentication with the IDP even if there is an active session with the IdP for the user. Default: false.
	ForceAuthn bool `json:"force_authn,omitempty"`
	// A unique identifier for the IdP instance
	IdpEntityid string `json:"idp_entityid,omitempty"`
	// Use local metadata. Supply local metadata XML.
	IdpMetadata string `json:"idp_metadata,omitempty"`
	// Use metadata located at specified remote URL. For example: 'https://dev-12914105.okta.com/app/exke7ia133bKXWP2g5d7/sso/saml/metadata'
	IdpMetadataUrl string `json:"idp_metadata_url,omitempty"`
	// Specifies the certificate file content to use for requiring signed responses from the IdP. Required if want_assertions_or_response_signed is true.
	SigningCert string `json:"signing_cert,omitempty"`
	// Specifies the key file content to use for requiring signed responses from the IdP. Required if want_assertions_or_response_signed is true.
	SigningKey string `json:"signing_key,omitempty"`
	// Set to true to require a signed response or assertion from the IdP. VMS then fails the user authentication if an unsigned response is received. If true, upload a certificate and key. Use signing_cert and signing_key to provide certificate and key. Default: false.
	WantAssertionsOrResponseSigned bool `json:"want_assertions_or_response_signed,omitempty"`
}
