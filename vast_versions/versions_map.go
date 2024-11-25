package vast_versions

import (
	version_4_6_0 "github.com/vast-data/terraform-provider-vastdata/codegen/4.6.0"
	version_4_7_0 "github.com/vast-data/terraform-provider-vastdata/codegen/4.7.0"
	version_5_0_0 "github.com/vast-data/terraform-provider-vastdata/codegen/5.0.0"
	version_5_1_0 "github.com/vast-data/terraform-provider-vastdata/codegen/5.1.0"
	version_5_2_0 "github.com/vast-data/terraform-provider-vastdata/codegen/5.2.0"
	version_5_3_0 "github.com/vast-data/terraform-provider-vastdata/codegen/5.3.0"
	"reflect"
)

var vast_versions map[string]map[string]reflect.Type = map[string]map[string]reflect.Type{
	"4.6.0": map[string]reflect.Type{
		"ActiveDirectory":                 reflect.TypeOf((*version_4_6_0.ActiveDirectory)(nil)).Elem(),
		"Cnode":                           reflect.TypeOf((*version_4_6_0.Cnode)(nil)).Elem(),
		"DefaultQuota":                    reflect.TypeOf((*version_4_6_0.DefaultQuota)(nil)).Elem(),
		"Dns":                             reflect.TypeOf((*version_4_6_0.Dns)(nil)).Elem(),
		"GlobalSnapshot":                  reflect.TypeOf((*version_4_6_0.GlobalSnapshot)(nil)).Elem(),
		"GlobalSnapshotOwnerRootSnapshot": reflect.TypeOf((*version_4_6_0.GlobalSnapshotOwnerRootSnapshot)(nil)).Elem(),
		"GlobalSnapshotOwnerTenant":       reflect.TypeOf((*version_4_6_0.GlobalSnapshotOwnerTenant)(nil)).Elem(),
		"Group":                           reflect.TypeOf((*version_4_6_0.Group)(nil)).Elem(),
		"Ldap":                            reflect.TypeOf((*version_4_6_0.Ldap)(nil)).Elem(),
		"Nis":                             reflect.TypeOf((*version_4_6_0.Nis)(nil)).Elem(),
		"ProtectedPath":                   reflect.TypeOf((*version_4_6_0.ProtectedPath)(nil)).Elem(),
		"ProtectionPolicy":                reflect.TypeOf((*version_4_6_0.ProtectionPolicy)(nil)).Elem(),
		"ProtectionPolicySchedule":        reflect.TypeOf((*version_4_6_0.ProtectionPolicySchedule)(nil)).Elem(),
		"QosDynamicLimits":                reflect.TypeOf((*version_4_6_0.QosDynamicLimits)(nil)).Elem(),
		"QosPolicy":                       reflect.TypeOf((*version_4_6_0.QosPolicy)(nil)).Elem(),
		"QosStaticLimits":                 reflect.TypeOf((*version_4_6_0.QosStaticLimits)(nil)).Elem(),
		"Quota":                           reflect.TypeOf((*version_4_6_0.Quota)(nil)).Elem(),
		"QuotaEntityInfo":                 reflect.TypeOf((*version_4_6_0.QuotaEntityInfo)(nil)).Elem(),
		"ReplicationPeers":                reflect.TypeOf((*version_4_6_0.ReplicationPeers)(nil)).Elem(),
		"S3LifeCycleRule":                 reflect.TypeOf((*version_4_6_0.S3LifeCycleRule)(nil)).Elem(),
		"S3Policy":                        reflect.TypeOf((*version_4_6_0.S3Policy)(nil)).Elem(),
		"S3replicationPeers":              reflect.TypeOf((*version_4_6_0.S3replicationPeers)(nil)).Elem(),
		"ShareAcl":                        reflect.TypeOf((*version_4_6_0.ShareAcl)(nil)).Elem(),
		"Snapshot":                        reflect.TypeOf((*version_4_6_0.Snapshot)(nil)).Elem(),
		"Tenant":                          reflect.TypeOf((*version_4_6_0.Tenant)(nil)).Elem(),
		"User":                            reflect.TypeOf((*version_4_6_0.User)(nil)).Elem(),
		"UserQuota":                       reflect.TypeOf((*version_4_6_0.UserQuota)(nil)).Elem(),
		"View":                            reflect.TypeOf((*version_4_6_0.View)(nil)).Elem(),
		"ViewCreate":                      reflect.TypeOf((*version_4_6_0.ViewCreate)(nil)).Elem(),
		"ViewPolicy":                      reflect.TypeOf((*version_4_6_0.ViewPolicy)(nil)).Elem(),
		"ViewShareAcl":                    reflect.TypeOf((*version_4_6_0.ViewShareAcl)(nil)).Elem(),
		"VipPool":                         reflect.TypeOf((*version_4_6_0.VipPool)(nil)).Elem(),
	},
	"4.7.0": map[string]reflect.Type{
		"ActiveDirectory":                 reflect.TypeOf((*version_4_7_0.ActiveDirectory)(nil)).Elem(),
		"Cnode":                           reflect.TypeOf((*version_4_7_0.Cnode)(nil)).Elem(),
		"DefaultQuota":                    reflect.TypeOf((*version_4_7_0.DefaultQuota)(nil)).Elem(),
		"Dns":                             reflect.TypeOf((*version_4_7_0.Dns)(nil)).Elem(),
		"GlobalLocalSnapshot":             reflect.TypeOf((*version_4_7_0.GlobalLocalSnapshot)(nil)).Elem(),
		"GlobalSnapshot":                  reflect.TypeOf((*version_4_7_0.GlobalSnapshot)(nil)).Elem(),
		"GlobalSnapshotOwnerRootSnapshot": reflect.TypeOf((*version_4_7_0.GlobalSnapshotOwnerRootSnapshot)(nil)).Elem(),
		"GlobalSnapshotOwnerTenant":       reflect.TypeOf((*version_4_7_0.GlobalSnapshotOwnerTenant)(nil)).Elem(),
		"Group":                           reflect.TypeOf((*version_4_7_0.Group)(nil)).Elem(),
		"Ldap":                            reflect.TypeOf((*version_4_7_0.Ldap)(nil)).Elem(),
		"Nis":                             reflect.TypeOf((*version_4_7_0.Nis)(nil)).Elem(),
		"ProtectedPath":                   reflect.TypeOf((*version_4_7_0.ProtectedPath)(nil)).Elem(),
		"ProtectionPolicy":                reflect.TypeOf((*version_4_7_0.ProtectionPolicy)(nil)).Elem(),
		"ProtectionPolicySchedule":        reflect.TypeOf((*version_4_7_0.ProtectionPolicySchedule)(nil)).Elem(),
		"QosDynamicLimits":                reflect.TypeOf((*version_4_7_0.QosDynamicLimits)(nil)).Elem(),
		"QosPolicy":                       reflect.TypeOf((*version_4_7_0.QosPolicy)(nil)).Elem(),
		"QosStaticLimits":                 reflect.TypeOf((*version_4_7_0.QosStaticLimits)(nil)).Elem(),
		"Quota":                           reflect.TypeOf((*version_4_7_0.Quota)(nil)).Elem(),
		"QuotaEntityInfo":                 reflect.TypeOf((*version_4_7_0.QuotaEntityInfo)(nil)).Elem(),
		"ReplicationPeers":                reflect.TypeOf((*version_4_7_0.ReplicationPeers)(nil)).Elem(),
		"S3LifeCycleRule":                 reflect.TypeOf((*version_4_7_0.S3LifeCycleRule)(nil)).Elem(),
		"S3Policy":                        reflect.TypeOf((*version_4_7_0.S3Policy)(nil)).Elem(),
		"S3replicationPeers":              reflect.TypeOf((*version_4_7_0.S3replicationPeers)(nil)).Elem(),
		"ShareAcl":                        reflect.TypeOf((*version_4_7_0.ShareAcl)(nil)).Elem(),
		"Snapshot":                        reflect.TypeOf((*version_4_7_0.Snapshot)(nil)).Elem(),
		"Tenant":                          reflect.TypeOf((*version_4_7_0.Tenant)(nil)).Elem(),
		"User":                            reflect.TypeOf((*version_4_7_0.User)(nil)).Elem(),
		"UserQuota":                       reflect.TypeOf((*version_4_7_0.UserQuota)(nil)).Elem(),
		"View":                            reflect.TypeOf((*version_4_7_0.View)(nil)).Elem(),
		"ViewCreate":                      reflect.TypeOf((*version_4_7_0.ViewCreate)(nil)).Elem(),
		"ViewPolicy":                      reflect.TypeOf((*version_4_7_0.ViewPolicy)(nil)).Elem(),
		"ViewShareAcl":                    reflect.TypeOf((*version_4_7_0.ViewShareAcl)(nil)).Elem(),
		"VipPool":                         reflect.TypeOf((*version_4_7_0.VipPool)(nil)).Elem(),
	},
	"5.0.0": map[string]reflect.Type{
		"ActiveDirectory":                 reflect.TypeOf((*version_5_0_0.ActiveDirectory)(nil)).Elem(),
		"Cnode":                           reflect.TypeOf((*version_5_0_0.Cnode)(nil)).Elem(),
		"DefaultQuota":                    reflect.TypeOf((*version_5_0_0.DefaultQuota)(nil)).Elem(),
		"Dns":                             reflect.TypeOf((*version_5_0_0.Dns)(nil)).Elem(),
		"GlobalLocalSnapshot":             reflect.TypeOf((*version_5_0_0.GlobalLocalSnapshot)(nil)).Elem(),
		"GlobalSnapshot":                  reflect.TypeOf((*version_5_0_0.GlobalSnapshot)(nil)).Elem(),
		"GlobalSnapshotOwnerRootSnapshot": reflect.TypeOf((*version_5_0_0.GlobalSnapshotOwnerRootSnapshot)(nil)).Elem(),
		"GlobalSnapshotOwnerTenant":       reflect.TypeOf((*version_5_0_0.GlobalSnapshotOwnerTenant)(nil)).Elem(),
		"Group":                           reflect.TypeOf((*version_5_0_0.Group)(nil)).Elem(),
		"Ldap":                            reflect.TypeOf((*version_5_0_0.Ldap)(nil)).Elem(),
		"Manager":                         reflect.TypeOf((*version_5_0_0.Manager)(nil)).Elem(),
		"Nis":                             reflect.TypeOf((*version_5_0_0.Nis)(nil)).Elem(),
		"ProtectedPath":                   reflect.TypeOf((*version_5_0_0.ProtectedPath)(nil)).Elem(),
		"ProtectionPolicy":                reflect.TypeOf((*version_5_0_0.ProtectionPolicy)(nil)).Elem(),
		"ProtectionPolicySchedule":        reflect.TypeOf((*version_5_0_0.ProtectionPolicySchedule)(nil)).Elem(),
		"QosDynamicLimits":                reflect.TypeOf((*version_5_0_0.QosDynamicLimits)(nil)).Elem(),
		"QosPolicy":                       reflect.TypeOf((*version_5_0_0.QosPolicy)(nil)).Elem(),
		"QosStaticLimits":                 reflect.TypeOf((*version_5_0_0.QosStaticLimits)(nil)).Elem(),
		"Quota":                           reflect.TypeOf((*version_5_0_0.Quota)(nil)).Elem(),
		"QuotaEntityInfo":                 reflect.TypeOf((*version_5_0_0.QuotaEntityInfo)(nil)).Elem(),
		"ReplicationPeers":                reflect.TypeOf((*version_5_0_0.ReplicationPeers)(nil)).Elem(),
		"Role":                            reflect.TypeOf((*version_5_0_0.Role)(nil)).Elem(),
		"S3LifeCycleRule":                 reflect.TypeOf((*version_5_0_0.S3LifeCycleRule)(nil)).Elem(),
		"S3Policy":                        reflect.TypeOf((*version_5_0_0.S3Policy)(nil)).Elem(),
		"S3replicationPeers":              reflect.TypeOf((*version_5_0_0.S3replicationPeers)(nil)).Elem(),
		"ShareAcl":                        reflect.TypeOf((*version_5_0_0.ShareAcl)(nil)).Elem(),
		"Snapshot":                        reflect.TypeOf((*version_5_0_0.Snapshot)(nil)).Elem(),
		"Tenant":                          reflect.TypeOf((*version_5_0_0.Tenant)(nil)).Elem(),
		"User":                            reflect.TypeOf((*version_5_0_0.User)(nil)).Elem(),
		"UserKey":                         reflect.TypeOf((*version_5_0_0.UserKey)(nil)).Elem(),
		"UserQuota":                       reflect.TypeOf((*version_5_0_0.UserQuota)(nil)).Elem(),
		"View":                            reflect.TypeOf((*version_5_0_0.View)(nil)).Elem(),
		"ViewCreate":                      reflect.TypeOf((*version_5_0_0.ViewCreate)(nil)).Elem(),
		"ViewPolicy":                      reflect.TypeOf((*version_5_0_0.ViewPolicy)(nil)).Elem(),
		"ViewShareAcl":                    reflect.TypeOf((*version_5_0_0.ViewShareAcl)(nil)).Elem(),
		"VipPool":                         reflect.TypeOf((*version_5_0_0.VipPool)(nil)).Elem(),
	},
	"5.1.0": map[string]reflect.Type{
		"ActiveDirectory":                 reflect.TypeOf((*version_5_1_0.ActiveDirectory)(nil)).Elem(),
		"ActiveDirectory2":                reflect.TypeOf((*version_5_1_0.ActiveDirectory2)(nil)).Elem(),
		"Cnode":                           reflect.TypeOf((*version_5_1_0.Cnode)(nil)).Elem(),
		"DefaultQuota":                    reflect.TypeOf((*version_5_1_0.DefaultQuota)(nil)).Elem(),
		"Dns":                             reflect.TypeOf((*version_5_1_0.Dns)(nil)).Elem(),
		"GlobalLocalSnapshot":             reflect.TypeOf((*version_5_1_0.GlobalLocalSnapshot)(nil)).Elem(),
		"GlobalSnapshot":                  reflect.TypeOf((*version_5_1_0.GlobalSnapshot)(nil)).Elem(),
		"GlobalSnapshotOwnerRootSnapshot": reflect.TypeOf((*version_5_1_0.GlobalSnapshotOwnerRootSnapshot)(nil)).Elem(),
		"GlobalSnapshotOwnerTenant":       reflect.TypeOf((*version_5_1_0.GlobalSnapshotOwnerTenant)(nil)).Elem(),
		"Group":                           reflect.TypeOf((*version_5_1_0.Group)(nil)).Elem(),
		"Ldap":                            reflect.TypeOf((*version_5_1_0.Ldap)(nil)).Elem(),
		"Manager":                         reflect.TypeOf((*version_5_1_0.Manager)(nil)).Elem(),
		"Nis":                             reflect.TypeOf((*version_5_1_0.Nis)(nil)).Elem(),
		"PermissionsPerVipPool":           reflect.TypeOf((*version_5_1_0.PermissionsPerVipPool)(nil)).Elem(),
		"ProtectedPath":                   reflect.TypeOf((*version_5_1_0.ProtectedPath)(nil)).Elem(),
		"ProtectionPolicy":                reflect.TypeOf((*version_5_1_0.ProtectionPolicy)(nil)).Elem(),
		"ProtectionPolicySchedule":        reflect.TypeOf((*version_5_1_0.ProtectionPolicySchedule)(nil)).Elem(),
		"ProtocolsAudit":                  reflect.TypeOf((*version_5_1_0.ProtocolsAudit)(nil)).Elem(),
		"QosDynamicLimits":                reflect.TypeOf((*version_5_1_0.QosDynamicLimits)(nil)).Elem(),
		"QosPolicy":                       reflect.TypeOf((*version_5_1_0.QosPolicy)(nil)).Elem(),
		"QosStaticLimits":                 reflect.TypeOf((*version_5_1_0.QosStaticLimits)(nil)).Elem(),
		"Quota":                           reflect.TypeOf((*version_5_1_0.Quota)(nil)).Elem(),
		"QuotaEntityInfo":                 reflect.TypeOf((*version_5_1_0.QuotaEntityInfo)(nil)).Elem(),
		"ReplicationPeers":                reflect.TypeOf((*version_5_1_0.ReplicationPeers)(nil)).Elem(),
		"Role":                            reflect.TypeOf((*version_5_1_0.Role)(nil)).Elem(),
		"S3LifeCycleRule":                 reflect.TypeOf((*version_5_1_0.S3LifeCycleRule)(nil)).Elem(),
		"S3Policy":                        reflect.TypeOf((*version_5_1_0.S3Policy)(nil)).Elem(),
		"S3replicationPeers":              reflect.TypeOf((*version_5_1_0.S3replicationPeers)(nil)).Elem(),
		"ShareAcl":                        reflect.TypeOf((*version_5_1_0.ShareAcl)(nil)).Elem(),
		"Snapshot":                        reflect.TypeOf((*version_5_1_0.Snapshot)(nil)).Elem(),
		"Tenant":                          reflect.TypeOf((*version_5_1_0.Tenant)(nil)).Elem(),
		"User":                            reflect.TypeOf((*version_5_1_0.User)(nil)).Elem(),
		"UserKey":                         reflect.TypeOf((*version_5_1_0.UserKey)(nil)).Elem(),
		"UserQuota":                       reflect.TypeOf((*version_5_1_0.UserQuota)(nil)).Elem(),
		"View":                            reflect.TypeOf((*version_5_1_0.View)(nil)).Elem(),
		"ViewCreate":                      reflect.TypeOf((*version_5_1_0.ViewCreate)(nil)).Elem(),
		"ViewPolicy":                      reflect.TypeOf((*version_5_1_0.ViewPolicy)(nil)).Elem(),
		"ViewShareAcl":                    reflect.TypeOf((*version_5_1_0.ViewShareAcl)(nil)).Elem(),
		"VipPool":                         reflect.TypeOf((*version_5_1_0.VipPool)(nil)).Elem(),
	},
	"5.2.0": map[string]reflect.Type{
		"ActiveDirectory":                 reflect.TypeOf((*version_5_2_0.ActiveDirectory)(nil)).Elem(),
		"ActiveDirectory2":                reflect.TypeOf((*version_5_2_0.ActiveDirectory2)(nil)).Elem(),
		"BucketLogging":                   reflect.TypeOf((*version_5_2_0.BucketLogging)(nil)).Elem(),
		"Cnode":                           reflect.TypeOf((*version_5_2_0.Cnode)(nil)).Elem(),
		"DefaultQuota":                    reflect.TypeOf((*version_5_2_0.DefaultQuota)(nil)).Elem(),
		"Dns":                             reflect.TypeOf((*version_5_2_0.Dns)(nil)).Elem(),
		"GlobalLocalSnapshot":             reflect.TypeOf((*version_5_2_0.GlobalLocalSnapshot)(nil)).Elem(),
		"GlobalSnapshot":                  reflect.TypeOf((*version_5_2_0.GlobalSnapshot)(nil)).Elem(),
		"GlobalSnapshotOwnerRootSnapshot": reflect.TypeOf((*version_5_2_0.GlobalSnapshotOwnerRootSnapshot)(nil)).Elem(),
		"GlobalSnapshotOwnerTenant":       reflect.TypeOf((*version_5_2_0.GlobalSnapshotOwnerTenant)(nil)).Elem(),
		"Group":                           reflect.TypeOf((*version_5_2_0.Group)(nil)).Elem(),
		"Ldap":                            reflect.TypeOf((*version_5_2_0.Ldap)(nil)).Elem(),
		"Manager":                         reflect.TypeOf((*version_5_2_0.Manager)(nil)).Elem(),
		"Nis":                             reflect.TypeOf((*version_5_2_0.Nis)(nil)).Elem(),
		"PermissionsPerVipPool":           reflect.TypeOf((*version_5_2_0.PermissionsPerVipPool)(nil)).Elem(),
		"ProtectedPath":                   reflect.TypeOf((*version_5_2_0.ProtectedPath)(nil)).Elem(),
		"ProtectionPolicy":                reflect.TypeOf((*version_5_2_0.ProtectionPolicy)(nil)).Elem(),
		"ProtectionPolicySchedule":        reflect.TypeOf((*version_5_2_0.ProtectionPolicySchedule)(nil)).Elem(),
		"ProtocolsAudit":                  reflect.TypeOf((*version_5_2_0.ProtocolsAudit)(nil)).Elem(),
		"QoSDynamicTotalLimits":           reflect.TypeOf((*version_5_2_0.QoSDynamicTotalLimits)(nil)).Elem(),
		"QoSStaticTotalLimits":            reflect.TypeOf((*version_5_2_0.QoSStaticTotalLimits)(nil)).Elem(),
		"QosDynamicLimits":                reflect.TypeOf((*version_5_2_0.QosDynamicLimits)(nil)).Elem(),
		"QosPolicy":                       reflect.TypeOf((*version_5_2_0.QosPolicy)(nil)).Elem(),
		"QosStaticLimits":                 reflect.TypeOf((*version_5_2_0.QosStaticLimits)(nil)).Elem(),
		"QosUser":                         reflect.TypeOf((*version_5_2_0.QosUser)(nil)).Elem(),
		"Quota":                           reflect.TypeOf((*version_5_2_0.Quota)(nil)).Elem(),
		"QuotaEntityInfo":                 reflect.TypeOf((*version_5_2_0.QuotaEntityInfo)(nil)).Elem(),
		"Realm":                           reflect.TypeOf((*version_5_2_0.Realm)(nil)).Elem(),
		"ReplicationPeers":                reflect.TypeOf((*version_5_2_0.ReplicationPeers)(nil)).Elem(),
		"Role":                            reflect.TypeOf((*version_5_2_0.Role)(nil)).Elem(),
		"S3LifeCycleRule":                 reflect.TypeOf((*version_5_2_0.S3LifeCycleRule)(nil)).Elem(),
		"S3Policy":                        reflect.TypeOf((*version_5_2_0.S3Policy)(nil)).Elem(),
		"S3replicationPeers":              reflect.TypeOf((*version_5_2_0.S3replicationPeers)(nil)).Elem(),
		"ShareAcl":                        reflect.TypeOf((*version_5_2_0.ShareAcl)(nil)).Elem(),
		"Snapshot":                        reflect.TypeOf((*version_5_2_0.Snapshot)(nil)).Elem(),
		"Tenant":                          reflect.TypeOf((*version_5_2_0.Tenant)(nil)).Elem(),
		"User":                            reflect.TypeOf((*version_5_2_0.User)(nil)).Elem(),
		"UserKey":                         reflect.TypeOf((*version_5_2_0.UserKey)(nil)).Elem(),
		"UserQuota":                       reflect.TypeOf((*version_5_2_0.UserQuota)(nil)).Elem(),
		"View":                            reflect.TypeOf((*version_5_2_0.View)(nil)).Elem(),
		"ViewCreate":                      reflect.TypeOf((*version_5_2_0.ViewCreate)(nil)).Elem(),
		"ViewPolicy":                      reflect.TypeOf((*version_5_2_0.ViewPolicy)(nil)).Elem(),
		"ViewShareAcl":                    reflect.TypeOf((*version_5_2_0.ViewShareAcl)(nil)).Elem(),
		"VipPool":                         reflect.TypeOf((*version_5_2_0.VipPool)(nil)).Elem(),
	},
	"5.3.0": map[string]reflect.Type{
		"ActiveDirectory":                 reflect.TypeOf((*version_5_3_0.ActiveDirectory)(nil)).Elem(),
		"ActiveDirectory2":                reflect.TypeOf((*version_5_3_0.ActiveDirectory2)(nil)).Elem(),
		"BucketLogging":                   reflect.TypeOf((*version_5_3_0.BucketLogging)(nil)).Elem(),
		"Cnode":                           reflect.TypeOf((*version_5_3_0.Cnode)(nil)).Elem(),
		"DefaultQuota":                    reflect.TypeOf((*version_5_3_0.DefaultQuota)(nil)).Elem(),
		"Dns":                             reflect.TypeOf((*version_5_3_0.Dns)(nil)).Elem(),
		"GlobalLocalSnapshot":             reflect.TypeOf((*version_5_3_0.GlobalLocalSnapshot)(nil)).Elem(),
		"GlobalSnapshot":                  reflect.TypeOf((*version_5_3_0.GlobalSnapshot)(nil)).Elem(),
		"GlobalSnapshotOwnerRootSnapshot": reflect.TypeOf((*version_5_3_0.GlobalSnapshotOwnerRootSnapshot)(nil)).Elem(),
		"GlobalSnapshotOwnerTenant":       reflect.TypeOf((*version_5_3_0.GlobalSnapshotOwnerTenant)(nil)).Elem(),
		"Group":                           reflect.TypeOf((*version_5_3_0.Group)(nil)).Elem(),
		"KafkaBrokerAddressParams":        reflect.TypeOf((*version_5_3_0.KafkaBrokerAddressParams)(nil)).Elem(),
		"Kafkabroker":                     reflect.TypeOf((*version_5_3_0.Kafkabroker)(nil)).Elem(),
		"Ldap":                            reflect.TypeOf((*version_5_3_0.Ldap)(nil)).Elem(),
		"Manager":                         reflect.TypeOf((*version_5_3_0.Manager)(nil)).Elem(),
		"Nis":                             reflect.TypeOf((*version_5_3_0.Nis)(nil)).Elem(),
		"PermissionsPerVipPool":           reflect.TypeOf((*version_5_3_0.PermissionsPerVipPool)(nil)).Elem(),
		"ProtectedPath":                   reflect.TypeOf((*version_5_3_0.ProtectedPath)(nil)).Elem(),
		"ProtectionPolicy":                reflect.TypeOf((*version_5_3_0.ProtectionPolicy)(nil)).Elem(),
		"ProtectionPolicySchedule":        reflect.TypeOf((*version_5_3_0.ProtectionPolicySchedule)(nil)).Elem(),
		"ProtocolsAudit":                  reflect.TypeOf((*version_5_3_0.ProtocolsAudit)(nil)).Elem(),
		"QoSDynamicTotalLimits":           reflect.TypeOf((*version_5_3_0.QoSDynamicTotalLimits)(nil)).Elem(),
		"QoSStaticTotalLimits":            reflect.TypeOf((*version_5_3_0.QoSStaticTotalLimits)(nil)).Elem(),
		"QosDynamicLimits":                reflect.TypeOf((*version_5_3_0.QosDynamicLimits)(nil)).Elem(),
		"QosPolicy":                       reflect.TypeOf((*version_5_3_0.QosPolicy)(nil)).Elem(),
		"QosStaticLimits":                 reflect.TypeOf((*version_5_3_0.QosStaticLimits)(nil)).Elem(),
		"QosUser":                         reflect.TypeOf((*version_5_3_0.QosUser)(nil)).Elem(),
		"Quota":                           reflect.TypeOf((*version_5_3_0.Quota)(nil)).Elem(),
		"QuotaEntityInfo":                 reflect.TypeOf((*version_5_3_0.QuotaEntityInfo)(nil)).Elem(),
		"Realm":                           reflect.TypeOf((*version_5_3_0.Realm)(nil)).Elem(),
		"ReplicationPeers":                reflect.TypeOf((*version_5_3_0.ReplicationPeers)(nil)).Elem(),
		"Role":                            reflect.TypeOf((*version_5_3_0.Role)(nil)).Elem(),
		"S3LifeCycleRule":                 reflect.TypeOf((*version_5_3_0.S3LifeCycleRule)(nil)).Elem(),
		"S3Policy":                        reflect.TypeOf((*version_5_3_0.S3Policy)(nil)).Elem(),
		"S3replicationPeers":              reflect.TypeOf((*version_5_3_0.S3replicationPeers)(nil)).Elem(),
		"ShareAcl":                        reflect.TypeOf((*version_5_3_0.ShareAcl)(nil)).Elem(),
		"Snapshot":                        reflect.TypeOf((*version_5_3_0.Snapshot)(nil)).Elem(),
		"Tenant":                          reflect.TypeOf((*version_5_3_0.Tenant)(nil)).Elem(),
		"User":                            reflect.TypeOf((*version_5_3_0.User)(nil)).Elem(),
		"UserKey":                         reflect.TypeOf((*version_5_3_0.UserKey)(nil)).Elem(),
		"UserQuota":                       reflect.TypeOf((*version_5_3_0.UserQuota)(nil)).Elem(),
		"View":                            reflect.TypeOf((*version_5_3_0.View)(nil)).Elem(),
		"ViewCreate":                      reflect.TypeOf((*version_5_3_0.ViewCreate)(nil)).Elem(),
		"ViewPolicy":                      reflect.TypeOf((*version_5_3_0.ViewPolicy)(nil)).Elem(),
		"ViewShareAcl":                    reflect.TypeOf((*version_5_3_0.ViewShareAcl)(nil)).Elem(),
		"VipPool":                         reflect.TypeOf((*version_5_3_0.VipPool)(nil)).Elem(),
		"Volume":                          reflect.TypeOf((*version_5_3_0.Volume)(nil)).Elem(),
	},
}
