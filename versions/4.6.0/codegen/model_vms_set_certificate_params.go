/*
 * VAST API Swagger Schema
 *
 * VAST Management API definition
 *
 * API version: 1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type VmsSetCertificateParams struct {
	// SSL Certificate file content, including the BEGIN CERTIFICATE and END CERTIFICATE lines
	SslCertificate string `json:"ssl_certificate,omitempty"`
	// SSL private key file content, include the BEGIN PRIVATE KEY and END PRIVATE KEY lines
	SslKeyfile string `json:"ssl_keyfile,omitempty"`
}
