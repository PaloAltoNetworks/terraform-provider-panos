package provider

import (
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	rsschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/PaloAltoNetworks/pango/rule"
)

type TerraformPositionObject struct {
	Where    types.String `tfsdk:"where"`
	Pivot    types.String `tfsdk:"pivot"`
	Directly types.Bool   `tfsdk:"directly"`
}

func TerraformPositionObjectSchema() rsschema.SingleNestedAttribute {
	return rsschema.SingleNestedAttribute{
		Required: true,
		Attributes: map[string]rsschema.Attribute{
			"where": rsschema.StringAttribute{
				Required: true,
			},
			"pivot": rsschema.StringAttribute{
				Optional: true,
			},
			"directly": rsschema.BoolAttribute{
				Optional: true,
			},
		},
	}
}

func (o *TerraformPositionObject) CopyToPango() rule.Position {
	trueVal := true
	switch o.Where.ValueString() {
	case "first":
		return rule.Position{
			First: &trueVal,
		}
	case "last":
		return rule.Position{
			Last: &trueVal,
		}
	case "before":
		if o.Directly.ValueBool() == true {
			return rule.Position{
				DirectlyBefore: o.Pivot.ValueStringPointer(),
			}
		} else {
			return rule.Position{
				SomewhereBefore: o.Pivot.ValueStringPointer(),
			}
		}
	case "after":
		if o.Directly.ValueBool() == true {
			return rule.Position{
				DirectlyAfter: o.Pivot.ValueStringPointer(),
			}
		} else {
			return rule.Position{
				SomewhereAfter: o.Pivot.ValueStringPointer(),
			}
		}
	default:
		panic("unreachable")
	}
}

func (o *TerraformPositionObject) ValidateConfig(resp *resource.ValidateConfigResponse) {
	allowedPositions := []string{"first", "last", "before", "after"}

	if !slices.Contains(allowedPositions, o.Where.ValueString()) {
		resp.Diagnostics.AddAttributeError(
			path.Root("position").AtName("directly"),
			"Missing attribute configuration",
			"where attribute must be one of the valid values: first, last, before, after")
	}

	if !o.Pivot.IsNull() && o.Directly.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("position").AtName("directly"),
			"Missing attribute configuration",
			"Expected directly to be configured with pivot")
	}

	if o.Pivot.IsNull() && !o.Directly.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("position").AtName("pivot"),
			"Missing attribute configuration",
			"Expected pivot to be configured with directly")
	}
}
