package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api_latest "github.com/vast-data/terraform-provider-vastdata/codegen/latest"
	utils "github.com/vast-data/terraform-provider-vastdata/utils"
)

var resources_templates = []ResourceTemplateV2{
	ResourceTemplateV2{
		ResourceName:             "User",
		Path:                     ToStringPointer("/api/users/"),
		Model:                    api_latest.User{},
		DestFile:                 ToStringPointer("user.go"),
		IgnoreFields:             NewStringSet("Id"),
		RequiredIdentifierFields: NewStringSet("name"),
		OptionalIdentifierFields: NewStringSet(),
		ListsNamesMap:            map[string][]string{"access_keys": []string{"access_key", "enabled"}},
		Generate:                 true,
		ResponseGetByURL:         false,
		DataSourceName:           "vastdata_user",
	},
	ResourceTemplateV2{
		ResourceName:             "Group",
		Path:                     ToStringPointer("/api/groups/"),
		Model:                    api_latest.Group{},
		DestFile:                 ToStringPointer("groups.go"),
		IgnoreFields:             NewStringSet("Id"),
		RequiredIdentifierFields: NewStringSet("name", "gid"),
		OptionalIdentifierFields: NewStringSet(),
		ListsNamesMap:            map[string][]string{"s3_policies_ids": []string{"policy"}},
		Generate:                 true,
		ResponseGetByURL:         false,
		DataSourceName:           "vastdata_group",
	},
	ResourceTemplateV2{
		ResourceName:             "VipPool",
		Path:                     ToStringPointer("/api/vippools/"),
		Model:                    api_latest.VipPool{},
		DestFile:                 ToStringPointer("vippools.go"),
		IgnoreFields:             NewStringSet("Id"),
		RequiredIdentifierFields: NewStringSet("name", "subnet_cidr", "role", "ip_ranges"),
		OptionalIdentifierFields: NewStringSet(),
		ListsNamesMap:            map[string][]string{"ip_ranges": []string{"start_ip", "end_ip"}, "cnode_ids": []string{"id"}},
		Generate:                 true,
		ResponseGetByURL:         false,
		DataSourceName:           "vastdata_vip_pool",
	},
	ResourceTemplateV2{
		ResourceName:             "Tenant",
		Path:                     ToStringPointer("/api/tenants/"),
		Model:                    api_latest.Tenant{},
		DestFile:                 ToStringPointer("tenants.go"),
		IgnoreFields:             NewStringSet("Created", "SyncTime", "Id"),
		RequiredIdentifierFields: NewStringSet("name"),
		OptionalIdentifierFields: NewStringSet(),
		ListsNamesMap:            map[string][]string{"client_ip_ranges": []string{"start_ip", "end_ip"}},
		Generate:                 true,
		ResponseGetByURL:         false,
		DataSourceName:           "vastdata_tenant",
		ListFields: map[string][]FakeField{"client_ip_ranges": []FakeField{
			FakeField{Name: "start_ip", Description: "The first ip of the range"},
			FakeField{Name: "end_ip", Description: "The last ip of the range"}}},
	}, ResourceTemplateV2{
		ResourceName:             "QosPolicy",
		Path:                     ToStringPointer("/api/qospolicies/"),
		Model:                    api_latest.QosPolicy{},
		DestFile:                 ToStringPointer("qospolicies.go"),
		IgnoreFields:             NewStringSet("Created", "SyncTime", "Id"),
		RequiredIdentifierFields: NewStringSet("name"),
		OptionalIdentifierFields: NewStringSet(),
		ListsNamesMap:            map[string][]string{},
		Generate:                 true,
		DataSourceName:           "vastdata_qos_policy",
	},
	ResourceTemplateV2{
		ResourceName:             "QosDynamicLimits",
		Path:                     nil,
		Model:                    api_latest.QosDynamicLimits{},
		DestFile:                 nil,
		IgnoreFields:             NewStringSet("Id"),
		RequiredIdentifierFields: NewStringSet(),
		OptionalIdentifierFields: NewStringSet(),
		ListsNamesMap:            map[string][]string{},
		Generate:                 false,
		DataSourceName:           "",
	},
	ResourceTemplateV2{
		ResourceName:             "QosStaticLimits",
		Path:                     nil,
		Model:                    api_latest.QosStaticLimits{},
		DestFile:                 nil,
		IgnoreFields:             NewStringSet("Id"),
		RequiredIdentifierFields: NewStringSet(),
		OptionalIdentifierFields: NewStringSet(),
		ListsNamesMap:            map[string][]string{},
		Generate:                 false,
		DataSourceName:           "",
	},
	ResourceTemplateV2{
		ResourceName:             "ProtectionPolicy",
		Path:                     ToStringPointer("/api/protectionpolicies/"),
		Model:                    api_latest.ProtectionPolicy{},
		DestFile:                 ToStringPointer("protectionpolicy.go"),
		IgnoreFields:             NewStringSet("Id"),
		RequiredIdentifierFields: NewStringSet("name", "prefix", "clone_type"),
		OptionalIdentifierFields: NewStringSet(),
		ListsNamesMap:            map[string][]string{},
		Generate:                 true,
		DataSourceName:           "vastdata_protection_policy",
	},
	ResourceTemplateV2{
		ResourceName:             "ProtectionPolicySchedule",
		Path:                     nil,
		Model:                    api_latest.ProtectionPolicySchedule{},
		DestFile:                 nil,
		IgnoreFields:             NewStringSet("Id"),
		RequiredIdentifierFields: NewStringSet(),
		OptionalIdentifierFields: NewStringSet(),
		ListsNamesMap:            map[string][]string{},
		Generate:                 false,
		DataSourceName:           "",
		FieldsValidators: map[string]schema.SchemaValidateDiagFunc{"start_at": utils.ProtectionPolicyStartAt,
			"every":       utils.ProtectionPolicyTimeIntervalValidation,
			"keep_local":  utils.ProtectionPolicyTimeIntervalValidation,
			"keep_remote": utils.ProtectionPolicyTimeIntervalValidation,
		},
	},
	ResourceTemplateV2{
		ResourceName:             "Quota",
		Path:                     ToStringPointer("/api/latest/quotas/"),
		Model:                    api_latest.Quota{},
		DestFile:                 ToStringPointer("quotas.go"),
		IgnoreFields:             NewStringSet("LastUserQuotasUpdate", "Id"),
		RequiredIdentifierFields: NewStringSet("name"),
		OptionalIdentifierFields: NewStringSet("tenant_id"),
		ListsNamesMap:            map[string][]string{},
		Generate:                 true,
		ResponseGetByURL:         true,
		BeforePostFunc:           utils.EntityMergeToUserQuotas,
		BeforePatchFunc:          utils.EntityMergeToUserQuotas,
		DataSourceName:           "vastdata_quota",
		FieldsValidators:         map[string]schema.SchemaValidateDiagFunc{"grace_period": utils.GracePeriodFormatValidation},
	},
	ResourceTemplateV2{
		ResourceName:             "DefaultQuota",
		Path:                     ToStringPointer("/api/latest/quotas/"),
		Model:                    api_latest.DefaultQuota{},
		DestFile:                 ToStringPointer("quotas.go"),
		IgnoreFields:             NewStringSet("LastUserQuotasUpdate", "Id"),
		RequiredIdentifierFields: NewStringSet("name"),
		OptionalIdentifierFields: NewStringSet(),
		ListsNamesMap:            map[string][]string{},
		Generate:                 false,
		ResponseGetByURL:         false,
		DataSourceName:           "vastdata_quota",
		FieldsValidators:         map[string]schema.SchemaValidateDiagFunc{"grace_period": utils.GracePeriodFormatValidation},
	},
	ResourceTemplateV2{
		ResourceName:             "UserQuota",
		Path:                     ToStringPointer("/api/latest/quotas/"),
		Model:                    api_latest.UserQuota{},
		DestFile:                 ToStringPointer("quotas.go"),
		IgnoreFields:             NewStringSet("LastUserQuotasUpdate", "Id"),
		RequiredIdentifierFields: NewStringSet("name"),
		OptionalIdentifierFields: NewStringSet(),
		ListsNamesMap:            map[string][]string{},
		Generate:                 false,
		ResponseGetByURL:         false,
		DataSourceName:           "vastdata_quota",
		FieldsValidators:         map[string]schema.SchemaValidateDiagFunc{"grace_period": utils.GracePeriodFormatValidation},
	},
	ResourceTemplateV2{
		ResourceName:             "QuotaEntityInfo",
		Path:                     ToStringPointer("/api/latest/quotas/"),
		Model:                    api_latest.QuotaEntityInfo{},
		DestFile:                 ToStringPointer("quotas.go"),
		IgnoreFields:             NewStringSet("LastUserQuotasUpdate", "Id"),
		RequiredIdentifierFields: NewStringSet("identifier"),
		OptionalIdentifierFields: NewStringSet("name"),
		ListsNamesMap:            map[string][]string{},
		Generate:                 false,
		ResponseGetByURL:         false,
		DataSourceName:           "vastdata_quota",
		FieldsValidators:         map[string]schema.SchemaValidateDiagFunc{"grace_period": utils.GracePeriodFormatValidation},
	},
	ResourceTemplateV2{
		ResourceName:             "Dns",
		Path:                     ToStringPointer("/api/latest/dns/"),
		Model:                    api_latest.Dns{},
		DestFile:                 ToStringPointer("dns.go"),
		IgnoreFields:             NewStringSet("Id"),
		RequiredIdentifierFields: NewStringSet("name"),
		OptionalIdentifierFields: NewStringSet(),
		ListsNamesMap:            map[string][]string{"cnode_ids": []string{"id"}},
		Generate:                 true,
		ResponseGetByURL:         false,
		DataSourceName:           "vastdata_dns",
	},
	ResourceTemplateV2{
		ResourceName:             "ViewPolicy",
		Path:                     ToStringPointer("/api/viewpolicies/"),
		Model:                    api_latest.ViewPolicy{},
		DestFile:                 ToStringPointer("viewpolicies.go"),
		IgnoreFields:             NewStringSet("RemoteMapping", "ProtocolsAudit", "Created", "Id"),
		RequiredIdentifierFields: NewStringSet("name"),
		OptionalIdentifierFields: NewStringSet(),
		ListsNamesMap: map[string][]string{"nfs_read_write": []string{"address"},
			"nfs_root_squash": []string{"address"},
			"read_write":      []string{"address"},
			"s3_read_write":   []string{"address"},
			"smb_read_write":  []string{"address"}},
		Generate:         true,
		ResponseGetByURL: false,
		DataSourceName:   "vastdata_view_policy",
	},
	ResourceTemplateV2{
		ResourceName:             "View",
		Path:                     ToStringPointer("/api/views/"),
		Model:                    api_latest.View{},
		DestFile:                 ToStringPointer("views.go"),
		IgnoreFields:             NewStringSet("Id"),
		RequiredIdentifierFields: NewStringSet("path", "policy_id"),
		OptionalIdentifierFields: NewStringSet("tenant_id"),
		ListsNamesMap:            map[string][]string{},
		Generate:                 true,
		ResponseGetByURL:         false,
		DataSourceName:           "vastdata_view",
		//		IgnoreUpdates:            NewStringSet("create_dir"),
		BeforePatchFunc:  utils.AlwaysSendCreateDir,
		BeforeCreateFunc: utils.AlwaysStoreCreateDir,
		AfterPatchFunc:   utils.AlwaysStoreCreateDir,
		AfterReadFunc:    utils.KeepCreateDirState,
	},
	ResourceTemplateV2{
		ResourceName:             "ViewShareAcl",
		Path:                     ToStringPointer("/api/views/"),
		Model:                    api_latest.ViewShareAcl{},
		DestFile:                 ToStringPointer("views.go"),
		IgnoreFields:             NewStringSet("Id"),
		RequiredIdentifierFields: NewStringSet("name"),
		OptionalIdentifierFields: NewStringSet(),
		ListsNamesMap:            map[string][]string{},
		Generate:                 false,
		ResponseGetByURL:         false,
		DataSourceName:           "",
	},
	ResourceTemplateV2{
		ResourceName:             "ShareAcl",
		Path:                     ToStringPointer("/api/views/"),
		Model:                    api_latest.ShareAcl{},
		DestFile:                 ToStringPointer("views.go"),
		IgnoreFields:             NewStringSet(),
		RequiredIdentifierFields: NewStringSet("name"),
		OptionalIdentifierFields: NewStringSet(),
		ListsNamesMap:            map[string][]string{},
		Generate:                 false,
		ResponseGetByURL:         false,
		DataSourceName:           "",
	},

	ResourceTemplateV2{
		ResourceName:             "Nis",
		Path:                     ToStringPointer("/api/nis/"),
		Model:                    api_latest.Nis{},
		DestFile:                 ToStringPointer("nis.go"),
		IgnoreFields:             NewStringSet("Id"),
		RequiredIdentifierFields: NewStringSet("domain_name"),
		OptionalIdentifierFields: NewStringSet(),
		ListsNamesMap:            map[string][]string{"hosts": []string{"host"}},
		Generate:                 true,
		ResponseGetByURL:         false,
		DataSourceName:           "vastdata_nis",
	},
	ResourceTemplateV2{
		ResourceName:             "Ldap",
		Path:                     ToStringPointer("/api/ldaps/"),
		Model:                    api_latest.Ldap{},
		DestFile:                 ToStringPointer("ldaps.go"),
		IgnoreFields:             NewStringSet("Id"),
		RequiredIdentifierFields: NewStringSet("domain_name"),
		OptionalIdentifierFields: NewStringSet(),
		ListsNamesMap:            map[string][]string{"urls": []string{"url"}},
		SensitiveFields:          NewStringSet("bindpw"),
		Generate:                 true,
		IgnoreUpdates:            NewStringSet("bindpw"),
		DataSourceName:           "vastdata_ldap",
	},
	ResourceTemplateV2{
		ResourceName:             "S3LifeCycleRule",
		Path:                     ToStringPointer("/api/s3lifecyclerules/"),
		Model:                    api_latest.S3LifeCycleRule{},
		DestFile:                 ToStringPointer("s3lifecyclerules.go"),
		IgnoreFields:             NewStringSet("Id", "expiration_date", "view_path"),
		RequiredIdentifierFields: NewStringSet("name", "prefix"),
		OptionalIdentifierFields: NewStringSet(),
		ListsNamesMap:            map[string][]string{},
		Generate:                 true,
		DataSourceName:           "vastdata_s3_life_cycle_rule",
		ResponseProcessingFunc:   "ProcessingResultsListResponse",
		BeforePatchFunc:          utils.EnabledMustBeSet,
		BeforePostFunc:           utils.EnabledMustBeSet,
	},
	ResourceTemplateV2{
		ResourceName:             "ActiveDirectory",
		Path:                     ToStringPointer("/api/activedirectory/"),
		Model:                    api_latest.ActiveDirectory{},
		DestFile:                 ToStringPointer("activedirectory.go"),
		IgnoreFields:             NewStringSet("Id"),
		RequiredIdentifierFields: NewStringSet("machine_account_name", "ldap_id"),
		OptionalIdentifierFields: NewStringSet(),
		ListsNamesMap:            map[string][]string{},
		Generate:                 true,
		BeforeDeleteFunc:         utils.AlwaysSkipDeleteLdap,
		DataSourceName:           "vastdata_active_directory",
	},
	ResourceTemplateV2{
		ResourceName:             "S3Policy",
		Path:                     ToStringPointer("/api/s3policies/"),
		Model:                    api_latest.S3Policy{},
		DestFile:                 ToStringPointer("s3userpolicies.go"),
		IgnoreFields:             NewStringSet("Id", "Users", "Groups", "IsReplicated"),
		RequiredIdentifierFields: NewStringSet("name", "policy"),
		OptionalIdentifierFields: NewStringSet(),
		ListsNamesMap:            map[string][]string{"users": []string{"user"}, "groups": []string{"group"}},
		Generate:                 true,
		DataSourceName:           "vastdata_s3_policy",
		BeforePostFunc:           utils.EnabledMustBeSet,
		BeforePatchFunc:          utils.EnabledMustBeSet,
	},
	ResourceTemplateV2{
		ResourceName:             "ProtectedPath",
		Path:                     ToStringPointer("/api/protectedpaths/"),
		Model:                    api_latest.ProtectedPath{},
		DestFile:                 ToStringPointer("protectedpaths.go"),
		IgnoreFields:             NewStringSet("Id"),
		RequiredIdentifierFields: NewStringSet("name"),
		OptionalIdentifierFields: NewStringSet(),
		ListsNamesMap:            map[string][]string{},
		Generate:                 true,
		DataSourceName:           "vastdata_protected_path",
	},
	ResourceTemplateV2{
		ResourceName:             "Snapshot",
		Path:                     ToStringPointer("/api/snapshots/"),
		Model:                    api_latest.Snapshot{},
		DestFile:                 ToStringPointer("snapshots.go"),
		IgnoreFields:             NewStringSet("Id"),
		RequiredIdentifierFields: NewStringSet("name"),
		OptionalIdentifierFields: NewStringSet("tenant_id"),
		ListsNamesMap:            map[string][]string{},
		Generate:                 true,
		DataSourceName:           "vastdata_snapshot",
		FieldsValidators: map[string]schema.SchemaValidateDiagFunc{
			"expiration_time": utils.SnapshotExpirationFormatValidation,
		},
	},
	ResourceTemplateV2{
		ResourceName:             "GlobalSnapshot",
		Path:                     ToStringPointer("/api/globalsnapstreams/"),
		Model:                    api_latest.GlobalSnapshot{},
		DestFile:                 ToStringPointer("globalsnapshots.go"),
		IgnoreFields:             NewStringSet("Id"),
		RequiredIdentifierFields: NewStringSet("name", "loanee_root_path", "remote_target_id", "remote_target_guid", "owner_tenant", "owner_root_snapshot", "remote_target_path", "loanee_tenant_id"),
		OptionalIdentifierFields: NewStringSet(),
		ListsNamesMap:            map[string][]string{},
		Generate:                 true,
		DataSourceName:           "vastdata_global_snapshot",
		BeforePostFunc:           utils.AddStreamInfo,
		BeforePatchFunc:          utils.UpdateStreamInfo,
		IgnoreUpdates:            NewStringSet("loanee_tenant_id"),
	},
	ResourceTemplateV2{
		ResourceName:             "GlobalLocalSnapshot",
		Path:                     ToStringPointer("/api/globalsnapstreams/"),
		Model:                    api_latest.GlobalLocalSnapshot{},
		DestFile:                 ToStringPointer("globallocalsnapshots.go"),
		IgnoreFields:             NewStringSet("Id"),
		RequiredIdentifierFields: NewStringSet("name", "loanee_root_path", "owner_tenant", "loanee_snapshot_id", "loanee_tenant_id"),
		OptionalIdentifierFields: NewStringSet(),
		ListsNamesMap:            map[string][]string{},
		Generate:                 true,
		DataSourceName:           "vastdata_global_local_snapshot",
	},

	ResourceTemplateV2{
		ResourceName:             "GlobalSnapshotOwnerRootSnapshot",
		Path:                     ToStringPointer("/api/globalsnapstreams/"),
		Model:                    api_latest.GlobalSnapshotOwnerRootSnapshot{},
		DestFile:                 ToStringPointer("globalsnapshots.go"),
		IgnoreFields:             NewStringSet(),
		RequiredIdentifierFields: NewStringSet("name"),
		OptionalIdentifierFields: NewStringSet(),
		ComputedFields:           NewStringSet("clone_id", "parent_handle_ehandle"),
		ListsNamesMap:            map[string][]string{},
		Generate:                 false,
		DataSourceName:           "",
	},
	ResourceTemplateV2{
		ResourceName:             "GlobalSnapshotOwnerTenant",
		Path:                     ToStringPointer("/api/globalsnapstreams/"),
		Model:                    api_latest.GlobalSnapshotOwnerTenant{},
		DestFile:                 ToStringPointer("globalsnapshots.go"),
		IgnoreFields:             NewStringSet(),
		RequiredIdentifierFields: NewStringSet("name", "guid"),
		OptionalIdentifierFields: NewStringSet(),
		ListsNamesMap:            map[string][]string{},
		Generate:                 false,
		DataSourceName:           "",
	},

	ResourceTemplateV2{
		ResourceName:             "ReplicationPeers",
		Path:                     ToStringPointer("/api/nativereplicationremotetargets/"),
		Model:                    api_latest.ReplicationPeers{},
		DestFile:                 ToStringPointer("replicationpeers.go"),
		IgnoreFields:             NewStringSet("Id"),
		RequiredIdentifierFields: NewStringSet("name"),
		OptionalIdentifierFields: NewStringSet(),
		ListsNamesMap:            map[string][]string{},
		Generate:                 true,
		DataSourceName:           "vastdata_replication_peers",
	},
	ResourceTemplateV2{
		ResourceName:             "S3replicationPeers",
		Path:                     ToStringPointer("/api/replicationtargets/"),
		Model:                    api_latest.S3replicationPeers{},
		DestFile:                 ToStringPointer("s3replicationpeers.go"),
		IgnoreFields:             NewStringSet("Id"),
		RequiredIdentifierFields: NewStringSet("name"),
		OptionalIdentifierFields: NewStringSet(),
		ListsNamesMap:            map[string][]string{},
		Generate:                 true,
		DataSourceName:           "vastdata_s3_replication_peers",
		SensitiveFields:          NewStringSet("secret_key", "access_key"),
		IgnoreUpdates:            NewStringSet("secret_key", "access_key"),
	},
}
