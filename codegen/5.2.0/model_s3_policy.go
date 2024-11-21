/*
 * VastData API
 *
 * A API document representing VastData API
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger

type S3Policy struct {
	Id int32 `json:"id,omitempty"`
	// GUID
	Guid string `json:"guid,omitempty"`
	Name string `json:"name,omitempty"`
	Policy string `json:"policy,omitempty"`
	// List of group names associated with this policy
	Users []string `json:"users,omitempty"`
	// List of group names associated with this policy
	Groups []string `json:"groups,omitempty"`
	IsReplicated bool `json:"is_replicated,omitempty"`
	Enabled bool `json:"enabled"`
	TenantId int64 `json:"tenant_id,omitempty"`
}
