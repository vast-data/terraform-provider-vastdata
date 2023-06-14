/*
 * VAST API Swagger Schema
 *
 * VAST Management API definition
 *
 * API version: 1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type BigCatalogIndexedColumnRemove struct {
	// Name of the Big catalog indexed column
	Name string `json:"name"`
	// Type of indexed column
	ColumnType string `json:"column_type"`
}
