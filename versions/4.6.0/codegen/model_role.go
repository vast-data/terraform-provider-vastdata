/*
 * VAST API Swagger Schema
 *
 * VAST Management API definition
 *
 * API version: 1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type Role struct {
	// 
	Id int32 `json:"id,omitempty"`
	// The friendly name of the role
	Name string `json:"name,omitempty"`
	// 
	Managers *interface{} `json:"managers,omitempty"`
	// 
	Permissions *interface{} `json:"permissions,omitempty"`
	// Is the role is a default role
	IsDefault bool `json:"is_default,omitempty"`
	// Ldap groups which will be granted this role in login
	LdapGroups []string `json:"ldap_groups,omitempty"`
	// Tenants for that role
	Tenants []string `json:"tenants,omitempty"`
}