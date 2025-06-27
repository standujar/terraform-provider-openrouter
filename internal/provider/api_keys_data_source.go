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

var _ datasource.DataSource = &ApiKeysDataSource{}

func NewApiKeysDataSource() datasource.DataSource {
	return &ApiKeysDataSource{}
}

type ApiKeysDataSource struct {
	client *client.Client
}

type ApiKeysDataSourceModel struct {
	IncludeDisabled types.Bool `tfsdk:"include_disabled"`
	Keys            []ApiKeyModel `tfsdk:"keys"`
}

type ApiKeyModel struct {
	ID             types.String  `tfsdk:"id"`
	Name           types.String  `tfsdk:"name"`
	IsProvisioner  types.Bool    `tfsdk:"is_provisioner"`
	Limit          types.Float64 `tfsdk:"limit"`
	LimitMinutes   types.Int64   `tfsdk:"limit_minutes"`
	Usage          types.Float64 `tfsdk:"usage"`
	IsDisabled     types.Bool    `tfsdk:"is_disabled"`
	CreatedAt      types.String  `tfsdk:"created_at"`
}

func (d *ApiKeysDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_api_keys"
}

func (d *ApiKeysDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves a list of OpenRouter API keys.",

		Attributes: map[string]schema.Attribute{
			"include_disabled": schema.BoolAttribute{
				MarkdownDescription: "Whether to include disabled API keys in the list.",
				Optional:            true,
			},
			"keys": schema.ListNestedAttribute{
				MarkdownDescription: "List of API keys.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "The hash identifier of the API key.",
							Computed:            true,
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
				},
			},
		},
	}
}

func (d *ApiKeysDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ApiKeysDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ApiKeysDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "reading API keys data source")

	params := &client.ListApiKeysRequest{
		IncludeDisabled: data.IncludeDisabled.ValueBool(),
	}

	apiKeys, err := d.client.ListApiKeys(ctx, params)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read API keys, got error: %s", err))
		return
	}

	data.Keys = make([]ApiKeyModel, len(apiKeys))
	for i, apiKey := range apiKeys {
		key := ApiKeyModel{
			ID:            types.StringValue(apiKey.ID),
			Name:          types.StringValue(apiKey.Name),
			IsProvisioner: types.BoolValue(apiKey.IsProvisioner),
			Usage:         types.Float64Value(apiKey.Usage),
			IsDisabled:    types.BoolValue(apiKey.IsDisabled),
		}

		if apiKey.Limit != nil {
			key.Limit = types.Float64Value(*apiKey.Limit)
		} else {
			key.Limit = types.Float64Null()
		}

		if apiKey.LimitMinutes != nil {
			key.LimitMinutes = types.Int64Value(int64(*apiKey.LimitMinutes))
		} else {
			key.LimitMinutes = types.Int64Null()
		}

		if apiKey.CreatedAt != nil {
			key.CreatedAt = types.StringValue(apiKey.CreatedAt.Format("2006-01-02T15:04:05Z"))
		} else {
			key.CreatedAt = types.StringNull()
		}

		data.Keys[i] = key
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}