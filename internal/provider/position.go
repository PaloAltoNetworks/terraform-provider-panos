package provider

import (
	"encoding/json"
	"fmt"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/attr"
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

func (o *TerraformPositionObject) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"where":    types.StringType,
		"pivot":    types.StringType,
		"directly": types.BoolType,
	}
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

	var where string
	if !o.Where.IsUnknown() {
		where = o.Where.ValueString()
		if !slices.Contains(allowedPositions, o.Where.ValueString()) {
			resp.Diagnostics.AddAttributeError(
				path.Root("position").AtName("where"),
				"Missing attribute configuration",
				fmt.Sprintf("where attribute must be one of the valid values: first, last, before, after, found: '%s'", o.Where.ValueString()))
			return
		}
	}

	// where is either a valid position, or an empty string at this point. If where position requires a valid pivot point, and
	// o.Pivot is known at this time, validate that o.Pivot is neither null nor an empty string.
	if (where == "after" || where == "before") && !o.Pivot.IsUnknown() && (o.Pivot.IsNull() || o.Pivot.ValueString() == "") {
		resp.Diagnostics.AddAttributeError(
			path.Root("position").AtName("pivot"),
			"Missing attribute configuration",
			"position pivot attribute must be set to a valid object when where attribute is set to either 'after' or 'before'")
		return
	}

	if where == "first" || where == "last" {
		if !o.Pivot.IsUnknown() && !o.Pivot.IsNull() {
			resp.Diagnostics.AddAttributeWarning(
				path.Root("position").AtName("pivot"),
				"Unexpected attribute configuration",
				"pivot attribute is ignored when where is set to 'first' or 'last'")
		}

		if !o.Directly.IsUnknown() && !o.Directly.IsNull() {
			resp.Diagnostics.AddAttributeWarning(
				path.Root("position").AtName("directly"),
				"Unexpected attribute configuration",
				"directly attribute is ignored when where is set to 'first' or 'last'")
		}
	}

	// If either pivot or direclty are unknown we can't validate that they are both set correcly.
	if o.Pivot.IsUnknown() || o.Directly.IsUnknown() {
		return
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
