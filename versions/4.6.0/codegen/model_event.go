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

type Event struct {
	// 
	Id int32 `json:"id,omitempty"`
	// 
	Timestamp time.Time `json:"timestamp,omitempty"`
	// 
	Guid string `json:"guid,omitempty"`
	// 
	Cluster string `json:"cluster,omitempty"`
	// 
	ObjectType string `json:"object_type,omitempty"`
	// 
	ObjectId string `json:"object_id,omitempty"`
	// 
	ObjectGuid string `json:"object_guid,omitempty"`
	// 
	EventType string `json:"event_type,omitempty"`
	// 
	EventOrigin string `json:"event_origin,omitempty"`
	// 
	Severity string `json:"severity,omitempty"`
	// 
	EventMessage string `json:"event_message,omitempty"`
	// event definition name
	EventName string `json:"event_name,omitempty"`
	// 
	Metadata *interface{} `json:"metadata,omitempty"`
}
