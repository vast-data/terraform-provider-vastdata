/*
 * VAST API Swagger Schema
 *
 * VAST Management API definition
 *
 * API version: 1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type Snapshot struct {
	// 
	Id int32 `json:"id,omitempty"`
	// 
	Guid string `json:"guid,omitempty"`
	// 
	Name string `json:"name"`
	// Snapshot path
	Path string `json:"path,omitempty"`
	// Snapshot expiration time UTC
	ExpirationTime string `json:"expiration_time,omitempty"`
	// Snapshot stats
	State string `json:"state,omitempty"`
	// Associated snapshot policy
	Policy string `json:"policy,omitempty"`
	// Associated snapshot policy ID
	PolicyId int32 `json:"policy_id,omitempty"`
	// Parent Cluster
	Cluster string `json:"cluster,omitempty"`
	// Parent Cluster ID
	ClusterId int32 `json:"cluster_id,omitempty"`
	// Parent handle
	Handle string `json:"handle,omitempty"`
	// Snapshot created time
	Created string `json:"created,omitempty"`
	// Lock the snapshot from being deleted by cleanup
	Locked bool `json:"locked,omitempty"`
	// 
	CloneId int32 `json:"clone_id,omitempty"`
	// The usable capacity reclaimable by deleting the snapshot and all older snapshots on the protected path
	AggrPhysEstimation int64 `json:"aggr_phys_estimation,omitempty"`
	// The usable capacity reclaimable by deleting the snapshot without deleting other snapshots on the path
	UniquePhysEstimation int64 `json:"unique_phys_estimation,omitempty"`
	// 
	ProtectionPolicyId int32 `json:"protection_policy_id,omitempty"`
	// Protection Policy Name
	ProtectionPolicy string `json:"protection_policy,omitempty"`
	Type_ string `json:"type,omitempty"`
	// Prevent the snapshot from being deleted
	Indestructible bool `json:"indestructible,omitempty"`
	// Tenant ID
	TenantId int32 `json:"tenant_id,omitempty"`
}