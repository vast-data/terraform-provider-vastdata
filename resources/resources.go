package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var Resources map[string]*schema.Resource = map[string]*schema.Resource{
	"vastdata_user":                   ResourceUser(),
	"vastdata_group":                  ResourceGroup(),
	"vastdata_vip_pool":               ResourceVipPool(),
	"vastdata_tenant":                 ResourceTenant(),
	"vastdata_qos_policy":             ResourceQosPolicy(),
	"vastdata_protection_policy":      ResourceProtectionPolicy(),
	"vastdata_quota":                  ResourceQuota(),
	"vastdata_dns":                    ResourceDns(),
	"vastdata_view_policy":            ResourceViewPolicy(),
	"vastdata_view":                   ResourceView(),
	"vastdata_nis":                    ResourceNis(),
	"vastdata_ldap":                   ResourceLdap(),
	"vastdata_s3_life_cycle_rule":     ResourceS3LifeCycleRule(),
	"vastdata_active_directory":       ResourceActiveDirectory(),
	"vastdata_s3_policy":              ResourceS3Policy(),
	"vastdata_protected_path":         ResourceProtectedPath(),
	"vastdata_snapshot":               ResourceSnapshot(),
	"vastdata_global_snapshot":        ResourceGlobalSnapshot(),
	"vastdata_global_local_snapshot":  ResourceGlobalLocalSnapshot(),
	"vastdata_replication_peers":      ResourceReplicationPeers(),
	"vastdata_s3_replication_peers":   ResourceS3replicationPeers(),
	"vastdata_user_key":               ResourceUserKey(),
	"vastdata_active_directory2":      ResourceActiveDirectory2(),
	"vastdata_administators_realms":   ResourceRealm(),
	"vastdata_administators_roles":    ResourceRole(),
	"vastdata_administators_managers": ResourceManager(),
}
