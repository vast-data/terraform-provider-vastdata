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

func ResourceProtectionPolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceProtectionPolicyRead,
		DeleteContext: resourceProtectionPolicyDelete,
		CreateContext: resourceProtectionPolicyCreate,
		UpdateContext: resourceProtectionPolicyUpdate,

		Importer: &schema.ResourceImporter{
			StateContext: resourceProtectionPolicyImporter,
		},

		Description: ``,
		Schema:      getResourceProtectionPolicySchema(),
	}
}

func getResourceProtectionPolicySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{

		"guid": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("ProtectionPolicy").GetConflictingFields("guid"),

			Computed:    true,
			Optional:    false,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) A unique guid given to the  replication peer configuration`,
		},

		"name": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("ProtectionPolicy").GetConflictingFields("name"),

			Required:    true,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The name of the replication peer configuration`,
		},

		"url": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("ProtectionPolicy").GetConflictingFields("url"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Direct link to the replication policy`,
		},

		"target_name": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("ProtectionPolicy").GetConflictingFields("target_name"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The target peer name`,
		},

		"target_object_id": {
			Type:          schema.TypeInt,
			ConflictsWith: codegen_configs.GetResourceByName("ProtectionPolicy").GetConflictingFields("target_object_id"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The id of the target peer`,
		},

		"prefix": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("ProtectionPolicy").GetConflictingFields("prefix"),

			Required:    true,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The prefix to be given to the replicated data`,
		},

		"clone_type": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("ProtectionPolicy").GetConflictingFields("clone_type"),

			Required: true,

			ValidateDiagFunc: utils.OneOf([]string{"NATIVE_REPLICATION", "LOCAL"}),
			Description:      `(Valid for versions: 5.0.0,5.1.0,5.2.0) The type the replication Allowed Values are [NATIVE_REPLICATION LOCAL]`,
		},

		"frames": {
			Type:          schema.TypeList,
			ConflictsWith: codegen_configs.GetResourceByName("ProtectionPolicy").GetConflictingFields("frames"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) List of snapshots schedules`,

			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{

					"every": {
						Type:          schema.TypeString,
						ConflictsWith: codegen_configs.GetResourceByName("ProtectionPolicySchedule").GetConflictingFields("every"),

						DiffSuppressOnRefresh: false,
						DiffSuppressFunc:      codegen_configs.GetResourceByName("ProtectionPolicySchedule").GetAttributeDiffFunc("every"),
						Computed:              true,
						Optional:              true,
						Sensitive:             false,

						ValidateDiagFunc: utils.ProtectionPolicyTimeIntervalValidation,
					},

					"start_at": {
						Type:          schema.TypeString,
						ConflictsWith: codegen_configs.GetResourceByName("ProtectionPolicySchedule").GetConflictingFields("start_at"),

						Computed:  true,
						Optional:  true,
						Sensitive: false,

						ValidateDiagFunc: utils.ProtectionPolicyStartAt,
					},

					"keep_local": {
						Type:          schema.TypeString,
						ConflictsWith: codegen_configs.GetResourceByName("ProtectionPolicySchedule").GetConflictingFields("keep_local"),

						DiffSuppressOnRefresh: false,
						DiffSuppressFunc:      codegen_configs.GetResourceByName("ProtectionPolicySchedule").GetAttributeDiffFunc("keep_local"),
						Computed:              true,
						Optional:              true,
						Sensitive:             false,

						ValidateDiagFunc: utils.ProtectionPolicyTimeIntervalValidation,
					},

					"keep_remote": {
						Type:          schema.TypeString,
						ConflictsWith: codegen_configs.GetResourceByName("ProtectionPolicySchedule").GetConflictingFields("keep_remote"),

						DiffSuppressOnRefresh: false,
						DiffSuppressFunc:      codegen_configs.GetResourceByName("ProtectionPolicySchedule").GetAttributeDiffFunc("keep_remote"),
						Computed:              true,
						Optional:              true,
						Sensitive:             false,

						ValidateDiagFunc: utils.ProtectionPolicyTimeIntervalValidation,
					},
				},
			},
		},

		"indestructible": {
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("ProtectionPolicy").GetConflictingFields("indestructible"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Is the snapshot indestructable`,
		},
	}
}

var ProtectionPolicyNamesMapping = map[string][]string{}

func ResourceProtectionPolicyReadStructIntoSchema(ctx context.Context, resource api_latest.ProtectionPolicy, d *schema.ResourceData) diag.Diagnostics {
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

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Url", resource.Url))

	err = d.Set("url", resource.Url)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"url\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "TargetName", resource.TargetName))

	err = d.Set("target_name", resource.TargetName)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"target_name\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "TargetObjectId", resource.TargetObjectId))

	err = d.Set("target_object_id", resource.TargetObjectId)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"target_object_id\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Prefix", resource.Prefix))

	err = d.Set("prefix", resource.Prefix)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"prefix\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "CloneType", resource.CloneType))

	err = d.Set("clone_type", resource.CloneType)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"clone_type\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Frames", resource.Frames))

	err = d.Set("frames", utils.FlattenListOfModelsToList(ctx, resource.Frames))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"frames\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Indestructible", resource.Indestructible))

	err = d.Set("indestructible", resource.Indestructible)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"indestructible\"",
			Detail:   err.Error(),
		})
	}

	return diags

}
func resourceProtectionPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*vast_client.VMSSession)
	resourceConfig := codegen_configs.GetResourceByName("ProtectionPolicy")
	attrs := map[string]interface{}{"path": utils.GenPath("protectionpolicies"), "id": d.Id()}
	tflog.Debug(ctx, fmt.Sprintf("[resourceProtectionPolicyRead] Calling Get Function : %v for resource ProtectionPolicy", utils.GetFuncName(resourceConfig.GetFunc)))
	response, err := resourceConfig.GetFunc(ctx, client, attrs, d, map[string]string{})
	utils.VastVersionsWarn(ctx)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred while obtaining data from the VAST Data cluster",
			Detail:   err.Error(),
		})
		return diags

	}
	tflog.Info(ctx, response.Request.URL.String())
	resource := api_latest.ProtectionPolicy{}
	body, err := resourceConfig.ResponseProcessingFunc(ctx, response)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred reading data received from VAST Data cluster",
			Detail:   err.Error(),
		})
		return diags

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
	diags = ResourceProtectionPolicyReadStructIntoSchema(ctx, resource, d)

	return diags
}

func resourceProtectionPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*vast_client.VMSSession)
	resourceConfig := codegen_configs.GetResourceByName("ProtectionPolicy")
	attrs := map[string]interface{}{"path": utils.GenPath("protectionpolicies"), "id": d.Id()}

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

func resourceProtectionPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	namesMapping := utils.ContextKey("namesMapping")
	newCtx := context.WithValue(ctx, namesMapping, ProtectionPolicyNamesMapping)
	var diags diag.Diagnostics
	data := make(map[string]interface{})
	client := m.(*vast_client.VMSSession)
	resourceConfig := codegen_configs.GetResourceByName("ProtectionPolicy")
	tflog.Info(ctx, fmt.Sprintf("Creating Resource ProtectionPolicy"))
	reflectProtectionPolicy := reflect.TypeOf((*api_latest.ProtectionPolicy)(nil))
	utils.PopulateResourceMap(newCtx, reflectProtectionPolicy.Elem(), d, &data, "", false)

	versionsEqual := utils.VastVersionsWarn(ctx)

	if versionsEqual != metadata.CLUSTER_VERSION_EQUALS {
		clusterVersion := metadata.ClusterVersionString()
		t, typeExists := vast_versions.GetVersionedType(clusterVersion, "ProtectionPolicy")
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
			tflog.Warn(ctx, fmt.Sprintf("Could have not found resource %s in version %s , things might not work properly", "ProtectionPolicy", clusterVersion))
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
	attrs := map[string]interface{}{"path": utils.GenPath("protectionpolicies")}
	response, createErr := resourceConfig.CreateFunc(ctx, client, attrs, data, map[string]string{})
	tflog.Info(ctx, fmt.Sprintf("Server Error for  ProtectionPolicy %v", createErr))

	if createErr != nil {
		errorMessage := createErr.Error() + " Server Response: " + utils.GetResponseBodyAsStr(response)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Object Creation Failed",
			Detail:   errorMessage,
		})
		return diags
	}
	responseBody, _ := io.ReadAll(response.Body)
	tflog.Debug(ctx, fmt.Sprintf("Object created , server response %v", string(responseBody)))
	resource := api_latest.ProtectionPolicy{}
	err = json.Unmarshal(responseBody, &resource)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to convert response body into ProtectionPolicy",
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
	resourceProtectionPolicyRead(ctxWithResource, d, m)

	return diags
}

func resourceProtectionPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	namesMapping := utils.ContextKey("namesMapping")
	newCtx := context.WithValue(ctx, namesMapping, ProtectionPolicyNamesMapping)
	var diags diag.Diagnostics
	data := make(map[string]interface{})
	versionsEqual := utils.VastVersionsWarn(ctx)
	resourceConfig := codegen_configs.GetResourceByName("ProtectionPolicy")
	if versionsEqual != metadata.CLUSTER_VERSION_EQUALS {
		clusterVersion := metadata.ClusterVersionString()
		t, typeExists := vast_versions.GetVersionedType(clusterVersion, "ProtectionPolicy")
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
			tflog.Warn(ctx, fmt.Sprintf("Could have not found resource %s in version %s , things might not work properly", "ProtectionPolicy", clusterVersion))
		}
	}

	client := m.(*vast_client.VMSSession)
	tflog.Info(ctx, fmt.Sprintf("Updating Resource ProtectionPolicy"))
	reflectProtectionPolicy := reflect.TypeOf((*api_latest.ProtectionPolicy)(nil))
	utils.PopulateResourceMap(newCtx, reflectProtectionPolicy.Elem(), d, &data, "", false)

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
	attrs := map[string]interface{}{"path": utils.GenPath("protectionpolicies"), "id": d.Id()}
	response, patchErr := resourceConfig.UpdateFunc(ctx, client, attrs, data, d, map[string]string{})
	tflog.Info(ctx, fmt.Sprintf("Server Error for  ProtectionPolicy %v", patchErr))
	if patchErr != nil {
		errorMessage := patchErr.Error() + " Server Response: " + utils.GetResponseBodyAsStr(response)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Object Creation Failed",
			Detail:   errorMessage,
		})
		return diags
	}
	resourceProtectionPolicyRead(ctx, d, m)

	return diags

}

func resourceProtectionPolicyImporter(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {

	var result []*schema.ResourceData
	client := m.(*vast_client.VMSSession)
	resourceConfig := codegen_configs.GetResourceByName("ProtectionPolicy")
	attrs := map[string]interface{}{"path": utils.GenPath("protectionpolicies")}
	response, err := resourceConfig.ImportFunc(ctx, client, attrs, d, resourceConfig.Importer.GetFunc())

	if err != nil {
		return result, err
	}

	var resourceList []api_latest.ProtectionPolicy
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

	diags := ResourceProtectionPolicyReadStructIntoSchema(ctx, resource, d)
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
