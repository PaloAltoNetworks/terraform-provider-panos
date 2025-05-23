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
	"github.com/PaloAltoNetworks/pango/panorama/devicegroup"

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
	_ datasource.DataSource              = &DeviceGroupDataSource{}
	_ datasource.DataSourceWithConfigure = &DeviceGroupDataSource{}
)

func NewDeviceGroupDataSource() datasource.DataSource {
	return &DeviceGroupDataSource{}
}

type DeviceGroupDataSource struct {
	client  *pango.Client
	manager *sdkmanager.EntryObjectManager[*devicegroup.Entry, devicegroup.Location, *devicegroup.Service]
}

type DeviceGroupDataSourceFilter struct {
	// TODO: Generate Data Source filter via function
}

type DeviceGroupDataSourceModel struct {
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

func (o *DeviceGroupDataSourceModel) CopyToPango(ctx context.Context, obj **devicegroup.Entry, encrypted *map[string]types.String) diag.Diagnostics {
	var diags diag.Diagnostics
	description_value := o.Description.ValueStringPointer()
	templates_pango_entries := make([]string, 0)
	diags.Append(o.Templates.ElementsAs(ctx, &templates_pango_entries, false)...)
	if diags.HasError() {
		return diags
	}
	var devices_tf_entries []DeviceGroupDataSourceDevicesObject
	var devices_pango_entries []devicegroup.Devices
	{
		d := o.Devices.ElementsAs(ctx, &devices_tf_entries, false)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}
		for _, elt := range devices_tf_entries {
			var entry *devicegroup.Devices
			diags.Append(elt.CopyToPango(ctx, &entry, encrypted)...)
			if diags.HasError() {
				return diags
			}
			devices_pango_entries = append(devices_pango_entries, *entry)
		}
	}
	authorizationCode_value := o.AuthorizationCode.ValueStringPointer()

	if (*obj) == nil {
		*obj = new(devicegroup.Entry)
	}
	(*obj).Name = o.Name.ValueString()
	(*obj).Description = description_value
	(*obj).Templates = templates_pango_entries
	(*obj).Devices = devices_pango_entries
	(*obj).AuthorizationCode = authorizationCode_value

	return diags
}
func (o *DeviceGroupDataSourceDevicesObject) CopyToPango(ctx context.Context, obj **devicegroup.Devices, encrypted *map[string]types.String) diag.Diagnostics {
	var diags diag.Diagnostics
	vsys_pango_entries := make([]string, 0)
	diags.Append(o.Vsys.ElementsAs(ctx, &vsys_pango_entries, false)...)
	if diags.HasError() {
		return diags
	}

	if (*obj) == nil {
		*obj = new(devicegroup.Devices)
	}
	(*obj).Name = o.Name.ValueString()
	(*obj).Vsys = vsys_pango_entries

	return diags
}

func (o *DeviceGroupDataSourceModel) CopyFromPango(ctx context.Context, obj *devicegroup.Entry, encrypted *map[string]types.String) diag.Diagnostics {
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

func (o *DeviceGroupDataSourceDevicesObject) CopyFromPango(ctx context.Context, obj *devicegroup.Devices, encrypted *map[string]types.String) diag.Diagnostics {
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

			"templates": dsschema.ListAttribute{
				Description: "List of reference templates",
				Required:    false,
				Optional:    true,
				Computed:    true,
				Sensitive:   false,
				ElementType: types.StringType,
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

	providerData := req.ProviderData.(*ProviderData)
	d.client = providerData.Client
	specifier, _, err := devicegroup.Versioning(d.client.Versioning())
	if err != nil {
		resp.Diagnostics.AddError("Failed to configure SDK client", err.Error())
		return
	}
	batchSize := providerData.MultiConfigBatchSize
	d.manager = sdkmanager.NewEntryObjectManager(d.client, devicegroup.NewService(d.client), batchSize, specifier, devicegroup.SpecMatches)
}
func (o *DeviceGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var savestate, state DeviceGroupDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &savestate)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var location devicegroup.Location

	if savestate.Location.Panorama != nil {
		location.Panorama = &devicegroup.PanoramaLocation{

			PanoramaDevice: savestate.Location.Panorama.PanoramaDevice.ValueString(),
		}
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource read", map[string]any{
		"resource_name": "panos_device_group_resource",
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
	_ resource.Resource                = &DeviceGroupResource{}
	_ resource.ResourceWithConfigure   = &DeviceGroupResource{}
	_ resource.ResourceWithImportState = &DeviceGroupResource{}
)

func NewDeviceGroupResource() resource.Resource {
	if _, found := resourceFuncMap["panos_device_group"]; !found {
		resourceFuncMap["panos_device_group"] = resourceFuncs{
			CreateImportId: DeviceGroupImportStateCreator,
		}
	}
	return &DeviceGroupResource{}
}

type DeviceGroupResource struct {
	client  *pango.Client
	manager *sdkmanager.EntryObjectManager[*devicegroup.Entry, devicegroup.Location, *devicegroup.Service]
}

func DeviceGroupResourceLocationSchema() rsschema.Attribute {
	return DeviceGroupLocationSchema()
}

type DeviceGroupResourceModel struct {
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

func (r *DeviceGroupResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
}

// <ResourceSchema>

func DeviceGroupResourceSchema() rsschema.Schema {
	return rsschema.Schema{
		Attributes: map[string]rsschema.Attribute{

			"location": DeviceGroupResourceLocationSchema(),

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

func (r *DeviceGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device_group"
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

	providerData := req.ProviderData.(*ProviderData)
	r.client = providerData.Client
	specifier, _, err := devicegroup.Versioning(r.client.Versioning())
	if err != nil {
		resp.Diagnostics.AddError("Failed to configure SDK client", err.Error())
		return
	}
	batchSize := providerData.MultiConfigBatchSize
	r.manager = sdkmanager.NewEntryObjectManager(r.client, devicegroup.NewService(r.client), batchSize, specifier, devicegroup.SpecMatches)
}

func (o *DeviceGroupResourceModel) CopyToPango(ctx context.Context, obj **devicegroup.Entry, encrypted *map[string]types.String) diag.Diagnostics {
	var diags diag.Diagnostics
	description_value := o.Description.ValueStringPointer()
	templates_pango_entries := make([]string, 0)
	diags.Append(o.Templates.ElementsAs(ctx, &templates_pango_entries, false)...)
	if diags.HasError() {
		return diags
	}
	var devices_tf_entries []DeviceGroupResourceDevicesObject
	var devices_pango_entries []devicegroup.Devices
	{
		d := o.Devices.ElementsAs(ctx, &devices_tf_entries, false)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}
		for _, elt := range devices_tf_entries {
			var entry *devicegroup.Devices
			diags.Append(elt.CopyToPango(ctx, &entry, encrypted)...)
			if diags.HasError() {
				return diags
			}
			devices_pango_entries = append(devices_pango_entries, *entry)
		}
	}
	authorizationCode_value := o.AuthorizationCode.ValueStringPointer()

	if (*obj) == nil {
		*obj = new(devicegroup.Entry)
	}
	(*obj).Name = o.Name.ValueString()
	(*obj).Description = description_value
	(*obj).Templates = templates_pango_entries
	(*obj).Devices = devices_pango_entries
	(*obj).AuthorizationCode = authorizationCode_value

	return diags
}
func (o *DeviceGroupResourceDevicesObject) CopyToPango(ctx context.Context, obj **devicegroup.Devices, encrypted *map[string]types.String) diag.Diagnostics {
	var diags diag.Diagnostics
	vsys_pango_entries := make([]string, 0)
	diags.Append(o.Vsys.ElementsAs(ctx, &vsys_pango_entries, false)...)
	if diags.HasError() {
		return diags
	}

	if (*obj) == nil {
		*obj = new(devicegroup.Devices)
	}
	(*obj).Name = o.Name.ValueString()
	(*obj).Vsys = vsys_pango_entries

	return diags
}

func (o *DeviceGroupResourceModel) CopyFromPango(ctx context.Context, obj *devicegroup.Entry, encrypted *map[string]types.String) diag.Diagnostics {
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

func (o *DeviceGroupResourceDevicesObject) CopyFromPango(ctx context.Context, obj *devicegroup.Devices, encrypted *map[string]types.String) diag.Diagnostics {
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

	var location devicegroup.Location

	if state.Location.Panorama != nil {
		location.Panorama = &devicegroup.PanoramaLocation{

			PanoramaDevice: state.Location.Panorama.PanoramaDevice.ValueString(),
		}
	}

	if err := location.IsValid(); err != nil {
		resp.Diagnostics.AddError("Invalid location", err.Error())
		return
	}

	// Load the desired config.
	var obj *devicegroup.Entry

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
func (o *DeviceGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var savestate, state DeviceGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &savestate)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var location devicegroup.Location

	if savestate.Location.Panorama != nil {
		location.Panorama = &devicegroup.PanoramaLocation{

			PanoramaDevice: savestate.Location.Panorama.PanoramaDevice.ValueString(),
		}
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource read", map[string]any{
		"resource_name": "panos_device_group_resource",
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
func (r *DeviceGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var plan, state DeviceGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var location devicegroup.Location

	if state.Location.Panorama != nil {
		location.Panorama = &devicegroup.PanoramaLocation{

			PanoramaDevice: state.Location.Panorama.PanoramaDevice.ValueString(),
		}
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource update", map[string]any{
		"resource_name": "panos_device_group_resource",
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
func (r *DeviceGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var state DeviceGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource delete", map[string]any{
		"resource_name": "panos_device_group_resource",
		"function":      "Delete",
		"name":          state.Name.ValueString(),
	})

	// Verify mode.
	if r.client.Hostname == "" {
		resp.Diagnostics.AddError("Invalid mode error", InspectionModeError)
		return
	}

	var location devicegroup.Location

	if state.Location.Panorama != nil {
		location.Panorama = &devicegroup.PanoramaLocation{

			PanoramaDevice: state.Location.Panorama.PanoramaDevice.ValueString(),
		}
	}

	err := r.manager.Delete(ctx, location, []string{state.Name.ValueString()})
	if err != nil && !errors.Is(err, sdkmanager.ErrObjectNotFound) {
		resp.Diagnostics.AddError("Error in delete", err.Error())
	}

}

type DeviceGroupImportState struct {
	Location DeviceGroupLocation `json:"location"`
	Name     string              `json:"name"`
}

func DeviceGroupImportStateCreator(ctx context.Context, resource types.Object) ([]byte, error) {
	attrs := resource.Attributes()
	if attrs == nil {
		return nil, fmt.Errorf("Object has no attributes")
	}

	locationAttr, ok := attrs["location"]
	if !ok {
		return nil, fmt.Errorf("location attribute missing")
	}

	var location DeviceGroupLocation
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

	importStruct := DeviceGroupImportState{
		Location: location,
		Name:     name,
	}

	return json.Marshal(importStruct)
}

func (r *DeviceGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	var obj DeviceGroupImportState
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

func (o DeviceGroupPanoramaLocation) MarshalJSON() ([]byte, error) {
	obj := struct {
		PanoramaDevice *string `json:"panorama_device"`
	}{
		PanoramaDevice: o.PanoramaDevice.ValueStringPointer(),
	}

	return json.Marshal(obj)
}

func (o *DeviceGroupPanoramaLocation) UnmarshalJSON(data []byte) error {
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
func (o DeviceGroupLocation) MarshalJSON() ([]byte, error) {
	obj := struct {
		Panorama *DeviceGroupPanoramaLocation `json:"panorama"`
	}{
		Panorama: o.Panorama,
	}

	return json.Marshal(obj)
}

func (o *DeviceGroupLocation) UnmarshalJSON(data []byte) error {
	var shadow struct {
		Panorama *DeviceGroupPanoramaLocation `json:"panorama"`
	}

	err := json.Unmarshal(data, &shadow)
	if err != nil {
		return err
	}
	o.Panorama = shadow.Panorama

	return nil
}
