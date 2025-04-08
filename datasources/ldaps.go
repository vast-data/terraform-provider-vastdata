package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api_latest "github.com/vast-data/terraform-provider-vastdata/codegen/latest"
	codegen_configs "github.com/vast-data/terraform-provider-vastdata/codegen_tools/configs"
	utils "github.com/vast-data/terraform-provider-vastdata/utils"
	vast_client "github.com/vast-data/terraform-provider-vastdata/vast-client"
	"net/url"
	"strconv"
)

func DataSourceLdap() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLdapRead,
		Description: ``,
		Schema: map[string]*schema.Schema{

			"id": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,
			},

			"guid": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,
			},

			"urls": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) List of URIs of LDAP servers (Domain Controllers (DCs) in Active Directory), in priority order. The URI with highest priority that has a good health status is used. Specify each URI in the format <scheme>://<address>. <address> can be either a DNS name or an IP address. e.g. ldap://ldap.company.com, ldaps://ldaps.company.com, ldap://192.0.2.2`,

				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"port": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    false,
				Required:    false,
				Optional:    true,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) LDAP server port. 389 (LDAP)  636 (LDAPS)`,
			},

			"binddn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Distinguished name of LDAP superuser`,
			},

			"bindpw": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Password for the LDAP superuser`,
			},

			"searchbase": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The Base DN is the starting point the LDAP provider uses when searching for users and groups. If the Group Base DN is configured it will be used instead of the Base DN, for groups only`,
			},

			"group_searchbase": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Base DN for group queries within the joined domain only. When auto discovery is enabled, group queries outside the joined domain use auto-discovered Base DNs.`,
			},

			"method": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Bind Authentication Method Allowed Values are [simple sasl anonymous]`,
			},

			"gid_number": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    false,
				Required:    false,
				Optional:    true,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Attrirbute mapping for gid number`,
			},

			"uid": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    false,
				Required:    false,
				Optional:    true,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Attrirbute mapping for uid`,
			},

			"uid_number": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    false,
				Required:    false,
				Optional:    true,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Attrirbute mapping for uid number`,
			},

			"match_user": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    false,
				Required:    false,
				Optional:    true,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Attribute mapping for user matching`,
			},

			"uid_member": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    false,
				Required:    false,
				Optional:    true,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Attrirbute mapping for uid member`,
			},

			"posix_account": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    false,
				Required:    false,
				Optional:    true,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Attrirbute mapping for posix account`,
			},

			"posix_group": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    false,
				Required:    false,
				Optional:    true,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Attrirbute mapping for posix account`,
			},

			"use_tls": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    false,
				Required:    false,
				Optional:    true,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) configure LDAP with TLS`,
			},

			"posix_primary_provider": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) POSIX primary provider`,
			},

			"posix_attributes_source": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    false,
				Required:    false,
				Optional:    true,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,
			},

			"reverse_lookup": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    false,
				Required:    false,
				Optional:    true,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,
			},

			"tls_certificate": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,
			},

			"active_directory": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,
			},

			"query_groups_mode": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Query group mode Allowed Values are [COMPATIBLE RFC2307BIS RFC2307 NONE]`,
			},

			"username_property_name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    false,
				Required:    false,
				Optional:    true,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Username property name`,
			},

			"domain_name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    false,
				Required:    true,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) FQDN of the domain.`,
			},

			"user_login_name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    false,
				Required:    false,
				Optional:    true,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The attribute used to query AD for the user login name in NFS ID mapping. Applicable only with AD and NFSv4.1.`,
			},

			"group_login_name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    false,
				Required:    false,
				Optional:    true,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The attribute used to query AD for the group login name in NFS ID mapping. Applicable only with AD and NFSv4.1.`,
			},

			"mail_property_name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    false,
				Required:    false,
				Optional:    true,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,
			},

			"uid_member_value_property_name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    false,
				Required:    false,
				Optional:    true,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,
			},

			"use_auto_discovery": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) When enabled, Active Directory Domain Controllers (DCs) and Active Directory domains are auto discovered. Queries extend beyond the joined domain to all domains in the forest. When disabled, queries are restricted to the joined domain and DCs must be provided in the URLs field.`,
			},

			"use_ldaps": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Use LDAPS for Auto-Discovery`,
			},

			"is_vms_auth_provider": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    false,
				Required:    false,
				Optional:    true,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Whether the LDAP should be used for VMS auth. There is only two LDAPs allowed for VMS auth: one with AD and one w/o.`,
			},

			"query_posix_attributes_from_gc": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: false,
				Required: false,
				Optional: true,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) When set to True - users/groups from non-joined domain POSIX attributes are supported,
when set to False - Posix attributes of users/groups from non-joined domain are not supported.
As a condition Global catalog needs to be configured to support Posix attributes.
`,
			},
		},
	}
}

func dataSourceLdapRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*vast_client.VMSSession)
	values := url.Values{}
	datasource_config := codegen_configs.GetDataSourceByName("Ldap")

	domain_name := d.Get("domain_name")
	values.Add("domain_name", fmt.Sprintf("%v", domain_name))

	response, err := client.Get(ctx, utils.GenPath("ldaps"), values.Encode(), map[string]string{})
	tflog.Info(ctx, response.Request.URL.String())
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured while obtaining data from the vastdata cluster",
			Detail:   err.Error(),
		})
		return diags

	}
	resource_l := []api_latest.Ldap{}
	body, err := datasource_config.ResponseProcessingFunc(ctx, response)

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

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Id", resource.Id))

	err = d.Set("id", resource.Id)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"id\"",
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

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Urls", resource.Urls))

	err = d.Set("urls", utils.FlattenListOfPrimitives(&resource.Urls))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"urls\"",
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

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Bindpw", resource.Bindpw))

	err = d.Set("bindpw", resource.Bindpw)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"bindpw\"",
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

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "GroupSearchbase", resource.GroupSearchbase))

	err = d.Set("group_searchbase", resource.GroupSearchbase)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"group_searchbase\"",
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

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "GidNumber", resource.GidNumber))

	err = d.Set("gid_number", resource.GidNumber)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"gid_number\"",
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

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "UseTls", resource.UseTls))

	err = d.Set("use_tls", resource.UseTls)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"use_tls\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "PosixPrimaryProvider", resource.PosixPrimaryProvider))

	err = d.Set("posix_primary_provider", resource.PosixPrimaryProvider)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"posix_primary_provider\"",
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

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "ReverseLookup", resource.ReverseLookup))

	err = d.Set("reverse_lookup", resource.ReverseLookup)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"reverse_lookup\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "TlsCertificate", resource.TlsCertificate))

	err = d.Set("tls_certificate", resource.TlsCertificate)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"tls_certificate\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "ActiveDirectory", resource.ActiveDirectory))

	err = d.Set("active_directory", resource.ActiveDirectory)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"active_directory\"",
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

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "UsernamePropertyName", resource.UsernamePropertyName))

	err = d.Set("username_property_name", resource.UsernamePropertyName)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"username_property_name\"",
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

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "UidMemberValuePropertyName", resource.UidMemberValuePropertyName))

	err = d.Set("uid_member_value_property_name", resource.UidMemberValuePropertyName)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"uid_member_value_property_name\"",
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

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "IsVmsAuthProvider", resource.IsVmsAuthProvider))

	err = d.Set("is_vms_auth_provider", resource.IsVmsAuthProvider)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"is_vms_auth_provider\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "QueryPosixAttributesFromGc", resource.QueryPosixAttributesFromGc))

	err = d.Set("query_posix_attributes_from_gc", resource.QueryPosixAttributesFromGc)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"query_posix_attributes_from_gc\"",
			Detail:   err.Error(),
		})
	}

	Id := (int64)(resource.Id)
	d.SetId(strconv.FormatInt(Id, 10))
	return diags
}
