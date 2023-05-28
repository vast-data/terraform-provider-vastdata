/*
 * VAST API Swagger Schema
 *
 * VAST Management API definition
 *
 * API version: 1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type QosStaticLimits struct {
	// Minimal amount of performance to provide when there is resource contention
	MinReadsBwMbps int32 `json:"min_reads_bw_mbps"`
	// Maximal amount of performance to provide when there is no resource contention
	MaxReadsBwMbps int32 `json:"max_reads_bw_mbps"`
	// Minimal amount of performance to provide when there is resource contention
	MinWritesBwMbps int32 `json:"min_writes_bw_mbps"`
	// Maximal amount of performance to provide when there is no resource contention
	MaxWritesBwMbps int32 `json:"max_writes_bw_mbps"`
	// Minimal amount of performance to provide when there is resource contention
	MinReadsIops int32 `json:"min_reads_iops"`
	// Maximal amount of performance to provide when there is no resource contention
	MaxReadsIops int32 `json:"max_reads_iops"`
	// Minimal amount of performance to provide when there is resource contention
	MinWritesIops int32 `json:"min_writes_iops"`
	// Maximal amount of performance to provide when there is no resource contention
	MaxWritesIops int32 `json:"max_writes_iops"`
}