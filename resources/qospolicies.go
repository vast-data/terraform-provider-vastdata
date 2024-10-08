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

func ResourceQosPolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceQosPolicyRead,
		DeleteContext: resourceQosPolicyDelete,
		CreateContext: resourceQosPolicyCreate,
		UpdateContext: resourceQosPolicyUpdate,

		Importer: &schema.ResourceImporter{
			StateContext: resourceQosPolicyImporter,
		},

		Description: ``,
		Schema:      getResourceQosPolicySchema(),
	}
}

func getResourceQosPolicySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{

		"guid": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("QosPolicy").GetConflictingFields("guid"),

			Computed:    true,
			Optional:    false,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0) QoS Policy guid`,
		},

		"name": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("QosPolicy").GetConflictingFields("name"),

			Required: true,
		},

		"mode": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("QosPolicy").GetConflictingFields("mode"),

			Computed:  true,
			Optional:  true,
			Sensitive: false,

			ValidateDiagFunc: utils.OneOf([]string{"STATIC", "USED_CAPACITY", "PROVISIONED_CAPACITY"}),
			Description:      `(Valid for versions: 5.0.0,5.1.0) QoS provisioning mode Allowed Values are [STATIC USED_CAPACITY PROVISIONED_CAPACITY]`,
		},

		"io_size_bytes": &schema.Schema{
			Type:          schema.TypeInt,
			ConflictsWith: codegen_configs.GetResourceByName("QosPolicy").GetConflictingFields("io_size_bytes"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0) Sets the size of IO for static and capacity limit definitions. The number of IOs per request is obtained by dividing request size by IO size. Default: 64K, Recommended range: 4K - 1M`,
		},

		"static_limits": &schema.Schema{
			Type:          schema.TypeList,
			ConflictsWith: codegen_configs.GetResourceByName("QosPolicy").GetConflictingFields("static_limits"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0) `,

			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{

					"min_reads_bw_mbps": &schema.Schema{
						Type:          schema.TypeInt,
						ConflictsWith: codegen_configs.GetResourceByName("QosStaticLimits").GetConflictingFields("min_reads_bw_mbps"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.0.0,5.1.0) Minimal amount of performance to provide when there is resource contention`,
					},

					"max_reads_bw_mbps": &schema.Schema{
						Type:          schema.TypeInt,
						ConflictsWith: codegen_configs.GetResourceByName("QosStaticLimits").GetConflictingFields("max_reads_bw_mbps"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.0.0,5.1.0) Maximal amount of performance to provide when there is no resource contention`,
					},

					"min_writes_bw_mbps": &schema.Schema{
						Type:          schema.TypeInt,
						ConflictsWith: codegen_configs.GetResourceByName("QosStaticLimits").GetConflictingFields("min_writes_bw_mbps"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.0.0,5.1.0) Minimal amount of performance to provide when there is resource contention`,
					},

					"max_writes_bw_mbps": &schema.Schema{
						Type:          schema.TypeInt,
						ConflictsWith: codegen_configs.GetResourceByName("QosStaticLimits").GetConflictingFields("max_writes_bw_mbps"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.0.0,5.1.0) Maximal amount of performance to provide when there is no resource contention`,
					},

					"min_reads_iops": &schema.Schema{
						Type:          schema.TypeInt,
						ConflictsWith: codegen_configs.GetResourceByName("QosStaticLimits").GetConflictingFields("min_reads_iops"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.0.0,5.1.0) Minimal amount of performance to provide when there is resource contention`,
					},

					"max_reads_iops": &schema.Schema{
						Type:          schema.TypeInt,
						ConflictsWith: codegen_configs.GetResourceByName("QosStaticLimits").GetConflictingFields("max_reads_iops"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.0.0,5.1.0) Maximal amount of performance to provide when there is no resource contention`,
					},

					"min_writes_iops": &schema.Schema{
						Type:          schema.TypeInt,
						ConflictsWith: codegen_configs.GetResourceByName("QosStaticLimits").GetConflictingFields("min_writes_iops"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.0.0,5.1.0) Minimal amount of performance to provide when there is resource contention`,
					},

					"max_writes_iops": &schema.Schema{
						Type:          schema.TypeInt,
						ConflictsWith: codegen_configs.GetResourceByName("QosStaticLimits").GetConflictingFields("max_writes_iops"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.0.0,5.1.0) Maximal amount of performance to provide when there is no resource contention`,
					},
				},
			},
		},

		"capacity_limits": &schema.Schema{
			Type:          schema.TypeList,
			ConflictsWith: codegen_configs.GetResourceByName("QosPolicy").GetConflictingFields("capacity_limits"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0) `,

			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{

					"max_reads_bw_mbps_per_gb_capacity": &schema.Schema{
						Type:          schema.TypeInt,
						ConflictsWith: codegen_configs.GetResourceByName("QosDynamicLimits").GetConflictingFields("max_reads_bw_mbps_per_gb_capacity"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.0.0,5.1.0) Maximal amount of performance per GB to provide when there is no resource contention`,
					},

					"max_writes_bw_mbps_per_gb_capacity": &schema.Schema{
						Type:          schema.TypeInt,
						ConflictsWith: codegen_configs.GetResourceByName("QosDynamicLimits").GetConflictingFields("max_writes_bw_mbps_per_gb_capacity"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.0.0,5.1.0) Maximal amount of performance per GB to provide when there is no resource contention`,
					},

					"max_reads_iops_per_gb_capacity": &schema.Schema{
						Type:          schema.TypeInt,
						ConflictsWith: codegen_configs.GetResourceByName("QosDynamicLimits").GetConflictingFields("max_reads_iops_per_gb_capacity"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.0.0,5.1.0) Maximal amount of performance per GB to provide when there is no resource contention`,
					},

					"max_writes_iops_per_gb_capacity": &schema.Schema{
						Type:          schema.TypeInt,
						ConflictsWith: codegen_configs.GetResourceByName("QosDynamicLimits").GetConflictingFields("max_writes_iops_per_gb_capacity"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.0.0,5.1.0) Maximal amount of performance per GB to provide when there is no resource contention`,
					},
				},
			},
		},
	}
}

var QosPolicy_names_mapping map[string][]string = map[string][]string{}

func ResourceQosPolicyReadStructIntoSchema(ctx context.Context, resource api_latest.QosPolicy, d *schema.ResourceData) diag.Diagnostics {
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

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Mode", resource.Mode))

	err = d.Set("mode", resource.Mode)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"mode\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "IoSizeBytes", resource.IoSizeBytes))

	err = d.Set("io_size_bytes", resource.IoSizeBytes)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"io_size_bytes\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "StaticLimits", resource.StaticLimits))

	tflog.Debug(ctx, fmt.Sprintf("Found a pointer object %v", resource.StaticLimits))
	err = d.Set("static_limits", utils.FlattenModelAsList(ctx, resource.StaticLimits))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"static_limits\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "CapacityLimits", resource.CapacityLimits))

	tflog.Debug(ctx, fmt.Sprintf("Found a pointer object %v", resource.CapacityLimits))
	err = d.Set("capacity_limits", utils.FlattenModelAsList(ctx, resource.CapacityLimits))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"capacity_limits\"",
			Detail:   err.Error(),
		})
	}

	return diags

}
func resourceQosPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(vast_client.JwtSession)
	resource_config := codegen_configs.GetResourceByName("QosPolicy")
	attrs := map[string]interface{}{"path": utils.GenPath("qospolicies"), "id": d.Id()}
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
	resource := api_latest.QosPolicy{}
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
	diags = ResourceQosPolicyReadStructIntoSchema(ctx, resource, d)

	return diags
}

func resourceQosPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(vast_client.JwtSession)
	resource_config := codegen_configs.GetResourceByName("QosPolicy")
	attrs := map[string]interface{}{"path": utils.GenPath("qospolicies"), "id": d.Id()}

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

func resourceQosPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	names_mapping := utils.ContextKey("names_mapping")
	new_ctx := context.WithValue(ctx, names_mapping, QosPolicy_names_mapping)
	var diags diag.Diagnostics
	data := make(map[string]interface{})
	client := m.(vast_client.JwtSession)
	resource_config := codegen_configs.GetResourceByName("QosPolicy")
	tflog.Info(ctx, fmt.Sprintf("Creating Resource QosPolicy"))
	reflect_QosPolicy := reflect.TypeOf((*api_latest.QosPolicy)(nil))
	utils.PopulateResourceMap(new_ctx, reflect_QosPolicy.Elem(), d, &data, "", false)

	version_compare := utils.VastVersionsWarn(ctx)

	if version_compare != metadata.CLUSTER_VERSION_EQUALS {
		cluster_version := metadata.ClusterVersionString()
		t, t_exists := vast_versions.GetVersionedType(cluster_version, "QosPolicy")
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
			tflog.Warn(ctx, fmt.Sprintf("Could have not found resource %s in version %s , things might not work properly", "QosPolicy", cluster_version))
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
	attrs := map[string]interface{}{"path": utils.GenPath("qospolicies")}
	response, create_err := resource_config.CreateFunc(ctx, client, attrs, data, map[string]string{})
	tflog.Info(ctx, fmt.Sprintf("Server Error for  QosPolicy %v", create_err))

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
	resource := api_latest.QosPolicy{}
	err = json.Unmarshal(response_body, &resource)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to convert response body into QosPolicy",
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
	resourceQosPolicyRead(ctx_with_resource, d, m)

	return diags
}

func resourceQosPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	names_mapping := utils.ContextKey("names_mapping")
	new_ctx := context.WithValue(ctx, names_mapping, QosPolicy_names_mapping)
	var diags diag.Diagnostics
	data := make(map[string]interface{})
	version_compare := utils.VastVersionsWarn(ctx)
	resource_config := codegen_configs.GetResourceByName("QosPolicy")
	if version_compare != metadata.CLUSTER_VERSION_EQUALS {
		cluster_version := metadata.ClusterVersionString()
		t, t_exists := vast_versions.GetVersionedType(cluster_version, "QosPolicy")
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
			tflog.Warn(ctx, fmt.Sprintf("Could have not found resource %s in version %s , things might not work properly", "QosPolicy", cluster_version))
		}
	}

	client := m.(vast_client.JwtSession)
	tflog.Info(ctx, fmt.Sprintf("Updating Resource QosPolicy"))
	reflect_QosPolicy := reflect.TypeOf((*api_latest.QosPolicy)(nil))
	utils.PopulateResourceMap(new_ctx, reflect_QosPolicy.Elem(), d, &data, "", false)

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
	attrs := map[string]interface{}{"path": utils.GenPath("qospolicies"), "id": d.Id()}
	response, patch_err := resource_config.UpdateFunc(ctx, client, attrs, data, d, map[string]string{})
	tflog.Info(ctx, fmt.Sprintf("Server Error for  QosPolicy %v", patch_err))
	if patch_err != nil {
		error_message := patch_err.Error() + " Server Response: " + utils.GetResponseBodyAsStr(response)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Object Creation Failed",
			Detail:   error_message,
		})
		return diags
	}
	resourceQosPolicyRead(ctx, d, m)

	return diags

}

func resourceQosPolicyImporter(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {

	result := []*schema.ResourceData{}
	client := m.(vast_client.JwtSession)
	resource_config := codegen_configs.GetResourceByName("QosPolicy")
	attrs := map[string]interface{}{"path": utils.GenPath("qospolicies")}
	response, err := resource_config.ImportFunc(ctx, client, attrs, d, resource_config.Importer.GetFunc())

	if err != nil {
		return result, err
	}

	resource_l := []api_latest.QosPolicy{}
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

	diags := ResourceQosPolicyReadStructIntoSchema(ctx, resource, d)
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
