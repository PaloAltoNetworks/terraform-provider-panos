package provider

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/PaloAltoNetworks/pango"
	pangoutil "github.com/PaloAltoNetworks/pango/util"

	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
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
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	zone_protection "github.com/PaloAltoNetworks/pango/network/profiles/zone_protection"
	sdkmanager "github.com/PaloAltoNetworks/terraform-provider-panos/internal/manager"
)

// -----------------------------------------------------------------------
// DataSource
// -----------------------------------------------------------------------

var (
	_ datasource.DataSource              = &ZoneProtectionProfileDataSource{}
	_ datasource.DataSourceWithConfigure = &ZoneProtectionProfileDataSource{}
)

func NewZoneProtectionProfileDataSource() datasource.DataSource {
	return &ZoneProtectionProfileDataSource{}
}

type ZoneProtectionProfileDataSource struct {
	client  *pango.Client
	manager *sdkmanager.EntryObjectManager[*zone_protection.Entry, zone_protection.Location, *zone_protection.Service]
}

type ZoneProtectionProfileDataSourceModel struct {
	Location                   types.Object `tfsdk:"location"`
	Name                       types.String `tfsdk:"name"`
	Description                types.String `tfsdk:"description"`
	Flood                      types.Object `tfsdk:"flood"`
	Scan                       types.List   `tfsdk:"scan"`
	DiscardIpSpoof             types.Bool   `tfsdk:"discard_ip_spoof"`
	DiscardStrictSourceRouting types.Bool   `tfsdk:"discard_strict_source_routing"`
	DiscardLooseSourceRouting  types.Bool   `tfsdk:"discard_loose_source_routing"`
	DiscardMalformedOption     types.Bool   `tfsdk:"discard_malformed_option"`
	RemoveTcpTimestamp         types.Bool   `tfsdk:"remove_tcp_timestamp"`
	DiscardIpFrag              types.Bool   `tfsdk:"discard_ip_frag"`
	TcpSynWithData             types.Bool   `tfsdk:"tcp_syn_with_data"`
	StripTcpFastOpenAndData    types.Bool   `tfsdk:"strip_tcp_fast_open_and_data"`
	StripMptcpOption           types.String `tfsdk:"strip_mptcp_option"`
}

func (o *ZoneProtectionProfileDataSourceModel) AttributeTypes() map[string]attr.Type {
	var locationObj ZoneProtectionProfileLocation
	var floodObj *ZoneProtectionProfileFloodObject
	var scanObj *ZoneProtectionProfileScanObject
	return map[string]attr.Type{
		"location":    types.ObjectType{AttrTypes: locationObj.AttributeTypes()},
		"name":        types.StringType,
		"description": types.StringType,
		"flood":       types.ObjectType{AttrTypes: floodObj.AttributeTypes()},
		"scan":        types.ListType{ElemType: types.ObjectType{AttrTypes: scanObj.AttributeTypes()}},
		"discard_ip_spoof":              types.BoolType,
		"discard_strict_source_routing": types.BoolType,
		"discard_loose_source_routing":  types.BoolType,
		"discard_malformed_option":      types.BoolType,
		"remove_tcp_timestamp":          types.BoolType,
		"discard_ip_frag":               types.BoolType,
		"tcp_syn_with_data":             types.BoolType,
		"strip_tcp_fast_open_and_data":  types.BoolType,
		"strip_mptcp_option":            types.StringType,
	}
}

func (o ZoneProtectionProfileDataSourceModel) AncestorName() string { return "" }
func (o ZoneProtectionProfileDataSourceModel) EntryName() *string    { return nil }

func (o *ZoneProtectionProfileDataSourceModel) CopyToPango(ctx context.Context, client pangoutil.PangoClient, ancestors []Ancestor, obj **zone_protection.Entry, ev *EncryptedValuesManager) diag.Diagnostics {
	var diags diag.Diagnostics

	flood, d := copyFloodToPango(ctx, o.Flood)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	scan, d := copyScanToPango(ctx, o.Scan)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	if *obj == nil {
		*obj = new(zone_protection.Entry)
	}
	(*obj).Name = o.Name.ValueString()
	(*obj).Description = o.Description.ValueStringPointer()
	(*obj).Flood = flood
	(*obj).Scan = scan
	(*obj).DiscardIpSpoof = o.DiscardIpSpoof.ValueBoolPointer()
	(*obj).DiscardStrictSourceRouting = o.DiscardStrictSourceRouting.ValueBoolPointer()
	(*obj).DiscardLooseSourceRouting = o.DiscardLooseSourceRouting.ValueBoolPointer()
	(*obj).DiscardMalformedOption = o.DiscardMalformedOption.ValueBoolPointer()
	(*obj).RemoveTcpTimestamp = o.RemoveTcpTimestamp.ValueBoolPointer()
	(*obj).DiscardIpFrag = o.DiscardIpFrag.ValueBoolPointer()
	(*obj).TcpSynWithData = o.TcpSynWithData.ValueBoolPointer()
	(*obj).StripTcpFastOpenAndData = o.StripTcpFastOpenAndData.ValueBoolPointer()
	(*obj).StripMptcpOption = o.StripMptcpOption.ValueStringPointer()

	return diags
}

func (o *ZoneProtectionProfileDataSourceModel) CopyFromPango(ctx context.Context, client pangoutil.PangoClient, ancestors []Ancestor, obj *zone_protection.Entry, ev *EncryptedValuesManager) diag.Diagnostics {
	var diags diag.Diagnostics

	floodVal, d := copyFloodFromPango(ctx, obj.Flood, o.Flood)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	scanVal, d := copyScanFromPango(ctx, obj.Scan)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	o.Name = types.StringValue(obj.Name)
	o.Description = types.StringPointerValue(obj.Description)
	o.Flood = floodVal
	o.Scan = scanVal
	setBoolFromPtr(&o.DiscardIpSpoof, obj.DiscardIpSpoof)
	setBoolFromPtr(&o.DiscardStrictSourceRouting, obj.DiscardStrictSourceRouting)
	setBoolFromPtr(&o.DiscardLooseSourceRouting, obj.DiscardLooseSourceRouting)
	setBoolFromPtr(&o.DiscardMalformedOption, obj.DiscardMalformedOption)
	setBoolFromPtr(&o.RemoveTcpTimestamp, obj.RemoveTcpTimestamp)
	setBoolFromPtr(&o.DiscardIpFrag, obj.DiscardIpFrag)
	setBoolFromPtr(&o.TcpSynWithData, obj.TcpSynWithData)
	setBoolFromPtr(&o.StripTcpFastOpenAndData, obj.StripTcpFastOpenAndData)
	o.StripMptcpOption = types.StringPointerValue(obj.StripMptcpOption)

	return diags
}

func (o *ZoneProtectionProfileDataSourceModel) resourceXpathParentComponents() ([]string, error) {
	return []string{}, nil
}

func (o *ZoneProtectionProfileDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_zone_protection_profile"
}

func (o *ZoneProtectionProfileDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ZoneProtectionProfileDataSourceSchema()
}

func (o *ZoneProtectionProfileDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerData := req.ProviderData.(*ProviderData)
	o.client = providerData.Client
	specifier, _, err := zone_protection.Versioning(o.client.Versioning())
	if err != nil {
		resp.Diagnostics.AddError("Failed to configure SDK client", err.Error())
		return
	}
	batchSize := providerData.MultiConfigBatchSize
	o.manager = sdkmanager.NewEntryObjectManager[*zone_protection.Entry, zone_protection.Location, *zone_protection.Service](
		o.client, zone_protection.NewService(o.client), batchSize, specifier, zone_protection.SpecMatches,
	)
}

func (o *ZoneProtectionProfileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ZoneProtectionProfileDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var location zone_protection.Location
	resp.Diagnostics.Append(zoneProtectionProfileLocationFromTF(ctx, state.Location, &location)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "performing datasource read", map[string]any{
		"resource_name": "panos_zone_protection_profile_datasource",
		"function":      "Read",
		"name":          state.Name.ValueString(),
	})

	components, err := state.resourceXpathParentComponents()
	if err != nil {
		resp.Diagnostics.AddError("Error creating resource xpath", err.Error())
		return
	}
	object, err := o.manager.Read(ctx, location, components, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading entry", err.Error())
		return
	}

	resp.Diagnostics.Append(state.CopyFromPango(ctx, o.client, nil, object, nil)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func ZoneProtectionProfileDataSourceSchema() dsschema.Schema {
	return dsschema.Schema{
		Attributes: map[string]dsschema.Attribute{
			"location": ZoneProtectionProfileDataSourceLocationSchema(),
			"name": dsschema.StringAttribute{
				Description: "Zone protection profile name.",
				Required:    true,
			},
			"description": dsschema.StringAttribute{
				Description: "Description of the zone protection profile.",
				Optional:    true,
				Computed:    true,
			},
			"flood": dsschema.SingleNestedAttribute{
				Description: "Flood protection settings.",
				Optional:    true,
				Computed:    true,
				Attributes:  ZoneProtectionProfileDataSourceFloodSchema(),
			},
			"scan": dsschema.ListNestedAttribute{
				Description: "Reconnaissance protection (port scan / host sweep) entries.",
				Optional:    true,
				Computed:    true,
				NestedObject: dsschema.NestedAttributeObject{
					Attributes: ZoneProtectionProfileDataSourceScanSchema(),
				},
			},
			"discard_ip_spoof": dsschema.BoolAttribute{
				Description: "Discard IP spoofed packets.",
				Optional:    true,
				Computed:    true,
			},
			"discard_strict_source_routing": dsschema.BoolAttribute{
				Description: "Discard packets with strict source routing IP option.",
				Optional:    true,
				Computed:    true,
			},
			"discard_loose_source_routing": dsschema.BoolAttribute{
				Description: "Discard packets with loose source routing IP option.",
				Optional:    true,
				Computed:    true,
			},
			"discard_malformed_option": dsschema.BoolAttribute{
				Description: "Discard packets with malformed IP options.",
				Optional:    true,
				Computed:    true,
			},
			"remove_tcp_timestamp": dsschema.BoolAttribute{
				Description: "Remove TCP timestamp option.",
				Optional:    true,
				Computed:    true,
			},
			"discard_ip_frag": dsschema.BoolAttribute{
				Description: "Discard IP fragments.",
				Optional:    true,
				Computed:    true,
			},
			"tcp_syn_with_data": dsschema.BoolAttribute{
				Description: "Discard TCP SYN packets with data.",
				Optional:    true,
				Computed:    true,
			},
			"strip_tcp_fast_open_and_data": dsschema.BoolAttribute{
				Description: "Strip TCP fast open option and data.",
				Optional:    true,
				Computed:    true,
			},
			"strip_mptcp_option": dsschema.StringAttribute{
				Description: "Strip MPTCP option: global or never.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func ZoneProtectionProfileDataSourceFloodSchema() map[string]dsschema.Attribute {
	ratesAttrs := map[string]dsschema.Attribute{
		"alarm_rate": dsschema.Int64Attribute{
			Description: "Alarm rate (packets/sec).",
			Optional:    true,
			Computed:    true,
		},
		"activate_rate": dsschema.Int64Attribute{
			Description: "Activate rate (packets/sec).",
			Optional:    true,
			Computed:    true,
		},
		"maximal_rate": dsschema.Int64Attribute{
			Description: "Maximal rate (packets/sec).",
			Optional:    true,
			Computed:    true,
		},
	}
	protocolAttrs := map[string]dsschema.Attribute{
		"enable": dsschema.BoolAttribute{
			Description: "Enable flood protection for this protocol.",
			Optional:    true,
			Computed:    true,
		},
		"red": dsschema.SingleNestedAttribute{
			Description: "Random Early Drop settings.",
			Optional:    true,
			Computed:    true,
			Attributes:  ratesAttrs,
		},
	}
	synAttrs := map[string]dsschema.Attribute{
		"enable": dsschema.BoolAttribute{
			Description: "Enable SYN flood protection.",
			Optional:    true,
			Computed:    true,
		},
		"red": dsschema.SingleNestedAttribute{
			Description: "Random Early Drop settings.",
			Optional:    true,
			Computed:    true,
			Attributes:  ratesAttrs,
		},
		"syn_cookies": dsschema.SingleNestedAttribute{
			Description: "SYN Cookies settings.",
			Optional:    true,
			Computed:    true,
			Attributes:  ratesAttrs,
		},
	}
	return map[string]dsschema.Attribute{
		"syn": dsschema.SingleNestedAttribute{
			Description: "TCP SYN flood protection.",
			Optional:    true,
			Computed:    true,
			Attributes:  synAttrs,
		},
		"icmp": dsschema.SingleNestedAttribute{
			Description: "ICMP flood protection.",
			Optional:    true,
			Computed:    true,
			Attributes:  protocolAttrs,
		},
		"icmpv6": dsschema.SingleNestedAttribute{
			Description: "ICMPv6 flood protection.",
			Optional:    true,
			Computed:    true,
			Attributes:  protocolAttrs,
		},
		"udp": dsschema.SingleNestedAttribute{
			Description: "UDP flood protection.",
			Optional:    true,
			Computed:    true,
			Attributes:  protocolAttrs,
		},
		"other": dsschema.SingleNestedAttribute{
			Description: "Other IP flood protection.",
			Optional:    true,
			Computed:    true,
			Attributes:  protocolAttrs,
		},
	}
}

func ZoneProtectionProfileDataSourceLocationSchema() dsschema.Attribute {
	return dsschema.SingleNestedAttribute{
		Description: "The location of this object.",
		Required:    true,
		Attributes: map[string]dsschema.Attribute{
			"ngfw": dsschema.SingleNestedAttribute{
				Description: "Located in a specific NGFW device",
				Optional:    true,
				Attributes: map[string]dsschema.Attribute{
					"ngfw_device": dsschema.StringAttribute{
						Description: "The NGFW device",
						Optional:    true,
						Computed:    true,
					},
				},
			},
			"template": dsschema.SingleNestedAttribute{
				Description: "Located in a specific template",
				Optional:    true,
				Attributes: map[string]dsschema.Attribute{
					"panorama_device": dsschema.StringAttribute{
						Description: "Specific Panorama device",
						Optional:    true,
						Computed:    true,
					},
					"name": dsschema.StringAttribute{
						Description: "Specific Panorama template",
						Optional:    true,
						Computed:    true,
					},
					"ngfw_device": dsschema.StringAttribute{
						Description: "The NGFW device",
						Optional:    true,
						Computed:    true,
					},
				},
			},
			"template_stack": dsschema.SingleNestedAttribute{
				Description: "Located in a specific template stack",
				Optional:    true,
				Attributes: map[string]dsschema.Attribute{
					"panorama_device": dsschema.StringAttribute{
						Description: "Specific Panorama device",
						Optional:    true,
						Computed:    true,
					},
					"name": dsschema.StringAttribute{
						Description: "Specific Panorama template stack",
						Optional:    true,
						Computed:    true,
					},
					"ngfw_device": dsschema.StringAttribute{
						Description: "The NGFW device",
						Optional:    true,
						Computed:    true,
					},
				},
			},
		},
	}
}

// -----------------------------------------------------------------------
// Resource
// -----------------------------------------------------------------------

var (
	_ resource.Resource                = &ZoneProtectionProfileResource{}
	_ resource.ResourceWithConfigure   = &ZoneProtectionProfileResource{}
	_ resource.ResourceWithImportState = &ZoneProtectionProfileResource{}
)

func NewZoneProtectionProfileResource() resource.Resource {
	return &ZoneProtectionProfileResource{}
}

type ZoneProtectionProfileResource struct {
	client  *pango.Client
	manager *sdkmanager.EntryObjectManager[*zone_protection.Entry, zone_protection.Location, *zone_protection.Service]
}

type ZoneProtectionProfileResourceModel struct {
	Location                   types.Object `tfsdk:"location"`
	Name                       types.String `tfsdk:"name"`
	Description                types.String `tfsdk:"description"`
	Flood                      types.Object `tfsdk:"flood"`
	Scan                       types.List   `tfsdk:"scan"`
	DiscardIpSpoof             types.Bool   `tfsdk:"discard_ip_spoof"`
	DiscardStrictSourceRouting types.Bool   `tfsdk:"discard_strict_source_routing"`
	DiscardLooseSourceRouting  types.Bool   `tfsdk:"discard_loose_source_routing"`
	DiscardMalformedOption     types.Bool   `tfsdk:"discard_malformed_option"`
	RemoveTcpTimestamp         types.Bool   `tfsdk:"remove_tcp_timestamp"`
	DiscardIpFrag              types.Bool   `tfsdk:"discard_ip_frag"`
	TcpSynWithData             types.Bool   `tfsdk:"tcp_syn_with_data"`
	StripTcpFastOpenAndData    types.Bool   `tfsdk:"strip_tcp_fast_open_and_data"`
	StripMptcpOption           types.String `tfsdk:"strip_mptcp_option"`
}

func (o *ZoneProtectionProfileResourceModel) AttributeTypes() map[string]attr.Type {
	var locationObj ZoneProtectionProfileLocation
	var floodObj *ZoneProtectionProfileFloodObject
	var scanObj *ZoneProtectionProfileScanObject
	return map[string]attr.Type{
		"location":    types.ObjectType{AttrTypes: locationObj.AttributeTypes()},
		"name":        types.StringType,
		"description": types.StringType,
		"flood":       types.ObjectType{AttrTypes: floodObj.AttributeTypes()},
		"scan":        types.ListType{ElemType: types.ObjectType{AttrTypes: scanObj.AttributeTypes()}},
		"discard_ip_spoof":              types.BoolType,
		"discard_strict_source_routing": types.BoolType,
		"discard_loose_source_routing":  types.BoolType,
		"discard_malformed_option":      types.BoolType,
		"remove_tcp_timestamp":          types.BoolType,
		"discard_ip_frag":               types.BoolType,
		"tcp_syn_with_data":             types.BoolType,
		"strip_tcp_fast_open_and_data":  types.BoolType,
		"strip_mptcp_option":            types.StringType,
	}
}

func (o ZoneProtectionProfileResourceModel) AncestorName() string { return "" }
func (o ZoneProtectionProfileResourceModel) EntryName() *string    { return nil }

func (o *ZoneProtectionProfileResourceModel) ValidateConfig(ctx context.Context, resp *resource.ValidateConfigResponse, p path.Path) {
}

func (o *ZoneProtectionProfileResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var r ZoneProtectionProfileResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &r)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.ValidateConfig(ctx, resp, path.Empty())
}

func (o *ZoneProtectionProfileResourceModel) CopyToPango(ctx context.Context, client pangoutil.PangoClient, ancestors []Ancestor, obj **zone_protection.Entry, ev *EncryptedValuesManager) diag.Diagnostics {
	var diags diag.Diagnostics

	flood, d := copyFloodToPango(ctx, o.Flood)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	scan, d := copyScanToPango(ctx, o.Scan)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	if *obj == nil {
		*obj = new(zone_protection.Entry)
	}
	(*obj).Name = o.Name.ValueString()
	(*obj).Description = o.Description.ValueStringPointer()
	(*obj).Flood = flood
	(*obj).Scan = scan
	(*obj).DiscardIpSpoof = o.DiscardIpSpoof.ValueBoolPointer()
	(*obj).DiscardStrictSourceRouting = o.DiscardStrictSourceRouting.ValueBoolPointer()
	(*obj).DiscardLooseSourceRouting = o.DiscardLooseSourceRouting.ValueBoolPointer()
	(*obj).DiscardMalformedOption = o.DiscardMalformedOption.ValueBoolPointer()
	(*obj).RemoveTcpTimestamp = o.RemoveTcpTimestamp.ValueBoolPointer()
	(*obj).DiscardIpFrag = o.DiscardIpFrag.ValueBoolPointer()
	(*obj).TcpSynWithData = o.TcpSynWithData.ValueBoolPointer()
	(*obj).StripTcpFastOpenAndData = o.StripTcpFastOpenAndData.ValueBoolPointer()
	(*obj).StripMptcpOption = o.StripMptcpOption.ValueStringPointer()

	return diags
}

func (o *ZoneProtectionProfileResourceModel) CopyFromPango(ctx context.Context, client pangoutil.PangoClient, ancestors []Ancestor, obj *zone_protection.Entry, ev *EncryptedValuesManager) diag.Diagnostics {
	var diags diag.Diagnostics

	floodVal, d := copyFloodFromPango(ctx, obj.Flood, o.Flood)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	scanVal, d := copyScanFromPango(ctx, obj.Scan)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	o.Name = types.StringValue(obj.Name)
	o.Description = types.StringPointerValue(obj.Description)
	o.Flood = floodVal
	o.Scan = scanVal
	setBoolFromPtr(&o.DiscardIpSpoof, obj.DiscardIpSpoof)
	setBoolFromPtr(&o.DiscardStrictSourceRouting, obj.DiscardStrictSourceRouting)
	setBoolFromPtr(&o.DiscardLooseSourceRouting, obj.DiscardLooseSourceRouting)
	setBoolFromPtr(&o.DiscardMalformedOption, obj.DiscardMalformedOption)
	setBoolFromPtr(&o.RemoveTcpTimestamp, obj.RemoveTcpTimestamp)
	setBoolFromPtr(&o.DiscardIpFrag, obj.DiscardIpFrag)
	setBoolFromPtr(&o.TcpSynWithData, obj.TcpSynWithData)
	setBoolFromPtr(&o.StripTcpFastOpenAndData, obj.StripTcpFastOpenAndData)
	o.StripMptcpOption = types.StringPointerValue(obj.StripMptcpOption)

	return diags
}

func (o *ZoneProtectionProfileResourceModel) resourceXpathParentComponents() ([]string, error) {
	return []string{}, nil
}

// <ResourceSchema>

func ZoneProtectionProfileResourceSchema() rsschema.Schema {
	return rsschema.Schema{
		Attributes: map[string]rsschema.Attribute{
			"location": ZoneProtectionProfileResourceLocationSchema(),
			"name": rsschema.StringAttribute{
				Description: "Zone protection profile name.",
				Required:    true,
			},
			"description": rsschema.StringAttribute{
				Description: "Description of the zone protection profile.",
				Optional:    true,
			},
			"flood": rsschema.SingleNestedAttribute{
				Description: "Flood protection settings.",
				Optional:    true,
				Attributes:  ZoneProtectionProfileResourceFloodSchema(),
			},
			"scan": rsschema.ListNestedAttribute{
				Description: "Reconnaissance protection (port scan / host sweep) entries.",
				Optional:    true,
				NestedObject: rsschema.NestedAttributeObject{
					Attributes: ZoneProtectionProfileResourceScanSchema(),
				},
			},
			"discard_ip_spoof": rsschema.BoolAttribute{
				Description: "Discard IP spoofed packets.",
				Optional:    true,
			},
			"discard_strict_source_routing": rsschema.BoolAttribute{
				Description: "Discard packets with strict source routing IP option.",
				Optional:    true,
			},
			"discard_loose_source_routing": rsschema.BoolAttribute{
				Description: "Discard packets with loose source routing IP option.",
				Optional:    true,
			},
			"discard_malformed_option": rsschema.BoolAttribute{
				Description: "Discard packets with malformed IP options.",
				Optional:    true,
			},
			"remove_tcp_timestamp": rsschema.BoolAttribute{
				Description: "Remove TCP timestamp option.",
				Optional:    true,
			},
			"discard_ip_frag": rsschema.BoolAttribute{
				Description: "Discard IP fragments.",
				Optional:    true,
			},
			"tcp_syn_with_data": rsschema.BoolAttribute{
				Description: "Discard TCP SYN packets with data.",
				Optional:    true,
			},
			"strip_tcp_fast_open_and_data": rsschema.BoolAttribute{
				Description: "Strip TCP fast open option and data.",
				Optional:    true,
			},
			"strip_mptcp_option": rsschema.StringAttribute{
				Description: "Strip MPTCP option: global or never.",
				Optional:    true,
			},
		},
	}
}

func ZoneProtectionProfileResourceFloodSchema() map[string]rsschema.Attribute {
	ratesAttrs := map[string]rsschema.Attribute{
		"alarm_rate": rsschema.Int64Attribute{
			Description: "Alarm rate (packets/sec).",
			Optional:    true,
		},
		"activate_rate": rsschema.Int64Attribute{
			Description: "Activate rate (packets/sec).",
			Optional:    true,
		},
		"maximal_rate": rsschema.Int64Attribute{
			Description: "Maximal rate (packets/sec).",
			Optional:    true,
		},
	}
	protocolAttrs := map[string]rsschema.Attribute{
		"enable": rsschema.BoolAttribute{
			Description: "Enable flood protection for this protocol.",
			Optional:    true,
		},
		"red": rsschema.SingleNestedAttribute{
			Description: "Random Early Drop settings.",
			Optional:    true,
			Attributes:  ratesAttrs,
		},
	}
	synAttrs := map[string]rsschema.Attribute{
		"enable": rsschema.BoolAttribute{
			Description: "Enable SYN flood protection.",
			Optional:    true,
		},
		"red": rsschema.SingleNestedAttribute{
			Description: "Random Early Drop settings.",
			Optional:    true,
			Attributes:  ratesAttrs,
		},
		"syn_cookies": rsschema.SingleNestedAttribute{
			Description: "SYN Cookies settings.",
			Optional:    true,
			Attributes:  ratesAttrs,
		},
	}
	return map[string]rsschema.Attribute{
		"syn": rsschema.SingleNestedAttribute{
			Description: "TCP SYN flood protection.",
			Optional:    true,
			Attributes:  synAttrs,
		},
		"icmp": rsschema.SingleNestedAttribute{
			Description: "ICMP flood protection.",
			Optional:    true,
			Attributes:  protocolAttrs,
		},
		"icmpv6": rsschema.SingleNestedAttribute{
			Description: "ICMPv6 flood protection.",
			Optional:    true,
			Attributes:  protocolAttrs,
		},
		"udp": rsschema.SingleNestedAttribute{
			Description: "UDP flood protection.",
			Optional:    true,
			Attributes:  protocolAttrs,
		},
		"other": rsschema.SingleNestedAttribute{
			Description: "Other IP flood protection.",
			Optional:    true,
			Attributes:  protocolAttrs,
		},
	}
}

// </ResourceSchema>

func (o *ZoneProtectionProfileResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_zone_protection_profile"
}

func (o *ZoneProtectionProfileResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ZoneProtectionProfileResourceSchema()
}

func (o *ZoneProtectionProfileResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerData := req.ProviderData.(*ProviderData)
	o.client = providerData.Client
	specifier, _, err := zone_protection.Versioning(o.client.Versioning())
	if err != nil {
		resp.Diagnostics.AddError("Failed to configure SDK client", err.Error())
		return
	}
	batchSize := providerData.MultiConfigBatchSize
	o.manager = sdkmanager.NewEntryObjectManager[*zone_protection.Entry, zone_protection.Location, *zone_protection.Service](
		o.client, zone_protection.NewService(o.client), batchSize, specifier, zone_protection.SpecMatches,
	)
}

func (o *ZoneProtectionProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state ZoneProtectionProfileResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ev, err := NewEncryptedValuesManager(nil, false)
	if err != nil {
		resp.Diagnostics.AddError("Failed to init encrypted values manager", err.Error())
		return
	}

	var location zone_protection.Location
	resp.Diagnostics.Append(zoneProtectionProfileLocationFromTF(ctx, state.Location, &location)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := location.IsValid(); err != nil {
		resp.Diagnostics.AddError("Invalid location", err.Error())
		return
	}

	var obj *zone_protection.Entry
	resp.Diagnostics.Append(state.CopyToPango(ctx, o.client, nil, &obj, ev)...)
	if resp.Diagnostics.HasError() {
		return
	}

	components, err := state.resourceXpathParentComponents()
	if err != nil {
		resp.Diagnostics.AddError("Error creating resource xpath", err.Error())
		return
	}
	created, err := o.manager.Create(ctx, location, components, obj)
	if err != nil {
		resp.Diagnostics.AddError("Error in create", err.Error())
		return
	}

	resp.Diagnostics.Append(state.CopyFromPango(ctx, o.client, nil, created, ev)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload, err := json.Marshal(ev)
	if err != nil {
		resp.Diagnostics.AddError("Failed to marshal encrypted values state", err.Error())
		return
	}
	resp.Private.SetKey(ctx, "encrypted_values", payload)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (o *ZoneProtectionProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ZoneProtectionProfileResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	encryptedValues, diags := req.Private.GetKey(ctx, "encrypted_values")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	ev, err := NewEncryptedValuesManager(encryptedValues, true)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read encrypted values from private state", err.Error())
		return
	}

	var location zone_protection.Location
	resp.Diagnostics.Append(zoneProtectionProfileLocationFromTF(ctx, state.Location, &location)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "performing resource read", map[string]any{
		"resource_name": "panos_zone_protection_profile_resource",
		"function":      "Read",
		"name":          state.Name.ValueString(),
	})

	components, err := state.resourceXpathParentComponents()
	if err != nil {
		resp.Diagnostics.AddError("Error creating resource xpath", err.Error())
		return
	}
	object, err := o.manager.Read(ctx, location, components, state.Name.ValueString())
	if err != nil {
		if errors.Is(err, sdkmanager.ErrObjectNotFound) {
			resp.State.RemoveResource(ctx)
		} else {
			resp.Diagnostics.AddError("Error reading entry", err.Error())
		}
		return
	}

	resp.Diagnostics.Append(state.CopyFromPango(ctx, o.client, nil, object, ev)...)

	state.Location = state.Location

	payload, err := json.Marshal(ev)
	if err != nil {
		resp.Diagnostics.AddError("Failed to marshal encrypted values state", err.Error())
		return
	}
	resp.Private.SetKey(ctx, "encrypted_values", payload)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (o *ZoneProtectionProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state ZoneProtectionProfileResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	encryptedValues, diags := req.Private.GetKey(ctx, "encrypted_values")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	ev, err := NewEncryptedValuesManager(encryptedValues, false)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read encrypted values from private state", err.Error())
		return
	}

	var location zone_protection.Location
	resp.Diagnostics.Append(zoneProtectionProfileLocationFromTF(ctx, state.Location, &location)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "performing resource update", map[string]any{
		"resource_name": "panos_zone_protection_profile_resource",
		"function":      "Update",
	})

	if o.client.Hostname == "" {
		resp.Diagnostics.AddError("Invalid mode error", InspectionModeError)
		return
	}

	components, err := state.resourceXpathParentComponents()
	if err != nil {
		resp.Diagnostics.AddError("Error creating resource xpath", err.Error())
		return
	}

	var obj *zone_protection.Entry
	if state.Name.ValueString() != plan.Name.ValueString() {
		obj, err = o.manager.Read(ctx, location, components, state.Name.ValueString())
	} else {
		obj, err = o.manager.Read(ctx, location, components, plan.Name.ValueString())
	}
	if err != nil {
		resp.Diagnostics.AddError("Error in update", err.Error())
		return
	}

	resp.Diagnostics.Append(plan.CopyToPango(ctx, o.client, nil, &obj, ev)...)
	if resp.Diagnostics.HasError() {
		return
	}

	components, err = plan.resourceXpathParentComponents()
	if err != nil {
		resp.Diagnostics.AddError("Error creating resource xpath", err.Error())
		return
	}

	var newName string
	if state.Name.ValueString() != plan.Name.ValueString() {
		newName = plan.Name.ValueString()
		obj.Name = state.Name.ValueString()
	}

	updated, err := o.manager.Update(ctx, location, components, obj, newName)
	if err != nil {
		resp.Diagnostics.AddError("Error in update", err.Error())
		return
	}

	resp.Diagnostics.Append(plan.CopyFromPango(ctx, o.client, nil, updated, ev)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.Location = state.Location

	payload, err := json.Marshal(ev)
	if err != nil {
		resp.Diagnostics.AddError("Failed to marshal encrypted values state", err.Error())
		return
	}
	resp.Private.SetKey(ctx, "encrypted_values", payload)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (o *ZoneProtectionProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ZoneProtectionProfileResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var location zone_protection.Location
	resp.Diagnostics.Append(zoneProtectionProfileLocationFromTF(ctx, state.Location, &location)...)
	if resp.Diagnostics.HasError() {
		return
	}

	components, err := state.resourceXpathParentComponents()
	if err != nil {
		resp.Diagnostics.AddError("Error creating resource xpath", err.Error())
		return
	}
	err = o.manager.Delete(ctx, location, components, []string{state.Name.ValueString()})
	if err != nil && !errors.Is(err, sdkmanager.ErrObjectNotFound) {
		resp.Diagnostics.AddError("Error in delete", err.Error())
	}
}

// ImportState

type ZoneProtectionProfileImportState struct {
	Location types.Object `json:"location"`
	Name     types.String `json:"name"`
}

func (o ZoneProtectionProfileImportState) MarshalJSON() ([]byte, error) {
	type shadow struct {
		Location interface{} `json:"location"`
		Name     *string     `json:"name"`
	}
	var location_object interface{}
	{
		var err error
		location_object, err = TypesObjectToMap(o.Location, ZoneProtectionProfileResourceLocationSchema())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal location into JSON document: %w", err)
		}
	}
	obj := shadow{
		Location: location_object,
		Name:     o.Name.ValueStringPointer(),
	}
	return json.Marshal(obj)
}

func (o *ZoneProtectionProfileImportState) UnmarshalJSON(data []byte) error {
	var shadow struct {
		Location interface{} `json:"location"`
		Name     *string     `json:"name"`
	}
	err := json.Unmarshal(data, &shadow)
	if err != nil {
		return err
	}
	var location_object types.Object
	{
		location_map, ok := shadow.Location.(map[string]interface{})
		if !ok {
			return NewDiagnosticsError("Failed to unmarshal JSON document into location: expected map[string]interface{}", nil)
		}
		var err error
		location_object, err = MapToTypesObject(location_map, ZoneProtectionProfileResourceLocationSchema())
		if err != nil {
			return fmt.Errorf("failed to unmarshal location from JSON: %w", err)
		}
	}
	o.Location = location_object
	o.Name = types.StringPointerValue(shadow.Name)
	return nil
}

func ZoneProtectionProfileImportStateCreator(ctx context.Context, resource types.Object) ([]byte, error) {
	attrs := resource.Attributes()
	if attrs == nil {
		return nil, fmt.Errorf("Object has no attributes")
	}

	locationAttr, ok := attrs["location"]
	if !ok {
		return nil, fmt.Errorf("location attribute missing")
	}
	var location types.Object
	switch value := locationAttr.(type) {
	case types.Object:
		location = value
	default:
		return nil, fmt.Errorf("location attribute expected to be an object")
	}

	nameAttr, ok := attrs["name"]
	if !ok {
		return nil, fmt.Errorf("name attribute missing")
	}
	var name types.String
	switch value := nameAttr.(type) {
	case types.String:
		name = value
	default:
		return nil, fmt.Errorf("name attribute expected to be a string")
	}

	importStruct := ZoneProtectionProfileImportState{
		Location: location,
		Name:     name,
	}
	return json.Marshal(importStruct)
}

func (o *ZoneProtectionProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var obj ZoneProtectionProfileImportState
	data, err := base64.StdEncoding.DecodeString(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Failed to decode Import ID", err.Error())
		return
	}

	err = json.Unmarshal(data, &obj)
	if err != nil {
		var diagsErr *DiagnosticsError
		if errors.As(err, &diagsErr) {
			resp.Diagnostics.Append(diagsErr.Diagnostics()...)
		} else {
			resp.Diagnostics.AddError("Failed to unmarshal Import ID", err.Error())
		}
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("location"), obj.Location)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), obj.Name)...)
}

// -----------------------------------------------------------------------
// Location types (shared between DS and Resource)
// -----------------------------------------------------------------------

type ZoneProtectionProfileNgfwLocation struct {
	NgfwDevice types.String `tfsdk:"ngfw_device"`
}

type ZoneProtectionProfileTemplateLocation struct {
	PanoramaDevice types.String `tfsdk:"panorama_device"`
	Name           types.String `tfsdk:"name"`
	NgfwDevice     types.String `tfsdk:"ngfw_device"`
}

type ZoneProtectionProfileTemplateStackLocation struct {
	PanoramaDevice types.String `tfsdk:"panorama_device"`
	Name           types.String `tfsdk:"name"`
	NgfwDevice     types.String `tfsdk:"ngfw_device"`
}

type ZoneProtectionProfileLocation struct {
	Ngfw          types.Object `tfsdk:"ngfw"`
	Template      types.Object `tfsdk:"template"`
	TemplateStack types.Object `tfsdk:"template_stack"`
}

func ZoneProtectionProfileResourceLocationSchema() rsschema.Attribute {
	return rsschema.SingleNestedAttribute{
		Description: "The location of this object.",
		Required:    true,
		Attributes: map[string]rsschema.Attribute{
			"ngfw": rsschema.SingleNestedAttribute{
				Description: "Located in a specific NGFW device",
				Optional:    true,
				Attributes: map[string]rsschema.Attribute{
					"ngfw_device": rsschema.StringAttribute{
						Description: "The NGFW device",
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
				Validators: []validator.Object{
					objectvalidator.ExactlyOneOf(path.Expressions{
						path.MatchRelative().AtParent().AtName("ngfw"),
						path.MatchRelative().AtParent().AtName("template"),
						path.MatchRelative().AtParent().AtName("template_stack"),
					}...),
				},
			},
			"template": rsschema.SingleNestedAttribute{
				Description: "Located in a specific template",
				Optional:    true,
				Attributes: map[string]rsschema.Attribute{
					"panorama_device": rsschema.StringAttribute{
						Description: "Specific Panorama device",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString("localhost.localdomain"),
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"name": rsschema.StringAttribute{
						Description: "Specific Panorama template",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(""),
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"ngfw_device": rsschema.StringAttribute{
						Description: "The NGFW device",
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
			"template_stack": rsschema.SingleNestedAttribute{
				Description: "Located in a specific template stack",
				Optional:    true,
				Attributes: map[string]rsschema.Attribute{
					"panorama_device": rsschema.StringAttribute{
						Description: "Specific Panorama device",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString("localhost.localdomain"),
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"name": rsschema.StringAttribute{
						Description: "Specific Panorama template stack",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(""),
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"ngfw_device": rsschema.StringAttribute{
						Description: "The NGFW device",
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

func (o ZoneProtectionProfileNgfwLocation) MarshalJSON() ([]byte, error) {
	type shadow struct {
		NgfwDevice *string `json:"ngfw_device,omitempty"`
	}
	return json.Marshal(shadow{NgfwDevice: o.NgfwDevice.ValueStringPointer()})
}

func (o *ZoneProtectionProfileNgfwLocation) UnmarshalJSON(data []byte) error {
	var shadow struct {
		NgfwDevice *string `json:"ngfw_device,omitempty"`
	}
	if err := json.Unmarshal(data, &shadow); err != nil {
		return err
	}
	o.NgfwDevice = types.StringPointerValue(shadow.NgfwDevice)
	return nil
}

func (o ZoneProtectionProfileTemplateLocation) MarshalJSON() ([]byte, error) {
	type shadow struct {
		PanoramaDevice *string `json:"panorama_device,omitempty"`
		Name           *string `json:"name,omitempty"`
		NgfwDevice     *string `json:"ngfw_device,omitempty"`
	}
	return json.Marshal(shadow{
		PanoramaDevice: o.PanoramaDevice.ValueStringPointer(),
		Name:           o.Name.ValueStringPointer(),
		NgfwDevice:     o.NgfwDevice.ValueStringPointer(),
	})
}

func (o *ZoneProtectionProfileTemplateLocation) UnmarshalJSON(data []byte) error {
	var shadow struct {
		PanoramaDevice *string `json:"panorama_device,omitempty"`
		Name           *string `json:"name,omitempty"`
		NgfwDevice     *string `json:"ngfw_device,omitempty"`
	}
	if err := json.Unmarshal(data, &shadow); err != nil {
		return err
	}
	o.PanoramaDevice = types.StringPointerValue(shadow.PanoramaDevice)
	o.Name = types.StringPointerValue(shadow.Name)
	o.NgfwDevice = types.StringPointerValue(shadow.NgfwDevice)
	return nil
}

func (o ZoneProtectionProfileTemplateStackLocation) MarshalJSON() ([]byte, error) {
	type shadow struct {
		PanoramaDevice *string `json:"panorama_device,omitempty"`
		Name           *string `json:"name,omitempty"`
		NgfwDevice     *string `json:"ngfw_device,omitempty"`
	}
	return json.Marshal(shadow{
		PanoramaDevice: o.PanoramaDevice.ValueStringPointer(),
		Name:           o.Name.ValueStringPointer(),
		NgfwDevice:     o.NgfwDevice.ValueStringPointer(),
	})
}

func (o *ZoneProtectionProfileTemplateStackLocation) UnmarshalJSON(data []byte) error {
	var shadow struct {
		PanoramaDevice *string `json:"panorama_device,omitempty"`
		Name           *string `json:"name,omitempty"`
		NgfwDevice     *string `json:"ngfw_device,omitempty"`
	}
	if err := json.Unmarshal(data, &shadow); err != nil {
		return err
	}
	o.PanoramaDevice = types.StringPointerValue(shadow.PanoramaDevice)
	o.Name = types.StringPointerValue(shadow.Name)
	o.NgfwDevice = types.StringPointerValue(shadow.NgfwDevice)
	return nil
}

func (o ZoneProtectionProfileLocation) AttributeTypes() map[string]attr.Type {
	var ngfwObj ZoneProtectionProfileNgfwLocation
	var templateObj ZoneProtectionProfileTemplateLocation
	var templateStackObj ZoneProtectionProfileTemplateStackLocation
	return map[string]attr.Type{
		"ngfw":           types.ObjectType{AttrTypes: ngfwObj.AttributeTypes()},
		"template":       types.ObjectType{AttrTypes: templateObj.AttributeTypes()},
		"template_stack": types.ObjectType{AttrTypes: templateStackObj.AttributeTypes()},
	}
}

func (o ZoneProtectionProfileNgfwLocation) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{"ngfw_device": types.StringType}
}

func (o ZoneProtectionProfileTemplateLocation) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"panorama_device": types.StringType,
		"name":            types.StringType,
		"ngfw_device":     types.StringType,
	}
}

func (o ZoneProtectionProfileTemplateStackLocation) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"panorama_device": types.StringType,
		"name":            types.StringType,
		"ngfw_device":     types.StringType,
	}
}

// zoneProtectionProfileLocationFromTF converts a TF location object into the SDK location type.
func zoneProtectionProfileLocationFromTF(ctx context.Context, locationObj types.Object, location *zone_protection.Location) diag.Diagnostics {
	var diags diag.Diagnostics

	var terraformLocation ZoneProtectionProfileLocation
	diags.Append(locationObj.As(ctx, &terraformLocation, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return diags
	}

	if !terraformLocation.Ngfw.IsNull() {
		location.Ngfw = &zone_protection.NgfwLocation{}
		var inner ZoneProtectionProfileNgfwLocation
		diags.Append(terraformLocation.Ngfw.As(ctx, &inner, basetypes.ObjectAsOptions{})...)
		if diags.HasError() {
			return diags
		}
		location.Ngfw.NgfwDevice = inner.NgfwDevice.ValueString()
	}

	if !terraformLocation.Template.IsNull() {
		location.Template = &zone_protection.TemplateLocation{}
		var inner ZoneProtectionProfileTemplateLocation
		diags.Append(terraformLocation.Template.As(ctx, &inner, basetypes.ObjectAsOptions{})...)
		if diags.HasError() {
			return diags
		}
		location.Template.PanoramaDevice = inner.PanoramaDevice.ValueString()
		location.Template.Template = inner.Name.ValueString()
		location.Template.NgfwDevice = inner.NgfwDevice.ValueString()
	}

	if !terraformLocation.TemplateStack.IsNull() {
		location.TemplateStack = &zone_protection.TemplateStackLocation{}
		var inner ZoneProtectionProfileTemplateStackLocation
		diags.Append(terraformLocation.TemplateStack.As(ctx, &inner, basetypes.ObjectAsOptions{})...)
		if diags.HasError() {
			return diags
		}
		location.TemplateStack.PanoramaDevice = inner.PanoramaDevice.ValueString()
		location.TemplateStack.TemplateStack = inner.Name.ValueString()
		location.TemplateStack.NgfwDevice = inner.NgfwDevice.ValueString()
	}

	return diags
}

// -----------------------------------------------------------------------
// Flood object types
// -----------------------------------------------------------------------

type ZoneProtectionProfileFloodObject struct {
	Syn    types.Object `tfsdk:"syn"`
	Icmp   types.Object `tfsdk:"icmp"`
	Icmpv6 types.Object `tfsdk:"icmpv6"`
	Udp    types.Object `tfsdk:"udp"`
	Other  types.Object `tfsdk:"other"`
}

type ZoneProtectionProfileFloodSynObject struct {
	Enable     types.Bool   `tfsdk:"enable"`
	Red        types.Object `tfsdk:"red"`
	SynCookies types.Object `tfsdk:"syn_cookies"`
}

type ZoneProtectionProfileFloodProtocolObject struct {
	Enable types.Bool   `tfsdk:"enable"`
	Red    types.Object `tfsdk:"red"`
}

type ZoneProtectionProfileFloodRatesObject struct {
	AlarmRate    types.Int64 `tfsdk:"alarm_rate"`
	ActivateRate types.Int64 `tfsdk:"activate_rate"`
	MaximalRate  types.Int64 `tfsdk:"maximal_rate"`
}

func (o *ZoneProtectionProfileFloodRatesObject) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"alarm_rate":    types.Int64Type,
		"activate_rate": types.Int64Type,
		"maximal_rate":  types.Int64Type,
	}
}

func (o *ZoneProtectionProfileFloodProtocolObject) AttributeTypes() map[string]attr.Type {
	var ratesObj *ZoneProtectionProfileFloodRatesObject
	return map[string]attr.Type{
		"enable": types.BoolType,
		"red":    types.ObjectType{AttrTypes: ratesObj.AttributeTypes()},
	}
}

func (o *ZoneProtectionProfileFloodSynObject) AttributeTypes() map[string]attr.Type {
	var ratesObj *ZoneProtectionProfileFloodRatesObject
	return map[string]attr.Type{
		"enable":      types.BoolType,
		"red":         types.ObjectType{AttrTypes: ratesObj.AttributeTypes()},
		"syn_cookies": types.ObjectType{AttrTypes: ratesObj.AttributeTypes()},
	}
}

func (o *ZoneProtectionProfileFloodObject) AttributeTypes() map[string]attr.Type {
	var synObj *ZoneProtectionProfileFloodSynObject
	var protocolObj *ZoneProtectionProfileFloodProtocolObject
	return map[string]attr.Type{
		"syn":    types.ObjectType{AttrTypes: synObj.AttributeTypes()},
		"icmp":   types.ObjectType{AttrTypes: protocolObj.AttributeTypes()},
		"icmpv6": types.ObjectType{AttrTypes: protocolObj.AttributeTypes()},
		"udp":    types.ObjectType{AttrTypes: protocolObj.AttributeTypes()},
		"other":  types.ObjectType{AttrTypes: protocolObj.AttributeTypes()},
	}
}

// -----------------------------------------------------------------------
// Flood conversion helpers
// -----------------------------------------------------------------------

func copyFloodToPango(ctx context.Context, floodVal types.Object) (*zone_protection.Flood, diag.Diagnostics) {
	var diags diag.Diagnostics
	if floodVal.IsNull() || floodVal.IsUnknown() {
		return nil, diags
	}

	var floodObj ZoneProtectionProfileFloodObject
	diags.Append(floodVal.As(ctx, &floodObj, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil, diags
	}

	flood := &zone_protection.Flood{}

	// SYN
	if !floodObj.Syn.IsNull() && !floodObj.Syn.IsUnknown() {
		var synObj ZoneProtectionProfileFloodSynObject
		diags.Append(floodObj.Syn.As(ctx, &synObj, basetypes.ObjectAsOptions{})...)
		if diags.HasError() {
			return nil, diags
		}
		flood.Syn = &zone_protection.FloodSyn{
			Enable: synObj.Enable.ValueBoolPointer(),
		}
		if !synObj.Red.IsNull() && !synObj.Red.IsUnknown() {
			var ratesObj ZoneProtectionProfileFloodRatesObject
			diags.Append(synObj.Red.As(ctx, &ratesObj, basetypes.ObjectAsOptions{})...)
			if diags.HasError() {
				return nil, diags
			}
			flood.Syn.Red = &zone_protection.FloodRates{
				AlarmRate:    ratesObj.AlarmRate.ValueInt64Pointer(),
				ActivateRate: ratesObj.ActivateRate.ValueInt64Pointer(),
				MaximalRate:  ratesObj.MaximalRate.ValueInt64Pointer(),
			}
		}
		if !synObj.SynCookies.IsNull() && !synObj.SynCookies.IsUnknown() {
			var ratesObj ZoneProtectionProfileFloodRatesObject
			diags.Append(synObj.SynCookies.As(ctx, &ratesObj, basetypes.ObjectAsOptions{})...)
			if diags.HasError() {
				return nil, diags
			}
			flood.Syn.SynCookies = &zone_protection.FloodRates{
				AlarmRate:    ratesObj.AlarmRate.ValueInt64Pointer(),
				ActivateRate: ratesObj.ActivateRate.ValueInt64Pointer(),
				MaximalRate:  ratesObj.MaximalRate.ValueInt64Pointer(),
			}
		}
	}

	// ICMP, ICMPv6, UDP, Other
	if p, d := copyFloodProtocolToPango(ctx, floodObj.Icmp); d.HasError() {
		diags.Append(d...)
		return nil, diags
	} else {
		flood.Icmp = p
	}
	if p, d := copyFloodProtocolToPango(ctx, floodObj.Icmpv6); d.HasError() {
		diags.Append(d...)
		return nil, diags
	} else {
		flood.Icmpv6 = p
	}
	if p, d := copyFloodProtocolToPango(ctx, floodObj.Udp); d.HasError() {
		diags.Append(d...)
		return nil, diags
	} else {
		flood.Udp = p
	}
	if p, d := copyFloodProtocolToPango(ctx, floodObj.Other); d.HasError() {
		diags.Append(d...)
		return nil, diags
	} else {
		flood.Other = p
	}

	return flood, diags
}

func copyFloodProtocolToPango(ctx context.Context, obj types.Object) (*zone_protection.FloodProtocol, diag.Diagnostics) {
	var diags diag.Diagnostics
	if obj.IsNull() || obj.IsUnknown() {
		return nil, diags
	}

	var protocolObj ZoneProtectionProfileFloodProtocolObject
	diags.Append(obj.As(ctx, &protocolObj, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil, diags
	}

	p := &zone_protection.FloodProtocol{
		Enable: protocolObj.Enable.ValueBoolPointer(),
	}
	if !protocolObj.Red.IsNull() && !protocolObj.Red.IsUnknown() {
		var ratesObj ZoneProtectionProfileFloodRatesObject
		diags.Append(protocolObj.Red.As(ctx, &ratesObj, basetypes.ObjectAsOptions{})...)
		if diags.HasError() {
			return nil, diags
		}
		p.Red = &zone_protection.FloodRates{
			AlarmRate:    ratesObj.AlarmRate.ValueInt64Pointer(),
			ActivateRate: ratesObj.ActivateRate.ValueInt64Pointer(),
			MaximalRate:  ratesObj.MaximalRate.ValueInt64Pointer(),
		}
	}
	return p, diags
}

func copyFloodFromPango(ctx context.Context, flood *zone_protection.Flood, existing types.Object) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	var floodObj *ZoneProtectionProfileFloodObject

	if flood == nil {
		return types.ObjectNull(floodObj.AttributeTypes()), diags
	}

	// SYN
	var synVal types.Object
	{
		var synObj *ZoneProtectionProfileFloodSynObject
		if flood.Syn == nil {
			synVal = types.ObjectNull(synObj.AttributeTypes())
		} else {
			redVal, d := floodRatesFromPango(ctx, flood.Syn.Red)
			diags.Append(d...)
			if diags.HasError() {
				return types.ObjectNull(floodObj.AttributeTypes()), diags
			}
			synCookiesVal, d := floodRatesFromPango(ctx, flood.Syn.SynCookies)
			diags.Append(d...)
			if diags.HasError() {
				return types.ObjectNull(floodObj.AttributeTypes()), diags
			}
			synTF := ZoneProtectionProfileFloodSynObject{
				Enable:     boolPtrToType(flood.Syn.Enable),
				Red:        redVal,
				SynCookies: synCookiesVal,
			}
			var sv types.Object
			sv, d = types.ObjectValueFrom(ctx, synTF.AttributeTypes(), synTF)
			diags.Append(d...)
			if diags.HasError() {
				return types.ObjectNull(floodObj.AttributeTypes()), diags
			}
			synVal = sv
		}
	}

	icmpVal, d := floodProtocolFromPango(ctx, flood.Icmp)
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(floodObj.AttributeTypes()), diags
	}
	icmpv6Val, d := floodProtocolFromPango(ctx, flood.Icmpv6)
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(floodObj.AttributeTypes()), diags
	}
	udpVal, d := floodProtocolFromPango(ctx, flood.Udp)
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(floodObj.AttributeTypes()), diags
	}
	otherVal, d := floodProtocolFromPango(ctx, flood.Other)
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(floodObj.AttributeTypes()), diags
	}

	floodTF := ZoneProtectionProfileFloodObject{
		Syn:    synVal,
		Icmp:   icmpVal,
		Icmpv6: icmpv6Val,
		Udp:    udpVal,
		Other:  otherVal,
	}
	result, d := types.ObjectValueFrom(ctx, floodTF.AttributeTypes(), floodTF)
	diags.Append(d...)
	return result, diags
}

func floodProtocolFromPango(ctx context.Context, p *zone_protection.FloodProtocol) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	var protocolObj *ZoneProtectionProfileFloodProtocolObject
	if p == nil {
		return types.ObjectNull(protocolObj.AttributeTypes()), diags
	}

	redVal, d := floodRatesFromPango(ctx, p.Red)
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(protocolObj.AttributeTypes()), diags
	}

	tf := ZoneProtectionProfileFloodProtocolObject{
		Enable: boolPtrToType(p.Enable),
		Red:    redVal,
	}
	result, d := types.ObjectValueFrom(ctx, tf.AttributeTypes(), tf)
	diags.Append(d...)
	return result, diags
}

func floodRatesFromPango(ctx context.Context, r *zone_protection.FloodRates) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	var ratesObj *ZoneProtectionProfileFloodRatesObject
	if r == nil {
		return types.ObjectNull(ratesObj.AttributeTypes()), diags
	}

	tf := ZoneProtectionProfileFloodRatesObject{
		AlarmRate:    int64PtrToType(r.AlarmRate),
		ActivateRate: int64PtrToType(r.ActivateRate),
		MaximalRate:  int64PtrToType(r.MaximalRate),
	}
	result, d := types.ObjectValueFrom(ctx, tf.AttributeTypes(), tf)
	diags.Append(d...)
	return result, diags
}

// -----------------------------------------------------------------------
// Scan object types
// -----------------------------------------------------------------------

// ZoneProtectionProfileScanBlockIpObject maps to pango ScanBlockIp.
type ZoneProtectionProfileScanBlockIpObject struct {
	TrackBy  types.String `tfsdk:"track_by"`
	Duration types.Int64  `tfsdk:"duration"`
}

func (o *ZoneProtectionProfileScanBlockIpObject) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"track_by": types.StringType,
		"duration": types.Int64Type,
	}
}

// ZoneProtectionProfileScanObject maps to pango ScanEntry.
// Name is the scan type, e.g. "tcp-port-scan", "udp-port-scan", "host-sweep".
// Action is mutually exclusive: exactly one of alert, block, block_ip should be set.
type ZoneProtectionProfileScanObject struct {
	Name      types.String `tfsdk:"name"`
	Interval  types.Int64  `tfsdk:"interval"`
	Threshold types.Int64  `tfsdk:"threshold"`
	Alert     types.Bool   `tfsdk:"alert"`
	Block     types.Bool   `tfsdk:"block"`
	BlockIp   types.Object `tfsdk:"block_ip"`
}

func (o *ZoneProtectionProfileScanObject) AttributeTypes() map[string]attr.Type {
	var blockIpObj *ZoneProtectionProfileScanBlockIpObject
	return map[string]attr.Type{
		"name":      types.StringType,
		"interval":  types.Int64Type,
		"threshold": types.Int64Type,
		"alert":     types.BoolType,
		"block":     types.BoolType,
		"block_ip":  types.ObjectType{AttrTypes: blockIpObj.AttributeTypes()},
	}
}

func ZoneProtectionProfileDataSourceScanSchema() map[string]dsschema.Attribute {
	var blockIpObj *ZoneProtectionProfileScanBlockIpObject
	return map[string]dsschema.Attribute{
		"name": dsschema.StringAttribute{
			Description: "Scan type: tcp-port-scan, udp-port-scan, or host-sweep.",
			Required:    true,
		},
		"interval": dsschema.Int64Attribute{
			Description: "Interval in seconds.",
			Optional:    true,
			Computed:    true,
		},
		"threshold": dsschema.Int64Attribute{
			Description: "Threshold (number of scan attempts).",
			Optional:    true,
			Computed:    true,
		},
		"alert": dsschema.BoolAttribute{
			Description: "Action: alert only.",
			Optional:    true,
			Computed:    true,
		},
		"block": dsschema.BoolAttribute{
			Description: "Action: block.",
			Optional:    true,
			Computed:    true,
		},
		"block_ip": dsschema.SingleNestedAttribute{
			Description: "Action: block IP.",
			Optional:    true,
			Computed:    true,
			Attributes: func() map[string]dsschema.Attribute {
				_ = blockIpObj
				return map[string]dsschema.Attribute{
					"track_by": dsschema.StringAttribute{
						Description: "Track by attacker or attacker-and-victim.",
						Optional:    true,
						Computed:    true,
					},
					"duration": dsschema.Int64Attribute{
						Description: "Block duration in seconds.",
						Optional:    true,
						Computed:    true,
					},
				}
			}(),
		},
	}
}

func ZoneProtectionProfileResourceScanSchema() map[string]rsschema.Attribute {
	return map[string]rsschema.Attribute{
		"name": rsschema.StringAttribute{
			Description: "Scan type: tcp-port-scan, udp-port-scan, or host-sweep.",
			Required:    true,
		},
		"interval": rsschema.Int64Attribute{
			Description: "Interval in seconds.",
			Optional:    true,
		},
		"threshold": rsschema.Int64Attribute{
			Description: "Threshold (number of scan attempts).",
			Optional:    true,
		},
		"alert": rsschema.BoolAttribute{
			Description: "Action: alert only.",
			Optional:    true,
		},
		"block": rsschema.BoolAttribute{
			Description: "Action: block.",
			Optional:    true,
		},
		"block_ip": rsschema.SingleNestedAttribute{
			Description: "Action: block IP.",
			Optional:    true,
			Attributes: map[string]rsschema.Attribute{
				"track_by": rsschema.StringAttribute{
					Description: "Track by attacker or attacker-and-victim.",
					Optional:    true,
				},
				"duration": rsschema.Int64Attribute{
					Description: "Block duration in seconds.",
					Optional:    true,
				},
			},
		},
	}
}

// -----------------------------------------------------------------------
// Scan conversion helpers
// -----------------------------------------------------------------------

func copyScanToPango(ctx context.Context, scanList types.List) ([]zone_protection.ScanEntry, diag.Diagnostics) {
	var diags diag.Diagnostics
	if scanList.IsNull() || scanList.IsUnknown() {
		return nil, diags
	}

	var scanObjs []ZoneProtectionProfileScanObject
	diags.Append(scanList.ElementsAs(ctx, &scanObjs, false)...)
	if diags.HasError() {
		return nil, diags
	}

	entries := make([]zone_protection.ScanEntry, 0, len(scanObjs))
	for _, s := range scanObjs {
		entry := zone_protection.ScanEntry{
			Name:      s.Name.ValueString(),
			Interval:  s.Interval.ValueInt64Pointer(),
			Threshold: s.Threshold.ValueInt64Pointer(),
		}

		if s.Alert.ValueBool() {
			entry.Action.Alert = true
		} else if s.Block.ValueBool() {
			entry.Action.Block = true
		} else if !s.BlockIp.IsNull() && !s.BlockIp.IsUnknown() {
			var blockIpObj ZoneProtectionProfileScanBlockIpObject
			diags.Append(s.BlockIp.As(ctx, &blockIpObj, basetypes.ObjectAsOptions{})...)
			if diags.HasError() {
				return nil, diags
			}
			entry.Action.BlockIp = &zone_protection.ScanBlockIp{
				TrackBy:  blockIpObj.TrackBy.ValueStringPointer(),
				Duration: blockIpObj.Duration.ValueInt64Pointer(),
			}
		}
		entries = append(entries, entry)
	}
	return entries, diags
}

func copyScanFromPango(ctx context.Context, entries []zone_protection.ScanEntry) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	var scanObj *ZoneProtectionProfileScanObject
	elemType := types.ObjectType{AttrTypes: scanObj.AttributeTypes()}

	if len(entries) == 0 {
		return types.ListNull(elemType), diags
	}

	tfEntries := make([]ZoneProtectionProfileScanObject, 0, len(entries))
	for _, e := range entries {
		var blockIpVal types.Object
		var blockIpObjT *ZoneProtectionProfileScanBlockIpObject
		if e.Action.BlockIp != nil {
			bip := ZoneProtectionProfileScanBlockIpObject{
				TrackBy:  types.StringPointerValue(e.Action.BlockIp.TrackBy),
				Duration: int64PtrToType(e.Action.BlockIp.Duration),
			}
			var d diag.Diagnostics
			blockIpVal, d = types.ObjectValueFrom(ctx, bip.AttributeTypes(), bip)
			diags.Append(d...)
			if diags.HasError() {
				return types.ListNull(elemType), diags
			}
		} else {
			blockIpVal = types.ObjectNull(blockIpObjT.AttributeTypes())
		}

		tfEntries = append(tfEntries, ZoneProtectionProfileScanObject{
			Name:      types.StringValue(e.Name),
			Interval:  int64PtrToType(e.Interval),
			Threshold: int64PtrToType(e.Threshold),
			Alert:     types.BoolValue(e.Action.Alert),
			Block:     types.BoolValue(e.Action.Block),
			BlockIp:   blockIpVal,
		})
	}

	result, d := types.ListValueFrom(ctx, elemType, tfEntries)
	diags.Append(d...)
	return result, diags
}

// -----------------------------------------------------------------------
// Small helpers
// -----------------------------------------------------------------------

func setBoolFromPtr(dst *types.Bool, src *bool) {
	if src == nil {
		*dst = types.BoolNull()
	} else {
		*dst = types.BoolValue(*src)
	}
}

func boolPtrToType(v *bool) types.Bool {
	if v == nil {
		return types.BoolNull()
	}
	return types.BoolValue(*v)
}

func int64PtrToType(v *int64) types.Int64 {
	if v == nil {
		return types.Int64Null()
	}
	return types.Int64Value(*v)
}
