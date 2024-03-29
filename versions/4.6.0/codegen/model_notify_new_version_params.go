/*
 * VAST API Swagger Schema
 *
 * VAST Management API definition
 *
 * API version: 1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type NotifyNewVersionParams struct {
	// link to download new version
	S3Url string `json:"s3_url,omitempty"`
	// build
	Build string `json:"build,omitempty"`
	// sys_version
	SysVersion string `json:"sys_version,omitempty"`
	// os_version
	OsVersion string `json:"os_version,omitempty"`
	// ssd_version
	SsdVersion string `json:"ssd_version,omitempty"`
	// nvram_version
	NvramVersion string `json:"nvram_version,omitempty"`
	// bmc_fw_version
	BmcFwVersion string `json:"bmc_fw_version,omitempty"`
}
