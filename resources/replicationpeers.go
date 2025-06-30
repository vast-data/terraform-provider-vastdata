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

func ResourceReplicationPeers() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceReplicationPeersRead,
		DeleteContext: resourceReplicationPeersDelete,
		CreateContext: resourceReplicationPeersCreate,
		UpdateContext: resourceReplicationPeersUpdate,

		Importer: &schema.ResourceImporter{
			StateContext: resourceReplicationPeersImporter,
		},

		Description: ``,
		Schema:      getResourceReplicationPeersSchema(),
	}
}

func getResourceReplicationPeersSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{

		"guid": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("ReplicationPeers").GetConflictingFields("guid"),

			Computed:    true,
			Optional:    false,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) A unique guid given to the  replication peer configuration`,
		},

		"name": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("ReplicationPeers").GetConflictingFields("name"),

			Required:    true,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The name of the replication peer configuration`,
		},

		"url": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("ReplicationPeers").GetConflictingFields("url"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Direct url of the replication peer configurations`,
		},

		"leading_vip": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("ReplicationPeers").GetConflictingFields("leading_vip"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The vip provided for the replication peer configuration`,
		},

		"remote_vip_range": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("ReplicationPeers").GetConflictingFields("remote_vip_range"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The vip range which were reported by the peer`,
		},

		"version": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("ReplicationPeers").GetConflictingFields("version"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The version of the source`,
		},

		"remote_version": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("ReplicationPeers").GetConflictingFields("remote_version"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The version of the remote peer`,
		},

		"is_local": {
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("ReplicationPeers").GetConflictingFields("is_local"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Is the source of the replication local (this host is the source)`,
		},

		"peer_name": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("ReplicationPeers").GetConflictingFields("peer_name"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The name of the peer cluster`,
		},

		"secure_mode": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("ReplicationPeers").GetConflictingFields("secure_mode"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Is the connection secure`,
		},

		"pool_id": {
			Type:          schema.TypeInt,
			ConflictsWith: codegen_configs.GetResourceByName("ReplicationPeers").GetConflictingFields("pool_id"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The replication Vippool id`,
		},
	}
}

var ReplicationPeersNamesMapping = map[string][]string{}

func ResourceReplicationPeersReadStructIntoSchema(ctx context.Context, resource api_latest.ReplicationPeers, d *schema.ResourceData) diag.Diagnostics {
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

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "LeadingVip", resource.LeadingVip))

	err = d.Set("leading_vip", resource.LeadingVip)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"leading_vip\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "RemoteVipRange", resource.RemoteVipRange))

	err = d.Set("remote_vip_range", resource.RemoteVipRange)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"remote_vip_range\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Version", resource.Version))

	err = d.Set("version", resource.Version)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"version\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "RemoteVersion", resource.RemoteVersion))

	err = d.Set("remote_version", resource.RemoteVersion)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"remote_version\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "IsLocal", resource.IsLocal))

	err = d.Set("is_local", resource.IsLocal)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"is_local\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "PeerName", resource.PeerName))

	err = d.Set("peer_name", resource.PeerName)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"peer_name\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "SecureMode", resource.SecureMode))

	err = d.Set("secure_mode", resource.SecureMode)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"secure_mode\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "PoolId", resource.PoolId))

	err = d.Set("pool_id", resource.PoolId)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"pool_id\"",
			Detail:   err.Error(),
		})
	}

	return diags

}
func resourceReplicationPeersRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*vast_client.VMSSession)
	resourceConfig := codegen_configs.GetResourceByName("ReplicationPeers")
	attrs := map[string]interface{}{"path": utils.GenPath("nativereplicationremotetargets"), "id": d.Id()}
	tflog.Debug(ctx, fmt.Sprintf("[resourceReplicationPeersRead] Calling Get Function : %v for resource ReplicationPeers", utils.GetFuncName(resourceConfig.GetFunc)))
	response, err := resourceConfig.GetFunc(ctx, client, attrs, d, map[string]string{})
	utils.VastVersionsWarn(ctx)

	var body []byte
	var resource api_latest.ReplicationPeers
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
	diags = ResourceReplicationPeersReadStructIntoSchema(ctx, resource, d)

	return diags
}

func resourceReplicationPeersDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*vast_client.VMSSession)
	resourceConfig := codegen_configs.GetResourceByName("ReplicationPeers")
	attrs := map[string]interface{}{"path": utils.GenPath("nativereplicationremotetargets"), "id": d.Id()}

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

func resourceReplicationPeersCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	namesMapping := utils.ContextKey("names_mapping")
	newCtx := context.WithValue(ctx, namesMapping, ReplicationPeersNamesMapping)
	var diags diag.Diagnostics
	data := make(map[string]interface{})
	client := m.(*vast_client.VMSSession)
	resourceConfig := codegen_configs.GetResourceByName("ReplicationPeers")
	tflog.Info(ctx, fmt.Sprintf("Creating Resource ReplicationPeers"))
	reflectReplicationPeers := reflect.TypeOf((*api_latest.ReplicationPeers)(nil))
	utils.PopulateResourceMap(newCtx, reflectReplicationPeers.Elem(), d, &data, "", false)

	versionsEqual := utils.VastVersionsWarn(ctx)

	if versionsEqual != metadata.CLUSTER_VERSION_EQUALS {
		clusterVersion := metadata.ClusterVersionString()
		t, typeExists := vast_versions.GetVersionedType(clusterVersion, "ReplicationPeers")
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
			tflog.Warn(ctx, fmt.Sprintf("Could have not found resource %s in version %s, things might not work properly", "ReplicationPeers", clusterVersion))
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
	attrs := map[string]interface{}{"path": utils.GenPath("nativereplicationremotetargets")}
	response, createErr := resourceConfig.CreateFunc(ctx, client, attrs, data, map[string]string{})
	tflog.Info(ctx, fmt.Sprintf("Server Error for  ReplicationPeers %v", createErr))

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
	resource := api_latest.ReplicationPeers{}
	err = json.Unmarshal(responseBody, &resource)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to convert response body into ReplicationPeers",
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
	resourceReplicationPeersRead(ctxWithResource, d, m)

	return diags
}

func resourceReplicationPeersUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	namesMapping := utils.ContextKey("names_mapping")
	newCtx := context.WithValue(ctx, namesMapping, ReplicationPeersNamesMapping)
	var diags diag.Diagnostics
	data := make(map[string]interface{})
	versionsEqual := utils.VastVersionsWarn(ctx)
	resourceConfig := codegen_configs.GetResourceByName("ReplicationPeers")
	if versionsEqual != metadata.CLUSTER_VERSION_EQUALS {
		clusterVersion := metadata.ClusterVersionString()
		t, typeExists := vast_versions.GetVersionedType(clusterVersion, "ReplicationPeers")
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
			tflog.Warn(ctx, fmt.Sprintf("Could have not found resource %s in version %s, things might not work properly", "ReplicationPeers", clusterVersion))
		}
	}

	client := m.(*vast_client.VMSSession)
	tflog.Info(ctx, fmt.Sprintf("Updating Resource ReplicationPeers"))
	reflectReplicationPeers := reflect.TypeOf((*api_latest.ReplicationPeers)(nil))
	utils.PopulateResourceMap(newCtx, reflectReplicationPeers.Elem(), d, &data, "", false)

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
	attrs := map[string]interface{}{"path": utils.GenPath("nativereplicationremotetargets"), "id": d.Id()}
	response, patchErr := resourceConfig.UpdateFunc(ctx, client, attrs, data, d, map[string]string{})
	tflog.Info(ctx, fmt.Sprintf("Server Error for  ReplicationPeers %v", patchErr))
	if patchErr != nil {
		errorMessage := fmt.Sprintf("server response:\n%v\nUnderlying error:\n%v", utils.GetResponseBodyAsStr(response), patchErr.Error())
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Object Creation Failed",
			Detail:   errorMessage,
		})
		return diags
	}
	resourceReplicationPeersRead(ctx, d, m)

	return diags

}

func resourceReplicationPeersImporter(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {

	var result []*schema.ResourceData
	client := m.(*vast_client.VMSSession)
	resourceConfig := codegen_configs.GetResourceByName("ReplicationPeers")
	attrs := map[string]interface{}{"path": utils.GenPath("nativereplicationremotetargets")}
	response, err := resourceConfig.ImportFunc(ctx, client, attrs, d, resourceConfig.Importer.GetFunc())

	if err != nil {
		return result, err
	}

	var resourceList []api_latest.ReplicationPeers
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

	diags := ResourceReplicationPeersReadStructIntoSchema(ctx, resource, d)
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
