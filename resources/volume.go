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

func ResourceVolume() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceVolumeRead,
		DeleteContext: resourceVolumeDelete,
		CreateContext: resourceVolumeCreate,
		UpdateContext: resourceVolumeUpdate,

		Importer: &schema.ResourceImporter{
			StateContext: resourceVolumeImporter,
		},

		Description: ``,
		Schema:      getResourceVolumeSchema(),
	}
}

func getResourceVolumeSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{

		"guid": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Volume").GetConflictingFields("guid"),

			Computed:    true,
			Optional:    false,
			Sensitive:   false,
			Description: `(Valid for versions: 5.3.0) A uniqe GUID assigned to the manager`,
		},

		"name": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Volume").GetConflictingFields("name"),

			Required:    true,
			Description: `(Valid for versions: 5.3.0) A uniqe name given to the volume`,
		},

		"size": &schema.Schema{
			Type:          schema.TypeInt,
			ConflictsWith: codegen_configs.GetResourceByName("Volume").GetConflictingFields("size"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.3.0) The volume size of the volume in bytes`,
		},

		"view_id": &schema.Schema{
			Type:          schema.TypeInt,
			ConflictsWith: codegen_configs.GetResourceByName("Volume").GetConflictingFields("view_id"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.3.0) The View ID to relate this volume with , must be a View with protocol defined as BLOCK`,
		},
	}
}

var Volume_names_mapping map[string][]string = map[string][]string{}

func ResourceVolumeReadStructIntoSchema(ctx context.Context, resource api_latest.Volume, d *schema.ResourceData) diag.Diagnostics {
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

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Size", resource.Size))

	err = d.Set("size", resource.Size)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"size\"",
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

	return diags

}
func resourceVolumeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(vast_client.JwtSession)
	resource_config := codegen_configs.GetResourceByName("Volume")
	attrs := map[string]interface{}{"path": utils.GenPath("volumes"), "id": d.Id()}
	tflog.Debug(ctx, fmt.Sprintf("[resourceVolumeRead] Calling Get Function : %v for resource Volume", utils.GetFuncName(resource_config.GetFunc)))
	response, err := resource_config.GetFunc(ctx, client, attrs, d, map[string]string{})
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
	resource := api_latest.Volume{}
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
	diags = ResourceVolumeReadStructIntoSchema(ctx, resource, d)

	return diags
}

func resourceVolumeDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(vast_client.JwtSession)
	resource_config := codegen_configs.GetResourceByName("Volume")
	attrs := map[string]interface{}{"path": utils.GenPath("volumes"), "id": d.Id()}

	response, err := resource_config.DeleteFunc(ctx, client, attrs, nil, map[string]string{})

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

func resourceVolumeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	names_mapping := utils.ContextKey("names_mapping")
	new_ctx := context.WithValue(ctx, names_mapping, Volume_names_mapping)
	var diags diag.Diagnostics
	data := make(map[string]interface{})
	client := m.(vast_client.JwtSession)
	resource_config := codegen_configs.GetResourceByName("Volume")
	tflog.Info(ctx, fmt.Sprintf("Creating Resource Volume"))
	reflect_Volume := reflect.TypeOf((*api_latest.Volume)(nil))
	utils.PopulateResourceMap(new_ctx, reflect_Volume.Elem(), d, &data, "", false)

	version_compare := utils.VastVersionsWarn(ctx)

	if version_compare != metadata.CLUSTER_VERSION_EQUALS {
		cluster_version := metadata.ClusterVersionString()
		t, t_exists := vast_versions.GetVersionedType(cluster_version, "Volume")
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
			tflog.Warn(ctx, fmt.Sprintf("Could have not found resource %s in version %s , things might not work properly", "Volume", cluster_version))
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
	attrs := map[string]interface{}{"path": utils.GenPath("volumes")}
	response, create_err := resource_config.CreateFunc(ctx, client, attrs, data, map[string]string{})
	tflog.Info(ctx, fmt.Sprintf("Server Error for  Volume %v", create_err))

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
	resource := api_latest.Volume{}
	err = json.Unmarshal(response_body, &resource)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to convert response body into Volume",
			Detail:   err.Error(),
		})
		return diags
	}

	id_err := resource_config.IdFunc(ctx, client, resource.Id, d)
	if id_err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to set Id",
			Detail:   err.Error(),
		})
		return diags
	}
	ctx_with_resource := context.WithValue(ctx, utils.ContextKey("resource"), resource)
	resourceVolumeRead(ctx_with_resource, d, m)

	return diags
}

func resourceVolumeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	names_mapping := utils.ContextKey("names_mapping")
	new_ctx := context.WithValue(ctx, names_mapping, Volume_names_mapping)
	var diags diag.Diagnostics
	data := make(map[string]interface{})
	version_compare := utils.VastVersionsWarn(ctx)
	resource_config := codegen_configs.GetResourceByName("Volume")
	if version_compare != metadata.CLUSTER_VERSION_EQUALS {
		cluster_version := metadata.ClusterVersionString()
		t, t_exists := vast_versions.GetVersionedType(cluster_version, "Volume")
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
			tflog.Warn(ctx, fmt.Sprintf("Could have not found resource %s in version %s , things might not work properly", "Volume", cluster_version))
		}
	}

	client := m.(vast_client.JwtSession)
	tflog.Info(ctx, fmt.Sprintf("Updating Resource Volume"))
	reflect_Volume := reflect.TypeOf((*api_latest.Volume)(nil))
	utils.PopulateResourceMap(new_ctx, reflect_Volume.Elem(), d, &data, "", false)

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
	attrs := map[string]interface{}{"path": utils.GenPath("volumes"), "id": d.Id()}
	response, patch_err := resource_config.UpdateFunc(ctx, client, attrs, data, d, map[string]string{})
	tflog.Info(ctx, fmt.Sprintf("Server Error for  Volume %v", patch_err))
	if patch_err != nil {
		error_message := patch_err.Error() + " Server Response: " + utils.GetResponseBodyAsStr(response)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Object Creation Failed",
			Detail:   error_message,
		})
		return diags
	}
	resourceVolumeRead(ctx, d, m)

	return diags

}

func resourceVolumeImporter(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {

	result := []*schema.ResourceData{}
	client := m.(vast_client.JwtSession)
	resource_config := codegen_configs.GetResourceByName("Volume")
	attrs := map[string]interface{}{"path": utils.GenPath("volumes")}
	response, err := resource_config.ImportFunc(ctx, client, attrs, d, resource_config.Importer.GetFunc())

	if err != nil {
		return result, err
	}

	resource_l := []api_latest.Volume{}
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

	diags := ResourceVolumeReadStructIntoSchema(ctx, resource, d)
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
