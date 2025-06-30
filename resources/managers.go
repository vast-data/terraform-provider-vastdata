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

func ResourceManager() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceManagerRead,
		DeleteContext: resourceManagerDelete,
		CreateContext: resourceManagerCreate,
		UpdateContext: resourceManagerUpdate,

		Importer: &schema.ResourceImporter{
			StateContext: resourceManagerImporter,
		},

		Description: ``,
		Schema:      getResourceManagerSchema(),
	}
}

func getResourceManagerSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{

		"guid": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Manager").GetConflictingFields("guid"),

			Computed:    true,
			Optional:    false,
			Sensitive:   false,
			Description: `(Valid for versions: 5.2.0) A uniqe GUID assigned to the manager`,
		},

		"username": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Manager").GetConflictingFields("username"),

			Required:    true,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The username of the manager`,
		},

		"password": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Manager").GetConflictingFields("password"),

			DiffSuppressOnRefresh: false,
			DiffSuppressFunc:      codegen_configs.GetResourceByName("Manager").GetAttributeDiffFunc("password"),
			Computed:              true,
			Optional:              true,
			Sensitive:             true,

			ValidateDiagFunc: utils.ValidateManagerPassword,
		},

		"first_name": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Manager").GetConflictingFields("first_name"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The user firstname`,
		},

		"last_name": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Manager").GetConflictingFields("last_name"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The user last name`,
		},

		"permissions_list": {
			Type:          schema.TypeList,
			ConflictsWith: codegen_configs.GetResourceByName("Manager").GetConflictingFields("permissions_list"),

			DiffSuppressOnRefresh: false,
			DiffSuppressFunc:      codegen_configs.GetResourceByName("Manager").GetAttributeDiffFunc("permissions_list"),
			Computed:              true,
			Optional:              true,
			Sensitive:             false,
			Description:           `(Valid for versions: 5.0.0,5.1.0,5.2.0) List of allowed permissions Allowed Values are [create_support create_settings create_security create_monitoring create_logical create_hardware create_events create_database create_applications view_support view_settings view_security view_monitoring view_logical view_hardware view_events view_applications view_database edit_support edit_settings edit_security edit_monitoring edit_logical edit_hardware edit_events edit_database edit_applications delete_support delete_settings delete_security delete_monitoring delete_logical delete_hardware delete_events delete_applications delete_database]`,

			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},

		"roles": {
			Type:          schema.TypeList,
			ConflictsWith: codegen_configs.GetResourceByName("Manager").GetConflictingFields("roles"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) List of roles ids`,

			Elem: &schema.Schema{
				Type: schema.TypeInt,
			},
		},

		"password_expiration_disabled": {
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("Manager").GetConflictingFields("password_expiration_disabled"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Disable password expiration`,
		},

		"is_temporary_password": {
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("Manager").GetConflictingFields("is_temporary_password"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) If this set to true next time that a user will login he will be promped to replace his password`,
		},

		"permissions": {
			Type:          schema.TypeList,
			ConflictsWith: codegen_configs.GetResourceByName("Manager").GetConflictingFields("permissions"),

			Computed:    true,
			Optional:    false,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) List of allowed permissions returned from the VMS`,

			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}

var ManagerNamesMapping = map[string][]string{}

func ResourceManagerReadStructIntoSchema(ctx context.Context, resource api_latest.Manager, d *schema.ResourceData) diag.Diagnostics {
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

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Username", resource.Username))

	err = d.Set("username", resource.Username)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"username\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Password", resource.Password))

	err = d.Set("password", resource.Password)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"password\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "FirstName", resource.FirstName))

	err = d.Set("first_name", resource.FirstName)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"first_name\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "LastName", resource.LastName))

	err = d.Set("last_name", resource.LastName)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"last_name\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "PermissionsList", resource.PermissionsList))

	err = d.Set("permissions_list", utils.FlattenListOfPrimitives(&resource.PermissionsList))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"permissions_list\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Roles", resource.Roles))

	err = d.Set("roles", utils.FlattenListOfPrimitives(&resource.Roles))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"roles\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "PasswordExpirationDisabled", resource.PasswordExpirationDisabled))

	err = d.Set("password_expiration_disabled", resource.PasswordExpirationDisabled)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"password_expiration_disabled\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "IsTemporaryPassword", resource.IsTemporaryPassword))

	err = d.Set("is_temporary_password", resource.IsTemporaryPassword)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"is_temporary_password\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Permissions", resource.Permissions))

	err = d.Set("permissions", utils.FlattenListOfPrimitives(&resource.Permissions))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"permissions\"",
			Detail:   err.Error(),
		})
	}

	return diags

}
func resourceManagerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*vast_client.VMSSession)
	resourceConfig := codegen_configs.GetResourceByName("Manager")
	attrs := map[string]interface{}{"path": utils.GenPath("managers"), "id": d.Id()}
	tflog.Debug(ctx, fmt.Sprintf("[resourceManagerRead] Calling Get Function : %v for resource Manager", utils.GetFuncName(resourceConfig.GetFunc)))
	response, err := resourceConfig.GetFunc(ctx, client, attrs, d, map[string]string{})
	utils.VastVersionsWarn(ctx)

	var body []byte
	var resource api_latest.Manager
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
		body, err = resourceConfig.ResponseProcessingFunc(ctx, response, d)
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
	diags = ResourceManagerReadStructIntoSchema(ctx, resource, d)

	var after_read_error error
	after_read_error = resourceConfig.AfterReadFunc(client, ctx, d)
	if after_read_error != nil {
		return diag.FromErr(after_read_error)
	}

	return diags
}

func resourceManagerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*vast_client.VMSSession)
	resourceConfig := codegen_configs.GetResourceByName("Manager")
	attrs := map[string]interface{}{"path": utils.GenPath("managers"), "id": d.Id()}

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

func resourceManagerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	namesMapping := utils.ContextKey("names_mapping")
	newCtx := context.WithValue(ctx, namesMapping, ManagerNamesMapping)
	var diags diag.Diagnostics
	data := make(map[string]interface{})
	client := m.(*vast_client.VMSSession)
	resourceConfig := codegen_configs.GetResourceByName("Manager")
	tflog.Info(ctx, fmt.Sprintf("Creating Resource Manager"))
	reflectManager := reflect.TypeOf((*api_latest.Manager)(nil))
	utils.PopulateResourceMap(newCtx, reflectManager.Elem(), d, &data, "", false)

	versionsEqual := utils.VastVersionsWarn(ctx)

	if versionsEqual != metadata.CLUSTER_VERSION_EQUALS {
		clusterVersion := metadata.ClusterVersionString()
		t, typeExists := vast_versions.GetVersionedType(clusterVersion, "Manager")
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
			tflog.Warn(ctx, fmt.Sprintf("Could have not found resource %s in version %s, things might not work properly", "Manager", clusterVersion))
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
	attrs := map[string]interface{}{"path": utils.GenPath("managers")}
	response, createErr := resourceConfig.CreateFunc(ctx, client, attrs, data, map[string]string{})
	tflog.Info(ctx, fmt.Sprintf("Server Error for  Manager %v", createErr))

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
	resource := api_latest.Manager{}
	err = json.Unmarshal(responseBody, &resource)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to convert response body into Manager",
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
	resourceManagerRead(ctxWithResource, d, m)

	return diags
}

func resourceManagerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	namesMapping := utils.ContextKey("names_mapping")
	newCtx := context.WithValue(ctx, namesMapping, ManagerNamesMapping)
	var diags diag.Diagnostics
	data := make(map[string]interface{})
	versionsEqual := utils.VastVersionsWarn(ctx)
	resourceConfig := codegen_configs.GetResourceByName("Manager")
	if versionsEqual != metadata.CLUSTER_VERSION_EQUALS {
		clusterVersion := metadata.ClusterVersionString()
		t, typeExists := vast_versions.GetVersionedType(clusterVersion, "Manager")
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
			tflog.Warn(ctx, fmt.Sprintf("Could have not found resource %s in version %s, things might not work properly", "Manager", clusterVersion))
		}
	}

	client := m.(*vast_client.VMSSession)
	tflog.Info(ctx, fmt.Sprintf("Updating Resource Manager"))
	reflectManager := reflect.TypeOf((*api_latest.Manager)(nil))
	utils.PopulateResourceMap(newCtx, reflectManager.Elem(), d, &data, "", false)

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
	attrs := map[string]interface{}{"path": utils.GenPath("managers"), "id": d.Id()}
	response, patchErr := resourceConfig.UpdateFunc(ctx, client, attrs, data, d, map[string]string{})
	tflog.Info(ctx, fmt.Sprintf("Server Error for  Manager %v", patchErr))
	if patchErr != nil {
		errorMessage := fmt.Sprintf("server response:\n%v\nUnderlying error:\n%v", utils.GetResponseBodyAsStr(response), patchErr.Error())
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Object Creation Failed",
			Detail:   errorMessage,
		})
		return diags
	}
	resourceManagerRead(ctx, d, m)

	return diags

}

func resourceManagerImporter(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {

	var result []*schema.ResourceData
	client := m.(*vast_client.VMSSession)
	resourceConfig := codegen_configs.GetResourceByName("Manager")
	attrs := map[string]interface{}{"path": utils.GenPath("managers")}
	response, err := resourceConfig.ImportFunc(ctx, client, attrs, d, resourceConfig.Importer.GetFunc())

	if err != nil {
		return result, err
	}

	var resourceList []api_latest.Manager
	body, err := resourceConfig.ResponseProcessingFunc(ctx, response, d)

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

	diags := ResourceManagerReadStructIntoSchema(ctx, resource, d)
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
