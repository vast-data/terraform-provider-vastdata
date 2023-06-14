/*
 * VAST API Swagger Schema
 *
 * VAST Management API definition
 *
 * API version: 1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type ClusterRpc struct {
	// rpc name to execute
	Rpc string `json:"rpc,omitempty"`
	// params for rpc call
	Params string `json:"params,omitempty"`
	// Module type for the commander connection
	ModuleType string `json:"module_type,omitempty"`
}
