package configs

import (
	api_latest "github.com/vast-data/terraform-provider-vastdata/codegen/latest"
	utils "github.com/vast-data/terraform-provider-vastdata/utils"
)

var datasources_map map[string]ResourceTemplateV2 = map[string]ResourceTemplateV2{}

var DatasourcesTemplates = []ResourceTemplateV2{
	ResourceTemplateV2{
		ResourceName:             "Cnode",
		Path:                     ToStringPointer("cnodes"),
		Model:                    api_latest.Cnode{},
		DestFile:                 ToStringPointer("cnodes.go"),
		IgnoreFields:             NewStringSet(),
		RequiredIdentifierFields: NewStringSet("name"),
		OptionalIdentifierFields: NewStringSet(),
		ListsNamesMap:            map[string][]string{},
		Generate:                 true,
		DataSourceName:           "vastdata_cnode",
	},
	ResourceTemplateV2{
		ResourceName:             "QosPolicy",
		Path:                     ToStringPointer("qospolicies"),
		Model:                    api_latest.QosPolicy{},
		DestFile:                 ToStringPointer("qospolicies.go"),
		IgnoreFields:             NewStringSet("Created", "SyncTime"),
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
		IgnoreFields:             NewStringSet(),
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
		IgnoreFields:             NewStringSet(),
		RequiredIdentifierFields: NewStringSet(),
		OptionalIdentifierFields: NewStringSet(),
		ListsNamesMap:            map[string][]string{},
		Generate:                 false,
		DataSourceName:           "",
	},
	ResourceTemplateV2{
		ResourceName:             "Quota",
		Path:                     ToStringPointer("quotas"),
		Model:                    api_latest.Quota{},
		DestFile:                 ToStringPointer("quotas.go"),
		IgnoreFields:             NewStringSet("LastUserQuotasUpdate"),
		RequiredIdentifierFields: NewStringSet("name"),
		OptionalIdentifierFields: NewStringSet("tenant_id"),
		ListsNamesMap:            map[string][]string{},
		Generate:                 true,
		ResponseGetByURL:         true,
		DataSourceName:           "vastdata_quota",
	},
	ResourceTemplateV2{
		ResourceName:             "DefaultQuota",
		Path:                     ToStringPointer("quotas"),
		Model:                    api_latest.DefaultQuota{},
		DestFile:                 ToStringPointer("quotas.go"),
		IgnoreFields:             NewStringSet("LastUserQuotasUpdate"),
		RequiredIdentifierFields: NewStringSet("name"),
		OptionalIdentifierFields: NewStringSet(),
		ListsNamesMap:            map[string][]string{},
		Generate:                 false,
		ResponseGetByURL:         false,
		DataSourceName:           "vastdata_quota",
	},
	ResourceTemplateV2{
		ResourceName:             "UserQuota",
		Path:                     ToStringPointer("quotas"),
		Model:                    api_latest.UserQuota{},
		DestFile:                 ToStringPointer("quotas.go"),
		IgnoreFields:             NewStringSet("LastUserQuotasUpdate"),
		RequiredIdentifierFields: NewStringSet("name"),
		OptionalIdentifierFields: NewStringSet(),
		ListsNamesMap:            map[string][]string{},
		Generate:                 false,
		ResponseGetByURL:         false,
		DataSourceName:           "vastdata_quota",
	},
	ResourceTemplateV2{
		ResourceName:             "QuotaEntityInfo",
		Path:                     ToStringPointer("quotas"),
		Model:                    api_latest.QuotaEntityInfo{},
		DestFile:                 ToStringPointer("quotas.go"),
		IgnoreFields:             NewStringSet("LastUserQuotasUpdate"),
		RequiredIdentifierFields: NewStringSet("name"),
		OptionalIdentifierFields: NewStringSet(),
		ListsNamesMap:            map[string][]string{},
		Generate:                 false,
		ResponseGetByURL:         false,
		DataSourceName:           "vastdata_quota",
	},
	ResourceTemplateV2{
		ResourceName:             "Dns",
		Path:                     ToStringPointer("dns"),
		Model:                    api_latest.Dns{},
		DestFile:                 ToStringPointer("dns.go"),
		IgnoreFields:             NewStringSet(),
		RequiredIdentifierFields: NewStringSet("name"),
		OptionalIdentifierFields: NewStringSet(),
		ListsNamesMap:            map[string][]string{"cnode_ids": []string{"id"}},
		Generate:                 true,
		ResponseGetByURL:         false,
		DataSourceName:           "vastdata_dns",
	},
	ResourceTemplateV2{
		ResourceName:             "VipPool",
		Path:                     ToStringPointer("vippools"),
		Model:                    api_latest.VipPool{},
		DestFile:                 ToStringPointer("vippools.go"),
		IgnoreFields:             NewStringSet(),
		RequiredIdentifierFields: NewStringSet("name"),
		OptionalIdentifierFields: NewStringSet(),
		ListsNamesMap:            map[string][]string{"ip_ranges": []string{"start_ip", "end_ip"}, "cnode_ids": []string{"id"}},
		Generate:                 true,
		ResponseGetByURL:         false,
		DataSourceName:           "vastdata_vip_pool",
	},
	ResourceTemplateV2{
		ResourceName:             "ViewPolicy",
		Path:                     ToStringPointer("viewpolicies"),
		Model:                    api_latest.ViewPolicy{},
		DestFile:                 ToStringPointer("viewpolicies.go"),
		IgnoreFields:             NewStringSet("RemoteMapping", "ProtocolsAudit", "Created"),
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
		ResourceName:             "Group",
		Path:                     ToStringPointer("groups"),
		Model:                    api_latest.Group{},
		DestFile:                 ToStringPointer("groups.go"),
		IgnoreFields:             NewStringSet(),
		RequiredIdentifierFields: NewStringSet("name"),
		OptionalIdentifierFields: NewStringSet(),
		ListsNamesMap:            map[string][]string{"s3_policies_ids": []string{"policy"}},
		Generate:                 true,
		ResponseGetByURL:         false,
		DataSourceName:           "vastdata_group",
	},
	ResourceTemplateV2{
		ResourceName:             "User",
		Path:                     ToStringPointer("users"),
		Model:                    api_latest.User{},
		DestFile:                 ToStringPointer("user.go"),
		IgnoreFields:             NewStringSet(),
		RequiredIdentifierFields: NewStringSet("name"),
		OptionalIdentifierFields: NewStringSet(),
		ListsNamesMap:            map[string][]string{"access_keys": []string{"access_key", "enabled"}},
		Generate:                 true,
		ResponseGetByURL:         false,
		DataSourceName:           "vastdata_user",
	},
	ResourceTemplateV2{
		ResourceName:             "View",
		Path:                     ToStringPointer("views"),
		Model:                    api_latest.View{},
		DestFile:                 ToStringPointer("views.go"),
		IgnoreFields:             NewStringSet(),
		RequiredIdentifierFields: NewStringSet("path"),
		OptionalIdentifierFields: NewStringSet("tenant_id"),
		ListsNamesMap:            map[string][]string{},
		Generate:                 true,
		ResponseGetByURL:         false,
		DataSourceName:           "vastdata_view",
	},
	ResourceTemplateV2{
		ResourceName:             "ViewShareAcl",
		Path:                     ToStringPointer("views"),
		Model:                    api_latest.ViewShareAcl{},
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
		ResourceName:             "ShareAcl",
		Path:                     ToStringPointer("views"),
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
		Path:                     ToStringPointer("nis"),
		Model:                    api_latest.Nis{},
		DestFile:                 ToStringPointer("nis.go"),
		IgnoreFields:             NewStringSet(),
		RequiredIdentifierFields: NewStringSet("domain_name"),
		OptionalIdentifierFields: NewStringSet(),
		ListsNamesMap:            map[string][]string{"hosts": []string{"host"}},
		Generate:                 true,
		ResponseGetByURL:         false,
		DataSourceName:           "vastdata_nis",
	},
	ResourceTemplateV2{
		ResourceName:             "Tenant",
		Path:                     ToStringPointer("tenants"),
		Model:                    api_latest.Tenant{},
		DestFile:                 ToStringPointer("tenants.go"),
		IgnoreFields:             NewStringSet("Created", "SyncTime"),
		RequiredIdentifierFields: NewStringSet("name"),
		OptionalIdentifierFields: NewStringSet(),
		ListsNamesMap:            map[string][]string{"client_ip_ranges": []string{"start_ip", "end_ip"}},
		Generate:                 true,
		ResponseGetByURL:         false,
		DataSourceName:           "vastdata_tenant",
		ListFields: map[string][]FakeField{"client_ip_ranges": []FakeField{
			FakeField{Name: "start_ip", Description: "The first ip of the range"},
			FakeField{Name: "end_ip", Description: "The last ip of the range"}}},
	},

	ResourceTemplateV2{
		ResourceName:             "Ldap",
		Path:                     ToStringPointer("ldaps"),
		Model:                    api_latest.Ldap{},
		DestFile:                 ToStringPointer("ldaps.go"),
		IgnoreFields:             NewStringSet(),
		RequiredIdentifierFields: NewStringSet("domain_name"),
		OptionalIdentifierFields: NewStringSet(),
		ListsNamesMap:            map[string][]string{"urls": []string{"url"}},
		Generate:                 true,
		DataSourceName:           "vastdata_ldap",
	},
	ResourceTemplateV2{
		ResourceName:             "S3LifeCycleRule",
		Path:                     ToStringPointer("s3lifecyclerules"),
		Model:                    api_latest.S3LifeCycleRule{},
		DestFile:                 ToStringPointer("s3lifecyclerules.go"),
		IgnoreFields:             NewStringSet(),
		RequiredIdentifierFields: NewStringSet("name"),
		OptionalIdentifierFields: NewStringSet(),
		ListsNamesMap:            map[string][]string{},
		Generate:                 true,
		DataSourceName:           "vastdata_s3_life_cycle_rule",
		ResponseProcessingFunc:   utils.ProcessingResultsListResponse,
	},
	ResourceTemplateV2{
		ResourceName:             "ActiveDirectory",
		Path:                     ToStringPointer("activedirectory"),
		Model:                    api_latest.ActiveDirectory{},
		DestFile:                 ToStringPointer("activedirectory.go"),
		IgnoreFields:             NewStringSet(),
		RequiredIdentifierFields: NewStringSet("machine_account_name"),
		OptionalIdentifierFields: NewStringSet(),
		ListsNamesMap:            map[string][]string{"preferred_dc_list": []string{"dc"}},
		Generate:                 true,
		DataSourceName:           "vastdata_active_directory",
	},
	ResourceTemplateV2{
		ResourceName:             "S3Policy",
		Path:                     ToStringPointer("s3userpolicies"),
		Model:                    api_latest.S3Policy{},
		DestFile:                 ToStringPointer("s3userpolicies.go"),
		IgnoreFields:             NewStringSet(),
		RequiredIdentifierFields: NewStringSet("name"),
		OptionalIdentifierFields: NewStringSet(),
		ListsNamesMap:            map[string][]string{"users": []string{"user"}, "groups": []string{"group"}},
		Generate:                 true,
		DataSourceName:           "vastdata_s3_policy",
	},
	ResourceTemplateV2{
		ResourceName:             "ProtectedPath",
		Path:                     ToStringPointer("protectedpaths"),
		Model:                    api_latest.ProtectedPath{},
		DestFile:                 ToStringPointer("protectedpaths.go"),
		IgnoreFields:             NewStringSet(),
		RequiredIdentifierFields: NewStringSet("name"),
		OptionalIdentifierFields: NewStringSet(),
		ListsNamesMap:            map[string][]string{},
		Generate:                 true,
		DataSourceName:           "vastdata_protected_path",
	},
	ResourceTemplateV2{
		ResourceName:             "Snapshot",
		Path:                     ToStringPointer("snapshots"),
		Model:                    api_latest.Snapshot{},
		DestFile:                 ToStringPointer("snapshots.go"),
		IgnoreFields:             NewStringSet(),
		RequiredIdentifierFields: NewStringSet("name"),
		OptionalIdentifierFields: NewStringSet("tenant_id"),
		ListsNamesMap:            map[string][]string{},
		Generate:                 true,
		DataSourceName:           "vastdata_snapshot",
	},
	ResourceTemplateV2{
		ResourceName:             "GlobalSnapshot",
		Path:                     ToStringPointer("globalsnapstreams"),
		Model:                    api_latest.GlobalSnapshot{},
		DestFile:                 ToStringPointer("globalsnapshots.go"),
		IgnoreFields:             NewStringSet(),
		RequiredIdentifierFields: NewStringSet("name"),
		OptionalIdentifierFields: NewStringSet(),
		ListsNamesMap:            map[string][]string{},
		Generate:                 true,
		DataSourceName:           "vastdata_global_snapshot",
	},
	ResourceTemplateV2{
		ResourceName:             "ReplicationPeers",
		Path:                     ToStringPointer("nativereplicationremotetargets"),
		Model:                    api_latest.ReplicationPeers{},
		DestFile:                 ToStringPointer("replicationpeers.go"),
		IgnoreFields:             NewStringSet(),
		RequiredIdentifierFields: NewStringSet("name"),
		OptionalIdentifierFields: NewStringSet(),
		ListsNamesMap:            map[string][]string{},
		Generate:                 true,
		DataSourceName:           "vastdata_replication_peers",
	},
	ResourceTemplateV2{
		ResourceName:             "ProtectionPolicy",
		Path:                     ToStringPointer("protectionpolicy"),
		Model:                    api_latest.ProtectionPolicy{},
		DestFile:                 ToStringPointer("protectionpolicy.go"),
		IgnoreFields:             NewStringSet(),
		RequiredIdentifierFields: NewStringSet("name"),
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
		IgnoreFields:             NewStringSet(),
		RequiredIdentifierFields: NewStringSet(),
		OptionalIdentifierFields: NewStringSet(),
		ListsNamesMap:            map[string][]string{},
		Generate:                 false,
		DataSourceName:           "",
	},
	ResourceTemplateV2{
		ResourceName:             "S3replicationPeers",
		Path:                     ToStringPointer("replicationtargets"),
		Model:                    api_latest.S3replicationPeers{},
		DestFile:                 ToStringPointer("s3replicationpeers.go"),
		IgnoreFields:             NewStringSet(),
		RequiredIdentifierFields: NewStringSet("name"),
		OptionalIdentifierFields: NewStringSet(),
		ListsNamesMap:            map[string][]string{},
		Generate:                 true,
		DataSourceName:           "vastdata_s3_replication_peers",
	},
	ResourceTemplateV2{
		ResourceName:             "ActiveDirectory2",
		Path:                     ToStringPointer("activedirectory"),
		Model:                    api_latest.ActiveDirectory2{},
		DestFile:                 ToStringPointer("activedirectory2.go"),
		IgnoreFields:             NewStringSet("Id"),
		RequiredIdentifierFields: NewStringSet("machine_account_name"),
		GetFunc:                  utils.ActiveDirectory2GetFunc,
		OptionalIdentifierFields: NewStringSet(),
		ListsNamesMap:            map[string][]string{},
		Generate:                 true,
		IgnoreUpdates:            NewStringSet("bindpw"),
		DataSourceName:           "vastdata_active_directory2",
	},
}

func init() {
	for _, r := range DatasourcesTemplates {
		r.SetFunctions()
		datasources_map[r.ResourceName] = r
	}
}

func GetDataSourceByName(name string) *ResourceTemplateV2 {
	resource, exists := datasources_map[name]
	if exists {
		return &resource
	}
	return nil
}
