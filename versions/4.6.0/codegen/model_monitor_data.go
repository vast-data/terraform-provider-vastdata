/*
 * VAST API Swagger Schema
 *
 * VAST Management API definition
 *
 * API version: 1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type MonitorData struct {
	// list of object ids that appears in the query results
	ObjectIds *interface{} `json:"object_ids,omitempty"`
	// list of metrics that appears in the result
	PropList *interface{} `json:"prop_list,omitempty"`
	// list of data samples, each sample in form of [timestamp, obj_id, *data]
	Data *interface{} `json:"data,omitempty"`
}
