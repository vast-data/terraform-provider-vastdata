/*
 * VastData API
 *
 * A API document representing VastData API
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger

type View struct {
	// A uniqe ID used to identify the View
	Id int32 `json:"id,omitempty"`
	// A uniqe GUID assigned to the View
	Guid string `json:"guid,omitempty"`
	// A uniq name given to the view
	Name string `json:"name,omitempty"`
	// File system path. Begin with '/'. Do not include a trailing slash
	Path string `json:"path"`
	// Creates the directory specified by the path
	CreateDir bool `json:"create_dir,omitempty"`
	// Alias for NFS export, must start with '/' and only ASCII characters are allowed. If configured, this supersedes the exposed NFS export path
	Alias string `json:"alias,omitempty"`
	// S3 Bucket name
	Bucket string `json:"bucket,omitempty"`
	// Associated view policy ID
	PolicyId int32 `json:"policy_id,omitempty"`
	// Parent Cluster
	Cluster string `json:"cluster,omitempty"`
	// Parent Cluster ID
	ClusterId int32 `json:"cluster_id,omitempty"`
	// The tenant ID related to this view
	TenantId int32 `json:"tenant_id,omitempty"`
	// Create the directory if it does not exist
	Directory bool `json:"directory,omitempty"`
	// Trun on S3 Versioning
	S3Versioning bool `json:"s3_versioning,omitempty"`
	// Allow S3 Unverified Lookup
	S3UnverifiedLookup bool `json:"s3_unverified_lookup,omitempty"`
	// Allow S3 anonymous access
	AllowAnonymousAccess bool `json:"allow_anonymous_access,omitempty"`
	// Allow S3 anonymous access
	AllowS3AnonymousAccess bool `json:"allow_s3_anonymous_access,omitempty"`
	// Protocols exposed by this view
	Protocols []string `json:"protocols,omitempty"`
	// Name of the SMB Share. Must not include the following characters: \" \\ / [ ] : | < > + = ; , * ?
	Share string `json:"share,omitempty"`
	// S3 Bucket owner
	BucketOwner string `json:"bucket_owner,omitempty"`
	// List of bucket creators users
	BucketCreators []string `json:"bucket_creators,omitempty"`
	// List of bucket creators groups
	BucketCreatorsGroups []string `json:"bucket_creators_groups,omitempty"`
	// S3 Object Lock
	S3Locks bool `json:"s3_locks,omitempty"`
	// S3 Locks retention mode
	S3LocksRetentionMode string `json:"s3_locks_retention_mode,omitempty"`
	// Period should be positive in format like 0d|2d|1y|2y
	S3LocksRetentionPeriod string `json:"s3_locks_retention_period,omitempty"`
	// Physical Capacity
	PhysicalCapacity int64 `json:"physical_capacity,omitempty"`
	// Logical Capacity
	LogicalCapacity int64 `json:"logical_capacity,omitempty"`
	// Indicates whether the view should support simultaneous access to NFS3/NFS4/SMB protocols.
	NfsInteropFlags string `json:"nfs_interop_flags,omitempty"`
	IsRemote bool `json:"is_remote,omitempty"`
	ShareAcl *ViewShareAcl `json:"share_acl,omitempty"`
	// QoS Policy ID
	QosPolicyId int32 `json:"qos_policy_id,omitempty"`
	// Supports seamless failover between replication peers by syncing file handles between the view and remote views on the replicated path on replication peers. This enables NFSv3 client users to retain the same mount point to the view in the event of a failover of the view path to a replication peer. This feature enables NFSv3 client users to retain the same mount point to the view in the event of a failover of the view path to a replication peer. Enabling this option may cause overhead and should only be enabled when the use case is relevant. To complete the configuration for seamless failover between any two peers, a seamless view must be created on each peer.
	IsSeamless bool `json:"is_seamless,omitempty"`
	// Applicable if locking is enabled. Sets a maximum retention period for files that are locked in the view. Files cannot be locked for longer than this period, whether they are locked manually (by setting the atime) or automatically, using auto-commit. Specify as an integer value followed by a letter for the unit (m - minutes, h - hours, d - days, y - years). Example: 2y (2 years).
	MaxRetentionPeriod string `json:"max_retention_period,omitempty"`
	// Applicable if locking is enabled. Sets a minimum retention period for files that are locked in the view. Files cannot be locked for less than this period, whether locked manually (by setting the atime) or automatically, using auto-commit. Specify as an integer value followed by a letter for the unit (h - hours, d - days, m - months, y - years). Example: 1d (1 day).
	MinRetentionPeriod string `json:"min_retention_period,omitempty"`
	// Applicable if locking is enabled. The retention mode for new files. For views enabled for NFSv3 or SMB, if locking is enabled, files_retention_mode must be set to GOVERNANCE or COMPLIANCE. If the view is enabled for S3 and not for NFSv3 or SMB, files_retention_mode can be set to NONE. If GOVERNANCE, locked files cannot be deleted or changed. The Retention settings can be shortened or extended by users with sufficient permissions. If COMPLIANCE, locked files cannot be deleted or changed. Retention settings can be extended, but not shortened, by users with sufficient permissions. If NONE (S3 only), the retention mode is not set for the view; it is set individually for each object.
	FilesRetentionMode string `json:"files_retention_mode,omitempty"`
	// Relevant if locking is enabled. Required if s3_locks_retention_mode is set to governance or compliance. Specifies a default retention period for objects in the bucket. If set, object versions that are placed in the bucket are automatically protected with the specified retention lock. Otherwise, by default, each object version has no automatic protection but can be configured with a retention period or legal hold. Specify as an integer followed by h for hours, d for days, m for months, or y for years. For example: 2d or 1y.
	DefaultRetentionPeriod string `json:"default_retention_period,omitempty"`
	// Applicable if locking is enabled. Sets the auto-commit time for files that are locked automatically. These files are locked automatically after the auto-commit period elapses from the time the file is saved. Files locked automatically are locked for the default-retention-period, after which they are unlocked. Specify as an integer value followed by a letter for the unit (h - hours, d - days, y - years). Example: 2h (2 hours).
	AutoCommit string `json:"auto_commit,omitempty"`
	S3ObjectOwnershipRule string `json:"s3_object_ownership_rule,omitempty"`
	// Write Once Read Many (WORM) locking enabled
	Locking bool `json:"locking,omitempty"`
	// Ignore oos
	IgnoreOos bool `json:"ignore_oos,omitempty"`
	BucketLogging *BucketLogging `json:"bucket_logging,omitempty"`
	// List of attribute based access control tags, this option can be used only when using SMB/NFSv4 protocols
	AbacTags []string `json:"abac_tags,omitempty"`
	// Restricts ABE to a specified path depth. For example, if max depth is 3, ABE does not affect paths deeper than three levels. If not specified, ABE affects all path depths.
	AbeMaxDepth int32 `json:"abe_max_depth,omitempty"`
	// The protocols for which Access-Based Enumeration (ABE) is enabled , allowed values [ NFS, SMB, NFS4, S3 ]
	AbeProtocols []string `json:"abe_protocols,omitempty"`
	// Set as the default subsystem view for block devices (sub-system)
	IsDefaultSubsystem bool `json:"is_default_subsystem,omitempty"`
}
