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

func ResourceQosPolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceQosPolicyRead,
		DeleteContext: resourceQosPolicyDelete,
		CreateContext: resourceQosPolicyCreate,
		UpdateContext: resourceQosPolicyUpdate,

		Importer: &schema.ResourceImporter{
			StateContext: resourceQosPolicyImporter,
		},

		Description: ``,
		Schema:      getResourceQosPolicySchema(),
	}
}

func getResourceQosPolicySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{

		"guid": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("QosPolicy").GetConflictingFields("guid"),

			Computed:    true,
			Optional:    false,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) QoS Policy guid`,
		},

		"name": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("QosPolicy").GetConflictingFields("name"),

			Required:    true,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,
		},

		"mode": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("QosPolicy").GetConflictingFields("mode"),

			Computed:  true,
			Optional:  true,
			Sensitive: false,

			ValidateDiagFunc: utils.OneOf([]string{"STATIC", "USED_CAPACITY", "PROVISIONED_CAPACITY"}),
			Description:      `(Valid for versions: 5.0.0,5.1.0,5.2.0) QoS provisioning mode Allowed Values are [STATIC USED_CAPACITY PROVISIONED_CAPACITY]`,
		},

		"policy_type": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("QosPolicy").GetConflictingFields("policy_type"),

			Computed:  true,
			Optional:  true,
			Sensitive: false,

			ValidateDiagFunc: utils.OneOf([]string{"VIEW", "USER"}),
			Description:      `(Valid for versions: 5.2.0) The QoS type Allowed Values are [VIEW USER]`,
		},

		"limit_by": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("QosPolicy").GetConflictingFields("limit_by"),

			Computed:  false,
			Optional:  true,
			Sensitive: false,

			ValidateDiagFunc: utils.OneOf([]string{"BW_IOPS", "BW", "IOPS"}),
			Description:      `(Valid for versions: 5.2.0) What attributes are setting the limitations. Allowed Values are [BW_IOPS BW IOPS]`,

			Default: "BW_IOPS",
		},

		"tenant_id": &schema.Schema{
			Type:          schema.TypeInt,
			ConflictsWith: codegen_configs.GetResourceByName("QosPolicy").GetConflictingFields("tenant_id"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.2.0) When setting is_default this is the tenant which will take affect`,
		},

		"attached_users_identifiers": &schema.Schema{
			Type:          schema.TypeList,
			ConflictsWith: codegen_configs.GetResourceByName("QosPolicy").GetConflictingFields("attached_users_identifiers"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.2.0) List of local user IDs to which this QoS Policy is affective.`,

			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},

		"is_default": &schema.Schema{
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("QosPolicy").GetConflictingFields("is_default"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.2.0) Should this QoS Policy be the default QoS per user for this tenant ?, tnenat_id should be also provided when settingthis attribute`,
		},

		"io_size_bytes": &schema.Schema{
			Type:          schema.TypeInt,
			ConflictsWith: codegen_configs.GetResourceByName("QosPolicy").GetConflictingFields("io_size_bytes"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Sets the size of IO for static and capacity limit definitions. The number of IOs per request is obtained by dividing request size by IO size. Default: 64K, Recommended range: 4K - 1M`,
		},

		"static_limits": &schema.Schema{
			Type:          schema.TypeList,
			ConflictsWith: codegen_configs.GetResourceByName("QosPolicy").GetConflictingFields("static_limits"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,

			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{

					"min_reads_bw_mbps": &schema.Schema{
						Type:          schema.TypeInt,
						ConflictsWith: codegen_configs.GetResourceByName("QosStaticLimits").GetConflictingFields("min_reads_bw_mbps"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Minimal amount of performance to provide when there is resource contention`,
					},

					"max_reads_bw_mbps": &schema.Schema{
						Type:          schema.TypeInt,
						ConflictsWith: codegen_configs.GetResourceByName("QosStaticLimits").GetConflictingFields("max_reads_bw_mbps"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Maximal amount of performance to provide when there is no resource contention`,
					},

					"min_writes_bw_mbps": &schema.Schema{
						Type:          schema.TypeInt,
						ConflictsWith: codegen_configs.GetResourceByName("QosStaticLimits").GetConflictingFields("min_writes_bw_mbps"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Minimal amount of performance to provide when there is resource contention`,
					},

					"max_writes_bw_mbps": &schema.Schema{
						Type:          schema.TypeInt,
						ConflictsWith: codegen_configs.GetResourceByName("QosStaticLimits").GetConflictingFields("max_writes_bw_mbps"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Maximal amount of performance to provide when there is no resource contention`,
					},

					"min_reads_iops": &schema.Schema{
						Type:          schema.TypeInt,
						ConflictsWith: codegen_configs.GetResourceByName("QosStaticLimits").GetConflictingFields("min_reads_iops"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Minimal amount of performance to provide when there is resource contention`,
					},

					"max_reads_iops": &schema.Schema{
						Type:          schema.TypeInt,
						ConflictsWith: codegen_configs.GetResourceByName("QosStaticLimits").GetConflictingFields("max_reads_iops"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Maximal amount of performance to provide when there is no resource contention`,
					},

					"min_writes_iops": &schema.Schema{
						Type:          schema.TypeInt,
						ConflictsWith: codegen_configs.GetResourceByName("QosStaticLimits").GetConflictingFields("min_writes_iops"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Minimal amount of performance to provide when there is resource contention`,
					},

					"max_writes_iops": &schema.Schema{
						Type:          schema.TypeInt,
						ConflictsWith: codegen_configs.GetResourceByName("QosStaticLimits").GetConflictingFields("max_writes_iops"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Maximal amount of performance to provide when there is no resource contention`,
					},

					"burst_reads_bw_mb": &schema.Schema{
						Type:          schema.TypeInt,
						ConflictsWith: codegen_configs.GetResourceByName("QosStaticLimits").GetConflictingFields("burst_reads_bw_mb"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.2.0) Burst reads BW Mb`,
					},

					"burst_reads_loan_mb": &schema.Schema{
						Type:          schema.TypeInt,
						ConflictsWith: codegen_configs.GetResourceByName("QosStaticLimits").GetConflictingFields("burst_reads_loan_mb"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.2.0) Burst reads loan Mb`,
					},

					"burst_writes_bw_mb": &schema.Schema{
						Type:          schema.TypeInt,
						ConflictsWith: codegen_configs.GetResourceByName("QosStaticLimits").GetConflictingFields("burst_writes_bw_mb"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.2.0) Burst writes BW Mb`,
					},

					"burst_writes_loan_mb": &schema.Schema{
						Type:          schema.TypeInt,
						ConflictsWith: codegen_configs.GetResourceByName("QosStaticLimits").GetConflictingFields("burst_writes_loan_mb"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.2.0) Burst writes loan Mb`,
					},

					"burst_reads_iops": &schema.Schema{
						Type:          schema.TypeInt,
						ConflictsWith: codegen_configs.GetResourceByName("QosStaticLimits").GetConflictingFields("burst_reads_iops"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.2.0) Burst reads IOPS`,
					},

					"burst_reads_loan_iops": &schema.Schema{
						Type:          schema.TypeInt,
						ConflictsWith: codegen_configs.GetResourceByName("QosStaticLimits").GetConflictingFields("burst_reads_loan_iops"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.2.0) Burst reads loan IOPS`,
					},

					"burst_writes_iops": &schema.Schema{
						Type:          schema.TypeInt,
						ConflictsWith: codegen_configs.GetResourceByName("QosStaticLimits").GetConflictingFields("burst_writes_iops"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.2.0) Burst writes IOPS`,
					},

					"burst_writes_loan_iops": &schema.Schema{
						Type:          schema.TypeInt,
						ConflictsWith: codegen_configs.GetResourceByName("QosStaticLimits").GetConflictingFields("burst_writes_loan_iops"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.2.0) Burst writes loan IOPS`,
					},
				},
			},
		},

		"capacity_limits": &schema.Schema{
			Type:          schema.TypeList,
			ConflictsWith: codegen_configs.GetResourceByName("QosPolicy").GetConflictingFields("capacity_limits"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,

			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{

					"max_reads_bw_mbps_per_gb_capacity": &schema.Schema{
						Type:          schema.TypeInt,
						ConflictsWith: codegen_configs.GetResourceByName("QosDynamicLimits").GetConflictingFields("max_reads_bw_mbps_per_gb_capacity"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Maximal amount of performance per GB to provide when there is no resource contention`,
					},

					"max_writes_bw_mbps_per_gb_capacity": &schema.Schema{
						Type:          schema.TypeInt,
						ConflictsWith: codegen_configs.GetResourceByName("QosDynamicLimits").GetConflictingFields("max_writes_bw_mbps_per_gb_capacity"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Maximal amount of performance per GB to provide when there is no resource contention`,
					},

					"max_reads_iops_per_gb_capacity": &schema.Schema{
						Type:          schema.TypeInt,
						ConflictsWith: codegen_configs.GetResourceByName("QosDynamicLimits").GetConflictingFields("max_reads_iops_per_gb_capacity"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Maximal amount of performance per GB to provide when there is no resource contention`,
					},

					"max_writes_iops_per_gb_capacity": &schema.Schema{
						Type:          schema.TypeInt,
						ConflictsWith: codegen_configs.GetResourceByName("QosDynamicLimits").GetConflictingFields("max_writes_iops_per_gb_capacity"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Maximal amount of performance per GB to provide when there is no resource contention`,
					},
				},
			},
		},

		"static_total_limits": &schema.Schema{
			Type:          schema.TypeList,
			ConflictsWith: codegen_configs.GetResourceByName("QosPolicy").GetConflictingFields("static_total_limits"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.2.0) `,

			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{

					"max_bw_mbps": &schema.Schema{
						Type:          schema.TypeInt,
						ConflictsWith: codegen_configs.GetResourceByName("QoSStaticTotalLimits").GetConflictingFields("max_bw_mbps"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.2.0) Maximal BW Mb/s`,
					},

					"burst_bw_mb": &schema.Schema{
						Type:          schema.TypeInt,
						ConflictsWith: codegen_configs.GetResourceByName("QoSStaticTotalLimits").GetConflictingFields("burst_bw_mb"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.2.0) Burst BW Mb`,
					},

					"burst_loan_mb": &schema.Schema{
						Type:          schema.TypeInt,
						ConflictsWith: codegen_configs.GetResourceByName("QoSStaticTotalLimits").GetConflictingFields("burst_loan_mb"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.2.0) Burst loan Mb`,
					},

					"max_iops": &schema.Schema{
						Type:          schema.TypeInt,
						ConflictsWith: codegen_configs.GetResourceByName("QoSStaticTotalLimits").GetConflictingFields("max_iops"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.2.0) Maximal IOPS`,
					},

					"burst_iops": &schema.Schema{
						Type:          schema.TypeInt,
						ConflictsWith: codegen_configs.GetResourceByName("QoSStaticTotalLimits").GetConflictingFields("burst_iops"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.2.0) Burst IOPS`,
					},

					"burst_loan_iops": &schema.Schema{
						Type:          schema.TypeInt,
						ConflictsWith: codegen_configs.GetResourceByName("QoSStaticTotalLimits").GetConflictingFields("burst_loan_iops"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.2.0) Burst loan IOPS`,
					},
				},
			},
		},

		"capacity_total_limits": &schema.Schema{
			Type:          schema.TypeList,
			ConflictsWith: codegen_configs.GetResourceByName("QosPolicy").GetConflictingFields("capacity_total_limits"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.2.0) `,

			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{

					"max_bw_mbps_per_gb_capacity": &schema.Schema{
						Type:          schema.TypeInt,
						ConflictsWith: codegen_configs.GetResourceByName("QoSDynamicTotalLimits").GetConflictingFields("max_bw_mbps_per_gb_capacity"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.2.0) Maximal amount of performance per GB to provide when there is no resource contention`,
					},

					"max_iops_per_gb_capacity": &schema.Schema{
						Type:          schema.TypeInt,
						ConflictsWith: codegen_configs.GetResourceByName("QoSDynamicTotalLimits").GetConflictingFields("max_iops_per_gb_capacity"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.2.0) Maximal amount of performance per GB to provide when there is no resource contention`,
					},
				},
			},
		},

		"attached_users": &schema.Schema{
			Type:          schema.TypeList,
			ConflictsWith: codegen_configs.GetResourceByName("QosPolicy").GetConflictingFields("attached_users"),

			Computed:    true,
			Optional:    false,
			Sensitive:   false,
			Description: `(Valid for versions: 5.2.0) `,

			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{

					"fqdn": &schema.Schema{
						Type:          schema.TypeString,
						ConflictsWith: codegen_configs.GetResourceByName("QosUser").GetConflictingFields("fqdn"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.2.0) User FQDN`,
					},

					"is_sid": &schema.Schema{
						Type:          schema.TypeBool,
						ConflictsWith: codegen_configs.GetResourceByName("QosUser").GetConflictingFields("is_sid"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.2.0) `,
					},

					"sid_str": &schema.Schema{
						Type:          schema.TypeString,
						ConflictsWith: codegen_configs.GetResourceByName("QosUser").GetConflictingFields("sid_str"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.2.0) The user SID`,
					},

					"uid_or_gid": &schema.Schema{
						Type:          schema.TypeInt,
						ConflictsWith: codegen_configs.GetResourceByName("QosUser").GetConflictingFields("uid_or_gid"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.2.0) `,
					},

					"label": &schema.Schema{
						Type:          schema.TypeString,
						ConflictsWith: codegen_configs.GetResourceByName("QosUser").GetConflictingFields("label"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.2.0) How to display the user`,
					},

					"value": &schema.Schema{
						Type:          schema.TypeString,
						ConflictsWith: codegen_configs.GetResourceByName("QosUser").GetConflictingFields("value"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.2.0) The user name`,
					},

					"login_name": &schema.Schema{
						Type:          schema.TypeString,
						ConflictsWith: codegen_configs.GetResourceByName("QosUser").GetConflictingFields("login_name"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.2.0) The user login name`,
					},

					"name": &schema.Schema{
						Type:          schema.TypeString,
						ConflictsWith: codegen_configs.GetResourceByName("QosUser").GetConflictingFields("name"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.2.0) The user name`,
					},

					"identifier_type": &schema.Schema{
						Type:          schema.TypeString,
						ConflictsWith: codegen_configs.GetResourceByName("QosUser").GetConflictingFields("identifier_type"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.2.0) The user type of idetify`,
					},

					"identifier_value": &schema.Schema{
						Type:          schema.TypeString,
						ConflictsWith: codegen_configs.GetResourceByName("QosUser").GetConflictingFields("identifier_value"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.2.0) The value to use fo the identifier_type`,
					},
				},
			},
		},
	}
}

var QosPolicy_names_mapping map[string][]string = map[string][]string{}

func ResourceQosPolicyReadStructIntoSchema(ctx context.Context, resource api_latest.QosPolicy, d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics
	var err error

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Guid", resource.Guid))

	err = d.Set("guid", resource.Guid)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"guid\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Name", resource.Name))

	err = d.Set("name", resource.Name)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"name\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Mode", resource.Mode))

	err = d.Set("mode", resource.Mode)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"mode\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "PolicyType", resource.PolicyType))

	err = d.Set("policy_type", resource.PolicyType)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"policy_type\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "LimitBy", resource.LimitBy))

	err = d.Set("limit_by", resource.LimitBy)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"limit_by\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "TenantId", resource.TenantId))

	err = d.Set("tenant_id", resource.TenantId)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"tenant_id\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "AttachedUsersIdentifiers", resource.AttachedUsersIdentifiers))

	err = d.Set("attached_users_identifiers", utils.FlattenListOfPrimitives(&resource.AttachedUsersIdentifiers))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"attached_users_identifiers\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "IsDefault", resource.IsDefault))

	err = d.Set("is_default", resource.IsDefault)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"is_default\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "IoSizeBytes", resource.IoSizeBytes))

	err = d.Set("io_size_bytes", resource.IoSizeBytes)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"io_size_bytes\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "StaticLimits", resource.StaticLimits))

	tflog.Debug(ctx, fmt.Sprintf("Found a pointer object %v", resource.StaticLimits))
	err = d.Set("static_limits", utils.FlattenModelAsList(ctx, resource.StaticLimits))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"static_limits\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "CapacityLimits", resource.CapacityLimits))

	tflog.Debug(ctx, fmt.Sprintf("Found a pointer object %v", resource.CapacityLimits))
	err = d.Set("capacity_limits", utils.FlattenModelAsList(ctx, resource.CapacityLimits))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"capacity_limits\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "StaticTotalLimits", resource.StaticTotalLimits))

	tflog.Debug(ctx, fmt.Sprintf("Found a pointer object %v", resource.StaticTotalLimits))
	err = d.Set("static_total_limits", utils.FlattenModelAsList(ctx, resource.StaticTotalLimits))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"static_total_limits\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "CapacityTotalLimits", resource.CapacityTotalLimits))

	tflog.Debug(ctx, fmt.Sprintf("Found a pointer object %v", resource.CapacityTotalLimits))
	err = d.Set("capacity_total_limits", utils.FlattenModelAsList(ctx, resource.CapacityTotalLimits))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"capacity_total_limits\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "AttachedUsers", resource.AttachedUsers))

	err = d.Set("attached_users", utils.FlattenListOfModelsToList(ctx, resource.AttachedUsers))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"attached_users\"",
			Detail:   err.Error(),
		})
	}

	return diags

}
func resourceQosPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*vast_client.VMSSession)
	resource_config := codegen_configs.GetResourceByName("QosPolicy")
	attrs := map[string]interface{}{"path": utils.GenPath("qospolicies"), "id": d.Id()}
	tflog.Debug(ctx, fmt.Sprintf("[resourceQosPolicyRead] Calling Get Function : %v for resource QosPolicy", utils.GetFuncName(resource_config.GetFunc)))
	response, err := resource_config.GetFunc(ctx, client, attrs, d, map[string]string{})
	utils.VastVersionsWarn(ctx)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured while obtaining data from the vastdata cluster",
			Detail:   err.Error(),
		})
		return diags

	}
	tflog.Info(ctx, response.Request.URL.String())
	resource := api_latest.QosPolicy{}
	body, err := resource_config.ResponseProcessingFunc(ctx, response)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured reading data recived from VastData cluster",
			Detail:   err.Error(),
		})
		return diags

	}
	err = json.Unmarshal(body, &resource)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured while parsing data recived from VastData cluster",
			Detail:   err.Error(),
		})
		return diags

	}
	diags = ResourceQosPolicyReadStructIntoSchema(ctx, resource, d)

	var after_read_error error
	after_read_error = resource_config.AfterReadFunc(client, ctx, d)
	if after_read_error != nil {
		return diag.FromErr(after_read_error)
	}

	return diags
}

func resourceQosPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*vast_client.VMSSession)
	resource_config := codegen_configs.GetResourceByName("QosPolicy")
	attrs := map[string]interface{}{"path": utils.GenPath("qospolicies"), "id": d.Id()}

	response, err := resource_config.DeleteFunc(ctx, client, attrs, nil, map[string]string{})

	tflog.Info(ctx, fmt.Sprintf("Removing Resource"))
	if response != nil {
		tflog.Info(ctx, response.Request.URL.String())
		tflog.Info(ctx, utils.GetResponseBodyAsStr(response))
	}

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured while deleting a resource from the vastdata cluster",
			Detail:   err.Error(),
		})

	}

	return diags

}

func resourceQosPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	names_mapping := utils.ContextKey("names_mapping")
	new_ctx := context.WithValue(ctx, names_mapping, QosPolicy_names_mapping)
	var diags diag.Diagnostics
	data := make(map[string]interface{})
	client := m.(*vast_client.VMSSession)
	resource_config := codegen_configs.GetResourceByName("QosPolicy")
	tflog.Info(ctx, fmt.Sprintf("Creating Resource QosPolicy"))
	reflect_QosPolicy := reflect.TypeOf((*api_latest.QosPolicy)(nil))
	utils.PopulateResourceMap(new_ctx, reflect_QosPolicy.Elem(), d, &data, "", false)

	var before_post_error error
	data, before_post_error = resource_config.BeforePostFunc(data, client, ctx, d)
	if before_post_error != nil {
		return diag.FromErr(before_post_error)
	}

	version_compare := utils.VastVersionsWarn(ctx)

	if version_compare != metadata.CLUSTER_VERSION_EQUALS {
		cluster_version := metadata.ClusterVersionString()
		t, t_exists := vast_versions.GetVersionedType(cluster_version, "QosPolicy")
		if t_exists {
			versions_error := utils.VersionMatch(t, data)
			if versions_error != nil {
				tflog.Warn(ctx, versions_error.Error())
				version_validation_mode, version_validation_mode_exists := metadata.GetClusterConfig("version_validation_mode")
				tflog.Warn(ctx, fmt.Sprintf("Version Validation Mode Detected %s", version_validation_mode))
				if version_validation_mode_exists && version_validation_mode == "strict" {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "Cluster Version & Build Version Are Too Differant",
						Detail:   versions_error.Error(),
					})
					return diags
				}
			}
		} else {
			tflog.Warn(ctx, fmt.Sprintf("Could have not found resource %s in version %s , things might not work properly", "QosPolicy", cluster_version))
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
	attrs := map[string]interface{}{"path": utils.GenPath("qospolicies")}
	response, create_err := resource_config.CreateFunc(ctx, client, attrs, data, map[string]string{})
	tflog.Info(ctx, fmt.Sprintf("Server Error for  QosPolicy %v", create_err))

	if create_err != nil {
		error_message := create_err.Error() + " Server Response: " + utils.GetResponseBodyAsStr(response)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Object Creation Failed",
			Detail:   error_message,
		})
		return diags
	}
	response_body, _ := io.ReadAll(response.Body)
	tflog.Debug(ctx, fmt.Sprintf("Object created , server response %v", string(response_body)))
	resource := api_latest.QosPolicy{}
	err = json.Unmarshal(response_body, &resource)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to convert response body into QosPolicy",
			Detail:   err.Error(),
		})
		return diags
	}

	err = resource_config.IdFunc(ctx, client, resource.Id, d)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to set Id",
			Detail:   err.Error(),
		})
		return diags
	}
	ctx_with_resource := context.WithValue(ctx, utils.ContextKey("resource"), resource)
	resourceQosPolicyRead(ctx_with_resource, d, m)

	return diags
}

func resourceQosPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	names_mapping := utils.ContextKey("names_mapping")
	new_ctx := context.WithValue(ctx, names_mapping, QosPolicy_names_mapping)
	var diags diag.Diagnostics
	data := make(map[string]interface{})
	version_compare := utils.VastVersionsWarn(ctx)
	resource_config := codegen_configs.GetResourceByName("QosPolicy")
	if version_compare != metadata.CLUSTER_VERSION_EQUALS {
		cluster_version := metadata.ClusterVersionString()
		t, t_exists := vast_versions.GetVersionedType(cluster_version, "QosPolicy")
		if t_exists {
			versions_error := utils.VersionMatch(t, data)
			if versions_error != nil {
				tflog.Warn(ctx, versions_error.Error())
				version_validation_mode, version_validation_mode_exists := metadata.GetClusterConfig("version_validation_mode")
				tflog.Warn(ctx, fmt.Sprintf("Version Validation Mode Detected %s", version_validation_mode))
				if version_validation_mode_exists && version_validation_mode == "strict" {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "Cluster Version & Build Version Are Too Differant",
						Detail:   versions_error.Error(),
					})
					return diags
				}
			}
		} else {
			tflog.Warn(ctx, fmt.Sprintf("Could have not found resource %s in version %s , things might not work properly", "QosPolicy", cluster_version))
		}
	}

	client := m.(*vast_client.VMSSession)
	tflog.Info(ctx, fmt.Sprintf("Updating Resource QosPolicy"))
	reflect_QosPolicy := reflect.TypeOf((*api_latest.QosPolicy)(nil))
	utils.PopulateResourceMap(new_ctx, reflect_QosPolicy.Elem(), d, &data, "", false)

	var before_patch_error error
	data, before_patch_error = resource_config.BeforePatchFunc(data, client, ctx, d)
	if before_patch_error != nil {
		return diag.FromErr(before_patch_error)
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
	attrs := map[string]interface{}{"path": utils.GenPath("qospolicies"), "id": d.Id()}
	response, patch_err := resource_config.UpdateFunc(ctx, client, attrs, data, d, map[string]string{})
	tflog.Info(ctx, fmt.Sprintf("Server Error for  QosPolicy %v", patch_err))
	if patch_err != nil {
		error_message := patch_err.Error() + " Server Response: " + utils.GetResponseBodyAsStr(response)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Object Creation Failed",
			Detail:   error_message,
		})
		return diags
	}
	resourceQosPolicyRead(ctx, d, m)

	return diags

}

func resourceQosPolicyImporter(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {

	result := []*schema.ResourceData{}
	client := m.(*vast_client.VMSSession)
	resource_config := codegen_configs.GetResourceByName("QosPolicy")
	attrs := map[string]interface{}{"path": utils.GenPath("qospolicies")}
	response, err := resource_config.ImportFunc(ctx, client, attrs, d, resource_config.Importer.GetFunc())

	if err != nil {
		return result, err
	}

	resource_l := []api_latest.QosPolicy{}
	body, err := resource_config.ResponseProcessingFunc(ctx, response)

	if err != nil {
		return result, err
	}
	err = json.Unmarshal(body, &resource_l)
	if err != nil {
		return result, err
	}

	if len(resource_l) == 0 {
		return result, errors.New("Cluster provided 0 elements matchin gthis guid")
	}

	resource := resource_l[0]
	id_err := resource_config.IdFunc(ctx, client, resource.Id, d)
	if id_err != nil {
		return result, id_err
	}

	diags := ResourceQosPolicyReadStructIntoSchema(ctx, resource, d)
	if diags.HasError() {
		all_errors := "Errors occured while importing:\n"
		for _, dig := range diags {
			all_errors += fmt.Sprintf("Summary:%s\nDetails:%s\n", dig.Summary, dig.Detail)
		}
		return result, errors.New(all_errors)
	}
	result = append(result, d)

	return result, err

}
