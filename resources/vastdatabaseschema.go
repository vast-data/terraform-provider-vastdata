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

func ResourceVastDatabaseSchema() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceVastDatabaseSchemaRead,
		DeleteContext: resourceVastDatabaseSchemaDelete,
		CreateContext: resourceVastDatabaseSchemaCreate,
		UpdateContext: resourceVastDatabaseSchemaUpdate,

		Description: ``,
		Schema:      getResourceVastDatabaseSchemaSchema(),
	}
}

func getResourceVastDatabaseSchemaSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{

		"guid": &schema.Schema{
			Type: schema.TypeString,

			Computed:    true,
			Optional:    false,
			Sensitive:   false,
			Description: `A unique guid given to the databse`,
		},

		"database_name": &schema.Schema{
			Type: schema.TypeString,

			Required: true,
		},

		"name": &schema.Schema{
			Type: schema.TypeString,

			Required: true,
		},

		"identifier": &schema.Schema{
			Type: schema.TypeString,

			Computed:    true,
			Optional:    false,
			Sensitive:   false,
			Description: `A computed identifier of the table`,
		},
	}
}

var VastDatabaseSchema_names_mapping map[string][]string = map[string][]string{}

func ResourceVastDatabaseSchemaReadStructIntoSchema(ctx context.Context, resource api_latest.VastDatabaseSchema, d *schema.ResourceData) diag.Diagnostics {
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

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "DatabaseName", resource.DatabaseName))

	err = d.Set("database_name", resource.DatabaseName)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"database_name\"",
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

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Identifier", resource.Identifier))

	err = d.Set("identifier", resource.Identifier)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"identifier\"",
			Detail:   err.Error(),
		})
	}

	return diags

}
func resourceVastDatabaseSchemaRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	response, err := utils.ReadVastDataDatabseSchema(ctx, m, "/api/latest/schemas/", "", map[string]string{}, d)

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
	resource := api_latest.VastDatabaseSchema{}

	body, read_before_unmarshall_err := utils.ReadResultField(ctx, response)
	tflog.Debug(ctx, fmt.Sprintf("Body VastDatabaseSchema returned after processing response %v", string(body)))
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
	diags = ResourceVastDatabaseSchemaReadStructIntoSchema(ctx, resource, d)

	var after_read_error error
	after_read_error = utils.AddVastdDatabaseSchemaIdentifierFieled(resource, ctx, d)
	if after_read_error != nil {
		return diag.FromErr(after_read_error)
	}

	return diags
}

func resourceVastDatabaseSchemaDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	response, err := utils.DeleteVastDatabaseSchema(ctx, m, "/api/latest/schemas/", "", map[string]string{}, d)

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

func resourceVastDatabaseSchemaCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	names_mapping := utils.ContextKey("names_mapping")
	new_ctx := context.WithValue(ctx, names_mapping, VastDatabaseSchema_names_mapping)
	var diags diag.Diagnostics
	data := make(map[string]interface{})
	client := m.(vast_client.JwtSession)
	tflog.Info(ctx, fmt.Sprintf("Creating Resource VastDatabaseSchema"))
	reflect_VastDatabaseSchema := reflect.TypeOf((*api_latest.VastDatabaseSchema)(nil))
	utils.PopulateResourceMap(new_ctx, reflect_VastDatabaseSchema.Elem(), d, &data, "", false)

	version_compare := utils.VastVersionsWarn(ctx)

	if version_compare != metadata.CLUSTER_VERSION_EQUALS {
		cluster_version := metadata.ClusterVersionString()
		t, t_exists := vast_versions.GetVersionedType(cluster_version, "VastDatabaseSchema")
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
			tflog.Warn(ctx, fmt.Sprintf("Could have not found resource %s in version %s , things might not work properly", "VastDatabaseSchema", cluster_version))
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
	response, create_err := client.Post(ctx, "/api/latest/schemas/", bytes.NewReader(b), map[string]string{})
	tflog.Info(ctx, fmt.Sprintf("Server Error for  VastDatabaseSchema %v", create_err))

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
	tflog.Debug(ctx, fmt.Sprintf("Object type VastDatabaseSchema created , server response %v", string(response_body)))
	resource := api_latest.VastDatabaseSchema{}

	response_body, err = utils.VastDatabaseSchemaBeforeCreateUnmarshel(ctx, response_body, d)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed before unmarshl func for VastDatabaseSchema",
			Detail:   err.Error(),
		})
		return diags
	}

	err = json.Unmarshal(response_body, &resource)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to convert response body into VastDatabaseSchema",
			Detail:   err.Error(),
		})
		return diags
	}

	d.SetId(resource.Id)

	resourceVastDatabaseSchemaRead(ctx, d, m)

	return diags
}

func resourceVastDatabaseSchemaUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	names_mapping := utils.ContextKey("names_mapping")
	new_ctx := context.WithValue(ctx, names_mapping, VastDatabaseSchema_names_mapping)
	var diags diag.Diagnostics
	data := make(map[string]interface{})
	version_compare := utils.VastVersionsWarn(ctx)
	if version_compare != metadata.CLUSTER_VERSION_EQUALS {
		cluster_version := metadata.ClusterVersionString()
		t, t_exists := vast_versions.GetVersionedType(cluster_version, "VastDatabaseSchema")
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
			tflog.Warn(ctx, fmt.Sprintf("Could have not found resource %s in version %s , things might not work properly", "VastDatabaseSchema", cluster_version))
		}
	}

	tflog.Info(ctx, fmt.Sprintf("Updating Resource VastDatabaseSchema"))
	reflect_VastDatabaseSchema := reflect.TypeOf((*api_latest.VastDatabaseSchema)(nil))
	utils.PopulateResourceMap(new_ctx, reflect_VastDatabaseSchema.Elem(), d, &data, "", false)

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

	response, patch_err := utils.UpdateVastDatabaseSchemaName(ctx, m, "/api/latest/schemas/", "application/json", map[string]string{}, d)

	tflog.Info(ctx, fmt.Sprintf("Server Error for  VastDatabaseSchema %v", patch_err))
	if patch_err != nil {
		error_message := patch_err.Error() + " Server Response: " + utils.GetResponseBodyAsStr(response)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Object Creation Failed",
			Detail:   error_message,
		})
		return diags
	}
	resourceVastDatabaseSchemaRead(ctx, d, m)

	return diags

}
