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
	"github.com/PaloAltoNetworks/pango/objects/profiles/customurlcategory"

	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	rsschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	sdkmanager "github.com/PaloAltoNetworks/terraform-provider-panos/internal/manager"
)

// Generate Terraform Data Source object.
var (
	_ datasource.DataSource              = &CustomUrlCategoryDataSource{}
	_ datasource.DataSourceWithConfigure = &CustomUrlCategoryDataSource{}
)

func NewCustomUrlCategoryDataSource() datasource.DataSource {
	return &CustomUrlCategoryDataSource{}
}

type CustomUrlCategoryDataSource struct {
	client  *pango.Client
	manager *sdkmanager.EntryObjectManager[*customurlcategory.Entry, customurlcategory.Location, *customurlcategory.Service]
}

type CustomUrlCategoryDataSourceFilter struct {
	// TODO: Generate Data Source filter via function
}

type CustomUrlCategoryDataSourceModel struct {
	Location        CustomUrlCategoryLocation `tfsdk:"location"`
	Name            types.String              `tfsdk:"name"`
	Description     types.String              `tfsdk:"description"`
	DisableOverride types.String              `tfsdk:"disable_override"`
	List            types.List                `tfsdk:"list"`
	Type            types.String              `tfsdk:"type"`
}

func (o *CustomUrlCategoryDataSourceModel) CopyToPango(ctx context.Context, obj **customurlcategory.Entry, encrypted *map[string]types.String) diag.Diagnostics {
	var diags diag.Diagnostics
	type_value := o.Type.ValueStringPointer()
	description_value := o.Description.ValueStringPointer()
	disableOverride_value := o.DisableOverride.ValueStringPointer()
	list_pango_entries := make([]string, 0)
	diags.Append(o.List.ElementsAs(ctx, &list_pango_entries, false)...)
	if diags.HasError() {
		return diags
	}

	if (*obj) == nil {
		*obj = new(customurlcategory.Entry)
	}
	(*obj).Name = o.Name.ValueString()
	(*obj).Type = type_value
	(*obj).Description = description_value
	(*obj).DisableOverride = disableOverride_value
	(*obj).List = list_pango_entries

	return diags
}

func (o *CustomUrlCategoryDataSourceModel) CopyFromPango(ctx context.Context, obj *customurlcategory.Entry, encrypted *map[string]types.String) diag.Diagnostics {
	var diags diag.Diagnostics
	var list_list types.List
	{
		var list_diags diag.Diagnostics
		list_list, list_diags = types.ListValueFrom(ctx, types.StringType, obj.List)
		diags.Append(list_diags...)
	}

	var disableOverride_value types.String
	if obj.DisableOverride != nil {
		disableOverride_value = types.StringValue(*obj.DisableOverride)
	}
	var type_value types.String
	if obj.Type != nil {
		type_value = types.StringValue(*obj.Type)
	}
	var description_value types.String
	if obj.Description != nil {
		description_value = types.StringValue(*obj.Description)
	}
	o.Name = types.StringValue(obj.Name)
	o.DisableOverride = disableOverride_value
	o.List = list_list
	o.Type = type_value
	o.Description = description_value

	return diags
}

func CustomUrlCategoryDataSourceSchema() dsschema.Schema {
	return dsschema.Schema{
		Attributes: map[string]dsschema.Attribute{

			"location": CustomUrlCategoryDataSourceLocationSchema(),

			"name": dsschema.StringAttribute{
				Description: "",
				Computed:    false,
				Required:    true,
				Optional:    false,
				Sensitive:   false,
			},

			"description": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
				Required:    false,
				Optional:    true,
				Sensitive:   false,
			},

			"disable_override": dsschema.StringAttribute{
				Description: "disable object override in child device groups",
				Computed:    true,
				Required:    false,
				Optional:    true,
				Sensitive:   false,
			},

			"list": dsschema.ListAttribute{
				Description: "",
				Required:    false,
				Optional:    true,
				Computed:    true,
				Sensitive:   false,
				ElementType: types.StringType,
			},

			"type": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
				Required:    false,
				Optional:    true,
				Sensitive:   false,
			},
		},
	}
}

func (o *CustomUrlCategoryDataSourceModel) getTypeFor(name string) attr.Type {
	schema := CustomUrlCategoryDataSourceSchema()
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

func CustomUrlCategoryDataSourceLocationSchema() rsschema.Attribute {
	return CustomUrlCategoryLocationSchema()
}

// Metadata returns the data source type name.
func (d *CustomUrlCategoryDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_custom_url_category"
}

// Schema defines the schema for this data source.
func (d *CustomUrlCategoryDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = CustomUrlCategoryDataSourceSchema()
}

// Configure prepares the struct.
func (d *CustomUrlCategoryDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*pango.Client)
	specifier, _, err := customurlcategory.Versioning(d.client.Versioning())
	if err != nil {
		resp.Diagnostics.AddError("Failed to configure SDK client", err.Error())
		return
	}
	d.manager = sdkmanager.NewEntryObjectManager(d.client, customurlcategory.NewService(d.client), specifier, customurlcategory.SpecMatches)
}

func (o *CustomUrlCategoryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var savestate, state CustomUrlCategoryDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &savestate)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var location customurlcategory.Location

	if !savestate.Location.Shared.IsNull() && savestate.Location.Shared.ValueBool() {
		location.Shared = true
	}
	if savestate.Location.Vsys != nil {
		location.Vsys = &customurlcategory.VsysLocation{

			NgfwDevice: savestate.Location.Vsys.NgfwDevice.ValueString(),
			Vsys:       savestate.Location.Vsys.Name.ValueString(),
		}
	}
	if savestate.Location.DeviceGroup != nil {
		location.DeviceGroup = &customurlcategory.DeviceGroupLocation{

			PanoramaDevice: savestate.Location.DeviceGroup.PanoramaDevice.ValueString(),
			DeviceGroup:    savestate.Location.DeviceGroup.Name.ValueString(),
		}
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource read", map[string]any{
		"resource_name": "panos_custom_url_category_resource",
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
	_ resource.Resource                = &CustomUrlCategoryResource{}
	_ resource.ResourceWithConfigure   = &CustomUrlCategoryResource{}
	_ resource.ResourceWithImportState = &CustomUrlCategoryResource{}
)

func NewCustomUrlCategoryResource() resource.Resource {
	if _, found := resourceFuncMap["panos_custom_url_category"]; !found {
		resourceFuncMap["panos_custom_url_category"] = resourceFuncs{
			CreateImportId: CustomUrlCategoryImportStateCreator,
		}
	}
	return &CustomUrlCategoryResource{}
}

type CustomUrlCategoryResource struct {
	client  *pango.Client
	manager *sdkmanager.EntryObjectManager[*customurlcategory.Entry, customurlcategory.Location, *customurlcategory.Service]
}

func CustomUrlCategoryResourceLocationSchema() rsschema.Attribute {
	return CustomUrlCategoryLocationSchema()
}

type CustomUrlCategoryResourceModel struct {
	Location        CustomUrlCategoryLocation `tfsdk:"location"`
	Name            types.String              `tfsdk:"name"`
	Description     types.String              `tfsdk:"description"`
	DisableOverride types.String              `tfsdk:"disable_override"`
	List            types.List                `tfsdk:"list"`
	Type            types.String              `tfsdk:"type"`
}

func (r *CustomUrlCategoryResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_custom_url_category"
}

func (r *CustomUrlCategoryResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
}

// <ResourceSchema>

func CustomUrlCategoryResourceSchema() rsschema.Schema {
	return rsschema.Schema{
		Attributes: map[string]rsschema.Attribute{

			"location": CustomUrlCategoryResourceLocationSchema(),

			"name": rsschema.StringAttribute{
				Description: "",
				Computed:    false,
				Required:    true,
				Optional:    false,
				Sensitive:   false,
			},

			"description": rsschema.StringAttribute{
				Description: "",
				Computed:    false,
				Required:    false,
				Optional:    true,
				Sensitive:   false,
			},

			"disable_override": rsschema.StringAttribute{
				Description: "disable object override in child device groups",
				Computed:    true,
				Required:    false,
				Optional:    true,
				Sensitive:   false,
				Default:     stringdefault.StaticString("no"),

				Validators: []validator.String{
					stringvalidator.OneOf([]string{
						"no",
					}...),
				},
			},

			"list": rsschema.ListAttribute{
				Description: "",
				Required:    false,
				Optional:    true,
				Computed:    false,
				Sensitive:   false,
				ElementType: types.StringType,
			},

			"type": rsschema.StringAttribute{
				Description: "",
				Computed:    true,
				Required:    false,
				Optional:    true,
				Sensitive:   false,
				Default:     stringdefault.StaticString("URL List"),
			},
		},
	}
}

func (o *CustomUrlCategoryResourceModel) getTypeFor(name string) attr.Type {
	schema := CustomUrlCategoryResourceSchema()
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

func (r *CustomUrlCategoryResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = CustomUrlCategoryResourceSchema()
}

// </ResourceSchema>

func (r *CustomUrlCategoryResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*pango.Client)
	specifier, _, err := customurlcategory.Versioning(r.client.Versioning())
	if err != nil {
		resp.Diagnostics.AddError("Failed to configure SDK client", err.Error())
		return
	}
	r.manager = sdkmanager.NewEntryObjectManager(r.client, customurlcategory.NewService(r.client), specifier, customurlcategory.SpecMatches)
}

func (o *CustomUrlCategoryResourceModel) CopyToPango(ctx context.Context, obj **customurlcategory.Entry, encrypted *map[string]types.String) diag.Diagnostics {
	var diags diag.Diagnostics
	description_value := o.Description.ValueStringPointer()
	disableOverride_value := o.DisableOverride.ValueStringPointer()
	list_pango_entries := make([]string, 0)
	diags.Append(o.List.ElementsAs(ctx, &list_pango_entries, false)...)
	if diags.HasError() {
		return diags
	}
	type_value := o.Type.ValueStringPointer()

	if (*obj) == nil {
		*obj = new(customurlcategory.Entry)
	}
	(*obj).Name = o.Name.ValueString()
	(*obj).Description = description_value
	(*obj).DisableOverride = disableOverride_value
	(*obj).List = list_pango_entries
	(*obj).Type = type_value

	return diags
}

func (o *CustomUrlCategoryResourceModel) CopyFromPango(ctx context.Context, obj *customurlcategory.Entry, encrypted *map[string]types.String) diag.Diagnostics {
	var diags diag.Diagnostics
	var list_list types.List
	{
		var list_diags diag.Diagnostics
		list_list, list_diags = types.ListValueFrom(ctx, types.StringType, obj.List)
		diags.Append(list_diags...)
	}

	var description_value types.String
	if obj.Description != nil {
		description_value = types.StringValue(*obj.Description)
	}
	var disableOverride_value types.String
	if obj.DisableOverride != nil {
		disableOverride_value = types.StringValue(*obj.DisableOverride)
	}
	var type_value types.String
	if obj.Type != nil {
		type_value = types.StringValue(*obj.Type)
	}
	o.Name = types.StringValue(obj.Name)
	o.Description = description_value
	o.DisableOverride = disableOverride_value
	o.List = list_list
	o.Type = type_value

	return diags
}

func (r *CustomUrlCategoryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state CustomUrlCategoryResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource create", map[string]any{
		"resource_name": "panos_custom_url_category_resource",
		"function":      "Create",
		"name":          state.Name.ValueString(),
	})

	// Verify mode.
	if r.client.Hostname == "" {
		resp.Diagnostics.AddError("Invalid mode error", InspectionModeError)
		return
	}

	// Determine the location.

	var location customurlcategory.Location

	if state.Location.Vsys != nil {
		location.Vsys = &customurlcategory.VsysLocation{

			NgfwDevice: state.Location.Vsys.NgfwDevice.ValueString(),
			Vsys:       state.Location.Vsys.Name.ValueString(),
		}
	}
	if state.Location.DeviceGroup != nil {
		location.DeviceGroup = &customurlcategory.DeviceGroupLocation{

			PanoramaDevice: state.Location.DeviceGroup.PanoramaDevice.ValueString(),
			DeviceGroup:    state.Location.DeviceGroup.Name.ValueString(),
		}
	}
	if !state.Location.Shared.IsNull() && state.Location.Shared.ValueBool() {
		location.Shared = true
	}

	if err := location.IsValid(); err != nil {
		resp.Diagnostics.AddError("Invalid location", err.Error())
		return
	}

	// Load the desired config.
	var obj *customurlcategory.Entry

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

func (o *CustomUrlCategoryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var savestate, state CustomUrlCategoryResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &savestate)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var location customurlcategory.Location

	if !savestate.Location.Shared.IsNull() && savestate.Location.Shared.ValueBool() {
		location.Shared = true
	}
	if savestate.Location.Vsys != nil {
		location.Vsys = &customurlcategory.VsysLocation{

			NgfwDevice: savestate.Location.Vsys.NgfwDevice.ValueString(),
			Vsys:       savestate.Location.Vsys.Name.ValueString(),
		}
	}
	if savestate.Location.DeviceGroup != nil {
		location.DeviceGroup = &customurlcategory.DeviceGroupLocation{

			DeviceGroup:    savestate.Location.DeviceGroup.Name.ValueString(),
			PanoramaDevice: savestate.Location.DeviceGroup.PanoramaDevice.ValueString(),
		}
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource read", map[string]any{
		"resource_name": "panos_custom_url_category_resource",
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

func (r *CustomUrlCategoryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var plan, state CustomUrlCategoryResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var location customurlcategory.Location

	if !state.Location.Shared.IsNull() && state.Location.Shared.ValueBool() {
		location.Shared = true
	}
	if state.Location.Vsys != nil {
		location.Vsys = &customurlcategory.VsysLocation{

			NgfwDevice: state.Location.Vsys.NgfwDevice.ValueString(),
			Vsys:       state.Location.Vsys.Name.ValueString(),
		}
	}
	if state.Location.DeviceGroup != nil {
		location.DeviceGroup = &customurlcategory.DeviceGroupLocation{

			PanoramaDevice: state.Location.DeviceGroup.PanoramaDevice.ValueString(),
			DeviceGroup:    state.Location.DeviceGroup.Name.ValueString(),
		}
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource update", map[string]any{
		"resource_name": "panos_custom_url_category_resource",
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

func (r *CustomUrlCategoryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var state CustomUrlCategoryResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource delete", map[string]any{
		"resource_name": "panos_custom_url_category_resource",
		"function":      "Delete",
		"name":          state.Name.ValueString(),
	})

	// Verify mode.
	if r.client.Hostname == "" {
		resp.Diagnostics.AddError("Invalid mode error", InspectionModeError)
		return
	}

	var location customurlcategory.Location

	if !state.Location.Shared.IsNull() && state.Location.Shared.ValueBool() {
		location.Shared = true
	}
	if state.Location.Vsys != nil {
		location.Vsys = &customurlcategory.VsysLocation{

			NgfwDevice: state.Location.Vsys.NgfwDevice.ValueString(),
			Vsys:       state.Location.Vsys.Name.ValueString(),
		}
	}
	if state.Location.DeviceGroup != nil {
		location.DeviceGroup = &customurlcategory.DeviceGroupLocation{

			PanoramaDevice: state.Location.DeviceGroup.PanoramaDevice.ValueString(),
			DeviceGroup:    state.Location.DeviceGroup.Name.ValueString(),
		}
	}

	err := r.manager.Delete(ctx, location, []string{state.Name.ValueString()})
	if err != nil && !errors.Is(err, sdkmanager.ErrObjectNotFound) {
		resp.Diagnostics.AddError("Error in delete", err.Error())
	}

}

type CustomUrlCategoryImportState struct {
	Location CustomUrlCategoryLocation `json:"location"`
	Name     string                    `json:"name"`
}

func CustomUrlCategoryImportStateCreator(ctx context.Context, resource types.Object) ([]byte, error) {
	attrs := resource.Attributes()
	if attrs == nil {
		return nil, fmt.Errorf("Object has no attributes")
	}

	locationAttr, ok := attrs["location"]
	if !ok {
		return nil, fmt.Errorf("location attribute missing")
	}

	var location CustomUrlCategoryLocation
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

	importStruct := CustomUrlCategoryImportState{
		Location: location,
		Name:     name,
	}

	return json.Marshal(importStruct)
}

func (r *CustomUrlCategoryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	var obj CustomUrlCategoryImportState
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
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), obj.Name)...)

}

type CustomUrlCategoryVsysLocation struct {
	NgfwDevice types.String `tfsdk:"ngfw_device"`
	Name       types.String `tfsdk:"name"`
}
type CustomUrlCategoryDeviceGroupLocation struct {
	PanoramaDevice types.String `tfsdk:"panorama_device"`
	Name           types.String `tfsdk:"name"`
}
type CustomUrlCategoryLocation struct {
	Shared      types.Bool                            `tfsdk:"shared"`
	Vsys        *CustomUrlCategoryVsysLocation        `tfsdk:"vsys"`
	DeviceGroup *CustomUrlCategoryDeviceGroupLocation `tfsdk:"device_group"`
}

func CustomUrlCategoryLocationSchema() rsschema.Attribute {
	return rsschema.SingleNestedAttribute{
		Description: "The location of this object.",
		Required:    true,
		Attributes: map[string]rsschema.Attribute{
			"shared": rsschema.BoolAttribute{
				Description: "Location in Shared Panorama",
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},

				Validators: []validator.Bool{
					boolvalidator.ExactlyOneOf(path.Expressions{
						path.MatchRelative().AtParent().AtName("shared"),
						path.MatchRelative().AtParent().AtName("vsys"),
						path.MatchRelative().AtParent().AtName("device_group"),
					}...),
				},
			},
			"vsys": rsschema.SingleNestedAttribute{
				Description: "Located in a specific Virtual System",
				Optional:    true,
				Attributes: map[string]rsschema.Attribute{
					"ngfw_device": rsschema.StringAttribute{
						Description: "The NGFW device name",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString("localhost.localdomain"),
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"name": rsschema.StringAttribute{
						Description: "The Virtual System name",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString("vsys1"),
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
			},
			"device_group": rsschema.SingleNestedAttribute{
				Description: "Located in a specific Device Group",
				Optional:    true,
				Attributes: map[string]rsschema.Attribute{
					"panorama_device": rsschema.StringAttribute{
						Description: "Panorama device name",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString("localhost.localdomain"),
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"name": rsschema.StringAttribute{
						Description: "Device Group name",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(""),
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

func (o CustomUrlCategoryVsysLocation) MarshalJSON() ([]byte, error) {
	obj := struct {
		Name       *string `json:"name"`
		NgfwDevice *string `json:"ngfw_device"`
	}{
		Name:       o.Name.ValueStringPointer(),
		NgfwDevice: o.NgfwDevice.ValueStringPointer(),
	}

	return json.Marshal(obj)
}

func (o *CustomUrlCategoryVsysLocation) UnmarshalJSON(data []byte) error {
	var shadow struct {
		Name       *string `json:"name"`
		NgfwDevice *string `json:"ngfw_device"`
	}

	err := json.Unmarshal(data, &shadow)
	if err != nil {
		return err
	}
	o.Name = types.StringPointerValue(shadow.Name)
	o.NgfwDevice = types.StringPointerValue(shadow.NgfwDevice)

	return nil
}
func (o CustomUrlCategoryDeviceGroupLocation) MarshalJSON() ([]byte, error) {
	obj := struct {
		PanoramaDevice *string `json:"panorama_device"`
		Name           *string `json:"name"`
	}{
		PanoramaDevice: o.PanoramaDevice.ValueStringPointer(),
		Name:           o.Name.ValueStringPointer(),
	}

	return json.Marshal(obj)
}

func (o *CustomUrlCategoryDeviceGroupLocation) UnmarshalJSON(data []byte) error {
	var shadow struct {
		PanoramaDevice *string `json:"panorama_device"`
		Name           *string `json:"name"`
	}

	err := json.Unmarshal(data, &shadow)
	if err != nil {
		return err
	}
	o.PanoramaDevice = types.StringPointerValue(shadow.PanoramaDevice)
	o.Name = types.StringPointerValue(shadow.Name)

	return nil
}
func (o CustomUrlCategoryLocation) MarshalJSON() ([]byte, error) {
	obj := struct {
		Shared      *bool                                 `json:"shared"`
		Vsys        *CustomUrlCategoryVsysLocation        `json:"vsys"`
		DeviceGroup *CustomUrlCategoryDeviceGroupLocation `json:"device_group"`
	}{
		Shared:      o.Shared.ValueBoolPointer(),
		Vsys:        o.Vsys,
		DeviceGroup: o.DeviceGroup,
	}

	return json.Marshal(obj)
}

func (o *CustomUrlCategoryLocation) UnmarshalJSON(data []byte) error {
	var shadow struct {
		Shared      *bool                                 `json:"shared"`
		Vsys        *CustomUrlCategoryVsysLocation        `json:"vsys"`
		DeviceGroup *CustomUrlCategoryDeviceGroupLocation `json:"device_group"`
	}

	err := json.Unmarshal(data, &shadow)
	if err != nil {
		return err
	}
	o.Shared = types.BoolPointerValue(shadow.Shared)
	o.Vsys = shadow.Vsys
	o.DeviceGroup = shadow.DeviceGroup

	return nil
}
