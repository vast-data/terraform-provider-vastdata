/*
 * VAST API Swagger Schema
 *
 * VAST Management API definition
 *
 * API version: 1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type SystemShardExpandParams struct {
	// EStore shard count
	EstoreShardCount int32 `json:"estore_shard_count,omitempty"`
	// DR shard count
	DrShardCount int32 `json:"dr_shard_count,omitempty"`
	// DR WB shard count
	DrWbShardCount int32 `json:"dr_wb_shard_count,omitempty"`
	// Force shard expansion. Use if you want to run shard expansion even though shards are denylisted (indicated by MAINTENANCE_DENYLIST_EXISTS in error code when running without 'force').
	Force bool `json:"force,omitempty"`
}
