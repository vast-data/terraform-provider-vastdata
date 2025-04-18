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
	// The QoS type
	PolicyType string `json:"policy_type,omitempty"`
	// What attributes are setting the limitations.
	LimitBy string `json:"limit_by,omitempty"`
	// When setting is_default this is the tenant which will take affect
	TenantId int64 `json:"tenant_id,omitempty"`
	// List of local user IDs to which this QoS Policy is affective.
	AttachedUsersIdentifiers []string `json:"attached_users_identifiers,omitempty"`
	// Should this QoS Policy be the default QoS per user for this tenant ?, tnenat_id should be also provided when settingthis attribute
	IsDefault bool `json:"is_default,omitempty"`
	// Sets the size of IO for static and capacity limit definitions. The number of IOs per request is obtained by dividing request size by IO size. Default: 64K, Recommended range: 4K - 1M
	IoSizeBytes int64 `json:"io_size_bytes,omitempty"`
	StaticLimits *QosStaticLimits `json:"static_limits,omitempty"`
	CapacityLimits *QosDynamicLimits `json:"capacity_limits,omitempty"`
	StaticTotalLimits *QoSStaticTotalLimits `json:"static_total_limits,omitempty"`
	CapacityTotalLimits *QoSDynamicTotalLimits `json:"capacity_total_limits,omitempty"`
	AttachedUsers []QosUser `json:"attached_users,omitempty"`
}
