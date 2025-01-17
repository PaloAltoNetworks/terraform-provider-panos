package provider

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ function.Function = &ImportStateCreator{}
)

type ImportStateCreator struct{}

func NewCreateImportIdFunction() function.Function {
	return &ImportStateCreator{}
}

func (o *ImportStateCreator) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "generate_import_id"
}

func (o *ImportStateCreator) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "Generate Import ID",
		Description: "Generate Import ID for the given resource that can be used to import resources into the state.",

		Parameters: []function.Parameter{
			function.StringParameter{
				Name:        "resource_asn",
				Description: "Name of the resource",
			},
			function.DynamicParameter{
				Name:        "resource_data",
				Description: "Resource data",
			},
		},
		Return: function.StringReturn{},
	}
}

func (o *ImportStateCreator) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var resourceAsn string
	var dynamicResource types.Dynamic

	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &resourceAsn, &dynamicResource))
	if resp.Error != nil {
		return
	}

	var resource types.Object
	switch value := dynamicResource.UnderlyingValue().(type) {
	case types.Object:
		resource = value
	default:
		resp.Error = function.ConcatFuncErrors(resp.Error, function.NewArgumentFuncError(1, fmt.Sprintf("Wrong resource type: must be an object")))
		return
	}

	var data []byte

	if resourceFuncs, found := resourceFuncMap[resourceAsn]; !found {
		resp.Error = function.ConcatFuncErrors(resp.Error, function.NewArgumentFuncError(0, fmt.Sprintf("Unsupported resource type: %s'", resourceAsn)))
		return
	} else {
		var err error
		data, err = resourceFuncs.CreateImportId(ctx, resource)
		if err != nil {
			resp.Error = function.ConcatFuncErrors(resp.Error, function.NewFuncError(err.Error()))
			return
		}

	}

	result := base64.StdEncoding.EncodeToString(data)
	resp.Error = function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, result))
}
