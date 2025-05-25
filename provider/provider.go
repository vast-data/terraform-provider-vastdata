package vastdata

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	datasources "github.com/vast-data/terraform-provider-vastdata/datasources"
	resources "github.com/vast-data/terraform-provider-vastdata/resources"

	metadata "github.com/vast-data/terraform-provider-vastdata/metadata"
	vast_client "github.com/vast-data/terraform-provider-vastdata/vast-client"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap:   resources.Resources,
		DataSourcesMap: datasources.DataSources,
		Schema: map[string]*schema.Schema{
			"host": {
				Type:        schema.TypeString,
				Required:    true,
				Optional:    false,
				Description: `The VAST cluster hostname/address. If the VASTDATA_HOST environment variable exists, it will be used.`,
				DefaultFunc: schema.EnvDefaultFunc("VASTDATA_HOST", nil),
			},
			"port": {
				Type:        schema.TypeInt,
				Required:    false,
				Optional:    true,
				Description: `The server API port (default is 443). If the VASTDATA_PORT environment variable exists, it will be used.`,
				DefaultFunc: schema.EnvDefaultFunc("VASTDATA_PORT", 443),
			},
			"skip_ssl_verify": {
				Type:        schema.TypeBool,
				Required:    false,
				Optional:    true,
				Description: `A flag to determine whether the SSL certificate is to be verified (default is False). If the VASTDATA_VERIFY_SSL environment variable exists, it will be used.`,
				DefaultFunc: schema.EnvDefaultFunc("VASTDATA_VERIFY_SSL", false),
			},

			"username": {
				Type:          schema.TypeString,
				Required:      false,
				Optional:      true,
				Sensitive:     true,
				Description:   `The VAST cluster user name. If the VASTDATA_CLUSTER_USERNAME environment variable exists, it will be used.`,
				DefaultFunc:   schema.EnvDefaultFunc("VASTDATA_CLUSTER_USERNAME", nil),
				ConflictsWith: []string{"api_token"},
				RequiredWith:  []string{"password"},
				AtLeastOneOf:  []string{"api_token", "username"},
			},
			"password": {
				Type:          schema.TypeString,
				Required:      false,
				Optional:      true,
				Sensitive:     true,
				Description:   `The VAST cluster password. If the VASTDATA_CLUSTER_PASSWORD environment variable exists, it will be used.`,
				DefaultFunc:   schema.EnvDefaultFunc("VASTDATA_CLUSTER_PASSWORD", nil),
				ConflictsWith: []string{"api_token"},
				RequiredWith:  []string{"username"},
			},
			"api_token": {
				Type:          schema.TypeString,
				Required:      false,
				Optional:      true,
				Sensitive:     true,
				Description:   `The VAST cluster API token. If the VASTDATA_API_TOKEN environment variable exists, it will be used.`,
				DefaultFunc:   schema.EnvDefaultFunc("VASTDATA_API_TOKEN", nil),
				ConflictsWith: []string{"username", "password"},
				AtLeastOneOf:  []string{"api_token", "username"},
			},
			"version_validation_mode": {
				Type:      schema.TypeString,
				Required:  false,
				Optional:  true,
				Sensitive: false,
				Description: `The version validation mode to use. Version validation checks if a resource request will work with the current cluster version. Depending on the value, the operation will be aborted if it won't work with the current version. Valid values: strict - to abort the operation before it starts, warn - to issue a warning without aborting the operation.`,
				DefaultFunc:  schema.EnvDefaultFunc("VERSION_VALIDATION_MODE", "warn"),
				ValidateFunc: validation.StringInSlice([]string{"warn", "strict"}, true),
			},
		},
		ConfigureContextFunc: providerConfigure,
	}

}

/*
	Provide a VastData client which have been started ,
	this will validate the following:
        *) the cluster is up and responding
	*) the username & password are valid
*/

func providerConfigure(ctx context.Context, r *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	config := &vast_client.RestClientConfig{
		Host:      r.Get("host").(string),
		Port:      uint64(r.Get("port").(int)),
		Username:  r.Get("username").(string),
		Password:  r.Get("password").(string),
		ApiToken:  r.Get("api_token").(string),
		SslVerify: !r.Get("skip_ssl_verify").(bool),
	}
	client := vast_client.NewSession(ctx, config)
	err := client.Start()

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to start a session to the VAST cluster",
			Detail:   err.Error(),
		})
		return client, diags
	}
	clusterVersion, _, err := client.ClusterVersion(ctx)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error obtaning cluster version",
			Detail:   err.Error(),
		})

	}
	tflog.Info(ctx, fmt.Sprintf("Cluster version found %s", clusterVersion))
	clusterVersion, truncated := metadata.SanitizeVersion(clusterVersion)
	if truncated {
		tflog.Info(ctx, fmt.Sprintf("Cluster version truncated to: %s", clusterVersion))
	}
	err = metadata.UpdateClusterVersion(clusterVersion)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred while updating cluster version",
			Detail:   err.Error(),
		})
		return client, diags
	}
	if metadata.ClusterVersionCompare() != metadata.CLUSTER_VERSION_EQUALS {
		tflog.Warn(ctx, "Cluster Version & Build Version are not matching, some actions might fail")
	}
	metadata.SetClusterConfig("version_validation_mode", r.Get("version_validation_mode").(string))
	v := metadata.FindVastVersion(clusterVersion)
	metadata.SetClusterConfig("vast_version", v)
	tflog.Debug(ctx, fmt.Sprintf("API Version than will be used %v", v))
	return client, diags
}
