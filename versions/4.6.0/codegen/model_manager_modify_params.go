/*
 * VAST API Swagger Schema
 *
 * VAST Management API definition
 *
 * API version: 1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type ManagerModifyParams struct {
	// 
	Username string `json:"username,omitempty"`
	// 
	Password string `json:"password,omitempty"`
	// 
	FirstName string `json:"first_name,omitempty"`
	// 
	LastName string `json:"last_name,omitempty"`
	// Specify all role IDs to update the roles assigned to the manager. Separate role IDs with commas.
	Roles string `json:"roles,omitempty"`
	// Realm type
	Realm string `json:"realm,omitempty"`
	// permission type
	Permissions string `json:"permissions,omitempty"`
	// Specify permissions list as an array of permission codenames in the format <permission>-<realm>. To list permission codenames, run /permissions/get.
	PermissionsList string `json:"permissions_list,omitempty"`
	// object type
	ObjectType string `json:"object_type,omitempty"`
	// object ID
	ObjectId string `json:"object_id,omitempty"`
}
