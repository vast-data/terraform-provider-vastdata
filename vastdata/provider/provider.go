// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"os"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	vsd "github.com/vast-data/terraform-provider-vastdata/vastdata"
	"github.com/vast-data/terraform-provider-vastdata/vastdata/client"
)

var _ provider.Provider = &VastProvider{}

type VastProvider struct {
	version string
}

type VastProviderModel struct {
	Host                  types.String `tfsdk:"host"`
	Port                  types.Int64  `tfsdk:"port"`
	SkipSSLVerify         types.Bool   `tfsdk:"skip_ssl_verify"`
	Username              types.String `tfsdk:"username"`
	Password              types.String `tfsdk:"password"`
	ApiToken              types.String `tfsdk:"api_token"`
	Tenant                types.String `tfsdk:"tenant"`
	VersionValidationMode types.String `tfsdk:"version_validation_mode"`
}

func New(
	version string,
) func() provider.Provider {
	return func() provider.Provider {
		return &VastProvider{
			version: version,
		}
	}
}

func (p *VastProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "vastdata"
	resp.Version = p.version
}

func (p *VastProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The VastData Cluster hostname/address , if environment variable VASTDATA_HOST exists it will be used",
			},
			"port": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "The server API port (Default is 443) ,if environment variable VASTDATA_PORT exists it will be used",
			},
			"skip_ssl_verify": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Whether to skip SSL certificate verification.",
			},
			"username": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "VastData Cluster username (conflicts with api_token).",
			},
			"password": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "VastData Cluster password (conflicts with api_token).",
			},
			"api_token": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "VastData Cluster API token (conflicts with username/password).",
			},
			"tenant": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Tenant name for tenant-scoped authentication (tenant admin). If environment variable VASTDATA_TENANT exists it will be used. This is passed via the X-Tenant-Name header during login.",
			},
			"version_validation_mode": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Version validation mode: 'strict' or 'warn'.",
			},
		},
	}
}

func (p *VastProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config VastProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Fail on unknowns
	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown Host",
			"The provider cannot create the VastData client because host is unknown.",
		)
	}

	// Environment variable fallback
	host := getenvOr(config.Host, "VASTDATA_HOST")
	port := int64Or(config.Port, "VASTDATA_PORT", 443)
	skipSSL := boolOr(config.SkipSSLVerify, "VASTDATA_VERIFY_SSL", false)
	username := getenvOr(config.Username, "VASTDATA_CLUSTER_USERNAME")
	password := getenvOr(config.Password, "VASTDATA_CLUSTER_PASSWORD")
	apiToken := getenvOr(config.ApiToken, "VASTDATA_API_TOKEN")
	tenant := getenvOr(config.Tenant, "VASTDATA_TENANT")
	validationMode := getenvOr(config.VersionValidationMode, "VERSION_VALIDATION_MODE")
	if validationMode == "" {
		validationMode = "warn"
	}

	// Generic timeout. Should be enough for all API operations.
	restTimeout := time.Minute * 4
	vmsRest, err := client.NewRest(host, port, username, password, apiToken, tenant, !skipSSL, p.version, restTimeout)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create VAST API Client",
			"An unexpected error occurred when creating the VAST API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"VAST Client Error: "+err.Error(),
		)
		return
	}

	migrateMode := os.Getenv("VASTDATA_MIGRATE_MODE")
	isMigrateMode := migrateMode == "1" || migrateMode == "true" || migrateMode == "TRUE"
	if isMigrateMode {
		tflog.Warn(ctx, "╔══════════════════════════════════════════════════════════════════╗")
		tflog.Warn(ctx, "║  VASTDATA MIGRATE MODE ENABLED                                  ║")
		tflog.Warn(ctx, "║  Create operations will READ existing resources from the cluster ║")
		tflog.Warn(ctx, "║  Update and Delete operations will be BLOCKED (empty state only) ║")
		tflog.Warn(ctx, "║  Use this mode to populate state from existing infrastructure    ║")
		tflog.Warn(ctx, "╚══════════════════════════════════════════════════════════════════╝")
	}

	providerData := &vsd.ProviderData{
		Client:      vmsRest,
		MigrateMode: isMigrateMode,
	}
	resp.ResourceData = providerData
	resp.DataSourceData = providerData
}

func (p *VastProvider) Resources(_ context.Context) []func() resource.Resource {
	return vsd.GetResourceFactories()
}

func (p *VastProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return vsd.GetDatasourceFactories()
}
