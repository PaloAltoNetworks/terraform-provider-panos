package provider

import (
	"encoding/json"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	rsschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/PaloAltoNetworks/pango/movement"
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

func (o *TerraformPositionObject) CopyToPango() movement.Position {
	switch o.Where.ValueString() {
	case "first":
		return movement.PositionFirst{}
	case "last":
		return movement.PositionLast{}
	case "before":
		return movement.PositionBefore{
			Pivot:    o.Pivot.ValueString(),
			Directly: o.Directly.ValueBool(),
		}
	case "after":
		return movement.PositionAfter{
			Pivot:    o.Pivot.ValueString(),
			Directly: o.Directly.ValueBool(),
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

func (o TerraformPositionObject) MarshalJSON() ([]byte, error) {
	obj := struct {
		Where    *string `json:"where,omitempty"`
		Directly *bool   `json:"directly,omitempty"`
		Pivot    *string `json:"pivot,omitempty"`
	}{
		Where:    o.Where.ValueStringPointer(),
		Directly: o.Directly.ValueBoolPointer(),
		Pivot:    o.Pivot.ValueStringPointer(),
	}

	return json.Marshal(obj)
}

func (o *TerraformPositionObject) UnmarshalJSON(data []byte) error {
	var shadow struct {
		Where    *string `json:"where"`
		Directly *bool   `json:"directly"`
		Pivot    *string `json:"pivot"`
	}

	err := json.Unmarshal(data, &shadow)
	if err != nil {
		return err
	}

	o.Where = types.StringPointerValue(shadow.Where)
	o.Directly = types.BoolPointerValue(shadow.Directly)
	o.Pivot = types.StringPointerValue(shadow.Pivot)

	return nil
}
