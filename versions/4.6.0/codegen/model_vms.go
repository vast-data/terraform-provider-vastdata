/*
 * VAST API Swagger Schema
 *
 * VAST Management API definition
 *
 * API version: 1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

import (
	"time"
)

type Vms struct {
	// 
	Id int32 `json:"id,omitempty"`
	// Data Bond IP
	Ip string `json:"ip,omitempty"`
	// Data IP1
	Ip1 string `json:"ip1,omitempty"`
	// Data IP2
	Ip2 string `json:"ip2,omitempty"`
	// Management IP
	MgmtIp string `json:"mgmt_ip,omitempty"`
	// Management VIP
	MgmtVip string `json:"mgmt_vip,omitempty"`
	// Management IPv6 VIP
	MgmtVipIpv6 string `json:"mgmt_vip_ipv6,omitempty"`
	// Management CNode
	MgmtCnode string `json:"mgmt_cnode,omitempty"`
	// VMS inner (NFS) vip
	MgmtInnerVip string `json:"mgmt_inner_vip,omitempty"`
	// VMS inner VIP CNode
	MgmtInnerVipCnode string `json:"mgmt_inner_vip_cnode,omitempty"`
	// 
	Name string `json:"name"`
	// 
	Build string `json:"build,omitempty"`
	// 
	SwVersion string `json:"sw_version,omitempty"`
	// 
	AutoLogoutTimeout string `json:"auto_logout_timeout,omitempty"`
	// 
	Created time.Time `json:"created,omitempty"`
	// 
	State string `json:"state,omitempty"`
	// 
	Url string `json:"url,omitempty"`
	// Is Management HA disabled?
	DisableMgmtHa bool `json:"disable_mgmt_ha,omitempty"`
	// Disable vms metrics collection
	DisableVmsMetrics bool `json:"disable_vms_metrics,omitempty"`
	// Format capacity properties to base 10
	CapacityBase10 bool `json:"capacity_base_10,omitempty"`
	// Format performance properties to base 10
	PerformanceBase10 bool `json:"performance_base_10,omitempty"`
	// Minimum supported TLS version (e.g.: 1.2)
	MinTlsVersion string `json:"min_tls_version,omitempty"`
	// The reason for VMS degraded state
	DegradedReason string `json:"degraded_reason,omitempty"`
	// Validity duration for JWT access token, specify as [DD [HH:[MM:]]]ss
	AccessTokenLifetime string `json:"access_token_lifetime,omitempty"`
	// Validity duration for JWT refresh token, specify as [DD [HH:[MM:]]]ss
	RefreshTokenLifetime string `json:"refresh_token_lifetime,omitempty"`
	// Customize login banner for VMS Web UI and CLI
	LoginBanner string `json:"login_banner,omitempty"`
	// the remaining capacity out of the active capacity
	TotalUsageCapacityPercentage string `json:"total_usage_capacity_percentage,omitempty"`
	// sum of all active licenses capacity
	TotalActiveCapacity string `json:"total_active_capacity,omitempty"`
	// the total active capacity minus the system capacity
	TotalRemainingCapacity string `json:"total_remaining_capacity,omitempty"`
	// Parameter that controls everything related to database
	TabularSupport string `json:"tabular_support,omitempty"`
	// Parameter that controls visibility of ipv6 fields for VIP Pools
	Ipv6Support bool `json:"ipv6_support,omitempty"`
}
