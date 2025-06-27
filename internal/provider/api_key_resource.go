package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/standujar/terraform-provider-openrouter/internal/client"
)

var _ resource.Resource = &ApiKeyResource{}
var _ resource.ResourceWithImportState = &ApiKeyResource{}

func NewApiKeyResource() resource.Resource {
	return &ApiKeyResource{}
}

type ApiKeyResource struct {
	client *client.Client
}

type ApiKeyResourceModel struct {
	ID           types.String  `tfsdk:"id"`
	Key          types.String  `tfsdk:"key"`
	Name         types.String  `tfsdk:"name"`
	Limit        types.Float64 `tfsdk:"limit"`
	LimitMinutes types.Int64   `tfsdk:"limit_minutes"`
	IsDisabled   types.Bool    `tfsdk:"is_disabled"`
	Usage        types.Float64 `tfsdk:"usage"`
	CreatedAt    types.String  `tfsdk:"created_at"`
}

func (r *ApiKeyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_api_key"
}

func (r *ApiKeyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an OpenRouter API key.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The hash identifier of the API key.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"key": schema.StringAttribute{
				MarkdownDescription: "The API key value. Only available during creation.",
				Computed:            true,
				Sensitive:           true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the API key.",
				Required:            true,
			},
			"limit": schema.Float64Attribute{
				MarkdownDescription: "The spend limit for the API key in USD.",
				Optional:            true,
			},
			"limit_minutes": schema.Int64Attribute{
				MarkdownDescription: "The time limit for the API key in minutes.",
				Optional:            true,
			},
			"is_disabled": schema.BoolAttribute{
				MarkdownDescription: "Whether the API key is disabled.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"usage": schema.Float64Attribute{
				MarkdownDescription: "The current usage of the API key in USD.",
				Computed:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "The creation timestamp of the API key.",
				Computed:            true,
			},
		},
	}
}

func (r *ApiKeyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *ApiKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ApiKeyResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "creating API key")

	createReq := &client.CreateApiKeyRequest{
		Name: data.Name.ValueString(),
	}

	if !data.Limit.IsNull() {
		limit := data.Limit.ValueFloat64()
		createReq.Limit = &limit
	}

	if !data.LimitMinutes.IsNull() {
		limitMinutes := int(data.LimitMinutes.ValueInt64())
		createReq.LimitMinutes = &limitMinutes
	}

	apiKey, err := r.client.CreateApiKey(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create API key, got error: %s", err))
		return
	}

	data.ID = types.StringValue(apiKey.Data.ID)
	data.Key = types.StringValue(apiKey.Key)
	data.Usage = types.Float64Value(apiKey.Data.Usage)
	data.IsDisabled = types.BoolValue(apiKey.Data.IsDisabled)
	
	if apiKey.Data.CreatedAt != nil {
		data.CreatedAt = types.StringValue(apiKey.Data.CreatedAt.Format("2006-01-02T15:04:05Z"))
	}

	tflog.Trace(ctx, "created API key")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ApiKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ApiKeyResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "reading API key", map[string]interface{}{
		"id": data.ID.ValueString(),
	})

	apiKey, err := r.client.GetApiKey(ctx, data.ID.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read API key, got error: %s", err))
		return
	}

	data.Name = types.StringValue(apiKey.Name)
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

func (r *ApiKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ApiKeyResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "updating API key", map[string]interface{}{
		"id": data.ID.ValueString(),
	})

	updateReq := &client.UpdateApiKeyRequest{}

	var state ApiKeyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !data.Name.Equal(state.Name) {
		name := data.Name.ValueString()
		updateReq.Name = &name
	}

	if !data.Limit.Equal(state.Limit) {
		if !data.Limit.IsNull() {
			limit := data.Limit.ValueFloat64()
			updateReq.Limit = &limit
		} else {
			var zeroLimit float64 = 0
			updateReq.Limit = &zeroLimit
		}
	}

	if !data.IsDisabled.Equal(state.IsDisabled) {
		isDisabled := data.IsDisabled.ValueBool()
		updateReq.IsDisabled = &isDisabled
	}

	apiKey, err := r.client.UpdateApiKey(ctx, data.ID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update API key, got error: %s", err))
		return
	}

	data.Usage = types.Float64Value(apiKey.Usage)

	tflog.Trace(ctx, "updated API key")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ApiKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ApiKeyResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "deleting API key", map[string]interface{}{
		"id": data.ID.ValueString(),
	})

	err := r.client.DeleteApiKey(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete API key, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "deleted API key")
}

func (r *ApiKeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}