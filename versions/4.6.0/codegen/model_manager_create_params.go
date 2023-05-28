/*
 * VAST API Swagger Schema
 *
 * VAST Management API definition
 *
 * API version: 1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type ManagerCreateParams struct {
	// Manager's user name
	Username string `json:"username"`
	// Manager's password
	Password string `json:"password"`
	// Manager's first name
	FirstName string `json:"first_name,omitempty"`
	// Manager's last name
	LastName string `json:"last_name,omitempty"`
	// Joins manager to specified roles. Specify as an array of role IDs, separated by commas.
	Roles string `json:"roles,omitempty"`
}