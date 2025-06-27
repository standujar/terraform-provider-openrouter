package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/standujar/terraform-provider-openrouter/internal/client"
)

var _ datasource.DataSource = &ApiKeyDataSource{}

func NewApiKeyDataSource() datasource.DataSource {
	return &ApiKeyDataSource{}
}

type ApiKeyDataSource struct {
	client *client.Client
}

type ApiKeyDataSourceModel struct {
	ID             types.String  `tfsdk:"id"`
	Name           types.String  `tfsdk:"name"`
	IsProvisioner  types.Bool    `tfsdk:"is_provisioner"`
	Limit          types.Float64 `tfsdk:"limit"`
	LimitMinutes   types.Int64   `tfsdk:"limit_minutes"`
	Usage          types.Float64 `tfsdk:"usage"`
	IsDisabled     types.Bool    `tfsdk:"is_disabled"`
	CreatedAt      types.String  `tfsdk:"created_at"`
}

func (d *ApiKeyDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_api_key"
}

func (d *ApiKeyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves information about an OpenRouter API key.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The hash identifier of the API key.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the API key.",
				Computed:            true,
			},
			"is_provisioner": schema.BoolAttribute{
				MarkdownDescription: "Whether the API key is a provisioner key.",
				Computed:            true,
			},
			"limit": schema.Float64Attribute{
				MarkdownDescription: "The spend limit for the API key in USD.",
				Computed:            true,
			},
			"limit_minutes": schema.Int64Attribute{
				MarkdownDescription: "The time limit for the API key in minutes.",
				Computed:            true,
			},
			"usage": schema.Float64Attribute{
				MarkdownDescription: "The current usage of the API key in USD.",
				Computed:            true,
			},
			"is_disabled": schema.BoolAttribute{
				MarkdownDescription: "Whether the API key is disabled.",
				Computed:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "The creation timestamp of the API key.",
				Computed:            true,
			},
		},
	}
}

func (d *ApiKeyDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *ApiKeyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ApiKeyDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "reading API key data source", map[string]interface{}{
		"id": data.ID.ValueString(),
	})

	apiKey, err := d.client.GetApiKey(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read API key, got error: %s", err))
		return
	}

	data.Name = types.StringValue(apiKey.Name)
	data.IsProvisioner = types.BoolValue(apiKey.IsProvisioner)
	data.Usage = types.Float64Value(apiKey.Usage)
	data.IsDisabled = types.BoolValue(apiKey.IsDisabled)
	
	if apiKey.Limit != nil {
		data.Limit = types.Float64Value(*apiKey.Limit)
	} else {
		data.Limit = types.Float64Null()
	}

	if apiKey.LimitMinutes != nil {
		data.LimitMinutes = types.Int64Value(int64(*apiKey.LimitMinutes))
	} else {
		data.LimitMinutes = types.Int64Null()
	}

	if apiKey.CreatedAt != nil {
		data.CreatedAt = types.StringValue(apiKey.CreatedAt.Format("2006-01-02T15:04:05Z"))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}