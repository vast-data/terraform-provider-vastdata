/*
 * VastData API
 *
 * A API document representing VastData API
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger

type GroupPost struct {
	// A uniq name given to the group
	Name string `json:"name,omitempty"`
	// The group linux gid
	Gid int32 `json:"gid,omitempty"`
	// The group SID
	Sid string `json:"sid,omitempty"`
	// List of S3 policies IDs
	S3PoliciesIds []int32 `json:"s3_policies_ids,omitempty"`
}
