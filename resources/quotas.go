package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"reflect"

	//        "net/url"
	"errors"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api_latest "github.com/vast-data/terraform-provider-vastdata/codegen/latest"
	codegen_configs "github.com/vast-data/terraform-provider-vastdata/codegen_tools/configs"
	metadata "github.com/vast-data/terraform-provider-vastdata/metadata"
	utils "github.com/vast-data/terraform-provider-vastdata/utils"
	vast_client "github.com/vast-data/terraform-provider-vastdata/vast-client"
	vast_versions "github.com/vast-data/terraform-provider-vastdata/vast_versions"
)

func ResourceQuota() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceQuotaRead,
		DeleteContext: resourceQuotaDelete,
		CreateContext: resourceQuotaCreate,
		UpdateContext: resourceQuotaUpdate,

		Importer: &schema.ResourceImporter{
			StateContext: resourceQuotaImporter,
		},

		Description: `This is a quota`,
		Schema:      getResourceQuotaSchema(),
	}
}

func getResourceQuotaSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{

		"guid": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Quota").GetConflictingFields("guid"),

			Computed:    true,
			Optional:    false,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Quota guid`,
		},

		"name": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Quota").GetConflictingFields("name"),

			Required:    true,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The name`,
		},

		"state": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Quota").GetConflictingFields("state"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,
		},

		"pretty_state": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Quota").GetConflictingFields("pretty_state"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,
		},

		"path": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Quota").GetConflictingFields("path"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Directory path`,
		},

		"pretty_grace_period": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Quota").GetConflictingFields("pretty_grace_period"),

			Computed:    true,
			Optional:    false,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Quota enforcement pretty grace period in seconds, minutes, hours or days. Example: 90m`,
		},

		"grace_period": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Quota").GetConflictingFields("grace_period"),

			Computed:  true,
			Optional:  true,
			Sensitive: false,

			ValidateDiagFunc: utils.GracePeriodFormatValidation,
		},

		"time_to_block": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Quota").GetConflictingFields("time_to_block"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Grace period expiration time`,
		},

		"soft_limit": {
			Type:          schema.TypeInt,
			ConflictsWith: codegen_configs.GetResourceByName("Quota").GetConflictingFields("soft_limit"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Soft quota limit`,
		},

		"hard_limit": {
			Type:          schema.TypeInt,
			ConflictsWith: codegen_configs.GetResourceByName("Quota").GetConflictingFields("hard_limit"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Hard quota limit`,
		},

		"hard_limit_inodes": {
			Type:          schema.TypeInt,
			ConflictsWith: codegen_configs.GetResourceByName("Quota").GetConflictingFields("hard_limit_inodes"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Hard inodes quota limit`,
		},

		"soft_limit_inodes": {
			Type:          schema.TypeInt,
			ConflictsWith: codegen_configs.GetResourceByName("Quota").GetConflictingFields("soft_limit_inodes"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Soft inodes quota limit`,
		},

		"used_inodes": {
			Type:          schema.TypeInt,
			ConflictsWith: codegen_configs.GetResourceByName("Quota").GetConflictingFields("used_inodes"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Used inodes`,
		},

		"used_capacity": {
			Type:          schema.TypeInt,
			ConflictsWith: codegen_configs.GetResourceByName("Quota").GetConflictingFields("used_capacity"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Used capacity in bytes`,
		},

		"used_capacity_tb": {
			Type:          schema.TypeFloat,
			ConflictsWith: codegen_configs.GetResourceByName("Quota").GetConflictingFields("used_capacity_tb"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Used capacity in TB`,
		},

		"used_effective_capacity": {
			Type:          schema.TypeInt,
			ConflictsWith: codegen_configs.GetResourceByName("Quota").GetConflictingFields("used_effective_capacity"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Used effective capacity in bytes`,
		},

		"used_effective_capacity_tb": {
			Type:          schema.TypeFloat,
			ConflictsWith: codegen_configs.GetResourceByName("Quota").GetConflictingFields("used_effective_capacity_tb"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Used effective capacity in TB`,
		},

		"tenant_id": {
			Type:          schema.TypeInt,
			ConflictsWith: codegen_configs.GetResourceByName("Quota").GetConflictingFields("tenant_id"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Tenant ID`,
		},

		"tenant_name": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Quota").GetConflictingFields("tenant_name"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Tenant Name`,
		},

		"cluster": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Quota").GetConflictingFields("cluster"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Parent Cluster`,
		},

		"cluster_id": {
			Type:          schema.TypeInt,
			ConflictsWith: codegen_configs.GetResourceByName("Quota").GetConflictingFields("cluster_id"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Parent Cluster ID`,
		},

		"system_id": {
			Type:          schema.TypeInt,
			ConflictsWith: codegen_configs.GetResourceByName("Quota").GetConflictingFields("system_id"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,
		},

		"is_user_quota": {
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("Quota").GetConflictingFields("is_user_quota"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,
		},

		"enable_email_providers": {
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("Quota").GetConflictingFields("enable_email_providers"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,
		},

		"num_exceeded_users": {
			Type:          schema.TypeInt,
			ConflictsWith: codegen_configs.GetResourceByName("Quota").GetConflictingFields("num_exceeded_users"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,
		},

		"num_blocked_users": {
			Type:          schema.TypeInt,
			ConflictsWith: codegen_configs.GetResourceByName("Quota").GetConflictingFields("num_blocked_users"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,
		},

		"enable_alarms": {
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("Quota").GetConflictingFields("enable_alarms"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Enable alarms when users or groups are exceeding their limit`,
		},

		"default_email": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Quota").GetConflictingFields("default_email"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The default Email if there is no suffix and no address in the providers`,
		},

		"percent_inodes": {
			Type:          schema.TypeInt,
			ConflictsWith: codegen_configs.GetResourceByName("Quota").GetConflictingFields("percent_inodes"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Percent of used inodes out of the hard limit`,
		},

		"percent_capacity": {
			Type:          schema.TypeInt,
			ConflictsWith: codegen_configs.GetResourceByName("Quota").GetConflictingFields("percent_capacity"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Percent of used capacity out of the hard limit`,
		},

		"default_user_quota": {
			Type:          schema.TypeList,
			ConflictsWith: codegen_configs.GetResourceByName("Quota").GetConflictingFields("default_user_quota"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,

			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{

					"quota_system_id": {
						Type:          schema.TypeInt,
						ConflictsWith: codegen_configs.GetResourceByName("DefaultQuota").GetConflictingFields("quota_system_id"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The system ID of the quota`,
					},

					"soft_limit": {
						Type:          schema.TypeInt,
						ConflictsWith: codegen_configs.GetResourceByName("DefaultQuota").GetConflictingFields("soft_limit"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The size soft limit in bytes`,
					},

					"hard_limit": {
						Type:          schema.TypeInt,
						ConflictsWith: codegen_configs.GetResourceByName("DefaultQuota").GetConflictingFields("hard_limit"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The size hard limit in bytes`,
					},

					"sof_limit_inodes": {
						Type:          schema.TypeInt,
						ConflictsWith: codegen_configs.GetResourceByName("DefaultQuota").GetConflictingFields("sof_limit_inodes"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The sof limit of inodes number`,
					},

					"hard_limit_inodes": {
						Type:          schema.TypeInt,
						ConflictsWith: codegen_configs.GetResourceByName("DefaultQuota").GetConflictingFields("hard_limit_inodes"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The hard limit in inode number`,
					},

					"grace_period": {
						Type:          schema.TypeString,
						ConflictsWith: codegen_configs.GetResourceByName("DefaultQuota").GetConflictingFields("grace_period"),

						Computed:  true,
						Optional:  true,
						Sensitive: false,

						ValidateDiagFunc: utils.GracePeriodFormatValidation,
					},
				},
			},
		},

		"default_group_quota": {
			Type:          schema.TypeList,
			ConflictsWith: codegen_configs.GetResourceByName("Quota").GetConflictingFields("default_group_quota"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,

			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{

					"quota_system_id": {
						Type:          schema.TypeInt,
						ConflictsWith: codegen_configs.GetResourceByName("DefaultQuota").GetConflictingFields("quota_system_id"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The system ID of the quota`,
					},

					"soft_limit": {
						Type:          schema.TypeInt,
						ConflictsWith: codegen_configs.GetResourceByName("DefaultQuota").GetConflictingFields("soft_limit"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The size soft limit in bytes`,
					},

					"hard_limit": {
						Type:          schema.TypeInt,
						ConflictsWith: codegen_configs.GetResourceByName("DefaultQuota").GetConflictingFields("hard_limit"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The size hard limit in bytes`,
					},

					"sof_limit_inodes": {
						Type:          schema.TypeInt,
						ConflictsWith: codegen_configs.GetResourceByName("DefaultQuota").GetConflictingFields("sof_limit_inodes"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The sof limit of inodes number`,
					},

					"hard_limit_inodes": {
						Type:          schema.TypeInt,
						ConflictsWith: codegen_configs.GetResourceByName("DefaultQuota").GetConflictingFields("hard_limit_inodes"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The hard limit in inode number`,
					},

					"grace_period": {
						Type:          schema.TypeString,
						ConflictsWith: codegen_configs.GetResourceByName("DefaultQuota").GetConflictingFields("grace_period"),

						Computed:  true,
						Optional:  true,
						Sensitive: false,

						ValidateDiagFunc: utils.GracePeriodFormatValidation,
					},
				},
			},
		},

		"user_quotas": {
			Type:          schema.TypeList,
			ConflictsWith: codegen_configs.GetResourceByName("Quota").GetConflictingFields("user_quotas"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,

			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{

					"grace_period": {
						Type:          schema.TypeString,
						ConflictsWith: codegen_configs.GetResourceByName("UserQuota").GetConflictingFields("grace_period"),

						Computed:  true,
						Optional:  true,
						Sensitive: false,

						ValidateDiagFunc: utils.GracePeriodFormatValidation,
					},

					"time_to_block": {
						Type:          schema.TypeString,
						ConflictsWith: codegen_configs.GetResourceByName("UserQuota").GetConflictingFields("time_to_block"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Grace period expiration time`,
					},

					"soft_limit": {
						Type:          schema.TypeInt,
						ConflictsWith: codegen_configs.GetResourceByName("UserQuota").GetConflictingFields("soft_limit"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Soft quota limit`,
					},

					"hard_limit": {
						Type:          schema.TypeInt,
						ConflictsWith: codegen_configs.GetResourceByName("UserQuota").GetConflictingFields("hard_limit"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Hard quota limit`,
					},

					"hard_limit_inodes": {
						Type:          schema.TypeInt,
						ConflictsWith: codegen_configs.GetResourceByName("UserQuota").GetConflictingFields("hard_limit_inodes"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Hard inodes quota limit`,
					},

					"soft_limit_inodes": {
						Type:          schema.TypeInt,
						ConflictsWith: codegen_configs.GetResourceByName("UserQuota").GetConflictingFields("soft_limit_inodes"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Soft inodes quota limit`,
					},

					"used_inodes": {
						Type:          schema.TypeInt,
						ConflictsWith: codegen_configs.GetResourceByName("UserQuota").GetConflictingFields("used_inodes"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Used inodes`,
					},

					"used_capacity": {
						Type:          schema.TypeInt,
						ConflictsWith: codegen_configs.GetResourceByName("UserQuota").GetConflictingFields("used_capacity"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Used capacity in bytes`,
					},

					"is_accountable": {
						Type:          schema.TypeBool,
						ConflictsWith: codegen_configs.GetResourceByName("UserQuota").GetConflictingFields("is_accountable"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,
					},

					"quota_system_id": {
						Type:          schema.TypeInt,
						ConflictsWith: codegen_configs.GetResourceByName("UserQuota").GetConflictingFields("quota_system_id"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,
					},

					"entity": {
						Type:          schema.TypeList,
						ConflictsWith: codegen_configs.GetResourceByName("UserQuota").GetConflictingFields("entity"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,

						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{

								"name": {
									Type:          schema.TypeString,
									ConflictsWith: codegen_configs.GetResourceByName("QuotaEntityInfo").GetConflictingFields("name"),

									Computed:    true,
									Optional:    true,
									Sensitive:   false,
									Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The name of the entity`,
								},

								"vast_id": {
									Type:          schema.TypeInt,
									ConflictsWith: codegen_configs.GetResourceByName("QuotaEntityInfo").GetConflictingFields("vast_id"),

									Computed:    true,
									Optional:    true,
									Sensitive:   false,
									Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,
								},

								"email": {
									Type:          schema.TypeString,
									ConflictsWith: codegen_configs.GetResourceByName("QuotaEntityInfo").GetConflictingFields("email"),

									Computed:    true,
									Optional:    true,
									Sensitive:   false,
									Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,
								},

								"is_group": {
									Type:          schema.TypeBool,
									ConflictsWith: codegen_configs.GetResourceByName("QuotaEntityInfo").GetConflictingFields("is_group"),

									Computed:    true,
									Optional:    true,
									Sensitive:   false,
									Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,
								},

								"identifier": {
									Type:          schema.TypeString,
									ConflictsWith: codegen_configs.GetResourceByName("QuotaEntityInfo").GetConflictingFields("identifier"),

									Required:    true,
									Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,
								},

								"identifier_type": {
									Type:          schema.TypeString,
									ConflictsWith: codegen_configs.GetResourceByName("QuotaEntityInfo").GetConflictingFields("identifier_type"),

									Computed:    true,
									Optional:    true,
									Sensitive:   false,
									Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,
								},
							},
						},
					},
				},
			},
		},

		"group_quotas": {
			Type:          schema.TypeList,
			ConflictsWith: codegen_configs.GetResourceByName("Quota").GetConflictingFields("group_quotas"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,

			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{

					"grace_period": {
						Type:          schema.TypeString,
						ConflictsWith: codegen_configs.GetResourceByName("UserQuota").GetConflictingFields("grace_period"),

						Computed:  true,
						Optional:  true,
						Sensitive: false,

						ValidateDiagFunc: utils.GracePeriodFormatValidation,
					},

					"time_to_block": {
						Type:          schema.TypeString,
						ConflictsWith: codegen_configs.GetResourceByName("UserQuota").GetConflictingFields("time_to_block"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Grace period expiration time`,
					},

					"soft_limit": {
						Type:          schema.TypeInt,
						ConflictsWith: codegen_configs.GetResourceByName("UserQuota").GetConflictingFields("soft_limit"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Soft quota limit`,
					},

					"hard_limit": {
						Type:          schema.TypeInt,
						ConflictsWith: codegen_configs.GetResourceByName("UserQuota").GetConflictingFields("hard_limit"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Hard quota limit`,
					},

					"hard_limit_inodes": {
						Type:          schema.TypeInt,
						ConflictsWith: codegen_configs.GetResourceByName("UserQuota").GetConflictingFields("hard_limit_inodes"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Hard inodes quota limit`,
					},

					"soft_limit_inodes": {
						Type:          schema.TypeInt,
						ConflictsWith: codegen_configs.GetResourceByName("UserQuota").GetConflictingFields("soft_limit_inodes"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Soft inodes quota limit`,
					},

					"used_inodes": {
						Type:          schema.TypeInt,
						ConflictsWith: codegen_configs.GetResourceByName("UserQuota").GetConflictingFields("used_inodes"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Used inodes`,
					},

					"used_capacity": {
						Type:          schema.TypeInt,
						ConflictsWith: codegen_configs.GetResourceByName("UserQuota").GetConflictingFields("used_capacity"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Used capacity in bytes`,
					},

					"is_accountable": {
						Type:          schema.TypeBool,
						ConflictsWith: codegen_configs.GetResourceByName("UserQuota").GetConflictingFields("is_accountable"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,
					},

					"quota_system_id": {
						Type:          schema.TypeInt,
						ConflictsWith: codegen_configs.GetResourceByName("UserQuota").GetConflictingFields("quota_system_id"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,
					},

					"entity": {
						Type:          schema.TypeList,
						ConflictsWith: codegen_configs.GetResourceByName("UserQuota").GetConflictingFields("entity"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,

						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{

								"name": {
									Type:          schema.TypeString,
									ConflictsWith: codegen_configs.GetResourceByName("QuotaEntityInfo").GetConflictingFields("name"),

									Computed:    true,
									Optional:    true,
									Sensitive:   false,
									Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The name of the entity`,
								},

								"vast_id": {
									Type:          schema.TypeInt,
									ConflictsWith: codegen_configs.GetResourceByName("QuotaEntityInfo").GetConflictingFields("vast_id"),

									Computed:    true,
									Optional:    true,
									Sensitive:   false,
									Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,
								},

								"email": {
									Type:          schema.TypeString,
									ConflictsWith: codegen_configs.GetResourceByName("QuotaEntityInfo").GetConflictingFields("email"),

									Computed:    true,
									Optional:    true,
									Sensitive:   false,
									Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,
								},

								"is_group": {
									Type:          schema.TypeBool,
									ConflictsWith: codegen_configs.GetResourceByName("QuotaEntityInfo").GetConflictingFields("is_group"),

									Computed:    true,
									Optional:    true,
									Sensitive:   false,
									Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,
								},

								"identifier": {
									Type:          schema.TypeString,
									ConflictsWith: codegen_configs.GetResourceByName("QuotaEntityInfo").GetConflictingFields("identifier"),

									Required:    true,
									Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,
								},

								"identifier_type": {
									Type:          schema.TypeString,
									ConflictsWith: codegen_configs.GetResourceByName("QuotaEntityInfo").GetConflictingFields("identifier_type"),

									Computed:    true,
									Optional:    true,
									Sensitive:   false,
									Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,
								},
							},
						},
					},
				},
			},
		},
	}
}

var QuotaNamesMapping = map[string][]string{}

func ResourceQuotaReadStructIntoSchema(ctx context.Context, resource api_latest.Quota, d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics
	var err error

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Guid", resource.Guid))

	err = d.Set("guid", resource.Guid)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"guid\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Name", resource.Name))

	err = d.Set("name", resource.Name)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"name\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "State", resource.State))

	err = d.Set("state", resource.State)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"state\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "PrettyState", resource.PrettyState))

	err = d.Set("pretty_state", resource.PrettyState)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"pretty_state\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Path", resource.Path))

	err = d.Set("path", resource.Path)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"path\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "PrettyGracePeriod", resource.PrettyGracePeriod))

	err = d.Set("pretty_grace_period", resource.PrettyGracePeriod)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"pretty_grace_period\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "GracePeriod", resource.GracePeriod))

	err = d.Set("grace_period", resource.GracePeriod)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"grace_period\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "TimeToBlock", resource.TimeToBlock))

	err = d.Set("time_to_block", resource.TimeToBlock)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"time_to_block\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "SoftLimit", resource.SoftLimit))

	err = d.Set("soft_limit", resource.SoftLimit)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"soft_limit\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "HardLimit", resource.HardLimit))

	err = d.Set("hard_limit", resource.HardLimit)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"hard_limit\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "HardLimitInodes", resource.HardLimitInodes))

	err = d.Set("hard_limit_inodes", resource.HardLimitInodes)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"hard_limit_inodes\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "SoftLimitInodes", resource.SoftLimitInodes))

	err = d.Set("soft_limit_inodes", resource.SoftLimitInodes)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"soft_limit_inodes\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "UsedInodes", resource.UsedInodes))

	err = d.Set("used_inodes", resource.UsedInodes)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"used_inodes\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "UsedCapacity", resource.UsedCapacity))

	err = d.Set("used_capacity", resource.UsedCapacity)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"used_capacity\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "UsedCapacityTb", resource.UsedCapacityTb))

	err = d.Set("used_capacity_tb", resource.UsedCapacityTb)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"used_capacity_tb\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "UsedEffectiveCapacity", resource.UsedEffectiveCapacity))

	err = d.Set("used_effective_capacity", resource.UsedEffectiveCapacity)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"used_effective_capacity\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "UsedEffectiveCapacityTb", resource.UsedEffectiveCapacityTb))

	err = d.Set("used_effective_capacity_tb", resource.UsedEffectiveCapacityTb)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"used_effective_capacity_tb\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "TenantId", resource.TenantId))

	err = d.Set("tenant_id", resource.TenantId)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"tenant_id\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "TenantName", resource.TenantName))

	err = d.Set("tenant_name", resource.TenantName)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"tenant_name\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Cluster", resource.Cluster))

	err = d.Set("cluster", resource.Cluster)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"cluster\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "ClusterId", resource.ClusterId))

	err = d.Set("cluster_id", resource.ClusterId)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"cluster_id\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "SystemId", resource.SystemId))

	err = d.Set("system_id", resource.SystemId)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"system_id\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "IsUserQuota", resource.IsUserQuota))

	err = d.Set("is_user_quota", resource.IsUserQuota)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"is_user_quota\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "EnableEmailProviders", resource.EnableEmailProviders))

	err = d.Set("enable_email_providers", resource.EnableEmailProviders)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"enable_email_providers\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "NumExceededUsers", resource.NumExceededUsers))

	err = d.Set("num_exceeded_users", resource.NumExceededUsers)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"num_exceeded_users\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "NumBlockedUsers", resource.NumBlockedUsers))

	err = d.Set("num_blocked_users", resource.NumBlockedUsers)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"num_blocked_users\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "EnableAlarms", resource.EnableAlarms))

	err = d.Set("enable_alarms", resource.EnableAlarms)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"enable_alarms\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "DefaultEmail", resource.DefaultEmail))

	err = d.Set("default_email", resource.DefaultEmail)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"default_email\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "PercentInodes", resource.PercentInodes))

	err = d.Set("percent_inodes", resource.PercentInodes)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"percent_inodes\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "PercentCapacity", resource.PercentCapacity))

	err = d.Set("percent_capacity", resource.PercentCapacity)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"percent_capacity\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "DefaultUserQuota", resource.DefaultUserQuota))

	tflog.Debug(ctx, fmt.Sprintf("Found a pointer object %v", resource.DefaultUserQuota))
	err = d.Set("default_user_quota", utils.FlattenModelAsList(ctx, resource.DefaultUserQuota))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"default_user_quota\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "DefaultGroupQuota", resource.DefaultGroupQuota))

	tflog.Debug(ctx, fmt.Sprintf("Found a pointer object %v", resource.DefaultGroupQuota))
	err = d.Set("default_group_quota", utils.FlattenModelAsList(ctx, resource.DefaultGroupQuota))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"default_group_quota\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "UserQuotas", resource.UserQuotas))

	err = d.Set("user_quotas", utils.FlattenListOfModelsToList(ctx, resource.UserQuotas))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"user_quotas\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "GroupQuotas", resource.GroupQuotas))

	err = d.Set("group_quotas", utils.FlattenListOfModelsToList(ctx, resource.GroupQuotas))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"group_quotas\"",
			Detail:   err.Error(),
		})
	}

	return diags

}
func resourceQuotaRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*vast_client.VMSSession)
	resourceConfig := codegen_configs.GetResourceByName("Quota")
	attrs := map[string]interface{}{"path": utils.GenPath("quotas"), "id": d.Id()}
	tflog.Debug(ctx, fmt.Sprintf("[resourceQuotaRead] Calling Get Function : %v for resource Quota", utils.GetFuncName(resourceConfig.GetFunc)))
	response, err := resourceConfig.GetFunc(ctx, client, attrs, d, map[string]string{})
	utils.VastVersionsWarn(ctx)

	var body []byte
	var resource api_latest.Quota
	if err != nil && response != nil && response.StatusCode == 404 && !resourceConfig.DisableFallbackRequest {
		var fallbackErr error
		body, fallbackErr = utils.HandleFallback(ctx, client, attrs, d, resourceConfig.IdFunc)
		if fallbackErr != nil {
			errorMessage := fmt.Sprintf("Initial request failed:\n%v\nFallback request also failed:\n%v", err.Error(), fallbackErr.Error())
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error occurred while obtaining data from the VAST Data cluster",
				Detail:   errorMessage,
			})
			return diags
		}
	} else if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred while obtaining data from the VAST Data cluster",
			Detail:   err.Error(),
		})
		return diags
	} else {
		tflog.Info(ctx, response.Request.URL.String())
		body, err = resourceConfig.ResponseProcessingFunc(ctx, response)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error occurred reading data received from VAST Data cluster",
				Detail:   err.Error(),
			})
			return diags
		}
	}
	err = json.Unmarshal(body, &resource)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred while parsing data received from VAST Data cluster",
			Detail:   err.Error(),
		})
		return diags

	}
	diags = ResourceQuotaReadStructIntoSchema(ctx, resource, d)

	return diags
}

func resourceQuotaDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*vast_client.VMSSession)
	resourceConfig := codegen_configs.GetResourceByName("Quota")
	attrs := map[string]interface{}{"path": utils.GenPath("quotas"), "id": d.Id()}

	response, err := resourceConfig.DeleteFunc(ctx, client, attrs, nil, map[string]string{})

	tflog.Info(ctx, fmt.Sprintf("Removing Resource"))
	if response != nil {
		tflog.Info(ctx, response.Request.URL.String())
		tflog.Info(ctx, utils.GetResponseBodyAsStr(response))
	}

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred while deleting a resource from the VAST Data cluster",
			Detail:   err.Error(),
		})

	}

	return diags

}

func resourceQuotaCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	namesMapping := utils.ContextKey("names_mapping")
	newCtx := context.WithValue(ctx, namesMapping, QuotaNamesMapping)
	var diags diag.Diagnostics
	data := make(map[string]interface{})
	client := m.(*vast_client.VMSSession)
	resourceConfig := codegen_configs.GetResourceByName("Quota")
	tflog.Info(ctx, fmt.Sprintf("Creating Resource Quota"))
	reflectQuota := reflect.TypeOf((*api_latest.Quota)(nil))
	utils.PopulateResourceMap(newCtx, reflectQuota.Elem(), d, &data, "", false)

	var before_post_error error
	data, before_post_error = resourceConfig.BeforePostFunc(data, client, ctx, d)
	if before_post_error != nil {
		return diag.FromErr(before_post_error)
	}

	versionsEqual := utils.VastVersionsWarn(ctx)

	if versionsEqual != metadata.CLUSTER_VERSION_EQUALS {
		clusterVersion := metadata.ClusterVersionString()
		t, typeExists := vast_versions.GetVersionedType(clusterVersion, "Quota")
		if typeExists {
			versionError := utils.VersionMatch(t, data)
			if versionError != nil {
				tflog.Warn(ctx, versionError.Error())
				versionValidationMode, versionValidationModeExists := metadata.GetClusterConfig("version_validation_mode")
				tflog.Warn(ctx, fmt.Sprintf("Version Validation Mode Detected %s", versionValidationMode))
				if versionValidationModeExists && versionValidationMode == "strict" {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "Cluster Version & Build Version Are Too Different",
						Detail:   versionError.Error(),
					})
					return diags
				}
			}
		} else {
			tflog.Warn(ctx, fmt.Sprintf("Could have not found resource %s in version %s, things might not work properly", "Quota", clusterVersion))
		}
	}
	tflog.Debug(ctx, fmt.Sprintf("Data %v", data))
	b, err := json.MarshalIndent(data, "", "   ")
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Could have not generate request json",
			Detail:   err.Error(),
		})
		return diags
	}
	tflog.Debug(ctx, fmt.Sprintf("Request json created %v", string(b)))
	attrs := map[string]interface{}{"path": utils.GenPath("quotas")}
	response, createErr := resourceConfig.CreateFunc(ctx, client, attrs, data, map[string]string{})
	tflog.Info(ctx, fmt.Sprintf("Server Error for  Quota %v", createErr))

	if createErr != nil {
		errorMessage := fmt.Sprintf("server response:\n%v\nUnderlying error:\n%v", utils.GetResponseBodyAsStr(response), createErr.Error())
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Object Creation Failed",
			Detail:   errorMessage,
		})
		return diags
	}
	responseBody, _ := io.ReadAll(response.Body)
	tflog.Debug(ctx, fmt.Sprintf("Object created, server response %v", string(responseBody)))
	resource := api_latest.Quota{}
	err = json.Unmarshal(responseBody, &resource)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to convert response body into Quota",
			Detail:   err.Error(),
		})
		return diags
	}

	err = resourceConfig.IdFunc(ctx, client, resource.Id, d)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to set Id",
			Detail:   err.Error(),
		})
		return diags
	}
	ctxWithResource := context.WithValue(ctx, utils.ContextKey("resource"), resource)
	resourceQuotaRead(ctxWithResource, d, m)

	return diags
}

func resourceQuotaUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	namesMapping := utils.ContextKey("names_mapping")
	newCtx := context.WithValue(ctx, namesMapping, QuotaNamesMapping)
	var diags diag.Diagnostics
	data := make(map[string]interface{})
	versionsEqual := utils.VastVersionsWarn(ctx)
	resourceConfig := codegen_configs.GetResourceByName("Quota")
	if versionsEqual != metadata.CLUSTER_VERSION_EQUALS {
		clusterVersion := metadata.ClusterVersionString()
		t, typeExists := vast_versions.GetVersionedType(clusterVersion, "Quota")
		if typeExists {
			versionError := utils.VersionMatch(t, data)
			if versionError != nil {
				tflog.Warn(ctx, versionError.Error())
				versionValidationMode, versionValidationModeExists := metadata.GetClusterConfig("version_validation_mode")
				tflog.Warn(ctx, fmt.Sprintf("Version Validation Mode Detected %s", versionValidationMode))
				if versionValidationModeExists && versionValidationMode == "strict" {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "Cluster Version & Build Version Are Too Different",
						Detail:   versionError.Error(),
					})
					return diags
				}
			}
		} else {
			tflog.Warn(ctx, fmt.Sprintf("Could have not found resource %s in version %s, things might not work properly", "Quota", clusterVersion))
		}
	}

	client := m.(*vast_client.VMSSession)
	tflog.Info(ctx, fmt.Sprintf("Updating Resource Quota"))
	reflectQuota := reflect.TypeOf((*api_latest.Quota)(nil))
	utils.PopulateResourceMap(newCtx, reflectQuota.Elem(), d, &data, "", false)

	var beforePatchError error
	data, beforePatchError = resourceConfig.BeforePatchFunc(data, client, ctx, d)
	if beforePatchError != nil {
		return diag.FromErr(beforePatchError)
	}

	tflog.Debug(ctx, fmt.Sprintf("Data %v", data))
	b, err := json.MarshalIndent(data, "", "   ")
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Could have not generate request json",
			Detail:   err.Error(),
		})
		return diags
	}
	tflog.Debug(ctx, fmt.Sprintf("Request json created %v", string(b)))
	attrs := map[string]interface{}{"path": utils.GenPath("quotas"), "id": d.Id()}
	response, patchErr := resourceConfig.UpdateFunc(ctx, client, attrs, data, d, map[string]string{})
	tflog.Info(ctx, fmt.Sprintf("Server Error for  Quota %v", patchErr))
	if patchErr != nil {
		errorMessage := fmt.Sprintf("server response:\n%v\nUnderlying error:\n%v", utils.GetResponseBodyAsStr(response), patchErr.Error())
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Object Creation Failed",
			Detail:   errorMessage,
		})
		return diags
	}
	resourceQuotaRead(ctx, d, m)

	return diags

}

func resourceQuotaImporter(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {

	var result []*schema.ResourceData
	client := m.(*vast_client.VMSSession)
	resourceConfig := codegen_configs.GetResourceByName("Quota")
	attrs := map[string]interface{}{"path": utils.GenPath("quotas")}
	response, err := resourceConfig.ImportFunc(ctx, client, attrs, d, resourceConfig.Importer.GetFunc())

	if err != nil {
		return result, err
	}

	var resourceList []api_latest.Quota
	body, err := resourceConfig.ResponseProcessingFunc(ctx, response)

	if err != nil {
		return result, err
	}

	body, err = utils.ResponseGetByURL(ctx, body, client)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(body, &resourceList)
	if err != nil {
		return result, err
	}

	if len(resourceList) == 0 {
		return result, errors.New("cluster returned 0 elements matching provided guid")
	}

	resource := resourceList[0]
	idErr := resourceConfig.IdFunc(ctx, client, resource.Id, d)
	if idErr != nil {
		return result, idErr
	}

	diags := ResourceQuotaReadStructIntoSchema(ctx, resource, d)
	if diags.HasError() {
		allErrors := "Errors occurred while importing:\n"
		for _, dig := range diags {
			allErrors += fmt.Sprintf("Summary:%s\nDetails:%s\n", dig.Summary, dig.Detail)
		}
		return result, errors.New(allErrors)
	}
	result = append(result, d)

	return result, err

}
