/*
 * VastData API
 *
 * A API document representing VastData API
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger

type RealmPermission struct {
	// The name of the realm
	RealmName string `json:"realm_name,omitempty"`
	// Should allow create related to permissions associated with this realm
	Create bool `json:"create,omitempty"`
	// Should allow view related to permissions associated with this realm
	View bool `json:"view,omitempty"`
	// Should allow delete related to permissions associated with this realm
	Delete bool `json:"delete,omitempty"`
	// Should allow edit related to permissions associated with this realm
	Edit bool `json:"edit,omitempty"`
}
