package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

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
