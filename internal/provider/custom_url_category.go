package provider

// Note:  This file is automatically generated.  Manually made changes
// will be overwritten when the provider is generated.

import (
	"context"
	"errors"
	"fmt"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/objects/profiles"

	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
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
	manager *sdkmanager.EntryObjectManager[*profiles.Entry, profiles.Location, *profiles.Service]
}

type CustomUrlCategoryDataSourceFilter struct {
	// TODO: Generate Data Source filter via function
}
type CustomUrlCategoryDataSourceTfid struct {
	Name     string            `json:"name"`
	Location profiles.Location `json:"location"`
}

func (o *CustomUrlCategoryDataSourceTfid) IsValid() error {
	if o.Name == "" {
		return fmt.Errorf("name is unspecified")
	}
	return o.Location.IsValid()
}

type CustomUrlCategoryDataSourceModel struct {
	Tfid            types.String              `tfsdk:"tfid"`
	Location        CustomUrlCategoryLocation `tfsdk:"location"`
	Name            types.String              `tfsdk:"name"`
	Description     types.String              `tfsdk:"description"`
	List            types.List                `tfsdk:"list"`
	Type            types.String              `tfsdk:"type"`
	DisableOverride types.Bool                `tfsdk:"disable_override"`
}

func (o *CustomUrlCategoryDataSourceModel) CopyToPango(ctx context.Context, obj **profiles.Entry, encrypted *map[string]types.String) diag.Diagnostics {
	var diags diag.Diagnostics
	list_pango_entries := make([]string, 0)
	diags.Append(o.List.ElementsAs(ctx, &list_pango_entries, false)...)
	if diags.HasError() {
		return diags
	}
	type_value := o.Type.ValueStringPointer()
	disableOverride_value := o.DisableOverride.ValueBoolPointer()
	description_value := o.Description.ValueStringPointer()

	if (*obj) == nil {
		*obj = new(profiles.Entry)
	}
	(*obj).Name = o.Name.ValueString()
	(*obj).List = list_pango_entries
	(*obj).Type = type_value
	(*obj).DisableOverride = disableOverride_value
	(*obj).Description = description_value

	return diags
}

func (o *CustomUrlCategoryDataSourceModel) CopyFromPango(ctx context.Context, obj *profiles.Entry, encrypted *map[string]types.String) diag.Diagnostics {
	var diags diag.Diagnostics
	var list_list types.List
	{
		var list_diags diag.Diagnostics
		list_list, list_diags = types.ListValueFrom(ctx, types.StringType, obj.List)
		diags.Append(list_diags...)
	}
	var type_value types.String
	if obj.Type != nil {
		type_value = types.StringValue(*obj.Type)
	}
	var disableOverride_value types.Bool
	if obj.DisableOverride != nil {
		disableOverride_value = types.BoolValue(*obj.DisableOverride)
	}
	var description_value types.String
	if obj.Description != nil {
		description_value = types.StringValue(*obj.Description)
	}
	o.Name = types.StringValue(obj.Name)
	o.Type = type_value
	o.DisableOverride = disableOverride_value
	o.Description = description_value
	o.List = list_list

	return diags
}

func CustomUrlCategoryDataSourceSchema() dsschema.Schema {
	return dsschema.Schema{
		Attributes: map[string]dsschema.Attribute{

			"location": CustomUrlCategoryDataSourceLocationSchema(),

			"tfid": dsschema.StringAttribute{
				Description: "The Terraform ID.",
				Computed:    true,
				Required:    false,
				Optional:    false,
				Sensitive:   false,
			},

			"name": dsschema.StringAttribute{
				Description: "Name of the custom category",
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

			"disable_override": dsschema.BoolAttribute{
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
	specifier, _, err := profiles.Versioning(d.client.Versioning())
	if err != nil {
		resp.Diagnostics.AddError("Failed to configure SDK client", err.Error())
		return
	}
	d.manager = sdkmanager.NewEntryObjectManager(d.client, profiles.NewService(d.client), specifier, profiles.SpecMatches)
}

func (o *CustomUrlCategoryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var savestate, state CustomUrlCategoryDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &savestate)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var loc CustomUrlCategoryDataSourceTfid
	loc.Name = *savestate.Name.ValueStringPointer()

	if !savestate.Location.Shared.IsNull() && savestate.Location.Shared.ValueBool() {
		loc.Location.Shared = true
	}
	if savestate.Location.Vsys != nil {
		loc.Location.Vsys = &profiles.VsysLocation{

			NgfwDevice: savestate.Location.Vsys.NgfwDevice.ValueString(),
			Vsys:       savestate.Location.Vsys.Name.ValueString(),
		}
	}
	if !savestate.Location.FromPanoramaShared.IsNull() && savestate.Location.FromPanoramaShared.ValueBool() {
		loc.Location.FromPanoramaShared = true
	}
	if savestate.Location.FromPanoramaVsys != nil {
		loc.Location.FromPanoramaVsys = &profiles.FromPanoramaVsysLocation{

			Vsys: savestate.Location.FromPanoramaVsys.Vsys.ValueString(),
		}
	}
	if savestate.Location.DeviceGroup != nil {
		loc.Location.DeviceGroup = &profiles.DeviceGroupLocation{

			PanoramaDevice: savestate.Location.DeviceGroup.PanoramaDevice.ValueString(),
			DeviceGroup:    savestate.Location.DeviceGroup.Name.ValueString(),
		}
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource read", map[string]any{
		"resource_name": "panos_custom_url_category_resource",
		"function":      "Read",
		"name":          loc.Name,
	})

	// Perform the operation.
	object, err := o.manager.Read(ctx, loc.Location, loc.Name)
	if err != nil {
		tflog.Warn(ctx, "KK: HERE3-1", map[string]any{"Error": err.Error()})
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
	// Save tfid to state.
	state.Tfid = savestate.Tfid

	// Save the answer to state.

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
	return &CustomUrlCategoryResource{}
}

type CustomUrlCategoryResource struct {
	client  *pango.Client
	manager *sdkmanager.EntryObjectManager[*profiles.Entry, profiles.Location, *profiles.Service]
}
type CustomUrlCategoryResourceTfid struct {
	Name     string            `json:"name"`
	Location profiles.Location `json:"location"`
}

func (o *CustomUrlCategoryResourceTfid) IsValid() error {
	if o.Name == "" {
		return fmt.Errorf("name is unspecified")
	}
	return o.Location.IsValid()
}

func CustomUrlCategoryResourceLocationSchema() rsschema.Attribute {
	return CustomUrlCategoryLocationSchema()
}

type CustomUrlCategoryResourceModel struct {
	Tfid            types.String              `tfsdk:"tfid"`
	Location        CustomUrlCategoryLocation `tfsdk:"location"`
	Name            types.String              `tfsdk:"name"`
	DisableOverride types.Bool                `tfsdk:"disable_override"`
	Description     types.String              `tfsdk:"description"`
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

			"tfid": rsschema.StringAttribute{
				Description: "The Terraform ID.",
				Computed:    true,
				Required:    false,
				Optional:    false,
				Sensitive:   false,
			},

			"name": rsschema.StringAttribute{
				Description: "Name of the custom category",
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
				Computed:    false,
				Required:    false,
				Optional:    true,
				Sensitive:   false,
			},

			"disable_override": rsschema.BoolAttribute{
				Description: "",
				Computed:    false,
				Required:    false,
				Optional:    true,
				Sensitive:   false,
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
	specifier, _, err := profiles.Versioning(r.client.Versioning())
	if err != nil {
		resp.Diagnostics.AddError("Failed to configure SDK client", err.Error())
		return
	}
	r.manager = sdkmanager.NewEntryObjectManager(r.client, profiles.NewService(r.client), specifier, profiles.SpecMatches)
}

func (o *CustomUrlCategoryResourceModel) CopyToPango(ctx context.Context, obj **profiles.Entry, encrypted *map[string]types.String) diag.Diagnostics {
	var diags diag.Diagnostics
	description_value := o.Description.ValueStringPointer()
	list_pango_entries := make([]string, 0)
	diags.Append(o.List.ElementsAs(ctx, &list_pango_entries, false)...)
	if diags.HasError() {
		return diags
	}
	type_value := o.Type.ValueStringPointer()
	disableOverride_value := o.DisableOverride.ValueBoolPointer()

	if (*obj) == nil {
		*obj = new(profiles.Entry)
	}
	(*obj).Name = o.Name.ValueString()
	(*obj).Description = description_value
	(*obj).List = list_pango_entries
	(*obj).Type = type_value
	(*obj).DisableOverride = disableOverride_value

	return diags
}

func (o *CustomUrlCategoryResourceModel) CopyFromPango(ctx context.Context, obj *profiles.Entry, encrypted *map[string]types.String) diag.Diagnostics {
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
	var type_value types.String
	if obj.Type != nil {
		type_value = types.StringValue(*obj.Type)
	}
	var disableOverride_value types.Bool
	if obj.DisableOverride != nil {
		disableOverride_value = types.BoolValue(*obj.DisableOverride)
	}
	o.Name = types.StringValue(obj.Name)
	o.Description = description_value
	o.List = list_list
	o.Type = type_value
	o.DisableOverride = disableOverride_value

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
	loc := CustomUrlCategoryResourceTfid{Name: state.Name.ValueString()}

	// TODO: this needs to handle location structure for UUID style shared has nested structure type

	if !state.Location.Shared.IsNull() && state.Location.Shared.ValueBool() {
		loc.Location.Shared = true
	}
	if state.Location.Vsys != nil {
		loc.Location.Vsys = &profiles.VsysLocation{

			NgfwDevice: state.Location.Vsys.NgfwDevice.ValueString(),
			Vsys:       state.Location.Vsys.Name.ValueString(),
		}
	}
	if !state.Location.FromPanoramaShared.IsNull() && state.Location.FromPanoramaShared.ValueBool() {
		loc.Location.FromPanoramaShared = true
	}
	if state.Location.FromPanoramaVsys != nil {
		loc.Location.FromPanoramaVsys = &profiles.FromPanoramaVsysLocation{

			Vsys: state.Location.FromPanoramaVsys.Vsys.ValueString(),
		}
	}
	if state.Location.DeviceGroup != nil {
		loc.Location.DeviceGroup = &profiles.DeviceGroupLocation{

			PanoramaDevice: state.Location.DeviceGroup.PanoramaDevice.ValueString(),
			DeviceGroup:    state.Location.DeviceGroup.Name.ValueString(),
		}
	}

	if err := loc.IsValid(); err != nil {
		resp.Diagnostics.AddError("Invalid location", err.Error())
		return
	}

	// Load the desired config.
	var obj *profiles.Entry

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
	created, err := r.manager.Create(ctx, loc.Location, obj)
	if err != nil {
		resp.Diagnostics.AddError("Error in create", err.Error())
		return
	}

	// Tfid handling.
	tfid, err := EncodeLocation(&loc)
	if err != nil {
		resp.Diagnostics.AddError("Error creating tfid", err.Error())
		return
	}

	// Save the state.
	state.Tfid = types.StringValue(tfid)

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
	var loc CustomUrlCategoryResourceTfid
	// Parse the location from tfid.
	if err := DecodeLocation(savestate.Tfid.ValueString(), &loc); err != nil {
		resp.Diagnostics.AddError("Error parsing tfid", err.Error())
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource read", map[string]any{
		"resource_name": "panos_custom_url_category_resource",
		"function":      "Read",
		"name":          loc.Name,
	})

	// Perform the operation.
	object, err := o.manager.Read(ctx, loc.Location, loc.Name)
	if err != nil {
		tflog.Warn(ctx, "KK: HERE3-1", map[string]any{"Error": err.Error()})
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
	// Save tfid to state.
	state.Tfid = savestate.Tfid

	// Save the answer to state.

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

	var loc CustomUrlCategoryResourceTfid
	if err := DecodeLocation(state.Tfid.ValueString(), &loc); err != nil {
		resp.Diagnostics.AddError("Error parsing tfid", err.Error())
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource update", map[string]any{
		"resource_name": "panos_custom_url_category_resource",
		"function":      "Update",
		"tfid":          state.Tfid.ValueString(),
	})

	// Verify mode.
	if r.client.Hostname == "" {
		resp.Diagnostics.AddError("Invalid mode error", InspectionModeError)
		return
	}
	obj, err := r.manager.Read(ctx, loc.Location, loc.Name)
	if err != nil {
		resp.Diagnostics.AddError("Error in update", err.Error())
		return
	}

	resp.Diagnostics.Append(plan.CopyToPango(ctx, &obj, nil)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Perform the operation.
	updated, err := r.manager.Update(ctx, loc.Location, obj, loc.Name)
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

	// Save the tfid.
	loc.Name = obj.Name
	tfid, err := EncodeLocation(&loc)
	if err != nil {
		resp.Diagnostics.AddError("error creating tfid", err.Error())
		return
	}
	state.Tfid = types.StringValue(tfid)

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

	// Parse the location from tfid.
	var loc CustomUrlCategoryResourceTfid
	if err := DecodeLocation(state.Tfid.ValueString(), &loc); err != nil {
		resp.Diagnostics.AddError("error parsing tfid", err.Error())
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource delete", map[string]any{
		"resource_name": "panos_custom_url_category_resource",
		"function":      "Delete",
		"name":          loc.Name,
	})

	// Verify mode.
	if r.client.Hostname == "" {
		resp.Diagnostics.AddError("Invalid mode error", InspectionModeError)
		return
	}
	err := r.manager.Delete(ctx, loc.Location, []string{loc.Name})
	if err != nil && !errors.Is(err, sdkmanager.ErrObjectNotFound) {
		resp.Diagnostics.AddError("Error in delete", err.Error())
	}

}

func (r *CustomUrlCategoryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("tfid"), req, resp)
}

type CustomUrlCategoryFromPanoramaVsysLocation struct {
	Vsys types.String `tfsdk:"vsys"`
}
type CustomUrlCategoryDeviceGroupLocation struct {
	PanoramaDevice types.String `tfsdk:"panorama_device"`
	Name           types.String `tfsdk:"name"`
}
type CustomUrlCategoryVsysLocation struct {
	NgfwDevice types.String `tfsdk:"ngfw_device"`
	Name       types.String `tfsdk:"name"`
}
type CustomUrlCategoryLocation struct {
	FromPanoramaVsys   *CustomUrlCategoryFromPanoramaVsysLocation `tfsdk:"from_panorama_vsys"`
	DeviceGroup        *CustomUrlCategoryDeviceGroupLocation      `tfsdk:"device_group"`
	Shared             types.Bool                                 `tfsdk:"shared"`
	Vsys               *CustomUrlCategoryVsysLocation             `tfsdk:"vsys"`
	FromPanoramaShared types.Bool                                 `tfsdk:"from_panorama_shared"`
}

func CustomUrlCategoryLocationSchema() rsschema.Attribute {
	return rsschema.SingleNestedAttribute{
		Description: "The location of this object.",
		Required:    true,
		Attributes: map[string]rsschema.Attribute{
			"shared": rsschema.BoolAttribute{
				Description: "Located in shared.",
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},

				Validators: []validator.Bool{
					boolvalidator.ExactlyOneOf(path.Expressions{
						path.MatchRelative().AtParent().AtName("shared"),
						path.MatchRelative().AtParent().AtName("vsys"),
						path.MatchRelative().AtParent().AtName("from_panorama_shared"),
						path.MatchRelative().AtParent().AtName("from_panorama_vsys"),
						path.MatchRelative().AtParent().AtName("device_group"),
					}...),
				},
			},
			"vsys": rsschema.SingleNestedAttribute{
				Description: "Located in a specific vsys.",
				Optional:    true,
				Attributes: map[string]rsschema.Attribute{
					"ngfw_device": rsschema.StringAttribute{
						Description: "The NGFW device.",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString("localhost.localdomain"),
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"name": rsschema.StringAttribute{
						Description: "The vsys.",
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
			"from_panorama_shared": rsschema.BoolAttribute{
				Description: "Located in shared in the config pushed from Panorama.",
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"from_panorama_vsys": rsschema.SingleNestedAttribute{
				Description: "Located in a specific vsys in the config pushed from Panorama.",
				Optional:    true,
				Attributes: map[string]rsschema.Attribute{
					"vsys": rsschema.StringAttribute{
						Description: "The vsys.",
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
				Description: "Located in a specific device group.",
				Optional:    true,
				Attributes: map[string]rsschema.Attribute{
					"name": rsschema.StringAttribute{
						Description: "The device group.",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(""),
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"panorama_device": rsschema.StringAttribute{
						Description: "The panorama device.",
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
