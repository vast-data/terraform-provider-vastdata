/*
 * VAST Data API
 *
 * An API document representing VAST Data API
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger

type QosDynamicLimits struct {
	// Maximum read bandwidth (in MB/s) per GB to provide when there is no resource contention.
	MaxReadsBwMbpsPerGbCapacity int64 `json:"max_reads_bw_mbps_per_gb_capacity,omitempty"`
	// Maximum write bandwidth (in MB/s) per GB to provide when there is no resource contention.
	MaxWritesBwMbpsPerGbCapacity int64 `json:"max_writes_bw_mbps_per_gb_capacity,omitempty"`
	// Maximum read IOPS per GB to provide when there is no resource contention.
	MaxReadsIopsPerGbCapacity int64 `json:"max_reads_iops_per_gb_capacity,omitempty"`
	// Maximum write IOPS per GB to provide when there is no resource contention.
	MaxWritesIopsPerGbCapacity int64 `json:"max_writes_iops_per_gb_capacity,omitempty"`
}
