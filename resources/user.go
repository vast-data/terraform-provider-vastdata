package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"reflect"

	//        "net/url"
	"errors"
	codegen_configs "github.com/vast-data/terraform-provider-vastdata/codegen_tools/configs"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api_latest "github.com/vast-data/terraform-provider-vastdata/codegen/latest"
	metadata "github.com/vast-data/terraform-provider-vastdata/metadata"
	utils "github.com/vast-data/terraform-provider-vastdata/utils"
	vast_client "github.com/vast-data/terraform-provider-vastdata/vast-client"
	vast_versions "github.com/vast-data/terraform-provider-vastdata/vast_versions"
)

func ResourceUser() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceUserRead,
		DeleteContext: resourceUserDelete,
		CreateContext: resourceUserCreate,
		UpdateContext: resourceUserUpdate,

		Importer: &schema.ResourceImporter{
			StateContext: resourceUserImporter,
		},

		Description: ``,
		Schema:      getResourceUserSchema(),
	}
}

func getResourceUserSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{

		"guid": &schema.Schema{
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    false,
			Sensitive:   false,
			Description: `A uniq guid given to the user`,
		},

		"name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},

		"uid": &schema.Schema{
			Type:        schema.TypeInt,
			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `The user unix UID`,
		},

		"leading_gid": &schema.Schema{
			Type:        schema.TypeInt,
			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `The user leading unix GID`,
		},

		"gids": &schema.Schema{
			Type:        schema.TypeList,
			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `List of supplementary GID list`,

			Elem: &schema.Schema{
				Type: schema.TypeInt,
			},
		},

		"groups": &schema.Schema{
			Type:        schema.TypeList,
			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `List of supplementary Group list`,

			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},

		"group_count": &schema.Schema{
			Type:        schema.TypeInt,
			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `Group Count`,
		},

		"leading_group_name": &schema.Schema{
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `Leading Group Name`,
		},

		"leading_group_gid": &schema.Schema{
			Type:        schema.TypeInt,
			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `Leading Group GID`,
		},

		"sid": &schema.Schema{
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `The user SID`,
		},

		"primary_group_sid": &schema.Schema{
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `The user primary group SID`,
		},

		"sids": &schema.Schema{
			Type:        schema.TypeList,
			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `supplementary SID list`,

			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},

		"local": &schema.Schema{
			Type:        schema.TypeBool,
			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `IS this a local user`,
		},

		"access_keys": &schema.Schema{
			Type:        schema.TypeList,
			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `List of User Access Keys`,

			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{

					"access_key": &schema.Schema{
						Type:        schema.TypeString,
						Computed:    true,
						Optional:    true,
						Description: "",
					},

					"enabled": &schema.Schema{
						Type:        schema.TypeString,
						Computed:    true,
						Optional:    true,
						Description: "",
					},
				},
			},
		},

		"allow_create_bucket": &schema.Schema{
			Type:        schema.TypeBool,
			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `Allow create bucket`,
		},

		"allow_delete_bucket": &schema.Schema{
			Type:        schema.TypeBool,
			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `Allow delete bucket`,
		},

		"s3_superuser": &schema.Schema{
			Type:        schema.TypeBool,
			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `Is S3 superuser`,
		},

		"s3_policies_ids": &schema.Schema{
			Type:        schema.TypeList,
			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `List S3 policies IDs`,

			Elem: &schema.Schema{
				Type: schema.TypeInt,
			},
		},
	}
}

var User_names_mapping map[string][]string = map[string][]string{
	"access_keys": []string{"access_key", "enabled"},
}

func ResourceUserReadStructIntoSchema(ctx context.Context, resource api_latest.User, d *schema.ResourceData) diag.Diagnostics {
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

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Uid", resource.Uid))

	err = d.Set("uid", resource.Uid)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"uid\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "LeadingGid", resource.LeadingGid))

	err = d.Set("leading_gid", resource.LeadingGid)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"leading_gid\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Gids", resource.Gids))

	err = d.Set("gids", utils.FlattenListOfPrimitives(&resource.Gids))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"gids\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Groups", resource.Groups))

	err = d.Set("groups", utils.FlattenListOfPrimitives(&resource.Groups))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"groups\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "GroupCount", resource.GroupCount))

	err = d.Set("group_count", resource.GroupCount)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"group_count\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "LeadingGroupName", resource.LeadingGroupName))

	err = d.Set("leading_group_name", resource.LeadingGroupName)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"leading_group_name\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "LeadingGroupGid", resource.LeadingGroupGid))

	err = d.Set("leading_group_gid", resource.LeadingGroupGid)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"leading_group_gid\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Sid", resource.Sid))

	err = d.Set("sid", resource.Sid)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"sid\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "PrimaryGroupSid", resource.PrimaryGroupSid))

	err = d.Set("primary_group_sid", resource.PrimaryGroupSid)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"primary_group_sid\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Sids", resource.Sids))

	err = d.Set("sids", utils.FlattenListOfPrimitives(&resource.Sids))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"sids\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Local", resource.Local))

	err = d.Set("local", resource.Local)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"local\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "AccessKeys", resource.AccessKeys))

	err = d.Set("access_keys", utils.FlattenListOfStringsList(&resource.AccessKeys, []string{"access_key", "enabled"}))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"access_keys\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "AllowCreateBucket", resource.AllowCreateBucket))

	err = d.Set("allow_create_bucket", resource.AllowCreateBucket)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"allow_create_bucket\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "AllowDeleteBucket", resource.AllowDeleteBucket))

	err = d.Set("allow_delete_bucket", resource.AllowDeleteBucket)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"allow_delete_bucket\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "S3Superuser", resource.S3Superuser))

	err = d.Set("s3_superuser", resource.S3Superuser)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"s3_superuser\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "S3PoliciesIds", resource.S3PoliciesIds))

	err = d.Set("s3_policies_ids", utils.FlattenListOfPrimitives(&resource.S3PoliciesIds))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"s3_policies_ids\"",
			Detail:   err.Error(),
		})
	}

	return diags

}
func resourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(vast_client.JwtSession)

	attrs := map[string]interface{}{"path": utils.GenPath("users"), "id": d.Id()}
	response, err := utils.DefaultGetFunc(ctx, client, attrs, d, map[string]string{})
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
	resource := api_latest.User{}
	body, err := utils.DefaultProcessingFunc(ctx, response)

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
	diags = ResourceUserReadStructIntoSchema(ctx, resource, d)

	return diags
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(vast_client.JwtSession)
	attrs := map[string]interface{}{"path": utils.GenPath("users"), "id": d.Id()}

	response, err := utils.DefaultDeleteFunc(ctx, client, attrs, nil, map[string]string{})

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

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	names_mapping := utils.ContextKey("names_mapping")
	new_ctx := context.WithValue(ctx, names_mapping, User_names_mapping)
	var diags diag.Diagnostics
	data := make(map[string]interface{})
	client := m.(vast_client.JwtSession)
	tflog.Info(ctx, fmt.Sprintf("Creating Resource User"))
	reflect_User := reflect.TypeOf((*api_latest.User)(nil))
	utils.PopulateResourceMap(new_ctx, reflect_User.Elem(), d, &data, "", false)

	version_compare := utils.VastVersionsWarn(ctx)

	if version_compare != metadata.CLUSTER_VERSION_EQUALS {
		cluster_version := metadata.ClusterVersionString()
		t, t_exists := vast_versions.GetVersionedType(cluster_version, "User")
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
			tflog.Warn(ctx, fmt.Sprintf("Could have not found resource %s in version %s , things might not work properly", "User", cluster_version))
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
	attrs := map[string]interface{}{"path": utils.GenPath("users")}
	response, create_err := utils.DefaultCreateFunc(ctx, client, attrs, data, map[string]string{})
	tflog.Info(ctx, fmt.Sprintf("Server Error for  User %v", create_err))

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
	resource := api_latest.User{}
	err = json.Unmarshal(response_body, &resource)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to convert response body into User",
			Detail:   err.Error(),
		})
		return diags
	}

	id_err := utils.DefaultIdFunc(ctx, client, resource.Id, d)
	if id_err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to set Id",
			Detail:   err.Error(),
		})
		return diags
	}
	ctx_with_resource := context.WithValue(ctx, utils.ContextKey("resource"), resource)
	resourceUserRead(ctx_with_resource, d, m)

	return diags
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	names_mapping := utils.ContextKey("names_mapping")
	new_ctx := context.WithValue(ctx, names_mapping, User_names_mapping)
	var diags diag.Diagnostics
	data := make(map[string]interface{})
	version_compare := utils.VastVersionsWarn(ctx)
	if version_compare != metadata.CLUSTER_VERSION_EQUALS {
		cluster_version := metadata.ClusterVersionString()
		t, t_exists := vast_versions.GetVersionedType(cluster_version, "User")
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
			tflog.Warn(ctx, fmt.Sprintf("Could have not found resource %s in version %s , things might not work properly", "User", cluster_version))
		}
	}

	client := m.(vast_client.JwtSession)
	tflog.Info(ctx, fmt.Sprintf("Updating Resource User"))
	reflect_User := reflect.TypeOf((*api_latest.User)(nil))
	utils.PopulateResourceMap(new_ctx, reflect_User.Elem(), d, &data, "", false)

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
	attrs := map[string]interface{}{"path": "users", "id": d.Id()}
	response, patch_err := utils.DefaultUpdateFunc(ctx, client, attrs, data, d, map[string]string{})
	tflog.Info(ctx, fmt.Sprintf("Server Error for  User %v", patch_err))
	if patch_err != nil {
		error_message := patch_err.Error() + " Server Response: " + utils.GetResponseBodyAsStr(response)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Object Creation Failed",
			Detail:   error_message,
		})
		return diags
	}
	resourceUserRead(ctx, d, m)

	return diags

}

func resourceUserImporter(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {

	result := []*schema.ResourceData{}
	client := m.(vast_client.JwtSession)
	resource_config := codegen_configs.GetResourceByName("User")
	attrs := map[string]interface{}{"path": utils.GenPath("users")}
	response, err := utils.DefaultImportFunc(ctx, client, attrs, d, resource_config.Importer.GetFunc())

	if err != nil {
		return result, err
	}

	resource_l := []api_latest.User{}

	body, err := utils.DefaultProcessingFunc(ctx, response)
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
	id_err := utils.DefaultIdFunc(ctx, client, resource.Id, d)
	if id_err != nil {
		return result, id_err
	}

	diags := ResourceUserReadStructIntoSchema(ctx, resource, d)
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
