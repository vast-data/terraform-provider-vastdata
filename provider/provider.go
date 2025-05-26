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
				Description: `The VastData Cluster hostname/address , if environment variable VASTDATA_HOST exists it will be used`,
				DefaultFunc: schema.EnvDefaultFunc("VASTDATA_HOST", nil),
			},
			"port": {
				Type:        schema.TypeInt,
				Required:    false,
				Optional:    true,
				Description: `The server API port (Default is 443) ,if environment variable VASTDATA_PORT exists it will be used`,
				DefaultFunc: schema.EnvDefaultFunc("VASTDATA_PORT", 443),
			},
			"skip_ssl_verify": {
				Type:        schema.TypeBool,
				Required:    false,
				Optional:    true,
				Description: `A boolean representing should SSL certificate be verified (Default is False) , if environmnet variable VASTDATA_VERIFY_SSL exists it will be used`,
				DefaultFunc: schema.EnvDefaultFunc("VASTDATA_VERIFY_SSL", false),
			},

			"username": {
				Type:          schema.TypeString,
				Required:      false,
				Optional:      true,
				Sensitive:     true,
				Description:   `The VastData Cluster username, if environment variable VASTDATA_CLUSTER_USERNAME exists it will be used`,
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
				Description:   `The VastData Cluster password, if environment variable VASTDATA_CLUSTER_PASSWORD exists it will be used`,
				DefaultFunc:   schema.EnvDefaultFunc("VASTDATA_CLUSTER_PASSWORD", nil),
				ConflictsWith: []string{"api_token"},
				RequiredWith:  []string{"username"},
			},
			"api_token": {
				Type:          schema.TypeString,
				Required:      false,
				Optional:      true,
				Sensitive:     true,
				Description:   `The VastData Cluster API token. If environment variable VASTDATA_API_TOKEN exists it will be used`,
				DefaultFunc:   schema.EnvDefaultFunc("VASTDATA_API_TOKEN", nil),
				ConflictsWith: []string{"username", "password"},
				AtLeastOneOf:  []string{"api_token", "username"},
			},
			"version_validation_mode": {
				Type:      schema.TypeString,
				Required:  false,
				Optional:  true,
				Sensitive: false,
				Description: `The version validation mode to use , version validation checks if a resource request will work with the current cluster version
			Depending on the value the operation will abort from happening if according to the version the operation might not work.
			2 options are valid for this attribute
			1. strict - abort the operation before it starts
			2. warn - Just issue a warning `,
				DefaultFunc:  schema.EnvDefaultFunc("VERSION_VALIDATION_MODE", "warn"),
				ValidateFunc: validation.StringInSlice([]string{"warn", "strict"}, true),
			},
		},
		ConfigureContextFunc: providerConfigure,
	}

}

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
			Summary:  "Unable to start a session to the vastdata cluster",
			Detail:   err.Error(),
		})
		return client, diags
	}
	clusterVersion, _, err := client.ClusterVersion(ctx)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error obtaining VAST Data cluster version",
			Detail:   fmt.Sprintf("This is the first hit to the API. Did you provide correct credentials?\nUnderlying error is:\n%v", err.Error()),
		})
		return client, diags
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
