package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ function.Function = &AddressValueFunction{}

type AddressValueFunction struct{}

func NewAddressValueFunction() function.Function {
	return &AddressValueFunction{}
}

func (f *AddressValueFunction) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "address_value"
}

func (f *AddressValueFunction) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "Return value of a given address resource",
		Description: "Given an address object resource, return its value.",

		Parameters: []function.Parameter{
			function.ObjectParameter{
				Name:        "address",
				Description: "address resource to get value from",
				AttributeTypes: map[string]attr.Type{
					"ip_netmask":  types.StringType,
					"ip_range":    types.StringType,
					"ip_wildcard": types.StringType,
					"fqdn":        types.StringType,
				},
			},
		},
		Return: function.StringReturn{},
	}
}

func (f *AddressValueFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var address struct {
		IpNetmask  *string `tfsdk:"ip_netmask"`
		IpRange    *string `tfsdk:"ip_range"`
		IpWildcard *string `tfsdk:"ip_wildcard"`
		Fqdn       *string `tfsdk:"fqdn"`
	}

	// Read Terraform argument data into the variable
	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &address))
	if resp.Error != nil {
		return
	}

	var value string
	if address.IpNetmask != nil {
		value = *address.IpNetmask
	} else if address.IpRange != nil {
		value = *address.IpRange
	} else if address.IpWildcard != nil {
		value = *address.IpWildcard
	} else if address.Fqdn != nil {
		value = *address.Fqdn
	} else {
		resp.Error = function.ConcatFuncErrors(resp.Error, function.NewFuncError("given address has no value set"))
		return
	}

	// Set the result to the same data
	resp.Error = function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, value))
}
