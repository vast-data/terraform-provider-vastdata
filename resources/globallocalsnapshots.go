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

func ResourceGlobalLocalSnapshot() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceGlobalLocalSnapshotRead,
		DeleteContext: resourceGlobalLocalSnapshotDelete,
		CreateContext: resourceGlobalLocalSnapshotCreate,
		UpdateContext: resourceGlobalLocalSnapshotUpdate,

		Importer: &schema.ResourceImporter{
			StateContext: resourceGlobalLocalSnapshotImporter,
		},

		Description: ``,
		Schema:      getResourceGlobalLocalSnapshotSchema(),
	}
}

func getResourceGlobalLocalSnapshotSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{

		"guid": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("GlobalLocalSnapshot").GetConflictingFields("guid"),

			Computed:    true,
			Optional:    false,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) A unique guid given to the global snapshot`,
		},

		"name": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("GlobalLocalSnapshot").GetConflictingFields("name"),

			Required:    true,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The name of the snapshot`,
		},

		"loanee_tenant_id": {
			Type:          schema.TypeInt,
			ConflictsWith: codegen_configs.GetResourceByName("GlobalLocalSnapshot").GetConflictingFields("loanee_tenant_id"),

			Required:    true,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The tenant ID of the target`,
		},

		"loanee_root_path": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("GlobalLocalSnapshot").GetConflictingFields("loanee_root_path"),

			Required:    true,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The path where to store the snapshot on a Target`,
		},

		"loanee_snapshot_id": {
			Type:          schema.TypeInt,
			ConflictsWith: codegen_configs.GetResourceByName("GlobalLocalSnapshot").GetConflictingFields("loanee_snapshot_id"),

			Required:    true,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The id of the local snapshot`,
		},

		"enabled": {
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("GlobalLocalSnapshot").GetConflictingFields("enabled"),

			Computed:    false,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Is the snapshot enabled`,

			Default: true,
		},

		"owner_tenant": {
			Type:          schema.TypeList,
			ConflictsWith: codegen_configs.GetResourceByName("GlobalLocalSnapshot").GetConflictingFields("owner_tenant"),

			Required:    true,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,

			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{

					"name": {
						Type:          schema.TypeString,
						ConflictsWith: codegen_configs.GetResourceByName("GlobalSnapshotOwnerTenant").GetConflictingFields("name"),

						Required:    true,
						Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Ten name of the remote Tenant`,
					},

					"guid": {
						Type:          schema.TypeString,
						ConflictsWith: codegen_configs.GetResourceByName("GlobalSnapshotOwnerTenant").GetConflictingFields("guid"),

						Required:    true,
						Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The remote tenant guid`,
					},
				},
			},
		},
	}
}

var GlobalLocalSnapshotNamesMapping = map[string][]string{}

func ResourceGlobalLocalSnapshotReadStructIntoSchema(ctx context.Context, resource api_latest.GlobalLocalSnapshot, d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics
	var err error

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Guid", resource.Guid))

	err = d.Set("guid", resource.Guid)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"guid\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Name", resource.Name))

	err = d.Set("name", resource.Name)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"name\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "LoaneeTenantId", resource.LoaneeTenantId))

	err = d.Set("loanee_tenant_id", resource.LoaneeTenantId)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"loanee_tenant_id\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "LoaneeRootPath", resource.LoaneeRootPath))

	err = d.Set("loanee_root_path", resource.LoaneeRootPath)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"loanee_root_path\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "LoaneeSnapshotId", resource.LoaneeSnapshotId))

	err = d.Set("loanee_snapshot_id", resource.LoaneeSnapshotId)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"loanee_snapshot_id\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Enabled", resource.Enabled))

	err = d.Set("enabled", resource.Enabled)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"enabled\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "OwnerTenant", resource.OwnerTenant))

	tflog.Debug(ctx, fmt.Sprintf("Found a pointer object %v", resource.OwnerTenant))
	err = d.Set("owner_tenant", utils.FlattenModelAsList(ctx, resource.OwnerTenant))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"owner_tenant\"",
			Detail:   err.Error(),
		})
	}

	return diags

}
func resourceGlobalLocalSnapshotRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*vast_client.VMSSession)
	resourceConfig := codegen_configs.GetResourceByName("GlobalLocalSnapshot")
	attrs := map[string]interface{}{"path": utils.GenPath("globalsnapstreams"), "id": d.Id()}
	tflog.Debug(ctx, fmt.Sprintf("[resourceGlobalLocalSnapshotRead] Calling Get Function : %v for resource GlobalLocalSnapshot", utils.GetFuncName(resourceConfig.GetFunc)))
	response, err := resourceConfig.GetFunc(ctx, client, attrs, d, map[string]string{})
	utils.VastVersionsWarn(ctx)

	var body []byte
	var resource api_latest.GlobalLocalSnapshot
	if err != nil && response != nil && response.StatusCode == 404 {
		var fallbackErr error
		body, fallbackErr = utils.HandleFallback(ctx, client, attrs, d, resourceConfig.IdFunc)
		if fallbackErr != nil {
			errorMessage := fmt.Sprintf("Initial request failed:\n%v\nFallback request also failed:\n%v", err.Error(), fallbackErr.Error())
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error occurred while obtaining data from the VAST Data cluster",
				Detail:   errorMessage,
			})
			return diags
		}
	} else if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred while obtaining data from the VAST Data cluster",
			Detail:   err.Error(),
		})
		return diags
	} else {
		tflog.Info(ctx, response.Request.URL.String())
		body, err = resourceConfig.ResponseProcessingFunc(ctx, response)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error occurred reading data received from VAST Data cluster",
				Detail:   err.Error(),
			})
			return diags
		}
	}
	err = json.Unmarshal(body, &resource)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred while parsing data received from VAST Data cluster",
			Detail:   err.Error(),
		})
		return diags

	}
	diags = ResourceGlobalLocalSnapshotReadStructIntoSchema(ctx, resource, d)

	return diags
}

func resourceGlobalLocalSnapshotDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*vast_client.VMSSession)
	resourceConfig := codegen_configs.GetResourceByName("GlobalLocalSnapshot")
	attrs := map[string]interface{}{"path": utils.GenPath("globalsnapstreams"), "id": d.Id()}

	response, err := resourceConfig.DeleteFunc(ctx, client, attrs, nil, map[string]string{})

	tflog.Info(ctx, fmt.Sprintf("Removing Resource"))
	if response != nil {
		tflog.Info(ctx, response.Request.URL.String())
		tflog.Info(ctx, utils.GetResponseBodyAsStr(response))
	}

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred while deleting a resource from the VAST Data cluster",
			Detail:   err.Error(),
		})

	}

	return diags

}

func resourceGlobalLocalSnapshotCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	namesMapping := utils.ContextKey("namesMapping")
	newCtx := context.WithValue(ctx, namesMapping, GlobalLocalSnapshotNamesMapping)
	var diags diag.Diagnostics
	data := make(map[string]interface{})
	client := m.(*vast_client.VMSSession)
	resourceConfig := codegen_configs.GetResourceByName("GlobalLocalSnapshot")
	tflog.Info(ctx, fmt.Sprintf("Creating Resource GlobalLocalSnapshot"))
	reflectGlobalLocalSnapshot := reflect.TypeOf((*api_latest.GlobalLocalSnapshot)(nil))
	utils.PopulateResourceMap(newCtx, reflectGlobalLocalSnapshot.Elem(), d, &data, "", false)

	versionsEqual := utils.VastVersionsWarn(ctx)

	if versionsEqual != metadata.CLUSTER_VERSION_EQUALS {
		clusterVersion := metadata.ClusterVersionString()
		t, typeExists := vast_versions.GetVersionedType(clusterVersion, "GlobalLocalSnapshot")
		if typeExists {
			versionError := utils.VersionMatch(t, data)
			if versionError != nil {
				tflog.Warn(ctx, versionError.Error())
				versionValidationMode, versionValidationModeExists := metadata.GetClusterConfig("version_validation_mode")
				tflog.Warn(ctx, fmt.Sprintf("Version Validation Mode Detected %s", versionValidationMode))
				if versionValidationModeExists && versionValidationMode == "strict" {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "Cluster Version & Build Version Are Too Different",
						Detail:   versionError.Error(),
					})
					return diags
				}
			}
		} else {
			tflog.Warn(ctx, fmt.Sprintf("Could have not found resource %s in version %s, things might not work properly", "GlobalLocalSnapshot", clusterVersion))
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
	attrs := map[string]interface{}{"path": utils.GenPath("globalsnapstreams")}
	response, createErr := resourceConfig.CreateFunc(ctx, client, attrs, data, map[string]string{})
	tflog.Info(ctx, fmt.Sprintf("Server Error for  GlobalLocalSnapshot %v", createErr))

	if createErr != nil {
		errorMessage := fmt.Sprintf("server response:\n%v\nUnderlying error:\n%v", utils.GetResponseBodyAsStr(response), createErr.Error())
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Object Creation Failed",
			Detail:   errorMessage,
		})
		return diags
	}
	responseBody, _ := io.ReadAll(response.Body)
	tflog.Debug(ctx, fmt.Sprintf("Object created, server response %v", string(responseBody)))
	resource := api_latest.GlobalLocalSnapshot{}
	err = json.Unmarshal(responseBody, &resource)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to convert response body into GlobalLocalSnapshot",
			Detail:   err.Error(),
		})
		return diags
	}

	err = resourceConfig.IdFunc(ctx, client, resource.Id, d)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to set Id",
			Detail:   err.Error(),
		})
		return diags
	}
	ctxWithResource := context.WithValue(ctx, utils.ContextKey("resource"), resource)
	resourceGlobalLocalSnapshotRead(ctxWithResource, d, m)

	return diags
}

func resourceGlobalLocalSnapshotUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	namesMapping := utils.ContextKey("namesMapping")
	newCtx := context.WithValue(ctx, namesMapping, GlobalLocalSnapshotNamesMapping)
	var diags diag.Diagnostics
	data := make(map[string]interface{})
	versionsEqual := utils.VastVersionsWarn(ctx)
	resourceConfig := codegen_configs.GetResourceByName("GlobalLocalSnapshot")
	if versionsEqual != metadata.CLUSTER_VERSION_EQUALS {
		clusterVersion := metadata.ClusterVersionString()
		t, typeExists := vast_versions.GetVersionedType(clusterVersion, "GlobalLocalSnapshot")
		if typeExists {
			versionError := utils.VersionMatch(t, data)
			if versionError != nil {
				tflog.Warn(ctx, versionError.Error())
				versionValidationMode, versionValidationModeExists := metadata.GetClusterConfig("version_validation_mode")
				tflog.Warn(ctx, fmt.Sprintf("Version Validation Mode Detected %s", versionValidationMode))
				if versionValidationModeExists && versionValidationMode == "strict" {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "Cluster Version & Build Version Are Too Different",
						Detail:   versionError.Error(),
					})
					return diags
				}
			}
		} else {
			tflog.Warn(ctx, fmt.Sprintf("Could have not found resource %s in version %s, things might not work properly", "GlobalLocalSnapshot", clusterVersion))
		}
	}

	client := m.(*vast_client.VMSSession)
	tflog.Info(ctx, fmt.Sprintf("Updating Resource GlobalLocalSnapshot"))
	reflectGlobalLocalSnapshot := reflect.TypeOf((*api_latest.GlobalLocalSnapshot)(nil))
	utils.PopulateResourceMap(newCtx, reflectGlobalLocalSnapshot.Elem(), d, &data, "", false)

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
	attrs := map[string]interface{}{"path": utils.GenPath("globalsnapstreams"), "id": d.Id()}
	response, patchErr := resourceConfig.UpdateFunc(ctx, client, attrs, data, d, map[string]string{})
	tflog.Info(ctx, fmt.Sprintf("Server Error for  GlobalLocalSnapshot %v", patchErr))
	if patchErr != nil {
		errorMessage := fmt.Sprintf("server response:\n%v\nUnderlying error:\n%v", utils.GetResponseBodyAsStr(response), patchErr.Error())
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Object Creation Failed",
			Detail:   errorMessage,
		})
		return diags
	}
	resourceGlobalLocalSnapshotRead(ctx, d, m)

	return diags

}

func resourceGlobalLocalSnapshotImporter(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {

	var result []*schema.ResourceData
	client := m.(*vast_client.VMSSession)
	resourceConfig := codegen_configs.GetResourceByName("GlobalLocalSnapshot")
	attrs := map[string]interface{}{"path": utils.GenPath("globalsnapstreams")}
	response, err := resourceConfig.ImportFunc(ctx, client, attrs, d, resourceConfig.Importer.GetFunc())

	if err != nil {
		return result, err
	}

	var resourceList []api_latest.GlobalLocalSnapshot
	body, err := resourceConfig.ResponseProcessingFunc(ctx, response)

	if err != nil {
		return result, err
	}
	err = json.Unmarshal(body, &resourceList)
	if err != nil {
		return result, err
	}

	if len(resourceList) == 0 {
		return result, errors.New("cluster returned 0 elements matching provided guid")
	}

	resource := resourceList[0]
	idErr := resourceConfig.IdFunc(ctx, client, resource.Id, d)
	if idErr != nil {
		return result, idErr
	}

	diags := ResourceGlobalLocalSnapshotReadStructIntoSchema(ctx, resource, d)
	if diags.HasError() {
		allErrors := "Errors occurred while importing:\n"
		for _, dig := range diags {
			allErrors += fmt.Sprintf("Summary:%s\nDetails:%s\n", dig.Summary, dig.Detail)
		}
		return result, errors.New(allErrors)
	}
	result = append(result, d)

	return result, err

}
