package resources

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api_latest "github.com/vast-data/terraform-provider-vastdata/codegen/latest"
	metadata "github.com/vast-data/terraform-provider-vastdata/metadata"
	utils "github.com/vast-data/terraform-provider-vastdata/utils"
	vast_client "github.com/vast-data/terraform-provider-vastdata/vast-client"
	vast_versions "github.com/vast-data/terraform-provider-vastdata/vast_versions"
	"io"
	"net/url"
	"reflect"
	"strconv"
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
			Type:     schema.TypeString,
			Required: true,
		},

		"vip": &schema.Schema{
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `Assigns a IP to the DNS service. DNS requests from your external DNS server must be delegated to this IP.`,
		},

		"domain_suffix": &schema.Schema{
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `Specifies a suffix to append to domain names of each VIP pool. The suffix should complete each domain name to form a valid FQDN for DNS requests to target.`,
		},

		"vip_gateway": &schema.Schema{
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `Specifies a gateway IP to external DNS server if on different subnet. Must be on same subnet as the IP and reachable from the relevant nework interface.`,
		},

		"enabled": &schema.Schema{
			Type:        schema.TypeBool,
			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `Enable the VAST DNS server configurations`,
		},

		"guid": &schema.Schema{
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    false,
			Sensitive:   false,
			Description: `A uniqe guid assigned to the VAST DNS server configurations`,
		},

		"vip_subnet_cidr": &schema.Schema{
			Type:        schema.TypeInt,
			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `Specifies the subnet, as a CIDR index, on which the DNS resides.`,
		},

		"vip_vlan": &schema.Schema{
			Type:        schema.TypeInt,
			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `Specifies a VLAN if needed to enable communication with external DNS server(s).`,
		},

		"cnode_ids": &schema.Schema{
			Type:        schema.TypeList,
			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: ``,

			Elem: &schema.Schema{
				Type: schema.TypeInt,
			},
		},

		"vip_ipv6": &schema.Schema{
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `Assigns an IPv6 to the DNS service.`,
		},

		"vip_ipv6_subnet_cidr": &schema.Schema{
			Type:        schema.TypeInt,
			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `Specifies the subnet, as a CIDR index, on which the DNS resides. [1..128]`,
		},

		"vip_ipv6_gateway": &schema.Schema{
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `Specifies a gateway IPv6 to external DNS server if on different subnet.`,
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

	return diags

}
func resourceDnsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(vast_client.JwtSession)

	attrs := map[string]interface{}{"path": "/api/latest/dns/", "id": d.Id()}
	response, err := utils.DefaultGetFunc(ctx, client, attrs, map[string]string{})
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
	resource := api_latest.Dns{}
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
	diags = ResourceDnsReadStructIntoSchema(ctx, resource, d)

	return diags
}

func resourceDnsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(vast_client.JwtSession)
	attrs := map[string]interface{}{"path": "/api/latest/dns/", "id": d.Id()}

	response, err := utils.DefaultDeleteFunc(ctx, client, attrs, nil, map[string]string{})

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

func resourceDnsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	names_mapping := utils.ContextKey("names_mapping")
	new_ctx := context.WithValue(ctx, names_mapping, Dns_names_mapping)
	var diags diag.Diagnostics
	data := make(map[string]interface{})
	client := m.(vast_client.JwtSession)
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
	attrs := map[string]interface{}{"path": "/api/latest/dns/"}
	response, create_err := utils.DefaultCreateFunc(ctx, client, attrs, data, map[string]string{})
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

	id_err := utils.DefaultIdFunc(ctx, client, resource.Id, d)
	if id_err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to set Id",
			Detail:   err.Error(),
		})
		return diags
	}
	resourceDnsRead(ctx, d, m)

	return diags
}

func resourceDnsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	names_mapping := utils.ContextKey("names_mapping")
	new_ctx := context.WithValue(ctx, names_mapping, Dns_names_mapping)
	var diags diag.Diagnostics
	data := make(map[string]interface{})
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

	client := m.(vast_client.JwtSession)
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
	attrs := map[string]interface{}{"path": "/api/latest/dns/", "id": d.Id()}
	response, patch_err := utils.DefaultUpdateFunc(ctx, client, attrs, data, map[string]string{})
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
	client := m.(vast_client.JwtSession)
	guid := d.Id()
	values := url.Values{}
	values.Add("guid", fmt.Sprintf("%v", guid))
	attrs := map[string]interface{}{"path": "/api/latest/dns/", "query": values.Encode()}
	response, err := utils.DefaultGetFunc(ctx, client, attrs, map[string]string{})

	if err != nil {
		return result, err
	}

	resource_l := []api_latest.Dns{}

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

	Id := (int64)(resource.Id)
	d.SetId(strconv.FormatInt(Id, 10))
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
