/*
 * VastData API
 *
 * A API document representing VastData API
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger

type UserKey struct {
	// The Access key unique identifier
	Id string `json:"id,omitempty"`
	// The user id to create the Key for
	UserId int32 `json:"user_id,omitempty"`
	// The access id of the user key
	AccessKey string `json:"access_key,omitempty"`
	// The secret id of the user key, please note that that the secret id is not encrypted and should be kept in an highly secure backend ,this field will only be returned if pgp_public_key is not provided
	SecretKey string `json:"secret_key,omitempty"`
	// The PGP public key at ascii armor format to encrypt the secret id returned from vast cluster, if this option is set than the encrypted_secret_key will be returned and secret_key will be empty, changing it after apply will have no affect
	PgpPublicKey string `json:"pgp_public_key,omitempty"`
	// The secret id returned from the vast cluster encrypted with the public key provided at pgp_public_key
	EncryptedSecretKey string `json:"encrypted_secret_key,omitempty"`
	// Should the key be enabled or disabled
	Enabled bool `json:"enabled,omitempty"`
}
