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

- `alias` (String) (Valid for versions: 5.0.0,5.1.0) Alias for NFS export, must start with '/' and only ASCII characters are allowed. If configured, this supersedes the exposed NFS export path
- `allow_anonymous_access` (Boolean) (Valid for versions: 5.0.0,5.1.0) Allow S3 anonymous access
- `allow_s3_anonymous_access` (Boolean) (Valid for versions: 5.0.0,5.1.0) Allow S3 anonymous access
- `auto_commit` (String) (Valid for versions: 5.1.0) Applicable if locking is enabled. Sets the auto-commit time for files that are locked automatically. These files are locked automatically after the auto-commit period elapses from the time the file is saved. Files locked automatically are locked for the default-retention-period, after which they are unlocked. Specify as an integer value followed by a letter for the unit (h - hours, d - days, y - years). Example: 2h (2 hours).
- `bucket` (String) (Valid for versions: 5.0.0,5.1.0) S3 Bucket name
- `bucket_creators` (List of String) (Valid for versions: 5.0.0,5.1.0) List of bucket creators users
- `bucket_creators_groups` (List of String) (Valid for versions: 5.0.0,5.1.0) List of bucket creators groups
- `bucket_owner` (String) (Valid for versions: 5.0.0,5.1.0) S3 Bucket owner
- `cluster` (String) (Valid for versions: 5.0.0,5.1.0) Parent Cluster
- `cluster_id` (Number) (Valid for versions: 5.0.0,5.1.0) Parent Cluster ID
- `create_dir` (Boolean) (Valid for versions: 5.0.0,5.1.0) Creates the directory specified by the path
- `default_retention_period` (String) (Valid for versions: 5.1.0) Relevant if locking is enabled. Required if s3_locks_retention_mode is set to governance or compliance. Specifies a default retention period for objects in the bucket. If set, object versions that are placed in the bucket are automatically protected with the specified retention lock. Otherwise, by default, each object version has no automatic protection but can be configured with a retention period or legal hold. Specify as an integer followed by h for hours, d for days, m for months, or y for years. For example: 2d or 1y.
- `directory` (Boolean) (Valid for versions: 5.0.0,5.1.0) Create the directory if it does not exist
- `files_retention_mode` (String) (Valid for versions: 5.1.0) Applicable if locking is enabled. The retention mode for new files. For views enabled for NFSv3 or SMB, if locking is enabled, files_retention_mode must be set to GOVERNANCE or COMPLIANCE. If the view is enabled for S3 and not for NFSv3 or SMB, files_retention_mode can be set to NONE. If GOVERNANCE, locked files cannot be deleted or changed. The Retention settings can be shortened or extended by users with sufficient permissions. If COMPLIANCE, locked files cannot be deleted or changed. Retention settings can be extended, but not shortened, by users with sufficient permissions. If NONE (S3 only), the retention mode is not set for the view; it is set individually for each object. Allowed Values are [GOVERNANCE COMPLIANCE NONE]
- `ignore_oos` (Boolean) (Valid for versions: 5.1.0) Ignore oos
- `is_remote` (Boolean) (Valid for versions: 5.0.0,5.1.0)
- `is_seamless` (Boolean) (Valid for versions: 5.1.0) Supports seamless failover between replication peers by syncing file handles between the view and remote views on the replicated path on replication peers. This enables NFSv3 client users to retain the same mount point to the view in the event of a failover of the view path to a replication peer. This feature enables NFSv3 client users to retain the same mount point to the view in the event of a failover of the view path to a replication peer. Enabling this option may cause overhead and should only be enabled when the use case is relevant. To complete the configuration for seamless failover between any two peers, a seamless view must be created on each peer.
- `locking` (Boolean) (Valid for versions: 5.1.0) Write Once Read Many (WORM) locking enabled
- `logical_capacity` (Number) (Valid for versions: 5.0.0,5.1.0) Logical Capacity
- `max_retention_period` (String) (Valid for versions: 5.1.0) Applicable if locking is enabled. Sets a maximum retention period for files that are locked in the view. Files cannot be locked for longer than this period, whether they are locked manually (by setting the atime) or automatically, using auto-commit. Specify as an integer value followed by a letter for the unit (m - minutes, h - hours, d - days, y - years). Example: 2y (2 years).
- `min_retention_period` (String) (Valid for versions: 5.1.0) Applicable if locking is enabled. Sets a minimum retention period for files that are locked in the view. Files cannot be locked for less than this period, whether locked manually (by setting the atime) or automatically, using auto-commit. Specify as an integer value followed by a letter for the unit (h - hours, d - days, m - months, y - years). Example: 1d (1 day).
- `name` (String) (Valid for versions: 5.0.0,5.1.0) A uniq name given to the view
- `nfs_interop_flags` (String) (Valid for versions: 5.0.0,5.1.0) Indicates whether the view should support simultaneous access to NFS3/NFS4/SMB protocols. Allowed Values are [BOTH_NFS3_AND_NFS4_INTEROP_DISABLED ONLY_NFS3_INTEROP_ENABLED ONLY_NFS4_INTEROP_ENABLED BOTH_NFS3_AND_NFS4_INTEROP_ENABLED]
- `physical_capacity` (Number) (Valid for versions: 5.0.0,5.1.0) Physical Capacity
- `protocols` (List of String) (Valid for versions: 5.0.0,5.1.0) Protocols exposed by this view
- `qos_policy_id` (Number) (Valid for versions: 5.0.0,5.1.0) QoS Policy ID
- `s3_locks` (Boolean) (Valid for versions: 5.0.0,5.1.0) S3 Object Lock
- `s3_locks_retention_mode` (String) (Valid for versions: 5.0.0,5.1.0) S3 Locks retention mode
- `s3_locks_retention_period` (String) (Valid for versions: 5.0.0,5.1.0) Period should be positive in format like 0d|2d|1y|2y
- `s3_object_ownership_rule` (String) (Valid for versions: 5.1.0) S3 Object Ownership lets you set ownership of objects uploaded to a given bucket and to determine whether ACLs are used to control access to objects within this bucket. A bucket can be configured with one of the following object ownership rules: BucketOwnerEnforced - The bucket owner has full control over any object in the bucket ObjectWriter - The user that uploads an object has full control over this object. ACLs can be used to let other users access the object. BucketOwnerPreferred - The bucket owner has full control over new objects uploaded to the bucket by other users. ACLs can be used to control access to the objects. None - S3 Object Ownership is disabled for the bucket.  Allowed Values are [None BucketOwnerPreferred ObjectWriter BucketOwnerEnforced]
- `s3_unverified_lookup` (Boolean) (Valid for versions: 5.0.0,5.1.0) Allow S3 Unverified Lookup
- `s3_versioning` (Boolean) (Valid for versions: 5.0.0,5.1.0) Trun on S3 Versioning
- `share` (String) (Valid for versions: 5.0.0,5.1.0) Name of the SMB Share. Must not include the following characters: " \ / [ ] : | < > + = ; , * ?
- `share_acl` (Block List) (Valid for versions: 5.0.0,5.1.0) Share-level ACL details (see [below for nested schema](#nestedblock--share_acl))
- `tenant_id` (Number) (Valid for versions: 5.0.0,5.1.0) The tenant ID related to this view

### Read-Only

- `guid` (String) (Valid for versions: 5.0.0,5.1.0) A uniqe GUID assigned to the View
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

- `fqdn` (String) (Valid for versions: 5.0.0,5.1.0)
- `grantee` (String) (Valid for versions: 5.0.0,5.1.0)  Allowed Values are [users groups]
- `permissions` (String) (Valid for versions: 5.0.0,5.1.0)  Allowed Values are [FULL CHANGE READ]
- `sid_str` (String) (Valid for versions: 5.0.0,5.1.0)
- `uid_or_gid` (Number) (Valid for versions: 5.0.0,5.1.0)

## Import

Import is supported using the following syntax:

```shell
terraform import vastdata_view.example <guid>
terraform import vastdata_view.example <Path>|<Tenant Name>
```
