/*
 * VAST API Swagger Schema
 *
 * VAST Management API definition
 *
 * API version: 1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type TableAddColumnsParams struct {
	// Name of the Database
	DatabaseName string `json:"database_name"`
	// Name of the Schema
	SchemaName string `json:"schema_name"`
	// Name of the Table
	Name string `json:"name"`
	// Description of table columns
	ArrowSchema []interface{} `json:"arrow_schema"`
}
