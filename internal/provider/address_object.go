package provider

import (
	"context"
	"fmt"
	"regexp"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/objects/address"

	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	rsschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Data source (listing).
var (
	_ datasource.DataSource              = &AddressObjectListDataSource{}
	_ datasource.DataSourceWithConfigure = &AddressObjectListDataSource{}
)

func NewAddressObjectListDataSource() datasource.DataSource {
	return &AddressObjectListDataSource{}
}

type AddressObjectListDataSource struct {
	client *pango.XmlApiClient
}

type AddressObjectDsListModel struct {
	// Input.
	Filter   *AddressObjectDsListFilter  `tfsdk:"filter"`
	Location AddressObjectDsListLocation `tfsdk:"location"`

	// Output.
	Data []AddressObjectDsListEntry `tfsdk:"data"`
}

type AddressObjectDsListFilter struct {
	Config types.String `tfsdk:"config"`
	Value  types.String `tfsdk:"value"`
	Quote  types.String `tfsdk:"quote"`
}

type AddressObjectDsListLocation struct {
	Shared       types.Bool                              `tfsdk:"shared"`
	Vsys         *AddressObjectDsListVsysLocation        `tfsdk:"vsys"`
	DeviceGroup  *AddressObjectDsListDeviceGroupLocation `tfsdk:"device_group"`
	FromPanorama types.Bool                              `tfsdk:"from_panorama"`
}

type AddressObjectDsListVsysLocation struct {
	NgfwDevice types.String `tfsdk:"ngfw_device"`
	Name       types.String `tfsdk:"name"`
}

type AddressObjectDsListDeviceGroupLocation struct {
	PanoramaDevice types.String `tfsdk:"panorama_device"`
	Name           types.String `tfsdk:"name"`
}

type AddressObjectDsListEntry struct {
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Tags        types.List   `tfsdk:"tags"`
	IpNetmask   types.String `tfsdk:"ip_netmask"`
	IpRange     types.String `tfsdk:"ip_range"`
	Fqdn        types.String `tfsdk:"fqdn"`
	IpWildcard  types.String `tfsdk:"ip_wildcard"`
}

func (d *AddressObjectListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_address_object_list"
}

func (d *AddressObjectListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Description: "Returns a list of address objects.",
		Attributes: map[string]dsschema.Attribute{
			"filter": dsschema.SingleNestedAttribute{
				Description: "Specify various properties about the read operation.",
				Optional:    true,
				Attributes: map[string]dsschema.Attribute{
					"config": dsschema.StringAttribute{
						Description: "Which type of config the data source should read from. If the provider is in local inspection mode, this param is ignored. Valid values are \"running\" or \"candidate\". Default: `\"candidate\"`.",
						Optional:    true,
					},
					"value": dsschema.StringAttribute{
						Description: "A filter to limit which objects are returned in the listing. Refer to the filter guide for more information.",
						Optional:    true,
					},
					"quote": dsschema.StringAttribute{
						Description: "The quote character for the given filter. Default: `'`.",
						Optional:    true,
					},
				},
			},
			"location": dsschema.SingleNestedAttribute{
				Description: "The location of this object. One and only one of the locations should be specified.",
				Required:    true,
				Attributes: map[string]dsschema.Attribute{
					"device_group": dsschema.SingleNestedAttribute{
						Description: "(Panorama) The given device group.",
						Optional:    true,
						Attributes: map[string]dsschema.Attribute{
							"name": dsschema.StringAttribute{
								Description: "The device group name.",
								Required:    true,
							},
							"panorama_device": dsschema.StringAttribute{
								Description: "The Panorama device. Default: `localhost.localdomain`.",
								Optional:    true,
							},
						},
					},
					"from_panorama": dsschema.BoolAttribute{
						Description: "(NGFW) Pushed from Panorama.",
						Optional:    true,
					},
					"shared": dsschema.BoolAttribute{
						Description: "(NGFW and Panorama) Located in shared.",
						Optional:    true,
					},
					"vsys": dsschema.SingleNestedAttribute{
						Description: "(NGFW) The given vsys.",
						Optional:    true,
						Attributes: map[string]dsschema.Attribute{
							"name": dsschema.StringAttribute{
								Description: "The vsys name.",
								Optional:    true,
							},
							"ngfw_device": dsschema.StringAttribute{
								Description: "The NGFW device. Default: `localhost.localdomain`.",
								Optional:    true,
							},
						},
					},
				},
			},
			"data": dsschema.ListNestedAttribute{
				Description: "The list of objects.",
				Computed:    true,
				NestedObject: dsschema.NestedAttributeObject{
					Attributes: map[string]dsschema.Attribute{
						"description": dsschema.StringAttribute{
							Description: "The description.",
							Computed:    true,
						},
						"fqdn": dsschema.StringAttribute{
							Description: "The Fqdn param.",
							Computed:    true,
						},
						"ip_netmask": dsschema.StringAttribute{
							Description: "The IpNetmask param.",
							Computed:    true,
						},
						"ip_range": dsschema.StringAttribute{
							Description: "The IpRange param.",
							Computed:    true,
						},
						"ip_wildcard": dsschema.StringAttribute{
							Description: "The IpWildcard param.",
							Computed:    true,
						},
						"name": dsschema.StringAttribute{
							Description: "Alphanumeric string [ 0-9a-zA-Z._-].",
							Computed:    true,
						},
						"tags": dsschema.ListAttribute{
							Description: "Tags for address object.",
							Computed:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
		},
	}
}

func (d *AddressObjectListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*pango.XmlApiClient)
}

func (d *AddressObjectListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state AddressObjectDsListModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Determine the location.
	var loc address.Location
	if !state.Location.Shared.IsNull() && state.Location.Shared.ValueBool() {
		loc.Shared = true
	} else if state.Location.Vsys != nil {
		loc.Vsys = &address.VsysLocation{}

		// NgfwDevice.
		if state.Location.Vsys.NgfwDevice.IsNull() {
			loc.Vsys.NgfwDevice = "localhost.localdomain"
		} else {
			loc.Vsys.NgfwDevice = state.Location.Vsys.NgfwDevice.ValueString()
		}

		// Name.
		if state.Location.Vsys.Name.IsNull() {
			loc.Vsys.Name = "vsys1"
		} else {
			loc.Vsys.Name = state.Location.Vsys.Name.ValueString()
		}
	} else if state.Location.DeviceGroup != nil {
		loc.DeviceGroup = &address.DeviceGroupLocation{}

		// PanoramaDevice.
		if state.Location.DeviceGroup.PanoramaDevice.IsNull() {
			loc.DeviceGroup.PanoramaDevice = "localhost.localdomain"
		} else {
			loc.DeviceGroup.PanoramaDevice = state.Location.DeviceGroup.PanoramaDevice.ValueString()
		}

		// Name.
		if state.Location.DeviceGroup.Name.IsNull() {
			resp.Diagnostics.AddError("Invalid location", "The device group name must be specified.")
			return
		}
		loc.DeviceGroup.Name = state.Location.DeviceGroup.Name.ValueString()
	} else if !state.Location.FromPanorama.IsNull() && state.Location.FromPanorama.ValueBool() {
		loc.FromPanorama = true
	} else {
		resp.Diagnostics.AddError("Unknown location", "Location for object is unknown")
		return
	}

	// Determine the rest of the List params.
	var action, filter, quote string
	if state.Filter == nil {
		action = "get"
	} else {
		// Action.
		if state.Filter.Config.IsNull() || state.Filter.Config.ValueString() == "candidate" {
			action = "get"
		} else if state.Filter.Config.ValueString() == "running" {
			action = "show"
		} else {
			resp.Diagnostics.AddError("Invalid filter.config", `The "filter.config" must be "candidate" or "running" if it is specified`)
			return
		}

		// Filter.
		filter = state.Filter.Value.ValueString()

		// Quote.
		if state.Filter.Quote.IsNull() {
			quote = `'`
		} else {
			quote = state.Filter.Quote.ValueString()
		}
	}

	// Create the service.
	svc := address.NewService(d.client)

	var err error
	var list []address.Entry

	// Perform the operation.
	if d.client.Hostname != "" {
		list, err = svc.List(ctx, loc, action, filter, quote)
	} else {
		list, err = svc.ListFromConfig(ctx, loc, filter, quote)
	}

	if err != nil {
		resp.Diagnostics.AddError("Error in read", err.Error())
		return
	}

	// Save to state.
	if len(list) == 0 {
		state.Data = nil
	} else {
		state.Data = make([]AddressObjectDsListEntry, 0, len(list))
		for _, var0 := range list {
			var1 := AddressObjectDsListEntry{}
			var1.Name = types.StringValue(var0.Name)
			var1.Description = types.StringPointerValue(var0.Description)
			var2, var3 := types.ListValueFrom(ctx, types.StringType, var0.Tags)
			var1.Tags = var2
			resp.Diagnostics.Append(var3.Errors()...)
			var1.IpNetmask = types.StringPointerValue(var0.IpNetmask)
			var1.IpRange = types.StringPointerValue(var0.IpRange)
			var1.Fqdn = types.StringPointerValue(var0.Fqdn)
			var1.IpWildcard = types.StringPointerValue(var0.IpWildcard)
			state.Data = append(state.Data, var1)
		}
	}

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Data source.
var (
	_ datasource.DataSource              = &AddressObjectDataSource{}
	_ datasource.DataSourceWithConfigure = &AddressObjectDataSource{}
)

func NewAddressObjectDataSource() datasource.DataSource {
	return &AddressObjectDataSource{}
}

type AddressObjectDataSource struct {
	client *pango.XmlApiClient
}

type AddressObjectDsModel struct {
	// Input.
	Filter   *AddressObjectDsFilter  `tfsdk:"filter"`
	Location AddressObjectDsLocation `tfsdk:"location"`
	Name     types.String            `tfsdk:"name"`

	// Output.
	Description types.String `tfsdk:"description"`
	Tags        types.List   `tfsdk:"tags"`
	IpNetmask   types.String `tfsdk:"ip_netmask"`
	IpRange     types.String `tfsdk:"ip_range"`
	Fqdn        types.String `tfsdk:"fqdn"`
	IpWildcard  types.String `tfsdk:"ip_wildcard"`
}

type AddressObjectDsFilter struct {
	Config types.String `tfsdk:"config"`
}

type AddressObjectDsLocation struct {
	Shared       types.Bool                          `tfsdk:"shared"`
	Vsys         *AddressObjectDsVsysLocation        `tfsdk:"vsys"`
	DeviceGroup  *AddressObjectDsDeviceGroupLocation `tfsdk:"device_group"`
	FromPanorama types.Bool                          `tfsdk:"from_panorama"`
}

type AddressObjectDsVsysLocation struct {
	NgfwDevice types.String `tfsdk:"ngfw_device"`
	Name       types.String `tfsdk:"name"`
}

type AddressObjectDsDeviceGroupLocation struct {
	PanoramaDevice types.String `tfsdk:"panorama_device"`
	Name           types.String `tfsdk:"name"`
}

func (d *AddressObjectDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_address_object"
}

func (d *AddressObjectDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Description: "Returns information about the given address object.",
		Attributes: map[string]dsschema.Attribute{
			// Input.
			"filter": dsschema.SingleNestedAttribute{
				Description: "Specify various properties about the read operation.",
				Optional:    true,
				Attributes: map[string]dsschema.Attribute{
					"config": dsschema.StringAttribute{
						Description: "Which type of config the data source should read from. If the provider is in local inspection mode, this param is ignored. Valid values are \"running\" or \"candidate\". Default: `\"candidate\"`.",
						Optional:    true,
					},
				},
			},
			"location": dsschema.SingleNestedAttribute{
				Description: "The location of this object. One and only one of the locations should be specified.",
				Required:    true,
				Attributes: map[string]dsschema.Attribute{
					"device_group": dsschema.SingleNestedAttribute{
						Description: "(Panorama) The given device group.",
						Optional:    true,
						Attributes: map[string]dsschema.Attribute{
							"name": dsschema.StringAttribute{
								Description: "The device group name.",
								Required:    true,
							},
							"panorama_device": dsschema.StringAttribute{
								Description: "The Panorama device. Default: `localhost.localdomain`.",
								Optional:    true,
							},
						},
					},
					"from_panorama": dsschema.BoolAttribute{
						Description: "(NGFW) Pushed from Panorama.",
						Optional:    true,
					},
					"shared": dsschema.BoolAttribute{
						Description: "(NGFW and Panorama) Located in shared.",
						Optional:    true,
					},
					"vsys": dsschema.SingleNestedAttribute{
						Description: "(NGFW) The given vsys.",
						Optional:    true,
						Attributes: map[string]dsschema.Attribute{
							"name": dsschema.StringAttribute{
								Description: "The vsys name. Default: `vsys1`.",
								Optional:    true,
							},
							"ngfw_device": dsschema.StringAttribute{
								Description: "The NGFW device. Default: `localhost.localdomain`.",
								Optional:    true,
							},
						},
					},
				},
			},
			"name": dsschema.StringAttribute{
				Description: "Alphanumeric string [ 0-9a-zA-Z._-].",
				Required:    true,
			},

			// Output.
			"description": dsschema.StringAttribute{
				Description: "The description.",
				Computed:    true,
			},
			"fqdn": dsschema.StringAttribute{
				Description: "The Fqdn param.",
				Computed:    true,
			},
			"ip_netmask": dsschema.StringAttribute{
				Description: "The IpNetmask param.",
				Computed:    true,
			},
			"ip_range": dsschema.StringAttribute{
				Description: "The IpRange param.",
				Computed:    true,
			},
			"ip_wildcard": dsschema.StringAttribute{
				Description: "The IpWildcard param.",
				Computed:    true,
			},
			"tags": dsschema.ListAttribute{
				Description: "Tags for address object.",
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (d *AddressObjectDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*pango.XmlApiClient)
}

func (d *AddressObjectDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state AddressObjectDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Determine the location.
	var loc address.Location
	if !state.Location.Shared.IsNull() && state.Location.Shared.ValueBool() {
		loc.Shared = true
	} else if state.Location.Vsys != nil {
		loc.Vsys = &address.VsysLocation{}

		// NgfwDevice.
		if state.Location.Vsys.NgfwDevice.IsNull() {
			loc.Vsys.NgfwDevice = "localhost.localdomain"
		} else {
			loc.Vsys.NgfwDevice = state.Location.Vsys.NgfwDevice.ValueString()
		}

		// Name.
		if state.Location.Vsys.Name.IsNull() {
			loc.Vsys.Name = "vsys1"
		} else {
			loc.Vsys.Name = state.Location.Vsys.Name.ValueString()
		}
	} else if state.Location.DeviceGroup != nil {
		loc.DeviceGroup = &address.DeviceGroupLocation{}

		// PanoramaDevice.
		if state.Location.DeviceGroup.PanoramaDevice.IsNull() {
			loc.DeviceGroup.PanoramaDevice = "localhost.localdomain"
		} else {
			loc.DeviceGroup.PanoramaDevice = state.Location.DeviceGroup.PanoramaDevice.ValueString()
		}

		// Name.
		if state.Location.DeviceGroup.Name.IsNull() {
			resp.Diagnostics.AddError("Invalid location", "The device group name must be specified.")
			return
		}
		loc.DeviceGroup.Name = state.Location.DeviceGroup.Name.ValueString()
	} else if !state.Location.FromPanorama.IsNull() && state.Location.FromPanorama.ValueBool() {
		loc.FromPanorama = true
	} else {
		resp.Diagnostics.AddError("Unknown location", "Location for object is unknown")
		return
	}

	// Determine the rest of the Read params.
	var action string
	if state.Filter == nil {
		action = "get"
	} else {
		if state.Filter.Config.IsNull() || state.Filter.Config.ValueString() == "candidate" {
			action = "get"
		} else if state.Filter.Config.ValueString() == "running" {
			action = "show"
		} else {
			resp.Diagnostics.AddError("Invalid filter.config", `The "filter.config" must be "candidate" or "running" if it is specified`)
			return
		}
	}

	// Create the service.
	svc := address.NewService(d.client)

	var err error
	var ans *address.Entry

	// Perform the operation.
	if d.client.Hostname != "" {
		ans, err = svc.Read(ctx, loc, state.Name.ValueString(), action)
	} else {
		ans, err = svc.ReadFromConfig(ctx, loc, state.Name.ValueString())
	}

	if err != nil {
		resp.Diagnostics.AddError("Error in read", err.Error())
		return
	}

	// Save the information to state.
	state.Name = types.StringValue(ans.Name)
	state.Description = types.StringPointerValue(ans.Description)
	var1, var2 := types.ListValueFrom(ctx, types.StringType, ans.Tags)
	state.Tags = var1
	resp.Diagnostics.Append(var2.Errors()...)
	state.IpNetmask = types.StringPointerValue(ans.IpNetmask)
	state.IpRange = types.StringPointerValue(ans.IpRange)
	state.Fqdn = types.StringPointerValue(ans.Fqdn)
	state.IpWildcard = types.StringPointerValue(ans.IpWildcard)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Resource.
var (
	_ resource.Resource                = &AddressObjectResource{}
	_ resource.ResourceWithConfigure   = &AddressObjectResource{}
	_ resource.ResourceWithImportState = &AddressObjectResource{}
)

func NewAddressObjectResource() resource.Resource {
	return &AddressObjectResource{}
}

type AddressObjectResource struct {
	client *pango.XmlApiClient
}

type AddressObjectTfid struct {
	Name     string           `json:"name"`
	Location address.Location `json:"location"`
}

func (o *AddressObjectTfid) IsValid() error {
	if o.Name == "" {
		return fmt.Errorf("name is unspecified")
	}

	return o.Location.IsValid()
}

type AddressObjectResourceModel struct {
	//Timeouts crudTimeouts `tfsdk:"timeouts"`
	Tfid types.String `tfsdk:"tfid"`

	Location AddressObjectResourceLocation `tfsdk:"location"`
	// old way:
	//Location nestedLocationModel `tfsdk:"location"`

	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Tags        types.List   `tfsdk:"tags"`
	IpNetmask   types.String `tfsdk:"ip_netmask"`
	IpRange     types.String `tfsdk:"ip_range"`
	Fqdn        types.String `tfsdk:"fqdn"`
	IpWildcard  types.String `tfsdk:"ip_wildcard"`
}

type AddressObjectResourceLocation struct {
	Shared       types.Bool                                `tfsdk:"shared"`
	FromPanorama types.Bool                                `tfsdk:"from_panorama"`
	Vsys         *AddressObjectResourceVsysLocation        `tfsdk:"vsys"`
	DeviceGroup  *AddressObjectResourceDeviceGroupLocation `tfsdk:"device_group"`
}

type AddressObjectResourceVsysLocation struct {
	Name       types.String `tfsdk:"name"`
	NgfwDevice types.String `tfsdk:"ngfw_device"`
}

type AddressObjectResourceDeviceGroupLocation struct {
	Name           types.String `tfsdk:"name"`
	PanoramaDevice types.String `tfsdk:"panorama_device"`
}

func (r *AddressObjectResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_address_object"
}

func (r *AddressObjectResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = rsschema.Schema{
		Description: "Manages an address object.  This is the \"nested\" style where the location is a struct.",

		Attributes: map[string]rsschema.Attribute{
			//"timeouts": CrudTimeoutsSchema("10m", "5m", "10m", "5m"),
			"location": rsschema.SingleNestedAttribute{
				Description: "The location of this object.",
				Required:    true,
				Attributes: map[string]rsschema.Attribute{
					"device_group": rsschema.SingleNestedAttribute{
						Description: "(Panorama) In the given device group. One of the following must be specified: `device_group`, `from_panorama`, `shared`, or `vsys`.",
						Optional:    true,
						Attributes: map[string]rsschema.Attribute{
							"name": rsschema.StringAttribute{
								Description: "The device group name.",
								Required:    true,
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.RequiresReplace(),
								},
							},
							"panorama_device": rsschema.StringAttribute{
								Description: "The Panorama device. Default: `localhost.localdomain`.",
								Optional:    true,
								Computed:    true,
								Default:     stringdefault.StaticString("localhost.localdomain"),
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.RequiresReplace(),
								},
							},
						},
					},
					"from_panorama": rsschema.BoolAttribute{
						Description: "(NGFW) Pushed from Panorama. This is a read-only location and only suitable for data sources. One of the following must be specified: `device_group`, `from_panorama`, `shared`, or `vsys`.",
						Optional:    true,
						Validators: []validator.Bool{
							boolvalidator.ExactlyOneOf(
								path.MatchRoot("location").AtName("from_panorama"),
								path.MatchRoot("location").AtName("device_group"),
								path.MatchRoot("location").AtName("vsys"),
								path.MatchRoot("location").AtName("shared"),
							),
						},
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.RequiresReplace(),
						},
					},
					"shared": rsschema.BoolAttribute{
						Description: "(NGFW and Panorama) Located in shared. One of the following must be specified: `device_group`, `from_panorama`, `shared`, or `vsys`.",
						Optional:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.RequiresReplace(),
						},
					},
					"vsys": rsschema.SingleNestedAttribute{
						Description: "(NGFW) In the given vsys. One of the following must be specified: `device_group`, `from_panorama`, `shared`, or `vsys`.",
						Optional:    true,
						Attributes: map[string]rsschema.Attribute{
							"name": rsschema.StringAttribute{
								Description: "The vsys name.",
								Optional:    true,
								Computed:    true,
								Default:     stringdefault.StaticString("vsys1"),
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.RequiresReplace(),
								},
							},
							"ngfw_device": rsschema.StringAttribute{
								Description: "The NGFW device.",
								Optional:    true,
								Computed:    true,
								Default:     stringdefault.StaticString("localhost.localdomain"),
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.RequiresReplace(),
								},
							},
						},
					},
				},
			},
			"tfid": rsschema.StringAttribute{
				Description: "The Terraform ID.",
				Computed:    true,
			},

			// Object properties.
			"name": rsschema.StringAttribute{
				Description: "Alphanumeric string [ 0-9a-zA-Z._-]. String length must not exceed 63 characters.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(63),
				},
			},
			"description": rsschema.StringAttribute{
				Description: "The description.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(1023),
				},
			},
			"fqdn": rsschema.StringAttribute{
				Description: "The Fqdn param. String length must be between 1 and 255 characters. String validation regex: `^[a-zA-Z0-9_]([a-zA-Z0-9._-])+[a-zA-Z0-9]$`. One of the following must be specified: `fqdn`, `ip_netmask`, `ip_range`, `ip_wildcard`",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 255),
					stringvalidator.RegexMatches(regexp.MustCompile("^[a-zA-Z0-9_]([a-zA-Z0-9._-])+[a-zA-Z0-9]$"), ""),
					stringvalidator.ExactlyOneOf(
						path.MatchRoot("fqdn"),
						path.MatchRoot("ip_netmask"),
						path.MatchRoot("ip_range"),
						path.MatchRoot("ip_wildcard"),
					),
				},
			},
			"ip_netmask": rsschema.StringAttribute{
				Description: "The IpNetmask param. One of the following must be specified: `fqdn`, `ip_netmask`, `ip_range`, `ip_wildcard`",
				Optional:    true,
			},
			"ip_range": rsschema.StringAttribute{
				Description: "The IpRange param. One of the following must be specified: `fqdn`, `ip_netmask`, `ip_range`, `ip_wildcard`",
				Optional:    true,
			},
			"ip_wildcard": rsschema.StringAttribute{
				Description: "The IpWildcard param. One of the following must be specified: `fqdn`, `ip_netmask`, `ip_range`, `ip_wildcard`",
				Optional:    true,
			},
			"tags": rsschema.ListAttribute{
				Description: "Tags for address object. List must contain at most 64 elements. Individual elements in this list are subject to additional validation. String length must not exceed 127 characters.",
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.List{
					listvalidator.SizeAtMost(64),
					listvalidator.ValueStringsAre(
						stringvalidator.LengthAtMost(127),
					),
				},
			},
		},
	}
}

func (r *AddressObjectResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*pango.XmlApiClient)
}

// Create.
func (r *AddressObjectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state AddressObjectResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource create", map[string]any{
		"resource_name": "panos_address_object",
		"function":      "Create",
		"name":          state.Name.ValueString(),
	})

	// Verify mode.
	if r.client.Hostname == "" {
		resp.Diagnostics.AddError("Invalid mode error", InspectionModeError)
		return
	}

	// Create the service.
	svc := address.NewService(r.client)

	// Determine the location.
	loc := AddressObjectTfid{Name: state.Name.ValueString()}
	if !state.Location.Shared.IsNull() && state.Location.Shared.ValueBool() {
		loc.Location.Shared = true
	}
	if !state.Location.FromPanorama.IsNull() && state.Location.FromPanorama.ValueBool() {
		loc.Location.FromPanorama = true
	}
	if state.Location.Vsys != nil {
		loc.Location.Vsys = &address.VsysLocation{}
		loc.Location.Vsys.NgfwDevice = state.Location.Vsys.NgfwDevice.ValueString()
		loc.Location.Vsys.Name = state.Location.Vsys.Name.ValueString()
	}
	if state.Location.DeviceGroup != nil {
		loc.Location.DeviceGroup = &address.DeviceGroupLocation{}
		loc.Location.DeviceGroup.Name = state.Location.DeviceGroup.Name.ValueString()
		loc.Location.DeviceGroup.PanoramaDevice = state.Location.DeviceGroup.PanoramaDevice.ValueString()
	}
	if err := loc.IsValid(); err != nil {
		resp.Diagnostics.AddError("Invalid location", err.Error())
		return
	}

	// Load the desired config.
	var obj address.Entry
	obj.Name = state.Name.ValueString()
	obj.Description = state.Description.ValueStringPointer()
	obj.IpNetmask = state.IpNetmask.ValueStringPointer()
	obj.IpRange = state.IpRange.ValueStringPointer()
	obj.Fqdn = state.Fqdn.ValueStringPointer()
	obj.IpWildcard = state.IpWildcard.ValueStringPointer()
	resp.Diagnostics.Append(state.Tags.ElementsAs(ctx, &obj.Tags, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	/*
		// Timeout handling.
		ctx, cancel := context.WithTimeout(ctx, GetTimeout(state.Timeouts.Create))
		defer cancel()
	*/

	// Perform the operation.
	ans, err := svc.Create(ctx, loc.Location, obj)
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
	state.Name = types.StringValue(ans.Name)
	state.Description = types.StringPointerValue(ans.Description)
	state.IpNetmask = types.StringPointerValue(ans.IpNetmask)
	state.IpRange = types.StringPointerValue(ans.IpRange)
	state.Fqdn = types.StringPointerValue(ans.Fqdn)
	state.IpWildcard = types.StringPointerValue(ans.IpWildcard)
	var1, var2 := types.ListValueFrom(ctx, types.StringType, ans.Tags)
	state.Tags = var1
	resp.Diagnostics.Append(var2.Errors()...)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Read performs Read for the struct.
func (r *AddressObjectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var savestate, state AddressObjectResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &savestate)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the location from tfid.
	var loc AddressObjectTfid
	if err := DecodeLocation(savestate.Tfid.ValueString(), &loc); err != nil {
		resp.Diagnostics.AddError("Error parsing tfid", err.Error())
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource read", map[string]any{
		"resource_name": "panos_address_object",
		"function":      "Read",
		"name":          loc.Name,
	})

	// Verify mode.
	if r.client.Hostname == "" {
		resp.Diagnostics.AddError("Invalid mode error", InspectionModeError)
		return
	}

	// Create the service.
	svc := address.NewService(r.client)

	/*
		// Timeout handling.
		ctx, cancel := context.WithTimeout(ctx, GetTimeout(savestate.Timeouts.Read))
		defer cancel()
	*/

	// Perform the operation.
	ans, err := svc.Read(ctx, loc.Location, loc.Name, "get")
	if err != nil {
		if IsObjectNotFound(err) {
			resp.State.RemoveResource(ctx)
		} else {
			resp.Diagnostics.AddError("Error reading config", err.Error())
		}
		return
	}

	// Save location to state.
	if loc.Location.Shared {
		state.Location.Shared = types.BoolValue(true)
	}
	if loc.Location.FromPanorama {
		state.Location.FromPanorama = types.BoolValue(true)
	}
	if loc.Location.Vsys != nil {
		state.Location.Vsys = &AddressObjectResourceVsysLocation{}
		state.Location.Vsys.Name = types.StringValue(loc.Location.Vsys.Name)
		state.Location.Vsys.NgfwDevice = types.StringValue(loc.Location.Vsys.NgfwDevice)
	}
	if loc.Location.DeviceGroup != nil {
		state.Location.DeviceGroup = &AddressObjectResourceDeviceGroupLocation{}
		state.Location.DeviceGroup.Name = types.StringValue(loc.Location.DeviceGroup.Name)
		state.Location.DeviceGroup.PanoramaDevice = types.StringValue(loc.Location.DeviceGroup.PanoramaDevice)
	}

	/*
		// Keep the timeouts.
	    // TODO: This won't work for state import.
		state.Timeouts = savestate.Timeouts
	*/

	// Save tfid to state.
	state.Tfid = savestate.Tfid

	// Save the answer to state.
	state.Name = types.StringValue(ans.Name)
	state.Description = types.StringPointerValue(ans.Description)
	state.IpNetmask = types.StringPointerValue(ans.IpNetmask)
	state.IpRange = types.StringPointerValue(ans.IpRange)
	state.Fqdn = types.StringPointerValue(ans.Fqdn)
	state.IpWildcard = types.StringPointerValue(ans.IpWildcard)
	var1, var2 := types.ListValueFrom(ctx, types.StringType, ans.Tags)
	state.Tags = var1
	resp.Diagnostics.Append(var2.Errors()...)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *AddressObjectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state AddressObjectResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var loc AddressObjectTfid
	if err := DecodeLocation(state.Tfid.ValueString(), &loc); err != nil {
		resp.Diagnostics.AddError("Error parsing tfid", err.Error())
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource update", map[string]any{
		"resource_name": "panos_address_object",
		"function":      "Update",
		"tfid":          state.Tfid.ValueString(),
	})

	// Verify mode.
	if r.client.Hostname == "" {
		resp.Diagnostics.AddError("Invalid mode error", InspectionModeError)
		return
	}

	// Create the service.
	svc := address.NewService(r.client)

	// Load the desired config.
	var obj address.Entry
	obj.Name = plan.Name.ValueString()
	obj.Description = plan.Description.ValueStringPointer()
	obj.IpNetmask = plan.IpNetmask.ValueStringPointer()
	obj.IpRange = plan.IpRange.ValueStringPointer()
	obj.Fqdn = plan.Fqdn.ValueStringPointer()
	obj.IpWildcard = plan.IpWildcard.ValueStringPointer()
	resp.Diagnostics.Append(plan.Tags.ElementsAs(ctx, &obj.Tags, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	/*
		// Timeout handling.
		ctx, cancel := context.WithTimeout(ctx, GetTimeout(plan.Timeouts.Update))
		defer cancel()
	*/

	// Perform the operation.
	ans, err := svc.Update(ctx, loc.Location, obj, loc.Name)
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

	// Save the state.
	state.Name = types.StringValue(ans.Name)
	state.Description = types.StringPointerValue(ans.Description)
	state.IpNetmask = types.StringPointerValue(ans.IpNetmask)
	state.IpRange = types.StringPointerValue(ans.IpRange)
	state.Fqdn = types.StringPointerValue(ans.Fqdn)
	state.IpWildcard = types.StringPointerValue(ans.IpWildcard)
	var1, var2 := types.ListValueFrom(ctx, types.StringType, ans.Tags)
	state.Tags = var1
	resp.Diagnostics.Append(var2.Errors()...)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *AddressObjectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state AddressObjectResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the location from tfid.
	var loc AddressObjectTfid
	if err := DecodeLocation(state.Tfid.ValueString(), &loc); err != nil {
		resp.Diagnostics.AddError("error parsing tfid", err.Error())
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource delete", map[string]any{
		"resource_name": "panos_address_object",
		"function":      "Delete",
		"name":          loc.Name,
	})

	// Verify mode.
	if r.client.Hostname == "" {
		resp.Diagnostics.AddError("Invalid mode error", InspectionModeError)
		return
	}

	// Create the service.
	svc := address.NewService(r.client)

	/*
		// Timeout handling.
		ctx, cancel := context.WithTimeout(ctx, GetTimeout(state.Timeouts.Delete))
		defer cancel()
	*/

	// Perform the operation.
	if err := svc.Delete(ctx, loc.Location, loc.Name); err != nil && !IsObjectNotFound(err) {
		resp.Diagnostics.AddError("Error in delete", err.Error())
	}
}

func (r *AddressObjectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("tfid"), req, resp)
}
