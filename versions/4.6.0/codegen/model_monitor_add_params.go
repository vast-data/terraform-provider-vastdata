/*
 * VAST API Swagger Schema
 *
 * VAST Management API definition
 *
 * API version: 1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type MonitorAddParams struct {
	// A name for the monitor.
	Name string `json:"name,omitempty"`
	// The type of object to monitor.
	ObjectType string `json:"object_type"`
	// 
	FromTime string `json:"from_time,omitempty"`
	// 
	ToTime string `json:"to_time,omitempty"`
	// Default time frame to report over. Examples: 2h (2 hours), 1D (1 Day), 10m (10 minutes), 1M (1 month)
	TimeFrame string `json:"time_frame,omitempty"`
	// Specific objects to include in the report, specified as a comma separated list of object IDs.
	ObjectIds string `json:"object_ids,omitempty"`
	// A list of metrics to query. To get the full list of metrics, use GET /metrics/.
	PropList string `json:"prop_list,omitempty"`
	// Data granularity: seconds (raw), minutes (five minute aggregated samples), hours (hourly aggregated samples), or days (daily aggregated samples)
	Granularity string `json:"granularity,omitempty"`
	// If data granularity is minutes, hours or days, the data is aggregated. This parameter selects which aggregation function to use.
	Aggregation string `json:"aggregation,omitempty"`
}
