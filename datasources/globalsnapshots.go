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

func DataSourceGlobalSnapshot() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGlobalSnapshotRead,
		Description: ``,
		Schema: map[string]*schema.Schema{

			"id": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `A unique id given to the global snapshot`,
			},

			"guid": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `A unique guid given to the global snapshot`,
			},

			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    false,
				Required:    true,
				Optional:    false,
				Description: `The name of the snapshot`,
			},

			"loanee_tenant_id": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `The tenant ID of the target`,
			},

			"loanee_root_path": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `The path where to store the snapshot on a Target`,
			},

			"remote_target_id": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `The remote replication peering id`,
			},

			"remote_target_guid": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `The remote replication peering guid`,
			},

			"remote_target_path": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `The path on the remote cluster`,
			},

			"enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `Is the snapshot enabled`,
			},

			"owner_root_snapshot": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: ``,

				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{},
				},
			},

			"owner_tenant": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: ``,

				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{},
				},
			},
		},
	}
}

func dataSourceGlobalSnapshotRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(vast_client.JwtSession)
	values := url.Values{}

	name := d.Get("name")
	values.Add("name", fmt.Sprintf("%v", name))

	response, err := client.Get(ctx, utils.GenPath("globalsnapstreams"), values.Encode(), map[string]string{})
	tflog.Info(ctx, response.Request.URL.String())
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured while obtaining data from the vastdata cluster",
			Detail:   err.Error(),
		})
		return diags

	}
	resource_l := []api_latest.GlobalSnapshot{}

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

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "LoaneeTenantId", resource.LoaneeTenantId))

	err = d.Set("loanee_tenant_id", resource.LoaneeTenantId)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"loanee_tenant_id\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "LoaneeRootPath", resource.LoaneeRootPath))

	err = d.Set("loanee_root_path", resource.LoaneeRootPath)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"loanee_root_path\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "RemoteTargetId", resource.RemoteTargetId))

	err = d.Set("remote_target_id", resource.RemoteTargetId)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"remote_target_id\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "RemoteTargetGuid", resource.RemoteTargetGuid))

	err = d.Set("remote_target_guid", resource.RemoteTargetGuid)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"remote_target_guid\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "RemoteTargetPath", resource.RemoteTargetPath))

	err = d.Set("remote_target_path", resource.RemoteTargetPath)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"remote_target_path\"",
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

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "OwnerRootSnapshot", resource.OwnerRootSnapshot))

	tflog.Debug(ctx, fmt.Sprintf("Found a pointer object %v", resource.OwnerRootSnapshot))
	err = d.Set("owner_root_snapshot", utils.FlattenModelAsList(ctx, resource.OwnerRootSnapshot))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"owner_root_snapshot\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "OwnerTenant", resource.OwnerTenant))

	tflog.Debug(ctx, fmt.Sprintf("Found a pointer object %v", resource.OwnerTenant))
	err = d.Set("owner_tenant", utils.FlattenModelAsList(ctx, resource.OwnerTenant))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"owner_tenant\"",
			Detail:   err.Error(),
		})
	}

	Id := (int64)(resource.Id)
	d.SetId(strconv.FormatInt(Id, 10))
	return diags
}
