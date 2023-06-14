/*
 * VAST API Swagger Schema
 *
 * VAST Management API definition
 *
 * API version: 1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type ActiveDirectoryDeleteParams struct {
	// Set to true to avoid removing the LDAP server configuration that is attached to the Active Directory configuration.
	SkipLdap bool `json:"skip_ldap,omitempty"`
}
