package provider

// Note:  This file is automatically generated.  Manually made changes
// will be overwritten when the provider is generated.

import (
	"context"
	"errors"
	"fmt"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/panorama/template_variable"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
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
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	sdkmanager "github.com/PaloAltoNetworks/terraform-provider-panos/internal/manager"
)

// Generate Terraform Data Source object.
var (
	_ datasource.DataSource              = &TemplateVariableDataSource{}
	_ datasource.DataSourceWithConfigure = &TemplateVariableDataSource{}
)

func NewTemplateVariableDataSource() datasource.DataSource {
	return &TemplateVariableDataSource{}
}

type TemplateVariableDataSource struct {
	client  *pango.Client
	manager *sdkmanager.EntryObjectManager[*template_variable.Entry, template_variable.Location, *template_variable.Service]
}

type TemplateVariableDataSourceFilter struct {
	// TODO: Generate Data Source filter via function
}
type TemplateVariableDataSourceTfid struct {
	Name     string                     `json:"name"`
	Location template_variable.Location `json:"location"`
}

func (o *TemplateVariableDataSourceTfid) IsValid() error {
	if o.Name == "" {
		return fmt.Errorf("name is unspecified")
	}
	return o.Location.IsValid()
}

type TemplateVariableDataSourceModel struct {
	Tfid        types.String                          `tfsdk:"tfid"`
	Location    TemplateVariableLocation              `tfsdk:"location"`
	Name        types.String                          `tfsdk:"name"`
	Description types.String                          `tfsdk:"description"`
	Type        *TemplateVariableDataSourceTypeObject `tfsdk:"type"`
}
type TemplateVariableDataSourceTypeObject struct {
	Fqdn           types.String `tfsdk:"fqdn"`
	DevicePriority types.String `tfsdk:"device_priority"`
	AsNumber       types.String `tfsdk:"as_number"`
	IpRange        types.String `tfsdk:"ip_range"`
	GroupId        types.String `tfsdk:"group_id"`
	DeviceId       types.String `tfsdk:"device_id"`
	Interface      types.String `tfsdk:"interface"`
	QosProfile     types.String `tfsdk:"qos_profile"`
	EgressMax      types.String `tfsdk:"egress_max"`
	LinkTag        types.String `tfsdk:"link_tag"`
	IpNetmask      types.String `tfsdk:"ip_netmask"`
}

func (o *TemplateVariableDataSourceModel) CopyToPango(ctx context.Context, obj **template_variable.Entry, encrypted *map[string]types.String) diag.Diagnostics {
	var diags diag.Diagnostics
	var type_entry *template_variable.Type
	if o.Type != nil {
		if *obj != nil && (*obj).Type != nil {
			type_entry = (*obj).Type
		} else {
			type_entry = new(template_variable.Type)
		}

		diags.Append(o.Type.CopyToPango(ctx, &type_entry, encrypted)...)
		if diags.HasError() {
			return diags
		}
	}
	description_value := o.Description.ValueStringPointer()

	if (*obj) == nil {
		*obj = new(template_variable.Entry)
	}
	(*obj).Name = o.Name.ValueString()
	(*obj).Type = type_entry
	(*obj).Description = description_value

	return diags
}
func (o *TemplateVariableDataSourceTypeObject) CopyToPango(ctx context.Context, obj **template_variable.Type, encrypted *map[string]types.String) diag.Diagnostics {
	var diags diag.Diagnostics
	deviceId_value := o.DeviceId.ValueStringPointer()
	interface_value := o.Interface.ValueStringPointer()
	qosProfile_value := o.QosProfile.ValueStringPointer()
	egressMax_value := o.EgressMax.ValueStringPointer()
	linkTag_value := o.LinkTag.ValueStringPointer()
	ipNetmask_value := o.IpNetmask.ValueStringPointer()
	groupId_value := o.GroupId.ValueStringPointer()
	devicePriority_value := o.DevicePriority.ValueStringPointer()
	asNumber_value := o.AsNumber.ValueStringPointer()
	ipRange_value := o.IpRange.ValueStringPointer()
	fqdn_value := o.Fqdn.ValueStringPointer()

	if (*obj) == nil {
		*obj = new(template_variable.Type)
	}
	(*obj).DeviceId = deviceId_value
	(*obj).Interface = interface_value
	(*obj).QosProfile = qosProfile_value
	(*obj).EgressMax = egressMax_value
	(*obj).LinkTag = linkTag_value
	(*obj).IpNetmask = ipNetmask_value
	(*obj).GroupId = groupId_value
	(*obj).DevicePriority = devicePriority_value
	(*obj).AsNumber = asNumber_value
	(*obj).IpRange = ipRange_value
	(*obj).Fqdn = fqdn_value

	return diags
}

func (o *TemplateVariableDataSourceModel) CopyFromPango(ctx context.Context, obj *template_variable.Entry, encrypted *map[string]types.String) diag.Diagnostics {
	var diags diag.Diagnostics
	var type_object *TemplateVariableDataSourceTypeObject
	if obj.Type != nil {
		type_object = new(TemplateVariableDataSourceTypeObject)

		diags.Append(type_object.CopyFromPango(ctx, obj.Type, encrypted)...)
		if diags.HasError() {
			return diags
		}
	}
	var description_value types.String
	if obj.Description != nil {
		description_value = types.StringValue(*obj.Description)
	}
	o.Name = types.StringValue(obj.Name)
	o.Type = type_object
	o.Description = description_value

	return diags
}

func (o *TemplateVariableDataSourceTypeObject) CopyFromPango(ctx context.Context, obj *template_variable.Type, encrypted *map[string]types.String) diag.Diagnostics {
	var diags diag.Diagnostics
	var asNumber_value types.String
	if obj.AsNumber != nil {
		asNumber_value = types.StringValue(*obj.AsNumber)
	}
	var ipRange_value types.String
	if obj.IpRange != nil {
		ipRange_value = types.StringValue(*obj.IpRange)
	}
	var fqdn_value types.String
	if obj.Fqdn != nil {
		fqdn_value = types.StringValue(*obj.Fqdn)
	}
	var devicePriority_value types.String
	if obj.DevicePriority != nil {
		devicePriority_value = types.StringValue(*obj.DevicePriority)
	}
	var interface_value types.String
	if obj.Interface != nil {
		interface_value = types.StringValue(*obj.Interface)
	}
	var qosProfile_value types.String
	if obj.QosProfile != nil {
		qosProfile_value = types.StringValue(*obj.QosProfile)
	}
	var egressMax_value types.String
	if obj.EgressMax != nil {
		egressMax_value = types.StringValue(*obj.EgressMax)
	}
	var linkTag_value types.String
	if obj.LinkTag != nil {
		linkTag_value = types.StringValue(*obj.LinkTag)
	}
	var ipNetmask_value types.String
	if obj.IpNetmask != nil {
		ipNetmask_value = types.StringValue(*obj.IpNetmask)
	}
	var groupId_value types.String
	if obj.GroupId != nil {
		groupId_value = types.StringValue(*obj.GroupId)
	}
	var deviceId_value types.String
	if obj.DeviceId != nil {
		deviceId_value = types.StringValue(*obj.DeviceId)
	}
	o.AsNumber = asNumber_value
	o.IpRange = ipRange_value
	o.Fqdn = fqdn_value
	o.DevicePriority = devicePriority_value
	o.Interface = interface_value
	o.QosProfile = qosProfile_value
	o.EgressMax = egressMax_value
	o.LinkTag = linkTag_value
	o.IpNetmask = ipNetmask_value
	o.GroupId = groupId_value
	o.DeviceId = deviceId_value

	return diags
}

func TemplateVariableDataSourceSchema() dsschema.Schema {
	return dsschema.Schema{
		Attributes: map[string]dsschema.Attribute{

			"location": TemplateVariableDataSourceLocationSchema(),

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

			"description": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
				Required:    false,
				Optional:    true,
				Sensitive:   false,
			},

			"type": TemplateVariableDataSourceTypeSchema(),
		},
	}
}

func (o *TemplateVariableDataSourceModel) getTypeFor(name string) attr.Type {
	schema := TemplateVariableDataSourceSchema()
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

func TemplateVariableDataSourceTypeSchema() dsschema.SingleNestedAttribute {
	return dsschema.SingleNestedAttribute{
		Description: "",
		Required:    false,
		Computed:    true,
		Optional:    true,
		Sensitive:   false,
		Attributes: map[string]dsschema.Attribute{

			"as_number": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
				Required:    false,
				Optional:    true,
				Sensitive:   false,

				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.Expressions{
						path.MatchRelative().AtParent().AtName("link_tag"),
						path.MatchRelative().AtParent().AtName("ip_netmask"),
						path.MatchRelative().AtParent().AtName("group_id"),
						path.MatchRelative().AtParent().AtName("device_id"),
						path.MatchRelative().AtParent().AtName("interface"),
						path.MatchRelative().AtParent().AtName("qos_profile"),
						path.MatchRelative().AtParent().AtName("egress_max"),
						path.MatchRelative().AtParent().AtName("ip_range"),
						path.MatchRelative().AtParent().AtName("fqdn"),
						path.MatchRelative().AtParent().AtName("device_priority"),
						path.MatchRelative().AtParent().AtName("as_number"),
					}...),
				},
			},

			"ip_range": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
				Required:    false,
				Optional:    true,
				Sensitive:   false,
			},

			"fqdn": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
				Required:    false,
				Optional:    true,
				Sensitive:   false,
			},

			"device_priority": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
				Required:    false,
				Optional:    true,
				Sensitive:   false,
			},

			"interface": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
				Required:    false,
				Optional:    true,
				Sensitive:   false,
			},

			"qos_profile": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
				Required:    false,
				Optional:    true,
				Sensitive:   false,
			},

			"egress_max": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
				Required:    false,
				Optional:    true,
				Sensitive:   false,
			},

			"link_tag": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
				Required:    false,
				Optional:    true,
				Sensitive:   false,
			},

			"ip_netmask": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
				Required:    false,
				Optional:    true,
				Sensitive:   false,
			},

			"group_id": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
				Required:    false,
				Optional:    true,
				Sensitive:   false,
			},

			"device_id": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
				Required:    false,
				Optional:    true,
				Sensitive:   false,
			},
		},
	}
}

func (o *TemplateVariableDataSourceTypeObject) getTypeFor(name string) attr.Type {
	schema := TemplateVariableDataSourceTypeSchema()
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

func TemplateVariableDataSourceLocationSchema() rsschema.Attribute {
	return TemplateVariableLocationSchema()
}

// Metadata returns the data source type name.
func (d *TemplateVariableDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_template_variable"
}

// Schema defines the schema for this data source.
func (d *TemplateVariableDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = TemplateVariableDataSourceSchema()
}

// Configure prepares the struct.
func (d *TemplateVariableDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*pango.Client)
	specifier, _, err := template_variable.Versioning(d.client.Versioning())
	if err != nil {
		resp.Diagnostics.AddError("Failed to configure SDK client", err.Error())
		return
	}
	d.manager = sdkmanager.NewEntryObjectManager(d.client, template_variable.NewService(d.client), specifier, template_variable.SpecMatches)
}

func (o *TemplateVariableDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var savestate, state TemplateVariableDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &savestate)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var loc TemplateVariableDataSourceTfid
	loc.Name = *savestate.Name.ValueStringPointer()

	if savestate.Location.Template != nil {
		loc.Location.Template = &template_variable.TemplateLocation{

			PanoramaDevice: savestate.Location.Template.PanoramaDevice.ValueString(),
			Template:       savestate.Location.Template.Name.ValueString(),
		}
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource read", map[string]any{
		"resource_name": "panos_template_variable_resource",
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
	_ resource.Resource                = &TemplateVariableResource{}
	_ resource.ResourceWithConfigure   = &TemplateVariableResource{}
	_ resource.ResourceWithImportState = &TemplateVariableResource{}
)

func NewTemplateVariableResource() resource.Resource {
	return &TemplateVariableResource{}
}

type TemplateVariableResource struct {
	client  *pango.Client
	manager *sdkmanager.EntryObjectManager[*template_variable.Entry, template_variable.Location, *template_variable.Service]
}
type TemplateVariableResourceTfid struct {
	Name     string                     `json:"name"`
	Location template_variable.Location `json:"location"`
}

func (o *TemplateVariableResourceTfid) IsValid() error {
	if o.Name == "" {
		return fmt.Errorf("name is unspecified")
	}
	return o.Location.IsValid()
}

func TemplateVariableResourceLocationSchema() rsschema.Attribute {
	return TemplateVariableLocationSchema()
}

type TemplateVariableResourceModel struct {
	Tfid        types.String                        `tfsdk:"tfid"`
	Location    TemplateVariableLocation            `tfsdk:"location"`
	Name        types.String                        `tfsdk:"name"`
	Description types.String                        `tfsdk:"description"`
	Type        *TemplateVariableResourceTypeObject `tfsdk:"type"`
}
type TemplateVariableResourceTypeObject struct {
	DevicePriority types.String `tfsdk:"device_priority"`
	AsNumber       types.String `tfsdk:"as_number"`
	IpRange        types.String `tfsdk:"ip_range"`
	Fqdn           types.String `tfsdk:"fqdn"`
	DeviceId       types.String `tfsdk:"device_id"`
	Interface      types.String `tfsdk:"interface"`
	QosProfile     types.String `tfsdk:"qos_profile"`
	EgressMax      types.String `tfsdk:"egress_max"`
	LinkTag        types.String `tfsdk:"link_tag"`
	IpNetmask      types.String `tfsdk:"ip_netmask"`
	GroupId        types.String `tfsdk:"group_id"`
}

func (r *TemplateVariableResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_template_variable"
}

func (r *TemplateVariableResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
}

// <ResourceSchema>

func TemplateVariableResourceSchema() rsschema.Schema {
	return rsschema.Schema{
		Attributes: map[string]rsschema.Attribute{

			"location": TemplateVariableResourceLocationSchema(),

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
				Description: "",
				Computed:    false,
				Required:    false,
				Optional:    true,
				Sensitive:   false,
			},

			"type": TemplateVariableResourceTypeSchema(),
		},
	}
}

func (o *TemplateVariableResourceModel) getTypeFor(name string) attr.Type {
	schema := TemplateVariableResourceSchema()
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

func TemplateVariableResourceTypeSchema() rsschema.SingleNestedAttribute {
	return rsschema.SingleNestedAttribute{
		Description: "",
		Required:    false,
		Computed:    false,
		Optional:    true,
		Sensitive:   false,
		Attributes: map[string]rsschema.Attribute{

			"link_tag": rsschema.StringAttribute{
				Description: "",
				Computed:    false,
				Required:    false,
				Optional:    true,
				Sensitive:   false,

				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.Expressions{
						path.MatchRelative().AtParent().AtName("group_id"),
						path.MatchRelative().AtParent().AtName("device_id"),
						path.MatchRelative().AtParent().AtName("interface"),
						path.MatchRelative().AtParent().AtName("qos_profile"),
						path.MatchRelative().AtParent().AtName("egress_max"),
						path.MatchRelative().AtParent().AtName("link_tag"),
						path.MatchRelative().AtParent().AtName("ip_netmask"),
						path.MatchRelative().AtParent().AtName("fqdn"),
						path.MatchRelative().AtParent().AtName("device_priority"),
						path.MatchRelative().AtParent().AtName("as_number"),
						path.MatchRelative().AtParent().AtName("ip_range"),
					}...),
				},
			},

			"ip_netmask": rsschema.StringAttribute{
				Description: "",
				Computed:    false,
				Required:    false,
				Optional:    true,
				Sensitive:   false,
			},

			"group_id": rsschema.StringAttribute{
				Description: "",
				Computed:    false,
				Required:    false,
				Optional:    true,
				Sensitive:   false,
			},

			"device_id": rsschema.StringAttribute{
				Description: "",
				Computed:    false,
				Required:    false,
				Optional:    true,
				Sensitive:   false,
			},

			"interface": rsschema.StringAttribute{
				Description: "",
				Computed:    false,
				Required:    false,
				Optional:    true,
				Sensitive:   false,
			},

			"qos_profile": rsschema.StringAttribute{
				Description: "",
				Computed:    false,
				Required:    false,
				Optional:    true,
				Sensitive:   false,
			},

			"egress_max": rsschema.StringAttribute{
				Description: "",
				Computed:    false,
				Required:    false,
				Optional:    true,
				Sensitive:   false,
			},

			"ip_range": rsschema.StringAttribute{
				Description: "",
				Computed:    false,
				Required:    false,
				Optional:    true,
				Sensitive:   false,
			},

			"fqdn": rsschema.StringAttribute{
				Description: "",
				Computed:    false,
				Required:    false,
				Optional:    true,
				Sensitive:   false,
			},

			"device_priority": rsschema.StringAttribute{
				Description: "",
				Computed:    false,
				Required:    false,
				Optional:    true,
				Sensitive:   false,
			},

			"as_number": rsschema.StringAttribute{
				Description: "",
				Computed:    false,
				Required:    false,
				Optional:    true,
				Sensitive:   false,
			},
		},
	}
}

func (o *TemplateVariableResourceTypeObject) getTypeFor(name string) attr.Type {
	schema := TemplateVariableResourceTypeSchema()
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

func (r *TemplateVariableResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = TemplateVariableResourceSchema()
}

// </ResourceSchema>

func (r *TemplateVariableResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*pango.Client)
	specifier, _, err := template_variable.Versioning(r.client.Versioning())
	if err != nil {
		resp.Diagnostics.AddError("Failed to configure SDK client", err.Error())
		return
	}
	r.manager = sdkmanager.NewEntryObjectManager(r.client, template_variable.NewService(r.client), specifier, template_variable.SpecMatches)
}

func (o *TemplateVariableResourceModel) CopyToPango(ctx context.Context, obj **template_variable.Entry, encrypted *map[string]types.String) diag.Diagnostics {
	var diags diag.Diagnostics
	description_value := o.Description.ValueStringPointer()
	var type_entry *template_variable.Type
	if o.Type != nil {
		if *obj != nil && (*obj).Type != nil {
			type_entry = (*obj).Type
		} else {
			type_entry = new(template_variable.Type)
		}

		diags.Append(o.Type.CopyToPango(ctx, &type_entry, encrypted)...)
		if diags.HasError() {
			return diags
		}
	}

	if (*obj) == nil {
		*obj = new(template_variable.Entry)
	}
	(*obj).Name = o.Name.ValueString()
	(*obj).Description = description_value
	(*obj).Type = type_entry

	return diags
}
func (o *TemplateVariableResourceTypeObject) CopyToPango(ctx context.Context, obj **template_variable.Type, encrypted *map[string]types.String) diag.Diagnostics {
	var diags diag.Diagnostics
	asNumber_value := o.AsNumber.ValueStringPointer()
	ipRange_value := o.IpRange.ValueStringPointer()
	fqdn_value := o.Fqdn.ValueStringPointer()
	devicePriority_value := o.DevicePriority.ValueStringPointer()
	interface_value := o.Interface.ValueStringPointer()
	qosProfile_value := o.QosProfile.ValueStringPointer()
	egressMax_value := o.EgressMax.ValueStringPointer()
	linkTag_value := o.LinkTag.ValueStringPointer()
	ipNetmask_value := o.IpNetmask.ValueStringPointer()
	groupId_value := o.GroupId.ValueStringPointer()
	deviceId_value := o.DeviceId.ValueStringPointer()

	if (*obj) == nil {
		*obj = new(template_variable.Type)
	}
	(*obj).AsNumber = asNumber_value
	(*obj).IpRange = ipRange_value
	(*obj).Fqdn = fqdn_value
	(*obj).DevicePriority = devicePriority_value
	(*obj).Interface = interface_value
	(*obj).QosProfile = qosProfile_value
	(*obj).EgressMax = egressMax_value
	(*obj).LinkTag = linkTag_value
	(*obj).IpNetmask = ipNetmask_value
	(*obj).GroupId = groupId_value
	(*obj).DeviceId = deviceId_value

	return diags
}

func (o *TemplateVariableResourceModel) CopyFromPango(ctx context.Context, obj *template_variable.Entry, encrypted *map[string]types.String) diag.Diagnostics {
	var diags diag.Diagnostics
	var type_object *TemplateVariableResourceTypeObject
	if obj.Type != nil {
		type_object = new(TemplateVariableResourceTypeObject)

		diags.Append(type_object.CopyFromPango(ctx, obj.Type, encrypted)...)
		if diags.HasError() {
			return diags
		}
	}
	var description_value types.String
	if obj.Description != nil {
		description_value = types.StringValue(*obj.Description)
	}
	o.Name = types.StringValue(obj.Name)
	o.Description = description_value
	o.Type = type_object

	return diags
}

func (o *TemplateVariableResourceTypeObject) CopyFromPango(ctx context.Context, obj *template_variable.Type, encrypted *map[string]types.String) diag.Diagnostics {
	var diags diag.Diagnostics
	var ipRange_value types.String
	if obj.IpRange != nil {
		ipRange_value = types.StringValue(*obj.IpRange)
	}
	var fqdn_value types.String
	if obj.Fqdn != nil {
		fqdn_value = types.StringValue(*obj.Fqdn)
	}
	var devicePriority_value types.String
	if obj.DevicePriority != nil {
		devicePriority_value = types.StringValue(*obj.DevicePriority)
	}
	var asNumber_value types.String
	if obj.AsNumber != nil {
		asNumber_value = types.StringValue(*obj.AsNumber)
	}
	var ipNetmask_value types.String
	if obj.IpNetmask != nil {
		ipNetmask_value = types.StringValue(*obj.IpNetmask)
	}
	var groupId_value types.String
	if obj.GroupId != nil {
		groupId_value = types.StringValue(*obj.GroupId)
	}
	var deviceId_value types.String
	if obj.DeviceId != nil {
		deviceId_value = types.StringValue(*obj.DeviceId)
	}
	var interface_value types.String
	if obj.Interface != nil {
		interface_value = types.StringValue(*obj.Interface)
	}
	var qosProfile_value types.String
	if obj.QosProfile != nil {
		qosProfile_value = types.StringValue(*obj.QosProfile)
	}
	var egressMax_value types.String
	if obj.EgressMax != nil {
		egressMax_value = types.StringValue(*obj.EgressMax)
	}
	var linkTag_value types.String
	if obj.LinkTag != nil {
		linkTag_value = types.StringValue(*obj.LinkTag)
	}
	o.IpRange = ipRange_value
	o.Fqdn = fqdn_value
	o.DevicePriority = devicePriority_value
	o.AsNumber = asNumber_value
	o.IpNetmask = ipNetmask_value
	o.GroupId = groupId_value
	o.DeviceId = deviceId_value
	o.Interface = interface_value
	o.QosProfile = qosProfile_value
	o.EgressMax = egressMax_value
	o.LinkTag = linkTag_value

	return diags
}

func (r *TemplateVariableResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state TemplateVariableResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource create", map[string]any{
		"resource_name": "panos_template_variable_resource",
		"function":      "Create",
		"name":          state.Name.ValueString(),
	})

	// Verify mode.
	if r.client.Hostname == "" {
		resp.Diagnostics.AddError("Invalid mode error", InspectionModeError)
		return
	}

	// Determine the location.
	loc := TemplateVariableResourceTfid{Name: state.Name.ValueString()}

	// TODO: this needs to handle location structure for UUID style shared has nested structure type

	if state.Location.Template != nil {
		loc.Location.Template = &template_variable.TemplateLocation{

			PanoramaDevice: state.Location.Template.PanoramaDevice.ValueString(),
			Template:       state.Location.Template.Name.ValueString(),
		}
	}

	if err := loc.IsValid(); err != nil {
		resp.Diagnostics.AddError("Invalid location", err.Error())
		return
	}

	// Load the desired config.
	var obj *template_variable.Entry

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

func (o *TemplateVariableResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var savestate, state TemplateVariableResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &savestate)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var loc TemplateVariableResourceTfid
	// Parse the location from tfid.
	if err := DecodeLocation(savestate.Tfid.ValueString(), &loc); err != nil {
		resp.Diagnostics.AddError("Error parsing tfid", err.Error())
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource read", map[string]any{
		"resource_name": "panos_template_variable_resource",
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

func (r *TemplateVariableResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var plan, state TemplateVariableResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var loc TemplateVariableResourceTfid
	if err := DecodeLocation(state.Tfid.ValueString(), &loc); err != nil {
		resp.Diagnostics.AddError("Error parsing tfid", err.Error())
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource update", map[string]any{
		"resource_name": "panos_template_variable_resource",
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

func (r *TemplateVariableResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var state TemplateVariableResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the location from tfid.
	var loc TemplateVariableResourceTfid
	if err := DecodeLocation(state.Tfid.ValueString(), &loc); err != nil {
		resp.Diagnostics.AddError("error parsing tfid", err.Error())
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource delete", map[string]any{
		"resource_name": "panos_template_variable_resource",
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

func (r *TemplateVariableResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("tfid"), req, resp)
}

type TemplateVariableTemplateLocation struct {
	PanoramaDevice types.String `tfsdk:"panorama_device"`
	Name           types.String `tfsdk:"name"`
}
type TemplateVariableLocation struct {
	Template *TemplateVariableTemplateLocation `tfsdk:"template"`
}

func TemplateVariableLocationSchema() rsschema.Attribute {
	return rsschema.SingleNestedAttribute{
		Description: "The location of this object.",
		Required:    true,
		Attributes: map[string]rsschema.Attribute{
			"template": rsschema.SingleNestedAttribute{
				Description: "Located in a specific template.",
				Optional:    true,
				Attributes: map[string]rsschema.Attribute{
					"panorama_device": rsschema.StringAttribute{
						Description: "The panorama device.",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString("localhost.localdomain"),
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"name": rsschema.StringAttribute{
						Description: "The template.",
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