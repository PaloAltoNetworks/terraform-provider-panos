package provider

// Note:  This file is automatically generated.  Manually made changes
// will be overwritten when the provider is generated.

import (
	"context"
	"errors"
	"fmt"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/panorama/device_group"

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
	"github.com/hashicorp/terraform-plugin-log/tflog"

	sdkmanager "github.com/PaloAltoNetworks/terraform-provider-panos/internal/manager"
)

// Generate Terraform Data Source object.
var (
	_ datasource.DataSource              = &DeviceGroupDataSource{}
	_ datasource.DataSourceWithConfigure = &DeviceGroupDataSource{}
)

func NewDeviceGroupDataSource() datasource.DataSource {
	return &DeviceGroupDataSource{}
}

type DeviceGroupDataSource struct {
	client  *pango.Client
	manager *sdkmanager.EntryObjectManager[*device_group.Entry, device_group.Location, *device_group.Service]
}

type DeviceGroupDataSourceFilter struct {
	// TODO: Generate Data Source filter via function
}
type DeviceGroupDataSourceTfid struct {
	Name     string                `json:"name"`
	Location device_group.Location `json:"location"`
}

func (o *DeviceGroupDataSourceTfid) IsValid() error {
	if o.Name == "" {
		return fmt.Errorf("name is unspecified")
	}
	return o.Location.IsValid()
}

type DeviceGroupDataSourceModel struct {
	Tfid              types.String        `tfsdk:"tfid"`
	Location          DeviceGroupLocation `tfsdk:"location"`
	Name              types.String        `tfsdk:"name"`
	Description       types.String        `tfsdk:"description"`
	Templates         types.List          `tfsdk:"templates"`
	Devices           types.List          `tfsdk:"devices"`
	AuthorizationCode types.String        `tfsdk:"authorization_code"`
}
type DeviceGroupDataSourceDevicesObject struct {
	Name types.String `tfsdk:"name"`
	Vsys types.List   `tfsdk:"vsys"`
}

func (o *DeviceGroupDataSourceModel) CopyToPango(ctx context.Context, obj **device_group.Entry, encrypted *map[string]types.String) diag.Diagnostics {
	var diags diag.Diagnostics
	authorizationCode_value := o.AuthorizationCode.ValueStringPointer()
	description_value := o.Description.ValueStringPointer()
	templates_pango_entries := make([]string, 0)
	diags.Append(o.Templates.ElementsAs(ctx, &templates_pango_entries, false)...)
	if diags.HasError() {
		return diags
	}
	var devices_tf_entries []DeviceGroupDataSourceDevicesObject
	var devices_pango_entries []device_group.Devices
	{
		d := o.Devices.ElementsAs(ctx, &devices_tf_entries, false)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}
		for _, elt := range devices_tf_entries {
			var entry *device_group.Devices
			diags.Append(elt.CopyToPango(ctx, &entry, encrypted)...)
			if diags.HasError() {
				return diags
			}
			devices_pango_entries = append(devices_pango_entries, *entry)
		}
	}

	if (*obj) == nil {
		*obj = new(device_group.Entry)
	}
	(*obj).Name = o.Name.ValueString()
	(*obj).AuthorizationCode = authorizationCode_value
	(*obj).Description = description_value
	(*obj).Templates = templates_pango_entries
	(*obj).Devices = devices_pango_entries

	return diags
}
func (o *DeviceGroupDataSourceDevicesObject) CopyToPango(ctx context.Context, obj **device_group.Devices, encrypted *map[string]types.String) diag.Diagnostics {
	var diags diag.Diagnostics
	vsys_pango_entries := make([]string, 0)
	diags.Append(o.Vsys.ElementsAs(ctx, &vsys_pango_entries, false)...)
	if diags.HasError() {
		return diags
	}

	if (*obj) == nil {
		*obj = new(device_group.Devices)
	}
	(*obj).Name = o.Name.ValueString()
	(*obj).Vsys = vsys_pango_entries

	return diags
}

func (o *DeviceGroupDataSourceModel) CopyFromPango(ctx context.Context, obj *device_group.Entry, encrypted *map[string]types.String) diag.Diagnostics {
	var diags diag.Diagnostics
	var templates_list types.List
	{
		var list_diags diag.Diagnostics
		templates_list, list_diags = types.ListValueFrom(ctx, types.StringType, obj.Templates)
		diags.Append(list_diags...)
	}
	var devices_list types.List
	{
		var devices_tf_entries []DeviceGroupDataSourceDevicesObject
		for _, elt := range obj.Devices {
			var entry DeviceGroupDataSourceDevicesObject
			entry_diags := entry.CopyFromPango(ctx, &elt, encrypted)
			diags.Append(entry_diags...)
			devices_tf_entries = append(devices_tf_entries, entry)
		}
		var list_diags diag.Diagnostics
		schemaType := o.getTypeFor("devices")
		devices_list, list_diags = types.ListValueFrom(ctx, schemaType, devices_tf_entries)
		diags.Append(list_diags...)
	}
	var authorizationCode_value types.String
	if obj.AuthorizationCode != nil {
		authorizationCode_value = types.StringValue(*obj.AuthorizationCode)
	}
	var description_value types.String
	if obj.Description != nil {
		description_value = types.StringValue(*obj.Description)
	}
	o.Name = types.StringValue(obj.Name)
	o.Templates = templates_list
	o.Devices = devices_list
	o.AuthorizationCode = authorizationCode_value
	o.Description = description_value

	return diags
}

func (o *DeviceGroupDataSourceDevicesObject) CopyFromPango(ctx context.Context, obj *device_group.Devices, encrypted *map[string]types.String) diag.Diagnostics {
	var diags diag.Diagnostics
	var vsys_list types.List
	{
		var list_diags diag.Diagnostics
		vsys_list, list_diags = types.ListValueFrom(ctx, types.StringType, obj.Vsys)
		diags.Append(list_diags...)
	}
	o.Name = types.StringValue(obj.Name)
	o.Vsys = vsys_list

	return diags
}

func DeviceGroupDataSourceSchema() dsschema.Schema {
	return dsschema.Schema{
		Attributes: map[string]dsschema.Attribute{

			"location": DeviceGroupDataSourceLocationSchema(),

			"tfid": dsschema.StringAttribute{
				Description: "The Terraform ID.",
				Computed:    true,
				Required:    false,
				Optional:    false,
				Sensitive:   false,
			},

			"name": dsschema.StringAttribute{
				Description: "The name of the service.",
				Computed:    false,
				Required:    true,
				Optional:    false,
				Sensitive:   false,
			},

			"devices": dsschema.ListNestedAttribute{
				Description:  "List of devices",
				Required:     false,
				Optional:     true,
				Computed:     true,
				Sensitive:    false,
				NestedObject: DeviceGroupDataSourceDevicesSchema(),
			},

			"authorization_code": dsschema.StringAttribute{
				Description: "Authorization code",
				Computed:    true,
				Required:    false,
				Optional:    true,
				Sensitive:   false,
			},

			"description": dsschema.StringAttribute{
				Description: "The description.",
				Computed:    true,
				Required:    false,
				Optional:    true,
				Sensitive:   false,
			},

			"templates": dsschema.ListAttribute{
				Description: "List of reference templates",
				Required:    false,
				Optional:    true,
				Computed:    true,
				Sensitive:   false,
				ElementType: types.StringType,
			},
		},
	}
}

func (o *DeviceGroupDataSourceModel) getTypeFor(name string) attr.Type {
	schema := DeviceGroupDataSourceSchema()
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

func DeviceGroupDataSourceDevicesSchema() dsschema.NestedAttributeObject {
	return dsschema.NestedAttributeObject{
		Attributes: map[string]dsschema.Attribute{

			"name": dsschema.StringAttribute{
				Description: "",
				Computed:    false,
				Required:    true,
				Optional:    false,
				Sensitive:   false,
			},

			"vsys": dsschema.ListAttribute{
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

func (o *DeviceGroupDataSourceDevicesObject) getTypeFor(name string) attr.Type {
	schema := DeviceGroupDataSourceDevicesSchema()
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

func DeviceGroupDataSourceLocationSchema() rsschema.Attribute {
	return DeviceGroupLocationSchema()
}

// Metadata returns the data source type name.
func (d *DeviceGroupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device_group"
}

// Schema defines the schema for this data source.
func (d *DeviceGroupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DeviceGroupDataSourceSchema()
}

// Configure prepares the struct.
func (d *DeviceGroupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*pango.Client)
	specifier, _, err := device_group.Versioning(d.client.Versioning())
	if err != nil {
		resp.Diagnostics.AddError("Failed to configure SDK client", err.Error())
		return
	}
	d.manager = sdkmanager.NewEntryObjectManager(d.client, device_group.NewService(d.client), specifier, device_group.SpecMatches)
}

func (o *DeviceGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var savestate, state DeviceGroupDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &savestate)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var loc DeviceGroupDataSourceTfid
	loc.Name = *savestate.Name.ValueStringPointer()

	if savestate.Location.Panorama != nil {
		loc.Location.Panorama = &device_group.PanoramaLocation{

			PanoramaDevice: savestate.Location.Panorama.PanoramaDevice.ValueString(),
		}
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource read", map[string]any{
		"resource_name": "panos_device_group_resource",
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
	_ resource.Resource                = &DeviceGroupResource{}
	_ resource.ResourceWithConfigure   = &DeviceGroupResource{}
	_ resource.ResourceWithImportState = &DeviceGroupResource{}
)

func NewDeviceGroupResource() resource.Resource {
	return &DeviceGroupResource{}
}

type DeviceGroupResource struct {
	client  *pango.Client
	manager *sdkmanager.EntryObjectManager[*device_group.Entry, device_group.Location, *device_group.Service]
}
type DeviceGroupResourceTfid struct {
	Name     string                `json:"name"`
	Location device_group.Location `json:"location"`
}

func (o *DeviceGroupResourceTfid) IsValid() error {
	if o.Name == "" {
		return fmt.Errorf("name is unspecified")
	}
	return o.Location.IsValid()
}

func DeviceGroupResourceLocationSchema() rsschema.Attribute {
	return DeviceGroupLocationSchema()
}

type DeviceGroupResourceModel struct {
	Tfid              types.String        `tfsdk:"tfid"`
	Location          DeviceGroupLocation `tfsdk:"location"`
	Name              types.String        `tfsdk:"name"`
	Description       types.String        `tfsdk:"description"`
	Templates         types.List          `tfsdk:"templates"`
	Devices           types.List          `tfsdk:"devices"`
	AuthorizationCode types.String        `tfsdk:"authorization_code"`
}
type DeviceGroupResourceDevicesObject struct {
	Name types.String `tfsdk:"name"`
	Vsys types.List   `tfsdk:"vsys"`
}

func (r *DeviceGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device_group"
}

func (r *DeviceGroupResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
}

// <ResourceSchema>

func DeviceGroupResourceSchema() rsschema.Schema {
	return rsschema.Schema{
		Attributes: map[string]rsschema.Attribute{

			"location": DeviceGroupResourceLocationSchema(),

			"tfid": rsschema.StringAttribute{
				Description: "The Terraform ID.",
				Computed:    true,
				Required:    false,
				Optional:    false,
				Sensitive:   false,
			},

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

			"templates": rsschema.ListAttribute{
				Description: "List of reference templates",
				Required:    false,
				Optional:    true,
				Computed:    false,
				Sensitive:   false,
				ElementType: types.StringType,
			},

			"devices": rsschema.ListNestedAttribute{
				Description:  "List of devices",
				Required:     false,
				Optional:     true,
				Computed:     false,
				Sensitive:    false,
				NestedObject: DeviceGroupResourceDevicesSchema(),
			},

			"authorization_code": rsschema.StringAttribute{
				Description: "Authorization code",
				Computed:    false,
				Required:    false,
				Optional:    true,
				Sensitive:   false,
			},
		},
	}
}

func (o *DeviceGroupResourceModel) getTypeFor(name string) attr.Type {
	schema := DeviceGroupResourceSchema()
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

func DeviceGroupResourceDevicesSchema() rsschema.NestedAttributeObject {
	return rsschema.NestedAttributeObject{
		Attributes: map[string]rsschema.Attribute{

			"name": rsschema.StringAttribute{
				Description: "",
				Computed:    false,
				Required:    true,
				Optional:    false,
				Sensitive:   false,
			},

			"vsys": rsschema.ListAttribute{
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

func (o *DeviceGroupResourceDevicesObject) getTypeFor(name string) attr.Type {
	schema := DeviceGroupResourceDevicesSchema()
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

func (r *DeviceGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = DeviceGroupResourceSchema()
}

// </ResourceSchema>

func (r *DeviceGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*pango.Client)
	specifier, _, err := device_group.Versioning(r.client.Versioning())
	if err != nil {
		resp.Diagnostics.AddError("Failed to configure SDK client", err.Error())
		return
	}
	r.manager = sdkmanager.NewEntryObjectManager(r.client, device_group.NewService(r.client), specifier, device_group.SpecMatches)
}

func (o *DeviceGroupResourceModel) CopyToPango(ctx context.Context, obj **device_group.Entry, encrypted *map[string]types.String) diag.Diagnostics {
	var diags diag.Diagnostics
	authorizationCode_value := o.AuthorizationCode.ValueStringPointer()
	description_value := o.Description.ValueStringPointer()
	templates_pango_entries := make([]string, 0)
	diags.Append(o.Templates.ElementsAs(ctx, &templates_pango_entries, false)...)
	if diags.HasError() {
		return diags
	}
	var devices_tf_entries []DeviceGroupResourceDevicesObject
	var devices_pango_entries []device_group.Devices
	{
		d := o.Devices.ElementsAs(ctx, &devices_tf_entries, false)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}
		for _, elt := range devices_tf_entries {
			var entry *device_group.Devices
			diags.Append(elt.CopyToPango(ctx, &entry, encrypted)...)
			if diags.HasError() {
				return diags
			}
			devices_pango_entries = append(devices_pango_entries, *entry)
		}
	}

	if (*obj) == nil {
		*obj = new(device_group.Entry)
	}
	(*obj).Name = o.Name.ValueString()
	(*obj).AuthorizationCode = authorizationCode_value
	(*obj).Description = description_value
	(*obj).Templates = templates_pango_entries
	(*obj).Devices = devices_pango_entries

	return diags
}
func (o *DeviceGroupResourceDevicesObject) CopyToPango(ctx context.Context, obj **device_group.Devices, encrypted *map[string]types.String) diag.Diagnostics {
	var diags diag.Diagnostics
	vsys_pango_entries := make([]string, 0)
	diags.Append(o.Vsys.ElementsAs(ctx, &vsys_pango_entries, false)...)
	if diags.HasError() {
		return diags
	}

	if (*obj) == nil {
		*obj = new(device_group.Devices)
	}
	(*obj).Name = o.Name.ValueString()
	(*obj).Vsys = vsys_pango_entries

	return diags
}

func (o *DeviceGroupResourceModel) CopyFromPango(ctx context.Context, obj *device_group.Entry, encrypted *map[string]types.String) diag.Diagnostics {
	var diags diag.Diagnostics
	var templates_list types.List
	{
		var list_diags diag.Diagnostics
		templates_list, list_diags = types.ListValueFrom(ctx, types.StringType, obj.Templates)
		diags.Append(list_diags...)
	}
	var devices_list types.List
	{
		var devices_tf_entries []DeviceGroupResourceDevicesObject
		for _, elt := range obj.Devices {
			var entry DeviceGroupResourceDevicesObject
			entry_diags := entry.CopyFromPango(ctx, &elt, encrypted)
			diags.Append(entry_diags...)
			devices_tf_entries = append(devices_tf_entries, entry)
		}
		var list_diags diag.Diagnostics
		schemaType := o.getTypeFor("devices")
		devices_list, list_diags = types.ListValueFrom(ctx, schemaType, devices_tf_entries)
		diags.Append(list_diags...)
	}
	var description_value types.String
	if obj.Description != nil {
		description_value = types.StringValue(*obj.Description)
	}
	var authorizationCode_value types.String
	if obj.AuthorizationCode != nil {
		authorizationCode_value = types.StringValue(*obj.AuthorizationCode)
	}
	o.Name = types.StringValue(obj.Name)
	o.Description = description_value
	o.Templates = templates_list
	o.Devices = devices_list
	o.AuthorizationCode = authorizationCode_value

	return diags
}

func (o *DeviceGroupResourceDevicesObject) CopyFromPango(ctx context.Context, obj *device_group.Devices, encrypted *map[string]types.String) diag.Diagnostics {
	var diags diag.Diagnostics
	var vsys_list types.List
	{
		var list_diags diag.Diagnostics
		vsys_list, list_diags = types.ListValueFrom(ctx, types.StringType, obj.Vsys)
		diags.Append(list_diags...)
	}
	o.Name = types.StringValue(obj.Name)
	o.Vsys = vsys_list

	return diags
}

func (r *DeviceGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state DeviceGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource create", map[string]any{
		"resource_name": "panos_device_group_resource",
		"function":      "Create",
		"name":          state.Name.ValueString(),
	})

	// Verify mode.
	if r.client.Hostname == "" {
		resp.Diagnostics.AddError("Invalid mode error", InspectionModeError)
		return
	}

	// Determine the location.
	loc := DeviceGroupResourceTfid{Name: state.Name.ValueString()}

	// TODO: this needs to handle location structure for UUID style shared has nested structure type

	if state.Location.Panorama != nil {
		loc.Location.Panorama = &device_group.PanoramaLocation{

			PanoramaDevice: state.Location.Panorama.PanoramaDevice.ValueString(),
		}
	}

	if err := loc.IsValid(); err != nil {
		resp.Diagnostics.AddError("Invalid location", err.Error())
		return
	}

	// Load the desired config.
	var obj *device_group.Entry

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

func (o *DeviceGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var savestate, state DeviceGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &savestate)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var loc DeviceGroupResourceTfid
	// Parse the location from tfid.
	if err := DecodeLocation(savestate.Tfid.ValueString(), &loc); err != nil {
		resp.Diagnostics.AddError("Error parsing tfid", err.Error())
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource read", map[string]any{
		"resource_name": "panos_device_group_resource",
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

func (r *DeviceGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var plan, state DeviceGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var loc DeviceGroupResourceTfid
	if err := DecodeLocation(state.Tfid.ValueString(), &loc); err != nil {
		resp.Diagnostics.AddError("Error parsing tfid", err.Error())
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource update", map[string]any{
		"resource_name": "panos_device_group_resource",
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

func (r *DeviceGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var state DeviceGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the location from tfid.
	var loc DeviceGroupResourceTfid
	if err := DecodeLocation(state.Tfid.ValueString(), &loc); err != nil {
		resp.Diagnostics.AddError("error parsing tfid", err.Error())
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource delete", map[string]any{
		"resource_name": "panos_device_group_resource",
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

func (r *DeviceGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("tfid"), req, resp)
}

type DeviceGroupPanoramaLocation struct {
	PanoramaDevice types.String `tfsdk:"panorama_device"`
}
type DeviceGroupLocation struct {
	Panorama *DeviceGroupPanoramaLocation `tfsdk:"panorama"`
}

func DeviceGroupLocationSchema() rsschema.Attribute {
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