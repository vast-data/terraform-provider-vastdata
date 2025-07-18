/*
 * VAST Data API
 *
 * A API document representing VAST Data API
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger

type ProtectedPath struct {
	Id   int32  `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	// GUID
	Guid string `json:"guid,omitempty"`
	// Protection policy ID.
	ProtectionPolicyId int32 `json:"protection_policy_id,omitempty"`
	// Path to replicate.
	SourceDir string `json:"source_dir,omitempty"`
	// The destination path to replicate to on the remote cluster.
	TargetExportedDir string `json:"target_exported_dir,omitempty"`
	// Local tenant ID.
	TenantId int32 `json:"tenant_id,omitempty"`
	// Remote tenant GUID.
	RemoteTenantGuid string `json:"remote_tenant_guid,omitempty"`
	// Remote target object ID.
	TargetId int32 `json:"target_id,omitempty"`
	// Available replication capabilities. Supported only for clusters v5.1 and later.
	Capabilities string `json:"capabilities,omitempty"`
	// Specifies whether the protected path is enabled.
	Enabled bool `json:"enabled,omitempty"`
}
