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

func ResourceGlobalSnapshot() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceGlobalSnapshotRead,
		DeleteContext: resourceGlobalSnapshotDelete,
		CreateContext: resourceGlobalSnapshotCreate,
		UpdateContext: resourceGlobalSnapshotUpdate,

		Importer: &schema.ResourceImporter{
			StateContext: resourceGlobalSnapshotImporter,
		},

		Description: ``,
		Schema:      getResourceGlobalSnapshotSchema(),
	}
}

func getResourceGlobalSnapshotSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{

		"guid": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("GlobalSnapshot").GetConflictingFields("guid"),

			Computed:    true,
			Optional:    false,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) A unique guid given to the global snapshot`,
		},

		"name": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("GlobalSnapshot").GetConflictingFields("name"),

			Required:    true,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The name of the snapshot`,
		},

		"loanee_tenant_id": {
			Type:          schema.TypeInt,
			ConflictsWith: codegen_configs.GetResourceByName("GlobalSnapshot").GetConflictingFields("loanee_tenant_id"),

			DiffSuppressOnRefresh: false,
			DiffSuppressFunc:      utils.DoNothingOnUpdate(),

			Required:    true,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The tenant ID of the target`,
		},

		"loanee_root_path": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("GlobalSnapshot").GetConflictingFields("loanee_root_path"),

			Required:    true,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The path where to store the snapshot on a Target`,
		},

		"remote_target_id": {
			Type:          schema.TypeInt,
			ConflictsWith: codegen_configs.GetResourceByName("GlobalSnapshot").GetConflictingFields("remote_target_id"),

			Required:    true,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The remote replication peering id`,
		},

		"remote_target_guid": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("GlobalSnapshot").GetConflictingFields("remote_target_guid"),

			Required:    true,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The remote replication peering guid`,
		},

		"remote_target_path": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("GlobalSnapshot").GetConflictingFields("remote_target_path"),

			Required:    true,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The path on the remote cluster`,
		},

		"enabled": {
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("GlobalSnapshot").GetConflictingFields("enabled"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Is the snapshot enabled`,
		},

		"owner_root_snapshot": {
			Type:          schema.TypeList,
			ConflictsWith: codegen_configs.GetResourceByName("GlobalSnapshot").GetConflictingFields("owner_root_snapshot"),

			Required:    true,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,

			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{

					"clone_id": {
						Type:          schema.TypeInt,
						ConflictsWith: codegen_configs.GetResourceByName("GlobalSnapshotOwnerRootSnapshot").GetConflictingFields("clone_id"),

						Computed:    true,
						Optional:    false,
						Sensitive:   false,
						Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The ID of the clone`,
					},

					"name": {
						Type:          schema.TypeString,
						ConflictsWith: codegen_configs.GetResourceByName("GlobalSnapshotOwnerRootSnapshot").GetConflictingFields("name"),

						Required:    true,
						Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Remote Clone Name`,
					},

					"parent_handle_ehandle": {
						Type:          schema.TypeString,
						ConflictsWith: codegen_configs.GetResourceByName("GlobalSnapshotOwnerRootSnapshot").GetConflictingFields("parent_handle_ehandle"),

						Computed:    true,
						Optional:    false,
						Sensitive:   false,
						Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The remote handle (inode)`,
					},
				},
			},
		},

		"owner_tenant": {
			Type:          schema.TypeList,
			ConflictsWith: codegen_configs.GetResourceByName("GlobalSnapshot").GetConflictingFields("owner_tenant"),

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

var GlobalSnapshotNamesMapping = map[string][]string{}

func ResourceGlobalSnapshotReadStructIntoSchema(ctx context.Context, resource api_latest.GlobalSnapshot, d *schema.ResourceData) diag.Diagnostics {
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

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "RemoteTargetId", resource.RemoteTargetId))

	err = d.Set("remote_target_id", resource.RemoteTargetId)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"remote_target_id\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "RemoteTargetGuid", resource.RemoteTargetGuid))

	err = d.Set("remote_target_guid", resource.RemoteTargetGuid)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"remote_target_guid\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "RemoteTargetPath", resource.RemoteTargetPath))

	err = d.Set("remote_target_path", resource.RemoteTargetPath)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"remote_target_path\"",
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

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "OwnerRootSnapshot", resource.OwnerRootSnapshot))

	tflog.Debug(ctx, fmt.Sprintf("Found a pointer object %v", resource.OwnerRootSnapshot))
	err = d.Set("owner_root_snapshot", utils.FlattenModelAsList(ctx, resource.OwnerRootSnapshot))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"owner_root_snapshot\"",
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
func resourceGlobalSnapshotRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*vast_client.VMSSession)
	resourceConfig := codegen_configs.GetResourceByName("GlobalSnapshot")
	attrs := map[string]interface{}{"path": utils.GenPath("globalsnapstreams"), "id": d.Id()}
	tflog.Debug(ctx, fmt.Sprintf("[resourceGlobalSnapshotRead] Calling Get Function : %v for resource GlobalSnapshot", utils.GetFuncName(resourceConfig.GetFunc)))
	response, err := resourceConfig.GetFunc(ctx, client, attrs, d, map[string]string{})
	utils.VastVersionsWarn(ctx)

	var body []byte
	var resource api_latest.GlobalSnapshot
	if err != nil && response != nil && response.StatusCode == 404 && !resourceConfig.DisableFallbackRequest {
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
	diags = ResourceGlobalSnapshotReadStructIntoSchema(ctx, resource, d)

	return diags
}

func resourceGlobalSnapshotDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*vast_client.VMSSession)
	resourceConfig := codegen_configs.GetResourceByName("GlobalSnapshot")
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

func resourceGlobalSnapshotCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	namesMapping := utils.ContextKey("names_mapping")
	newCtx := context.WithValue(ctx, namesMapping, GlobalSnapshotNamesMapping)
	var diags diag.Diagnostics
	data := make(map[string]interface{})
	client := m.(*vast_client.VMSSession)
	resourceConfig := codegen_configs.GetResourceByName("GlobalSnapshot")
	tflog.Info(ctx, fmt.Sprintf("Creating Resource GlobalSnapshot"))
	reflectGlobalSnapshot := reflect.TypeOf((*api_latest.GlobalSnapshot)(nil))
	utils.PopulateResourceMap(newCtx, reflectGlobalSnapshot.Elem(), d, &data, "", false)

	var before_post_error error
	data, before_post_error = resourceConfig.BeforePostFunc(data, client, ctx, d)
	if before_post_error != nil {
		return diag.FromErr(before_post_error)
	}

	versionsEqual := utils.VastVersionsWarn(ctx)

	if versionsEqual != metadata.CLUSTER_VERSION_EQUALS {
		clusterVersion := metadata.ClusterVersionString()
		t, typeExists := vast_versions.GetVersionedType(clusterVersion, "GlobalSnapshot")
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
			tflog.Warn(ctx, fmt.Sprintf("Could have not found resource %s in version %s, things might not work properly", "GlobalSnapshot", clusterVersion))
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
	tflog.Info(ctx, fmt.Sprintf("Server Error for  GlobalSnapshot %v", createErr))

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
	resource := api_latest.GlobalSnapshot{}
	err = json.Unmarshal(responseBody, &resource)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to convert response body into GlobalSnapshot",
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
	resourceGlobalSnapshotRead(ctxWithResource, d, m)

	return diags
}

func resourceGlobalSnapshotUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	namesMapping := utils.ContextKey("names_mapping")
	newCtx := context.WithValue(ctx, namesMapping, GlobalSnapshotNamesMapping)
	var diags diag.Diagnostics
	data := make(map[string]interface{})
	versionsEqual := utils.VastVersionsWarn(ctx)
	resourceConfig := codegen_configs.GetResourceByName("GlobalSnapshot")
	if versionsEqual != metadata.CLUSTER_VERSION_EQUALS {
		clusterVersion := metadata.ClusterVersionString()
		t, typeExists := vast_versions.GetVersionedType(clusterVersion, "GlobalSnapshot")
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
			tflog.Warn(ctx, fmt.Sprintf("Could have not found resource %s in version %s, things might not work properly", "GlobalSnapshot", clusterVersion))
		}
	}

	client := m.(*vast_client.VMSSession)
	tflog.Info(ctx, fmt.Sprintf("Updating Resource GlobalSnapshot"))
	reflectGlobalSnapshot := reflect.TypeOf((*api_latest.GlobalSnapshot)(nil))
	utils.PopulateResourceMap(newCtx, reflectGlobalSnapshot.Elem(), d, &data, "", false)

	var beforePatchError error
	data, beforePatchError = resourceConfig.BeforePatchFunc(data, client, ctx, d)
	if beforePatchError != nil {
		return diag.FromErr(beforePatchError)
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
	attrs := map[string]interface{}{"path": utils.GenPath("globalsnapstreams"), "id": d.Id()}
	response, patchErr := resourceConfig.UpdateFunc(ctx, client, attrs, data, d, map[string]string{})
	tflog.Info(ctx, fmt.Sprintf("Server Error for  GlobalSnapshot %v", patchErr))
	if patchErr != nil {
		errorMessage := fmt.Sprintf("server response:\n%v\nUnderlying error:\n%v", utils.GetResponseBodyAsStr(response), patchErr.Error())
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Object Creation Failed",
			Detail:   errorMessage,
		})
		return diags
	}
	resourceGlobalSnapshotRead(ctx, d, m)

	return diags

}

func resourceGlobalSnapshotImporter(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {

	var result []*schema.ResourceData
	client := m.(*vast_client.VMSSession)
	resourceConfig := codegen_configs.GetResourceByName("GlobalSnapshot")
	attrs := map[string]interface{}{"path": utils.GenPath("globalsnapstreams")}
	response, err := resourceConfig.ImportFunc(ctx, client, attrs, d, resourceConfig.Importer.GetFunc())

	if err != nil {
		return result, err
	}

	var resourceList []api_latest.GlobalSnapshot
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

	diags := ResourceGlobalSnapshotReadStructIntoSchema(ctx, resource, d)
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
