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

func ResourceRole() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceRoleRead,
		DeleteContext: resourceRoleDelete,
		CreateContext: resourceRoleCreate,
		UpdateContext: resourceRoleUpdate,

		Importer: &schema.ResourceImporter{
			StateContext: resourceRoleImporter,
		},

		Description: ``,
		Schema:      getResourceRoleSchema(),
	}
}

func getResourceRoleSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{

		"guid": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Role").GetConflictingFields("guid"),

			Computed:    true,
			Optional:    false,
			Sensitive:   false,
			Description: `(Valid for versions: 5.2.0) A uniqe GUID assigned to the role`,
		},

		"name": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Role").GetConflictingFields("name"),

			Required:    true,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) A uniqe name of the role`,
		},

		"permissions_list": {
			Type:          schema.TypeList,
			ConflictsWith: codegen_configs.GetResourceByName("Role").GetConflictingFields("permissions_list"),

			DiffSuppressOnRefresh: false,
			DiffSuppressFunc:      codegen_configs.GetResourceByName("Role").GetAttributeDiffFunc("permissions_list"),
			Computed:              true,
			Optional:              true,
			Sensitive:             false,
			Description:           `(Valid for versions: 5.0.0,5.1.0,5.2.0) List of allowed permissions Allowed Values are [create_support create_settings create_security create_monitoring create_logical create_hardware create_events create_database create_applications view_support view_settings view_security view_monitoring view_logical view_hardware view_events view_applications view_database edit_support edit_settings edit_security edit_monitoring edit_logical edit_hardware edit_events edit_database edit_applications delete_support delete_settings delete_security delete_monitoring delete_logical delete_hardware delete_events delete_applications delete_database]`,

			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},

		"permissions": {
			Type:          schema.TypeList,
			ConflictsWith: codegen_configs.GetResourceByName("Role").GetConflictingFields("permissions"),

			Computed:    true,
			Optional:    false,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) List of allowed permissions returned from the VMS`,

			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},

		"tenants": {
			Type:          schema.TypeList,
			ConflictsWith: codegen_configs.GetResourceByName("Role").GetConflictingFields("tenants"),

			Computed:    true,
			Optional:    false,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) List of tenants to which this role is associated with`,

			Elem: &schema.Schema{
				Type: schema.TypeInt,
			},
		},

		"is_admin": {
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("Role").GetConflictingFields("is_admin"),

			Computed:    true,
			Optional:    false,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Is the role is an admin role`,
		},

		"is_default": {
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("Role").GetConflictingFields("is_default"),

			Computed:    true,
			Optional:    false,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Is the role is a default role`,
		},

		"ldap_groups": {
			Type:          schema.TypeList,
			ConflictsWith: codegen_configs.GetResourceByName("Role").GetConflictingFields("ldap_groups"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) LDAP group(s) associated with the role. Members of the specified groups on a connected LDAP/Active Directory provider can access VMS and are granted whichever permissions are included in the role. A group can be associated with multiple roles.`,

			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}

var RoleNamesMapping = map[string][]string{}

func ResourceRoleReadStructIntoSchema(ctx context.Context, resource api_latest.Role, d *schema.ResourceData) diag.Diagnostics {
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

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "PermissionsList", resource.PermissionsList))

	err = d.Set("permissions_list", utils.FlattenListOfPrimitives(&resource.PermissionsList))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"permissions_list\"",
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

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Tenants", resource.Tenants))

	err = d.Set("tenants", utils.FlattenListOfPrimitives(&resource.Tenants))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"tenants\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "IsAdmin", resource.IsAdmin))

	err = d.Set("is_admin", resource.IsAdmin)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"is_admin\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "IsDefault", resource.IsDefault))

	err = d.Set("is_default", resource.IsDefault)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"is_default\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "LdapGroups", resource.LdapGroups))

	err = d.Set("ldap_groups", utils.FlattenListOfPrimitives(&resource.LdapGroups))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"ldap_groups\"",
			Detail:   err.Error(),
		})
	}

	return diags

}
func resourceRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*vast_client.VMSSession)
	resourceConfig := codegen_configs.GetResourceByName("Role")
	attrs := map[string]interface{}{"path": utils.GenPath("roles"), "id": d.Id()}
	tflog.Debug(ctx, fmt.Sprintf("[resourceRoleRead] Calling Get Function : %v for resource Role", utils.GetFuncName(resourceConfig.GetFunc)))
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
	resource := api_latest.Role{}
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
	diags = ResourceRoleReadStructIntoSchema(ctx, resource, d)

	var after_read_error error
	after_read_error = resourceConfig.AfterReadFunc(client, ctx, d)
	if after_read_error != nil {
		return diag.FromErr(after_read_error)
	}

	return diags
}

func resourceRoleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*vast_client.VMSSession)
	resourceConfig := codegen_configs.GetResourceByName("Role")
	attrs := map[string]interface{}{"path": utils.GenPath("roles"), "id": d.Id()}

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

func resourceRoleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	namesMapping := utils.ContextKey("namesMapping")
	newCtx := context.WithValue(ctx, namesMapping, RoleNamesMapping)
	var diags diag.Diagnostics
	data := make(map[string]interface{})
	client := m.(*vast_client.VMSSession)
	resourceConfig := codegen_configs.GetResourceByName("Role")
	tflog.Info(ctx, fmt.Sprintf("Creating Resource Role"))
	reflectRole := reflect.TypeOf((*api_latest.Role)(nil))
	utils.PopulateResourceMap(newCtx, reflectRole.Elem(), d, &data, "", false)

	versionsEqual := utils.VastVersionsWarn(ctx)

	if versionsEqual != metadata.CLUSTER_VERSION_EQUALS {
		clusterVersion := metadata.ClusterVersionString()
		t, typeExists := vast_versions.GetVersionedType(clusterVersion, "Role")
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
			tflog.Warn(ctx, fmt.Sprintf("Could have not found resource %s in version %s, things might not work properly", "Role", clusterVersion))
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
	attrs := map[string]interface{}{"path": utils.GenPath("roles")}
	response, createErr := resourceConfig.CreateFunc(ctx, client, attrs, data, map[string]string{})
	tflog.Info(ctx, fmt.Sprintf("Server Error for  Role %v", createErr))

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
	resource := api_latest.Role{}
	err = json.Unmarshal(responseBody, &resource)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to convert response body into Role",
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
	resourceRoleRead(ctxWithResource, d, m)

	return diags
}

func resourceRoleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	namesMapping := utils.ContextKey("namesMapping")
	newCtx := context.WithValue(ctx, namesMapping, RoleNamesMapping)
	var diags diag.Diagnostics
	data := make(map[string]interface{})
	versionsEqual := utils.VastVersionsWarn(ctx)
	resourceConfig := codegen_configs.GetResourceByName("Role")
	if versionsEqual != metadata.CLUSTER_VERSION_EQUALS {
		clusterVersion := metadata.ClusterVersionString()
		t, typeExists := vast_versions.GetVersionedType(clusterVersion, "Role")
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
			tflog.Warn(ctx, fmt.Sprintf("Could have not found resource %s in version %s, things might not work properly", "Role", clusterVersion))
		}
	}

	client := m.(*vast_client.VMSSession)
	tflog.Info(ctx, fmt.Sprintf("Updating Resource Role"))
	reflectRole := reflect.TypeOf((*api_latest.Role)(nil))
	utils.PopulateResourceMap(newCtx, reflectRole.Elem(), d, &data, "", false)

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
	attrs := map[string]interface{}{"path": utils.GenPath("roles"), "id": d.Id()}
	response, patchErr := resourceConfig.UpdateFunc(ctx, client, attrs, data, d, map[string]string{})
	tflog.Info(ctx, fmt.Sprintf("Server Error for  Role %v", patchErr))
	if patchErr != nil {
		errorMessage := fmt.Sprintf("server response:\n%v\nUnderlying error:\n%v", utils.GetResponseBodyAsStr(response), patchErr.Error())
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Object Creation Failed",
			Detail:   errorMessage,
		})
		return diags
	}
	resourceRoleRead(ctx, d, m)

	return diags

}

func resourceRoleImporter(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {

	var result []*schema.ResourceData
	client := m.(*vast_client.VMSSession)
	resourceConfig := codegen_configs.GetResourceByName("Role")
	attrs := map[string]interface{}{"path": utils.GenPath("roles")}
	response, err := resourceConfig.ImportFunc(ctx, client, attrs, d, resourceConfig.Importer.GetFunc())

	if err != nil {
		return result, err
	}

	var resourceList []api_latest.Role
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

	diags := ResourceRoleReadStructIntoSchema(ctx, resource, d)
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
