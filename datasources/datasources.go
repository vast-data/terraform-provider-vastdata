package datasources

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var DataSources map[string]*schema.Resource = map[string]*schema.Resource{
	"vastdata_cnode":                DataSourceCnode(),
	"vastdata_qos_policy":           DataSourceQosPolicy(),
	"vastdata_quota":                DataSourceQuota(),
	"vastdata_dns":                  DataSourceDns(),
	"vastdata_vip_pool":             DataSourceVipPool(),
	"vastdata_view_policy":          DataSourceViewPolicy(),
	"vastdata_group":                DataSourceGroup(),
	"vastdata_user":                 DataSourceUser(),
	"vastdata_view":                 DataSourceView(),
	"vastdata_nis":                  DataSourceNis(),
	"vastdata_tenant":               DataSourceTenant(),
	"vastdata_ldap":                 DataSourceLdap(),
	"vastdata_s3_life_cycle_rule":   DataSourceS3LifeCycleRule(),
	"vastdata_active_directory":     DataSourceActiveDirectory(),
	"vastdata_s3_policy":            DataSourceS3Policy(),
	"vastdata_protected_path":       DataSourceProtectedPath(),
	"vastdata_snapshot":             DataSourceSnapshot(),
	"vastdata_global_snapshot":      DataSourceGlobalSnapshot(),
	"vastdata_replication_peers":    DataSourceReplicationPeers(),
	"vastdata_protection_policy":    DataSourceProtectionPolicy(),
	"vastdata_s3_replication_peers": DataSourceS3replicationPeers(),
	"vastdata_active_directory2":    DataSourceActiveDirectory2(),
	"vastdata_administators_realms": DataSourceRealm(),
	"vastdata_administators_roles":  DataSourceRole(),
}
