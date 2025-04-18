/*
 * VastData API
 *
 * A API document representing VastData API
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger

type QosStaticLimits struct {
	// Minimal amount of performance to provide when there is resource contention
	MinReadsBwMbps int64 `json:"min_reads_bw_mbps,omitempty"`
	// Maximal amount of performance to provide when there is no resource contention
	MaxReadsBwMbps int64 `json:"max_reads_bw_mbps,omitempty"`
	// Minimal amount of performance to provide when there is resource contention
	MinWritesBwMbps int64 `json:"min_writes_bw_mbps,omitempty"`
	// Maximal amount of performance to provide when there is no resource contention
	MaxWritesBwMbps int64 `json:"max_writes_bw_mbps,omitempty"`
	// Minimal amount of performance to provide when there is resource contention
	MinReadsIops int64 `json:"min_reads_iops,omitempty"`
	// Maximal amount of performance to provide when there is no resource contention
	MaxReadsIops int64 `json:"max_reads_iops,omitempty"`
	// Minimal amount of performance to provide when there is resource contention
	MinWritesIops int64 `json:"min_writes_iops,omitempty"`
	// Maximal amount of performance to provide when there is no resource contention
	MaxWritesIops int64 `json:"max_writes_iops,omitempty"`
	// Burst reads BW Mb
	BurstReadsBwMb int64 `json:"burst_reads_bw_mb,omitempty"`
	// Burst reads loan Mb
	BurstReadsLoanMb int64 `json:"burst_reads_loan_mb,omitempty"`
	// Burst writes BW Mb
	BurstWritesBwMb int64 `json:"burst_writes_bw_mb,omitempty"`
	// Burst writes loan Mb
	BurstWritesLoanMb int64 `json:"burst_writes_loan_mb,omitempty"`
	// Burst reads IOPS
	BurstReadsIops int64 `json:"burst_reads_iops,omitempty"`
	// Burst reads loan IOPS
	BurstReadsLoanIops int64 `json:"burst_reads_loan_iops,omitempty"`
	// Burst writes IOPS
	BurstWritesIops int64 `json:"burst_writes_iops,omitempty"`
	// Burst writes loan IOPS
	BurstWritesLoanIops int64 `json:"burst_writes_loan_iops,omitempty"`
}
