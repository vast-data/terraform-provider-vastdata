/*
 * VAST API Swagger Schema
 *
 * VAST Management API definition
 *
 * API version: 1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type SchemaCreateParams struct {
	// Name of the Database
	DatabaseName string `json:"database_name"`
	// Name of the Schema
	Name string `json:"name"`
}
