/*
 * VAST API Swagger Schema
 *
 * VAST Management API definition
 *
 * API version: 1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type Table struct {
	// Name of the Database
	DatabaseName string `json:"database_name,omitempty"`
	// Name of the Schema
	SchemaName string `json:"schema_name,omitempty"`
	// Name of the Table
	Name string `json:"name,omitempty"`
}
