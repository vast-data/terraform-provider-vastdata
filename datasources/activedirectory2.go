package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api_latest "github.com/vast-data/terraform-provider-vastdata/codegen/latest"
	utils "github.com/vast-data/terraform-provider-vastdata/utils"
	vast_client "github.com/vast-data/terraform-provider-vastdata/vast-client"
	"net/url"
	"strconv"
)

func DataSourceActiveDirectory2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceActiveDirectory2Read,
		Description: ``,
		Schema: map[string]*schema.Schema{

			"guid": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: ``,
			},

			"machine_account_name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    false,
				Required:    true,
				Optional:    false,
				Description: ``,
			},

			"organizational_unit": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: ``,
			},

			"smb_allowed": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    false,
				Required:    false,
				Optional:    true,
				Description: ``,
			},

			"ntlm_enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    false,
				Required:    false,
				Optional:    true,
				Description: ``,
			},

			"use_auto_discovery": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    false,
				Required:    false,
				Optional:    true,
				Description: ``,
			},

			"use_ldaps": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    false,
				Required:    false,
				Optional:    true,
				Description: ``,
			},

			"port": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    false,
				Required:    false,
				Optional:    true,
				Description: ``,
			},

			"binddn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: ``,
			},

			"searchbase": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: ``,
			},

			"domain_name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: ``,
			},

			"method": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    false,
				Required:    false,
				Optional:    true,
				Description: ` Allowed Values are [simple anonymous sasl]`,
			},

			"query_groups_mode": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    false,
				Required:    false,
				Optional:    true,
				Description: ` Allowed Values are [COMPATIBLE NONE RFC2307BIS RFC2307]`,
			},

			"posix_attributes_source": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    false,
				Required:    false,
				Optional:    true,
				Description: ``,
			},

			"use_tls": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    false,
				Required:    false,
				Optional:    true,
				Description: ``,
			},

			"reverse_lookup": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    false,
				Required:    false,
				Optional:    true,
				Description: ``,
			},

			"gid_number": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: ``,
			},

			"use_multi_forest": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    false,
				Required:    false,
				Optional:    true,
				Description: ``,
			},

			"uid": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: ``,
			},

			"uid_number": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: ``,
			},

			"match_user": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: ``,
			},

			"uid_member_value_property_name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: ``,
			},

			"uid_member": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: ``,
			},

			"posix_account": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: ``,
			},

			"posix_group": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: ``,
			},

			"username_property_name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: ``,
			},

			"user_login_name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: ``,
			},

			"group_login_name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: ``,
			},

			"mail_property_name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: ``,
			},

			"is_vms_auth_provider": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    false,
				Required:    false,
				Optional:    true,
				Description: ``,
			},

			"bindpw": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: ``,
			},

			"urls": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `List of LDAP servers urls , starting with ldap://`,

				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"skip_ldap": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    false,
				Required:    false,
				Optional:    true,
				Description: `When creating an active directory using this resource an ldap configuration is also created if set to true than whed deleting this resource the ldap configuration will be kept otherwise it will also be deleted`,
			},
		},
	}
}

func dataSourceActiveDirectory2Read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(vast_client.JwtSession)
	values := url.Values{}

	machine_account_name := d.Get("machine_account_name")
	values.Add("machine_account_name", fmt.Sprintf("%v", machine_account_name))

	response, err := client.Get(ctx, "/api/activedirectory/", values.Encode(), map[string]string{})
	tflog.Info(ctx, response.Request.URL.String())
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured while obtaining data from the vastdata cluster",
			Detail:   err.Error(),
		})
		return diags

	}
	resource_l := []api_latest.ActiveDirectory2{}

	body, err := utils.DefaultProcessingFunc(ctx, response)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured reading data recived from VastData cluster",
			Detail:   err.Error(),
		})
		return diags

	}
	err = json.Unmarshal(body, &resource_l)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured while parsing data recived from VastData cluster",
			Detail:   err.Error(),
		})
		return diags

	}
	if len(resource_l) == 0 {
		d.SetId("")
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Could not find a resource that matches those attributes",
			Detail:   "Could not find a resource that matches those attributes",
		})
		return diags
	}
	if len(resource_l) > 1 {
		d.SetId("")
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Multiple results returned, you might want to add more attributes to get a specific resource",
			Detail:   "Multiple results returned, you might want to add more attributes to get a specific resource",
		})
		return diags
	}

	resource := resource_l[0]

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

	Id := (int64)(resource.Id)
	d.SetId(strconv.FormatInt(Id, 10))
	return diags
}
