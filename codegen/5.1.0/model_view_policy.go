/*
 * VastData API
 *
 * A API document representing VastData API
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger

type ViewPolicy struct {
	// A uniqe guid given to the view policy
	Id int32 `json:"id,omitempty"`
	// A uniqe guid given to the view policy
	Guid string `json:"guid,omitempty"`
	// A uniqe name given to the view policy.                         
	Name string `json:"name"`
	// Determine the way a file inherits GID
	GidInheritance string `json:"gid_inheritance,omitempty"`
	// Security flavor, which determines how file and directory permissions are applied in multiprotocol views.
	Flavor string `json:"flavor,omitempty"`
	// Applicable with MIXED_LAST_WINS security flavor (Access can be set via NFSv3 regardless of this option)
	AccessFlavor string `json:"access_flavor,omitempty"`
	// How to determine the maximum allowed path length
	PathLength string `json:"path_length,omitempty"`
	// How to determine the allowed characters in a path
	AllowedCharacters string `json:"allowed_characters,omitempty"`
	Use32bitFileid bool `json:"use_32bit_fileid,omitempty"`
	ExposeIdInFsid bool `json:"expose_id_in_fsid,omitempty"`
	// Use configured Auth Provider(s) to enforce group permissions when set to true , if set to ture with out specifing auth_source , the auth_source set to \"PROVIDERS\". if set to false than auth_source set to RPC. Due to the nature or terrafrom simply changing use_auth_provider from false to true or the other way around will not change the value auth_source as terrafrom will keep hold on the previous value. therefor it is adviasable to always specify the value of auth_source
	UseAuthProvider bool `json:"use_auth_provider,omitempty"`
	// The source of authentication
	AuthSource string `json:"auth_source,omitempty"`
	// Hosts with NFS read/write permissions
	ReadWrite []string `json:"read_write,omitempty"`
	// Hosts with NFS read only permissions
	ReadOnly []string `json:"read_only,omitempty"`
	// Hosts with NFS read/write permissions. when creating a new View Policy if the value is not set than an empty list is sent to the VastData cluster resulting in empty list of addresses However during update if nfs_all_squash is removed from the resource nothing is changed to preserve terraform default behaviour in such cases. If there is a need to change the value an empty list it must be secifed and set to [].
	NfsReadWrite []string `json:"nfs_read_write,omitempty"`
	// Hosts with NFS read only permissions. when creating a new View Policy if the value is not set than an empty list is sent to the VastData cluster resulting in empty list of addresses However during update if nfs_all_squash is removed from the resource nothing is changed to preserve terraform default behaviour in such cases. If there is a need to change the value an empty list it must be secifed and set to [].
	NfsReadOnly []string `json:"nfs_read_only,omitempty"`
	// Hosts with SMB read/write permissions. when creating a new View Policy if the value is not set than an empty list is sent to the VastData cluster resulting in empty list of addresses However during update if nfs_all_squash is removed from the resource nothing is changed to preserve terraform default behaviour in such cases. If there is a need to change the value an empty list it must be secifed and set to [].
	SmbReadWrite []string `json:"smb_read_write,omitempty"`
	// Hosts with SMB read only permissions. when creating a new View Policy if the value is not set than an empty list is sent to the VastData cluster resulting in empty list of addresses However during update if nfs_all_squash is removed from the resource nothing is changed to preserve terraform default behaviour in such cases. If there is a need to change the value an empty list it must be secifed and set to [].
	SmbReadOnly []string `json:"smb_read_only,omitempty"`
	// Hosts with S3 read/write permissions. when creating a new View Policy if the value is not set than an empty list is sent to the VastData cluster resulting in empty list of addresses However during update if nfs_all_squash is removed from the resource nothing is changed to preserve terraform default behaviour in such cases. If there is a need to change the value an empty list it must be secifed and set to [].
	S3ReadWrite []string `json:"s3_read_write,omitempty"`
	// Hosts with S3 read only permissions. when creating a new View Policy if the value is not set than an empty list is sent to the VastData cluster resulting in empty list of addresses However during update if nfs_all_squash is removed from the resource nothing is changed to preserve terraform default behaviour in such cases. If there is a need to change the value an empty list it must be secifed and set to [].
	S3ReadOnly []string `json:"s3_read_only,omitempty"`
	// Hosts with trash permissions
	TrashAccess []string `json:"trash_access,omitempty"`
	// Enable POSIX ACL
	NfsPosixAcl bool `json:"nfs_posix_acl,omitempty"`
	// when using smb use open permissions for files
	NfsReturnOpenPermissions bool `json:"nfs_return_open_permissions,omitempty"`
	// Hosts with no squash policy
	NfsNoSquash []string `json:"nfs_no_squash,omitempty"`
	// Hosts with root squash policy. when creating a new View Policy if the value is not set than an empty list is sent to the VastData cluster resulting in empty list of addresses However during update if nfs_all_squash is removed from the resource nothing is changed to preserve terraform default behaviour in such cases. If there is a need to change the value an empty list it must be secifed and set to [].
	NfsRootSquash []string `json:"nfs_root_squash,omitempty"`
	// Hosts with all squash policy. when creating a new View Policy if the value is not set than an empty list is sent to the VastData cluster resulting in empty list of addresses However during update if nfs_all_squash is removed from the resource nothing is changed to preserve terraform default behaviour in such cases. If there is a need to change the value an empty list it must be secifed and set to []
	NfsAllSquash []string `json:"nfs_all_squash,omitempty"`
	// Hosts with full permissions
	S3BucketFullControl string `json:"s3_bucket_full_control,omitempty"`
	// Hosts with full permissions
	S3BucketListing string `json:"s3_bucket_listing,omitempty"`
	// Hosts with full permissions
	S3BucketRead string `json:"s3_bucket_read,omitempty"`
	// Hosts with full permissions
	S3BucketReadAcp string `json:"s3_bucket_read_acp,omitempty"`
	// Hosts with full permissions
	S3BucketWrite string `json:"s3_bucket_write,omitempty"`
	// Hosts with full permissions
	S3BucketWriteAcp string `json:"s3_bucket_write_acp,omitempty"`
	// Hosts with full permissions
	S3ObjectFullControl string `json:"s3_object_full_control,omitempty"`
	// Hosts with full permissions
	S3ObjectRead string `json:"s3_object_read,omitempty"`
	// Hosts with full permissions
	S3ObjectReadAcp string `json:"s3_object_read_acp,omitempty"`
	// Hosts with full permissions
	S3ObjectWrite string `json:"s3_object_write,omitempty"`
	// Hosts with full permissions
	S3ObjectWriteAcp string `json:"s3_object_write_acp,omitempty"`
	// Default unix type permissions on new file
	SmbFileMode int32 `json:"smb_file_mode,omitempty"`
	// Default unix type permissions on new folder
	SmbDirectoryMode int32 `json:"smb_directory_mode,omitempty"`
	// Default unix type permissions on new file
	SmbFileModePadded string `json:"smb_file_mode_padded,omitempty"`
	// Default unix type permissions on new folder
	SmbDirectoryModePadded string `json:"smb_directory_mode_padded,omitempty"`
	// Parent Cluster
	Cluster string `json:"cluster,omitempty"`
	// Parent Cluster ID
	ClusterId int32 `json:"cluster_id,omitempty"`
	// Tenant ID
	TenantId int32 `json:"tenant_id,omitempty"`
	// Tenant Name
	TenantName string `json:"tenant_name,omitempty"`
	Url string `json:"url,omitempty"`
	// Frequency for updating the atime attribute of NFS files. atime is updated on read operations if the difference between the current time and the file's atime value is greater than the atime frequency. Specify as time in seconds.
	AtimeFrequency string `json:"atime_frequency,omitempty"`
	// Comma separated vip pool ids.
	VipPools []int32 `json:"vip_pools,omitempty"`
	// NFS 4.1 minimal protection level
	NfsMinimalProtectionLevel string `json:"nfs_minimal_protection_level,omitempty"`
	// A list of usernames for bucket listing permissions
	S3Visibility []string `json:"s3_visibility,omitempty"`
	// A list of group names for bucket listing permissions
	S3VisibilityGroups []string `json:"s3_visibility_groups,omitempty"`
	// Apple sid
	AppleSid bool `json:"apple_sid,omitempty"`
	// Map of protocols audit configurations
	ProtocolsAudit *interface{} `json:"protocols_audit,omitempty"`
	// Protocols to audit
	Protocols []string `json:"protocols,omitempty"`
	// Create/Delete Files/Directories/Objects
	DataCreateDelete bool `json:"data_create_delete,omitempty"`
	// Modify data/MD
	DataModify bool `json:"data_modify,omitempty"`
	// Read data
	DataRead bool `json:"data_read,omitempty"`
	// Log full path
	LogFullPath bool `json:"log_full_path,omitempty"`
	// Log hostname
	LogHostname bool `json:"log_hostname,omitempty"`
	// Log username
	LogUsername bool `json:"log_username,omitempty"`
	// Log deleted files/dirs from trash dir
	LogDeleted bool `json:"log_deleted,omitempty"`
	RemoteMapping *interface{} `json:"remote_mapping,omitempty"`
	// Number of Policy related Views
	CountViews int32 `json:"count_views,omitempty"`
	// Specifies whether to make the .snapshot directory accessible in subdirectories of the View.
	EnableSnapshotLookup bool `json:"enable_snapshot_lookup,omitempty"`
	// Specifies whether to make the .snapshot directory visible in subdirectories of the View.
	EnableListingOfSnapshotDir bool `json:"enable_listing_of_snapshot_dir,omitempty"`
}
