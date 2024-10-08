/*
 * VastData API
 *
 * A API document representing VastData API
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger

type Tenant struct {
	// A uniq id given to the tenant
	Id int32 `json:"id,omitempty"`
	// A uniq guid given to the tenant
	Guid string `json:"guid,omitempty"`
	// A uniq name given to the tenant
	Name string `json:"name,omitempty"`
	// Enables SMB privileged user
	UseSmbPrivilegedUser bool `json:"use_smb_privileged_user,omitempty"`
	// Optional custom username for the SMB privileged user. If not set, the SMB privileged user name is 'vastadmin'
	SmbPrivilegedUserName string `json:"smb_privileged_user_name,omitempty"`
	// Enables SMB privileged user group
	UseSmbPrivilegedGroup bool `json:"use_smb_privileged_group,omitempty"`
	// Optional custom SID to specify a non default SMB privileged group. If not set, SMB privileged group is the Backup Operators domain group.
	SmbPrivilegedGroupSid string `json:"smb_privileged_group_sid,omitempty"`
	// True=The SMB privileged user group has read and write control access. Members of the group can perform backup and restore operations on all files and directories, without requiring read or write access to the specific files and directories. False=the privileged group has read only access.
	SmbPrivilegedGroupFullAccess bool `json:"smb_privileged_group_full_access,omitempty"`
	// Optional custom name to specify a non default privileged group. If not set, privileged group is the Backup Operators domain group.
	SmbAdministratorsGroupName string `json:"smb_administrators_group_name,omitempty"`
	// Default Share-level permissions for Others
	DefaultOthersShareLevelPerm string `json:"default_others_share_level_perm,omitempty"`
	// GID with permissions to the trash folder
	TrashGid int32 `json:"trash_gid,omitempty"`
	// Array of source IP ranges to allow for the tenant.
	ClientIpRanges [][]string `json:"client_ip_ranges,omitempty"`
	// POSIX primary provider type
	PosixPrimaryProvider string `json:"posix_primary_provider,omitempty"`
	// AD provider ID
	AdProviderId int32 `json:"ad_provider_id,omitempty"`
	// Open-LDAP provider ID specified separately by the user
	LdapProviderId int32 `json:"ldap_provider_id,omitempty"`
	// NIS provider ID
	NisProviderId int32 `json:"nis_provider_id,omitempty"`
	// Tenant's encryption group unique identifier
	EncryptionCrn string `json:"encryption_crn,omitempty"`
	// Enable NFSv4.2
	IsNfsv42Supported bool `json:"is_nfsv42_supported,omitempty"`
	// Allow IO from users whose Active Directory accounts are locked out by lockout policies due to unsuccessful login attempts.
	AllowLockedUsers bool `json:"allow_locked_users,omitempty"`
	// Allow IO from users whose Active Directory accounts are explicitly disabled.
	AllowDisabledUsers bool `json:"allow_disabled_users,omitempty"`
	// Use native SMB authentication
	UseSmbNative bool `json:"use_smb_native,omitempty"`
	// An array of VIP Pool names attached to this tenant.
	VippoolNames []string `json:"vippool_names,omitempty"`
	// An array of VIP Pool ids to attach to tenant.
	VippoolIds []int64 `json:"vippool_ids,omitempty"`
}
