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

func DataSourceUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUserRead,
		Schema: map[string]*schema.Schema{

			"id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
				Required: false,
				Optional: false,
			},

			"guid": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Required: false,
				Optional: false,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: false,
				Required: true,
				Optional: false,
			},

			"uid": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
				Required: false,
				Optional: false,
			},

			"leading_gid": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
				Required: false,
				Optional: false,
			},

			"gids": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Required: false,
				Optional: false,

				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},

			"groups": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Required: false,
				Optional: false,

				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"group_count": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
				Required: false,
				Optional: false,
			},

			"leading_group_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Required: false,
				Optional: false,
			},

			"leading_group_gid": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
				Required: false,
				Optional: false,
			},

			"sid": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Required: false,
				Optional: false,
			},

			"primary_group_sid": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Required: false,
				Optional: false,
			},

			"sids": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Required: false,
				Optional: false,

				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"local": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
				Required: false,
				Optional: false,
			},

			"access_keys": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Required: false,
				Optional: false,

				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"access_key": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "",
						},

						"enabled": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "",
						},
					},
				},
			},

			"allow_create_bucket": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
				Required: false,
				Optional: false,
			},

			"allow_delete_bucket": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
				Required: false,
				Optional: false,
			},

			"s3_superuser": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
				Required: false,
				Optional: false,
			},

			"s3_policies_ids": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Required: false,
				Optional: false,

				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
		},
	}
}

func dataSourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(vast_client.JwtSession)
	values := url.Values{}

	name := d.Get("name")
	values.Add("name", fmt.Sprintf("%v", name))

	response, err := client.Get(ctx, "/api/users/", values.Encode(), map[string]string{})
	tflog.Info(ctx, response.Request.URL.String())
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured while obtaining data from the vastdata cluster",
			Detail:   err.Error(),
		})
		return diags

	}
	resource_l := []api_latest.User{}

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

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Uid", resource.Uid))

	err = d.Set("uid", resource.Uid)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"uid\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "LeadingGid", resource.LeadingGid))

	err = d.Set("leading_gid", resource.LeadingGid)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"leading_gid\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Gids", resource.Gids))

	err = d.Set("gids", utils.FlattenListOfPrimitives(&resource.Gids))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"gids\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Groups", resource.Groups))

	err = d.Set("groups", utils.FlattenListOfPrimitives(&resource.Groups))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"groups\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "GroupCount", resource.GroupCount))

	err = d.Set("group_count", resource.GroupCount)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"group_count\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "LeadingGroupName", resource.LeadingGroupName))

	err = d.Set("leading_group_name", resource.LeadingGroupName)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"leading_group_name\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "LeadingGroupGid", resource.LeadingGroupGid))

	err = d.Set("leading_group_gid", resource.LeadingGroupGid)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"leading_group_gid\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Sid", resource.Sid))

	err = d.Set("sid", resource.Sid)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"sid\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "PrimaryGroupSid", resource.PrimaryGroupSid))

	err = d.Set("primary_group_sid", resource.PrimaryGroupSid)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"primary_group_sid\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Sids", resource.Sids))

	err = d.Set("sids", utils.FlattenListOfPrimitives(&resource.Sids))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"sids\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Local", resource.Local))

	err = d.Set("local", resource.Local)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"local\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "AccessKeys", resource.AccessKeys))

	err = d.Set("access_keys", utils.FlattenListOfStringsList(&resource.AccessKeys, []string{"access_key", "enabled"}))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"access_keys\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "AllowCreateBucket", resource.AllowCreateBucket))

	err = d.Set("allow_create_bucket", resource.AllowCreateBucket)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"allow_create_bucket\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "AllowDeleteBucket", resource.AllowDeleteBucket))

	err = d.Set("allow_delete_bucket", resource.AllowDeleteBucket)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"allow_delete_bucket\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "S3Superuser", resource.S3Superuser))

	err = d.Set("s3_superuser", resource.S3Superuser)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"s3_superuser\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "S3PoliciesIds", resource.S3PoliciesIds))

	err = d.Set("s3_policies_ids", utils.FlattenListOfPrimitives(&resource.S3PoliciesIds))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"s3_policies_ids\"",
			Detail:   err.Error(),
		})
	}

	Id := (int64)(resource.Id)
	d.SetId(strconv.FormatInt(Id, 10))
	return diags
}
