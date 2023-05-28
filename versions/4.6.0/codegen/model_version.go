/*
 * VAST API Swagger Schema
 *
 * VAST Management API definition
 *
 * API version: 1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

import (
	"time"
)

type Version struct {
	// 
	Created time.Time `json:"created,omitempty"`
	// 
	Cluster string `json:"cluster,omitempty"`
	// 
	SysVersion string `json:"sys_version,omitempty"`
	// 
	Build string `json:"build,omitempty"`
	// 
	OsVersion string `json:"os_version,omitempty"`
	// 
	CarriersFwVersion string `json:"carriers_fw_version,omitempty"`
	// 
	SkipHwCheck bool `json:"skip_hw_check,omitempty"`
	// 
	Force bool `json:"force,omitempty"`
	// 
	SkipOsUpgrade bool `json:"skip_os_upgrade,omitempty"`
	// 
	EnableDr bool `json:"enable_dr,omitempty"`
	// 
	Status string `json:"status,omitempty"`
}