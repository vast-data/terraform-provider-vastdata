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

		"guid": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("ActiveDirectory2").GetConflictingFields("guid"),

			Computed:    true,
			Optional:    false,
			Sensitive:   false,
			Description: `(Valid for versions: 5.1.0,5.2.0) The unique GUID of the resource.`,
		},

		"machine_account_name": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("ActiveDirectory2").GetConflictingFields("machine_account_name"),

			Required:    true,
			Description: `(Valid for versions: 5.1.0,5.2.0) The name of the computer object/machine account to add. Recommended to use the name of the cluster.`,
			ForceNew:    true,
		},

		"organizational_unit": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("ActiveDirectory2").GetConflictingFields("organizational_unit"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.1.0,5.2.0) Organizational unit within Active Directory where the cluster's machine account will be created. If left empty, defaults to Computers OU.`,
			ForceNew:    true,
		},

		"smb_allowed": {
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("ActiveDirectory2").GetConflictingFields("smb_allowed"),

			Computed:    false,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.1.0,5.2.0) Indicates if Active Directory is allowed for SMB.`,

			Default: true,
		},

		"ntlm_enabled": {
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("ActiveDirectory2").GetConflictingFields("ntlm_enabled"),

			Computed:    false,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.1.0,5.2.0) Enables or disables support of NTLM authentication for SMB.`,

			Default: true,
		},

		"use_auto_discovery": {
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("ActiveDirectory2").GetConflictingFields("use_auto_discovery"),

			Computed:    false,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.1.0,5.2.0) If 'true', Active Directory Domain Controllers (DCs) and Active Directory domains are automatically discovered. Queries extend beyond the joined domain to all domains in the forest. If 'false', queries are restricted to the joined domain and DCs must be provided in the URLs field.`,

			Default: false,
		},

		"use_ldaps": {
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("ActiveDirectory2").GetConflictingFields("use_ldaps"),

			Computed:    false,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.1.0,5.2.0) Specifies whether to use LDAPS for auto-discovery. To enable use of LDAPS, also set 'use_auto_discovery' to 'true'.`,

			Default: false,
		},

		"port": {
			Type:          schema.TypeInt,
			ConflictsWith: codegen_configs.GetResourceByName("ActiveDirectory2").GetConflictingFields("port"),

			Computed:    false,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.1.0,5.2.0) Which port to use.`,

			Default: 389,
		},

		"binddn": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("ActiveDirectory2").GetConflictingFields("binddn"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.1.0,5.2.0) Distinguished name of the Active Directory superuser.`,
		},

		"searchbase": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("ActiveDirectory2").GetConflictingFields("searchbase"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.1.0,5.2.0) The base DN is the starting point that the Active Directory provider uses when searching for users and groups. If a group base DN is configured, it will be used instead of the base DN (for groups only).`,
		},

		"domain_name": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("ActiveDirectory2").GetConflictingFields("domain_name"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.1.0,5.2.0) FQDN of the domain.`,
		},

		"method": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("ActiveDirectory2").GetConflictingFields("method"),

			Computed:  false,
			Optional:  true,
			Sensitive: false,

			ValidateDiagFunc: utils.OneOf([]string{"simple", "anonymous", "sasl"}),
			Description:      `(Valid for versions: 5.1.0,5.2.0) Bind authentication method. Allowed Values are [simple anonymous sasl]`,

			Default: "simple",
		},

		"query_groups_mode": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("ActiveDirectory2").GetConflictingFields("query_groups_mode"),

			Computed:  false,
			Optional:  true,
			Sensitive: false,

			ValidateDiagFunc: utils.OneOf([]string{"COMPATIBLE", "NONE", "RFC2307BIS", "RFC2307"}),
			Description:      `(Valid for versions: 5.1.0,5.2.0) Query group mode. Allowed Values are [COMPATIBLE NONE RFC2307BIS RFC2307]`,

			Default: "COMPATIBLE",
		},

		"posix_attributes_source": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("ActiveDirectory2").GetConflictingFields("posix_attributes_source"),

			Computed:    false,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.1.0,5.2.0) Defines which domains POSIX attributes will be supported from.`,

			Default: "JOINED_DOMAIN",
		},

		"use_tls": {
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("ActiveDirectory2").GetConflictingFields("use_tls"),

			Computed:    false,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.1.0,5.2.0) Set to 'true' to enable use of TLS to secure communication between the VAST cluster and the Active Directory server.`,

			Default: false,
		},

		"tls_certificate": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("ActiveDirectory2").GetConflictingFields("tls_certificate"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.1.0,5.2.0) TLS certificate to use for verifying the remote Active Directory server’s TLS certificate.`,
		},

		"reverse_lookup": {
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("ActiveDirectory2").GetConflictingFields("reverse_lookup"),

			Computed:    false,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.1.0,5.2.0) Specifies whether to resolve Active Directory netgroups into hostnames.`,

			Default: false,
		},

		"gid_number": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("ActiveDirectory2").GetConflictingFields("gid_number"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.1.0,5.2.0) The attribute of a user entry on the Active Directory server that contains the UID number, if different from 'uidNumber'. Often, when binding the VAST cluster to Active Directory, this does not need to be set.`,
		},

		"use_multi_forest": {
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("ActiveDirectory2").GetConflictingFields("use_multi_forest"),

			Computed:    false,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.1.0,5.2.0) Allows or prohibits access for users from trusted domains on other forests.`,

			Default: false,
		},

		"uid": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("ActiveDirectory2").GetConflictingFields("uid"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.1.0,5.2.0) The attribute of a user entry on the Active Directory server that contains the user name, if different from 'uid'. When binding the VAST cluster to Active Directory, you may need to set this to 'sAMAccountname'.`,
		},

		"uid_number": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("ActiveDirectory2").GetConflictingFields("uid_number"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.1.0,5.2.0) The attribute of a user entry on the Active Directory server that contains the UID number, if different from 'uidNumber'. Often when binding the VAST cluster to Active Directory, this does not need to be set.`,
		},

		"match_user": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("ActiveDirectory2").GetConflictingFields("match_user"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.1.0,5.2.0) The attribute to use when querying a provider for a user that matches a user that has already been retrieved from another provider. A user entry that contains a matching value in this attribute will be considered the same user as the user previously retrieved.`,
		},

		"uid_member_value_property_name": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("ActiveDirectory2").GetConflictingFields("uid_member_value_property_name"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.1.0,5.2.0) Specifies the attribute which represents the value of the Active Directory group’s member property.`,
		},

		"uid_member": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("ActiveDirectory2").GetConflictingFields("uid_member"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.1.0,5.2.0) The attribute of a group entry on the Active Directory server that contains names of group members, if different from 'memberUid'. When binding the VAST cluster to Active Directory, you may need to set this to 'memberUID'.`,
		},

		"posix_account": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("ActiveDirectory2").GetConflictingFields("posix_account"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.1.0,5.2.0) The object class that defines a user entry on the Active Directory server, if different from 'posixAccount'. When binding the VAST cluster to Active Directory, set this parameter to 'user' to ensure that authorization works properly.`,
		},

		"posix_group": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("ActiveDirectory2").GetConflictingFields("posix_group"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.1.0,5.2.0) The object class that defines a group entry on the Active Directory server, if different from 'posixGroup'. When binding the VAST cluster to Active Directory, set this parameter to 'group' to ensure that authorization works properly.`,
		},

		"username_property_name": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("ActiveDirectory2").GetConflictingFields("username_property_name"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.1.0,5.2.0) The attribute to use for querying users in VMS user-initated user queries. Default is 'name'. Sometimes it can be set to 'cn'.`,
		},

		"user_login_name": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("ActiveDirectory2").GetConflictingFields("user_login_name"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.1.0,5.2.0) Specifies the attribute used to query Active Directory for the user login name in NFS ID mapping.`,
		},

		"group_login_name": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("ActiveDirectory2").GetConflictingFields("group_login_name"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.1.0,5.2.0) Specifies the attribute used to query Active Directory for the group login name in NFS ID mapping.`,
		},

		"mail_property_name": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("ActiveDirectory2").GetConflictingFields("mail_property_name"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.1.0,5.2.0) Specifies the attribute to use for the user’s email address.`,
		},

		"is_vms_auth_provider": {
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("ActiveDirectory2").GetConflictingFields("is_vms_auth_provider"),

			Computed:    false,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.1.0,5.2.0) Enables or disables use of the Active Directory for VMS authentication. Two Active Directory configurations per cluster can be used for VMS authentication: one with Active Directory and the other without Active Directory.`,

			Default: false,
		},

		"bindpw": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("ActiveDirectory2").GetConflictingFields("bindpw"),

			DiffSuppressOnRefresh: false,
			DiffSuppressFunc:      utils.DoNothingOnUpdate(),

			Computed:    true,
			Optional:    true,
			Sensitive:   true,
			Description: `(Valid for versions: 5.1.0,5.2.0) The password used with the Bind DN to authenticate to the Active Directory server.`,
		},

		"urls": {
			Type:          schema.TypeList,
			ConflictsWith: codegen_configs.GetResourceByName("ActiveDirectory2").GetConflictingFields("urls"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.1.0,5.2.0) A comma-separated list of URIs of Active Directory servers in the format 'SCHEME://ADDRESS'. The order of listing defines the priority order. The URI with the highest priority that has a good health status is used.`,

			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}

var ActiveDirectory2NamesMapping = map[string][]string{}

func ResourceActiveDirectory2ReadStructIntoSchema(ctx context.Context, resource api_latest.ActiveDirectory2, d *schema.ResourceData) diag.Diagnostics {
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

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "MachineAccountName", resource.MachineAccountName))

	err = d.Set("machine_account_name", resource.MachineAccountName)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"machine_account_name\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "OrganizationalUnit", resource.OrganizationalUnit))

	err = d.Set("organizational_unit", resource.OrganizationalUnit)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"organizational_unit\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "SmbAllowed", resource.SmbAllowed))

	err = d.Set("smb_allowed", resource.SmbAllowed)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"smb_allowed\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "NtlmEnabled", resource.NtlmEnabled))

	err = d.Set("ntlm_enabled", resource.NtlmEnabled)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"ntlm_enabled\"",
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

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Searchbase", resource.Searchbase))

	err = d.Set("searchbase", resource.Searchbase)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"searchbase\"",
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

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Method", resource.Method))

	err = d.Set("method", resource.Method)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"method\"",
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

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "PosixAttributesSource", resource.PosixAttributesSource))

	err = d.Set("posix_attributes_source", resource.PosixAttributesSource)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"posix_attributes_source\"",
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

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "TlsCertificate", resource.TlsCertificate))

	err = d.Set("tls_certificate", resource.TlsCertificate)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"tls_certificate\"",
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

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "GidNumber", resource.GidNumber))

	err = d.Set("gid_number", resource.GidNumber)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"gid_number\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "UseMultiForest", resource.UseMultiForest))

	err = d.Set("use_multi_forest", resource.UseMultiForest)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"use_multi_forest\"",
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

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "UidMemberValuePropertyName", resource.UidMemberValuePropertyName))

	err = d.Set("uid_member_value_property_name", resource.UidMemberValuePropertyName)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"uid_member_value_property_name\"",
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

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "UsernamePropertyName", resource.UsernamePropertyName))

	err = d.Set("username_property_name", resource.UsernamePropertyName)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"username_property_name\"",
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

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "IsVmsAuthProvider", resource.IsVmsAuthProvider))

	err = d.Set("is_vms_auth_provider", resource.IsVmsAuthProvider)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"is_vms_auth_provider\"",
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

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Urls", resource.Urls))

	err = d.Set("urls", utils.FlattenListOfPrimitives(&resource.Urls))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"urls\"",
			Detail:   err.Error(),
		})
	}

	return diags

}
func resourceActiveDirectory2Read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*vast_client.VMSSession)
	resourceConfig := codegen_configs.GetResourceByName("ActiveDirectory2")
	attrs := map[string]interface{}{"path": utils.GenPath("activedirectory"), "id": d.Id()}
	tflog.Debug(ctx, fmt.Sprintf("[resourceActiveDirectory2Read] Calling Get Function : %v for resource ActiveDirectory2", utils.GetFuncName(resourceConfig.GetFunc)))
	response, err := resourceConfig.GetFunc(ctx, client, attrs, d, map[string]string{})
	utils.VastVersionsWarn(ctx)

	var body []byte
	var resource api_latest.ActiveDirectory2
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
	diags = ResourceActiveDirectory2ReadStructIntoSchema(ctx, resource, d)

	return diags
}

func resourceActiveDirectory2Delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*vast_client.VMSSession)
	resourceConfig := codegen_configs.GetResourceByName("ActiveDirectory2")
	attrs := map[string]interface{}{"path": utils.GenPath("activedirectory"), "id": d.Id()}

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

func resourceActiveDirectory2Create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	namesMapping := utils.ContextKey("names_mapping")
	newCtx := context.WithValue(ctx, namesMapping, ActiveDirectory2NamesMapping)
	var diags diag.Diagnostics
	data := make(map[string]interface{})
	client := m.(*vast_client.VMSSession)
	resourceConfig := codegen_configs.GetResourceByName("ActiveDirectory2")
	tflog.Info(ctx, fmt.Sprintf("Creating Resource ActiveDirectory2"))
	reflectActiveDirectory2 := reflect.TypeOf((*api_latest.ActiveDirectory2)(nil))
	utils.PopulateResourceMap(newCtx, reflectActiveDirectory2.Elem(), d, &data, "", false)

	versionsEqual := utils.VastVersionsWarn(ctx)

	if versionsEqual != metadata.CLUSTER_VERSION_EQUALS {
		clusterVersion := metadata.ClusterVersionString()
		t, typeExists := vast_versions.GetVersionedType(clusterVersion, "ActiveDirectory2")
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
			tflog.Warn(ctx, fmt.Sprintf("Could have not found resource %s in version %s, things might not work properly", "ActiveDirectory2", clusterVersion))
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
	attrs := map[string]interface{}{"path": utils.GenPath("activedirectory")}
	response, createErr := resourceConfig.CreateFunc(ctx, client, attrs, data, map[string]string{})
	tflog.Info(ctx, fmt.Sprintf("Server Error for  ActiveDirectory2 %v", createErr))

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
	resource := api_latest.ActiveDirectory2{}
	err = json.Unmarshal(responseBody, &resource)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to convert response body into ActiveDirectory2",
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
	resourceActiveDirectory2Read(ctxWithResource, d, m)

	return diags
}

func resourceActiveDirectory2Update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	namesMapping := utils.ContextKey("names_mapping")
	newCtx := context.WithValue(ctx, namesMapping, ActiveDirectory2NamesMapping)
	var diags diag.Diagnostics
	data := make(map[string]interface{})
	versionsEqual := utils.VastVersionsWarn(ctx)
	resourceConfig := codegen_configs.GetResourceByName("ActiveDirectory2")
	if versionsEqual != metadata.CLUSTER_VERSION_EQUALS {
		clusterVersion := metadata.ClusterVersionString()
		t, typeExists := vast_versions.GetVersionedType(clusterVersion, "ActiveDirectory2")
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
			tflog.Warn(ctx, fmt.Sprintf("Could have not found resource %s in version %s, things might not work properly", "ActiveDirectory2", clusterVersion))
		}
	}

	client := m.(*vast_client.VMSSession)
	tflog.Info(ctx, fmt.Sprintf("Updating Resource ActiveDirectory2"))
	reflectActiveDirectory2 := reflect.TypeOf((*api_latest.ActiveDirectory2)(nil))
	utils.PopulateResourceMap(newCtx, reflectActiveDirectory2.Elem(), d, &data, "", false)

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
	attrs := map[string]interface{}{"path": utils.GenPath("activedirectory"), "id": d.Id()}
	response, patchErr := resourceConfig.UpdateFunc(ctx, client, attrs, data, d, map[string]string{})
	tflog.Info(ctx, fmt.Sprintf("Server Error for  ActiveDirectory2 %v", patchErr))
	if patchErr != nil {
		errorMessage := fmt.Sprintf("server response:\n%v\nUnderlying error:\n%v", utils.GetResponseBodyAsStr(response), patchErr.Error())
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Object Creation Failed",
			Detail:   errorMessage,
		})
		return diags
	}
	resourceActiveDirectory2Read(ctx, d, m)

	return diags

}

func resourceActiveDirectory2Importer(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {

	var result []*schema.ResourceData
	client := m.(*vast_client.VMSSession)
	resourceConfig := codegen_configs.GetResourceByName("ActiveDirectory2")
	attrs := map[string]interface{}{"path": utils.GenPath("activedirectory")}
	response, err := resourceConfig.ImportFunc(ctx, client, attrs, d, resourceConfig.Importer.GetFunc())

	if err != nil {
		return result, err
	}

	var resourceList []api_latest.ActiveDirectory2
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

	diags := ResourceActiveDirectory2ReadStructIntoSchema(ctx, resource, d)
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
