package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api_latest "github.com/vast-data/terraform-provider-vastdata.git/codegen/latest"
	utils "github.com/vast-data/terraform-provider-vastdata.git/utils"
	vast_client "github.com/vast-data/terraform-provider-vastdata.git/vast-client"
	"net/url"
	"strconv"
)

func DataSourceS3replicationPeers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceS3replicationPeersRead,
		Schema: map[string]*schema.Schema{

			"id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
				Required: false,
				Optional: false,
			},

			"guid": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Required: false,
				Optional: false,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: false,
				Required: true,
				Optional: false,
			},

			"url": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Required: false,
				Optional: false,
			},

			"bucket_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Required: false,
				Optional: false,
			},

			"http_protocol": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Required: false,
				Optional: false,
			},

			"type_": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Required: false,
				Optional: false,
			},

			"proxies": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Required: false,
				Optional: false,

				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"aws_region": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Required: false,
				Optional: false,
			},
		},
	}
}

func dataSourceS3replicationPeersRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(vast_client.JwtSession)
	values := url.Values{}

	name := d.Get("name")
	values.Add("name", fmt.Sprintf("%v", name))

	response, err := client.Get(ctx, "/api/replicationtargets/", values.Encode(), map[string]string{})
	tflog.Info(ctx, response.Request.URL.String())
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured while obtaining data from the vastdata cluster",
			Detail:   err.Error(),
		})
		return diags

	}
	resource_l := []api_latest.S3replicationPeers{}

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

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Url", resource.Url))

	err = d.Set("url", resource.Url)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"url\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "BucketName", resource.BucketName))

	err = d.Set("bucket_name", resource.BucketName)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"bucket_name\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "HttpProtocol", resource.HttpProtocol))

	err = d.Set("http_protocol", resource.HttpProtocol)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"http_protocol\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Type_", resource.Type_))

	err = d.Set("type_", resource.Type_)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"type_\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Proxies", resource.Proxies))

	err = d.Set("proxies", utils.FlattenListOfPrimitives(&resource.Proxies))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"proxies\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "AwsRegion", resource.AwsRegion))

	err = d.Set("aws_region", resource.AwsRegion)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"aws_region\"",
			Detail:   err.Error(),
		})
	}

	Id := (int64)(resource.Id)
	d.SetId(strconv.FormatInt(Id, 10))
	return diags
}