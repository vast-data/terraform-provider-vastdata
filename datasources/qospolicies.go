package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api_latest "github.com/vast-data/terraform-provider-vastdata/codegen/latest"
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
				Description: ``,
			},

			"guid": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `QoS Policy guid`,
			},

			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    false,
				Required:    true,
				Optional:    false,
				Description: ``,
			},

			"mode": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `QoS provisioning mode Allowed Values are [STATIC USED_CAPACITY PROVISIONED_CAPACITY]`,
			},

			"io_size_bytes": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `Sets the size of IO for static and capacity limit definitions. The number of IOs per request is obtained by dividing request size by IO size. Default: 64K, Recommended range: 4K - 1M`,
			},

			"static_limits": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: ``,

				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"min_reads_bw_mbps": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `Minimal amount of performance to provide when there is resource contention`,
						},

						"max_reads_bw_mbps": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `Maximal amount of performance to provide when there is no resource contention`,
						},

						"min_writes_bw_mbps": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `Minimal amount of performance to provide when there is resource contention`,
						},

						"max_writes_bw_mbps": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `Maximal amount of performance to provide when there is no resource contention`,
						},

						"min_reads_iops": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `Minimal amount of performance to provide when there is resource contention`,
						},

						"max_reads_iops": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `Maximal amount of performance to provide when there is no resource contention`,
						},

						"min_writes_iops": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `Minimal amount of performance to provide when there is resource contention`,
						},

						"max_writes_iops": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `Maximal amount of performance to provide when there is no resource contention`,
						},
					},
				},
			},

			"capacity_limits": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: ``,

				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"max_reads_bw_mbps_per_gb_capacity": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `Maximal amount of performance per GB to provide when there is no resource contention`,
						},

						"max_writes_bw_mbps_per_gb_capacity": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `Maximal amount of performance per GB to provide when there is no resource contention`,
						},

						"max_reads_iops_per_gb_capacity": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `Maximal amount of performance per GB to provide when there is no resource contention`,
						},

						"max_writes_iops_per_gb_capacity": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `Maximal amount of performance per GB to provide when there is no resource contention`,
						},
					},
				},
			},
		},
	}
}

func dataSourceQosPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(vast_client.JwtSession)
	values := url.Values{}

	name := d.Get("name")
	values.Add("name", fmt.Sprintf("%v", name))

	response, err := client.Get(ctx, "/api/qospolicies/", values.Encode(), map[string]string{})
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

	body, err := utils.DefaultProcessingFunc(ctx, response)
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

	Id := (int64)(resource.Id)
	d.SetId(strconv.FormatInt(Id, 10))
	return diags
}
