/*
 * VAST API Swagger Schema
 *
 * VAST Management API definition
 *
 * API version: 1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type GlobalPath struct {
	// 
	Id int32 `json:"id,omitempty"`
	// unique identifier
	Guid string `json:"guid,omitempty"`
	// 
	Name string `json:"name,omitempty"`
	// Global Path
	Path string `json:"path,omitempty"`
	// Exposed Path Policy
	Policy string `json:"policy,omitempty"`
	// Global Peer
	Peer string `json:"peer,omitempty"`
	// Exposed Path
	ExposedPath string `json:"exposed_path,omitempty"`
	// Enabed/Disabled
	Enabled bool `json:"enabled,omitempty"`
}