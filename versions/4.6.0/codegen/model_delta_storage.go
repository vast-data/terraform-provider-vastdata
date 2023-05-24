/*
 * VAST API Swagger Schema
 *
 * VAST Management API definition
 *
 * API version: 1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type DeltaStorage struct {
	// Runtime sequence id for delta. Reset if HA happened.
	CurrentSequence string `json:"current_sequence"`
	CurrentGeneration int32 `json:"current_generation"`
	Records []DeltaRecord `json:"records"`
	// ok - all ok, reset - journal was reset
	Status string `json:"status"`
}
