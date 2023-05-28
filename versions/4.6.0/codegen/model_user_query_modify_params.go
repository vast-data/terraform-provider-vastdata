/*
 * VAST API Swagger Schema
 *
 * VAST Management API definition
 *
 * API version: 1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type UserQueryModifyParams struct {
	// Set to true to give the user permission to create S3 buckets
	AllowCreateBucket bool `json:"allow_create_bucket,omitempty"`
	// Set to true to give the user permission to delete S3 buckets
	AllowDeleteBucket bool `json:"allow_delete_bucket,omitempty"`
	// Set to true for S3 superuser
	S3Superuser bool `json:"s3_superuser,omitempty"`
	// User name
	Username string `json:"username,omitempty"`
	// User UID
	Uid int32 `json:"uid,omitempty"`
	// User SID
	Sid string `json:"sid,omitempty"`
	// list of s3 policy ids
	S3PoliciesIds *interface{} `json:"s3_policies_ids,omitempty"`
	// Tenant ID
	TenantId int32 `json:"tenant_id,omitempty"`
}