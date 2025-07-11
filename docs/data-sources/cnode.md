---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "vastdata_cnode Data Source - terraform-provider-vastdata"
subcategory: ""
description: |-
  
---

# vastdata_cnode (Data Source)



## Example Usage

```terraform
data "vastdata_cnode" "cnode1" {
  name = "cnode1"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0)

### Read-Only

- `bmc_fw_version` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0) BMC FW version.
- `cbox` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0) Parent CBox.
- `cbox_id` (Number) (Valid for versions: 5.0.0,5.1.0,5.2.0) ID of the parent CBox.
- `cbox_uid` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0) UID of the parent CBox.
- `cluster` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0) Parent cluster.
- `data_rdma_port` (Number) (Valid for versions: 5.0.0,5.1.0,5.2.0)
- `data_tcp_port` (Number) (Valid for versions: 5.0.0,5.1.0,5.2.0)
- `display_state` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0)
- `enabled` (Boolean) (Valid for versions: 5.0.0,5.1.0,5.2.0) Sets the CNode to be enabled or disabled.
- `guid` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0)
- `host_label` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0) Host label, used to label the  container, e.g. 11.0.0.1-4000.
- `hostname` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0) Host name.
- `id` (Number) (Valid for versions: 5.0.0,5.1.0,5.2.0)
- `ip` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0) The IP address being used, bond of 'ip1' and 'ip2'.
- `ip1` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0) First internal IP address.
- `ip2` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0) Second internal IP address.
- `ipv6` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0) External IPv6 address.
- `is_mgmt` (Boolean) (Valid for versions: 5.0.0,5.1.0,5.2.0) Specifies whether the CNode is running VMS.
- `led_status` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0)
- `mgmt_ip` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0) Management IP.
- `new_name` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0)
- `os_version` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0) Node OS version.
- `platform_rdma_port` (Number) (Valid for versions: 5.0.0,5.1.0,5.2.0)
- `platform_tcp_port` (Number) (Valid for versions: 5.0.0,5.1.0,5.2.0)
- `sn` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0) Host serial number.
- `state` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0)
- `url` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0)
- `vms_preferred` (Boolean) (Valid for versions: 5.0.0,5.1.0,5.2.0) VMS preferred CNode.
