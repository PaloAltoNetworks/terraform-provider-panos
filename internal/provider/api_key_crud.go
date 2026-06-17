package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (o *ApiKeyResource) ImportStateCustom(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.AddError("Import not supported", "The panos_api_key resource does not support terraform import.")
}

type ApiKeyCustom struct{}

func NewApiKeyCustom(provider *ProviderData) (*ApiKeyCustom, error) {
	return &ApiKeyCustom{}, nil
}

func (o *ApiKeyResource) OpenCustom(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var data ApiKeyResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	username := data.Username.ValueString()
	password := data.Password.ValueString()

	apiKey, err := o.client.GenerateApiKey(ctx, username, password)
	if err != nil {
		resp.Diagnostics.AddError("failed to generate API key", err.Error())
		return
	}

	data.ApiKey = types.StringValue(apiKey)
	resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)
}
