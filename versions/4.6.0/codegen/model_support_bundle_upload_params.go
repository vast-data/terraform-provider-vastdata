/*
 * VAST API Swagger Schema
 *
 * VAST Management API definition
 *
 * API version: 1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type SupportBundleUploadParams struct {
	// Subdirectory in support bucket
	BucketSubdir string `json:"bucket_subdir,omitempty"`
	// S3 Bucket access key
	AccessKey string `json:"access_key,omitempty"`
	// S3 Bucket secret key
	SecretKey string `json:"secret_key,omitempty"`
	// S3 Bucket for upload
	BucketName string `json:"bucket_name,omitempty"`
	// If true, upload non-aggregated Support Bundle via VMS (requires proxy). Otherwise, upload from each node.
	UploadViaVms bool `json:"upload_via_vms,omitempty"`
}
