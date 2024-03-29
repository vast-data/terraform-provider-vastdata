/*
 * VAST API Swagger Schema
 *
 * VAST Management API definition
 *
 * API version: 1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type Host struct {
	// 
	Id int32 `json:"id,omitempty"`
	// 
	Guid string `json:"guid,omitempty"`
	// 
	VastInstallInfo *interface{} `json:"vast_install_info,omitempty"`
	// internal IP Bond Address
	Ip string `json:"ip,omitempty"`
	// 1st internal IP Address
	Ip1 string `json:"ip1,omitempty"`
	// 2nd internal IP Address
	Ip2 string `json:"ip2,omitempty"`
	// External IPv6 Address
	Ipv6 string `json:"ipv6,omitempty"`
	// Host label, used to label container, e.g. 11.0.0.1-4000
	HostLabel string `json:"host_label,omitempty"`
	// auto or manually host
	Auto bool `json:"auto,omitempty"`
	// Hostname
	Hostname string `json:"hostname"`
	// Node type: C-Node/D-Node
	NodeType string `json:"node_type"`
	// Unique h/w identifier
	DboxUid string `json:"dbox_uid,omitempty"`
	// Loopback (single node) installation
	Loopback bool `json:"loopback,omitempty"`
	// Check cluster performance
	PerfCheck bool `json:"perf_check,omitempty"`
	// Mock NVMeoF devices with dmsetup devices
	Dmsetup bool `json:"dmsetup,omitempty"`
	// enable data reduction
	EnableDr bool `json:"enable_dr,omitempty"`
	// Is half system
	HalfSystem bool `json:"half_system,omitempty"`
	// Is deep stripe system
	DeepStripe bool `json:"deep_stripe,omitempty"`
	// DR hash size in buckets
	DrHashSize int32 `json:"dr_hash_size,omitempty"`
	// NVRAM size for mocked devices
	NvramSize int32 `json:"nvram_size,omitempty"`
	// Drive size for mocked devices
	DriveSize int32 `json:"drive_size,omitempty"`
	// SSH User name
	SshUser string `json:"ssh_user,omitempty"`
	// Node state
	State string `json:"state,omitempty"`
	// 
	PlatformRdmaPort int32 `json:"platform_rdma_port,omitempty"`
	// General hardware info related to the host
	HwInfo *interface{} `json:"hw_info,omitempty"`
	// 
	PlatformTcpPort int32 `json:"platform_tcp_port,omitempty"`
	// 
	DataRdmaPort int32 `json:"data_rdma_port,omitempty"`
	// 
	DataTcpPort int32 `json:"data_tcp_port,omitempty"`
	// 
	BoxRdmaPort int32 `json:"box_rdma_port,omitempty"`
	// Node installation state
	InstallState string `json:"install_state,omitempty"`
	// Node s/w version
	SwVersion string `json:"sw_version,omitempty"`
	// Node OS version
	OsVersion string `json:"os_version,omitempty"`
	// Node build
	Build string `json:"build,omitempty"`
	// Parent Cluster
	Cluster string `json:"cluster,omitempty"`
	// 
	Url string `json:"url,omitempty"`
	// Host upgrade state
	UpgradeState string `json:"upgrade_state,omitempty"`
	// Hosts sibling node guid
	NodeGuid string `json:"node_guid,omitempty"`
	// Host has single NIC
	SingleNic string `json:"single_nic,omitempty"`
	NetType string `json:"net_type,omitempty"`
}
