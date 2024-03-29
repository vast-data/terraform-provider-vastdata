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

type Carrier struct {
	// 
	Id int32 `json:"id,omitempty"`
	// 
	Guid string `json:"guid,omitempty"`
	// NIC state
	State string `json:"state,omitempty"`
	InsertionTime time.Time `json:"insertion_time,omitempty"`
	// 
	LedStatus string `json:"led_status,omitempty"`
	// 
	Dbox string `json:"dbox,omitempty"`
	// 
	DboxId int32 `json:"dbox_id,omitempty"`
	// 
	Ssds *interface{} `json:"ssds,omitempty"`
	// 
	Nvrams *interface{} `json:"nvrams,omitempty"`
	// 
	CarrierIndex string `json:"carrier_index,omitempty"`
	// The shelf (left/right/front/rear)
	Shelf string `json:"shelf,omitempty"`
	// Disks type (SSD / NVRAM)
	CarrierType string `json:"carrier_type,omitempty"`
	// HW version
	HwVersion string `json:"hw_version,omitempty"`
	// SW version
	SwVersion string `json:"sw_version,omitempty"`
	// FW version
	FwVersion string `json:"fw_version,omitempty"`
	// Carrier serial number
	Sn string `json:"sn,omitempty"`
}
