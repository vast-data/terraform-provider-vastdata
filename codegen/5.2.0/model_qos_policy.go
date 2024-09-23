/*
 * VastData API
 *
 * A API document representing VastData API
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger

type QosPolicy struct {
	Id int32 `json:"id,omitempty"`
	// QoS Policy guid
	Guid string `json:"guid,omitempty"`
	Name string `json:"name,omitempty"`
	// QoS provisioning mode
	Mode string `json:"mode,omitempty"`
	// Sets the size of IO for static and capacity limit definitions. The number of IOs per request is obtained by dividing request size by IO size. Default: 64K, Recommended range: 4K - 1M
	IoSizeBytes int64 `json:"io_size_bytes,omitempty"`
	StaticLimits *QosStaticLimits `json:"static_limits,omitempty"`
	CapacityLimits *QosDynamicLimits `json:"capacity_limits,omitempty"`
}
