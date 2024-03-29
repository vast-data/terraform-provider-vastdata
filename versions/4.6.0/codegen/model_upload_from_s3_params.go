/*
 * VAST API Swagger Schema
 *
 * VAST Management API definition
 *
 * API version: 1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type UploadFromS3Params struct {
	// S3 URL to upgrade package. If not provided, will be taken from db
	S3Url string `json:"s3_url,omitempty"`
	// Skips preparing the cluster for upgrade, including: pre-upgrade validations, copying the bundle to other hosts, and pulling the image on all CNodes.
	SkipPrepare bool `json:"skip_prepare,omitempty"`
	// Skips validation of hardware component health. Use with caution since component redundancy is important in NDU. Do not use with OS upgrade.
	SkipHwCheck bool `json:"skip_hw_check,omitempty"`
}
