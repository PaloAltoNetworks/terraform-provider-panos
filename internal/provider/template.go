package provider

// Note:  This file is automatically generated.  Manually made changes
// will be overwritten when the provider is generated.

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/panorama/template"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	rsschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	sdkmanager "github.com/PaloAltoNetworks/terraform-provider-panos/internal/manager"
)

// Generate Terraform Data Source object.
var (
	_ datasource.DataSource              = &TemplateDataSource{}
	_ datasource.DataSourceWithConfigure = &TemplateDataSource{}
)

func NewTemplateDataSource() datasource.DataSource {
	return &TemplateDataSource{}
}

type TemplateDataSource struct {
	client  *pango.Client
	manager *sdkmanager.EntryObjectManager[*template.Entry, template.Location, *template.Service]
}

type TemplateDataSourceFilter struct {
	// TODO: Generate Data Source filter via function
}

type TemplateDataSourceModel struct {
	Location    TemplateLocation `tfsdk:"location"`
	Name        types.String     `tfsdk:"name"`
	Description types.String     `tfsdk:"description"`
}

func (o *TemplateDataSourceModel) CopyToPango(ctx context.Context, obj **template.Entry, encrypted *map[string]types.String) diag.Diagnostics {
	var diags diag.Diagnostics
	description_value := o.Description.ValueStringPointer()

	if (*obj) == nil {
		*obj = new(template.Entry)
	}
	(*obj).Name = o.Name.ValueString()
	(*obj).Description = description_value

	return diags
}

func (o *TemplateDataSourceModel) CopyFromPango(ctx context.Context, obj *template.Entry, encrypted *map[string]types.String) diag.Diagnostics {
	var diags diag.Diagnostics

	var description_value types.String
	if obj.Description != nil {
		description_value = types.StringValue(*obj.Description)
	}
	o.Name = types.StringValue(obj.Name)
	o.Description = description_value

	return diags
}

func TemplateDataSourceSchema() dsschema.Schema {
	return dsschema.Schema{
		Attributes: map[string]dsschema.Attribute{

			"location": TemplateDataSourceLocationSchema(),

			"name": dsschema.StringAttribute{
				Description: "The name of the service.",
				Computed:    false,
				Required:    true,
				Optional:    false,
				Sensitive:   false,
			},

			"description": dsschema.StringAttribute{
				Description: "The description.",
				Computed:    true,
				Required:    false,
				Optional:    true,
				Sensitive:   false,
			},
		},
	}
}

func (o *TemplateDataSourceModel) getTypeFor(name string) attr.Type {
	schema := TemplateDataSourceSchema()
	if attr, ok := schema.Attributes[name]; !ok {
		panic(fmt.Sprintf("could not resolve schema for attribute %s", name))
	} else {
		switch attr := attr.(type) {
		case dsschema.ListNestedAttribute:
			return attr.NestedObject.Type()
		case dsschema.MapNestedAttribute:
			return attr.NestedObject.Type()
		default:
			return attr.GetType()
		}
	}

	panic("unreachable")
}

func TemplateDataSourceLocationSchema() rsschema.Attribute {
	return TemplateLocationSchema()
}

// Metadata returns the data source type name.
func (d *TemplateDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_template"
}

// Schema defines the schema for this data source.
func (d *TemplateDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = TemplateDataSourceSchema()
}

// Configure prepares the struct.
func (d *TemplateDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerData := req.ProviderData.(*ProviderData)
	d.client = providerData.Client
	specifier, _, err := template.Versioning(d.client.Versioning())
	if err != nil {
		resp.Diagnostics.AddError("Failed to configure SDK client", err.Error())
		return
	}
	batchSize := providerData.MultiConfigBatchSize
	d.manager = sdkmanager.NewEntryObjectManager(d.client, template.NewService(d.client), batchSize, specifier, template.SpecMatches)
}
func (o *TemplateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var savestate, state TemplateDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &savestate)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var location template.Location

	if savestate.Location.Panorama != nil {
		location.Panorama = &template.PanoramaLocation{

			PanoramaDevice: savestate.Location.Panorama.PanoramaDevice.ValueString(),
		}
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource read", map[string]any{
		"resource_name": "panos_template_resource",
		"function":      "Read",
		"name":          savestate.Name.ValueString(),
	})

	// Perform the operation.
	object, err := o.manager.Read(ctx, location, savestate.Name.ValueString())
	if err != nil {
		if errors.Is(err, sdkmanager.ErrObjectNotFound) {
			resp.Diagnostics.AddError("Error reading data", err.Error())
		} else {
			resp.Diagnostics.AddError("Error reading entry", err.Error())
		}
		return
	}

	copy_diags := state.CopyFromPango(ctx, object, nil)
	resp.Diagnostics.Append(copy_diags...)

	/*
			// Keep the timeouts.
		    // TODO: This won't work for state import.
			state.Timeouts = savestate.Timeouts
	*/

	state.Location = savestate.Location

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

}

// Generate Terraform Resource object
var (
	_ resource.Resource                = &TemplateResource{}
	_ resource.ResourceWithConfigure   = &TemplateResource{}
	_ resource.ResourceWithImportState = &TemplateResource{}
)

func NewTemplateResource() resource.Resource {
	if _, found := resourceFuncMap["panos_template"]; !found {
		resourceFuncMap["panos_template"] = resourceFuncs{
			CreateImportId: TemplateImportStateCreator,
		}
	}
	return &TemplateResource{}
}

type TemplateResource struct {
	client  *pango.Client
	manager *sdkmanager.EntryObjectManager[*template.Entry, template.Location, *template.Service]
}

func TemplateResourceLocationSchema() rsschema.Attribute {
	return TemplateLocationSchema()
}

type TemplateResourceModel struct {
	Location    TemplateLocation `tfsdk:"location"`
	Name        types.String     `tfsdk:"name"`
	Description types.String     `tfsdk:"description"`
}

func (r *TemplateResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
}

// <ResourceSchema>

func TemplateResourceSchema() rsschema.Schema {
	return rsschema.Schema{
		Attributes: map[string]rsschema.Attribute{

			"location": TemplateResourceLocationSchema(),

			"name": rsschema.StringAttribute{
				Description: "The name of the service.",
				Computed:    false,
				Required:    true,
				Optional:    false,
				Sensitive:   false,
			},

			"description": rsschema.StringAttribute{
				Description: "The description.",
				Computed:    false,
				Required:    false,
				Optional:    true,
				Sensitive:   false,
			},
		},
	}
}

func (o *TemplateResourceModel) getTypeFor(name string) attr.Type {
	schema := TemplateResourceSchema()
	if attr, ok := schema.Attributes[name]; !ok {
		panic(fmt.Sprintf("could not resolve schema for attribute %s", name))
	} else {
		switch attr := attr.(type) {
		case rsschema.ListNestedAttribute:
			return attr.NestedObject.Type()
		case rsschema.MapNestedAttribute:
			return attr.NestedObject.Type()
		default:
			return attr.GetType()
		}
	}

	panic("unreachable")
}

func (r *TemplateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_template"
}

func (r *TemplateResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = TemplateResourceSchema()
}

// </ResourceSchema>

func (r *TemplateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	providerData := req.ProviderData.(*ProviderData)
	r.client = providerData.Client
	specifier, _, err := template.Versioning(r.client.Versioning())
	if err != nil {
		resp.Diagnostics.AddError("Failed to configure SDK client", err.Error())
		return
	}
	batchSize := providerData.MultiConfigBatchSize
	r.manager = sdkmanager.NewEntryObjectManager(r.client, template.NewService(r.client), batchSize, specifier, template.SpecMatches)
}

func (o *TemplateResourceModel) CopyToPango(ctx context.Context, obj **template.Entry, encrypted *map[string]types.String) diag.Diagnostics {
	var diags diag.Diagnostics
	description_value := o.Description.ValueStringPointer()

	if (*obj) == nil {
		*obj = new(template.Entry)
	}
	(*obj).Name = o.Name.ValueString()
	(*obj).Description = description_value

	return diags
}

func (o *TemplateResourceModel) CopyFromPango(ctx context.Context, obj *template.Entry, encrypted *map[string]types.String) diag.Diagnostics {
	var diags diag.Diagnostics

	var description_value types.String
	if obj.Description != nil {
		description_value = types.StringValue(*obj.Description)
	}
	o.Name = types.StringValue(obj.Name)
	o.Description = description_value

	return diags
}

func (r *TemplateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state TemplateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource create", map[string]any{
		"resource_name": "panos_template_resource",
		"function":      "Create",
		"name":          state.Name.ValueString(),
	})

	// Verify mode.
	if r.client.Hostname == "" {
		resp.Diagnostics.AddError("Invalid mode error", InspectionModeError)
		return
	}

	// Determine the location.

	var location template.Location

	if state.Location.Panorama != nil {
		location.Panorama = &template.PanoramaLocation{

			PanoramaDevice: state.Location.Panorama.PanoramaDevice.ValueString(),
		}
	}

	if err := location.IsValid(); err != nil {
		resp.Diagnostics.AddError("Invalid location", err.Error())
		return
	}

	// Load the desired config.
	var obj *template.Entry

	resp.Diagnostics.Append(state.CopyToPango(ctx, &obj, nil)...)
	if resp.Diagnostics.HasError() {
		return
	}

	/*
		// Timeout handling.
		ctx, cancel := context.WithTimeout(ctx, GetTimeout(state.Timeouts.Create))
		defer cancel()
	*/

	// Perform the operation.
	created, err := r.manager.Create(ctx, location, obj)
	if err != nil {
		resp.Diagnostics.AddError("Error in create", err.Error())
		return
	}

	resp.Diagnostics.Append(state.CopyFromPango(ctx, created, nil)...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.Name = types.StringValue(created.Name)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
func (o *TemplateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var savestate, state TemplateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &savestate)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var location template.Location

	if savestate.Location.Panorama != nil {
		location.Panorama = &template.PanoramaLocation{

			PanoramaDevice: savestate.Location.Panorama.PanoramaDevice.ValueString(),
		}
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource read", map[string]any{
		"resource_name": "panos_template_resource",
		"function":      "Read",
		"name":          savestate.Name.ValueString(),
	})

	// Perform the operation.
	object, err := o.manager.Read(ctx, location, savestate.Name.ValueString())
	if err != nil {
		if errors.Is(err, sdkmanager.ErrObjectNotFound) {
			resp.State.RemoveResource(ctx)
		} else {
			resp.Diagnostics.AddError("Error reading entry", err.Error())
		}
		return
	}

	copy_diags := state.CopyFromPango(ctx, object, nil)
	resp.Diagnostics.Append(copy_diags...)

	/*
			// Keep the timeouts.
		    // TODO: This won't work for state import.
			state.Timeouts = savestate.Timeouts
	*/

	state.Location = savestate.Location

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

}
func (r *TemplateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var plan, state TemplateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var location template.Location

	if state.Location.Panorama != nil {
		location.Panorama = &template.PanoramaLocation{

			PanoramaDevice: state.Location.Panorama.PanoramaDevice.ValueString(),
		}
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource update", map[string]any{
		"resource_name": "panos_template_resource",
		"function":      "Update",
	})

	// Verify mode.
	if r.client.Hostname == "" {
		resp.Diagnostics.AddError("Invalid mode error", InspectionModeError)
		return
	}
	obj, err := r.manager.Read(ctx, location, plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error in update", err.Error())
		return
	}

	resp.Diagnostics.Append(plan.CopyToPango(ctx, &obj, nil)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Perform the operation.
	updated, err := r.manager.Update(ctx, location, obj, obj.Name)
	if err != nil {
		resp.Diagnostics.AddError("Error in update", err.Error())
		return
	}

	// Save the location.
	state.Location = plan.Location

	/*
		// Keep the timeouts.
		state.Timeouts = plan.Timeouts
	*/

	copy_diags := state.CopyFromPango(ctx, updated, nil)
	resp.Diagnostics.Append(copy_diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

}
func (r *TemplateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var state TemplateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource delete", map[string]any{
		"resource_name": "panos_template_resource",
		"function":      "Delete",
		"name":          state.Name.ValueString(),
	})

	// Verify mode.
	if r.client.Hostname == "" {
		resp.Diagnostics.AddError("Invalid mode error", InspectionModeError)
		return
	}

	var location template.Location

	if state.Location.Panorama != nil {
		location.Panorama = &template.PanoramaLocation{

			PanoramaDevice: state.Location.Panorama.PanoramaDevice.ValueString(),
		}
	}

	err := r.manager.Delete(ctx, location, []string{state.Name.ValueString()})
	if err != nil && !errors.Is(err, sdkmanager.ErrObjectNotFound) {
		resp.Diagnostics.AddError("Error in delete", err.Error())
	}

}

type TemplateImportState struct {
	Location TemplateLocation `json:"location"`
	Name     string           `json:"name"`
}

func TemplateImportStateCreator(ctx context.Context, resource types.Object) ([]byte, error) {
	attrs := resource.Attributes()
	if attrs == nil {
		return nil, fmt.Errorf("Object has no attributes")
	}

	locationAttr, ok := attrs["location"]
	if !ok {
		return nil, fmt.Errorf("location attribute missing")
	}

	var location TemplateLocation
	switch value := locationAttr.(type) {
	case types.Object:
		value.As(ctx, &location, basetypes.ObjectAsOptions{})
	default:
		return nil, fmt.Errorf("location attribute expected to be an object")
	}
	nameAttr, ok := attrs["name"]
	if !ok {
		return nil, fmt.Errorf("name attribute missing")
	}

	var name string
	switch value := nameAttr.(type) {
	case types.String:
		name = value.ValueString()
	default:
		return nil, fmt.Errorf("name attribute expected to be a string")
	}

	importStruct := TemplateImportState{
		Location: location,
		Name:     name,
	}

	return json.Marshal(importStruct)
}

func (r *TemplateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	var obj TemplateImportState
	data, err := base64.StdEncoding.DecodeString(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Failed to decode Import ID", err.Error())
		return
	}

	err = json.Unmarshal(data, &obj)
	if err != nil {
		resp.Diagnostics.AddError("Failed to unmarshal Import ID", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("location"), obj.Location)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), obj.Name)...)
}

type TemplatePanoramaLocation struct {
	PanoramaDevice types.String `tfsdk:"panorama_device"`
}
type TemplateLocation struct {
	Panorama *TemplatePanoramaLocation `tfsdk:"panorama"`
}

func TemplateLocationSchema() rsschema.Attribute {
	return rsschema.SingleNestedAttribute{
		Description: "The location of this object.",
		Required:    true,
		Attributes: map[string]rsschema.Attribute{
			"panorama": rsschema.SingleNestedAttribute{
				Description: "Located in a specific Panorama.",
				Optional:    true,
				Attributes: map[string]rsschema.Attribute{
					"panorama_device": rsschema.StringAttribute{
						Description: "The Panorama device.",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString("localhost.localdomain"),
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (o TemplatePanoramaLocation) MarshalJSON() ([]byte, error) {
	obj := struct {
		PanoramaDevice *string `json:"panorama_device"`
	}{
		PanoramaDevice: o.PanoramaDevice.ValueStringPointer(),
	}

	return json.Marshal(obj)
}

func (o *TemplatePanoramaLocation) UnmarshalJSON(data []byte) error {
	var shadow struct {
		PanoramaDevice *string `json:"panorama_device"`
	}

	err := json.Unmarshal(data, &shadow)
	if err != nil {
		return err
	}
	o.PanoramaDevice = types.StringPointerValue(shadow.PanoramaDevice)

	return nil
}
func (o TemplateLocation) MarshalJSON() ([]byte, error) {
	obj := struct {
		Panorama *TemplatePanoramaLocation `json:"panorama"`
	}{
		Panorama: o.Panorama,
	}

	return json.Marshal(obj)
}

func (o *TemplateLocation) UnmarshalJSON(data []byte) error {
	var shadow struct {
		Panorama *TemplatePanoramaLocation `json:"panorama"`
	}

	err := json.Unmarshal(data, &shadow)
	if err != nil {
		return err
	}
	o.Panorama = shadow.Panorama

	return nil
}
