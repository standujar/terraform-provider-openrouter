package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/standujar/terraform-provider-openrouter/internal/client"
)

var _ provider.Provider = &OpenRouterProvider{}

type OpenRouterProvider struct {
	version string
}

type OpenRouterProviderModel struct {
	ApiKey   types.String `tfsdk:"api_key"`
	Endpoint types.String `tfsdk:"endpoint"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &OpenRouterProvider{
			version: version,
		}
	}
}

func (p *OpenRouterProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "openrouter"
	resp.Version = p.version
}

func (p *OpenRouterProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				MarkdownDescription: "API Key for OpenRouter. Can also be set via OPENROUTER_API_KEY environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "API endpoint for OpenRouter. Defaults to https://openrouter.ai/api/v1.",
				Optional:            true,
			},
		},
	}
}

func (p *OpenRouterProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring OpenRouter client")

	var config OpenRouterProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if config.ApiKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Unknown OpenRouter API Key",
			"The provider cannot create the OpenRouter client as there is an unknown configuration value for the OpenRouter API key. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the OPENROUTER_API_KEY environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	apiKey := os.Getenv("OPENROUTER_API_KEY")
	endpoint := "https://openrouter.ai/api/v1"

	if !config.ApiKey.IsNull() {
		apiKey = config.ApiKey.ValueString()
	}

	if !config.Endpoint.IsNull() {
		endpoint = config.Endpoint.ValueString()
	}

	if apiKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Missing OpenRouter API Key",
			"The provider cannot create the OpenRouter client as there is a missing or empty value for the OpenRouter API key. "+
				"Set the api_key value in the configuration or use the OPENROUTER_API_KEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating OpenRouter client")

	client := client.NewClient(apiKey, &endpoint)

	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured OpenRouter client", map[string]any{"endpoint": endpoint})
}

func (p *OpenRouterProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewApiKeyResource,
	}
}

func (p *OpenRouterProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewApiKeyDataSource,
		NewApiKeysDataSource,
	}
}