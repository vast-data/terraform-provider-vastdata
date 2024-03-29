/*
 * VAST API Swagger Schema
 *
 * VAST Management API definition
 *
 * API version: 1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type ObjectPermission struct {
	// 
	RoleId string `json:"role_id,omitempty"`
	// 
	UserId string `json:"user_id,omitempty"`
	// 
	Realm string `json:"realm,omitempty"`
	// 
	Permissions string `json:"permissions,omitempty"`
	// 
	ObjectType string `json:"object_type,omitempty"`
	// 
	ObjectId string `json:"object_id,omitempty"`
}
