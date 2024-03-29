/*
 * VAST API Swagger Schema
 *
 * VAST Management API definition
 *
 * API version: 1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type AddBoxesCBoxesHostsParamsInner struct {
	Id int32 `json:"id"`
	NetType string `json:"net_type"`
	// relevant only for large subnet
	NbEthMtu int32 `json:"nb_eth_mtu,omitempty"`
	// relevant only for large subnet
	NbIbMtu int32 `json:"nb_ib_mtu,omitempty"`
	// relevant only for large subnet
	NbIbMode string `json:"nb_ib_mode,omitempty"`
	ReverseNics bool `json:"reverse_nics"`
	SkipNic string `json:"skip_nic"`
}
