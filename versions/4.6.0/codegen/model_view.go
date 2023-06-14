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

type View struct {
	// 
	Id int32 `json:"id,omitempty"`
	// 
	Guid string `json:"guid,omitempty"`
	// 
	Name string `json:"name,omitempty"`
	Created time.Time `json:"created,omitempty"`
	// File system path. Begin with '/'. Do not include a trailing slash
	Path string `json:"path"`
	// Creates the directory specified by the path
	CreateDir string `json:"create_dir,omitempty"`
	// Alias for NFS export, must start with '/' and only ASCII characters are allowed. If configured, this supersedes the exposed NFS export path
	Alias string `json:"alias,omitempty"`
	// S3 Bucket name
	Bucket string `json:"bucket,omitempty"`
	// The view policy that applies to this view
	Policy string `json:"policy,omitempty"`
	// Associated view policy ID
	PolicyId int32 `json:"policy_id,omitempty"`
	// Parent Cluster
	Cluster string `json:"cluster,omitempty"`
	// Parent Cluster ID
	ClusterId int32 `json:"cluster_id,omitempty"`
	// Tenant ID
	TenantId int32 `json:"tenant_id,omitempty"`
	// 
	Url string `json:"url,omitempty"`
	// Create the directory if it does not exist
	Directory bool `json:"directory,omitempty"`
	// S3 Versioning
	S3Versioning bool `json:"s3_versioning,omitempty"`
	// S3 Unverified Lookup
	S3UnverifiedLookup bool `json:"s3_unverified_lookup,omitempty"`
	// Allow S3 anonymous access
	AllowAnonymousAccess bool `json:"allow_anonymous_access,omitempty"`
	// Allow S3 anonymous access
	AllowS3AnonymousAccess bool `json:"allow_s3_anonymous_access,omitempty"`
	// Protocols exposed by this view
	Protocols []string `json:"protocols,omitempty"`
	// Name of the SMB Share. Must not include the following characters: \" \\ / [ ] : | < > + = ; , * ?
	Share string `json:"share,omitempty"`
	// Synchronization state with leader
	Sync string `json:"sync,omitempty"`
	// Synchronization time with leader
	SyncTime string `json:"sync_time,omitempty"`
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
	PhysicalCapacity int32 `json:"physical_capacity,omitempty"`
	// Logical Capacity
	LogicalCapacity int32 `json:"logical_capacity,omitempty"`
	// Indicates whether the view should support simultaneous access to NFS3/NFS4/SMB protocols.
	NfsInteropFlags string `json:"nfs_interop_flags,omitempty"`
	IsRemote bool `json:"is_remote,omitempty"`
	ShareAcl *ViewShareAcl `json:"share_acl,omitempty"`
	SelectForLiveMonitoring bool `json:"select_for_live_monitoring,omitempty"`
	// QoS Policy ID
	QosPolicyId int32 `json:"qos_policy_id,omitempty"`
	// QoS Policy
	QosPolicy string `json:"qos_policy,omitempty"`
}
