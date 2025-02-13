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
	"github.com/PaloAltoNetworks/pango/objects/profiles/secgroup"
	pangoutil "github.com/PaloAltoNetworks/pango/util"

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
	_ datasource.DataSource              = &SecurityProfileGroupDataSource{}
	_ datasource.DataSourceWithConfigure = &SecurityProfileGroupDataSource{}
)

func NewSecurityProfileGroupDataSource() datasource.DataSource {
	return &SecurityProfileGroupDataSource{}
}

type SecurityProfileGroupDataSource struct {
	client  *pango.Client
	manager *sdkmanager.EntryObjectManager[*secgroup.Entry, secgroup.Location, *secgroup.Service]
}

type SecurityProfileGroupDataSourceFilter struct {
	// TODO: Generate Data Source filter via function
}

type SecurityProfileGroupDataSourceModel struct {
	Location         SecurityProfileGroupLocation `tfsdk:"location"`
	Name             types.String                 `tfsdk:"name"`
	Sctp             types.List                   `tfsdk:"sctp"`
	Spyware          types.List                   `tfsdk:"spyware"`
	WildfireAnalysis types.List                   `tfsdk:"wildfire_analysis"`
	DataFiltering    types.List                   `tfsdk:"data_filtering"`
	Gtp              types.List                   `tfsdk:"gtp"`
	UrlFiltering     types.List                   `tfsdk:"url_filtering"`
	Virus            types.List                   `tfsdk:"virus"`
	Vulnerability    types.List                   `tfsdk:"vulnerability"`
	DisableOverride  types.String                 `tfsdk:"disable_override"`
	FileBlocking     types.List                   `tfsdk:"file_blocking"`
}

func (o *SecurityProfileGroupDataSourceModel) CopyToPango(ctx context.Context, obj **secgroup.Entry, encrypted *map[string]types.String) diag.Diagnostics {
	var diags diag.Diagnostics
	vulnerability_pango_entries := make([]string, 0)
	diags.Append(o.Vulnerability.ElementsAs(ctx, &vulnerability_pango_entries, false)...)
	if diags.HasError() {
		return diags
	}
	disableOverride_value := o.DisableOverride.ValueStringPointer()
	fileBlocking_pango_entries := make([]string, 0)
	diags.Append(o.FileBlocking.ElementsAs(ctx, &fileBlocking_pango_entries, false)...)
	if diags.HasError() {
		return diags
	}
	urlFiltering_pango_entries := make([]string, 0)
	diags.Append(o.UrlFiltering.ElementsAs(ctx, &urlFiltering_pango_entries, false)...)
	if diags.HasError() {
		return diags
	}
	virus_pango_entries := make([]string, 0)
	diags.Append(o.Virus.ElementsAs(ctx, &virus_pango_entries, false)...)
	if diags.HasError() {
		return diags
	}
	wildfireAnalysis_pango_entries := make([]string, 0)
	diags.Append(o.WildfireAnalysis.ElementsAs(ctx, &wildfireAnalysis_pango_entries, false)...)
	if diags.HasError() {
		return diags
	}
	dataFiltering_pango_entries := make([]string, 0)
	diags.Append(o.DataFiltering.ElementsAs(ctx, &dataFiltering_pango_entries, false)...)
	if diags.HasError() {
		return diags
	}
	gtp_pango_entries := make([]string, 0)
	diags.Append(o.Gtp.ElementsAs(ctx, &gtp_pango_entries, false)...)
	if diags.HasError() {
		return diags
	}
	sctp_pango_entries := make([]string, 0)
	diags.Append(o.Sctp.ElementsAs(ctx, &sctp_pango_entries, false)...)
	if diags.HasError() {
		return diags
	}
	spyware_pango_entries := make([]string, 0)
	diags.Append(o.Spyware.ElementsAs(ctx, &spyware_pango_entries, false)...)
	if diags.HasError() {
		return diags
	}

	if (*obj) == nil {
		*obj = new(secgroup.Entry)
	}
	(*obj).Name = o.Name.ValueString()
	(*obj).Vulnerability = vulnerability_pango_entries
	(*obj).DisableOverride = disableOverride_value
	(*obj).FileBlocking = fileBlocking_pango_entries
	(*obj).UrlFiltering = urlFiltering_pango_entries
	(*obj).Virus = virus_pango_entries
	(*obj).WildfireAnalysis = wildfireAnalysis_pango_entries
	(*obj).DataFiltering = dataFiltering_pango_entries
	(*obj).Gtp = gtp_pango_entries
	(*obj).Sctp = sctp_pango_entries
	(*obj).Spyware = spyware_pango_entries

	return diags
}

func (o *SecurityProfileGroupDataSourceModel) CopyFromPango(ctx context.Context, obj *secgroup.Entry, encrypted *map[string]types.String) diag.Diagnostics {
	var diags diag.Diagnostics
	var fileBlocking_list types.List
	{
		var list_diags diag.Diagnostics
		fileBlocking_list, list_diags = types.ListValueFrom(ctx, types.StringType, obj.FileBlocking)
		diags.Append(list_diags...)
	}
	var urlFiltering_list types.List
	{
		var list_diags diag.Diagnostics
		urlFiltering_list, list_diags = types.ListValueFrom(ctx, types.StringType, obj.UrlFiltering)
		diags.Append(list_diags...)
	}
	var virus_list types.List
	{
		var list_diags diag.Diagnostics
		virus_list, list_diags = types.ListValueFrom(ctx, types.StringType, obj.Virus)
		diags.Append(list_diags...)
	}
	var vulnerability_list types.List
	{
		var list_diags diag.Diagnostics
		vulnerability_list, list_diags = types.ListValueFrom(ctx, types.StringType, obj.Vulnerability)
		diags.Append(list_diags...)
	}
	var dataFiltering_list types.List
	{
		var list_diags diag.Diagnostics
		dataFiltering_list, list_diags = types.ListValueFrom(ctx, types.StringType, obj.DataFiltering)
		diags.Append(list_diags...)
	}
	var gtp_list types.List
	{
		var list_diags diag.Diagnostics
		gtp_list, list_diags = types.ListValueFrom(ctx, types.StringType, obj.Gtp)
		diags.Append(list_diags...)
	}
	var sctp_list types.List
	{
		var list_diags diag.Diagnostics
		sctp_list, list_diags = types.ListValueFrom(ctx, types.StringType, obj.Sctp)
		diags.Append(list_diags...)
	}
	var spyware_list types.List
	{
		var list_diags diag.Diagnostics
		spyware_list, list_diags = types.ListValueFrom(ctx, types.StringType, obj.Spyware)
		diags.Append(list_diags...)
	}
	var wildfireAnalysis_list types.List
	{
		var list_diags diag.Diagnostics
		wildfireAnalysis_list, list_diags = types.ListValueFrom(ctx, types.StringType, obj.WildfireAnalysis)
		diags.Append(list_diags...)
	}

	var disableOverride_value types.String
	if obj.DisableOverride != nil {
		disableOverride_value = types.StringValue(*obj.DisableOverride)
	}
	o.Name = types.StringValue(obj.Name)
	o.DisableOverride = disableOverride_value
	o.FileBlocking = fileBlocking_list
	o.UrlFiltering = urlFiltering_list
	o.Virus = virus_list
	o.Vulnerability = vulnerability_list
	o.DataFiltering = dataFiltering_list
	o.Gtp = gtp_list
	o.Sctp = sctp_list
	o.Spyware = spyware_list
	o.WildfireAnalysis = wildfireAnalysis_list

	return diags
}

func (o *SecurityProfileGroupDataSourceModel) resourceXpathComponents() ([]string, error) {
	var components []string
	components = append(components, pangoutil.AsEntryXpath(
		[]string{o.Name.ValueString()},
	))
	return components, nil
}

func SecurityProfileGroupDataSourceSchema() dsschema.Schema {
	return dsschema.Schema{
		Attributes: map[string]dsschema.Attribute{

			"location": SecurityProfileGroupDataSourceLocationSchema(),

			"name": dsschema.StringAttribute{
				Description: "",
				Computed:    false,
				Required:    true,
				Optional:    false,
				Sensitive:   false,
			},

			"file_blocking": dsschema.ListAttribute{
				Description: "",
				Required:    false,
				Optional:    true,
				Computed:    true,
				Sensitive:   false,
				ElementType: types.StringType,
			},

			"url_filtering": dsschema.ListAttribute{
				Description: "",
				Required:    false,
				Optional:    true,
				Computed:    true,
				Sensitive:   false,
				ElementType: types.StringType,
			},

			"virus": dsschema.ListAttribute{
				Description: "",
				Required:    false,
				Optional:    true,
				Computed:    true,
				Sensitive:   false,
				ElementType: types.StringType,
			},

			"vulnerability": dsschema.ListAttribute{
				Description: "",
				Required:    false,
				Optional:    true,
				Computed:    true,
				Sensitive:   false,
				ElementType: types.StringType,
			},

			"disable_override": dsschema.StringAttribute{
				Description: "disable object override in child device groups",
				Computed:    true,
				Required:    false,
				Optional:    true,
				Sensitive:   false,
			},

			"gtp": dsschema.ListAttribute{
				Description: "",
				Required:    false,
				Optional:    true,
				Computed:    true,
				Sensitive:   false,
				ElementType: types.StringType,
			},

			"sctp": dsschema.ListAttribute{
				Description: "",
				Required:    false,
				Optional:    true,
				Computed:    true,
				Sensitive:   false,
				ElementType: types.StringType,
			},

			"spyware": dsschema.ListAttribute{
				Description: "",
				Required:    false,
				Optional:    true,
				Computed:    true,
				Sensitive:   false,
				ElementType: types.StringType,
			},

			"wildfire_analysis": dsschema.ListAttribute{
				Description: "",
				Required:    false,
				Optional:    true,
				Computed:    true,
				Sensitive:   false,
				ElementType: types.StringType,
			},

			"data_filtering": dsschema.ListAttribute{
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

func (o *SecurityProfileGroupDataSourceModel) getTypeFor(name string) attr.Type {
	schema := SecurityProfileGroupDataSourceSchema()
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

func SecurityProfileGroupDataSourceLocationSchema() rsschema.Attribute {
	return SecurityProfileGroupLocationSchema()
}

// Metadata returns the data source type name.
func (d *SecurityProfileGroupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_security_profile_group"
}

// Schema defines the schema for this data source.
func (d *SecurityProfileGroupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = SecurityProfileGroupDataSourceSchema()
}

// Configure prepares the struct.
func (d *SecurityProfileGroupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*pango.Client)
	specifier, _, err := secgroup.Versioning(d.client.Versioning())
	if err != nil {
		resp.Diagnostics.AddError("Failed to configure SDK client", err.Error())
		return
	}
	d.manager = sdkmanager.NewEntryObjectManager(d.client, secgroup.NewService(d.client), specifier, secgroup.SpecMatches)
}
func (o *SecurityProfileGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var savestate, state SecurityProfileGroupDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &savestate)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var location secgroup.Location

	if !savestate.Location.Shared.IsNull() && savestate.Location.Shared.ValueBool() {
		location.Shared = true
	}
	if savestate.Location.DeviceGroup != nil {
		location.DeviceGroup = &secgroup.DeviceGroupLocation{

			PanoramaDevice: savestate.Location.DeviceGroup.PanoramaDevice.ValueString(),
			DeviceGroup:    savestate.Location.DeviceGroup.Name.ValueString(),
		}
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource read", map[string]any{
		"resource_name": "panos_security_profile_group_resource",
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
	_ resource.Resource                = &SecurityProfileGroupResource{}
	_ resource.ResourceWithConfigure   = &SecurityProfileGroupResource{}
	_ resource.ResourceWithImportState = &SecurityProfileGroupResource{}
)

func NewSecurityProfileGroupResource() resource.Resource {
	if _, found := resourceFuncMap["panos_security_profile_group"]; !found {
		resourceFuncMap["panos_security_profile_group"] = resourceFuncs{
			CreateImportId: SecurityProfileGroupImportStateCreator,
		}
	}
	return &SecurityProfileGroupResource{}
}

type SecurityProfileGroupResource struct {
	client  *pango.Client
	manager *sdkmanager.EntryObjectManager[*secgroup.Entry, secgroup.Location, *secgroup.Service]
}

func SecurityProfileGroupResourceLocationSchema() rsschema.Attribute {
	return SecurityProfileGroupLocationSchema()
}

type SecurityProfileGroupResourceModel struct {
	Location         SecurityProfileGroupLocation `tfsdk:"location"`
	Name             types.String                 `tfsdk:"name"`
	DataFiltering    types.List                   `tfsdk:"data_filtering"`
	Gtp              types.List                   `tfsdk:"gtp"`
	Sctp             types.List                   `tfsdk:"sctp"`
	Spyware          types.List                   `tfsdk:"spyware"`
	WildfireAnalysis types.List                   `tfsdk:"wildfire_analysis"`
	DisableOverride  types.String                 `tfsdk:"disable_override"`
	FileBlocking     types.List                   `tfsdk:"file_blocking"`
	UrlFiltering     types.List                   `tfsdk:"url_filtering"`
	Virus            types.List                   `tfsdk:"virus"`
	Vulnerability    types.List                   `tfsdk:"vulnerability"`
}

func (r *SecurityProfileGroupResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
}

// <ResourceSchema>

func SecurityProfileGroupResourceSchema() rsschema.Schema {
	return rsschema.Schema{
		Attributes: map[string]rsschema.Attribute{

			"location": SecurityProfileGroupResourceLocationSchema(),

			"name": rsschema.StringAttribute{
				Description: "",
				Computed:    false,
				Required:    true,
				Optional:    false,
				Sensitive:   false,
			},

			"wildfire_analysis": rsschema.ListAttribute{
				Description: "",
				Required:    false,
				Optional:    true,
				Computed:    false,
				Sensitive:   false,
				ElementType: types.StringType,
			},

			"data_filtering": rsschema.ListAttribute{
				Description: "",
				Required:    false,
				Optional:    true,
				Computed:    false,
				Sensitive:   false,
				ElementType: types.StringType,
			},

			"gtp": rsschema.ListAttribute{
				Description: "",
				Required:    false,
				Optional:    true,
				Computed:    false,
				Sensitive:   false,
				ElementType: types.StringType,
			},

			"sctp": rsschema.ListAttribute{
				Description: "",
				Required:    false,
				Optional:    true,
				Computed:    false,
				Sensitive:   false,
				ElementType: types.StringType,
			},

			"spyware": rsschema.ListAttribute{
				Description: "",
				Required:    false,
				Optional:    true,
				Computed:    false,
				Sensitive:   false,
				ElementType: types.StringType,
			},

			"vulnerability": rsschema.ListAttribute{
				Description: "",
				Required:    false,
				Optional:    true,
				Computed:    false,
				Sensitive:   false,
				ElementType: types.StringType,
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
						"yes",
						"no",
					}...),
				},
			},

			"file_blocking": rsschema.ListAttribute{
				Description: "",
				Required:    false,
				Optional:    true,
				Computed:    false,
				Sensitive:   false,
				ElementType: types.StringType,
			},

			"url_filtering": rsschema.ListAttribute{
				Description: "",
				Required:    false,
				Optional:    true,
				Computed:    false,
				Sensitive:   false,
				ElementType: types.StringType,
			},

			"virus": rsschema.ListAttribute{
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

func (o *SecurityProfileGroupResourceModel) getTypeFor(name string) attr.Type {
	schema := SecurityProfileGroupResourceSchema()
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

func (r *SecurityProfileGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_security_profile_group"
}

func (r *SecurityProfileGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = SecurityProfileGroupResourceSchema()
}

// </ResourceSchema>

func (r *SecurityProfileGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*pango.Client)
	specifier, _, err := secgroup.Versioning(r.client.Versioning())
	if err != nil {
		resp.Diagnostics.AddError("Failed to configure SDK client", err.Error())
		return
	}
	r.manager = sdkmanager.NewEntryObjectManager(r.client, secgroup.NewService(r.client), specifier, secgroup.SpecMatches)
}

func (o *SecurityProfileGroupResourceModel) CopyToPango(ctx context.Context, obj **secgroup.Entry, encrypted *map[string]types.String) diag.Diagnostics {
	var diags diag.Diagnostics
	dataFiltering_pango_entries := make([]string, 0)
	diags.Append(o.DataFiltering.ElementsAs(ctx, &dataFiltering_pango_entries, false)...)
	if diags.HasError() {
		return diags
	}
	gtp_pango_entries := make([]string, 0)
	diags.Append(o.Gtp.ElementsAs(ctx, &gtp_pango_entries, false)...)
	if diags.HasError() {
		return diags
	}
	sctp_pango_entries := make([]string, 0)
	diags.Append(o.Sctp.ElementsAs(ctx, &sctp_pango_entries, false)...)
	if diags.HasError() {
		return diags
	}
	spyware_pango_entries := make([]string, 0)
	diags.Append(o.Spyware.ElementsAs(ctx, &spyware_pango_entries, false)...)
	if diags.HasError() {
		return diags
	}
	wildfireAnalysis_pango_entries := make([]string, 0)
	diags.Append(o.WildfireAnalysis.ElementsAs(ctx, &wildfireAnalysis_pango_entries, false)...)
	if diags.HasError() {
		return diags
	}
	disableOverride_value := o.DisableOverride.ValueStringPointer()
	fileBlocking_pango_entries := make([]string, 0)
	diags.Append(o.FileBlocking.ElementsAs(ctx, &fileBlocking_pango_entries, false)...)
	if diags.HasError() {
		return diags
	}
	urlFiltering_pango_entries := make([]string, 0)
	diags.Append(o.UrlFiltering.ElementsAs(ctx, &urlFiltering_pango_entries, false)...)
	if diags.HasError() {
		return diags
	}
	virus_pango_entries := make([]string, 0)
	diags.Append(o.Virus.ElementsAs(ctx, &virus_pango_entries, false)...)
	if diags.HasError() {
		return diags
	}
	vulnerability_pango_entries := make([]string, 0)
	diags.Append(o.Vulnerability.ElementsAs(ctx, &vulnerability_pango_entries, false)...)
	if diags.HasError() {
		return diags
	}

	if (*obj) == nil {
		*obj = new(secgroup.Entry)
	}
	(*obj).Name = o.Name.ValueString()
	(*obj).DataFiltering = dataFiltering_pango_entries
	(*obj).Gtp = gtp_pango_entries
	(*obj).Sctp = sctp_pango_entries
	(*obj).Spyware = spyware_pango_entries
	(*obj).WildfireAnalysis = wildfireAnalysis_pango_entries
	(*obj).DisableOverride = disableOverride_value
	(*obj).FileBlocking = fileBlocking_pango_entries
	(*obj).UrlFiltering = urlFiltering_pango_entries
	(*obj).Virus = virus_pango_entries
	(*obj).Vulnerability = vulnerability_pango_entries

	return diags
}

func (o *SecurityProfileGroupResourceModel) CopyFromPango(ctx context.Context, obj *secgroup.Entry, encrypted *map[string]types.String) diag.Diagnostics {
	var diags diag.Diagnostics
	var fileBlocking_list types.List
	{
		var list_diags diag.Diagnostics
		fileBlocking_list, list_diags = types.ListValueFrom(ctx, types.StringType, obj.FileBlocking)
		diags.Append(list_diags...)
	}
	var urlFiltering_list types.List
	{
		var list_diags diag.Diagnostics
		urlFiltering_list, list_diags = types.ListValueFrom(ctx, types.StringType, obj.UrlFiltering)
		diags.Append(list_diags...)
	}
	var virus_list types.List
	{
		var list_diags diag.Diagnostics
		virus_list, list_diags = types.ListValueFrom(ctx, types.StringType, obj.Virus)
		diags.Append(list_diags...)
	}
	var vulnerability_list types.List
	{
		var list_diags diag.Diagnostics
		vulnerability_list, list_diags = types.ListValueFrom(ctx, types.StringType, obj.Vulnerability)
		diags.Append(list_diags...)
	}
	var dataFiltering_list types.List
	{
		var list_diags diag.Diagnostics
		dataFiltering_list, list_diags = types.ListValueFrom(ctx, types.StringType, obj.DataFiltering)
		diags.Append(list_diags...)
	}
	var gtp_list types.List
	{
		var list_diags diag.Diagnostics
		gtp_list, list_diags = types.ListValueFrom(ctx, types.StringType, obj.Gtp)
		diags.Append(list_diags...)
	}
	var sctp_list types.List
	{
		var list_diags diag.Diagnostics
		sctp_list, list_diags = types.ListValueFrom(ctx, types.StringType, obj.Sctp)
		diags.Append(list_diags...)
	}
	var spyware_list types.List
	{
		var list_diags diag.Diagnostics
		spyware_list, list_diags = types.ListValueFrom(ctx, types.StringType, obj.Spyware)
		diags.Append(list_diags...)
	}
	var wildfireAnalysis_list types.List
	{
		var list_diags diag.Diagnostics
		wildfireAnalysis_list, list_diags = types.ListValueFrom(ctx, types.StringType, obj.WildfireAnalysis)
		diags.Append(list_diags...)
	}

	var disableOverride_value types.String
	if obj.DisableOverride != nil {
		disableOverride_value = types.StringValue(*obj.DisableOverride)
	}
	o.Name = types.StringValue(obj.Name)
	o.DisableOverride = disableOverride_value
	o.FileBlocking = fileBlocking_list
	o.UrlFiltering = urlFiltering_list
	o.Virus = virus_list
	o.Vulnerability = vulnerability_list
	o.DataFiltering = dataFiltering_list
	o.Gtp = gtp_list
	o.Sctp = sctp_list
	o.Spyware = spyware_list
	o.WildfireAnalysis = wildfireAnalysis_list

	return diags
}

func (o *SecurityProfileGroupResourceModel) resourceXpathComponents() ([]string, error) {
	var components []string
	components = append(components, pangoutil.AsEntryXpath(
		[]string{o.Name.ValueString()},
	))
	return components, nil
}

func (r *SecurityProfileGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state SecurityProfileGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource create", map[string]any{
		"resource_name": "panos_security_profile_group_resource",
		"function":      "Create",
		"name":          state.Name.ValueString(),
	})

	// Verify mode.
	if r.client.Hostname == "" {
		resp.Diagnostics.AddError("Invalid mode error", InspectionModeError)
		return
	}

	// Determine the location.

	var location secgroup.Location

	if !state.Location.Shared.IsNull() && state.Location.Shared.ValueBool() {
		location.Shared = true
	}
	if state.Location.DeviceGroup != nil {
		location.DeviceGroup = &secgroup.DeviceGroupLocation{

			PanoramaDevice: state.Location.DeviceGroup.PanoramaDevice.ValueString(),
			DeviceGroup:    state.Location.DeviceGroup.Name.ValueString(),
		}
	}

	if err := location.IsValid(); err != nil {
		resp.Diagnostics.AddError("Invalid location", err.Error())
		return
	}

	// Load the desired config.
	var obj *secgroup.Entry

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
func (o *SecurityProfileGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var savestate, state SecurityProfileGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &savestate)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var location secgroup.Location

	if !savestate.Location.Shared.IsNull() && savestate.Location.Shared.ValueBool() {
		location.Shared = true
	}
	if savestate.Location.DeviceGroup != nil {
		location.DeviceGroup = &secgroup.DeviceGroupLocation{

			DeviceGroup:    savestate.Location.DeviceGroup.Name.ValueString(),
			PanoramaDevice: savestate.Location.DeviceGroup.PanoramaDevice.ValueString(),
		}
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource read", map[string]any{
		"resource_name": "panos_security_profile_group_resource",
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
func (r *SecurityProfileGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var plan, state SecurityProfileGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var location secgroup.Location

	if !state.Location.Shared.IsNull() && state.Location.Shared.ValueBool() {
		location.Shared = true
	}
	if state.Location.DeviceGroup != nil {
		location.DeviceGroup = &secgroup.DeviceGroupLocation{

			DeviceGroup:    state.Location.DeviceGroup.Name.ValueString(),
			PanoramaDevice: state.Location.DeviceGroup.PanoramaDevice.ValueString(),
		}
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource update", map[string]any{
		"resource_name": "panos_security_profile_group_resource",
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
func (r *SecurityProfileGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var state SecurityProfileGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource delete", map[string]any{
		"resource_name": "panos_security_profile_group_resource",
		"function":      "Delete",
		"name":          state.Name.ValueString(),
	})

	// Verify mode.
	if r.client.Hostname == "" {
		resp.Diagnostics.AddError("Invalid mode error", InspectionModeError)
		return
	}

	var location secgroup.Location

	if !state.Location.Shared.IsNull() && state.Location.Shared.ValueBool() {
		location.Shared = true
	}
	if state.Location.DeviceGroup != nil {
		location.DeviceGroup = &secgroup.DeviceGroupLocation{

			PanoramaDevice: state.Location.DeviceGroup.PanoramaDevice.ValueString(),
			DeviceGroup:    state.Location.DeviceGroup.Name.ValueString(),
		}
	}

	err := r.manager.Delete(ctx, location, []string{state.Name.ValueString()})
	if err != nil && !errors.Is(err, sdkmanager.ErrObjectNotFound) {
		resp.Diagnostics.AddError("Error in delete", err.Error())
	}

}

type SecurityProfileGroupImportState struct {
	Location SecurityProfileGroupLocation `json:"location"`
	Name     string                       `json:"name"`
}

func SecurityProfileGroupImportStateCreator(ctx context.Context, resource types.Object) ([]byte, error) {
	attrs := resource.Attributes()
	if attrs == nil {
		return nil, fmt.Errorf("Object has no attributes")
	}

	locationAttr, ok := attrs["location"]
	if !ok {
		return nil, fmt.Errorf("location attribute missing")
	}

	var location SecurityProfileGroupLocation
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

	importStruct := SecurityProfileGroupImportState{
		Location: location,
		Name:     name,
	}

	return json.Marshal(importStruct)
}

func (r *SecurityProfileGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	var obj SecurityProfileGroupImportState
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

type SecurityProfileGroupDeviceGroupLocation struct {
	PanoramaDevice types.String `tfsdk:"panorama_device"`
	Name           types.String `tfsdk:"name"`
}
type SecurityProfileGroupLocation struct {
	Shared      types.Bool                               `tfsdk:"shared"`
	DeviceGroup *SecurityProfileGroupDeviceGroupLocation `tfsdk:"device_group"`
}

func SecurityProfileGroupLocationSchema() rsschema.Attribute {
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
						path.MatchRelative().AtParent().AtName("device_group"),
						path.MatchRelative().AtParent().AtName("shared"),
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

func (o SecurityProfileGroupDeviceGroupLocation) MarshalJSON() ([]byte, error) {
	obj := struct {
		PanoramaDevice *string `json:"panorama_device"`
		Name           *string `json:"name"`
	}{
		PanoramaDevice: o.PanoramaDevice.ValueStringPointer(),
		Name:           o.Name.ValueStringPointer(),
	}

	return json.Marshal(obj)
}

func (o *SecurityProfileGroupDeviceGroupLocation) UnmarshalJSON(data []byte) error {
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
func (o SecurityProfileGroupLocation) MarshalJSON() ([]byte, error) {
	obj := struct {
		Shared      *bool                                    `json:"shared"`
		DeviceGroup *SecurityProfileGroupDeviceGroupLocation `json:"device_group"`
	}{
		Shared:      o.Shared.ValueBoolPointer(),
		DeviceGroup: o.DeviceGroup,
	}

	return json.Marshal(obj)
}

func (o *SecurityProfileGroupLocation) UnmarshalJSON(data []byte) error {
	var shadow struct {
		Shared      *bool                                    `json:"shared"`
		DeviceGroup *SecurityProfileGroupDeviceGroupLocation `json:"device_group"`
	}

	err := json.Unmarshal(data, &shadow)
	if err != nil {
		return err
	}
	o.Shared = types.BoolPointerValue(shadow.Shared)
	o.DeviceGroup = shadow.DeviceGroup

	return nil
}
