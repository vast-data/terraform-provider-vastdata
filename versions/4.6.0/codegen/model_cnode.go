package swagger

type Cnode struct {
	//Tehc Cnode name
	Name string `json:"name,omitempty"`
	//The Cnode ip
	Ip     string `json:"ip,omitempty"`
	Ip1    string `json:"ip1,omitempty"`
	Ip2    string `json:"ip2,omitempty"`
	MgmtIp string `json:"mgmt_ip,omitempty"`
	Ipv6   string `json:"ipv6,omitempty"`
	//Is the cnode enabled
	Enabled bool `json:"enabled,omitempty"`
	//The id of the Vippool this node is assigned to
	Id int32 `json:"id,omitempty"`
	//The Guid of the Cnode
	Guid      string `json:"guid,omitempty"`
	OsVersion string `json:"os_version,omitempty"`
	HostLabel string `json:"host_label,omitempty"`
	NewName   string `json:"new_name,omitempty"`
}
