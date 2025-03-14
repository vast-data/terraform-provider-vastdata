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

		"guid": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Manager").GetConflictingFields("guid"),

			Computed:    true,
			Optional:    false,
			Sensitive:   false,
			Description: `(Valid for versions: 5.2.0) A uniqe GUID assigned to the manager`,
		},

		"username": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Manager").GetConflictingFields("username"),

			Required:    true,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The username of the manager`,
		},

		"password": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Manager").GetConflictingFields("password"),

			DiffSuppressOnRefresh: false,
			DiffSuppressFunc:      codegen_configs.GetResourceByName("Manager").GetAttributeDiffFunc("password"),
			Computed:              true,
			Optional:              true,
			Sensitive:             true,

			ValidateDiagFunc: utils.ValidateManagerPassword,
		},

		"first_name": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Manager").GetConflictingFields("first_name"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The user firstname`,
		},

		"last_name": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Manager").GetConflictingFields("last_name"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The user last name`,
		},

		"permissions_list": &schema.Schema{
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

		"roles": &schema.Schema{
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

		"password_expiration_disabled": &schema.Schema{
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("Manager").GetConflictingFields("password_expiration_disabled"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Disable password expiration`,
		},

		"is_temporary_password": &schema.Schema{
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("Manager").GetConflictingFields("is_temporary_password"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) If this set to true next time that a user will login he will be promped to replace his password`,
		},

		"permissions": &schema.Schema{
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

var Manager_names_mapping map[string][]string = map[string][]string{}

func ResourceManagerReadStructIntoSchema(ctx context.Context, resource api_latest.Manager, d *schema.ResourceData) diag.Diagnostics {
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

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Username", resource.Username))

	err = d.Set("username", resource.Username)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"username\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Password", resource.Password))

	err = d.Set("password", resource.Password)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"password\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "FirstName", resource.FirstName))

	err = d.Set("first_name", resource.FirstName)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"first_name\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "LastName", resource.LastName))

	err = d.Set("last_name", resource.LastName)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"last_name\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "PermissionsList", resource.PermissionsList))

	err = d.Set("permissions_list", utils.FlattenListOfPrimitives(&resource.PermissionsList))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"permissions_list\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Roles", resource.Roles))

	err = d.Set("roles", utils.FlattenListOfPrimitives(&resource.Roles))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"roles\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "PasswordExpirationDisabled", resource.PasswordExpirationDisabled))

	err = d.Set("password_expiration_disabled", resource.PasswordExpirationDisabled)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"password_expiration_disabled\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "IsTemporaryPassword", resource.IsTemporaryPassword))

	err = d.Set("is_temporary_password", resource.IsTemporaryPassword)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"is_temporary_password\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Permissions", resource.Permissions))

	err = d.Set("permissions", utils.FlattenListOfPrimitives(&resource.Permissions))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"permissions\"",
			Detail:   err.Error(),
		})
	}

	return diags

}
func resourceManagerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(vast_client.JwtSession)
	resource_config := codegen_configs.GetResourceByName("Manager")
	attrs := map[string]interface{}{"path": utils.GenPath("managers"), "id": d.Id()}
	tflog.Debug(ctx, fmt.Sprintf("[resourceManagerRead] Calling Get Function : %v for resource Manager", utils.GetFuncName(resource_config.GetFunc)))
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
	resource := api_latest.Manager{}
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
	diags = ResourceManagerReadStructIntoSchema(ctx, resource, d)

	var after_read_error error
	after_read_error = resource_config.AfterReadFunc(client, ctx, d)
	if after_read_error != nil {
		return diag.FromErr(after_read_error)
	}

	return diags
}

func resourceManagerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(vast_client.JwtSession)
	resource_config := codegen_configs.GetResourceByName("Manager")
	attrs := map[string]interface{}{"path": utils.GenPath("managers"), "id": d.Id()}

	response, err := resource_config.DeleteFunc(ctx, client, attrs, nil, map[string]string{})

	tflog.Info(ctx, fmt.Sprintf("Removing Resource"))
	if response != nil {
		tflog.Info(ctx, response.Request.URL.String())
		tflog.Info(ctx, utils.GetResponseBodyAsStr(response))
	}

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured while deleting a resource from the vastdata cluster",
			Detail:   err.Error(),
		})

	}

	return diags

}

func resourceManagerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	names_mapping := utils.ContextKey("names_mapping")
	new_ctx := context.WithValue(ctx, names_mapping, Manager_names_mapping)
	var diags diag.Diagnostics
	data := make(map[string]interface{})
	client := m.(vast_client.JwtSession)
	resource_config := codegen_configs.GetResourceByName("Manager")
	tflog.Info(ctx, fmt.Sprintf("Creating Resource Manager"))
	reflect_Manager := reflect.TypeOf((*api_latest.Manager)(nil))
	utils.PopulateResourceMap(new_ctx, reflect_Manager.Elem(), d, &data, "", false)

	version_compare := utils.VastVersionsWarn(ctx)

	if version_compare != metadata.CLUSTER_VERSION_EQUALS {
		cluster_version := metadata.ClusterVersionString()
		t, t_exists := vast_versions.GetVersionedType(cluster_version, "Manager")
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
			tflog.Warn(ctx, fmt.Sprintf("Could have not found resource %s in version %s , things might not work properly", "Manager", cluster_version))
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
	response, create_err := resource_config.CreateFunc(ctx, client, attrs, data, map[string]string{})
	tflog.Info(ctx, fmt.Sprintf("Server Error for  Manager %v", create_err))

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
	resource := api_latest.Manager{}
	err = json.Unmarshal(response_body, &resource)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to convert response body into Manager",
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
	resourceManagerRead(ctx_with_resource, d, m)

	return diags
}

func resourceManagerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	names_mapping := utils.ContextKey("names_mapping")
	new_ctx := context.WithValue(ctx, names_mapping, Manager_names_mapping)
	var diags diag.Diagnostics
	data := make(map[string]interface{})
	version_compare := utils.VastVersionsWarn(ctx)
	resource_config := codegen_configs.GetResourceByName("Manager")
	if version_compare != metadata.CLUSTER_VERSION_EQUALS {
		cluster_version := metadata.ClusterVersionString()
		t, t_exists := vast_versions.GetVersionedType(cluster_version, "Manager")
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
			tflog.Warn(ctx, fmt.Sprintf("Could have not found resource %s in version %s , things might not work properly", "Manager", cluster_version))
		}
	}

	client := m.(vast_client.JwtSession)
	tflog.Info(ctx, fmt.Sprintf("Updating Resource Manager"))
	reflect_Manager := reflect.TypeOf((*api_latest.Manager)(nil))
	utils.PopulateResourceMap(new_ctx, reflect_Manager.Elem(), d, &data, "", false)

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
	response, patch_err := resource_config.UpdateFunc(ctx, client, attrs, data, d, map[string]string{})
	tflog.Info(ctx, fmt.Sprintf("Server Error for  Manager %v", patch_err))
	if patch_err != nil {
		error_message := patch_err.Error() + " Server Response: " + utils.GetResponseBodyAsStr(response)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Object Creation Failed",
			Detail:   error_message,
		})
		return diags
	}
	resourceManagerRead(ctx, d, m)

	return diags

}

func resourceManagerImporter(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {

	result := []*schema.ResourceData{}
	client := m.(vast_client.JwtSession)
	resource_config := codegen_configs.GetResourceByName("Manager")
	attrs := map[string]interface{}{"path": utils.GenPath("managers")}
	response, err := resource_config.ImportFunc(ctx, client, attrs, d, resource_config.Importer.GetFunc())

	if err != nil {
		return result, err
	}

	resource_l := []api_latest.Manager{}
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

	diags := ResourceManagerReadStructIntoSchema(ctx, resource, d)
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
