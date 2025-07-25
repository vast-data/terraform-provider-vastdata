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

func ResourceTenant() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceTenantRead,
		DeleteContext: resourceTenantDelete,
		CreateContext: resourceTenantCreate,
		UpdateContext: resourceTenantUpdate,

		Importer: &schema.ResourceImporter{
			StateContext: resourceTenantImporter,
		},

		Description: ``,
		Schema:      getResourceTenantSchema(),
	}
}

func getResourceTenantSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{

		"guid": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Tenant").GetConflictingFields("guid"),

			Computed:    true,
			Optional:    false,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The unique GUID of the tenant.`,
		},

		"name": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Tenant").GetConflictingFields("name"),

			Required:    true,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The unique name of the tenant.`,
		},

		"use_smb_privileged_user": {
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("Tenant").GetConflictingFields("use_smb_privileged_user"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.1.0,5.2.0) Enables or disables SMB privileged user.`,
		},

		"smb_privileged_user_name": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Tenant").GetConflictingFields("smb_privileged_user_name"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) An optional custom username for the SMB privileged user. If not set, the SMB privileged user name is 'vastadmin'.`,
		},

		"use_smb_privileged_group": {
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("Tenant").GetConflictingFields("use_smb_privileged_group"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.1.0,5.2.0) Enables or disables SMB privileged user group.`,
		},

		"smb_privileged_group_sid": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Tenant").GetConflictingFields("smb_privileged_group_sid"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) An optional custom SID to specify a non-default SMB privileged group. If not set, the SMB privileged group is the Backup Operators domain group.`,
		},

		"smb_privileged_group_full_access": {
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("Tenant").GetConflictingFields("smb_privileged_group_full_access"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.1.0,5.2.0) If 'true', the SMB privileged user group has read and write control access. Members of the group can perform backup and restore operations on all files and directories, without requiring read or write access to the specific files and directories. If 'false', the privileged group has read-only access.`,
		},

		"smb_administrators_group_name": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Tenant").GetConflictingFields("smb_administrators_group_name"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) An optional custom name to specify a non-default privileged group. If not set, the privileged group is the Backup Operators domain group.`,
		},

		"default_others_share_level_perm": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Tenant").GetConflictingFields("default_others_share_level_perm"),

			Computed:  true,
			Optional:  true,
			Sensitive: false,

			ValidateDiagFunc: utils.OneOf([]string{"READ", "CHANGE", "FULL"}),
			Description:      `(Valid for versions: 5.0.0,5.1.0,5.2.0) Default share-level permissions for others. Allowed Values are [READ CHANGE FULL]`,
		},

		"trash_gid": {
			Type:          schema.TypeInt,
			ConflictsWith: codegen_configs.GetResourceByName("Tenant").GetConflictingFields("trash_gid"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) A GID with permissions to the trash folder.`,
		},

		"client_ip_ranges": {
			Type:          schema.TypeList,
			ConflictsWith: codegen_configs.GetResourceByName("Tenant").GetConflictingFields("client_ip_ranges"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) An array of source IP ranges to allow for the tenant.`,

			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{

					"start_ip": &schema.Schema{
						Type:        schema.TypeString,
						Computed:    true,
						Optional:    true,
						Description: "The first ip of the range",
					},

					"end_ip": &schema.Schema{
						Type:        schema.TypeString,
						Computed:    true,
						Optional:    true,
						Description: "The last ip of the range",
					},
				},
			},
		},

		"posix_primary_provider": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Tenant").GetConflictingFields("posix_primary_provider"),

			Computed:  true,
			Optional:  true,
			Sensitive: false,

			ValidateDiagFunc: utils.OneOf([]string{"NONE", "LDAP", "NIS", "AD", "LOCAL"}),
			Description:      `(Valid for versions: 5.0.0,5.1.0,5.2.0) The POSIX primary provider type. Allowed Values are [NONE LDAP NIS AD LOCAL]`,
		},

		"ad_provider_id": {
			Type:          schema.TypeInt,
			ConflictsWith: codegen_configs.GetResourceByName("Tenant").GetConflictingFields("ad_provider_id"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The ID of the Active Directory provider.`,
		},

		"ldap_provider_id": {
			Type:          schema.TypeInt,
			ConflictsWith: codegen_configs.GetResourceByName("Tenant").GetConflictingFields("ldap_provider_id"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The ID of the OpenLDAP provider specified separately by the user.`,
		},

		"nis_provider_id": {
			Type:          schema.TypeInt,
			ConflictsWith: codegen_configs.GetResourceByName("Tenant").GetConflictingFields("nis_provider_id"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The NIS provider ID.`,
		},

		"encryption_crn": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Tenant").GetConflictingFields("encryption_crn"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The unique ID of the tenant's encryption group.`,
		},

		"is_nfsv42_supported": {
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("Tenant").GetConflictingFields("is_nfsv42_supported"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.1.0,5.2.0) Enables or disables NFSv4.2.`,
		},

		"allow_locked_users": {
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("Tenant").GetConflictingFields("allow_locked_users"),

			Computed:    false,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.1.0,5.2.0) Allows or prohibits IO from users whose Active Directory accounts are locked out by lockout policies due to unsuccessful login attempts.`,

			Default: false,
		},

		"allow_disabled_users": {
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("Tenant").GetConflictingFields("allow_disabled_users"),

			Computed:    false,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.1.0,5.2.0) Allows or prohibits IO from users whose Active Directory accounts are explicitly disabled.`,

			Default: false,
		},

		"use_smb_native": {
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("Tenant").GetConflictingFields("use_smb_native"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.1.0,5.2.0) Enables or disables use of native SMB authentication.`,
		},

		"vippool_names": {
			Type:          schema.TypeList,
			ConflictsWith: codegen_configs.GetResourceByName("Tenant").GetConflictingFields("vippool_names"),

			Computed:    true,
			Optional:    false,
			Sensitive:   false,
			Description: `(Valid for versions: 5.1.0,5.2.0) An array of names of virtual IP pools attached to the tenant.`,

			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},

		"vippool_ids": {
			Type:          schema.TypeList,
			ConflictsWith: codegen_configs.GetResourceByName("Tenant").GetConflictingFields("vippool_ids"),

			DiffSuppressOnRefresh: false,
			DiffSuppressFunc:      codegen_configs.GetResourceByName("Tenant").GetAttributeDiffFunc("vippool_ids"),
			Computed:              true,
			Optional:              true,
			Sensitive:             false,
			Description:           `(Valid for versions: 5.1.0,5.2.0) An array of IDs of virtual IP pools attached to the tenant.`,

			Elem: &schema.Schema{
				Type: schema.TypeInt,
			},
		},
	}
}

var TenantNamesMapping = map[string][]string{
	"client_ip_ranges": []string{"start_ip", "end_ip"},
}

func ResourceTenantReadStructIntoSchema(ctx context.Context, resource api_latest.Tenant, d *schema.ResourceData) diag.Diagnostics {
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

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "UseSmbPrivilegedUser", resource.UseSmbPrivilegedUser))

	err = d.Set("use_smb_privileged_user", resource.UseSmbPrivilegedUser)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"use_smb_privileged_user\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "SmbPrivilegedUserName", resource.SmbPrivilegedUserName))

	err = d.Set("smb_privileged_user_name", resource.SmbPrivilegedUserName)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"smb_privileged_user_name\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "UseSmbPrivilegedGroup", resource.UseSmbPrivilegedGroup))

	err = d.Set("use_smb_privileged_group", resource.UseSmbPrivilegedGroup)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"use_smb_privileged_group\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "SmbPrivilegedGroupSid", resource.SmbPrivilegedGroupSid))

	err = d.Set("smb_privileged_group_sid", resource.SmbPrivilegedGroupSid)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"smb_privileged_group_sid\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "SmbPrivilegedGroupFullAccess", resource.SmbPrivilegedGroupFullAccess))

	err = d.Set("smb_privileged_group_full_access", resource.SmbPrivilegedGroupFullAccess)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"smb_privileged_group_full_access\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "SmbAdministratorsGroupName", resource.SmbAdministratorsGroupName))

	err = d.Set("smb_administrators_group_name", resource.SmbAdministratorsGroupName)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"smb_administrators_group_name\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "DefaultOthersShareLevelPerm", resource.DefaultOthersShareLevelPerm))

	err = d.Set("default_others_share_level_perm", resource.DefaultOthersShareLevelPerm)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"default_others_share_level_perm\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "TrashGid", resource.TrashGid))

	err = d.Set("trash_gid", resource.TrashGid)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"trash_gid\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "ClientIpRanges", resource.ClientIpRanges))

	err = d.Set("client_ip_ranges", utils.FlattenListOfStringsList(&resource.ClientIpRanges, []string{"start_ip", "end_ip"}))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"client_ip_ranges\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "PosixPrimaryProvider", resource.PosixPrimaryProvider))

	err = d.Set("posix_primary_provider", resource.PosixPrimaryProvider)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"posix_primary_provider\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "AdProviderId", resource.AdProviderId))

	err = d.Set("ad_provider_id", resource.AdProviderId)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"ad_provider_id\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "LdapProviderId", resource.LdapProviderId))

	err = d.Set("ldap_provider_id", resource.LdapProviderId)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"ldap_provider_id\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "NisProviderId", resource.NisProviderId))

	err = d.Set("nis_provider_id", resource.NisProviderId)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"nis_provider_id\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "EncryptionCrn", resource.EncryptionCrn))

	err = d.Set("encryption_crn", resource.EncryptionCrn)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"encryption_crn\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "IsNfsv42Supported", resource.IsNfsv42Supported))

	err = d.Set("is_nfsv42_supported", resource.IsNfsv42Supported)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"is_nfsv42_supported\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "AllowLockedUsers", resource.AllowLockedUsers))

	err = d.Set("allow_locked_users", resource.AllowLockedUsers)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"allow_locked_users\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "AllowDisabledUsers", resource.AllowDisabledUsers))

	err = d.Set("allow_disabled_users", resource.AllowDisabledUsers)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"allow_disabled_users\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "UseSmbNative", resource.UseSmbNative))

	err = d.Set("use_smb_native", resource.UseSmbNative)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"use_smb_native\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "VippoolNames", resource.VippoolNames))

	err = d.Set("vippool_names", utils.FlattenListOfPrimitives(&resource.VippoolNames))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"vippool_names\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "VippoolIds", resource.VippoolIds))

	err = d.Set("vippool_ids", utils.FlattenListOfPrimitives(&resource.VippoolIds))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"vippool_ids\"",
			Detail:   err.Error(),
		})
	}

	return diags

}
func resourceTenantRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*vast_client.VMSSession)
	resourceConfig := codegen_configs.GetResourceByName("Tenant")
	attrs := map[string]interface{}{"path": utils.GenPath("tenants"), "id": d.Id()}
	tflog.Debug(ctx, fmt.Sprintf("[resourceTenantRead] Calling Get Function : %v for resource Tenant", utils.GetFuncName(resourceConfig.GetFunc)))
	response, err := resourceConfig.GetFunc(ctx, client, attrs, d, map[string]string{})
	utils.VastVersionsWarn(ctx)

	var body []byte
	var resource api_latest.Tenant
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
	diags = ResourceTenantReadStructIntoSchema(ctx, resource, d)

	var after_read_error error
	after_read_error = resourceConfig.AfterReadFunc(client, ctx, d)
	if after_read_error != nil {
		return diag.FromErr(after_read_error)
	}

	return diags
}

func resourceTenantDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*vast_client.VMSSession)
	resourceConfig := codegen_configs.GetResourceByName("Tenant")
	attrs := map[string]interface{}{"path": utils.GenPath("tenants"), "id": d.Id()}

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

func resourceTenantCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	namesMapping := utils.ContextKey("names_mapping")
	newCtx := context.WithValue(ctx, namesMapping, TenantNamesMapping)
	var diags diag.Diagnostics
	data := make(map[string]interface{})
	client := m.(*vast_client.VMSSession)
	resourceConfig := codegen_configs.GetResourceByName("Tenant")
	tflog.Info(ctx, fmt.Sprintf("Creating Resource Tenant"))
	reflectTenant := reflect.TypeOf((*api_latest.Tenant)(nil))
	utils.PopulateResourceMap(newCtx, reflectTenant.Elem(), d, &data, "", false)

	var before_post_error error
	data, before_post_error = resourceConfig.BeforePostFunc(data, client, ctx, d)
	if before_post_error != nil {
		return diag.FromErr(before_post_error)
	}

	versionsEqual := utils.VastVersionsWarn(ctx)

	if versionsEqual != metadata.CLUSTER_VERSION_EQUALS {
		clusterVersion := metadata.ClusterVersionString()
		t, typeExists := vast_versions.GetVersionedType(clusterVersion, "Tenant")
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
			tflog.Warn(ctx, fmt.Sprintf("Could have not found resource %s in version %s, things might not work properly", "Tenant", clusterVersion))
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
	attrs := map[string]interface{}{"path": utils.GenPath("tenants")}
	response, createErr := resourceConfig.CreateFunc(ctx, client, attrs, data, map[string]string{})
	tflog.Info(ctx, fmt.Sprintf("Server Error for  Tenant %v", createErr))

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
	resource := api_latest.Tenant{}
	err = json.Unmarshal(responseBody, &resource)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to convert response body into Tenant",
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
	resourceTenantRead(ctxWithResource, d, m)

	return diags
}

func resourceTenantUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	namesMapping := utils.ContextKey("names_mapping")
	newCtx := context.WithValue(ctx, namesMapping, TenantNamesMapping)
	var diags diag.Diagnostics
	data := make(map[string]interface{})
	versionsEqual := utils.VastVersionsWarn(ctx)
	resourceConfig := codegen_configs.GetResourceByName("Tenant")
	if versionsEqual != metadata.CLUSTER_VERSION_EQUALS {
		clusterVersion := metadata.ClusterVersionString()
		t, typeExists := vast_versions.GetVersionedType(clusterVersion, "Tenant")
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
			tflog.Warn(ctx, fmt.Sprintf("Could have not found resource %s in version %s, things might not work properly", "Tenant", clusterVersion))
		}
	}

	client := m.(*vast_client.VMSSession)
	tflog.Info(ctx, fmt.Sprintf("Updating Resource Tenant"))
	reflectTenant := reflect.TypeOf((*api_latest.Tenant)(nil))
	utils.PopulateResourceMap(newCtx, reflectTenant.Elem(), d, &data, "", false)

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
	attrs := map[string]interface{}{"path": utils.GenPath("tenants"), "id": d.Id()}
	response, patchErr := resourceConfig.UpdateFunc(ctx, client, attrs, data, d, map[string]string{})
	tflog.Info(ctx, fmt.Sprintf("Server Error for  Tenant %v", patchErr))
	if patchErr != nil {
		errorMessage := fmt.Sprintf("server response:\n%v\nUnderlying error:\n%v", utils.GetResponseBodyAsStr(response), patchErr.Error())
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Object Creation Failed",
			Detail:   errorMessage,
		})
		return diags
	}
	resourceTenantRead(ctx, d, m)

	return diags

}

func resourceTenantImporter(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {

	var result []*schema.ResourceData
	client := m.(*vast_client.VMSSession)
	resourceConfig := codegen_configs.GetResourceByName("Tenant")
	attrs := map[string]interface{}{"path": utils.GenPath("tenants")}
	response, err := resourceConfig.ImportFunc(ctx, client, attrs, d, resourceConfig.Importer.GetFunc())

	if err != nil {
		return result, err
	}

	var resourceList []api_latest.Tenant
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

	diags := ResourceTenantReadStructIntoSchema(ctx, resource, d)
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
