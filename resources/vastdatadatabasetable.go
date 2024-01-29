package resources

import (
	"io"

	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api_latest "github.com/vast-data/terraform-provider-vastdata/codegen/latest"
	metadata "github.com/vast-data/terraform-provider-vastdata/metadata"
	utils "github.com/vast-data/terraform-provider-vastdata/utils"
	vast_client "github.com/vast-data/terraform-provider-vastdata/vast-client"
	vast_versions "github.com/vast-data/terraform-provider-vastdata/vast_versions"
)

func ResourceVastDatabaseTable() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceVastDatabaseTableRead,
		DeleteContext: resourceVastDatabaseTableDelete,
		CreateContext: resourceVastDatabaseTableCreate,
		UpdateContext: resourceVastDatabaseTableUpdate,

		Description: ``,
		Schema:      getResourceVastDatabaseTableSchema(),
	}
}

func getResourceVastDatabaseTableSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{

		"guid": &schema.Schema{
			Type: schema.TypeString,

			Computed:    true,
			Optional:    false,
			Sensitive:   false,
			Description: `A unique guid given to the databse`,
		},

		"name": &schema.Schema{
			Type: schema.TypeString,

			Required: true,
		},

		"schema_identifier": &schema.Schema{
			Type: schema.TypeString,

			Required: true,
		},

		"fields": &schema.Schema{
			Type: schema.TypeString,

			DiffSuppressOnRefresh: true,
			DiffSuppressFunc:      utils.JsonStructureCompare,
			Required:              true,
		},
	}
}

var VastDatabaseTable_names_mapping map[string][]string = map[string][]string{}

func ResourceVastDatabaseTableReadStructIntoSchema(ctx context.Context, resource api_latest.VastDatabaseTable, d *schema.ResourceData) diag.Diagnostics {
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

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "SchemaIdentifier", resource.SchemaIdentifier))

	err = d.Set("schema_identifier", resource.SchemaIdentifier)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"schema_identifier\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Fields", resource.Fields))

	err = d.Set("fields", resource.Fields)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"fields\"",
			Detail:   err.Error(),
		})
	}

	return diags

}
func resourceVastDatabaseTableRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	response, err := utils.ReadVastDataDatabseTable(ctx, m, "/api/latest/tables/", "", map[string]string{}, d)

	utils.VastVersionsWarn(ctx)

	tflog.Info(ctx, response.Request.URL.String())
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured while obtaining data from the vastdata cluster",
			Detail:   err.Error(),
		})
		return diags

	}
	resource := api_latest.VastDatabaseTable{}

	body, read_before_unmarshall_err := utils.GenerateTableReadResponse(ctx, response)
	tflog.Debug(ctx, fmt.Sprintf("Body VastDatabaseTable returned after processing response %v", string(body)))
	if read_before_unmarshall_err != nil {
		return diag.FromErr(read_before_unmarshall_err)
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
	diags = ResourceVastDatabaseTableReadStructIntoSchema(ctx, resource, d)

	return diags
}

func resourceVastDatabaseTableDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	response, err := utils.DeleteVastDatabaseTable(ctx, m, "/api/latest/tables/", "", map[string]string{}, d)

	tflog.Info(ctx, fmt.Sprintf("Removing Resource"))
	tflog.Info(ctx, response.Request.URL.String())
	tflog.Info(ctx, utils.GetResponseBodyAsStr(response))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured while deleting a resource from the vastdata cluster",
			Detail:   err.Error(),
		})

	}

	return diags

}

func resourceVastDatabaseTableCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	names_mapping := utils.ContextKey("names_mapping")
	new_ctx := context.WithValue(ctx, names_mapping, VastDatabaseTable_names_mapping)
	var diags diag.Diagnostics
	data := make(map[string]interface{})
	client := m.(vast_client.JwtSession)
	tflog.Info(ctx, fmt.Sprintf("Creating Resource VastDatabaseTable"))
	reflect_VastDatabaseTable := reflect.TypeOf((*api_latest.VastDatabaseTable)(nil))
	utils.PopulateResourceMap(new_ctx, reflect_VastDatabaseTable.Elem(), d, &data, "", false)

	var before_post_error error
	data, before_post_error = utils.ConstructTableFields(data, client, ctx, d)
	if before_post_error != nil {
		return diag.FromErr(before_post_error)
	}

	version_compare := utils.VastVersionsWarn(ctx)

	if version_compare != metadata.CLUSTER_VERSION_EQUALS {
		cluster_version := metadata.ClusterVersionString()
		t, t_exists := vast_versions.GetVersionedType(cluster_version, "VastDatabaseTable")
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
			tflog.Warn(ctx, fmt.Sprintf("Could have not found resource %s in version %s , things might not work properly", "VastDatabaseTable", cluster_version))
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
	response, create_err := client.Post(ctx, "/api/latest/tables/", bytes.NewReader(b), map[string]string{})
	tflog.Info(ctx, fmt.Sprintf("Server Error for  VastDatabaseTable %v", create_err))

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
	tflog.Debug(ctx, fmt.Sprintf("Object type VastDatabaseTable created , server response %v", string(response_body)))
	resource := api_latest.VastDatabaseTable{}

	response_body, err = utils.VastDatabaseTableBeforeCreateUnmarshel(ctx, response_body, d)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed before unmarshl func for VastDatabaseTable",
			Detail:   err.Error(),
		})
		return diags
	}

	err = json.Unmarshal(response_body, &resource)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to convert response body into VastDatabaseTable",
			Detail:   err.Error(),
		})
		return diags
	}

	d.SetId(resource.Id)

	resourceVastDatabaseTableRead(ctx, d, m)

	return diags
}

func resourceVastDatabaseTableUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	names_mapping := utils.ContextKey("names_mapping")
	new_ctx := context.WithValue(ctx, names_mapping, VastDatabaseTable_names_mapping)
	var diags diag.Diagnostics
	data := make(map[string]interface{})
	version_compare := utils.VastVersionsWarn(ctx)
	if version_compare != metadata.CLUSTER_VERSION_EQUALS {
		cluster_version := metadata.ClusterVersionString()
		t, t_exists := vast_versions.GetVersionedType(cluster_version, "VastDatabaseTable")
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
			tflog.Warn(ctx, fmt.Sprintf("Could have not found resource %s in version %s , things might not work properly", "VastDatabaseTable", cluster_version))
		}
	}

	client := m.(vast_client.JwtSession)
	VastDatabaseTableId := d.Id()

	tflog.Info(ctx, fmt.Sprintf("Updating Resource VastDatabaseTable"))
	reflect_VastDatabaseTable := reflect.TypeOf((*api_latest.VastDatabaseTable)(nil))
	utils.PopulateResourceMap(new_ctx, reflect_VastDatabaseTable.Elem(), d, &data, "", false)

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

	response, patch_err := client.Patch(ctx, fmt.Sprintf("/api/latest/tables//%v", VastDatabaseTableId), "application/json", bytes.NewReader(b), map[string]string{})

	tflog.Info(ctx, fmt.Sprintf("Server Error for  VastDatabaseTable %v", patch_err))
	if patch_err != nil {
		error_message := patch_err.Error() + " Server Response: " + utils.GetResponseBodyAsStr(response)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Object Creation Failed",
			Detail:   error_message,
		})
		return diags
	}
	resourceVastDatabaseTableRead(ctx, d, m)

	return diags

}
