/*
 * VAST API Swagger Schema
 *
 * VAST Management API definition
 *
 * API version: 1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type ModifyDeltaConfigParams struct {
	// enable/disable delta
	Enabled bool `json:"enabled,omitempty"`
	// set prefix to ignore objects in delta
	IgnoredObjectsPrefix bool `json:"ignored_objects_prefix,omitempty"`
}