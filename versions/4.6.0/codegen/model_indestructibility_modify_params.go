/*
 * VAST API Swagger Schema
 *
 * VAST Management API definition
 *
 * API version: 1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type IndestructibilityModifyParams struct {
	// 
	IndestructibilityPasswd string `json:"indestructibility_passwd,omitempty"`
	// 
	NewIndestructibilityPasswd string `json:"new_indestructibility_passwd,omitempty"`
	// 
	PasswdDelay string `json:"passwd_delay,omitempty"`
}
