/*
 * VAST API Swagger Schema
 *
 * VAST Management API definition
 *
 * API version: 1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type S3LifeCycleRuleCreateParams struct {
	// name
	Name string `json:"name,omitempty"`
	// View ID
	ViewId string `json:"view_id,omitempty"`
	Enabled bool `json:"enabled,omitempty"`
	Prefix string `json:"prefix,omitempty"`
	MinSize int32 `json:"min_size,omitempty"`
	MaxSize int32 `json:"max_size,omitempty"`
	ExpirationDays int32 `json:"expiration_days,omitempty"`
	ExpirationDate string `json:"expiration_date,omitempty"`
	ExpiredObjDeleteMarker bool `json:"expired_obj_delete_marker,omitempty"`
	NoncurrentDays int32 `json:"noncurrent_days,omitempty"`
	NewerNoncurrentVersions int32 `json:"newer_noncurrent_versions,omitempty"`
	AbortMpuDaysAfterInitiation int32 `json:"abort_mpu_days_after_initiation,omitempty"`
}
