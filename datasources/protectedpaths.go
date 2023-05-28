package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api_latest "github.com/vast-data/terraform-provider-vastdata.git/codegen/latest"
	utils "github.com/vast-data/terraform-provider-vastdata.git/utils"
	vast_client "github.com/vast-data/terraform-provider-vastdata.git/vast-client"
	"net/url"
	"strconv"
)

func DataSourceProtectedPath() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceProtectedPathRead,
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

			"protection_policy_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Required: false,
				Optional: false,
			},

			"protection_policy_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Required: false,
				Optional: false,
			},

			"source_dir": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Required: false,
				Optional: false,
			},

			"target_exported_dir": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Required: false,
				Optional: false,
			},

			"tenant_id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
				Required: false,
				Optional: false,
			},

			"remote_tenant_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Required: false,
				Optional: false,
			},
		},
	}
}

func dataSourceProtectedPathRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(vast_client.JwtSession)
	values := url.Values{}

	name := d.Get("name")
	values.Add("name", fmt.Sprintf("%v", name))

	response, err := client.Get(ctx, "/api/protectedpaths/", values.Encode(), map[string]string{})
	tflog.Info(ctx, response.Request.URL.String())
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured while obtaining data from the vastdata cluster",
			Detail:   err.Error(),
		})
		return diags

	}
	resource_l := []api_latest.ProtectedPath{}

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

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "ProtectionPolicyId", resource.ProtectionPolicyId))

	err = d.Set("protection_policy_id", resource.ProtectionPolicyId)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"protection_policy_id\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "ProtectionPolicyName", resource.ProtectionPolicyName))

	err = d.Set("protection_policy_name", resource.ProtectionPolicyName)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"protection_policy_name\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "SourceDir", resource.SourceDir))

	err = d.Set("source_dir", resource.SourceDir)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"source_dir\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "TargetExportedDir", resource.TargetExportedDir))

	err = d.Set("target_exported_dir", resource.TargetExportedDir)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"target_exported_dir\"",
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

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "RemoteTenantName", resource.RemoteTenantName))

	err = d.Set("remote_tenant_name", resource.RemoteTenantName)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"remote_tenant_name\"",
			Detail:   err.Error(),
		})
	}

	Id := (int64)(resource.Id)
	d.SetId(strconv.FormatInt(Id, 10))
	return diags
}