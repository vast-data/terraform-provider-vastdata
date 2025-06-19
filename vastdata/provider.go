package vastdata

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	vast_client "github.com/vast-data/go-vast-client"
	"os"
	"strconv"
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
	VersionValidationMode types.String `tfsdk:"version_validation_mode"`
}

func New(version string) func() provider.Provider {
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
	validationMode := getenvOr(config.VersionValidationMode, "VERSION_VALIDATION_MODE")
	if validationMode == "" {
		validationMode = "warn"
	}

	vmsConfig := &vast_client.VMSConfig{
		Host:      host,
		Port:      uint64(port),
		Username:  username,
		Password:  password,
		ApiToken:  apiToken,
		SslVerify: !skipSSL,
		UserAgent: "Terraform Provider VastData/" + p.version,

		BeforeRequestFn: BeforeRequestFnCallback,
		AfterRequestFn:  AfterRequestFnCallback,
	}

	client, err := vast_client.NewVMSRest(vmsConfig)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create VAST API Client",
			"An unexpected error occurred when creating the VAST API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"VAST Client Error: "+err.Error(),
		)
		return
	}

	resp.ResourceData = client
	resp.DataSourceData = client
}

func (p *VastProvider) Resources(ctx context.Context) []func() resource.Resource {
	var resourceFactories = []func() ResourceState{
		func() ResourceState { return &User{} },
		func() ResourceState { return &Group{} },
	}

	resources := make([]func() resource.Resource, 0, len(resourceFactories))
	for _, factory := range resourceFactories {
		resources = append(resources, RegisterResourceFor(factory))
	}
	return resources
}

func (p *VastProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	var dataSourceFactories = []func() DataSourceState{
		func() DataSourceState { return &User{} },
		func() DataSourceState { return &Group{} },
	}

	dataSources := make([]func() datasource.DataSource, 0, len(dataSourceFactories))
	for _, factory := range dataSourceFactories {
		dataSources = append(dataSources, RegisterDataSourceFor(factory))
	}
	return dataSources
}

func getenvOr(val types.String, envKey string) string {
	if !val.IsNull() {
		return val.ValueString()
	}
	if v := os.Getenv(envKey); v != "" {
		return v
	}
	return ""
}

func int64Or(val types.Int64, envKey string, def int64) int64 {
	if !val.IsNull() {
		return val.ValueInt64()
	}
	if v := os.Getenv(envKey); v != "" {
		if parsed, err := strconv.ParseInt(v, 10, 64); err == nil {
			return parsed
		}
	}
	return def
}

func boolOr(val types.Bool, envKey string, def bool) bool {
	if !val.IsNull() {
		return val.ValueBool()
	}
	if v := os.Getenv(envKey); v != "" {
		return v == "1" || v == "true" || v == "TRUE"
	}
	return def
}
