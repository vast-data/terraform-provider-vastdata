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

func DataSourceQuota() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceQuotaRead,
		Description: `This is a quota`,
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
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Quota guid`,
			},

			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    false,
				Required:    true,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The name`,
			},

			"state": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,
			},

			"pretty_state": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,
			},

			"path": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Directory path`,
			},

			"pretty_grace_period": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Quota enforcement pretty grace period in seconds, minutes, hours or days. Example: 90m`,
			},

			"grace_period": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Quota enforcement grace period in seconds, minutes, hours or days. Example: 90m`,
			},

			"time_to_block": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Grace period expiration time`,
			},

			"soft_limit": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Soft quota limit`,
			},

			"hard_limit": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Hard quota limit`,
			},

			"hard_limit_inodes": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Hard inodes quota limit`,
			},

			"soft_limit_inodes": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Soft inodes quota limit`,
			},

			"used_inodes": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Used inodes`,
			},

			"used_capacity": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Used capacity in bytes`,
			},

			"used_capacity_tb": &schema.Schema{
				Type:        schema.TypeFloat,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Used capacity in TB`,
			},

			"used_effective_capacity": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Used effective capacity in bytes`,
			},

			"used_effective_capacity_tb": &schema.Schema{
				Type:        schema.TypeFloat,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Used effective capacity in TB`,
			},

			"tenant_id": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Required:    false,
				Optional:    true,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Tenant ID`,
			},

			"tenant_name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Tenant Name`,
			},

			"cluster": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Parent Cluster`,
			},

			"cluster_id": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Parent Cluster ID`,
			},

			"system_id": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,
			},

			"is_user_quota": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,
			},

			"enable_email_providers": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,
			},

			"num_exceeded_users": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,
			},

			"num_blocked_users": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,
			},

			"enable_alarms": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Enable alarms when users or groups are exceeding their limit`,
			},

			"default_email": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The default Email if there is no suffix and no address in the providers`,
			},

			"percent_inodes": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Percent of used inodes out of the hard limit`,
			},

			"percent_capacity": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Percent of used capacity out of the hard limit`,
			},

			"default_user_quota": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,

				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"quota_system_id": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The system ID of the quota`,
						},

						"soft_limit": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The size soft limit in bytes`,
						},

						"hard_limit": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The size hard limit in bytes`,
						},

						"sof_limit_inodes": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The sof limit of inodes number`,
						},

						"hard_limit_inodes": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The hard limit in inode number`,
						},

						"grace_period": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Quota enforcement grace period at the format of HH:MM:SS`,
						},
					},
				},
			},

			"default_group_quota": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,

				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"quota_system_id": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The system ID of the quota`,
						},

						"soft_limit": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The size soft limit in bytes`,
						},

						"hard_limit": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The size hard limit in bytes`,
						},

						"sof_limit_inodes": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The sof limit of inodes number`,
						},

						"hard_limit_inodes": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The hard limit in inode number`,
						},

						"grace_period": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Quota enforcement grace period at the format of HH:MM:SS`,
						},
					},
				},
			},

			"user_quotas": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,

				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"grace_period": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Quota enforcement grace period at the format of HH:MM:SS`,
						},

						"time_to_block": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Grace period expiration time`,
						},

						"soft_limit": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Soft quota limit`,
						},

						"hard_limit": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Hard quota limit`,
						},

						"hard_limit_inodes": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Hard inodes quota limit`,
						},

						"soft_limit_inodes": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Soft inodes quota limit`,
						},

						"used_inodes": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Used inodes`,
						},

						"used_capacity": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Used capacity in bytes`,
						},

						"is_accountable": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,
						},

						"quota_system_id": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,
						},

						"entity": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,

							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{

									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    false,
										Required:    true,
										Optional:    false,
										Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The name of the entity`,
									},

									"vast_id": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Required:    false,
										Optional:    false,
										Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,
									},

									"email": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Required:    false,
										Optional:    false,
										Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,
									},

									"is_group": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Required:    false,
										Optional:    false,
										Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,
									},

									"identifier": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Required:    false,
										Optional:    false,
										Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,
									},

									"identifier_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Required:    false,
										Optional:    false,
										Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,
									},
								},
							},
						},
					},
				},
			},

			"group_quotas": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,

				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"grace_period": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Quota enforcement grace period at the format of HH:MM:SS`,
						},

						"time_to_block": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Grace period expiration time`,
						},

						"soft_limit": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Soft quota limit`,
						},

						"hard_limit": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Hard quota limit`,
						},

						"hard_limit_inodes": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Hard inodes quota limit`,
						},

						"soft_limit_inodes": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Soft inodes quota limit`,
						},

						"used_inodes": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Used inodes`,
						},

						"used_capacity": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Used capacity in bytes`,
						},

						"is_accountable": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,
						},

						"quota_system_id": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,
						},

						"entity": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,

							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{

									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    false,
										Required:    true,
										Optional:    false,
										Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The name of the entity`,
									},

									"vast_id": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Required:    false,
										Optional:    false,
										Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,
									},

									"email": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Required:    false,
										Optional:    false,
										Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,
									},

									"is_group": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Required:    false,
										Optional:    false,
										Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,
									},

									"identifier": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Required:    false,
										Optional:    false,
										Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,
									},

									"identifier_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Required:    false,
										Optional:    false,
										Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceQuotaRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*vast_client.VMSSession)
	values := url.Values{}
	datasource_config := codegen_configs.GetDataSourceByName("Quota")

	name := d.Get("name")
	values.Add("name", fmt.Sprintf("%v", name))

	if d.HasChanges("tenant_id") {
		tenant_id := d.Get("tenant_id")
		tflog.Debug(ctx, "Using optional attribute \"tenant_id\"")
		values.Add("tenant_id", fmt.Sprintf("%v", tenant_id))
	}

	response, err := client.Get(ctx, utils.GenPath("quotas"), values.Encode(), map[string]string{})
	tflog.Info(ctx, response.Request.URL.String())
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured while obtaining data from the vastdata cluster",
			Detail:   err.Error(),
		})
		return diags

	}
	resource_l := []api_latest.Quota{}
	body, err := datasource_config.ResponseProcessingFunc(ctx, response)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured reading data recived from VastData cluster",
			Detail:   err.Error(),
		})
		return diags

	}

	body, err = utils.ResponseGetByURL(ctx, body, client)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured reading urls from response",
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

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "State", resource.State))

	err = d.Set("state", resource.State)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"state\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "PrettyState", resource.PrettyState))

	err = d.Set("pretty_state", resource.PrettyState)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"pretty_state\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Path", resource.Path))

	err = d.Set("path", resource.Path)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"path\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "PrettyGracePeriod", resource.PrettyGracePeriod))

	err = d.Set("pretty_grace_period", resource.PrettyGracePeriod)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"pretty_grace_period\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "GracePeriod", resource.GracePeriod))

	err = d.Set("grace_period", resource.GracePeriod)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"grace_period\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "TimeToBlock", resource.TimeToBlock))

	err = d.Set("time_to_block", resource.TimeToBlock)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"time_to_block\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "SoftLimit", resource.SoftLimit))

	err = d.Set("soft_limit", resource.SoftLimit)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"soft_limit\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "HardLimit", resource.HardLimit))

	err = d.Set("hard_limit", resource.HardLimit)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"hard_limit\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "HardLimitInodes", resource.HardLimitInodes))

	err = d.Set("hard_limit_inodes", resource.HardLimitInodes)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"hard_limit_inodes\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "SoftLimitInodes", resource.SoftLimitInodes))

	err = d.Set("soft_limit_inodes", resource.SoftLimitInodes)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"soft_limit_inodes\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "UsedInodes", resource.UsedInodes))

	err = d.Set("used_inodes", resource.UsedInodes)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"used_inodes\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "UsedCapacity", resource.UsedCapacity))

	err = d.Set("used_capacity", resource.UsedCapacity)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"used_capacity\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "UsedCapacityTb", resource.UsedCapacityTb))

	err = d.Set("used_capacity_tb", resource.UsedCapacityTb)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"used_capacity_tb\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "UsedEffectiveCapacity", resource.UsedEffectiveCapacity))

	err = d.Set("used_effective_capacity", resource.UsedEffectiveCapacity)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"used_effective_capacity\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "UsedEffectiveCapacityTb", resource.UsedEffectiveCapacityTb))

	err = d.Set("used_effective_capacity_tb", resource.UsedEffectiveCapacityTb)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"used_effective_capacity_tb\"",
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

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "TenantName", resource.TenantName))

	err = d.Set("tenant_name", resource.TenantName)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"tenant_name\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Cluster", resource.Cluster))

	err = d.Set("cluster", resource.Cluster)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"cluster\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "ClusterId", resource.ClusterId))

	err = d.Set("cluster_id", resource.ClusterId)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"cluster_id\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "SystemId", resource.SystemId))

	err = d.Set("system_id", resource.SystemId)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"system_id\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "IsUserQuota", resource.IsUserQuota))

	err = d.Set("is_user_quota", resource.IsUserQuota)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"is_user_quota\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "EnableEmailProviders", resource.EnableEmailProviders))

	err = d.Set("enable_email_providers", resource.EnableEmailProviders)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"enable_email_providers\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "NumExceededUsers", resource.NumExceededUsers))

	err = d.Set("num_exceeded_users", resource.NumExceededUsers)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"num_exceeded_users\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "NumBlockedUsers", resource.NumBlockedUsers))

	err = d.Set("num_blocked_users", resource.NumBlockedUsers)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"num_blocked_users\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "EnableAlarms", resource.EnableAlarms))

	err = d.Set("enable_alarms", resource.EnableAlarms)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"enable_alarms\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "DefaultEmail", resource.DefaultEmail))

	err = d.Set("default_email", resource.DefaultEmail)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"default_email\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "PercentInodes", resource.PercentInodes))

	err = d.Set("percent_inodes", resource.PercentInodes)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"percent_inodes\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "PercentCapacity", resource.PercentCapacity))

	err = d.Set("percent_capacity", resource.PercentCapacity)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"percent_capacity\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "DefaultUserQuota", resource.DefaultUserQuota))

	tflog.Debug(ctx, fmt.Sprintf("Found a pointer object %v", resource.DefaultUserQuota))
	err = d.Set("default_user_quota", utils.FlattenModelAsList(ctx, resource.DefaultUserQuota))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"default_user_quota\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "DefaultGroupQuota", resource.DefaultGroupQuota))

	tflog.Debug(ctx, fmt.Sprintf("Found a pointer object %v", resource.DefaultGroupQuota))
	err = d.Set("default_group_quota", utils.FlattenModelAsList(ctx, resource.DefaultGroupQuota))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"default_group_quota\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "UserQuotas", resource.UserQuotas))

	err = d.Set("user_quotas", utils.FlattenListOfModelsToList(ctx, resource.UserQuotas))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"user_quotas\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "GroupQuotas", resource.GroupQuotas))

	err = d.Set("group_quotas", utils.FlattenListOfModelsToList(ctx, resource.GroupQuotas))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"group_quotas\"",
			Detail:   err.Error(),
		})
	}

	Id := (int64)(resource.Id)
	d.SetId(strconv.FormatInt(Id, 10))
	return diags
}
