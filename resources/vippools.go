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

func ResourceVipPool() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceVipPoolRead,
		DeleteContext: resourceVipPoolDelete,
		CreateContext: resourceVipPoolCreate,
		UpdateContext: resourceVipPoolUpdate,

		Importer: &schema.ResourceImporter{
			StateContext: resourceVipPoolImporter,
		},

		Description: ``,
		Schema:      getResourceVipPoolSchema(),
	}
}

func getResourceVipPoolSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{

		"guid": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("VipPool").GetConflictingFields("guid"),

			Computed:    true,
			Optional:    false,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The unique GUID of the virtual IP pool.`,
		},

		"name": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("VipPool").GetConflictingFields("name"),

			Required:    true,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The unique name of the virtual IP pool.`,
		},

		"subnet_cidr": {
			Type:          schema.TypeInt,
			ConflictsWith: codegen_configs.GetResourceByName("VipPool").GetConflictingFields("subnet_cidr"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) IPv4 subnet CIDR prefix (number of bits).`,
		},

		"subnet_cidr_ipv6": {
			Type:          schema.TypeInt,
			ConflictsWith: codegen_configs.GetResourceByName("VipPool").GetConflictingFields("subnet_cidr_ipv6"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) IPv6 subnet CIDR prefix (number of bits).`,
		},

		"gw_ip": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("VipPool").GetConflictingFields("gw_ip"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Gateway IPv4 address.`,
		},

		"gw_ipv6": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("VipPool").GetConflictingFields("gw_ipv6"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Gateway IPv6 address.`,
		},

		"vlan": {
			Type:          schema.TypeInt,
			ConflictsWith: codegen_configs.GetResourceByName("VipPool").GetConflictingFields("vlan"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The VLAN of the virtual IP pool.`,
		},

		"state": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("VipPool").GetConflictingFields("state"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,
		},

		"cnode_ids": {
			Type:          schema.TypeList,
			ConflictsWith: codegen_configs.GetResourceByName("VipPool").GetConflictingFields("cnode_ids"),

			DiffSuppressOnRefresh: false,
			DiffSuppressFunc:      codegen_configs.GetResourceByName("VipPool").GetAttributeDiffFunc("cnode_ids"),
			Computed:              true,
			Optional:              true,
			Sensitive:             false,
			Description:           `(Valid for versions: 5.0.0,5.1.0,5.2.0) IDs of CNodes comprising the CNode group.`,

			Elem: &schema.Schema{
				Type: schema.TypeInt,
			},
		},

		"cluster": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("VipPool").GetConflictingFields("cluster"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Parent cluster.`,
		},

		"url": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("VipPool").GetConflictingFields("url"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,
		},

		"domain_name": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("VipPool").GetConflictingFields("domain_name"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,
		},

		"role": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("VipPool").GetConflictingFields("role"),

			Required:    true,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Role.`,
		},

		"ip_ranges": {
			Type:          schema.TypeList,
			ConflictsWith: codegen_configs.GetResourceByName("VipPool").GetConflictingFields("ip_ranges"),

			Required:    true,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) IP ranges.`,

			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{

					"start_ip": &schema.Schema{
						Type:        schema.TypeString,
						Computed:    true,
						Optional:    true,
						Description: "",
					},

					"end_ip": &schema.Schema{
						Type:        schema.TypeString,
						Computed:    true,
						Optional:    true,
						Description: "",
					},
				},
			},
		},

		"vms_preferred": {
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("VipPool").GetConflictingFields("vms_preferred"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) If 'true', the CNodes included in this virtual IP pool are handled as preferred CNodes during VMS host election.`,
		},

		"enabled": {
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("VipPool").GetConflictingFields("enabled"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Enables or disables the virtual IP pool.`,
		},

		"port_membership": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("VipPool").GetConflictingFields("port_membership"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The port(s) on the CNode that this pool will use: 'Right', 'Left' or 'All'.`,
		},

		"active_interfaces": {
			Type:          schema.TypeInt,
			ConflictsWith: codegen_configs.GetResourceByName("VipPool").GetConflictingFields("active_interfaces"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The number of active interfaces.`,
		},

		"enable_l3": {
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("VipPool").GetConflictingFields("enable_l3"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Enables or disables L3 CNode access.`,
		},

		"vast_asn": {
			Type:          schema.TypeInt,
			ConflictsWith: codegen_configs.GetResourceByName("VipPool").GetConflictingFields("vast_asn"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) VAST ASN.`,
		},

		"peer_asn": {
			Type:          schema.TypeInt,
			ConflictsWith: codegen_configs.GetResourceByName("VipPool").GetConflictingFields("peer_asn"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Peer ASN.`,
		},

		"tenant_id": {
			Type:          schema.TypeInt,
			ConflictsWith: codegen_configs.GetResourceByName("VipPool").GetConflictingFields("tenant_id"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.1.0,5.2.0) The ID of the tenant associated with the virtual IP pool. An ID of '0' (zero) means the virtual IP pool is available for all tenants.`,
		},

		"active_cnode_ids": {
			Type:          schema.TypeList,
			ConflictsWith: codegen_configs.GetResourceByName("VipPool").GetConflictingFields("active_cnode_ids"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.1.0,5.2.0) IDs of active CNodes`,

			Elem: &schema.Schema{
				Type: schema.TypeInt,
			},
		},

		"cluster_id": {
			Type:          schema.TypeInt,
			ConflictsWith: codegen_configs.GetResourceByName("VipPool").GetConflictingFields("cluster_id"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.1.0,5.2.0) Cluster ID`,
		},

		"cnodes": {
			Type:          schema.TypeList,
			ConflictsWith: codegen_configs.GetResourceByName("VipPool").GetConflictingFields("cnodes"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.1.0,5.2.0) `,

			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},

		"enable_weighted_balancing": {
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("VipPool").GetConflictingFields("enable_weighted_balancing"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.1.0,5.2.0) Weighted Balancing Enabled`,
		},

		"ranges_summary": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("VipPool").GetConflictingFields("ranges_summary"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.1.0,5.2.0) IP ranges`,
		},

		"sync_time": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("VipPool").GetConflictingFields("sync_time"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.1.0,5.2.0) Synchronization time with leader`,
		},

		"sync": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("VipPool").GetConflictingFields("sync"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.1.0,5.2.0) Synchronization state with leader`,
		},

		"tenant_name": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("VipPool").GetConflictingFields("tenant_name"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.1.0,5.2.0) Tenant Name`,
		},

		"title": {
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("VipPool").GetConflictingFields("title"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.1.0,5.2.0) IP range of the VIP pool`,
		},
	}
}

var VipPoolNamesMapping = map[string][]string{
	"ip_ranges": []string{"start_ip", "end_ip"},
}

func ResourceVipPoolReadStructIntoSchema(ctx context.Context, resource api_latest.VipPool, d *schema.ResourceData) diag.Diagnostics {
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

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "SubnetCidr", resource.SubnetCidr))

	err = d.Set("subnet_cidr", resource.SubnetCidr)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"subnet_cidr\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "SubnetCidrIpv6", resource.SubnetCidrIpv6))

	err = d.Set("subnet_cidr_ipv6", resource.SubnetCidrIpv6)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"subnet_cidr_ipv6\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "GwIp", resource.GwIp))

	err = d.Set("gw_ip", resource.GwIp)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"gw_ip\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "GwIpv6", resource.GwIpv6))

	err = d.Set("gw_ipv6", resource.GwIpv6)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"gw_ipv6\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Vlan", resource.Vlan))

	err = d.Set("vlan", resource.Vlan)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"vlan\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "State", resource.State))

	err = d.Set("state", resource.State)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"state\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "CnodeIds", resource.CnodeIds))

	err = d.Set("cnode_ids", utils.FlattenListOfPrimitives(&resource.CnodeIds))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"cnode_ids\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Cluster", resource.Cluster))

	err = d.Set("cluster", resource.Cluster)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"cluster\"",
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

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "DomainName", resource.DomainName))

	err = d.Set("domain_name", resource.DomainName)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"domain_name\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Role", resource.Role))

	err = d.Set("role", resource.Role)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"role\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "IpRanges", resource.IpRanges))

	err = d.Set("ip_ranges", utils.FlattenListOfStringsList(&resource.IpRanges, []string{"start_ip", "end_ip"}))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"ip_ranges\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "VmsPreferred", resource.VmsPreferred))

	err = d.Set("vms_preferred", resource.VmsPreferred)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"vms_preferred\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Enabled", resource.Enabled))

	err = d.Set("enabled", resource.Enabled)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"enabled\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "PortMembership", resource.PortMembership))

	err = d.Set("port_membership", resource.PortMembership)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"port_membership\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "ActiveInterfaces", resource.ActiveInterfaces))

	err = d.Set("active_interfaces", resource.ActiveInterfaces)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"active_interfaces\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "EnableL3", resource.EnableL3))

	err = d.Set("enable_l3", resource.EnableL3)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"enable_l3\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "VastAsn", resource.VastAsn))

	err = d.Set("vast_asn", resource.VastAsn)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"vast_asn\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "PeerAsn", resource.PeerAsn))

	err = d.Set("peer_asn", resource.PeerAsn)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"peer_asn\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "TenantId", resource.TenantId))

	err = d.Set("tenant_id", resource.TenantId)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"tenant_id\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "ActiveCnodeIds", resource.ActiveCnodeIds))

	err = d.Set("active_cnode_ids", utils.FlattenListOfPrimitives(&resource.ActiveCnodeIds))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"active_cnode_ids\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "ClusterId", resource.ClusterId))

	err = d.Set("cluster_id", resource.ClusterId)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"cluster_id\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Cnodes", resource.Cnodes))

	err = d.Set("cnodes", utils.FlattenListOfPrimitives(&resource.Cnodes))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"cnodes\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "EnableWeightedBalancing", resource.EnableWeightedBalancing))

	err = d.Set("enable_weighted_balancing", resource.EnableWeightedBalancing)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"enable_weighted_balancing\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "RangesSummary", resource.RangesSummary))

	err = d.Set("ranges_summary", resource.RangesSummary)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"ranges_summary\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "SyncTime", resource.SyncTime))

	err = d.Set("sync_time", resource.SyncTime)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"sync_time\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Sync", resource.Sync))

	err = d.Set("sync", resource.Sync)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"sync\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "TenantName", resource.TenantName))

	err = d.Set("tenant_name", resource.TenantName)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"tenant_name\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Title", resource.Title))

	err = d.Set("title", resource.Title)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"title\"",
			Detail:   err.Error(),
		})
	}

	return diags

}
func resourceVipPoolRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*vast_client.VMSSession)
	resourceConfig := codegen_configs.GetResourceByName("VipPool")
	attrs := map[string]interface{}{"path": utils.GenPath("vippools"), "id": d.Id()}
	tflog.Debug(ctx, fmt.Sprintf("[resourceVipPoolRead] Calling Get Function : %v for resource VipPool", utils.GetFuncName(resourceConfig.GetFunc)))
	response, err := resourceConfig.GetFunc(ctx, client, attrs, d, map[string]string{})
	utils.VastVersionsWarn(ctx)

	var body []byte
	var resource api_latest.VipPool
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
	diags = ResourceVipPoolReadStructIntoSchema(ctx, resource, d)

	return diags
}

func resourceVipPoolDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*vast_client.VMSSession)
	resourceConfig := codegen_configs.GetResourceByName("VipPool")
	attrs := map[string]interface{}{"path": utils.GenPath("vippools"), "id": d.Id()}

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

func resourceVipPoolCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	namesMapping := utils.ContextKey("names_mapping")
	newCtx := context.WithValue(ctx, namesMapping, VipPoolNamesMapping)
	var diags diag.Diagnostics
	data := make(map[string]interface{})
	client := m.(*vast_client.VMSSession)
	resourceConfig := codegen_configs.GetResourceByName("VipPool")
	tflog.Info(ctx, fmt.Sprintf("Creating Resource VipPool"))
	reflectVipPool := reflect.TypeOf((*api_latest.VipPool)(nil))
	utils.PopulateResourceMap(newCtx, reflectVipPool.Elem(), d, &data, "", false)

	var before_post_error error
	data, before_post_error = resourceConfig.BeforePostFunc(data, client, ctx, d)
	if before_post_error != nil {
		return diag.FromErr(before_post_error)
	}

	versionsEqual := utils.VastVersionsWarn(ctx)

	if versionsEqual != metadata.CLUSTER_VERSION_EQUALS {
		clusterVersion := metadata.ClusterVersionString()
		t, typeExists := vast_versions.GetVersionedType(clusterVersion, "VipPool")
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
			tflog.Warn(ctx, fmt.Sprintf("Could have not found resource %s in version %s, things might not work properly", "VipPool", clusterVersion))
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
	attrs := map[string]interface{}{"path": utils.GenPath("vippools")}
	response, createErr := resourceConfig.CreateFunc(ctx, client, attrs, data, map[string]string{})
	tflog.Info(ctx, fmt.Sprintf("Server Error for  VipPool %v", createErr))

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
	resource := api_latest.VipPool{}
	err = json.Unmarshal(responseBody, &resource)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to convert response body into VipPool",
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
	resourceVipPoolRead(ctxWithResource, d, m)

	return diags
}

func resourceVipPoolUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	namesMapping := utils.ContextKey("names_mapping")
	newCtx := context.WithValue(ctx, namesMapping, VipPoolNamesMapping)
	var diags diag.Diagnostics
	data := make(map[string]interface{})
	versionsEqual := utils.VastVersionsWarn(ctx)
	resourceConfig := codegen_configs.GetResourceByName("VipPool")
	if versionsEqual != metadata.CLUSTER_VERSION_EQUALS {
		clusterVersion := metadata.ClusterVersionString()
		t, typeExists := vast_versions.GetVersionedType(clusterVersion, "VipPool")
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
			tflog.Warn(ctx, fmt.Sprintf("Could have not found resource %s in version %s, things might not work properly", "VipPool", clusterVersion))
		}
	}

	client := m.(*vast_client.VMSSession)
	tflog.Info(ctx, fmt.Sprintf("Updating Resource VipPool"))
	reflectVipPool := reflect.TypeOf((*api_latest.VipPool)(nil))
	utils.PopulateResourceMap(newCtx, reflectVipPool.Elem(), d, &data, "", false)

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
	attrs := map[string]interface{}{"path": utils.GenPath("vippools"), "id": d.Id()}
	response, patchErr := resourceConfig.UpdateFunc(ctx, client, attrs, data, d, map[string]string{})
	tflog.Info(ctx, fmt.Sprintf("Server Error for  VipPool %v", patchErr))
	if patchErr != nil {
		errorMessage := fmt.Sprintf("server response:\n%v\nUnderlying error:\n%v", utils.GetResponseBodyAsStr(response), patchErr.Error())
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Object Creation Failed",
			Detail:   errorMessage,
		})
		return diags
	}
	resourceVipPoolRead(ctx, d, m)

	return diags

}

func resourceVipPoolImporter(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {

	var result []*schema.ResourceData
	client := m.(*vast_client.VMSSession)
	resourceConfig := codegen_configs.GetResourceByName("VipPool")
	attrs := map[string]interface{}{"path": utils.GenPath("vippools")}
	response, err := resourceConfig.ImportFunc(ctx, client, attrs, d, resourceConfig.Importer.GetFunc())

	if err != nil {
		return result, err
	}

	var resourceList []api_latest.VipPool
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

	diags := ResourceVipPoolReadStructIntoSchema(ctx, resource, d)
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
