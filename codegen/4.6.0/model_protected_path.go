/*
 * VastData API
 *
 * A API document representing VastData API
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger

type ProtectedPath struct {
	Id int32 `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	// guid
	Guid string `json:"guid,omitempty"`
	// protection policy id
	ProtectionPolicyId string `json:"protection_policy_id,omitempty"`
	// protection policy name
	ProtectionPolicyName string `json:"protection_policy_name,omitempty"`
	// path to replicate
	SourceDir string `json:"source_dir,omitempty"`
	// where to replicate on the remote
	TargetExportedDir string `json:"target_exported_dir,omitempty"`
	// Local Tenant ID
	TenantId int32 `json:"tenant_id,omitempty"`
	// remote tenant name
	RemoteTenantName string `json:"remote_tenant_name,omitempty"`
}
