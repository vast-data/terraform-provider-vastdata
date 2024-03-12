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

func ResourceS3LifeCycleRule() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceS3LifeCycleRuleRead,
		DeleteContext: resourceS3LifeCycleRuleDelete,
		CreateContext: resourceS3LifeCycleRuleCreate,
		UpdateContext: resourceS3LifeCycleRuleUpdate,

		Importer: &schema.ResourceImporter{
			StateContext: resourceS3LifeCycleRuleImporter,
		},

		Description: ``,
		Schema:      getResourceS3LifeCycleRuleSchema(),
	}
}

func getResourceS3LifeCycleRuleSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{

		"name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},

		"guid": &schema.Schema{
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    false,
			Sensitive:   false,
			Description: ``,
		},

		"enabled": &schema.Schema{
			Type:        schema.TypeBool,
			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: ``,
		},

		"prefix": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},

		"min_size": &schema.Schema{
			Type:        schema.TypeInt,
			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `The minimum size of the object`,
		},

		"max_size": &schema.Schema{
			Type:        schema.TypeInt,
			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `The maximum size of the object`,
		},

		"expiration_days": &schema.Schema{
			Type:        schema.TypeInt,
			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `The number of days from creation until an object expires`,
		},

		"expiration_date": &schema.Schema{
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `The expiration date of the object`,
		},

		"expired_obj_delete_marker": &schema.Schema{
			Type:        schema.TypeBool,
			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `Remove expired objects delete markers`,
		},

		"noncurrent_days": &schema.Schema{
			Type:        schema.TypeInt,
			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `Number of days after objects become noncurrent`,
		},

		"newer_noncurrent_versions": &schema.Schema{
			Type:        schema.TypeInt,
			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `The number of newer versions to retain`,
		},

		"abort_mpu_days_after_initiation": &schema.Schema{
			Type:        schema.TypeInt,
			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `The number of days until expiration after an incomplete multipart upload`,
		},

		"view_path": &schema.Schema{
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `The path of the related View`,
		},

		"view_id": &schema.Schema{
			Type:        schema.TypeInt,
			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `The ID of the related View`,
		},
	}
}

var S3LifeCycleRule_names_mapping map[string][]string = map[string][]string{}

func ResourceS3LifeCycleRuleReadStructIntoSchema(ctx context.Context, resource api_latest.S3LifeCycleRule, d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics
	var err error

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Name", resource.Name))

	err = d.Set("name", resource.Name)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"name\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Guid", resource.Guid))

	err = d.Set("guid", resource.Guid)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"guid\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Enabled", resource.Enabled))

	err = d.Set("enabled", resource.Enabled)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"enabled\"",
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

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "MinSize", resource.MinSize))

	err = d.Set("min_size", resource.MinSize)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"min_size\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "MaxSize", resource.MaxSize))

	err = d.Set("max_size", resource.MaxSize)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"max_size\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "ExpirationDays", resource.ExpirationDays))

	err = d.Set("expiration_days", resource.ExpirationDays)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"expiration_days\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "ExpirationDate", resource.ExpirationDate))

	err = d.Set("expiration_date", resource.ExpirationDate)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"expiration_date\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "ExpiredObjDeleteMarker", resource.ExpiredObjDeleteMarker))

	err = d.Set("expired_obj_delete_marker", resource.ExpiredObjDeleteMarker)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"expired_obj_delete_marker\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "NoncurrentDays", resource.NoncurrentDays))

	err = d.Set("noncurrent_days", resource.NoncurrentDays)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"noncurrent_days\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "NewerNoncurrentVersions", resource.NewerNoncurrentVersions))

	err = d.Set("newer_noncurrent_versions", resource.NewerNoncurrentVersions)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"newer_noncurrent_versions\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "AbortMpuDaysAfterInitiation", resource.AbortMpuDaysAfterInitiation))

	err = d.Set("abort_mpu_days_after_initiation", resource.AbortMpuDaysAfterInitiation)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"abort_mpu_days_after_initiation\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "ViewPath", resource.ViewPath))

	err = d.Set("view_path", resource.ViewPath)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"view_path\"",
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
func resourceS3LifeCycleRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(vast_client.JwtSession)

	attrs := map[string]interface{}{"path": utils.GenPath("s3lifecyclerules"), "id": d.Id()}
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
	resource := api_latest.S3LifeCycleRule{}
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
	diags = ResourceS3LifeCycleRuleReadStructIntoSchema(ctx, resource, d)

	return diags
}

func resourceS3LifeCycleRuleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(vast_client.JwtSession)
	attrs := map[string]interface{}{"path": utils.GenPath("s3lifecyclerules"), "id": d.Id()}

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

func resourceS3LifeCycleRuleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	names_mapping := utils.ContextKey("names_mapping")
	new_ctx := context.WithValue(ctx, names_mapping, S3LifeCycleRule_names_mapping)
	var diags diag.Diagnostics
	data := make(map[string]interface{})
	client := m.(vast_client.JwtSession)
	tflog.Info(ctx, fmt.Sprintf("Creating Resource S3LifeCycleRule"))
	reflect_S3LifeCycleRule := reflect.TypeOf((*api_latest.S3LifeCycleRule)(nil))
	utils.PopulateResourceMap(new_ctx, reflect_S3LifeCycleRule.Elem(), d, &data, "", false)

	var before_post_error error
	data, before_post_error = utils.EnabledMustBeSet(data, client, ctx, d)
	if before_post_error != nil {
		return diag.FromErr(before_post_error)
	}

	version_compare := utils.VastVersionsWarn(ctx)

	if version_compare != metadata.CLUSTER_VERSION_EQUALS {
		cluster_version := metadata.ClusterVersionString()
		t, t_exists := vast_versions.GetVersionedType(cluster_version, "S3LifeCycleRule")
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
			tflog.Warn(ctx, fmt.Sprintf("Could have not found resource %s in version %s , things might not work properly", "S3LifeCycleRule", cluster_version))
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
	attrs := map[string]interface{}{"path": utils.GenPath("s3lifecyclerules")}
	response, create_err := utils.DefaultCreateFunc(ctx, client, attrs, data, map[string]string{})
	tflog.Info(ctx, fmt.Sprintf("Server Error for  S3LifeCycleRule %v", create_err))

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
	resource := api_latest.S3LifeCycleRule{}
	err = json.Unmarshal(response_body, &resource)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to convert response body into S3LifeCycleRule",
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
	resourceS3LifeCycleRuleRead(ctx_with_resource, d, m)

	return diags
}

func resourceS3LifeCycleRuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	names_mapping := utils.ContextKey("names_mapping")
	new_ctx := context.WithValue(ctx, names_mapping, S3LifeCycleRule_names_mapping)
	var diags diag.Diagnostics
	data := make(map[string]interface{})
	version_compare := utils.VastVersionsWarn(ctx)
	if version_compare != metadata.CLUSTER_VERSION_EQUALS {
		cluster_version := metadata.ClusterVersionString()
		t, t_exists := vast_versions.GetVersionedType(cluster_version, "S3LifeCycleRule")
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
			tflog.Warn(ctx, fmt.Sprintf("Could have not found resource %s in version %s , things might not work properly", "S3LifeCycleRule", cluster_version))
		}
	}

	client := m.(vast_client.JwtSession)
	tflog.Info(ctx, fmt.Sprintf("Updating Resource S3LifeCycleRule"))
	reflect_S3LifeCycleRule := reflect.TypeOf((*api_latest.S3LifeCycleRule)(nil))
	utils.PopulateResourceMap(new_ctx, reflect_S3LifeCycleRule.Elem(), d, &data, "", false)

	var before_patch_error error
	data, before_patch_error = utils.EnabledMustBeSet(data, client, ctx, d)
	if before_patch_error != nil {
		return diag.FromErr(before_patch_error)
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
	attrs := map[string]interface{}{"path": utils.GenPath("s3lifecyclerules"), "id": d.Id()}
	response, patch_err := utils.DefaultUpdateFunc(ctx, client, attrs, data, d, map[string]string{})
	tflog.Info(ctx, fmt.Sprintf("Server Error for  S3LifeCycleRule %v", patch_err))
	if patch_err != nil {
		error_message := patch_err.Error() + " Server Response: " + utils.GetResponseBodyAsStr(response)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Object Creation Failed",
			Detail:   error_message,
		})
		return diags
	}
	resourceS3LifeCycleRuleRead(ctx, d, m)

	return diags

}

func resourceS3LifeCycleRuleImporter(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {

	result := []*schema.ResourceData{}
	client := m.(vast_client.JwtSession)
	guid := d.Id()
	values := url.Values{}
	values.Add("guid", fmt.Sprintf("%v", guid))
	attrs := map[string]interface{}{"path": utils.GenPath("s3lifecyclerules"), "query": values.Encode()}
	response, err := utils.DefaultGetFunc(ctx, client, attrs, d, map[string]string{})

	if err != nil {
		return result, err
	}

	resource_l := []api_latest.S3LifeCycleRule{}

	body, err := utils.ProcessingResultsListResponse(ctx, response)
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

	diags := ResourceS3LifeCycleRuleReadStructIntoSchema(ctx, resource, d)
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
