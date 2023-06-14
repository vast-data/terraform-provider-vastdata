/*
 * VAST API Swagger Schema
 *
 * VAST Management API definition
 *
 * API version: 1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type ExposedPathCreateParams struct {
	// Do not specify this parameter.
	Id int32 `json:"id,omitempty"`
	// Do not specify this parameter.
	Guid string `json:"guid,omitempty"`
	// Exposed Path Name
	Name string `json:"name"`
	// Exposed Path
	Path string `json:"path,omitempty"`
	// Exposed Path Policy
	Policy string `json:"policy,omitempty"`
	// Exposed Path Policy ID
	PolicyId int32 `json:"policy_id,omitempty"`
}
