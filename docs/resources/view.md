---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "vastdata_view Resource - terraform-provider-vastdata"
subcategory: ""
description: |-
  
---

# vastdata_view (Resource)



## Example Usage

```terraform
# Create a view with NFS & NFSv4 protocols
resource "vastdata_view_policy" "example" {
  name   = "example"
  flavor = "NFS"
}


resource "vastdata_view" "example-view" {
  path       = "/example"
  policy_id  = vastdata_view_policy.example.id
  create_dir = "true"
  protocols  = ["NFS", "NFS4"]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `path` (String)
- `policy_id` (Number)

### Optional

- `alias` (String) Alias for NFS export, must start with '/' and only ASCII characters are allowed. If configured, this supersedes the exposed NFS export path
- `allow_anonymous_access` (Boolean) Allow S3 anonymous access
- `allow_s3_anonymous_access` (Boolean) Allow S3 anonymous access
- `bucket` (String) S3 Bucket name
- `bucket_creators` (List of String) List of bucket creators users
- `bucket_creators_groups` (List of String) List of bucket creators groups
- `bucket_owner` (String) S3 Bucket owner
- `cluster` (String) Parent Cluster
- `cluster_id` (Number) Parent Cluster ID
- `create_dir` (Boolean) Creates the directory specified by the path
- `directory` (Boolean) Create the directory if it does not exist
- `is_remote` (Boolean)
- `logical_capacity` (Number) Logical Capacity
- `name` (String) A uniq name given to the view
- `nfs_interop_flags` (String) Indicates whether the view should support simultaneous access to NFS3/NFS4/SMB protocols. Allowed Values are [BOTH_NFS3_AND_NFS4_INTEROP_DISABLED ONLY_NFS3_INTEROP_ENABLED ONLY_NFS4_INTEROP_ENABLED BOTH_NFS3_AND_NFS4_INTEROP_ENABLED]
- `physical_capacity` (Number) Physical Capacity
- `protocols` (List of String) Protocols exposed by this view
- `qos_policy_id` (Number) QoS Policy ID
- `s3_locks` (Boolean) S3 Object Lock
- `s3_locks_retention_mode` (String) S3 Locks retention mode
- `s3_locks_retention_period` (String) Period should be positive in format like 0d|2d|1y|2y
- `s3_unverified_lookup` (Boolean) Allow S3 Unverified Lookup
- `s3_versioning` (Boolean) Trun on S3 Versioning
- `share` (String) Name of the SMB Share. Must not include the following characters: " \ / [ ] : | < > + = ; , * ?
- `share_acl` (Block List) Share-level ACL details (see [below for nested schema](#nestedblock--share_acl))
- `tenant_id` (Number) The tenant ID related to this view

### Read-Only

- `guid` (String) A uniqe GUID assigned to the View
- `id` (String) The ID of this resource.

<a id="nestedblock--share_acl"></a>
### Nested Schema for `share_acl`

Optional:

- `acl` (Block List) (see [below for nested schema](#nestedblock--share_acl--acl))
- `enabled` (Boolean)

<a id="nestedblock--share_acl--acl"></a>
### Nested Schema for `share_acl.acl`

Required:

- `name` (String)

Optional:

- `fqdn` (String)
- `grantee` (String) Allowed Values are [users groups]
- `permissions` (String) Allowed Values are [FULL CHANGE READ]
- `sid_str` (String)
- `uid_or_gid` (String)

## Import

Import is supported using the following syntax:

```shell
terraform import vastdata_view.example <guid>
terraform import vastdata_view.example <Path>|<Tenant Name>
```
