/*
 * VAST Data API
 *
 * A API document representing VAST Data API
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger

type NonLocalUserKey struct {
	// The unique ID of the access key.
	Id string `json:"id,omitempty"`
	// The user's Unix UID.
	Uid int32 `json:"uid,omitempty"`
	// The user's access key.
	AccessKey string `json:"access_key,omitempty"`
	// The user's secret key.
	SecretKey string `json:"secret_key,omitempty"`
	// The PGP public key in the ASCII armor format to encrypt the secret key returned by the VAST cluster. If this option is set, the 'encrypted_secret_key' will be returned, while 'secret_key' will be empty. Changing this after apply will have no effect.
	PgpPublicKey string `json:"pgp_public_key,omitempty"`
	// The secret key returned from the VAST cluster. This key is encrypted with the public key that was supplied in 'pgp_public_key'.
	EncryptedSecretKey string `json:"encrypted_secret_key,omitempty"`
	// Tenant ID.
	TenantId int32 `json:"tenant_id,omitempty"`
	// Sets the key to be enabled or disabled.
	Enabled bool `json:"enabled,omitempty"`
}
