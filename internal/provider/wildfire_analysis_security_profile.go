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
	"github.com/PaloAltoNetworks/pango/objects/profiles/wildfireanalysis"

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
	_ datasource.DataSource              = &WildfireAnalysisSecurityProfileDataSource{}
	_ datasource.DataSourceWithConfigure = &WildfireAnalysisSecurityProfileDataSource{}
)

func NewWildfireAnalysisSecurityProfileDataSource() datasource.DataSource {
	return &WildfireAnalysisSecurityProfileDataSource{}
}

type WildfireAnalysisSecurityProfileDataSource struct {
	client  *pango.Client
	manager *sdkmanager.EntryObjectManager[*wildfireanalysis.Entry, wildfireanalysis.Location, *wildfireanalysis.Service]
}

type WildfireAnalysisSecurityProfileDataSourceFilter struct {
	// TODO: Generate Data Source filter via function
}

type WildfireAnalysisSecurityProfileDataSourceModel struct {
	Location        WildfireAnalysisSecurityProfileLocation `tfsdk:"location"`
	Name            types.String                            `tfsdk:"name"`
	Description     types.String                            `tfsdk:"description"`
	DisableOverride types.String                            `tfsdk:"disable_override"`
	Rules           types.List                              `tfsdk:"rules"`
}
type WildfireAnalysisSecurityProfileDataSourceRulesObject struct {
	Name        types.String `tfsdk:"name"`
	Application types.List   `tfsdk:"application"`
	FileType    types.List   `tfsdk:"file_type"`
	Direction   types.String `tfsdk:"direction"`
	Analysis    types.String `tfsdk:"analysis"`
}

func (o *WildfireAnalysisSecurityProfileDataSourceModel) CopyToPango(ctx context.Context, obj **wildfireanalysis.Entry, encrypted *map[string]types.String) diag.Diagnostics {
	var diags diag.Diagnostics
	description_value := o.Description.ValueStringPointer()
	disableOverride_value := o.DisableOverride.ValueStringPointer()
	var rules_tf_entries []WildfireAnalysisSecurityProfileDataSourceRulesObject
	var rules_pango_entries []wildfireanalysis.Rules
	{
		d := o.Rules.ElementsAs(ctx, &rules_tf_entries, false)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}
		for _, elt := range rules_tf_entries {
			var entry *wildfireanalysis.Rules
			diags.Append(elt.CopyToPango(ctx, &entry, encrypted)...)
			if diags.HasError() {
				return diags
			}
			rules_pango_entries = append(rules_pango_entries, *entry)
		}
	}

	if (*obj) == nil {
		*obj = new(wildfireanalysis.Entry)
	}
	(*obj).Name = o.Name.ValueString()
	(*obj).Description = description_value
	(*obj).DisableOverride = disableOverride_value
	(*obj).Rules = rules_pango_entries

	return diags
}
func (o *WildfireAnalysisSecurityProfileDataSourceRulesObject) CopyToPango(ctx context.Context, obj **wildfireanalysis.Rules, encrypted *map[string]types.String) diag.Diagnostics {
	var diags diag.Diagnostics
	application_pango_entries := make([]string, 0)
	diags.Append(o.Application.ElementsAs(ctx, &application_pango_entries, false)...)
	if diags.HasError() {
		return diags
	}
	fileType_pango_entries := make([]string, 0)
	diags.Append(o.FileType.ElementsAs(ctx, &fileType_pango_entries, false)...)
	if diags.HasError() {
		return diags
	}
	direction_value := o.Direction.ValueStringPointer()
	analysis_value := o.Analysis.ValueStringPointer()

	if (*obj) == nil {
		*obj = new(wildfireanalysis.Rules)
	}
	(*obj).Name = o.Name.ValueString()
	(*obj).Application = application_pango_entries
	(*obj).FileType = fileType_pango_entries
	(*obj).Direction = direction_value
	(*obj).Analysis = analysis_value

	return diags
}

func (o *WildfireAnalysisSecurityProfileDataSourceModel) CopyFromPango(ctx context.Context, obj *wildfireanalysis.Entry, encrypted *map[string]types.String) diag.Diagnostics {
	var diags diag.Diagnostics
	var rules_list types.List
	{
		var rules_tf_entries []WildfireAnalysisSecurityProfileDataSourceRulesObject
		for _, elt := range obj.Rules {
			var entry WildfireAnalysisSecurityProfileDataSourceRulesObject
			entry_diags := entry.CopyFromPango(ctx, &elt, encrypted)
			diags.Append(entry_diags...)
			rules_tf_entries = append(rules_tf_entries, entry)
		}
		var list_diags diag.Diagnostics
		schemaType := o.getTypeFor("rules")
		rules_list, list_diags = types.ListValueFrom(ctx, schemaType, rules_tf_entries)
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
	o.Name = types.StringValue(obj.Name)
	o.Description = description_value
	o.DisableOverride = disableOverride_value
	o.Rules = rules_list

	return diags
}

func (o *WildfireAnalysisSecurityProfileDataSourceRulesObject) CopyFromPango(ctx context.Context, obj *wildfireanalysis.Rules, encrypted *map[string]types.String) diag.Diagnostics {
	var diags diag.Diagnostics
	var application_list types.List
	{
		var list_diags diag.Diagnostics
		application_list, list_diags = types.ListValueFrom(ctx, types.StringType, obj.Application)
		diags.Append(list_diags...)
	}
	var fileType_list types.List
	{
		var list_diags diag.Diagnostics
		fileType_list, list_diags = types.ListValueFrom(ctx, types.StringType, obj.FileType)
		diags.Append(list_diags...)
	}

	var direction_value types.String
	if obj.Direction != nil {
		direction_value = types.StringValue(*obj.Direction)
	}
	var analysis_value types.String
	if obj.Analysis != nil {
		analysis_value = types.StringValue(*obj.Analysis)
	}
	o.Name = types.StringValue(obj.Name)
	o.Application = application_list
	o.FileType = fileType_list
	o.Direction = direction_value
	o.Analysis = analysis_value

	return diags
}

func WildfireAnalysisSecurityProfileDataSourceSchema() dsschema.Schema {
	return dsschema.Schema{
		Attributes: map[string]dsschema.Attribute{

			"location": WildfireAnalysisSecurityProfileDataSourceLocationSchema(),

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

			"rules": dsschema.ListNestedAttribute{
				Description:  "",
				Required:     false,
				Optional:     true,
				Computed:     true,
				Sensitive:    false,
				NestedObject: WildfireAnalysisSecurityProfileDataSourceRulesSchema(),
			},
		},
	}
}

func (o *WildfireAnalysisSecurityProfileDataSourceModel) getTypeFor(name string) attr.Type {
	schema := WildfireAnalysisSecurityProfileDataSourceSchema()
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

func WildfireAnalysisSecurityProfileDataSourceRulesSchema() dsschema.NestedAttributeObject {
	return dsschema.NestedAttributeObject{
		Attributes: map[string]dsschema.Attribute{

			"name": dsschema.StringAttribute{
				Description: "",
				Computed:    false,
				Required:    true,
				Optional:    false,
				Sensitive:   false,
			},

			"application": dsschema.ListAttribute{
				Description: "",
				Required:    false,
				Optional:    true,
				Computed:    true,
				Sensitive:   false,
				ElementType: types.StringType,
			},

			"file_type": dsschema.ListAttribute{
				Description: "",
				Required:    false,
				Optional:    true,
				Computed:    true,
				Sensitive:   false,
				ElementType: types.StringType,
			},

			"direction": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
				Required:    false,
				Optional:    true,
				Sensitive:   false,
			},

			"analysis": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
				Required:    false,
				Optional:    true,
				Sensitive:   false,
			},
		},
	}
}

func (o *WildfireAnalysisSecurityProfileDataSourceRulesObject) getTypeFor(name string) attr.Type {
	schema := WildfireAnalysisSecurityProfileDataSourceRulesSchema()
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

func WildfireAnalysisSecurityProfileDataSourceLocationSchema() rsschema.Attribute {
	return WildfireAnalysisSecurityProfileLocationSchema()
}

// Metadata returns the data source type name.
func (d *WildfireAnalysisSecurityProfileDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_wildfire_analysis_security_profile"
}

// Schema defines the schema for this data source.
func (d *WildfireAnalysisSecurityProfileDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = WildfireAnalysisSecurityProfileDataSourceSchema()
}

// Configure prepares the struct.
func (d *WildfireAnalysisSecurityProfileDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*pango.Client)
	specifier, _, err := wildfireanalysis.Versioning(d.client.Versioning())
	if err != nil {
		resp.Diagnostics.AddError("Failed to configure SDK client", err.Error())
		return
	}
	d.manager = sdkmanager.NewEntryObjectManager(d.client, wildfireanalysis.NewService(d.client), specifier, wildfireanalysis.SpecMatches)
}
func (o *WildfireAnalysisSecurityProfileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var savestate, state WildfireAnalysisSecurityProfileDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &savestate)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var location wildfireanalysis.Location

	if !savestate.Location.Shared.IsNull() && savestate.Location.Shared.ValueBool() {
		location.Shared = true
	}
	if savestate.Location.DeviceGroup != nil {
		location.DeviceGroup = &wildfireanalysis.DeviceGroupLocation{

			PanoramaDevice: savestate.Location.DeviceGroup.PanoramaDevice.ValueString(),
			DeviceGroup:    savestate.Location.DeviceGroup.Name.ValueString(),
		}
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource read", map[string]any{
		"resource_name": "panos_wildfire_analysis_security_profile_resource",
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
	_ resource.Resource                = &WildfireAnalysisSecurityProfileResource{}
	_ resource.ResourceWithConfigure   = &WildfireAnalysisSecurityProfileResource{}
	_ resource.ResourceWithImportState = &WildfireAnalysisSecurityProfileResource{}
)

func NewWildfireAnalysisSecurityProfileResource() resource.Resource {
	if _, found := resourceFuncMap["panos_wildfire_analysis_security_profile"]; !found {
		resourceFuncMap["panos_wildfire_analysis_security_profile"] = resourceFuncs{
			CreateImportId: WildfireAnalysisSecurityProfileImportStateCreator,
		}
	}
	return &WildfireAnalysisSecurityProfileResource{}
}

type WildfireAnalysisSecurityProfileResource struct {
	client  *pango.Client
	manager *sdkmanager.EntryObjectManager[*wildfireanalysis.Entry, wildfireanalysis.Location, *wildfireanalysis.Service]
}

func WildfireAnalysisSecurityProfileResourceLocationSchema() rsschema.Attribute {
	return WildfireAnalysisSecurityProfileLocationSchema()
}

type WildfireAnalysisSecurityProfileResourceModel struct {
	Location        WildfireAnalysisSecurityProfileLocation `tfsdk:"location"`
	Name            types.String                            `tfsdk:"name"`
	Description     types.String                            `tfsdk:"description"`
	DisableOverride types.String                            `tfsdk:"disable_override"`
	Rules           types.List                              `tfsdk:"rules"`
}
type WildfireAnalysisSecurityProfileResourceRulesObject struct {
	Name        types.String `tfsdk:"name"`
	Application types.List   `tfsdk:"application"`
	FileType    types.List   `tfsdk:"file_type"`
	Direction   types.String `tfsdk:"direction"`
	Analysis    types.String `tfsdk:"analysis"`
}

func (r *WildfireAnalysisSecurityProfileResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
}

// <ResourceSchema>

func WildfireAnalysisSecurityProfileResourceSchema() rsschema.Schema {
	return rsschema.Schema{
		Attributes: map[string]rsschema.Attribute{

			"location": WildfireAnalysisSecurityProfileResourceLocationSchema(),

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
						"yes",
						"no",
					}...),
				},
			},

			"rules": rsschema.ListNestedAttribute{
				Description:  "",
				Required:     false,
				Optional:     true,
				Computed:     false,
				Sensitive:    false,
				NestedObject: WildfireAnalysisSecurityProfileResourceRulesSchema(),
			},
		},
	}
}

func (o *WildfireAnalysisSecurityProfileResourceModel) getTypeFor(name string) attr.Type {
	schema := WildfireAnalysisSecurityProfileResourceSchema()
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

func WildfireAnalysisSecurityProfileResourceRulesSchema() rsschema.NestedAttributeObject {
	return rsschema.NestedAttributeObject{
		Attributes: map[string]rsschema.Attribute{

			"name": rsschema.StringAttribute{
				Description: "",
				Computed:    false,
				Required:    true,
				Optional:    false,
				Sensitive:   false,
			},

			"application": rsschema.ListAttribute{
				Description: "",
				Required:    false,
				Optional:    true,
				Computed:    false,
				Sensitive:   false,
				ElementType: types.StringType,
			},

			"file_type": rsschema.ListAttribute{
				Description: "",
				Required:    false,
				Optional:    true,
				Computed:    false,
				Sensitive:   false,
				ElementType: types.StringType,
			},

			"direction": rsschema.StringAttribute{
				Description: "",
				Computed:    false,
				Required:    false,
				Optional:    true,
				Sensitive:   false,
			},

			"analysis": rsschema.StringAttribute{
				Description: "",
				Computed:    true,
				Required:    false,
				Optional:    true,
				Sensitive:   false,
				Default:     stringdefault.StaticString("public-cloud"),
			},
		},
	}
}

func (o *WildfireAnalysisSecurityProfileResourceRulesObject) getTypeFor(name string) attr.Type {
	schema := WildfireAnalysisSecurityProfileResourceRulesSchema()
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

func (r *WildfireAnalysisSecurityProfileResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_wildfire_analysis_security_profile"
}

func (r *WildfireAnalysisSecurityProfileResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = WildfireAnalysisSecurityProfileResourceSchema()
}

// </ResourceSchema>

func (r *WildfireAnalysisSecurityProfileResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*pango.Client)
	specifier, _, err := wildfireanalysis.Versioning(r.client.Versioning())
	if err != nil {
		resp.Diagnostics.AddError("Failed to configure SDK client", err.Error())
		return
	}
	r.manager = sdkmanager.NewEntryObjectManager(r.client, wildfireanalysis.NewService(r.client), specifier, wildfireanalysis.SpecMatches)
}

func (o *WildfireAnalysisSecurityProfileResourceModel) CopyToPango(ctx context.Context, obj **wildfireanalysis.Entry, encrypted *map[string]types.String) diag.Diagnostics {
	var diags diag.Diagnostics
	description_value := o.Description.ValueStringPointer()
	disableOverride_value := o.DisableOverride.ValueStringPointer()
	var rules_tf_entries []WildfireAnalysisSecurityProfileResourceRulesObject
	var rules_pango_entries []wildfireanalysis.Rules
	{
		d := o.Rules.ElementsAs(ctx, &rules_tf_entries, false)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}
		for _, elt := range rules_tf_entries {
			var entry *wildfireanalysis.Rules
			diags.Append(elt.CopyToPango(ctx, &entry, encrypted)...)
			if diags.HasError() {
				return diags
			}
			rules_pango_entries = append(rules_pango_entries, *entry)
		}
	}

	if (*obj) == nil {
		*obj = new(wildfireanalysis.Entry)
	}
	(*obj).Name = o.Name.ValueString()
	(*obj).Description = description_value
	(*obj).DisableOverride = disableOverride_value
	(*obj).Rules = rules_pango_entries

	return diags
}
func (o *WildfireAnalysisSecurityProfileResourceRulesObject) CopyToPango(ctx context.Context, obj **wildfireanalysis.Rules, encrypted *map[string]types.String) diag.Diagnostics {
	var diags diag.Diagnostics
	application_pango_entries := make([]string, 0)
	diags.Append(o.Application.ElementsAs(ctx, &application_pango_entries, false)...)
	if diags.HasError() {
		return diags
	}
	fileType_pango_entries := make([]string, 0)
	diags.Append(o.FileType.ElementsAs(ctx, &fileType_pango_entries, false)...)
	if diags.HasError() {
		return diags
	}
	direction_value := o.Direction.ValueStringPointer()
	analysis_value := o.Analysis.ValueStringPointer()

	if (*obj) == nil {
		*obj = new(wildfireanalysis.Rules)
	}
	(*obj).Name = o.Name.ValueString()
	(*obj).Application = application_pango_entries
	(*obj).FileType = fileType_pango_entries
	(*obj).Direction = direction_value
	(*obj).Analysis = analysis_value

	return diags
}

func (o *WildfireAnalysisSecurityProfileResourceModel) CopyFromPango(ctx context.Context, obj *wildfireanalysis.Entry, encrypted *map[string]types.String) diag.Diagnostics {
	var diags diag.Diagnostics
	var rules_list types.List
	{
		var rules_tf_entries []WildfireAnalysisSecurityProfileResourceRulesObject
		for _, elt := range obj.Rules {
			var entry WildfireAnalysisSecurityProfileResourceRulesObject
			entry_diags := entry.CopyFromPango(ctx, &elt, encrypted)
			diags.Append(entry_diags...)
			rules_tf_entries = append(rules_tf_entries, entry)
		}
		var list_diags diag.Diagnostics
		schemaType := o.getTypeFor("rules")
		rules_list, list_diags = types.ListValueFrom(ctx, schemaType, rules_tf_entries)
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
	o.Name = types.StringValue(obj.Name)
	o.Description = description_value
	o.DisableOverride = disableOverride_value
	o.Rules = rules_list

	return diags
}

func (o *WildfireAnalysisSecurityProfileResourceRulesObject) CopyFromPango(ctx context.Context, obj *wildfireanalysis.Rules, encrypted *map[string]types.String) diag.Diagnostics {
	var diags diag.Diagnostics
	var application_list types.List
	{
		var list_diags diag.Diagnostics
		application_list, list_diags = types.ListValueFrom(ctx, types.StringType, obj.Application)
		diags.Append(list_diags...)
	}
	var fileType_list types.List
	{
		var list_diags diag.Diagnostics
		fileType_list, list_diags = types.ListValueFrom(ctx, types.StringType, obj.FileType)
		diags.Append(list_diags...)
	}

	var direction_value types.String
	if obj.Direction != nil {
		direction_value = types.StringValue(*obj.Direction)
	}
	var analysis_value types.String
	if obj.Analysis != nil {
		analysis_value = types.StringValue(*obj.Analysis)
	}
	o.Name = types.StringValue(obj.Name)
	o.Application = application_list
	o.FileType = fileType_list
	o.Direction = direction_value
	o.Analysis = analysis_value

	return diags
}

func (r *WildfireAnalysisSecurityProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state WildfireAnalysisSecurityProfileResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource create", map[string]any{
		"resource_name": "panos_wildfire_analysis_security_profile_resource",
		"function":      "Create",
		"name":          state.Name.ValueString(),
	})

	// Verify mode.
	if r.client.Hostname == "" {
		resp.Diagnostics.AddError("Invalid mode error", InspectionModeError)
		return
	}

	// Determine the location.

	var location wildfireanalysis.Location

	if !state.Location.Shared.IsNull() && state.Location.Shared.ValueBool() {
		location.Shared = true
	}
	if state.Location.DeviceGroup != nil {
		location.DeviceGroup = &wildfireanalysis.DeviceGroupLocation{

			PanoramaDevice: state.Location.DeviceGroup.PanoramaDevice.ValueString(),
			DeviceGroup:    state.Location.DeviceGroup.Name.ValueString(),
		}
	}

	if err := location.IsValid(); err != nil {
		resp.Diagnostics.AddError("Invalid location", err.Error())
		return
	}

	// Load the desired config.
	var obj *wildfireanalysis.Entry

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
func (o *WildfireAnalysisSecurityProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var savestate, state WildfireAnalysisSecurityProfileResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &savestate)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var location wildfireanalysis.Location

	if !savestate.Location.Shared.IsNull() && savestate.Location.Shared.ValueBool() {
		location.Shared = true
	}
	if savestate.Location.DeviceGroup != nil {
		location.DeviceGroup = &wildfireanalysis.DeviceGroupLocation{

			PanoramaDevice: savestate.Location.DeviceGroup.PanoramaDevice.ValueString(),
			DeviceGroup:    savestate.Location.DeviceGroup.Name.ValueString(),
		}
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource read", map[string]any{
		"resource_name": "panos_wildfire_analysis_security_profile_resource",
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
func (r *WildfireAnalysisSecurityProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var plan, state WildfireAnalysisSecurityProfileResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var location wildfireanalysis.Location

	if !state.Location.Shared.IsNull() && state.Location.Shared.ValueBool() {
		location.Shared = true
	}
	if state.Location.DeviceGroup != nil {
		location.DeviceGroup = &wildfireanalysis.DeviceGroupLocation{

			PanoramaDevice: state.Location.DeviceGroup.PanoramaDevice.ValueString(),
			DeviceGroup:    state.Location.DeviceGroup.Name.ValueString(),
		}
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource update", map[string]any{
		"resource_name": "panos_wildfire_analysis_security_profile_resource",
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
func (r *WildfireAnalysisSecurityProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var state WildfireAnalysisSecurityProfileResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource delete", map[string]any{
		"resource_name": "panos_wildfire_analysis_security_profile_resource",
		"function":      "Delete",
		"name":          state.Name.ValueString(),
	})

	// Verify mode.
	if r.client.Hostname == "" {
		resp.Diagnostics.AddError("Invalid mode error", InspectionModeError)
		return
	}

	var location wildfireanalysis.Location

	if !state.Location.Shared.IsNull() && state.Location.Shared.ValueBool() {
		location.Shared = true
	}
	if state.Location.DeviceGroup != nil {
		location.DeviceGroup = &wildfireanalysis.DeviceGroupLocation{

			PanoramaDevice: state.Location.DeviceGroup.PanoramaDevice.ValueString(),
			DeviceGroup:    state.Location.DeviceGroup.Name.ValueString(),
		}
	}

	err := r.manager.Delete(ctx, location, []string{state.Name.ValueString()})
	if err != nil && !errors.Is(err, sdkmanager.ErrObjectNotFound) {
		resp.Diagnostics.AddError("Error in delete", err.Error())
	}

}

type WildfireAnalysisSecurityProfileImportState struct {
	Location WildfireAnalysisSecurityProfileLocation `json:"location"`
	Name     string                                  `json:"name"`
}

func WildfireAnalysisSecurityProfileImportStateCreator(ctx context.Context, resource types.Object) ([]byte, error) {
	attrs := resource.Attributes()
	if attrs == nil {
		return nil, fmt.Errorf("Object has no attributes")
	}

	locationAttr, ok := attrs["location"]
	if !ok {
		return nil, fmt.Errorf("location attribute missing")
	}

	var location WildfireAnalysisSecurityProfileLocation
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

	importStruct := WildfireAnalysisSecurityProfileImportState{
		Location: location,
		Name:     name,
	}

	return json.Marshal(importStruct)
}

func (r *WildfireAnalysisSecurityProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	var obj WildfireAnalysisSecurityProfileImportState
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

type WildfireAnalysisSecurityProfileDeviceGroupLocation struct {
	PanoramaDevice types.String `tfsdk:"panorama_device"`
	Name           types.String `tfsdk:"name"`
}
type WildfireAnalysisSecurityProfileLocation struct {
	Shared      types.Bool                                          `tfsdk:"shared"`
	DeviceGroup *WildfireAnalysisSecurityProfileDeviceGroupLocation `tfsdk:"device_group"`
}

func WildfireAnalysisSecurityProfileLocationSchema() rsschema.Attribute {
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
						path.MatchRelative().AtParent().AtName("device_group"),
					}...),
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

func (o WildfireAnalysisSecurityProfileDeviceGroupLocation) MarshalJSON() ([]byte, error) {
	obj := struct {
		PanoramaDevice *string `json:"panorama_device"`
		Name           *string `json:"name"`
	}{
		PanoramaDevice: o.PanoramaDevice.ValueStringPointer(),
		Name:           o.Name.ValueStringPointer(),
	}

	return json.Marshal(obj)
}

func (o *WildfireAnalysisSecurityProfileDeviceGroupLocation) UnmarshalJSON(data []byte) error {
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
func (o WildfireAnalysisSecurityProfileLocation) MarshalJSON() ([]byte, error) {
	obj := struct {
		Shared      *bool                                               `json:"shared"`
		DeviceGroup *WildfireAnalysisSecurityProfileDeviceGroupLocation `json:"device_group"`
	}{
		Shared:      o.Shared.ValueBoolPointer(),
		DeviceGroup: o.DeviceGroup,
	}

	return json.Marshal(obj)
}

func (o *WildfireAnalysisSecurityProfileLocation) UnmarshalJSON(data []byte) error {
	var shadow struct {
		Shared      *bool                                               `json:"shared"`
		DeviceGroup *WildfireAnalysisSecurityProfileDeviceGroupLocation `json:"device_group"`
	}

	err := json.Unmarshal(data, &shadow)
	if err != nil {
		return err
	}
	o.Shared = types.BoolPointerValue(shadow.Shared)
	o.DeviceGroup = shadow.DeviceGroup

	return nil
}
