/*
 * VastData API
 *
 * A API document representing VastData API
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger

type User struct {
	// A uniq id given to user
	Id int32 `json:"id,omitempty"`
	// A uniq guid given to the user
	Guid string `json:"guid,omitempty"`
	// A uniq name given to the user
	Name string `json:"name,omitempty"`
	// The user unix UID
	Uid int32 `json:"uid,omitempty"`
	// The user leading unix GID
	LeadingGid int32 `json:"leading_gid,omitempty"`
	// List of supplementary GID list
	Gids []int32 `json:"gids,omitempty"`
	// List of supplementary Group list
	Groups []string `json:"groups,omitempty"`
	// Group Count
	GroupCount int32 `json:"group_count,omitempty"`
	// Leading Group Name
	LeadingGroupName string `json:"leading_group_name,omitempty"`
	// Leading Group GID
	LeadingGroupGid int `json:"leading_group_gid,omitempty"`
	// The user SID
	Sid string `json:"sid,omitempty"`
	// The user primary group SID
	PrimaryGroupSid string `json:"primary_group_sid,omitempty"`
	// supplementary SID list
	Sids []string `json:"sids,omitempty"`
	// IS this a local user
	Local bool `json:"local,omitempty"`
	// Allow create bucket
	AllowCreateBucket bool `json:"allow_create_bucket,omitempty"`
	// Allow delete bucket
	AllowDeleteBucket bool `json:"allow_delete_bucket,omitempty"`
	// Is S3 superuser
	S3Superuser bool `json:"s3_superuser,omitempty"`
	// List S3 policies IDs
	S3PoliciesIds []int32 `json:"s3_policies_ids,omitempty"`
}
