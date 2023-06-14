/*
 * VAST API Swagger Schema
 *
 * VAST Management API definition
 *
 * API version: 1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type User struct {
	//
	Id int32 `json:"id,omitempty"`
	//
	Guid string `json:"guid,omitempty"`
	// The name of the user
	Name string `json:"name,omitempty"`
	// UID
	Uid int32 `json:"uid,omitempty"`
	// Leading GID
	LeadingGid int32 `json:"leading_gid,omitempty"`
	// GID list
	Gids []int32 `json:"gids,omitempty"`
	// Group list
	Groups []string `json:"groups,omitempty"`
	// Group Count
	GroupCount int32 `json:"group_count,omitempty"`
	// Leading Group
	LeadingGroupName string `json:"leading_group_name,omitempty"`
	// Leading Group GID
	LeadingGroupGid int32 `json:"leading_group_gid,omitempty"`
	// SID
	Sid string `json:"sid,omitempty"`
	// Primary group SID
	PrimaryGroupSid string `json:"primary_group_sid,omitempty"`
	// SID list
	Sids []string `json:"sids,omitempty"`
	//
	Local bool `json:"local,omitempty"`
	// Access Keys
	AccessKeys [][]string `json:"access_keys,omitempty"`
	// Allow create bucket
	AllowCreateBucket bool `json:"allow_create_bucket,omitempty"`
	// Allow delete bucket
	AllowDeleteBucket bool `json:"allow_delete_bucket,omitempty"`
	// S3 superuser
	S3Superuser bool `json:"s3_superuser,omitempty"`
	// S3 policies IDs
	S3PoliciesIds []int32 `json:"s3_policies_ids,omitempty"`
}
