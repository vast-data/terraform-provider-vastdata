package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api_latest "github.com/vast-data/terraform-provider-vastdata/codegen/latest"
	codegen_configs "github.com/vast-data/terraform-provider-vastdata/codegen_tools/configs"
	utils "github.com/vast-data/terraform-provider-vastdata/utils"
	vast_client "github.com/vast-data/terraform-provider-vastdata/vast-client"
	"net/url"
	"strconv"
)

func DataSourceQosPolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceQosPolicyRead,
		Description: ``,
		Schema: map[string]*schema.Schema{

			"id": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,
			},

			"guid": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) QoS Policy guid`,
			},

			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    false,
				Required:    true,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,
			},

			"mode": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) QoS provisioning mode Allowed Values are [STATIC USED_CAPACITY PROVISIONED_CAPACITY]`,
			},

			"policy_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.2.0) The QoS type Allowed Values are [VIEW USER]`,
			},

			"limit_by": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    false,
				Required:    false,
				Optional:    true,
				Description: `(Valid for versions: 5.2.0) What attributes are setting the limitations. Allowed Values are [BW_IOPS BW IOPS]`,
			},

			"tenant_id": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.2.0) When setting is_default this is the tenant which will take affect`,
			},

			"attached_users_identifiers": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.2.0) List of local user IDs to which this QoS Policy is affective.`,

				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"is_default": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.2.0) Should this QoS Policy be the default QoS per user for this tenant ?, tnenat_id should be also provided when settingthis attribute`,
			},

			"io_size_bytes": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Sets the size of IO for static and capacity limit definitions. The number of IOs per request is obtained by dividing request size by IO size. Default: 64K, Recommended range: 4K - 1M`,
			},

			"static_limits": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,

				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"min_reads_bw_mbps": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Minimal amount of performance to provide when there is resource contention`,
						},

						"max_reads_bw_mbps": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Maximal amount of performance to provide when there is no resource contention`,
						},

						"min_writes_bw_mbps": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Minimal amount of performance to provide when there is resource contention`,
						},

						"max_writes_bw_mbps": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Maximal amount of performance to provide when there is no resource contention`,
						},

						"min_reads_iops": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Minimal amount of performance to provide when there is resource contention`,
						},

						"max_reads_iops": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Maximal amount of performance to provide when there is no resource contention`,
						},

						"min_writes_iops": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Minimal amount of performance to provide when there is resource contention`,
						},

						"max_writes_iops": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Maximal amount of performance to provide when there is no resource contention`,
						},

						"burst_reads_bw_mb": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.2.0) Burst reads BW Mb`,
						},

						"burst_reads_loan_mb": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.2.0) Burst reads loan Mb`,
						},

						"burst_writes_bw_mb": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.2.0) Burst writes BW Mb`,
						},

						"burst_writes_loan_mb": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.2.0) Burst writes loan Mb`,
						},

						"burst_reads_iops": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.2.0) Burst reads IOPS`,
						},

						"burst_reads_loan_iops": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.2.0) Burst reads loan IOPS`,
						},

						"burst_writes_iops": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.2.0) Burst writes IOPS`,
						},

						"burst_writes_loan_iops": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.2.0) Burst writes loan IOPS`,
						},
					},
				},
			},

			"capacity_limits": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,

				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"max_reads_bw_mbps_per_gb_capacity": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Maximal amount of performance per GB to provide when there is no resource contention`,
						},

						"max_writes_bw_mbps_per_gb_capacity": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Maximal amount of performance per GB to provide when there is no resource contention`,
						},

						"max_reads_iops_per_gb_capacity": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Maximal amount of performance per GB to provide when there is no resource contention`,
						},

						"max_writes_iops_per_gb_capacity": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Maximal amount of performance per GB to provide when there is no resource contention`,
						},
					},
				},
			},

			"static_total_limits": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.2.0) `,

				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"max_bw_mbps": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.2.0) Maximal BW Mb/s`,
						},

						"burst_bw_mb": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.2.0) Burst BW Mb`,
						},

						"burst_loan_mb": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.2.0) Burst loan Mb`,
						},

						"max_iops": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.2.0) Maximal IOPS`,
						},

						"burst_iops": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.2.0) Burst IOPS`,
						},

						"burst_loan_iops": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.2.0) Burst loan IOPS`,
						},
					},
				},
			},

			"capacity_total_limits": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.2.0) `,

				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"max_bw_mbps_per_gb_capacity": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.2.0) Maximal amount of performance per GB to provide when there is no resource contention`,
						},

						"max_iops_per_gb_capacity": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.2.0) Maximal amount of performance per GB to provide when there is no resource contention`,
						},
					},
				},
			},

			"attached_users": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.2.0) `,

				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"fqdn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.2.0) User FQDN`,
						},

						"is_sid": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.2.0) `,
						},

						"sid_str": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.2.0) The user SID`,
						},

						"uid_or_gid": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.2.0) `,
						},

						"label": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.2.0) How to display the user`,
						},

						"value": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.2.0) The user name`,
						},

						"login_name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.2.0) The user login name`,
						},

						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.2.0) The user name`,
						},

						"identifier_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.2.0) The user type of idetify`,
						},

						"identifier_value": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.2.0) The value to use fo the identifier_type`,
						},
					},
				},
			},
		},
	}
}

func dataSourceQosPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*vast_client.VMSSession)
	values := url.Values{}
	datasource_config := codegen_configs.GetDataSourceByName("QosPolicy")

	name := d.Get("name")
	values.Add("name", fmt.Sprintf("%v", name))

	response, err := client.Get(ctx, utils.GenPath("qospolicies"), values.Encode(), map[string]string{})
	tflog.Info(ctx, response.Request.URL.String())
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured while obtaining data from the vastdata cluster",
			Detail:   err.Error(),
		})
		return diags

	}
	resource_l := []api_latest.QosPolicy{}
	body, err := datasource_config.ResponseProcessingFunc(ctx, response)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured reading data recived from VastData cluster",
			Detail:   err.Error(),
		})
		return diags

	}
	err = json.Unmarshal(body, &resource_l)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured while parsing data recived from VastData cluster",
			Detail:   err.Error(),
		})
		return diags

	}
	if len(resource_l) == 0 {
		d.SetId("")
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Could not find a resource that matches those attributes",
			Detail:   "Could not find a resource that matches those attributes",
		})
		return diags
	}
	if len(resource_l) > 1 {
		d.SetId("")
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Multiple results returned, you might want to add more attributes to get a specific resource",
			Detail:   "Multiple results returned, you might want to add more attributes to get a specific resource",
		})
		return diags
	}

	resource := resource_l[0]

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Id", resource.Id))

	err = d.Set("id", resource.Id)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"id\"",
			Detail:   err.Error(),
		})
	}

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

	Id := (int64)(resource.Id)
	d.SetId(strconv.FormatInt(Id, 10))
	return diags
}
