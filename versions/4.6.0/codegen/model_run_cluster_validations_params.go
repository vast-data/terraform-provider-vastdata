/*
 * VAST API Swagger Schema
 *
 * VAST Management API definition
 *
 * API version: 1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type RunClusterValidationsParams struct {
	// Skip hardware related checks
	SkipHwCheck bool `json:"skip_hw_check,omitempty"`
	// skip os upgrade related checks
	SkipOsUpgrade bool `json:"skip_os_upgrade,omitempty"`
	// force
	Force bool `json:"force,omitempty"`
}
