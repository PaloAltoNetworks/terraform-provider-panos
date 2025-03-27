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
	"github.com/PaloAltoNetworks/pango/objects/address/group"

	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
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
	_ datasource.DataSource              = &AddressGroupDataSource{}
	_ datasource.DataSourceWithConfigure = &AddressGroupDataSource{}
)

func NewAddressGroupDataSource() datasource.DataSource {
	return &AddressGroupDataSource{}
}

type AddressGroupDataSource struct {
	client  *pango.Client
	manager *sdkmanager.EntryObjectManager[*group.Entry, group.Location, *group.Service]
}

type AddressGroupDataSourceFilter struct {
	// TODO: Generate Data Source filter via function
}

type AddressGroupDataSourceModel struct {
	Location        AddressGroupLocation                 `tfsdk:"location"`
	Name            types.String                         `tfsdk:"name"`
	Description     types.String                         `tfsdk:"description"`
	DisableOverride types.String                         `tfsdk:"disable_override"`
	Tag             types.List                           `tfsdk:"tag"`
	Dynamic         *AddressGroupDataSourceDynamicObject `tfsdk:"dynamic"`
	Static          types.List                           `tfsdk:"static"`
}
type AddressGroupDataSourceDynamicObject struct {
	Filter types.String `tfsdk:"filter"`
}

func (o *AddressGroupDataSourceModel) CopyToPango(ctx context.Context, obj **group.Entry, encrypted *map[string]types.String) diag.Diagnostics {
	var diags diag.Diagnostics
	description_value := o.Description.ValueStringPointer()
	disableOverride_value := o.DisableOverride.ValueStringPointer()
	tag_pango_entries := make([]string, 0)
	diags.Append(o.Tag.ElementsAs(ctx, &tag_pango_entries, false)...)
	if diags.HasError() {
		return diags
	}
	var dynamic_entry *group.Dynamic
	if o.Dynamic != nil {
		if *obj != nil && (*obj).Dynamic != nil {
			dynamic_entry = (*obj).Dynamic
		} else {
			dynamic_entry = new(group.Dynamic)
		}

		diags.Append(o.Dynamic.CopyToPango(ctx, &dynamic_entry, encrypted)...)
		if diags.HasError() {
			return diags
		}
	}
	static_pango_entries := make([]string, 0)
	diags.Append(o.Static.ElementsAs(ctx, &static_pango_entries, false)...)
	if diags.HasError() {
		return diags
	}

	if (*obj) == nil {
		*obj = new(group.Entry)
	}
	(*obj).Name = o.Name.ValueString()
	(*obj).Description = description_value
	(*obj).DisableOverride = disableOverride_value
	(*obj).Tag = tag_pango_entries
	(*obj).Dynamic = dynamic_entry
	(*obj).Static = static_pango_entries

	return diags
}
func (o *AddressGroupDataSourceDynamicObject) CopyToPango(ctx context.Context, obj **group.Dynamic, encrypted *map[string]types.String) diag.Diagnostics {
	var diags diag.Diagnostics
	filter_value := o.Filter.ValueStringPointer()

	if (*obj) == nil {
		*obj = new(group.Dynamic)
	}
	(*obj).Filter = filter_value

	return diags
}

func (o *AddressGroupDataSourceModel) CopyFromPango(ctx context.Context, obj *group.Entry, encrypted *map[string]types.String) diag.Diagnostics {
	var diags diag.Diagnostics
	var tag_list types.List
	{
		var list_diags diag.Diagnostics
		tag_list, list_diags = types.ListValueFrom(ctx, types.StringType, obj.Tag)
		diags.Append(list_diags...)
	}
	var static_list types.List
	{
		var list_diags diag.Diagnostics
		static_list, list_diags = types.ListValueFrom(ctx, types.StringType, obj.Static)
		diags.Append(list_diags...)
	}
	var dynamic_object *AddressGroupDataSourceDynamicObject
	if obj.Dynamic != nil {
		dynamic_object = new(AddressGroupDataSourceDynamicObject)

		diags.Append(dynamic_object.CopyFromPango(ctx, obj.Dynamic, encrypted)...)
		if diags.HasError() {
			return diags
		}
	}

	var description_value types.String
	if obj.Description != nil {
		description_value = types.StringValue(*obj.Description)
	}
	var disableOverride_value types.String
	if obj.DisableOverride != nil {
		disableOverride_value = types.StringValue(*obj.DisableOverride)
	}
	o.Name = types.StringValue(obj.Name)
	o.Description = description_value
	o.DisableOverride = disableOverride_value
	o.Tag = tag_list
	o.Dynamic = dynamic_object
	o.Static = static_list

	return diags
}

func (o *AddressGroupDataSourceDynamicObject) CopyFromPango(ctx context.Context, obj *group.Dynamic, encrypted *map[string]types.String) diag.Diagnostics {
	var diags diag.Diagnostics

	var filter_value types.String
	if obj.Filter != nil {
		filter_value = types.StringValue(*obj.Filter)
	}
	o.Filter = filter_value

	return diags
}

func AddressGroupDataSourceSchema() dsschema.Schema {
	return dsschema.Schema{
		Attributes: map[string]dsschema.Attribute{

			"location": AddressGroupDataSourceLocationSchema(),

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

			"tag": dsschema.ListAttribute{
				Description: "",
				Required:    false,
				Optional:    true,
				Computed:    true,
				Sensitive:   false,
				ElementType: types.StringType,
			},

			"dynamic": AddressGroupDataSourceDynamicSchema(),

			"static": dsschema.ListAttribute{
				Description: "",
				Required:    false,
				Optional:    true,
				Computed:    true,
				Sensitive:   false,
				ElementType: types.StringType,
			},
		},
	}
}

func (o *AddressGroupDataSourceModel) getTypeFor(name string) attr.Type {
	schema := AddressGroupDataSourceSchema()
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

func AddressGroupDataSourceDynamicSchema() dsschema.SingleNestedAttribute {
	return dsschema.SingleNestedAttribute{
		Description: "",
		Required:    false,
		Computed:    true,
		Optional:    true,
		Sensitive:   false,
		Attributes: map[string]dsschema.Attribute{

			"filter": dsschema.StringAttribute{
				Description: "tag-based filter",
				Computed:    true,
				Required:    false,
				Optional:    true,
				Sensitive:   false,
			},
		},
	}
}

func (o *AddressGroupDataSourceDynamicObject) getTypeFor(name string) attr.Type {
	schema := AddressGroupDataSourceDynamicSchema()
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

func AddressGroupDataSourceLocationSchema() rsschema.Attribute {
	return AddressGroupLocationSchema()
}

// Metadata returns the data source type name.
func (d *AddressGroupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_address_group"
}

// Schema defines the schema for this data source.
func (d *AddressGroupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = AddressGroupDataSourceSchema()
}

// Configure prepares the struct.
func (d *AddressGroupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*pango.Client)
	specifier, _, err := group.Versioning(d.client.Versioning())
	if err != nil {
		resp.Diagnostics.AddError("Failed to configure SDK client", err.Error())
		return
	}
	d.manager = sdkmanager.NewEntryObjectManager(d.client, group.NewService(d.client), specifier, group.SpecMatches)
}
func (o *AddressGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var savestate, state AddressGroupDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &savestate)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var location group.Location

	if !savestate.Location.Shared.IsNull() && savestate.Location.Shared.ValueBool() {
		location.Shared = true
	}
	if savestate.Location.Vsys != nil {
		location.Vsys = &group.VsysLocation{

			NgfwDevice: savestate.Location.Vsys.NgfwDevice.ValueString(),
			Vsys:       savestate.Location.Vsys.Name.ValueString(),
		}
	}
	if savestate.Location.DeviceGroup != nil {
		location.DeviceGroup = &group.DeviceGroupLocation{

			PanoramaDevice: savestate.Location.DeviceGroup.PanoramaDevice.ValueString(),
			DeviceGroup:    savestate.Location.DeviceGroup.Name.ValueString(),
		}
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource read", map[string]any{
		"resource_name": "panos_address_group_resource",
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
	_ resource.Resource                = &AddressGroupResource{}
	_ resource.ResourceWithConfigure   = &AddressGroupResource{}
	_ resource.ResourceWithImportState = &AddressGroupResource{}
)

func NewAddressGroupResource() resource.Resource {
	if _, found := resourceFuncMap["panos_address_group"]; !found {
		resourceFuncMap["panos_address_group"] = resourceFuncs{
			CreateImportId: AddressGroupImportStateCreator,
		}
	}
	return &AddressGroupResource{}
}

type AddressGroupResource struct {
	client  *pango.Client
	manager *sdkmanager.EntryObjectManager[*group.Entry, group.Location, *group.Service]
}

func AddressGroupResourceLocationSchema() rsschema.Attribute {
	return AddressGroupLocationSchema()
}

type AddressGroupResourceModel struct {
	Location        AddressGroupLocation               `tfsdk:"location"`
	Name            types.String                       `tfsdk:"name"`
	Description     types.String                       `tfsdk:"description"`
	DisableOverride types.String                       `tfsdk:"disable_override"`
	Tag             types.List                         `tfsdk:"tag"`
	Dynamic         *AddressGroupResourceDynamicObject `tfsdk:"dynamic"`
	Static          types.List                         `tfsdk:"static"`
}
type AddressGroupResourceDynamicObject struct {
	Filter types.String `tfsdk:"filter"`
}

func (r *AddressGroupResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
}

// <ResourceSchema>

func AddressGroupResourceSchema() rsschema.Schema {
	return rsschema.Schema{
		Attributes: map[string]rsschema.Attribute{

			"location": AddressGroupResourceLocationSchema(),

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
				Computed:    false,
				Required:    false,
				Optional:    true,
				Sensitive:   false,

				Validators: []validator.String{
					stringvalidator.OneOf([]string{
						"no",
					}...),
				},
			},

			"tag": rsschema.ListAttribute{
				Description: "",
				Required:    false,
				Optional:    true,
				Computed:    false,
				Sensitive:   false,
				ElementType: types.StringType,
			},

			"dynamic": AddressGroupResourceDynamicSchema(),

			"static": rsschema.ListAttribute{
				Description: "",
				Required:    false,
				Optional:    true,
				Computed:    false,
				Sensitive:   false,
				ElementType: types.StringType,
			},
		},
	}
}

func (o *AddressGroupResourceModel) getTypeFor(name string) attr.Type {
	schema := AddressGroupResourceSchema()
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

func AddressGroupResourceDynamicSchema() rsschema.SingleNestedAttribute {
	return rsschema.SingleNestedAttribute{
		Description: "",
		Required:    false,
		Computed:    false,
		Optional:    true,
		Sensitive:   false,

		Validators: []validator.Object{
			objectvalidator.ExactlyOneOf(path.Expressions{
				path.MatchRelative().AtParent().AtName("dynamic"),
				path.MatchRelative().AtParent().AtName("static"),
			}...),
		},
		Attributes: map[string]rsschema.Attribute{

			"filter": rsschema.StringAttribute{
				Description: "tag-based filter",
				Computed:    false,
				Required:    false,
				Optional:    true,
				Sensitive:   false,
			},
		},
	}
}

func (o *AddressGroupResourceDynamicObject) getTypeFor(name string) attr.Type {
	schema := AddressGroupResourceDynamicSchema()
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

func (r *AddressGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_address_group"
}

func (r *AddressGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = AddressGroupResourceSchema()
}

// </ResourceSchema>

func (r *AddressGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*pango.Client)
	specifier, _, err := group.Versioning(r.client.Versioning())
	if err != nil {
		resp.Diagnostics.AddError("Failed to configure SDK client", err.Error())
		return
	}
	r.manager = sdkmanager.NewEntryObjectManager(r.client, group.NewService(r.client), specifier, group.SpecMatches)
}

func (o *AddressGroupResourceModel) CopyToPango(ctx context.Context, obj **group.Entry, encrypted *map[string]types.String) diag.Diagnostics {
	var diags diag.Diagnostics
	description_value := o.Description.ValueStringPointer()
	disableOverride_value := o.DisableOverride.ValueStringPointer()
	tag_pango_entries := make([]string, 0)
	diags.Append(o.Tag.ElementsAs(ctx, &tag_pango_entries, false)...)
	if diags.HasError() {
		return diags
	}
	var dynamic_entry *group.Dynamic
	if o.Dynamic != nil {
		if *obj != nil && (*obj).Dynamic != nil {
			dynamic_entry = (*obj).Dynamic
		} else {
			dynamic_entry = new(group.Dynamic)
		}

		diags.Append(o.Dynamic.CopyToPango(ctx, &dynamic_entry, encrypted)...)
		if diags.HasError() {
			return diags
		}
	}
	static_pango_entries := make([]string, 0)
	diags.Append(o.Static.ElementsAs(ctx, &static_pango_entries, false)...)
	if diags.HasError() {
		return diags
	}

	if (*obj) == nil {
		*obj = new(group.Entry)
	}
	(*obj).Name = o.Name.ValueString()
	(*obj).Description = description_value
	(*obj).DisableOverride = disableOverride_value
	(*obj).Tag = tag_pango_entries
	(*obj).Dynamic = dynamic_entry
	(*obj).Static = static_pango_entries

	return diags
}
func (o *AddressGroupResourceDynamicObject) CopyToPango(ctx context.Context, obj **group.Dynamic, encrypted *map[string]types.String) diag.Diagnostics {
	var diags diag.Diagnostics
	filter_value := o.Filter.ValueStringPointer()

	if (*obj) == nil {
		*obj = new(group.Dynamic)
	}
	(*obj).Filter = filter_value

	return diags
}

func (o *AddressGroupResourceModel) CopyFromPango(ctx context.Context, obj *group.Entry, encrypted *map[string]types.String) diag.Diagnostics {
	var diags diag.Diagnostics
	var tag_list types.List
	{
		var list_diags diag.Diagnostics
		tag_list, list_diags = types.ListValueFrom(ctx, types.StringType, obj.Tag)
		diags.Append(list_diags...)
	}
	var static_list types.List
	{
		var list_diags diag.Diagnostics
		static_list, list_diags = types.ListValueFrom(ctx, types.StringType, obj.Static)
		diags.Append(list_diags...)
	}
	var dynamic_object *AddressGroupResourceDynamicObject
	if obj.Dynamic != nil {
		dynamic_object = new(AddressGroupResourceDynamicObject)

		diags.Append(dynamic_object.CopyFromPango(ctx, obj.Dynamic, encrypted)...)
		if diags.HasError() {
			return diags
		}
	}

	var description_value types.String
	if obj.Description != nil {
		description_value = types.StringValue(*obj.Description)
	}
	var disableOverride_value types.String
	if obj.DisableOverride != nil {
		disableOverride_value = types.StringValue(*obj.DisableOverride)
	}
	o.Name = types.StringValue(obj.Name)
	o.Description = description_value
	o.DisableOverride = disableOverride_value
	o.Tag = tag_list
	o.Dynamic = dynamic_object
	o.Static = static_list

	return diags
}

func (o *AddressGroupResourceDynamicObject) CopyFromPango(ctx context.Context, obj *group.Dynamic, encrypted *map[string]types.String) diag.Diagnostics {
	var diags diag.Diagnostics

	var filter_value types.String
	if obj.Filter != nil {
		filter_value = types.StringValue(*obj.Filter)
	}
	o.Filter = filter_value

	return diags
}

func (r *AddressGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state AddressGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource create", map[string]any{
		"resource_name": "panos_address_group_resource",
		"function":      "Create",
		"name":          state.Name.ValueString(),
	})

	// Verify mode.
	if r.client.Hostname == "" {
		resp.Diagnostics.AddError("Invalid mode error", InspectionModeError)
		return
	}

	// Determine the location.

	var location group.Location

	if !state.Location.Shared.IsNull() && state.Location.Shared.ValueBool() {
		location.Shared = true
	}
	if state.Location.Vsys != nil {
		location.Vsys = &group.VsysLocation{

			NgfwDevice: state.Location.Vsys.NgfwDevice.ValueString(),
			Vsys:       state.Location.Vsys.Name.ValueString(),
		}
	}
	if state.Location.DeviceGroup != nil {
		location.DeviceGroup = &group.DeviceGroupLocation{

			PanoramaDevice: state.Location.DeviceGroup.PanoramaDevice.ValueString(),
			DeviceGroup:    state.Location.DeviceGroup.Name.ValueString(),
		}
	}

	if err := location.IsValid(); err != nil {
		resp.Diagnostics.AddError("Invalid location", err.Error())
		return
	}

	// Load the desired config.
	var obj *group.Entry

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
func (o *AddressGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var savestate, state AddressGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &savestate)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var location group.Location

	if !savestate.Location.Shared.IsNull() && savestate.Location.Shared.ValueBool() {
		location.Shared = true
	}
	if savestate.Location.Vsys != nil {
		location.Vsys = &group.VsysLocation{

			NgfwDevice: savestate.Location.Vsys.NgfwDevice.ValueString(),
			Vsys:       savestate.Location.Vsys.Name.ValueString(),
		}
	}
	if savestate.Location.DeviceGroup != nil {
		location.DeviceGroup = &group.DeviceGroupLocation{

			PanoramaDevice: savestate.Location.DeviceGroup.PanoramaDevice.ValueString(),
			DeviceGroup:    savestate.Location.DeviceGroup.Name.ValueString(),
		}
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource read", map[string]any{
		"resource_name": "panos_address_group_resource",
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
func (r *AddressGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var plan, state AddressGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var location group.Location

	if !state.Location.Shared.IsNull() && state.Location.Shared.ValueBool() {
		location.Shared = true
	}
	if state.Location.Vsys != nil {
		location.Vsys = &group.VsysLocation{

			NgfwDevice: state.Location.Vsys.NgfwDevice.ValueString(),
			Vsys:       state.Location.Vsys.Name.ValueString(),
		}
	}
	if state.Location.DeviceGroup != nil {
		location.DeviceGroup = &group.DeviceGroupLocation{

			PanoramaDevice: state.Location.DeviceGroup.PanoramaDevice.ValueString(),
			DeviceGroup:    state.Location.DeviceGroup.Name.ValueString(),
		}
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource update", map[string]any{
		"resource_name": "panos_address_group_resource",
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
func (r *AddressGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var state AddressGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource delete", map[string]any{
		"resource_name": "panos_address_group_resource",
		"function":      "Delete",
		"name":          state.Name.ValueString(),
	})

	// Verify mode.
	if r.client.Hostname == "" {
		resp.Diagnostics.AddError("Invalid mode error", InspectionModeError)
		return
	}

	var location group.Location

	if !state.Location.Shared.IsNull() && state.Location.Shared.ValueBool() {
		location.Shared = true
	}
	if state.Location.Vsys != nil {
		location.Vsys = &group.VsysLocation{

			NgfwDevice: state.Location.Vsys.NgfwDevice.ValueString(),
			Vsys:       state.Location.Vsys.Name.ValueString(),
		}
	}
	if state.Location.DeviceGroup != nil {
		location.DeviceGroup = &group.DeviceGroupLocation{

			PanoramaDevice: state.Location.DeviceGroup.PanoramaDevice.ValueString(),
			DeviceGroup:    state.Location.DeviceGroup.Name.ValueString(),
		}
	}

	err := r.manager.Delete(ctx, location, []string{state.Name.ValueString()})
	if err != nil && !errors.Is(err, sdkmanager.ErrObjectNotFound) {
		resp.Diagnostics.AddError("Error in delete", err.Error())
	}

}

type AddressGroupImportState struct {
	Location AddressGroupLocation `json:"location"`
	Name     string               `json:"name"`
}

func AddressGroupImportStateCreator(ctx context.Context, resource types.Object) ([]byte, error) {
	attrs := resource.Attributes()
	if attrs == nil {
		return nil, fmt.Errorf("Object has no attributes")
	}

	locationAttr, ok := attrs["location"]
	if !ok {
		return nil, fmt.Errorf("location attribute missing")
	}

	var location AddressGroupLocation
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

	importStruct := AddressGroupImportState{
		Location: location,
		Name:     name,
	}

	return json.Marshal(importStruct)
}

func (r *AddressGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	var obj AddressGroupImportState
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

type AddressGroupVsysLocation struct {
	NgfwDevice types.String `tfsdk:"ngfw_device"`
	Name       types.String `tfsdk:"name"`
}
type AddressGroupDeviceGroupLocation struct {
	PanoramaDevice types.String `tfsdk:"panorama_device"`
	Name           types.String `tfsdk:"name"`
}
type AddressGroupLocation struct {
	Shared      types.Bool                       `tfsdk:"shared"`
	Vsys        *AddressGroupVsysLocation        `tfsdk:"vsys"`
	DeviceGroup *AddressGroupDeviceGroupLocation `tfsdk:"device_group"`
}

func AddressGroupLocationSchema() rsschema.Attribute {
	return rsschema.SingleNestedAttribute{
		Description: "The location of this object.",
		Required:    true,
		Attributes: map[string]rsschema.Attribute{
			"shared": rsschema.BoolAttribute{
				Description: "Panorama shared object",
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

func (o AddressGroupVsysLocation) MarshalJSON() ([]byte, error) {
	obj := struct {
		NgfwDevice *string `json:"ngfw_device"`
		Name       *string `json:"name"`
	}{
		NgfwDevice: o.NgfwDevice.ValueStringPointer(),
		Name:       o.Name.ValueStringPointer(),
	}

	return json.Marshal(obj)
}

func (o *AddressGroupVsysLocation) UnmarshalJSON(data []byte) error {
	var shadow struct {
		NgfwDevice *string `json:"ngfw_device"`
		Name       *string `json:"name"`
	}

	err := json.Unmarshal(data, &shadow)
	if err != nil {
		return err
	}
	o.NgfwDevice = types.StringPointerValue(shadow.NgfwDevice)
	o.Name = types.StringPointerValue(shadow.Name)

	return nil
}
func (o AddressGroupDeviceGroupLocation) MarshalJSON() ([]byte, error) {
	obj := struct {
		PanoramaDevice *string `json:"panorama_device"`
		Name           *string `json:"name"`
	}{
		PanoramaDevice: o.PanoramaDevice.ValueStringPointer(),
		Name:           o.Name.ValueStringPointer(),
	}

	return json.Marshal(obj)
}

func (o *AddressGroupDeviceGroupLocation) UnmarshalJSON(data []byte) error {
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
func (o AddressGroupLocation) MarshalJSON() ([]byte, error) {
	obj := struct {
		Shared      *bool                            `json:"shared"`
		Vsys        *AddressGroupVsysLocation        `json:"vsys"`
		DeviceGroup *AddressGroupDeviceGroupLocation `json:"device_group"`
	}{
		Shared:      o.Shared.ValueBoolPointer(),
		Vsys:        o.Vsys,
		DeviceGroup: o.DeviceGroup,
	}

	return json.Marshal(obj)
}

func (o *AddressGroupLocation) UnmarshalJSON(data []byte) error {
	var shadow struct {
		Shared      *bool                            `json:"shared"`
		Vsys        *AddressGroupVsysLocation        `json:"vsys"`
		DeviceGroup *AddressGroupDeviceGroupLocation `json:"device_group"`
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
