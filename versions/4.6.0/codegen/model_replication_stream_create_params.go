/*
 * VAST API Swagger Schema
 *
 * VAST Management API definition
 *
 * API version: 1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type ReplicationStreamCreateParams struct {
	// 
	Name string `json:"name,omitempty"`
	// path to replicate
	SourceDir string `json:"source_dir,omitempty"`
	// replication policy id
	PolicyId string `json:"policy_id,omitempty"`
	// enable/pause replication stream
	Enabled bool `json:"enabled,omitempty"`
	// protection policy id
	ProtectionPolicyId string `json:"protection_policy_id,omitempty"`
}