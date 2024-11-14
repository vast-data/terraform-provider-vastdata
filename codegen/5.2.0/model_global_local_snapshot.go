/*
 * VastData API
 *
 * A API document representing VastData API
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger

type GlobalLocalSnapshot struct {
	// A unique id given to the global snapshot
	Id int64 `json:"id,omitempty"`
	// A unique guid given to the global snapshot
	Guid string `json:"guid,omitempty"`
	// The name of the snapshot
	Name string `json:"name,omitempty"`
	// The tenant ID of the target
	LoaneeTenantId int `json:"loanee_tenant_id,omitempty"`
	// The path where to store the snapshot on a Target
	LoaneeRootPath string `json:"loanee_root_path,omitempty"`
	// The id of the local snapshot
	LoaneeSnapshotId int `json:"loanee_snapshot_id,omitempty"`
	// Is the snapshot enabled
	Enabled bool `json:"enabled,omitempty"`
	OwnerTenant *GlobalSnapshotOwnerTenant `json:"owner_tenant,omitempty"`
}