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

/*
	func validateVersionValidationMode(i any, c cty.Path) diag.Diagnostic {
		var d diag.Diagnostic
		s := strings.ToLower(fmt.Sprintf("%v", i))
		if (s != "warn") && (s != "strict") {
			return diag.FromErr(errors.New("Wrong value given when setting version_validation_mode %s, only possiable values are \"warn\" or \"strict\""))

		}
		return d

}
*/
func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap:   resources.Resources,
		DataSourcesMap: datasources.DataSources,
		Schema: map[string]*schema.Schema{
			"host": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    false,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("VASTDATA_HOST", nil),
			},
			"port": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Required:    false,
				DefaultFunc: schema.EnvDefaultFunc("VASTDATA_PORT", 443),
			},
			"skip_ssl_verify": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Required:    false,
				DefaultFunc: schema.EnvDefaultFunc("VASTDATA_VERIFY_SSL", false),
			},

			"username": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    false,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("VASTDATA_CLUSTER_USERNAME", nil),
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    false,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("VASTDATA_CLUSTER_USER", nil),
			},
			"version_validation_mode": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Required:     false,
				Sensitive:    false,
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

	client := vast_client.NewJwtSession(
		r.Get("host").(string),
		r.Get("username").(string),
		r.Get("password").(string),
		uint64(r.Get("port").(int)),
		r.Get("skip_ssl_verify").(bool))
	err := client.Start()

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to start a session to the vastdata cluser",
			Detail:   err.Error(),
		})
		return client, diags
	}
	cluster_version, _, version_get_error := client.ClusterVersion(ctx)
	if version_get_error != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error obtaning cluster version",
			Detail:   version_get_error.Error(),
		})

	}
	tflog.Info(ctx, fmt.Sprintf("Cluster version found %s", cluster_version))

	err = metadata.UpdateClusterVersion(cluster_version)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured while updating cluster version",
			Detail:   err.Error(),
		})
		return client, diags
	}
	if metadata.ClusterVersionCompare() != metadata.CLUSTER_VERSION_EQUALS {
		tflog.Warn(ctx, "Cluster Version & Build Version are not matching ,some actions might fail")
	}
	metadata.SetClusterConfig("version_validation_mode", r.Get("version_validation_mode").(string))

	return client, diags
}
