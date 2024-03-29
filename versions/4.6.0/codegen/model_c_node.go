/*
 * VAST API Swagger Schema
 *
 * VAST Management API definition
 *
 * API version: 1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type CNode struct {
	// 
	Id int32 `json:"id,omitempty"`
	// 
	Guid string `json:"guid,omitempty"`
	// 
	Name string `json:"name"`
	// 
	NewName string `json:"new_name,omitempty"`
	// Currently used IP Address, bond of ip1 and ip2
	Ip string `json:"ip,omitempty"`
	// 1st internal IP Address
	Ip1 string `json:"ip1,omitempty"`
	// 2nd internal IP Address
	Ip2 string `json:"ip2,omitempty"`
	// External IPv6 Address
	Ipv6 string `json:"ipv6,omitempty"`
	// Host label, used to label container, e.g. 11.0.0.1-4000
	HostLabel string `json:"host_label,omitempty"`
	// The host serial number
	Sn string `json:"sn,omitempty"`
	// 
	State string `json:"state,omitempty"`
	// 
	DisplayState string `json:"display_state,omitempty"`
	// 
	LedStatus string `json:"led_status,omitempty"`
	// 
	PlatformRdmaPort int32 `json:"platform_rdma_port,omitempty"`
	// 
	PlatformTcpPort int32 `json:"platform_tcp_port,omitempty"`
	// 
	DataRdmaPort int32 `json:"data_rdma_port,omitempty"`
	// 
	DataTcpPort int32 `json:"data_tcp_port,omitempty"`
	// cnode enabled
	Enabled bool `json:"enabled,omitempty"`
	// The cnode is running management
	IsMgmt bool `json:"is_mgmt,omitempty"`
	// Parent CBox
	Cbox string `json:"cbox,omitempty"`
	// Unique Parent CBox identifier
	CboxUid string `json:"cbox_uid,omitempty"`
	// Parent CBox id
	CboxId int32 `json:"cbox_id,omitempty"`
	// Parent Cluster
	Cluster string `json:"cluster,omitempty"`
	// Node OS version
	OsVersion string `json:"os_version,omitempty"`
	// BMC FW version
	BmcFwVersion string `json:"bmc_fw_version,omitempty"`
	// 
	Url string `json:"url,omitempty"`
	// Read IOPS
	RdIops int64 `json:"rd_iops,omitempty"`
	// Write IOPS
	WrIops int64 `json:"wr_iops,omitempty"`
	// IOPS
	Iops int64 `json:"iops,omitempty"`
	// Read Meta-data IOPS
	RdMdIops int64 `json:"rd_md_iops,omitempty"`
	// Write Meta-data IOPS
	WrMdIops int64 `json:"wr_md_iops,omitempty"`
	// Meta-data IOPS
	MdIops int64 `json:"md_iops,omitempty"`
	// Read Bandwidth
	RdBw int64 `json:"rd_bw,omitempty"`
	// Write Bandwidth
	WrBw int64 `json:"wr_bw,omitempty"`
	// Bandwidth
	Bw int64 `json:"bw,omitempty"`
	// Read Latency
	RdLatency int64 `json:"rd_latency,omitempty"`
	// Write Latency
	WrLatency int64 `json:"wr_latency,omitempty"`
	// Latency
	Latency int64 `json:"latency,omitempty"`
	// CNode cores
	Cores int32 `json:"cores,omitempty"`
	// 
	Build string `json:"build,omitempty"`
	// Management IP
	MgmtIp string `json:"mgmt_ip,omitempty"`
	// Host Name
	Hostname string `json:"hostname,omitempty"`
	// Synchronization state with leader
	Sync string `json:"sync,omitempty"`
	// Synchronization time with leader
	SyncTime string `json:"sync_time,omitempty"`
}
