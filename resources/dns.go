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

func ResourceDns() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceDnsRead,
		DeleteContext: resourceDnsDelete,
		CreateContext: resourceDnsCreate,
		UpdateContext: resourceDnsUpdate,

		Importer: &schema.ResourceImporter{
			StateContext: resourceDnsImporter,
		},

		Description: ``,
		Schema:      getResourceDnsSchema(),
	}
}

func getResourceDnsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{

		"name": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Dns").GetConflictingFields("name"),

			Required:    true,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Specifies a name for the VAST DNS server configuration`,
		},

		"vip": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Dns").GetConflictingFields("vip"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Assigns a IP to the DNS service. DNS requests from your external DNS server must be delegated to this IP.`,
		},

		"domain_suffix": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Dns").GetConflictingFields("domain_suffix"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Specifies a suffix to append to domain names of each VIP pool. The suffix should complete each domain name to form a valid FQDN for DNS requests to target.`,
		},

		"vip_gateway": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Dns").GetConflictingFields("vip_gateway"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Specifies a gateway IP to external DNS server if on different subnet. Must be on same subnet as the IP and reachable from the relevant nework interface.`,
		},

		"enabled": &schema.Schema{
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("Dns").GetConflictingFields("enabled"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Enable the VAST DNS server configurations`,
		},

		"guid": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Dns").GetConflictingFields("guid"),

			Computed:    true,
			Optional:    false,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) A uniqe guid assigned to the VAST DNS server configurations`,
		},

		"vip_subnet_cidr": &schema.Schema{
			Type:          schema.TypeInt,
			ConflictsWith: codegen_configs.GetResourceByName("Dns").GetConflictingFields("vip_subnet_cidr"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Specifies the subnet, as a CIDR index, on which the DNS resides.`,
		},

		"vip_vlan": &schema.Schema{
			Type:          schema.TypeInt,
			ConflictsWith: codegen_configs.GetResourceByName("Dns").GetConflictingFields("vip_vlan"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Specifies a VLAN if needed to enable communication with external DNS server(s).`,
		},

		"cnode_ids": &schema.Schema{
			Type:          schema.TypeList,
			ConflictsWith: codegen_configs.GetResourceByName("Dns").GetConflictingFields("cnode_ids"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,

			Elem: &schema.Schema{
				Type: schema.TypeInt,
			},
		},

		"vip_ipv6": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Dns").GetConflictingFields("vip_ipv6"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Assigns an IPv6 to the DNS service.`,
		},

		"vip_ipv6_subnet_cidr": &schema.Schema{
			Type:          schema.TypeInt,
			ConflictsWith: codegen_configs.GetResourceByName("Dns").GetConflictingFields("vip_ipv6_subnet_cidr"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Specifies the subnet, as a CIDR index, on which the DNS resides. [1..128]`,
		},

		"vip_ipv6_gateway": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Dns").GetConflictingFields("vip_ipv6_gateway"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Specifies a gateway IPv6 to external DNS server if on different subnet.`,
		},

		"net_type": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Dns").GetConflictingFields("net_type"),

			Computed:  true,
			Optional:  true,
			Sensitive: false,

			ValidateDiagFunc: utils.OneOf([]string{"NORTH_PORT", "SOUTH_PORT", "EXTERNAL_PORT"}),
			Description:      `(Valid for versions: 5.1.0,5.2.0) Select the interface, that listens for DNS service delegation requests Allowed Values are [NORTH_PORT SOUTH_PORT EXTERNAL_PORT]`,
		},

		"invalid_name_response": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Dns").GetConflictingFields("invalid_name_response"),

			Computed:  true,
			Optional:  true,
			Sensitive: false,

			ValidateDiagFunc: utils.OneOf([]string{"NXDOMAIN", "REFUSED", "SERVFAIL", "NOERROR"}),
			Description:      `(Valid for versions: 5.1.0,5.2.0) The response DNS type for invalid dns name Allowed Values are [NXDOMAIN REFUSED SERVFAIL NOERROR]`,
		},

		"invalid_type_response": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("Dns").GetConflictingFields("invalid_type_response"),

			Computed:  true,
			Optional:  true,
			Sensitive: false,

			ValidateDiagFunc: utils.OneOf([]string{"NXDOMAIN", "REFUSED", "SERVFAIL", "NOERROR"}),
			Description:      `(Valid for versions: 5.1.0,5.2.0) The response DNS type for invalid dns type Allowed Values are [NXDOMAIN REFUSED SERVFAIL NOERROR]`,
		},

		"ttl": &schema.Schema{
			Type:          schema.TypeInt,
			ConflictsWith: codegen_configs.GetResourceByName("Dns").GetConflictingFields("ttl"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.1.0,5.2.0) The reposne TTL in seconds`,
		},
	}
}

var Dns_names_mapping map[string][]string = map[string][]string{}

func ResourceDnsReadStructIntoSchema(ctx context.Context, resource api_latest.Dns, d *schema.ResourceData) diag.Diagnostics {
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

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Vip", resource.Vip))

	err = d.Set("vip", resource.Vip)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"vip\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "DomainSuffix", resource.DomainSuffix))

	err = d.Set("domain_suffix", resource.DomainSuffix)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"domain_suffix\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "VipGateway", resource.VipGateway))

	err = d.Set("vip_gateway", resource.VipGateway)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"vip_gateway\"",
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

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Guid", resource.Guid))

	err = d.Set("guid", resource.Guid)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"guid\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "VipSubnetCidr", resource.VipSubnetCidr))

	err = d.Set("vip_subnet_cidr", resource.VipSubnetCidr)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"vip_subnet_cidr\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "VipVlan", resource.VipVlan))

	err = d.Set("vip_vlan", resource.VipVlan)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"vip_vlan\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "CnodeIds", resource.CnodeIds))

	err = d.Set("cnode_ids", utils.FlattenListOfPrimitives(&resource.CnodeIds))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"cnode_ids\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "VipIpv6", resource.VipIpv6))

	err = d.Set("vip_ipv6", resource.VipIpv6)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"vip_ipv6\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "VipIpv6SubnetCidr", resource.VipIpv6SubnetCidr))

	err = d.Set("vip_ipv6_subnet_cidr", resource.VipIpv6SubnetCidr)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"vip_ipv6_subnet_cidr\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "VipIpv6Gateway", resource.VipIpv6Gateway))

	err = d.Set("vip_ipv6_gateway", resource.VipIpv6Gateway)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"vip_ipv6_gateway\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "NetType", resource.NetType))

	err = d.Set("net_type", resource.NetType)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"net_type\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "InvalidNameResponse", resource.InvalidNameResponse))

	err = d.Set("invalid_name_response", resource.InvalidNameResponse)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"invalid_name_response\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "InvalidTypeResponse", resource.InvalidTypeResponse))

	err = d.Set("invalid_type_response", resource.InvalidTypeResponse)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"invalid_type_response\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Ttl", resource.Ttl))

	err = d.Set("ttl", resource.Ttl)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"ttl\"",
			Detail:   err.Error(),
		})
	}

	return diags

}
func resourceDnsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*vast_client.VMSSession)
	resource_config := codegen_configs.GetResourceByName("Dns")
	attrs := map[string]interface{}{"path": utils.GenPath("dns"), "id": d.Id()}
	tflog.Debug(ctx, fmt.Sprintf("[resourceDnsRead] Calling Get Function : %v for resource Dns", utils.GetFuncName(resource_config.GetFunc)))
	response, err := resource_config.GetFunc(ctx, client, attrs, d, map[string]string{})
	utils.VastVersionsWarn(ctx)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured while obtaining data from the vastdata cluster",
			Detail:   err.Error(),
		})
		return diags

	}
	tflog.Info(ctx, response.Request.URL.String())
	resource := api_latest.Dns{}
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
	diags = ResourceDnsReadStructIntoSchema(ctx, resource, d)

	return diags
}

func resourceDnsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*vast_client.VMSSession)
	resource_config := codegen_configs.GetResourceByName("Dns")
	attrs := map[string]interface{}{"path": utils.GenPath("dns"), "id": d.Id()}

	response, err := resource_config.DeleteFunc(ctx, client, attrs, nil, map[string]string{})

	tflog.Info(ctx, fmt.Sprintf("Removing Resource"))
	if response != nil {
		tflog.Info(ctx, response.Request.URL.String())
		tflog.Info(ctx, utils.GetResponseBodyAsStr(response))
	}

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured while deleting a resource from the vastdata cluster",
			Detail:   err.Error(),
		})

	}

	return diags

}

func resourceDnsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	names_mapping := utils.ContextKey("names_mapping")
	new_ctx := context.WithValue(ctx, names_mapping, Dns_names_mapping)
	var diags diag.Diagnostics
	data := make(map[string]interface{})
	client := m.(*vast_client.VMSSession)
	resource_config := codegen_configs.GetResourceByName("Dns")
	tflog.Info(ctx, fmt.Sprintf("Creating Resource Dns"))
	reflect_Dns := reflect.TypeOf((*api_latest.Dns)(nil))
	utils.PopulateResourceMap(new_ctx, reflect_Dns.Elem(), d, &data, "", false)

	version_compare := utils.VastVersionsWarn(ctx)

	if version_compare != metadata.CLUSTER_VERSION_EQUALS {
		cluster_version := metadata.ClusterVersionString()
		t, t_exists := vast_versions.GetVersionedType(cluster_version, "Dns")
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
			tflog.Warn(ctx, fmt.Sprintf("Could have not found resource %s in version %s , things might not work properly", "Dns", cluster_version))
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
	attrs := map[string]interface{}{"path": utils.GenPath("dns")}
	response, create_err := resource_config.CreateFunc(ctx, client, attrs, data, map[string]string{})
	tflog.Info(ctx, fmt.Sprintf("Server Error for  Dns %v", create_err))

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
	resource := api_latest.Dns{}
	err = json.Unmarshal(response_body, &resource)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to convert response body into Dns",
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
	resourceDnsRead(ctx_with_resource, d, m)

	return diags
}

func resourceDnsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	names_mapping := utils.ContextKey("names_mapping")
	new_ctx := context.WithValue(ctx, names_mapping, Dns_names_mapping)
	var diags diag.Diagnostics
	data := make(map[string]interface{})
	version_compare := utils.VastVersionsWarn(ctx)
	resource_config := codegen_configs.GetResourceByName("Dns")
	if version_compare != metadata.CLUSTER_VERSION_EQUALS {
		cluster_version := metadata.ClusterVersionString()
		t, t_exists := vast_versions.GetVersionedType(cluster_version, "Dns")
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
			tflog.Warn(ctx, fmt.Sprintf("Could have not found resource %s in version %s , things might not work properly", "Dns", cluster_version))
		}
	}

	client := m.(*vast_client.VMSSession)
	tflog.Info(ctx, fmt.Sprintf("Updating Resource Dns"))
	reflect_Dns := reflect.TypeOf((*api_latest.Dns)(nil))
	utils.PopulateResourceMap(new_ctx, reflect_Dns.Elem(), d, &data, "", false)

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
	attrs := map[string]interface{}{"path": utils.GenPath("dns"), "id": d.Id()}
	response, patch_err := resource_config.UpdateFunc(ctx, client, attrs, data, d, map[string]string{})
	tflog.Info(ctx, fmt.Sprintf("Server Error for  Dns %v", patch_err))
	if patch_err != nil {
		error_message := patch_err.Error() + " Server Response: " + utils.GetResponseBodyAsStr(response)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Object Creation Failed",
			Detail:   error_message,
		})
		return diags
	}
	resourceDnsRead(ctx, d, m)

	return diags

}

func resourceDnsImporter(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {

	result := []*schema.ResourceData{}
	client := m.(*vast_client.VMSSession)
	resource_config := codegen_configs.GetResourceByName("Dns")
	attrs := map[string]interface{}{"path": utils.GenPath("dns")}
	response, err := resource_config.ImportFunc(ctx, client, attrs, d, resource_config.Importer.GetFunc())

	if err != nil {
		return result, err
	}

	resource_l := []api_latest.Dns{}
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

	diags := ResourceDnsReadStructIntoSchema(ctx, resource, d)
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
