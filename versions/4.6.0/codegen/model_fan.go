/*
 * VAST API Swagger Schema
 *
 * VAST Management API definition
 *
 * API version: 1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type Fan struct {
	// 
	Id int32 `json:"id,omitempty"`
	// 
	Guid string `json:"guid,omitempty"`
	// 
	Name string `json:"name,omitempty"`
	// Fan state
	State string `json:"state,omitempty"`
	// Fan serial number
	Sn string `json:"sn,omitempty"`
	// Fan model
	Model string `json:"model,omitempty"`
	// Parent Box name
	Box string `json:"box,omitempty"`
	// Parent Cluster
	Cluster string `json:"cluster,omitempty"`
	// Fan index
	Index int32 `json:"index,omitempty"`
	// 
	Location string `json:"location,omitempty"`
	// Fan speed
	Speed int32 `json:"speed,omitempty"`
}
