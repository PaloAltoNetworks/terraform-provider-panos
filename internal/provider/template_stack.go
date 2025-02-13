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
	"github.com/PaloAltoNetworks/pango/panorama/template_stack"
	pangoutil "github.com/PaloAltoNetworks/pango/util"

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
	_ datasource.DataSource              = &TemplateStackDataSource{}
	_ datasource.DataSourceWithConfigure = &TemplateStackDataSource{}
)

func NewTemplateStackDataSource() datasource.DataSource {
	return &TemplateStackDataSource{}
}

type TemplateStackDataSource struct {
	client  *pango.Client
	manager *sdkmanager.EntryObjectManager[*template_stack.Entry, template_stack.Location, *template_stack.Service]
}

type TemplateStackDataSourceFilter struct {
	// TODO: Generate Data Source filter via function
}

type TemplateStackDataSourceModel struct {
	Location        TemplateStackLocation                         `tfsdk:"location"`
	Name            types.String                                  `tfsdk:"name"`
	Description     types.String                                  `tfsdk:"description"`
	Templates       types.List                                    `tfsdk:"templates"`
	Devices         types.List                                    `tfsdk:"devices"`
	DefaultVsys     types.String                                  `tfsdk:"default_vsys"`
	UserGroupSource *TemplateStackDataSourceUserGroupSourceObject `tfsdk:"user_group_source"`
}
type TemplateStackDataSourceUserGroupSourceObject struct {
	MasterDevice types.String `tfsdk:"master_device"`
}

func (o *TemplateStackDataSourceModel) CopyToPango(ctx context.Context, obj **template_stack.Entry, encrypted *map[string]types.String) diag.Diagnostics {
	var diags diag.Diagnostics
	description_value := o.Description.ValueStringPointer()
	templates_pango_entries := make([]string, 0)
	diags.Append(o.Templates.ElementsAs(ctx, &templates_pango_entries, false)...)
	if diags.HasError() {
		return diags
	}
	devices_pango_entries := make([]string, 0)
	diags.Append(o.Devices.ElementsAs(ctx, &devices_pango_entries, false)...)
	if diags.HasError() {
		return diags
	}
	defaultVsys_value := o.DefaultVsys.ValueStringPointer()
	var userGroupSource_entry *template_stack.UserGroupSource
	if o.UserGroupSource != nil {
		if *obj != nil && (*obj).UserGroupSource != nil {
			userGroupSource_entry = (*obj).UserGroupSource
		} else {
			userGroupSource_entry = new(template_stack.UserGroupSource)
		}

		diags.Append(o.UserGroupSource.CopyToPango(ctx, &userGroupSource_entry, encrypted)...)
		if diags.HasError() {
			return diags
		}
	}

	if (*obj) == nil {
		*obj = new(template_stack.Entry)
	}
	(*obj).Name = o.Name.ValueString()
	(*obj).Description = description_value
	(*obj).Templates = templates_pango_entries
	(*obj).Devices = devices_pango_entries
	(*obj).DefaultVsys = defaultVsys_value
	(*obj).UserGroupSource = userGroupSource_entry

	return diags
}
func (o *TemplateStackDataSourceUserGroupSourceObject) CopyToPango(ctx context.Context, obj **template_stack.UserGroupSource, encrypted *map[string]types.String) diag.Diagnostics {
	var diags diag.Diagnostics
	masterDevice_value := o.MasterDevice.ValueStringPointer()

	if (*obj) == nil {
		*obj = new(template_stack.UserGroupSource)
	}
	(*obj).MasterDevice = masterDevice_value

	return diags
}

func (o *TemplateStackDataSourceModel) CopyFromPango(ctx context.Context, obj *template_stack.Entry, encrypted *map[string]types.String) diag.Diagnostics {
	var diags diag.Diagnostics
	var templates_list types.List
	{
		var list_diags diag.Diagnostics
		templates_list, list_diags = types.ListValueFrom(ctx, types.StringType, obj.Templates)
		diags.Append(list_diags...)
	}
	var devices_list types.List
	{
		var list_diags diag.Diagnostics
		devices_list, list_diags = types.ListValueFrom(ctx, types.StringType, obj.Devices)
		diags.Append(list_diags...)
	}
	var userGroupSource_object *TemplateStackDataSourceUserGroupSourceObject
	if obj.UserGroupSource != nil {
		userGroupSource_object = new(TemplateStackDataSourceUserGroupSourceObject)

		diags.Append(userGroupSource_object.CopyFromPango(ctx, obj.UserGroupSource, encrypted)...)
		if diags.HasError() {
			return diags
		}
	}

	var description_value types.String
	if obj.Description != nil {
		description_value = types.StringValue(*obj.Description)
	}
	var defaultVsys_value types.String
	if obj.DefaultVsys != nil {
		defaultVsys_value = types.StringValue(*obj.DefaultVsys)
	}
	o.Name = types.StringValue(obj.Name)
	o.Description = description_value
	o.Templates = templates_list
	o.Devices = devices_list
	o.DefaultVsys = defaultVsys_value
	o.UserGroupSource = userGroupSource_object

	return diags
}

func (o *TemplateStackDataSourceUserGroupSourceObject) CopyFromPango(ctx context.Context, obj *template_stack.UserGroupSource, encrypted *map[string]types.String) diag.Diagnostics {
	var diags diag.Diagnostics

	var masterDevice_value types.String
	if obj.MasterDevice != nil {
		masterDevice_value = types.StringValue(*obj.MasterDevice)
	}
	o.MasterDevice = masterDevice_value

	return diags
}

func (o *TemplateStackDataSourceModel) resourceXpathComponents() ([]string, error) {
	var components []string
	components = append(components, pangoutil.AsEntryXpath(
		[]string{o.Name.ValueString()},
	))
	return components, nil
}

func TemplateStackDataSourceSchema() dsschema.Schema {
	return dsschema.Schema{
		Attributes: map[string]dsschema.Attribute{

			"location": TemplateStackDataSourceLocationSchema(),

			"name": dsschema.StringAttribute{
				Description: "The name of the service.",
				Computed:    false,
				Required:    true,
				Optional:    false,
				Sensitive:   false,
			},

			"default_vsys": dsschema.StringAttribute{
				Description: "Default virtual system",
				Computed:    true,
				Required:    false,
				Optional:    true,
				Sensitive:   false,
			},

			"user_group_source": TemplateStackDataSourceUserGroupSourceSchema(),

			"description": dsschema.StringAttribute{
				Description: "The description.",
				Computed:    true,
				Required:    false,
				Optional:    true,
				Sensitive:   false,
			},

			"templates": dsschema.ListAttribute{
				Description: "List of templates",
				Required:    false,
				Optional:    true,
				Computed:    true,
				Sensitive:   false,
				ElementType: types.StringType,
			},

			"devices": dsschema.ListAttribute{
				Description: "List of devices",
				Required:    false,
				Optional:    true,
				Computed:    true,
				Sensitive:   false,
				ElementType: types.StringType,
			},
		},
	}
}

func (o *TemplateStackDataSourceModel) getTypeFor(name string) attr.Type {
	schema := TemplateStackDataSourceSchema()
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

func TemplateStackDataSourceUserGroupSourceSchema() dsschema.SingleNestedAttribute {
	return dsschema.SingleNestedAttribute{
		Description: "",
		Required:    false,
		Computed:    true,
		Optional:    true,
		Sensitive:   false,
		Attributes: map[string]dsschema.Attribute{

			"master_device": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
				Required:    false,
				Optional:    true,
				Sensitive:   false,
			},
		},
	}
}

func (o *TemplateStackDataSourceUserGroupSourceObject) getTypeFor(name string) attr.Type {
	schema := TemplateStackDataSourceUserGroupSourceSchema()
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

func TemplateStackDataSourceLocationSchema() rsschema.Attribute {
	return TemplateStackLocationSchema()
}

// Metadata returns the data source type name.
func (d *TemplateStackDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_template_stack"
}

// Schema defines the schema for this data source.
func (d *TemplateStackDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = TemplateStackDataSourceSchema()
}

// Configure prepares the struct.
func (d *TemplateStackDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*pango.Client)
	specifier, _, err := template_stack.Versioning(d.client.Versioning())
	if err != nil {
		resp.Diagnostics.AddError("Failed to configure SDK client", err.Error())
		return
	}
	d.manager = sdkmanager.NewEntryObjectManager(d.client, template_stack.NewService(d.client), specifier, template_stack.SpecMatches)
}
func (o *TemplateStackDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var savestate, state TemplateStackDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &savestate)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var location template_stack.Location

	if savestate.Location.Panorama != nil {
		location.Panorama = &template_stack.PanoramaLocation{

			PanoramaDevice: savestate.Location.Panorama.PanoramaDevice.ValueString(),
		}
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource read", map[string]any{
		"resource_name": "panos_template_stack_resource",
		"function":      "Read",
		"name":          savestate.Name.ValueString(),
	})

	components, err := savestate.resourceXpathComponents()
	if err != nil {
		resp.Diagnostics.AddError("Error creating resource xpath", err.Error())
		return
	}

	object, err := o.manager.Read(ctx, location, components)
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
	_ resource.Resource                = &TemplateStackResource{}
	_ resource.ResourceWithConfigure   = &TemplateStackResource{}
	_ resource.ResourceWithImportState = &TemplateStackResource{}
)

func NewTemplateStackResource() resource.Resource {
	if _, found := resourceFuncMap["panos_template_stack"]; !found {
		resourceFuncMap["panos_template_stack"] = resourceFuncs{
			CreateImportId: TemplateStackImportStateCreator,
		}
	}
	return &TemplateStackResource{}
}

type TemplateStackResource struct {
	client  *pango.Client
	manager *sdkmanager.EntryObjectManager[*template_stack.Entry, template_stack.Location, *template_stack.Service]
}

func TemplateStackResourceLocationSchema() rsschema.Attribute {
	return TemplateStackLocationSchema()
}

type TemplateStackResourceModel struct {
	Location        TemplateStackLocation                       `tfsdk:"location"`
	Name            types.String                                `tfsdk:"name"`
	Description     types.String                                `tfsdk:"description"`
	Templates       types.List                                  `tfsdk:"templates"`
	Devices         types.List                                  `tfsdk:"devices"`
	DefaultVsys     types.String                                `tfsdk:"default_vsys"`
	UserGroupSource *TemplateStackResourceUserGroupSourceObject `tfsdk:"user_group_source"`
}
type TemplateStackResourceUserGroupSourceObject struct {
	MasterDevice types.String `tfsdk:"master_device"`
}

func (r *TemplateStackResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
}

// <ResourceSchema>

func TemplateStackResourceSchema() rsschema.Schema {
	return rsschema.Schema{
		Attributes: map[string]rsschema.Attribute{

			"location": TemplateStackResourceLocationSchema(),

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
				Description: "List of templates",
				Required:    false,
				Optional:    true,
				Computed:    false,
				Sensitive:   false,
				ElementType: types.StringType,
			},

			"devices": rsschema.ListAttribute{
				Description: "List of devices",
				Required:    false,
				Optional:    true,
				Computed:    false,
				Sensitive:   false,
				ElementType: types.StringType,
			},

			"default_vsys": rsschema.StringAttribute{
				Description: "Default virtual system",
				Computed:    false,
				Required:    false,
				Optional:    true,
				Sensitive:   false,
			},

			"user_group_source": TemplateStackResourceUserGroupSourceSchema(),
		},
	}
}

func (o *TemplateStackResourceModel) getTypeFor(name string) attr.Type {
	schema := TemplateStackResourceSchema()
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

func TemplateStackResourceUserGroupSourceSchema() rsschema.SingleNestedAttribute {
	return rsschema.SingleNestedAttribute{
		Description: "",
		Required:    false,
		Computed:    false,
		Optional:    true,
		Sensitive:   false,
		Attributes: map[string]rsschema.Attribute{

			"master_device": rsschema.StringAttribute{
				Description: "",
				Computed:    false,
				Required:    false,
				Optional:    true,
				Sensitive:   false,
			},
		},
	}
}

func (o *TemplateStackResourceUserGroupSourceObject) getTypeFor(name string) attr.Type {
	schema := TemplateStackResourceUserGroupSourceSchema()
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

func (r *TemplateStackResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_template_stack"
}

func (r *TemplateStackResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = TemplateStackResourceSchema()
}

// </ResourceSchema>

func (r *TemplateStackResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*pango.Client)
	specifier, _, err := template_stack.Versioning(r.client.Versioning())
	if err != nil {
		resp.Diagnostics.AddError("Failed to configure SDK client", err.Error())
		return
	}
	r.manager = sdkmanager.NewEntryObjectManager(r.client, template_stack.NewService(r.client), specifier, template_stack.SpecMatches)
}

func (o *TemplateStackResourceModel) CopyToPango(ctx context.Context, obj **template_stack.Entry, encrypted *map[string]types.String) diag.Diagnostics {
	var diags diag.Diagnostics
	description_value := o.Description.ValueStringPointer()
	templates_pango_entries := make([]string, 0)
	diags.Append(o.Templates.ElementsAs(ctx, &templates_pango_entries, false)...)
	if diags.HasError() {
		return diags
	}
	devices_pango_entries := make([]string, 0)
	diags.Append(o.Devices.ElementsAs(ctx, &devices_pango_entries, false)...)
	if diags.HasError() {
		return diags
	}
	defaultVsys_value := o.DefaultVsys.ValueStringPointer()
	var userGroupSource_entry *template_stack.UserGroupSource
	if o.UserGroupSource != nil {
		if *obj != nil && (*obj).UserGroupSource != nil {
			userGroupSource_entry = (*obj).UserGroupSource
		} else {
			userGroupSource_entry = new(template_stack.UserGroupSource)
		}

		diags.Append(o.UserGroupSource.CopyToPango(ctx, &userGroupSource_entry, encrypted)...)
		if diags.HasError() {
			return diags
		}
	}

	if (*obj) == nil {
		*obj = new(template_stack.Entry)
	}
	(*obj).Name = o.Name.ValueString()
	(*obj).Description = description_value
	(*obj).Templates = templates_pango_entries
	(*obj).Devices = devices_pango_entries
	(*obj).DefaultVsys = defaultVsys_value
	(*obj).UserGroupSource = userGroupSource_entry

	return diags
}
func (o *TemplateStackResourceUserGroupSourceObject) CopyToPango(ctx context.Context, obj **template_stack.UserGroupSource, encrypted *map[string]types.String) diag.Diagnostics {
	var diags diag.Diagnostics
	masterDevice_value := o.MasterDevice.ValueStringPointer()

	if (*obj) == nil {
		*obj = new(template_stack.UserGroupSource)
	}
	(*obj).MasterDevice = masterDevice_value

	return diags
}

func (o *TemplateStackResourceModel) CopyFromPango(ctx context.Context, obj *template_stack.Entry, encrypted *map[string]types.String) diag.Diagnostics {
	var diags diag.Diagnostics
	var templates_list types.List
	{
		var list_diags diag.Diagnostics
		templates_list, list_diags = types.ListValueFrom(ctx, types.StringType, obj.Templates)
		diags.Append(list_diags...)
	}
	var devices_list types.List
	{
		var list_diags diag.Diagnostics
		devices_list, list_diags = types.ListValueFrom(ctx, types.StringType, obj.Devices)
		diags.Append(list_diags...)
	}
	var userGroupSource_object *TemplateStackResourceUserGroupSourceObject
	if obj.UserGroupSource != nil {
		userGroupSource_object = new(TemplateStackResourceUserGroupSourceObject)

		diags.Append(userGroupSource_object.CopyFromPango(ctx, obj.UserGroupSource, encrypted)...)
		if diags.HasError() {
			return diags
		}
	}

	var description_value types.String
	if obj.Description != nil {
		description_value = types.StringValue(*obj.Description)
	}
	var defaultVsys_value types.String
	if obj.DefaultVsys != nil {
		defaultVsys_value = types.StringValue(*obj.DefaultVsys)
	}
	o.Name = types.StringValue(obj.Name)
	o.Description = description_value
	o.Templates = templates_list
	o.Devices = devices_list
	o.DefaultVsys = defaultVsys_value
	o.UserGroupSource = userGroupSource_object

	return diags
}

func (o *TemplateStackResourceUserGroupSourceObject) CopyFromPango(ctx context.Context, obj *template_stack.UserGroupSource, encrypted *map[string]types.String) diag.Diagnostics {
	var diags diag.Diagnostics

	var masterDevice_value types.String
	if obj.MasterDevice != nil {
		masterDevice_value = types.StringValue(*obj.MasterDevice)
	}
	o.MasterDevice = masterDevice_value

	return diags
}

func (o *TemplateStackResourceModel) resourceXpathComponents() ([]string, error) {
	var components []string
	components = append(components, pangoutil.AsEntryXpath(
		[]string{o.Name.ValueString()},
	))
	return components, nil
}

func (r *TemplateStackResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state TemplateStackResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource create", map[string]any{
		"resource_name": "panos_template_stack_resource",
		"function":      "Create",
		"name":          state.Name.ValueString(),
	})

	// Verify mode.
	if r.client.Hostname == "" {
		resp.Diagnostics.AddError("Invalid mode error", InspectionModeError)
		return
	}

	// Determine the location.

	var location template_stack.Location

	if state.Location.Panorama != nil {
		location.Panorama = &template_stack.PanoramaLocation{

			PanoramaDevice: state.Location.Panorama.PanoramaDevice.ValueString(),
		}
	}

	if err := location.IsValid(); err != nil {
		resp.Diagnostics.AddError("Invalid location", err.Error())
		return
	}

	// Load the desired config.
	var obj *template_stack.Entry

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

	components, err := state.resourceXpathComponents()
	if err != nil {
		resp.Diagnostics.AddError("Error creating resource xpath", err.Error())
		return
	}
	created, err := r.manager.Create(ctx, location, components, obj)
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
func (o *TemplateStackResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var savestate, state TemplateStackResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &savestate)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var location template_stack.Location

	if savestate.Location.Panorama != nil {
		location.Panorama = &template_stack.PanoramaLocation{

			PanoramaDevice: savestate.Location.Panorama.PanoramaDevice.ValueString(),
		}
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource read", map[string]any{
		"resource_name": "panos_template_stack_resource",
		"function":      "Read",
		"name":          savestate.Name.ValueString(),
	})

	components, err := savestate.resourceXpathComponents()
	if err != nil {
		resp.Diagnostics.AddError("Error creating resource xpath", err.Error())
		return
	}

	object, err := o.manager.Read(ctx, location, components)
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
func (r *TemplateStackResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var plan, state TemplateStackResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var location template_stack.Location

	if state.Location.Panorama != nil {
		location.Panorama = &template_stack.PanoramaLocation{

			PanoramaDevice: state.Location.Panorama.PanoramaDevice.ValueString(),
		}
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource update", map[string]any{
		"resource_name": "panos_template_stack_resource",
		"function":      "Update",
	})

	// Verify mode.
	if r.client.Hostname == "" {
		resp.Diagnostics.AddError("Invalid mode error", InspectionModeError)
		return
	}

	components, err := state.resourceXpathComponents()
	if err != nil {
		resp.Diagnostics.AddError("Error creating resource xpath", err.Error())
		return
	}

	obj, err := r.manager.Read(ctx, location, components)
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
func (r *TemplateStackResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var state TemplateStackResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource delete", map[string]any{
		"resource_name": "panos_template_stack_resource",
		"function":      "Delete",
		"name":          state.Name.ValueString(),
	})

	// Verify mode.
	if r.client.Hostname == "" {
		resp.Diagnostics.AddError("Invalid mode error", InspectionModeError)
		return
	}

	var location template_stack.Location

	if state.Location.Panorama != nil {
		location.Panorama = &template_stack.PanoramaLocation{

			PanoramaDevice: state.Location.Panorama.PanoramaDevice.ValueString(),
		}
	}

	err := r.manager.Delete(ctx, location, []string{state.Name.ValueString()})
	if err != nil && !errors.Is(err, sdkmanager.ErrObjectNotFound) {
		resp.Diagnostics.AddError("Error in delete", err.Error())
	}

}

type TemplateStackImportState struct {
	Location TemplateStackLocation `json:"location"`
	Name     string                `json:"name"`
}

func TemplateStackImportStateCreator(ctx context.Context, resource types.Object) ([]byte, error) {
	attrs := resource.Attributes()
	if attrs == nil {
		return nil, fmt.Errorf("Object has no attributes")
	}

	locationAttr, ok := attrs["location"]
	if !ok {
		return nil, fmt.Errorf("location attribute missing")
	}

	var location TemplateStackLocation
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

	importStruct := TemplateStackImportState{
		Location: location,
		Name:     name,
	}

	return json.Marshal(importStruct)
}

func (r *TemplateStackResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	var obj TemplateStackImportState
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

type TemplateStackPanoramaLocation struct {
	PanoramaDevice types.String `tfsdk:"panorama_device"`
}
type TemplateStackLocation struct {
	Panorama *TemplateStackPanoramaLocation `tfsdk:"panorama"`
}

func TemplateStackLocationSchema() rsschema.Attribute {
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

func (o TemplateStackPanoramaLocation) MarshalJSON() ([]byte, error) {
	obj := struct {
		PanoramaDevice *string `json:"panorama_device"`
	}{
		PanoramaDevice: o.PanoramaDevice.ValueStringPointer(),
	}

	return json.Marshal(obj)
}

func (o *TemplateStackPanoramaLocation) UnmarshalJSON(data []byte) error {
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
func (o TemplateStackLocation) MarshalJSON() ([]byte, error) {
	obj := struct {
		Panorama *TemplateStackPanoramaLocation `json:"panorama"`
	}{
		Panorama: o.Panorama,
	}

	return json.Marshal(obj)
}

func (o *TemplateStackLocation) UnmarshalJSON(data []byte) error {
	var shadow struct {
		Panorama *TemplateStackPanoramaLocation `json:"panorama"`
	}

	err := json.Unmarshal(data, &shadow)
	if err != nil {
		return err
	}
	o.Panorama = shadow.Panorama

	return nil
}
