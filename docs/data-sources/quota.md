---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "vastdata_quota Data Source - terraform-provider-vastdata"
subcategory: ""
description: |-
  This is a quota
---

# vastdata_quota (Data Source)

This is a quota



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name

### Optional

- `cluster` (String) Parent Cluster
- `cluster_id` (Number) Parent Cluster ID
- `default_email` (String) The default Email if there is no suffix and no address in the providers
- `default_group_quota` (Block List) (see [below for nested schema](#nestedblock--default_group_quota))
- `default_user_quota` (Block List) (see [below for nested schema](#nestedblock--default_user_quota))
- `enable_alarms` (Boolean) Enable alarms when users or groups are exceeding their limit
- `enable_email_providers` (Boolean)
- `grace_period` (String) Quota enforcement grace period in seconds, minutes, hours or days. Example: 90m
- `group_quotas` (Block List) (see [below for nested schema](#nestedblock--group_quotas))
- `hard_limit` (Number) Hard quota limit
- `hard_limit_inodes` (Number) Hard inodes quota limit
- `is_user_quota` (Boolean)
- `num_blocked_users` (Number)
- `num_exceeded_users` (Number)
- `path` (String) Directory path
- `percent_capacity` (Number) Percent of used capacity out of the hard limit
- `percent_inodes` (Number) Percent of used inodes out of the hard limit
- `pretty_grace_period` (String) Quota enforcement pretty grace period in seconds, minutes, hours or days. Example: 90m
- `pretty_state` (String)
- `soft_limit` (Number) Soft quota limit
- `soft_limit_inodes` (Number) Soft inodes quota limit
- `state` (String)
- `system_id` (Number)
- `tenant_id` (Number) Tenant ID
- `tenant_name` (String) Tenant Name
- `time_to_block` (String) Grace period expiration time
- `used_capacity` (Number) Used capacity in bytes
- `used_capacity_tb` (Number) Used capacity in TB
- `used_effective_capacity` (Number) Used effective capacity in bytes
- `used_effective_capacity_tb` (Number) Used effective capacity in TB
- `used_inodes` (Number) Used inodes
- `user_quotas` (Block List) (see [below for nested schema](#nestedblock--user_quotas))

### Read-Only

- `guid` (String) Quota guid
- `id` (Number) The ID of this resource.

<a id="nestedblock--default_group_quota"></a>
### Nested Schema for `default_group_quota`

Optional:

- `grace_period` (String)
- `hard_limit` (Number) The size hard limit in bytes
- `hard_limit_inodes` (Number) The hard limit in inode number
- `quota_system_id` (Number) The system ID of the quota
- `sof_limit_inodes` (Number) The sof limit of inodes number
- `soft_limit` (Number) The size soft limit in bytes


<a id="nestedblock--default_user_quota"></a>
### Nested Schema for `default_user_quota`

Optional:

- `grace_period` (String)
- `hard_limit` (Number) The size hard limit in bytes
- `hard_limit_inodes` (Number) The hard limit in inode number
- `quota_system_id` (Number) The system ID of the quota
- `sof_limit_inodes` (Number) The sof limit of inodes number
- `soft_limit` (Number) The size soft limit in bytes


<a id="nestedblock--group_quotas"></a>
### Nested Schema for `group_quotas`

Optional:

- `entity` (Block List) (see [below for nested schema](#nestedblock--group_quotas--entity))
- `grace_period` (String) Quota enforcement grace period in seconds, minutes, hours or days. Example: 90m
- `hard_limit` (Number) Hard quota limit
- `hard_limit_inodes` (Number) Hard inodes quota limit
- `is_accountable` (Boolean)
- `quota_system_id` (Number)
- `soft_limit` (Number) Soft quota limit
- `soft_limit_inodes` (Number) Soft inodes quota limit
- `time_to_block` (String) Grace period expiration time
- `used_capacity` (Number) Used capacity in bytes
- `used_inodes` (Number) Used inodes

<a id="nestedblock--group_quotas--entity"></a>
### Nested Schema for `group_quotas.entity`

Required:

- `name` (String) The name of the entity

Optional:

- `email` (String)
- `identifier` (String)
- `identifier_type` (String)
- `is_group` (Boolean)
- `vast_id` (Number)



<a id="nestedblock--user_quotas"></a>
### Nested Schema for `user_quotas`

Optional:

- `entity` (Block List) (see [below for nested schema](#nestedblock--user_quotas--entity))
- `grace_period` (String) Quota enforcement grace period in seconds, minutes, hours or days. Example: 90m
- `hard_limit` (Number) Hard quota limit
- `hard_limit_inodes` (Number) Hard inodes quota limit
- `is_accountable` (Boolean)
- `quota_system_id` (Number)
- `soft_limit` (Number) Soft quota limit
- `soft_limit_inodes` (Number) Soft inodes quota limit
- `time_to_block` (String) Grace period expiration time
- `used_capacity` (Number) Used capacity in bytes
- `used_inodes` (Number) Used inodes

<a id="nestedblock--user_quotas--entity"></a>
### Nested Schema for `user_quotas.entity`

Required:

- `name` (String) The name of the entity

Optional:

- `email` (String)
- `identifier` (String)
- `identifier_type` (String)
- `is_group` (Boolean)
- `vast_id` (Number)