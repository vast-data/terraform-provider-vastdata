/*
 * VAST API Swagger Schema
 *
 * VAST Management API definition
 *
 * API version: 1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type VTask struct {
	// 
	Id int32 `json:"id,omitempty"`
	// Task state
	State string `json:"state,omitempty"`
	// Messages passed to polling clients.
	Messages *interface{} `json:"messages,omitempty"`
	// Extra task relate information passed to polling clients.
	Info *interface{} `json:"info,omitempty"`
	// Task should timeout after X seconds.
	TimeoutInSeconds int32 `json:"timeout_in_seconds,omitempty"`
	// Task name
	Name string `json:"name,omitempty"`
	// Task start time
	StartTime string `json:"start_time,omitempty"`
	// Task end time
	EndTime string `json:"end_time,omitempty"`
	// Task execution time in seconds
	ExecutionTime int32 `json:"execution_time,omitempty"`
}
