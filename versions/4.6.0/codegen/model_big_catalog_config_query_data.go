/*
 * VAST API Swagger Schema
 *
 * VAST Management API definition
 *
 * API version: 1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type BigCatalogConfigQueryData struct {
	// Target source path
	Path string `json:"path"`
	// Defines the filters
	Filters *interface{} `json:"filters"`
	// Tenant ID
	TenantId int32 `json:"tenant_id,omitempty"`
	// The name of the snapshot, the latest version is taken by default
	Snapshot string `json:"snapshot,omitempty"`
	// Defines which fields should be displayed
	Fields []string `json:"fields,omitempty"`
}
