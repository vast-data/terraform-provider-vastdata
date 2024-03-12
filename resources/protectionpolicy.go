package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"reflect"

	"errors"
	"net/url"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api_latest "github.com/vast-data/terraform-provider-vastdata/codegen/latest"
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

		"guid": &schema.Schema{
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    false,
			Sensitive:   false,
			Description: `A unique guid given to the  replication peer configuration`,
		},

		"name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},

		"url": &schema.Schema{
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `Direct link to the replication policy`,
		},

		"target_name": &schema.Schema{
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `The target peer name`,
		},

		"target_object_id": &schema.Schema{
			Type:        schema.TypeInt,
			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `The id of the target peer`,
		},

		"prefix": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},

		"clone_type": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},

		"frames": &schema.Schema{
			Type:        schema.TypeList,
			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `List of snapshots schedules`,

			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{

					"every": &schema.Schema{
						Type:             schema.TypeString,
						Computed:         true,
						Optional:         true,
						Sensitive:        false,
						ValidateDiagFunc: utils.ProtectionPolicyTimeIntervalValidation,
						Description:      `How often to make a snapshot, format is <integer><time period> , while time period can be D - Days ,W - Weeks ,s - Seconds ,m - Minutes, H - Hours, M - Months, Y - Years , Ex 1D = 1 Day`,
					},

					"start_at": &schema.Schema{
						Type:             schema.TypeString,
						Computed:         true,
						Optional:         true,
						Sensitive:        false,
						ValidateDiagFunc: utils.ProtectionPolicyStartAt,
						Description:      ``,
					},

					"keep_local": &schema.Schema{
						Type:             schema.TypeString,
						Computed:         true,
						Optional:         true,
						Sensitive:        false,
						ValidateDiagFunc: utils.ProtectionPolicyTimeIntervalValidation,
						Description:      ``,
					},

					"keep_remote": &schema.Schema{
						Type:             schema.TypeString,
						Computed:         true,
						Optional:         true,
						Sensitive:        false,
						ValidateDiagFunc: utils.ProtectionPolicyTimeIntervalValidation,
						Description:      ``,
					},
				},
			},
		},

		"indestructible": &schema.Schema{
			Type:        schema.TypeBool,
			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `Is the snapshot indestructable`,
		},
	}
}

var ProtectionPolicy_names_mapping map[string][]string = map[string][]string{}

func ResourceProtectionPolicyReadStructIntoSchema(ctx context.Context, resource api_latest.ProtectionPolicy, d *schema.ResourceData) diag.Diagnostics {
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

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Url", resource.Url))

	err = d.Set("url", resource.Url)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"url\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "TargetName", resource.TargetName))

	err = d.Set("target_name", resource.TargetName)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"target_name\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "TargetObjectId", resource.TargetObjectId))

	err = d.Set("target_object_id", resource.TargetObjectId)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"target_object_id\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Prefix", resource.Prefix))

	err = d.Set("prefix", resource.Prefix)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"prefix\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "CloneType", resource.CloneType))

	err = d.Set("clone_type", resource.CloneType)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"clone_type\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Frames", resource.Frames))

	err = d.Set("frames", utils.FlattenListOfModelsToList(ctx, resource.Frames))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"frames\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Indestructible", resource.Indestructible))

	err = d.Set("indestructible", resource.Indestructible)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"indestructible\"",
			Detail:   err.Error(),
		})
	}

	return diags

}
func resourceProtectionPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(vast_client.JwtSession)

	attrs := map[string]interface{}{"path": utils.GenPath("protectionpolicies"), "id": d.Id()}
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
	resource := api_latest.ProtectionPolicy{}
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
	diags = ResourceProtectionPolicyReadStructIntoSchema(ctx, resource, d)

	return diags
}

func resourceProtectionPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(vast_client.JwtSession)
	attrs := map[string]interface{}{"path": utils.GenPath("protectionpolicies"), "id": d.Id()}

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

func resourceProtectionPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	names_mapping := utils.ContextKey("names_mapping")
	new_ctx := context.WithValue(ctx, names_mapping, ProtectionPolicy_names_mapping)
	var diags diag.Diagnostics
	data := make(map[string]interface{})
	client := m.(vast_client.JwtSession)
	tflog.Info(ctx, fmt.Sprintf("Creating Resource ProtectionPolicy"))
	reflect_ProtectionPolicy := reflect.TypeOf((*api_latest.ProtectionPolicy)(nil))
	utils.PopulateResourceMap(new_ctx, reflect_ProtectionPolicy.Elem(), d, &data, "", false)

	version_compare := utils.VastVersionsWarn(ctx)

	if version_compare != metadata.CLUSTER_VERSION_EQUALS {
		cluster_version := metadata.ClusterVersionString()
		t, t_exists := vast_versions.GetVersionedType(cluster_version, "ProtectionPolicy")
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
			tflog.Warn(ctx, fmt.Sprintf("Could have not found resource %s in version %s , things might not work properly", "ProtectionPolicy", cluster_version))
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
	response, create_err := utils.DefaultCreateFunc(ctx, client, attrs, data, map[string]string{})
	tflog.Info(ctx, fmt.Sprintf("Server Error for  ProtectionPolicy %v", create_err))

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
	resource := api_latest.ProtectionPolicy{}
	err = json.Unmarshal(response_body, &resource)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to convert response body into ProtectionPolicy",
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
	resourceProtectionPolicyRead(ctx_with_resource, d, m)

	return diags
}

func resourceProtectionPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	names_mapping := utils.ContextKey("names_mapping")
	new_ctx := context.WithValue(ctx, names_mapping, ProtectionPolicy_names_mapping)
	var diags diag.Diagnostics
	data := make(map[string]interface{})
	version_compare := utils.VastVersionsWarn(ctx)
	if version_compare != metadata.CLUSTER_VERSION_EQUALS {
		cluster_version := metadata.ClusterVersionString()
		t, t_exists := vast_versions.GetVersionedType(cluster_version, "ProtectionPolicy")
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
			tflog.Warn(ctx, fmt.Sprintf("Could have not found resource %s in version %s , things might not work properly", "ProtectionPolicy", cluster_version))
		}
	}

	client := m.(vast_client.JwtSession)
	tflog.Info(ctx, fmt.Sprintf("Updating Resource ProtectionPolicy"))
	reflect_ProtectionPolicy := reflect.TypeOf((*api_latest.ProtectionPolicy)(nil))
	utils.PopulateResourceMap(new_ctx, reflect_ProtectionPolicy.Elem(), d, &data, "", false)

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
	response, patch_err := utils.DefaultUpdateFunc(ctx, client, attrs, data, d, map[string]string{})
	tflog.Info(ctx, fmt.Sprintf("Server Error for  ProtectionPolicy %v", patch_err))
	if patch_err != nil {
		error_message := patch_err.Error() + " Server Response: " + utils.GetResponseBodyAsStr(response)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Object Creation Failed",
			Detail:   error_message,
		})
		return diags
	}
	resourceProtectionPolicyRead(ctx, d, m)

	return diags

}

func resourceProtectionPolicyImporter(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {

	result := []*schema.ResourceData{}
	client := m.(vast_client.JwtSession)
	guid := d.Id()
	values := url.Values{}
	values.Add("guid", fmt.Sprintf("%v", guid))
	attrs := map[string]interface{}{"path": utils.GenPath("protectionpolicies"), "query": values.Encode()}
	response, err := utils.DefaultGetFunc(ctx, client, attrs, d, map[string]string{})

	if err != nil {
		return result, err
	}

	resource_l := []api_latest.ProtectionPolicy{}

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

	diags := ResourceProtectionPolicyReadStructIntoSchema(ctx, resource, d)
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
