/*
 * VAST API Swagger Schema
 *
 * VAST Management API definition
 *
 * API version: 1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type QuotaEntityInfo struct {
	// 
	Id int32 `json:"id,omitempty"`
	// The name
	Name string `json:"name,omitempty"`
	// 
	VastId int32 `json:"vast_id,omitempty"`
	Email string `json:"email,omitempty"`
	IsGroup bool `json:"is_group,omitempty"`
	Identifier string `json:"identifier,omitempty"`
	IdentifierType string `json:"identifier_type,omitempty"`
}
