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

func DataSourceCnode() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCnodeRead,
		Description: ``,
		Schema: map[string]*schema.Schema{

			"id": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: ``,
			},

			"guid": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: ``,
			},

			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    false,
				Required:    true,
				Optional:    false,
				Description: ``,
			},

			"new_name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: ``,
			},

			"ip": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `Currently used IP Address, bond of ip1 and ip2`,
			},

			"ip1": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `1st internal IP Address`,
			},

			"ip2": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `2nd internal IP Address`,
			},

			"ipv6": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `External IPv6 Address`,
			},

			"host_label": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `Host label, used to label container, e.g. 11.0.0.1-4000`,
			},

			"sn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `The host serial number`,
			},

			"state": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: ``,
			},

			"display_state": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: ``,
			},

			"led_status": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: ``,
			},

			"platform_rdma_port": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: ``,
			},

			"platform_tcp_port": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: ``,
			},

			"data_rdma_port": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: ``,
			},

			"data_tcp_port": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: ``,
			},

			"enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `cnode enabled`,
			},

			"is_mgmt": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `The cnode is running management`,
			},

			"cbox": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `Parent CBox`,
			},

			"cbox_uid": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `Unique Parent CBox identifier`,
			},

			"cbox_id": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `Parent CBox id`,
			},

			"cluster": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `Parent Cluster`,
			},

			"os_version": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `Node OS version`,
			},

			"bmc_fw_version": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `BMC FW version`,
			},

			"vms_preferred": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `VMS preferred CNode`,
			},

			"url": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: ``,
			},

			"mgmt_ip": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `Management IP`,
			},

			"hostname": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `Host Name`,
			},
		},
	}
}

func dataSourceCnodeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(vast_client.JwtSession)
	values := url.Values{}

	name := d.Get("name")
	values.Add("name", fmt.Sprintf("%v", name))

	response, err := client.Get(ctx, utils.GenPath("cnodes"), values.Encode(), map[string]string{})
	tflog.Info(ctx, response.Request.URL.String())
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured while obtaining data from the vastdata cluster",
			Detail:   err.Error(),
		})
		return diags

	}
	resource_l := []api_latest.Cnode{}

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

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Name", resource.Name))

	err = d.Set("name", resource.Name)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"name\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "NewName", resource.NewName))

	err = d.Set("new_name", resource.NewName)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"new_name\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Ip", resource.Ip))

	err = d.Set("ip", resource.Ip)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"ip\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Ip1", resource.Ip1))

	err = d.Set("ip1", resource.Ip1)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"ip1\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Ip2", resource.Ip2))

	err = d.Set("ip2", resource.Ip2)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"ip2\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Ipv6", resource.Ipv6))

	err = d.Set("ipv6", resource.Ipv6)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"ipv6\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "HostLabel", resource.HostLabel))

	err = d.Set("host_label", resource.HostLabel)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"host_label\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Sn", resource.Sn))

	err = d.Set("sn", resource.Sn)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"sn\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "State", resource.State))

	err = d.Set("state", resource.State)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"state\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "DisplayState", resource.DisplayState))

	err = d.Set("display_state", resource.DisplayState)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"display_state\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "LedStatus", resource.LedStatus))

	err = d.Set("led_status", resource.LedStatus)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"led_status\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "PlatformRdmaPort", resource.PlatformRdmaPort))

	err = d.Set("platform_rdma_port", resource.PlatformRdmaPort)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"platform_rdma_port\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "PlatformTcpPort", resource.PlatformTcpPort))

	err = d.Set("platform_tcp_port", resource.PlatformTcpPort)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"platform_tcp_port\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "DataRdmaPort", resource.DataRdmaPort))

	err = d.Set("data_rdma_port", resource.DataRdmaPort)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"data_rdma_port\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "DataTcpPort", resource.DataTcpPort))

	err = d.Set("data_tcp_port", resource.DataTcpPort)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"data_tcp_port\"",
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

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "IsMgmt", resource.IsMgmt))

	err = d.Set("is_mgmt", resource.IsMgmt)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"is_mgmt\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Cbox", resource.Cbox))

	err = d.Set("cbox", resource.Cbox)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"cbox\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "CboxUid", resource.CboxUid))

	err = d.Set("cbox_uid", resource.CboxUid)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"cbox_uid\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "CboxId", resource.CboxId))

	err = d.Set("cbox_id", resource.CboxId)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"cbox_id\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Cluster", resource.Cluster))

	err = d.Set("cluster", resource.Cluster)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"cluster\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "OsVersion", resource.OsVersion))

	err = d.Set("os_version", resource.OsVersion)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"os_version\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "BmcFwVersion", resource.BmcFwVersion))

	err = d.Set("bmc_fw_version", resource.BmcFwVersion)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"bmc_fw_version\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "VmsPreferred", resource.VmsPreferred))

	err = d.Set("vms_preferred", resource.VmsPreferred)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"vms_preferred\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Url", resource.Url))

	err = d.Set("url", resource.Url)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"url\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "MgmtIp", resource.MgmtIp))

	err = d.Set("mgmt_ip", resource.MgmtIp)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"mgmt_ip\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Hostname", resource.Hostname))

	err = d.Set("hostname", resource.Hostname)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"hostname\"",
			Detail:   err.Error(),
		})
	}

	Id := (int64)(resource.Id)
	d.SetId(strconv.FormatInt(Id, 10))
	return diags
}
