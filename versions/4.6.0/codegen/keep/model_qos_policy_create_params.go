/*
 * VAST API Swagger Schema
 *
 * VAST Management API definition
 *
 * API version: 1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type QosPolicyCreateParams struct {
	Name string `json:"name"`
	// Allocation of performance resources mode
	Mode string `json:"mode,omitempty"`
	// Size of a single IO, default is 64K
	IoSizeBytes int32 `json:"io_size_bytes,omitempty"`
	// Static mode limits
	StaticLimits *RequestQosStaticLimits `json:"static_limits,omitempty"`
	// Capacity mode limits
	CapacityLimits *RequestQosDynamicLimits `json:"capacity_limits,omitempty"`
}
