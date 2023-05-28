/*
 * VAST API Swagger Schema
 *
 * VAST Management API definition
 *
 * API version: 1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

// Share-level ACL details
type ViewShareAcl struct {
	// True if Share ACL is enabled on the view, otherwise False
	Enabled bool `json:"enabled,omitempty"`
	// Share-level ACL
	Acl []interface{} `json:"acl,omitempty"`
}