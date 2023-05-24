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

type ReplicationPolicyCreateParams struct {
	// 
	Name string `json:"name,omitempty"`
	// schedule frequency, in datetime format
	ScheduleFrequency time.Time `json:"schedule_frequency,omitempty"`
	// Schedule the first restore point after the initial sync
	ScheduleStartTime time.Time `json:"schedule_start_time,omitempty"`
	// replication target id
	ReplicationTargetId int32 `json:"replication_target_id,omitempty"`
	// bandwith limitation rules
	BandwidthLimitationRules string `json:"bandwidth_limitation_rules,omitempty"`
	// low / normal / high
	Priority string `json:"priority,omitempty"`
	// Amazon S3 / Amazon Glacier
	AwsPreferredStorage string `json:"aws_preferred_storage,omitempty"`
}
