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

func DataSourceS3LifeCycleRule() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceS3LifeCycleRuleRead,
		Schema: map[string]*schema.Schema{

			"id": &schema.Schema{
				Type:     schema.TypeInt,
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

			"guid": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Required: false,
				Optional: false,
			},

			"enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
				Required: false,
				Optional: false,
			},

			"prefix": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Required: false,
				Optional: false,
			},

			"min_size": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
				Required: false,
				Optional: false,
			},

			"max_size": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
				Required: false,
				Optional: false,
			},

			"expiration_days": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
				Required: false,
				Optional: false,
			},

			"expiration_date": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Required: false,
				Optional: false,
			},

			"expired_obj_delete_marker": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
				Required: false,
				Optional: false,
			},

			"noncurrent_days": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
				Required: false,
				Optional: false,
			},

			"newer_noncurrent_versions": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
				Required: false,
				Optional: false,
			},

			"abort_mpu_days_after_initiation": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
				Required: false,
				Optional: false,
			},

			"view_path": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Required: false,
				Optional: false,
			},

			"view_id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
				Required: false,
				Optional: false,
			},
		},
	}
}

func dataSourceS3LifeCycleRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(vast_client.JwtSession)
	values := url.Values{}

	name := d.Get("name")
	values.Add("name", fmt.Sprintf("%v", name))

	response, err := client.Get(ctx, "/api/s3lifecyclerules/", values.Encode(), map[string]string{})
	tflog.Info(ctx, response.Request.URL.String())
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured while obtaining data from the vastdata cluster",
			Detail:   err.Error(),
		})
		return diags

	}
	resource_l := []api_latest.S3LifeCycleRule{}

	body, err := utils.ProcessingResultsListResponse(ctx, response)
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

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Name", resource.Name))

	err = d.Set("name", resource.Name)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"name\"",
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

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Enabled", resource.Enabled))

	err = d.Set("enabled", resource.Enabled)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"enabled\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Prefix", resource.Prefix))

	err = d.Set("prefix", resource.Prefix)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"prefix\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "MinSize", resource.MinSize))

	err = d.Set("min_size", resource.MinSize)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"min_size\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "MaxSize", resource.MaxSize))

	err = d.Set("max_size", resource.MaxSize)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"max_size\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "ExpirationDays", resource.ExpirationDays))

	err = d.Set("expiration_days", resource.ExpirationDays)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"expiration_days\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "ExpirationDate", resource.ExpirationDate))

	err = d.Set("expiration_date", resource.ExpirationDate)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"expiration_date\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "ExpiredObjDeleteMarker", resource.ExpiredObjDeleteMarker))

	err = d.Set("expired_obj_delete_marker", resource.ExpiredObjDeleteMarker)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"expired_obj_delete_marker\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "NoncurrentDays", resource.NoncurrentDays))

	err = d.Set("noncurrent_days", resource.NoncurrentDays)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"noncurrent_days\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "NewerNoncurrentVersions", resource.NewerNoncurrentVersions))

	err = d.Set("newer_noncurrent_versions", resource.NewerNoncurrentVersions)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"newer_noncurrent_versions\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "AbortMpuDaysAfterInitiation", resource.AbortMpuDaysAfterInitiation))

	err = d.Set("abort_mpu_days_after_initiation", resource.AbortMpuDaysAfterInitiation)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"abort_mpu_days_after_initiation\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "ViewPath", resource.ViewPath))

	err = d.Set("view_path", resource.ViewPath)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"view_path\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "ViewId", resource.ViewId))

	err = d.Set("view_id", resource.ViewId)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"view_id\"",
			Detail:   err.Error(),
		})
	}

	Id := (int64)(resource.Id)
	d.SetId(strconv.FormatInt(Id, 10))
	return diags
}
