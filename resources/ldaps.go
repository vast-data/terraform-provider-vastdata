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

func ResourceLdap() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceLdapRead,
		DeleteContext: resourceLdapDelete,
		CreateContext: resourceLdapCreate,
		UpdateContext: resourceLdapUpdate,

		Importer: &schema.ResourceImporter{
			StateContext: resourceLdapImporter,
		},

		Description: ``,
		Schema:      getResourceLdapSchema(),
	}
}

func getResourceLdapSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{

		"guid": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Ldap").GetConflictingFields("guid"),

			Computed:    true,
			Optional:    false,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,
		},

		"urls": {
			Type:          schema.TypeList,
			ConflictsWith: codegen_configs.GetResourceByName("Ldap").GetConflictingFields("urls"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) List of URIs of LDAP servers (Domain Controllers (DCs) in Active Directory), in priority order. The URI with highest priority that has a good health status is used. Specify each URI in the format <scheme>://<address>. <address> can be either a DNS name or an IP address. e.g. ldap://ldap.company.com, ldaps://ldaps.company.com, ldap://192.0.2.2`,

			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},

		"port": {
			Type:          schema.TypeInt,
			ConflictsWith: codegen_configs.GetResourceByName("Ldap").GetConflictingFields("port"),

			Computed:    false,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) LDAP server port. 389 (LDAP)  636 (LDAPS)`,

			Default: 389,
		},

		"binddn": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Ldap").GetConflictingFields("binddn"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Distinguished name of LDAP superuser`,
		},

		"bindpw": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Ldap").GetConflictingFields("bindpw"),

			DiffSuppressOnRefresh: false,
			DiffSuppressFunc:      utils.DoNothingOnUpdate(),

			Computed:    true,
			Optional:    true,
			Sensitive:   true,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Password for the LDAP superuser`,
		},

		"searchbase": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Ldap").GetConflictingFields("searchbase"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The Base DN is the starting point the LDAP provider uses when searching for users and groups. If the Group Base DN is configured it will be used instead of the Base DN, for groups only`,
		},

		"group_searchbase": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Ldap").GetConflictingFields("group_searchbase"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Base DN for group queries within the joined domain only. When auto discovery is enabled, group queries outside the joined domain use auto-discovered Base DNs.`,
		},

		"method": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Ldap").GetConflictingFields("method"),

			Computed:  true,
			Optional:  true,
			Sensitive: false,

			ValidateDiagFunc: utils.OneOf([]string{"simple", "sasl", "anonymous"}),
			Description:      `(Valid for versions: 5.0.0,5.1.0,5.2.0) Bind Authentication Method Allowed Values are [simple sasl anonymous]`,
		},

		"gid_number": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Ldap").GetConflictingFields("gid_number"),

			Computed:    false,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Attrirbute mapping for gid number`,

			Default: "gidNumber",
		},

		"uid": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Ldap").GetConflictingFields("uid"),

			Computed:    false,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Attrirbute mapping for uid`,

			Default: "uid",
		},

		"uid_number": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Ldap").GetConflictingFields("uid_number"),

			Computed:    false,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Attrirbute mapping for uid number`,

			Default: "uidNumber",
		},

		"match_user": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Ldap").GetConflictingFields("match_user"),

			Computed:    false,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Attribute mapping for user matching`,

			Default: "uid",
		},

		"uid_member": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Ldap").GetConflictingFields("uid_member"),

			Computed:    false,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Attrirbute mapping for uid member`,

			Default: "memberUID",
		},

		"posix_account": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Ldap").GetConflictingFields("posix_account"),

			Computed:    false,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Attrirbute mapping for posix account`,

			Default: "posixAccount",
		},

		"posix_group": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Ldap").GetConflictingFields("posix_group"),

			Computed:    false,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Attrirbute mapping for posix account`,

			Default: "posixGroup",
		},

		"use_tls": {
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("Ldap").GetConflictingFields("use_tls"),

			Computed:    false,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) configure LDAP with TLS`,

			Default: false,
		},

		"posix_primary_provider": {
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("Ldap").GetConflictingFields("posix_primary_provider"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) POSIX primary provider`,
		},

		"posix_attributes_source": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Ldap").GetConflictingFields("posix_attributes_source"),

			Computed:    false,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,

			Default: "JOINED_DOMAIN",
		},

		"reverse_lookup": {
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("Ldap").GetConflictingFields("reverse_lookup"),

			Computed:    false,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,

			Default: false,
		},

		"tls_certificate": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Ldap").GetConflictingFields("tls_certificate"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,
		},

		"active_directory": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Ldap").GetConflictingFields("active_directory"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,
		},

		"query_groups_mode": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Ldap").GetConflictingFields("query_groups_mode"),

			Computed:  true,
			Optional:  true,
			Sensitive: false,

			ValidateDiagFunc: utils.OneOf([]string{"COMPATIBLE", "RFC2307BIS", "RFC2307", "NONE"}),
			Description:      `(Valid for versions: 5.0.0,5.1.0,5.2.0) Query group mode Allowed Values are [COMPATIBLE RFC2307BIS RFC2307 NONE]`,
		},

		"username_property_name": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Ldap").GetConflictingFields("username_property_name"),

			Computed:    false,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Username property name`,

			Default: "cn",
		},

		"domain_name": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Ldap").GetConflictingFields("domain_name"),

			Required:    true,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) FQDN of the domain.`,
		},

		"user_login_name": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Ldap").GetConflictingFields("user_login_name"),

			Computed:    false,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The attribute used to query AD for the user login name in NFS ID mapping. Applicable only with AD and NFSv4.1.`,

			Default: "uid",
		},

		"group_login_name": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Ldap").GetConflictingFields("group_login_name"),

			Computed:    false,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The attribute used to query AD for the group login name in NFS ID mapping. Applicable only with AD and NFSv4.1.`,

			Default: "cn",
		},

		"mail_property_name": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Ldap").GetConflictingFields("mail_property_name"),

			Computed:    false,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,

			Default: "mail",
		},

		"uid_member_value_property_name": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Ldap").GetConflictingFields("uid_member_value_property_name"),

			Computed:    false,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,

			Default: "uid",
		},

		"use_auto_discovery": {
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("Ldap").GetConflictingFields("use_auto_discovery"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) When enabled, Active Directory Domain Controllers (DCs) and Active Directory domains are auto discovered. Queries extend beyond the joined domain to all domains in the forest. When disabled, queries are restricted to the joined domain and DCs must be provided in the URLs field.`,
		},

		"use_ldaps": {
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("Ldap").GetConflictingFields("use_ldaps"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Use LDAPS for Auto-Discovery`,
		},

		"is_vms_auth_provider": {
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("Ldap").GetConflictingFields("is_vms_auth_provider"),

			Computed:    false,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Whether the LDAP should be used for VMS auth. There is only two LDAPs allowed for VMS auth: one with AD and one w/o.`,

			Default: false,
		},

		"query_posix_attributes_from_gc": {
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("Ldap").GetConflictingFields("query_posix_attributes_from_gc"),

			Computed:  false,
			Optional:  true,
			Sensitive: false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) When set to True - users/groups from non-joined domain POSIX attributes are supported,
when set to False - Posix attributes of users/groups from non-joined domain are not supported.
As a condition Global catalog needs to be configured to support Posix attributes.
`,

			Default: false,
		},
	}
}

var LdapNamesMapping = map[string][]string{}

func ResourceLdapReadStructIntoSchema(ctx context.Context, resource api_latest.Ldap, d *schema.ResourceData) diag.Diagnostics {
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

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Urls", resource.Urls))

	err = d.Set("urls", utils.FlattenListOfPrimitives(&resource.Urls))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"urls\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Port", resource.Port))

	err = d.Set("port", resource.Port)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"port\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Binddn", resource.Binddn))

	err = d.Set("binddn", resource.Binddn)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"binddn\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Bindpw", resource.Bindpw))

	err = d.Set("bindpw", resource.Bindpw)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"bindpw\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Searchbase", resource.Searchbase))

	err = d.Set("searchbase", resource.Searchbase)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"searchbase\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "GroupSearchbase", resource.GroupSearchbase))

	err = d.Set("group_searchbase", resource.GroupSearchbase)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"group_searchbase\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Method", resource.Method))

	err = d.Set("method", resource.Method)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"method\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "GidNumber", resource.GidNumber))

	err = d.Set("gid_number", resource.GidNumber)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"gid_number\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Uid", resource.Uid))

	err = d.Set("uid", resource.Uid)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"uid\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "UidNumber", resource.UidNumber))

	err = d.Set("uid_number", resource.UidNumber)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"uid_number\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "MatchUser", resource.MatchUser))

	err = d.Set("match_user", resource.MatchUser)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"match_user\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "UidMember", resource.UidMember))

	err = d.Set("uid_member", resource.UidMember)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"uid_member\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "PosixAccount", resource.PosixAccount))

	err = d.Set("posix_account", resource.PosixAccount)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"posix_account\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "PosixGroup", resource.PosixGroup))

	err = d.Set("posix_group", resource.PosixGroup)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"posix_group\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "UseTls", resource.UseTls))

	err = d.Set("use_tls", resource.UseTls)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"use_tls\"",
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

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "PosixAttributesSource", resource.PosixAttributesSource))

	err = d.Set("posix_attributes_source", resource.PosixAttributesSource)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"posix_attributes_source\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "ReverseLookup", resource.ReverseLookup))

	err = d.Set("reverse_lookup", resource.ReverseLookup)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"reverse_lookup\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "TlsCertificate", resource.TlsCertificate))

	err = d.Set("tls_certificate", resource.TlsCertificate)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"tls_certificate\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "ActiveDirectory", resource.ActiveDirectory))

	err = d.Set("active_directory", resource.ActiveDirectory)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"active_directory\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "QueryGroupsMode", resource.QueryGroupsMode))

	err = d.Set("query_groups_mode", resource.QueryGroupsMode)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"query_groups_mode\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "UsernamePropertyName", resource.UsernamePropertyName))

	err = d.Set("username_property_name", resource.UsernamePropertyName)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"username_property_name\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "DomainName", resource.DomainName))

	err = d.Set("domain_name", resource.DomainName)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"domain_name\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "UserLoginName", resource.UserLoginName))

	err = d.Set("user_login_name", resource.UserLoginName)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"user_login_name\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "GroupLoginName", resource.GroupLoginName))

	err = d.Set("group_login_name", resource.GroupLoginName)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"group_login_name\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "MailPropertyName", resource.MailPropertyName))

	err = d.Set("mail_property_name", resource.MailPropertyName)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"mail_property_name\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "UidMemberValuePropertyName", resource.UidMemberValuePropertyName))

	err = d.Set("uid_member_value_property_name", resource.UidMemberValuePropertyName)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"uid_member_value_property_name\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "UseAutoDiscovery", resource.UseAutoDiscovery))

	err = d.Set("use_auto_discovery", resource.UseAutoDiscovery)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"use_auto_discovery\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "UseLdaps", resource.UseLdaps))

	err = d.Set("use_ldaps", resource.UseLdaps)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"use_ldaps\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "IsVmsAuthProvider", resource.IsVmsAuthProvider))

	err = d.Set("is_vms_auth_provider", resource.IsVmsAuthProvider)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"is_vms_auth_provider\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "QueryPosixAttributesFromGc", resource.QueryPosixAttributesFromGc))

	err = d.Set("query_posix_attributes_from_gc", resource.QueryPosixAttributesFromGc)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"query_posix_attributes_from_gc\"",
			Detail:   err.Error(),
		})
	}

	return diags

}
func resourceLdapRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*vast_client.VMSSession)
	resourceConfig := codegen_configs.GetResourceByName("Ldap")
	attrs := map[string]interface{}{"path": utils.GenPath("ldaps"), "id": d.Id()}
	tflog.Debug(ctx, fmt.Sprintf("[resourceLdapRead] Calling Get Function : %v for resource Ldap", utils.GetFuncName(resourceConfig.GetFunc)))
	response, err := resourceConfig.GetFunc(ctx, client, attrs, d, map[string]string{})
	utils.VastVersionsWarn(ctx)

	var body []byte
	var resource api_latest.Ldap
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
		body, err = resourceConfig.ResponseProcessingFunc(ctx, response)
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
	diags = ResourceLdapReadStructIntoSchema(ctx, resource, d)

	return diags
}

func resourceLdapDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*vast_client.VMSSession)
	resourceConfig := codegen_configs.GetResourceByName("Ldap")
	attrs := map[string]interface{}{"path": utils.GenPath("ldaps"), "id": d.Id()}

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

func resourceLdapCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	namesMapping := utils.ContextKey("names_mapping")
	newCtx := context.WithValue(ctx, namesMapping, LdapNamesMapping)
	var diags diag.Diagnostics
	data := make(map[string]interface{})
	client := m.(*vast_client.VMSSession)
	resourceConfig := codegen_configs.GetResourceByName("Ldap")
	tflog.Info(ctx, fmt.Sprintf("Creating Resource Ldap"))
	reflectLdap := reflect.TypeOf((*api_latest.Ldap)(nil))
	utils.PopulateResourceMap(newCtx, reflectLdap.Elem(), d, &data, "", false)

	versionsEqual := utils.VastVersionsWarn(ctx)

	if versionsEqual != metadata.CLUSTER_VERSION_EQUALS {
		clusterVersion := metadata.ClusterVersionString()
		t, typeExists := vast_versions.GetVersionedType(clusterVersion, "Ldap")
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
			tflog.Warn(ctx, fmt.Sprintf("Could have not found resource %s in version %s, things might not work properly", "Ldap", clusterVersion))
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
	attrs := map[string]interface{}{"path": utils.GenPath("ldaps")}
	response, createErr := resourceConfig.CreateFunc(ctx, client, attrs, data, map[string]string{})
	tflog.Info(ctx, fmt.Sprintf("Server Error for  Ldap %v", createErr))

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
	resource := api_latest.Ldap{}
	err = json.Unmarshal(responseBody, &resource)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to convert response body into Ldap",
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
	resourceLdapRead(ctxWithResource, d, m)

	return diags
}

func resourceLdapUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	namesMapping := utils.ContextKey("names_mapping")
	newCtx := context.WithValue(ctx, namesMapping, LdapNamesMapping)
	var diags diag.Diagnostics
	data := make(map[string]interface{})
	versionsEqual := utils.VastVersionsWarn(ctx)
	resourceConfig := codegen_configs.GetResourceByName("Ldap")
	if versionsEqual != metadata.CLUSTER_VERSION_EQUALS {
		clusterVersion := metadata.ClusterVersionString()
		t, typeExists := vast_versions.GetVersionedType(clusterVersion, "Ldap")
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
			tflog.Warn(ctx, fmt.Sprintf("Could have not found resource %s in version %s, things might not work properly", "Ldap", clusterVersion))
		}
	}

	client := m.(*vast_client.VMSSession)
	tflog.Info(ctx, fmt.Sprintf("Updating Resource Ldap"))
	reflectLdap := reflect.TypeOf((*api_latest.Ldap)(nil))
	utils.PopulateResourceMap(newCtx, reflectLdap.Elem(), d, &data, "", false)

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
	attrs := map[string]interface{}{"path": utils.GenPath("ldaps"), "id": d.Id()}
	response, patchErr := resourceConfig.UpdateFunc(ctx, client, attrs, data, d, map[string]string{})
	tflog.Info(ctx, fmt.Sprintf("Server Error for  Ldap %v", patchErr))
	if patchErr != nil {
		errorMessage := fmt.Sprintf("server response:\n%v\nUnderlying error:\n%v", utils.GetResponseBodyAsStr(response), patchErr.Error())
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Object Creation Failed",
			Detail:   errorMessage,
		})
		return diags
	}
	resourceLdapRead(ctx, d, m)

	return diags

}

func resourceLdapImporter(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {

	var result []*schema.ResourceData
	client := m.(*vast_client.VMSSession)
	resourceConfig := codegen_configs.GetResourceByName("Ldap")
	attrs := map[string]interface{}{"path": utils.GenPath("ldaps")}
	response, err := resourceConfig.ImportFunc(ctx, client, attrs, d, resourceConfig.Importer.GetFunc())

	if err != nil {
		return result, err
	}

	var resourceList []api_latest.Ldap
	body, err := resourceConfig.ResponseProcessingFunc(ctx, response)

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

	diags := ResourceLdapReadStructIntoSchema(ctx, resource, d)
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
