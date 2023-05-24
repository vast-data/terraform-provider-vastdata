/*
 * VAST API Swagger Schema
 *
 * VAST Management API definition
 *
 * API version: 1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type DboxesModifyParams struct {
	// DBox description
	Description string `json:"description,omitempty"`
	// True to start replacement
	Replace bool `json:"replace,omitempty"`
	// True to conclude replacement
	Conclude bool `json:"conclude,omitempty"`
	// In case of concluding: do not verify that all devices have been moved. In case of replacing: support moving device while getting to degraded state during replacement
	Force bool `json:"force,omitempty"`
}
