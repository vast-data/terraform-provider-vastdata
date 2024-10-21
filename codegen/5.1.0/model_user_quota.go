/*
 * VastData API
 *
 * A API document representing VastData API
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger

type UserQuota struct {
	// Quota enforcement grace period at the format of DD HH:MM:SS or HH:MM:SS
	GracePeriod string `json:"grace_period,omitempty"`
	// Grace period expiration time
	TimeToBlock string `json:"time_to_block,omitempty"`
	// Soft quota limit
	SoftLimit int64 `json:"soft_limit,omitempty"`
	// Hard quota limit
	HardLimit int64 `json:"hard_limit,omitempty"`
	// Hard inodes quota limit
	HardLimitInodes int64 `json:"hard_limit_inodes,omitempty"`
	// Soft inodes quota limit
	SoftLimitInodes int64 `json:"soft_limit_inodes,omitempty"`
	// Used inodes
	UsedInodes int64 `json:"used_inodes,omitempty"`
	// Used capacity in bytes
	UsedCapacity int64 `json:"used_capacity,omitempty"`
	IsAccountable bool `json:"is_accountable,omitempty"`
	QuotaSystemId int32 `json:"quota_system_id,omitempty"`
	Entity *QuotaEntityInfo `json:"entity,omitempty"`
}
