---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "vastdata_qos_policy Data Source - terraform-provider-vastdata"
subcategory: ""
description: |-
  
---

# vastdata_qos_policy (Data Source)



## Example Usage

```terraform
data "vastdata_qos_policy" "qos1" {
  name = "qos1"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0)

### Optional

- `limit_by` (String) (Valid for versions: 5.2.0) Specifies which attributes are setting the limitations. Allowed Values are [BW_IOPS BW IOPS]

### Read-Only

- `attached_users` (List of Object) (Valid for versions: 5.2.0) (see [below for nested schema](#nestedatt--attached_users))
- `attached_users_identifiers` (List of String) (Valid for versions: 5.2.0) A list of local user IDs to which this QoS policy applies.
- `capacity_limits` (List of Object) (Valid for versions: 5.0.0,5.1.0,5.2.0) (see [below for nested schema](#nestedatt--capacity_limits))
- `capacity_total_limits` (List of Object) (Valid for versions: 5.2.0) (see [below for nested schema](#nestedatt--capacity_total_limits))
- `guid` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0) QoS policy GUID.
- `id` (Number) (Valid for versions: 5.0.0,5.1.0,5.2.0)
- `io_size_bytes` (Number) (Valid for versions: 5.0.0,5.1.0,5.2.0) Sets the size of IO for static and capacity limit definitions. The number of IOs per request is obtained by dividing the request size by IO size. Default: 64K. Recommended range: 4K - 1M.
- `is_default` (Boolean) (Valid for versions: 5.2.0) Specifies whether this QoS policy is to be used as the default QoS policy per user for this tenant. Setting this attribute requires that 'tenant_id' is also supplied.
- `mode` (String) (Valid for versions: 5.0.0,5.1.0,5.2.0) QoS provisioning mode. Allowed Values are [STATIC USED_CAPACITY PROVISIONED_CAPACITY]
- `policy_type` (String) (Valid for versions: 5.2.0) The QoS policy type. Allowed Values are [VIEW USER]
- `static_limits` (List of Object) (Valid for versions: 5.0.0,5.1.0,5.2.0) (see [below for nested schema](#nestedatt--static_limits))
- `static_total_limits` (List of Object) (Valid for versions: 5.2.0) (see [below for nested schema](#nestedatt--static_total_limits))
- `tenant_id` (Number) (Valid for versions: 5.2.0) When setting 'is_default', this is the tenant for which the policy will be used as the default user QoS policy.

<a id="nestedatt--attached_users"></a>
### Nested Schema for `attached_users`

Read-Only:

- `fqdn` (String)
- `identifier_type` (String)
- `identifier_value` (String)
- `is_sid` (Boolean)
- `label` (String)
- `login_name` (String)
- `name` (String)
- `sid_str` (String)
- `uid_or_gid` (Number)
- `value` (String)


<a id="nestedatt--capacity_limits"></a>
### Nested Schema for `capacity_limits`

Read-Only:

- `max_reads_bw_mbps_per_gb_capacity` (Number)
- `max_reads_iops_per_gb_capacity` (Number)
- `max_writes_bw_mbps_per_gb_capacity` (Number)
- `max_writes_iops_per_gb_capacity` (Number)


<a id="nestedatt--capacity_total_limits"></a>
### Nested Schema for `capacity_total_limits`

Read-Only:

- `max_bw_mbps_per_gb_capacity` (Number)
- `max_iops_per_gb_capacity` (Number)


<a id="nestedatt--static_limits"></a>
### Nested Schema for `static_limits`

Read-Only:

- `burst_reads_bw_mb` (Number)
- `burst_reads_iops` (Number)
- `burst_reads_loan_iops` (Number)
- `burst_reads_loan_mb` (Number)
- `burst_writes_bw_mb` (Number)
- `burst_writes_iops` (Number)
- `burst_writes_loan_iops` (Number)
- `burst_writes_loan_mb` (Number)
- `max_reads_bw_mbps` (Number)
- `max_reads_iops` (Number)
- `max_writes_bw_mbps` (Number)
- `max_writes_iops` (Number)
- `min_reads_bw_mbps` (Number)
- `min_reads_iops` (Number)
- `min_writes_bw_mbps` (Number)
- `min_writes_iops` (Number)


<a id="nestedatt--static_total_limits"></a>
### Nested Schema for `static_total_limits`

Read-Only:

- `burst_bw_mb` (Number)
- `burst_iops` (Number)
- `burst_loan_iops` (Number)
- `burst_loan_mb` (Number)
- `max_bw_mbps` (Number)
- `max_iops` (Number)
