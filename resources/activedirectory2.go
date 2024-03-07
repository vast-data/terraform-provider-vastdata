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

func ResourceActiveDirectory2() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceActiveDirectory2Read,
		DeleteContext: resourceActiveDirectory2Delete,
		CreateContext: resourceActiveDirectory2Create,
		UpdateContext: resourceActiveDirectory2Update,

		Importer: &schema.ResourceImporter{
			StateContext: resourceActiveDirectory2Importer,
		},

		Description: ``,
		Schema:      getResourceActiveDirectory2Schema(),
	}
}

func getResourceActiveDirectory2Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{

		"guid": &schema.Schema{
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    false,
			Sensitive:   false,
			Description: ``,
		},

		"machine_account_name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},

		"organizational_unit": &schema.Schema{
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: ``,
			ForceNew:    true,
		},

		"smb_allowed": &schema.Schema{
			Type:        schema.TypeBool,
			Computed:    false,
			Optional:    true,
			Sensitive:   false,
			Description: ``,

			Default: true,
		},

		"ntlm_enabled": &schema.Schema{
			Type:        schema.TypeBool,
			Computed:    false,
			Optional:    true,
			Sensitive:   false,
			Description: ``,

			Default: true,
		},

		"use_auto_discovery": &schema.Schema{
			Type:        schema.TypeBool,
			Computed:    false,
			Optional:    true,
			Sensitive:   false,
			Description: ``,

			Default: false,
		},

		"use_ldaps": &schema.Schema{
			Type:        schema.TypeBool,
			Computed:    false,
			Optional:    true,
			Sensitive:   false,
			Description: ``,

			Default: false,
		},

		"port": &schema.Schema{
			Type:        schema.TypeInt,
			Computed:    false,
			Optional:    true,
			Sensitive:   false,
			Description: ``,

			Default: 389,
		},

		"binddn": &schema.Schema{
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: ``,
		},

		"searchbase": &schema.Schema{
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: ``,
		},

		"domain_name": &schema.Schema{
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: ``,
		},

		"method": &schema.Schema{
			Type:      schema.TypeString,
			Computed:  false,
			Optional:  true,
			Sensitive: false,

			ValidateDiagFunc: utils.OneOf([]string{"simple", "anonymous", "sasl"}),
			Description:      ` Allowed Values are [simple anonymous sasl]`,

			Default: "simple",
		},

		"query_groups_mode": &schema.Schema{
			Type:      schema.TypeString,
			Computed:  false,
			Optional:  true,
			Sensitive: false,

			ValidateDiagFunc: utils.OneOf([]string{"COMPATIBLE", "NONE", "RFC2307BIS", "RFC2307"}),
			Description:      ` Allowed Values are [COMPATIBLE NONE RFC2307BIS RFC2307]`,

			Default: "COMPATIBLE",
		},

		"posix_attributes_source": &schema.Schema{
			Type:        schema.TypeString,
			Computed:    false,
			Optional:    true,
			Sensitive:   false,
			Description: ``,

			Default: "JOINED_DOMAIN",
		},

		"use_tls": &schema.Schema{
			Type:        schema.TypeBool,
			Computed:    false,
			Optional:    true,
			Sensitive:   false,
			Description: ``,

			Default: false,
		},

		"reverse_lookup": &schema.Schema{
			Type:        schema.TypeBool,
			Computed:    false,
			Optional:    true,
			Sensitive:   false,
			Description: ``,

			Default: false,
		},

		"gid_number": &schema.Schema{
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: ``,
		},

		"use_multi_forest": &schema.Schema{
			Type:        schema.TypeBool,
			Computed:    false,
			Optional:    true,
			Sensitive:   false,
			Description: ``,

			Default: false,
		},

		"uid": &schema.Schema{
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: ``,
		},

		"uid_number": &schema.Schema{
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: ``,
		},

		"match_user": &schema.Schema{
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: ``,
		},

		"uid_member_value_property_name": &schema.Schema{
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: ``,
		},

		"uid_member": &schema.Schema{
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: ``,
		},

		"posix_account": &schema.Schema{
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: ``,
		},

		"posix_group": &schema.Schema{
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: ``,
		},

		"username_property_name": &schema.Schema{
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: ``,
		},

		"user_login_name": &schema.Schema{
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: ``,
		},

		"group_login_name": &schema.Schema{
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: ``,
		},

		"mail_property_name": &schema.Schema{
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: ``,
		},

		"is_vms_auth_provider": &schema.Schema{
			Type:        schema.TypeBool,
			Computed:    false,
			Optional:    true,
			Sensitive:   false,
			Description: ``,

			Default: false,
		},

		"bindpw": &schema.Schema{
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: ``,
		},

		"urls": &schema.Schema{
			Type:        schema.TypeList,
			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `List of LDAP servers urls , starting with ldap://`,

			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},

		"skip_ldap": &schema.Schema{
			Type:        schema.TypeString,
			Computed:    false,
			Optional:    true,
			Sensitive:   false,
			Description: `When creating an active directory using this resource an ldap configuration is also created if set to true than whed deleting this resource the ldap configuration will be kept otherwise it will also be deleted`,

			Default: false,
		},
	}
}

var ActiveDirectory2_names_mapping map[string][]string = map[string][]string{}

func ResourceActiveDirectory2ReadStructIntoSchema(ctx context.Context, resource api_latest.ActiveDirectory2, d *schema.ResourceData) diag.Diagnostics {
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

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "MachineAccountName", resource.MachineAccountName))

	err = d.Set("machine_account_name", resource.MachineAccountName)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"machine_account_name\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "OrganizationalUnit", resource.OrganizationalUnit))

	err = d.Set("organizational_unit", resource.OrganizationalUnit)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"organizational_unit\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "SmbAllowed", resource.SmbAllowed))

	err = d.Set("smb_allowed", resource.SmbAllowed)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"smb_allowed\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "NtlmEnabled", resource.NtlmEnabled))

	err = d.Set("ntlm_enabled", resource.NtlmEnabled)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"ntlm_enabled\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "UseAutoDiscovery", resource.UseAutoDiscovery))

	err = d.Set("use_auto_discovery", resource.UseAutoDiscovery)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"use_auto_discovery\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "UseLdaps", resource.UseLdaps))

	err = d.Set("use_ldaps", resource.UseLdaps)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"use_ldaps\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Port", resource.Port))

	err = d.Set("port", resource.Port)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"port\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Binddn", resource.Binddn))

	err = d.Set("binddn", resource.Binddn)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"binddn\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Searchbase", resource.Searchbase))

	err = d.Set("searchbase", resource.Searchbase)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"searchbase\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "DomainName", resource.DomainName))

	err = d.Set("domain_name", resource.DomainName)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"domain_name\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Method", resource.Method))

	err = d.Set("method", resource.Method)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"method\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "QueryGroupsMode", resource.QueryGroupsMode))

	err = d.Set("query_groups_mode", resource.QueryGroupsMode)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"query_groups_mode\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "PosixAttributesSource", resource.PosixAttributesSource))

	err = d.Set("posix_attributes_source", resource.PosixAttributesSource)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"posix_attributes_source\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "UseTls", resource.UseTls))

	err = d.Set("use_tls", resource.UseTls)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"use_tls\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "ReverseLookup", resource.ReverseLookup))

	err = d.Set("reverse_lookup", resource.ReverseLookup)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"reverse_lookup\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "GidNumber", resource.GidNumber))

	err = d.Set("gid_number", resource.GidNumber)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"gid_number\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "UseMultiForest", resource.UseMultiForest))

	err = d.Set("use_multi_forest", resource.UseMultiForest)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"use_multi_forest\"",
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

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "UidNumber", resource.UidNumber))

	err = d.Set("uid_number", resource.UidNumber)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"uid_number\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "MatchUser", resource.MatchUser))

	err = d.Set("match_user", resource.MatchUser)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"match_user\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "UidMemberValuePropertyName", resource.UidMemberValuePropertyName))

	err = d.Set("uid_member_value_property_name", resource.UidMemberValuePropertyName)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"uid_member_value_property_name\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "UidMember", resource.UidMember))

	err = d.Set("uid_member", resource.UidMember)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"uid_member\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "PosixAccount", resource.PosixAccount))

	err = d.Set("posix_account", resource.PosixAccount)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"posix_account\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "PosixGroup", resource.PosixGroup))

	err = d.Set("posix_group", resource.PosixGroup)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"posix_group\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "UsernamePropertyName", resource.UsernamePropertyName))

	err = d.Set("username_property_name", resource.UsernamePropertyName)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"username_property_name\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "UserLoginName", resource.UserLoginName))

	err = d.Set("user_login_name", resource.UserLoginName)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"user_login_name\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "GroupLoginName", resource.GroupLoginName))

	err = d.Set("group_login_name", resource.GroupLoginName)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"group_login_name\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "MailPropertyName", resource.MailPropertyName))

	err = d.Set("mail_property_name", resource.MailPropertyName)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"mail_property_name\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "IsVmsAuthProvider", resource.IsVmsAuthProvider))

	err = d.Set("is_vms_auth_provider", resource.IsVmsAuthProvider)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"is_vms_auth_provider\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Bindpw", resource.Bindpw))

	err = d.Set("bindpw", resource.Bindpw)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"bindpw\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Urls", resource.Urls))

	err = d.Set("urls", utils.FlattenListOfPrimitives(&resource.Urls))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"urls\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "SkipLdap", resource.SkipLdap))

	err = d.Set("skip_ldap", resource.SkipLdap)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"skip_ldap\"",
			Detail:   err.Error(),
		})
	}

	return diags

}
func resourceActiveDirectory2Read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(vast_client.JwtSession)

	attrs := map[string]interface{}{"path": "/api/activedirectory/", "id": d.Id()}
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
	resource := api_latest.ActiveDirectory2{}
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
	diags = ResourceActiveDirectory2ReadStructIntoSchema(ctx, resource, d)

	return diags
}

func resourceActiveDirectory2Delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(vast_client.JwtSession)
	attrs := map[string]interface{}{"path": "/api/activedirectory/", "id": d.Id()}

	data, before_delete_error := utils.SkipDeleteAsUserSelected(ctx, d, m)
	if before_delete_error != nil {
		return diag.FromErr(before_delete_error)
	}
	unmarshaled_data := map[string]interface{}{}
	_data, _ := io.ReadAll(data)
	json.Unmarshal(_data, &unmarshaled_data)
	response, err := utils.DefaultDeleteFunc(ctx, client, attrs, unmarshaled_data, map[string]string{})

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

func resourceActiveDirectory2Create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	names_mapping := utils.ContextKey("names_mapping")
	new_ctx := context.WithValue(ctx, names_mapping, ActiveDirectory2_names_mapping)
	var diags diag.Diagnostics
	data := make(map[string]interface{})
	client := m.(vast_client.JwtSession)
	tflog.Info(ctx, fmt.Sprintf("Creating Resource ActiveDirectory2"))
	reflect_ActiveDirectory2 := reflect.TypeOf((*api_latest.ActiveDirectory2)(nil))
	utils.PopulateResourceMap(new_ctx, reflect_ActiveDirectory2.Elem(), d, &data, "", false)

	version_compare := utils.VastVersionsWarn(ctx)

	if version_compare != metadata.CLUSTER_VERSION_EQUALS {
		cluster_version := metadata.ClusterVersionString()
		t, t_exists := vast_versions.GetVersionedType(cluster_version, "ActiveDirectory2")
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
			tflog.Warn(ctx, fmt.Sprintf("Could have not found resource %s in version %s , things might not work properly", "ActiveDirectory2", cluster_version))
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
	attrs := map[string]interface{}{"path": "/api/activedirectory/"}
	response, create_err := utils.DefaultCreateFunc(ctx, client, attrs, data, map[string]string{})
	tflog.Info(ctx, fmt.Sprintf("Server Error for  ActiveDirectory2 %v", create_err))

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
	resource := api_latest.ActiveDirectory2{}
	err = json.Unmarshal(response_body, &resource)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to convert response body into ActiveDirectory2",
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
	resourceActiveDirectory2Read(ctx_with_resource, d, m)

	return diags
}

func resourceActiveDirectory2Update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	names_mapping := utils.ContextKey("names_mapping")
	new_ctx := context.WithValue(ctx, names_mapping, ActiveDirectory2_names_mapping)
	var diags diag.Diagnostics
	data := make(map[string]interface{})
	version_compare := utils.VastVersionsWarn(ctx)
	if version_compare != metadata.CLUSTER_VERSION_EQUALS {
		cluster_version := metadata.ClusterVersionString()
		t, t_exists := vast_versions.GetVersionedType(cluster_version, "ActiveDirectory2")
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
			tflog.Warn(ctx, fmt.Sprintf("Could have not found resource %s in version %s , things might not work properly", "ActiveDirectory2", cluster_version))
		}
	}

	client := m.(vast_client.JwtSession)
	tflog.Info(ctx, fmt.Sprintf("Updating Resource ActiveDirectory2"))
	reflect_ActiveDirectory2 := reflect.TypeOf((*api_latest.ActiveDirectory2)(nil))
	utils.PopulateResourceMap(new_ctx, reflect_ActiveDirectory2.Elem(), d, &data, "", false)

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
	attrs := map[string]interface{}{"path": "/api/activedirectory/", "id": d.Id()}
	response, patch_err := utils.DefaultUpdateFunc(ctx, client, attrs, data, d, map[string]string{})
	tflog.Info(ctx, fmt.Sprintf("Server Error for  ActiveDirectory2 %v", patch_err))
	if patch_err != nil {
		error_message := patch_err.Error() + " Server Response: " + utils.GetResponseBodyAsStr(response)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Object Creation Failed",
			Detail:   error_message,
		})
		return diags
	}
	resourceActiveDirectory2Read(ctx, d, m)

	return diags

}

func resourceActiveDirectory2Importer(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {

	result := []*schema.ResourceData{}
	client := m.(vast_client.JwtSession)
	guid := d.Id()
	values := url.Values{}
	values.Add("guid", fmt.Sprintf("%v", guid))
	attrs := map[string]interface{}{"path": "/api/activedirectory/", "query": values.Encode()}
	response, err := utils.DefaultGetFunc(ctx, client, attrs, d, map[string]string{})

	if err != nil {
		return result, err
	}

	resource_l := []api_latest.ActiveDirectory2{}

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

	diags := ResourceActiveDirectory2ReadStructIntoSchema(ctx, resource, d)
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
