package provider

import (
	"context"
	"encoding/xml"
	"regexp"

	"github.com/PaloAltoNetworks/pango/xmlapi"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type VmAuthKeyCustom struct{}

func NewVmAuthKeyCustom(data *ProviderData) (*VmAuthKeyCustom, error) {
	return &VmAuthKeyCustom{}, nil
}

type vmAuthKeyRequest struct {
	XMLName  xml.Name `xml:"request"`
	Lifetime int64    `xml:"bootstrap>vm-auth-key>generate>lifetime"`
}

type vmAuthKeyResponse struct {
	XMLName xml.Name `xml:"response"`
	Result  string   `xml:"result"`
}

func (o *VmAuthKeyResource) OpenCustom(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var data VmAuthKeyResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	lifetime := data.Lifetime.ValueInt64()

	cmd := &xmlapi.Op{
		Command: vmAuthKeyRequest{Lifetime: lifetime},
	}

	var serverResponse vmAuthKeyResponse
	if _, _, err := o.client.Communicate(ctx, cmd, false, &serverResponse); err != nil {
		resp.Diagnostics.AddError("Failed to generate Authenticaion Key", "Server returned an error: "+err.Error())
		return
	}

	vmAuthKeyRegexp := `VM auth key (?P<vmauthkey>.+) generated. Expires at: (?P<expiration>.+)`
	expr := regexp.MustCompile(vmAuthKeyRegexp)
	match := expr.FindStringSubmatch(serverResponse.Result)
	if match == nil {
		resp.Diagnostics.AddError("Failed to parse server response", "Server response did not match regular expression")
		return
	}

	groups := make(map[string]string)
	for i, name := range expr.SubexpNames() {
		if i != 0 && name != "" {
			groups[name] = match[i]
		}
	}

	if vmAuthKey, found := groups["vmauthkey"]; found {
		data.VmAuthKey = types.StringValue(vmAuthKey)
	} else {
		resp.Diagnostics.AddError("Failed to parse server response", "Server response did not contain matching authentication key")
		return
	}

	if expiration, found := groups["expiration"]; found {
		data.ExpirationDate = types.StringValue(expiration)
	} else {
		resp.Diagnostics.AddWarning("Incomplete server response", "Server response didn't contain a valid expiration date")
	}

	resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)
}
