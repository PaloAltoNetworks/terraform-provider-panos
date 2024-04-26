package provider

import (
	"context"
	"fmt"
	//"regexp"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/policies/rules/security"
	"github.com/PaloAltoNetworks/pango/rule"
	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/xmlapi"

	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	//"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	rsschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	//"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Data source (listing).
var (
	_ datasource.DataSource              = &SecurityRuleListDataSource{}
	_ datasource.DataSourceWithConfigure = &SecurityRuleListDataSource{}
)

func NewSecurityRuleListDataSource() datasource.DataSource {
	return &SecurityRuleListDataSource{}
}

type SecurityRuleListDataSource struct {
	client *pango.XmlApiClient
}

type SecurityRuleDsListModel struct {
	// Input.
	Filter   *SecurityRuleDsListFilter  `tfsdk:"filter"`
	Location SecurityRuleDsListLocation `tfsdk:"location"`

	// Output.
	Data []SecurityRuleDsListEntry `tfsdk:"data"`
}

type SecurityRuleDsListFilter struct {
	Config types.String `tfsdk:"config"`
	Value  types.String `tfsdk:"value"`
	Quote  types.String `tfsdk:"quote"`
}

type SecurityRuleDsListLocation struct {
	Shared      *SecurityRuleDsListSharedLocation      `tfsdk:"shared"`
	Vsys        *SecurityRuleDsListVsysLocation        `tfsdk:"vsys"`
	DeviceGroup *SecurityRuleDsListDeviceGroupLocation `tfsdk:"device_group"`
}

type SecurityRuleDsListSharedLocation struct {
	Rulebase types.String `tfsdk:"rulebase"`
}

type SecurityRuleDsListVsysLocation struct {
	NgfwDevice types.String `tfsdk:"ngfw_device"`
	Name       types.String `tfsdk:"name"`
}

type SecurityRuleDsListDeviceGroupLocation struct {
	PanoramaDevice types.String `tfsdk:"panorama_device"`
	Name           types.String `tfsdk:"name"`
	Rulebase       types.String `tfsdk:"rulebase"`
}

type SecurityRuleDsListEntry struct {
	Name                            types.String                             `tfsdk:"name"`
	Uuid                            types.String                             `tfsdk:"uuid"`
	SourceZones                     types.Set                                `tfsdk:"source_zones"`
	DestinationZones                types.Set                                `tfsdk:"destination_zones"`
	SourceAddresses                 types.Set                                `tfsdk:"source_addresses"`
	SourceUsers                     types.Set                                `tfsdk:"source_users"`
	DestinationAddresses            types.Set                                `tfsdk:"destination_addresses"`
	Services                        types.Set                                `tfsdk:"services"`
	Categories                      types.Set                                `tfsdk:"categories"`
	Applications                    types.Set                                `tfsdk:"applications"`
	SourceDevices                   types.Set                                `tfsdk:"source_devices"`
	DestinationDevices              types.Set                                `tfsdk:"destination_devices"`
	Schedule                        types.String                             `tfsdk:"schedule"`
	Tags                            types.List                               `tfsdk:"tags"`
	NegateSource                    types.Bool                               `tfsdk:"negate_source"`
	NegateDestination               types.Bool                               `tfsdk:"negate_destination"`
	Disabled                        types.Bool                               `tfsdk:"disabled"`
	Description                     types.String                             `tfsdk:"description"`
	GroupTag                        types.String                             `tfsdk:"group_tag"`
	Action                          types.String                             `tfsdk:"action"`
	IcmpUnreachable                 types.Bool                               `tfsdk:"icmp_unreachable"`
	Type                            types.String                             `tfsdk:"type"`
	DisableServerResponseInspection types.Bool                               `tfsdk:"disable_server_response_inspection"`
	LogSetting                      types.String                             `tfsdk:"log_setting"`
	LogStart                        types.Bool                               `tfsdk:"log_start"`
	LogEnd                          types.Bool                               `tfsdk:"log_end"`
	ProfileSettings                 *SecurityRuleDsListProfileSettingsObject `tfsdk:"profile_settings"`
	Qos                             *SecurityRuleDsListQosObject             `tfsdk:"qos"`
	// TODO: figure out Targets
	NegateTarget   types.Bool `tfsdk:"negate_target"`
	DisableInspect types.Bool `tfsdk:"disable_inspect"`
}

type SecurityRuleDsListProfileSettingsObject struct {
	Groups   types.List                        `tfsdk:"groups"`
	Profiles *SecurityRuleDsListProfilesObject `tfsdk:"profiles"`
}

type SecurityRuleDsListProfilesObject struct {
	UrlFilteringProfiles     types.List `tfsdk:"url_filtering_profiles"`
	DataFilteringProfiles    types.List `tfsdk:"data_filtering_profiles"`
	FileBlockingProfiles     types.List `tfsdk:"file_blocking_profiles"`
	WildfireAnalysisProfiles types.List `tfsdk:"wildfire_analysis_profiles"`
	AntiVirusProfiles        types.List `tfsdk:"anti_virus_profiles"`
	AntiSpywareProfiles      types.List `tfsdk:"anti_spyware_profiles"`
	VulnerabilityProfiles    types.List `tfsdk:"vulnerability_profiles"`
}

type SecurityRuleDsListQosObject struct {
	IpDscp                   types.String `tfsdk:"ip_dscp"`
	IpPrecedence             types.String `tfsdk:"ip_precedence"`
	FollowClientToServerFlow types.Bool   `tfsdk:"follow_client_to_server_flow"`
}

func (d *SecurityRuleListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_security_rule_list"
}

func (d *SecurityRuleListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Description: "Returns a list of security rules.",
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
							"rulebase": dsschema.StringAttribute{
								Description: "The rulebase. Valid values are `pre-rulebase` or `post-rulebase`.",
								Required:    true,
							},
						},
					},
					"shared": dsschema.SingleNestedAttribute{
						Description: "(Panorama) Located in shared.",
						Optional:    true,
						Attributes: map[string]dsschema.Attribute{
							"rulebase": dsschema.StringAttribute{
								Description: "The rulebase. Valid values are `pre-rulebase` or `post-rulebase`.",
								Required:    true,
							},
						},
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
			"data": dsschema.ListNestedAttribute{
				Description: "The list of objects.",
				Computed:    true,
				NestedObject: dsschema.NestedAttributeObject{
					Attributes: map[string]dsschema.Attribute{
						"name": dsschema.StringAttribute{
							Description: "Alphanumeric string [ 0-9a-zA-Z._-].",
							Computed:    true,
						},
						"uuid": dsschema.StringAttribute{
							Description: "The UUID.",
							Computed:    true,
						},
						"source_zones": dsschema.SetAttribute{
							Description: "The source zones.",
							Computed:    true,
							ElementType: types.StringType,
						},
						"destination_zones": dsschema.SetAttribute{
							Description: "The destination zones.",
							Computed:    true,
							ElementType: types.StringType,
						},
						"source_addresses": dsschema.SetAttribute{
							Description: "The source addresses.",
							Computed:    true,
							ElementType: types.StringType,
						},
						"source_users": dsschema.SetAttribute{
							Description: "The source users.",
							Computed:    true,
							ElementType: types.StringType,
						},
						"destination_addresses": dsschema.SetAttribute{
							Description: "The destination addresses.",
							Computed:    true,
							ElementType: types.StringType,
						},
						"services": dsschema.SetAttribute{
							Description: "The services.",
							Computed:    true,
							ElementType: types.StringType,
						},
						"categories": dsschema.SetAttribute{
							Description: "The categories.",
							Computed:    true,
							ElementType: types.StringType,
						},
						"applications": dsschema.SetAttribute{
							Description: "The applications.",
							Computed:    true,
							ElementType: types.StringType,
						},
						"source_devices": dsschema.SetAttribute{
							Description: "Source HIP devices.",
							Computed:    true,
							ElementType: types.StringType,
						},
						"destination_devices": dsschema.SetAttribute{
							Description: "Destination HIP devices.",
							Computed:    true,
							ElementType: types.StringType,
						},
						"schedule": dsschema.StringAttribute{
							Description: "Schedule.",
							Computed:    true,
						},
						"tags": dsschema.ListAttribute{
							Description: "Tags for address object.",
							Computed:    true,
							ElementType: types.StringType,
						},
						"negate_source": dsschema.BoolAttribute{
							Description: "Negate the source.",
							Computed:    true,
						},
						"negate_destination": dsschema.BoolAttribute{
							Description: "Negate the destination.",
							Computed:    true,
						},
						"disabled": dsschema.BoolAttribute{
							Description: "If the rule is disabled or not.",
							Computed:    true,
						},
						"description": dsschema.StringAttribute{
							Description: "The description.",
							Computed:    true,
						},
						"group_tag": dsschema.StringAttribute{
							Description: "The group.",
							Computed:    true,
						},
						"action": dsschema.StringAttribute{
							Description: "The rule action.",
							Computed:    true,
						},
						"icmp_unreachable": dsschema.BoolAttribute{
							Description: "ICMP unreachable.",
							Computed:    true,
						},
						"type": dsschema.StringAttribute{
							Description: "Rule type.",
							Computed:    true,
						},
						"disable_server_response_inspection": dsschema.BoolAttribute{
							Description: "Disable server response inspection.",
							Computed:    true,
						},
						"log_setting": dsschema.StringAttribute{
							Description: "Log setting.",
							Computed:    true,
						},
						"log_start": dsschema.BoolAttribute{
							Description: "Log at session start.",
							Computed:    true,
						},
						"log_end": dsschema.BoolAttribute{
							Description: "Log at session end.",
							Computed:    true,
						},
						"profile_settings": dsschema.SingleNestedAttribute{
							Description: "The profile settings.",
							Computed:    true,
							Attributes: map[string]dsschema.Attribute{
								"groups": dsschema.ListAttribute{
									Description: "The groups.",
									Computed:    true,
									ElementType: types.StringType,
								},
								"profiles": dsschema.SingleNestedAttribute{
									Description: "Profiles.",
									Computed:    true,
									Attributes: map[string]dsschema.Attribute{
										"url_filtering_profiles": dsschema.ListAttribute{
											Description: "URL filtering profiles.",
											Computed:    true,
											ElementType: types.StringType,
										},
										"data_filtering_profiles": dsschema.ListAttribute{
											Description: "Data filtering profiles.",
											Computed:    true,
											ElementType: types.StringType,
										},
										"file_blocking_profiles": dsschema.ListAttribute{
											Description: "File blocking profiles.",
											Computed:    true,
											ElementType: types.StringType,
										},
										"wildfire_analysis_profiles": dsschema.ListAttribute{
											Description: "Wildfire analysis profiles.",
											Computed:    true,
											ElementType: types.StringType,
										},
										"anti_virus_profiles": dsschema.ListAttribute{
											Description: "Anti-virus profiles.",
											Computed:    true,
											ElementType: types.StringType,
										},
										"anti_spyware_profiles": dsschema.ListAttribute{
											Description: "Anti-spyware profiles.",
											Computed:    true,
											ElementType: types.StringType,
										},
										"vulnerability_profiles": dsschema.ListAttribute{
											Description: "Vulnerability profiles.",
											Computed:    true,
											ElementType: types.StringType,
										},
									},
								},
							},
						},
						"qos": dsschema.SingleNestedAttribute{
							Description: "The QOS settings.",
							Computed:    true,
							Attributes: map[string]dsschema.Attribute{
								"ip_dscp": dsschema.StringAttribute{
									Description: "IP DSCP.",
									Computed:    true,
								},
								"ip_precedence": dsschema.StringAttribute{
									Description: "IP precedence.",
									Computed:    true,
								},
								"follow_client_to_server_flow": dsschema.BoolAttribute{
									Description: "Follow client to server flow.",
									Computed:    true,
								},
							},
						},
						// TODO: targets schema
						"negate_target": dsschema.BoolAttribute{
							Description: "Negate the target.",
							Computed:    true,
						},
						"disable_inspect": dsschema.BoolAttribute{
							Description: "(PAN-OS 10.2+) Disable inspect.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *SecurityRuleListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*pango.XmlApiClient)
}

func (d *SecurityRuleListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state SecurityRuleDsListModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Determine the location.
	var loc security.Location
	if state.Location.Shared != nil {
		loc.Shared = &security.SharedLocation{}

		// Rulebase.
		if state.Location.Shared.Rulebase.IsNull() {
			loc.Shared.Rulebase = "pre-rulebase"
		} else {
			loc.Shared.Rulebase = state.Location.Shared.Rulebase.ValueString()
		}
	} else if state.Location.Vsys != nil {
		loc.Vsys = &security.VsysLocation{}

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
		loc.DeviceGroup = &security.DeviceGroupLocation{}

		// PanoramaDevice.
		if state.Location.DeviceGroup.PanoramaDevice.IsNull() {
			loc.DeviceGroup.PanoramaDevice = "localhost.localdomain"
		} else {
			loc.DeviceGroup.PanoramaDevice = state.Location.DeviceGroup.PanoramaDevice.ValueString()
		}

		// Rulebase.
		if state.Location.DeviceGroup.Rulebase.IsNull() {
			loc.DeviceGroup.Rulebase = "pre-rulebase"
		} else {
			loc.DeviceGroup.Rulebase = state.Location.DeviceGroup.Rulebase.ValueString()
		}

		// Name.
		if state.Location.DeviceGroup.Name.IsNull() {
			resp.Diagnostics.AddError("Invalid location", "The device group name must be specified.")
			return
		}
		loc.DeviceGroup.Name = state.Location.DeviceGroup.Name.ValueString()
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
	svc := security.NewService(d.client)

	var err error
	var list []security.Entry

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
		state.Data = make([]SecurityRuleDsListEntry, 0, len(list))
		for _, var0 := range list {
			var1 := SecurityRuleDsListEntry{}
			var1.Name = types.StringValue(var0.Name)
			var1.Uuid = types.StringValue(var0.Uuid)
			var2, var3 := types.SetValueFrom(ctx, types.StringType, var0.SourceZones)
			var1.SourceZones = var2
			resp.Diagnostics.Append(var3.Errors()...)
			var4, var5 := types.SetValueFrom(ctx, types.StringType, var0.DestinationZones)
			var1.DestinationZones = var4
			resp.Diagnostics.Append(var5.Errors()...)
			var6, var7 := types.SetValueFrom(ctx, types.StringType, var0.SourceAddresses)
			var1.SourceAddresses = var6
			resp.Diagnostics.Append(var7.Errors()...)
			var8, var9 := types.SetValueFrom(ctx, types.StringType, var0.SourceUsers)
			var1.SourceUsers = var8
			resp.Diagnostics.Append(var9.Errors()...)
			var10, var11 := types.SetValueFrom(ctx, types.StringType, var0.DestinationAddresses)
			var1.DestinationAddresses = var10
			resp.Diagnostics.Append(var11.Errors()...)
			var12, var13 := types.SetValueFrom(ctx, types.StringType, var0.Services)
			var1.Services = var12
			resp.Diagnostics.Append(var13.Errors()...)
			var14, var15 := types.SetValueFrom(ctx, types.StringType, var0.Categories)
			var1.Categories = var14
			resp.Diagnostics.Append(var15.Errors()...)
			var16, var17 := types.SetValueFrom(ctx, types.StringType, var0.Applications)
			var1.Applications = var16
			resp.Diagnostics.Append(var17.Errors()...)
			var18, var19 := types.SetValueFrom(ctx, types.StringType, var0.SourceDevices)
			var1.SourceDevices = var18
			resp.Diagnostics.Append(var19.Errors()...)
			var20, var21 := types.SetValueFrom(ctx, types.StringType, var0.DestinationDevices)
			var1.DestinationDevices = var20
			resp.Diagnostics.Append(var21.Errors()...)
			var1.Schedule = types.StringPointerValue(var0.Schedule)
			var22, var23 := types.ListValueFrom(ctx, types.StringType, var0.Tags)
			var1.Tags = var22
			resp.Diagnostics.Append(var23.Errors()...)
			var1.NegateSource = types.BoolPointerValue(var0.NegateSource)
			var1.NegateDestination = types.BoolPointerValue(var0.NegateDestination)
			var1.Disabled = types.BoolPointerValue(var0.Disabled)
			var1.Description = types.StringPointerValue(var0.Description)
			var1.GroupTag = types.StringPointerValue(var0.GroupTag)
			var1.Action = types.StringValue(var0.Action)
			var1.IcmpUnreachable = types.BoolPointerValue(var0.IcmpUnreachable)
			var1.Type = types.StringPointerValue(var0.Type)
			var1.DisableServerResponseInspection = types.BoolPointerValue(var0.DisableServerResponseInspection)
			var1.LogSetting = types.StringPointerValue(var0.LogSetting)
			var1.LogStart = types.BoolPointerValue(var0.LogStart)
			var1.LogEnd = types.BoolPointerValue(var0.LogEnd)
			if var0.ProfileSettings != nil {
				var1.ProfileSettings = &SecurityRuleDsListProfileSettingsObject{}
				var24, var25 := types.ListValueFrom(ctx, types.StringType, var0.ProfileSettings.Groups)
				var1.ProfileSettings.Groups = var24
				resp.Diagnostics.Append(var25.Errors()...)
				if var0.ProfileSettings.Profiles != nil {
					var1.ProfileSettings.Profiles = &SecurityRuleDsListProfilesObject{}
					var26, var27 := types.ListValueFrom(ctx, types.StringType, var0.ProfileSettings.Profiles.UrlFilteringProfiles)
					var1.ProfileSettings.Profiles.UrlFilteringProfiles = var26
					resp.Diagnostics.Append(var27.Errors()...)
					var28, var29 := types.ListValueFrom(ctx, types.StringType, var0.ProfileSettings.Profiles.DataFilteringProfiles)
					var1.ProfileSettings.Profiles.DataFilteringProfiles = var28
					resp.Diagnostics.Append(var29.Errors()...)
					var30, var31 := types.ListValueFrom(ctx, types.StringType, var0.ProfileSettings.Profiles.FileBlockingProfiles)
					var1.ProfileSettings.Profiles.FileBlockingProfiles = var30
					resp.Diagnostics.Append(var31.Errors()...)
					var32, var33 := types.ListValueFrom(ctx, types.StringType, var0.ProfileSettings.Profiles.WildfireAnalysisProfiles)
					var1.ProfileSettings.Profiles.WildfireAnalysisProfiles = var32
					resp.Diagnostics.Append(var33.Errors()...)
					var34, var35 := types.ListValueFrom(ctx, types.StringType, var0.ProfileSettings.Profiles.AntiVirusProfiles)
					var1.ProfileSettings.Profiles.AntiVirusProfiles = var34
					resp.Diagnostics.Append(var35.Errors()...)
					var36, var37 := types.ListValueFrom(ctx, types.StringType, var0.ProfileSettings.Profiles.AntiSpywareProfiles)
					var1.ProfileSettings.Profiles.AntiSpywareProfiles = var36
					resp.Diagnostics.Append(var37.Errors()...)
					var38, var39 := types.ListValueFrom(ctx, types.StringType, var0.ProfileSettings.Profiles.VulnerabilityProfiles)
					var1.ProfileSettings.Profiles.VulnerabilityProfiles = var38
					resp.Diagnostics.Append(var39.Errors()...)
				}
			}
			if var0.Qos != nil {
				var1.Qos = &SecurityRuleDsListQosObject{}
				var1.Qos.IpDscp = types.StringPointerValue(var0.Qos.IpDscp)
				var1.Qos.IpPrecedence = types.StringPointerValue(var0.Qos.IpPrecedence)
				var1.Qos.FollowClientToServerFlow = types.BoolValue(var0.Qos.FollowClientToServerFlow != nil)
			}
			// TODO: targets
			var1.NegateTarget = types.BoolPointerValue(var0.NegateTarget)
			var1.DisableInspect = types.BoolPointerValue(var0.DisableInspect)

			state.Data = append(state.Data, var1)
		}
	}

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Data source.
var (
	_ datasource.DataSource              = &SecurityRuleDataSource{}
	_ datasource.DataSourceWithConfigure = &SecurityRuleDataSource{}
)

func NewSecurityRuleDataSource() datasource.DataSource {
	return &SecurityRuleDataSource{}
}

type SecurityRuleDataSource struct {
	client *pango.XmlApiClient
}

type SecurityRuleDsModel struct {
	// Input.
	Filter   *SecurityRuleDsFilter  `tfsdk:"filter"`
	Location SecurityRuleDsLocation `tfsdk:"location"`
	Name     types.String           `tfsdk:"name"`
	Uuid     types.String           `tfsdk:"uuid"`

	// Output.
	SourceZones                     types.Set                            `tfsdk:"source_zones"`
	DestinationZones                types.Set                            `tfsdk:"destination_zones"`
	SourceAddresses                 types.Set                            `tfsdk:"source_addresses"`
	SourceUsers                     types.Set                            `tfsdk:"source_users"`
	DestinationAddresses            types.Set                            `tfsdk:"destination_addresses"`
	Services                        types.Set                            `tfsdk:"services"`
	Categories                      types.Set                            `tfsdk:"categories"`
	Applications                    types.Set                            `tfsdk:"applications"`
	SourceDevices                   types.Set                            `tfsdk:"source_devices"`
	DestinationDevices              types.Set                            `tfsdk:"destination_devices"`
	Schedule                        types.String                         `tfsdk:"schedule"`
	Tags                            types.List                           `tfsdk:"tags"`
	NegateSource                    types.Bool                           `tfsdk:"negate_source"`
	NegateDestination               types.Bool                           `tfsdk:"negate_destination"`
	Disabled                        types.Bool                           `tfsdk:"disabled"`
	Description                     types.String                         `tfsdk:"description"`
	GroupTag                        types.String                         `tfsdk:"group_tag"`
	Action                          types.String                         `tfsdk:"action"`
	IcmpUnreachable                 types.Bool                           `tfsdk:"icmp_unreachable"`
	Type                            types.String                         `tfsdk:"type"`
	DisableServerResponseInspection types.Bool                           `tfsdk:"disable_server_response_inspection"`
	LogSetting                      types.String                         `tfsdk:"log_setting"`
	LogStart                        types.Bool                           `tfsdk:"log_start"`
	LogEnd                          types.Bool                           `tfsdk:"log_end"`
	ProfileSettings                 *SecurityRuleDsProfileSettingsObject `tfsdk:"profile_settings"`
	Qos                             *SecurityRuleDsQosObject             `tfsdk:"qos"`
	// TODO: figure out Targets
	NegateTarget   types.Bool `tfsdk:"negate_target"`
	DisableInspect types.Bool `tfsdk:"disable_inspect"`
}

type SecurityRuleDsFilter struct {
	Config types.String `tfsdk:"config"`
}

type SecurityRuleDsLocation struct {
	Shared      *SecurityRuleDsSharedLocation      `tfsdk:"shared"`
	Vsys        *SecurityRuleDsVsysLocation        `tfsdk:"vsys"`
	DeviceGroup *SecurityRuleDsDeviceGroupLocation `tfsdk:"device_group"`
}

type SecurityRuleDsSharedLocation struct {
	Rulebase types.String `tfsdk:"rulebase"`
}

type SecurityRuleDsVsysLocation struct {
	NgfwDevice types.String `tfsdk:"ngfw_device"`
	Name       types.String `tfsdk:"name"`
}

type SecurityRuleDsDeviceGroupLocation struct {
	PanoramaDevice types.String `tfsdk:"panorama_device"`
	Name           types.String `tfsdk:"name"`
	Rulebase       types.String `tfsdk:"rulebase"`
}

type SecurityRuleDsProfileSettingsObject struct {
	Groups   types.List                    `tfsdk:"groups"`
	Profiles *SecurityRuleDsProfilesObject `tfsdk:"profiles"`
}

type SecurityRuleDsProfilesObject struct {
	UrlFilteringProfiles     types.List `tfsdk:"url_filtering_profiles"`
	DataFilteringProfiles    types.List `tfsdk:"data_filtering_profiles"`
	FileBlockingProfiles     types.List `tfsdk:"file_blocking_profiles"`
	WildfireAnalysisProfiles types.List `tfsdk:"wildfire_analysis_profiles"`
	AntiVirusProfiles        types.List `tfsdk:"anti_virus_profiles"`
	AntiSpywareProfiles      types.List `tfsdk:"anti_spyware_profiles"`
	VulnerabilityProfiles    types.List `tfsdk:"vulnerability_profiles"`
}

type SecurityRuleDsQosObject struct {
	IpDscp                   types.String `tfsdk:"ip_dscp"`
	IpPrecedence             types.String `tfsdk:"ip_precedence"`
	FollowClientToServerFlow types.Bool   `tfsdk:"follow_client_to_server_flow"`
}

func (d *SecurityRuleDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_security_rule"
}

func (d *SecurityRuleDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Description: "Returns information about the given security rule.",
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
							"rulebase": dsschema.StringAttribute{
								Description: "The rulebase. Valid values are `pre-rulebase` or `post-rulebase`.",
								Required:    true,
							},
						},
					},
					"shared": dsschema.SingleNestedAttribute{
						Description: "(Panorama) Located in shared.",
						Optional:    true,
						Attributes: map[string]dsschema.Attribute{
							"rulebase": dsschema.StringAttribute{
								Description: "The rulebase. Valid values are `pre-rulebase` or `post-rulebase`.",
								Required:    true,
							},
						},
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
				Description: "Alphanumeric string [ 0-9a-zA-Z._-]. Either name or uuid must be specified.",
				Optional:    true,
				Computed:    true,
			},
			"uuid": dsschema.StringAttribute{
				Description: "The UUID. Either name or uuid must be specified.",
				Optional:    true,
				Computed:    true,
			},

			// Output.
			"source_zones": dsschema.SetAttribute{
				Description: "The source zones.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"destination_zones": dsschema.SetAttribute{
				Description: "The destination zones.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"source_addresses": dsschema.SetAttribute{
				Description: "The source addresses.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"source_users": dsschema.SetAttribute{
				Description: "The source users.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"destination_addresses": dsschema.SetAttribute{
				Description: "The destination addresses.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"services": dsschema.SetAttribute{
				Description: "The services.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"categories": dsschema.SetAttribute{
				Description: "The categories.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"applications": dsschema.SetAttribute{
				Description: "The applications.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"source_devices": dsschema.SetAttribute{
				Description: "Source HIP devices.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"destination_devices": dsschema.SetAttribute{
				Description: "Destination HIP devices.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"schedule": dsschema.StringAttribute{
				Description: "Schedule.",
				Computed:    true,
			},
			"tags": dsschema.ListAttribute{
				Description: "Tags for address object.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"negate_source": dsschema.BoolAttribute{
				Description: "Negate the source.",
				Computed:    true,
			},
			"negate_destination": dsschema.BoolAttribute{
				Description: "Negate the destination.",
				Computed:    true,
			},
			"disabled": dsschema.BoolAttribute{
				Description: "If the rule is disabled or not.",
				Computed:    true,
			},
			"description": dsschema.StringAttribute{
				Description: "The description.",
				Computed:    true,
			},
			"group_tag": dsschema.StringAttribute{
				Description: "The group.",
				Computed:    true,
			},
			"action": dsschema.StringAttribute{
				Description: "The rule action.",
				Computed:    true,
			},
			"icmp_unreachable": dsschema.BoolAttribute{
				Description: "ICMP unreachable.",
				Computed:    true,
			},
			"type": dsschema.StringAttribute{
				Description: "Rule type.",
				Computed:    true,
			},
			"disable_server_response_inspection": dsschema.BoolAttribute{
				Description: "Disable server response inspection.",
				Computed:    true,
			},
			"log_setting": dsschema.StringAttribute{
				Description: "Log setting.",
				Computed:    true,
			},
			"log_start": dsschema.BoolAttribute{
				Description: "Log at session start.",
				Computed:    true,
			},
			"log_end": dsschema.BoolAttribute{
				Description: "Log at session end.",
				Computed:    true,
			},
			"profile_settings": dsschema.SingleNestedAttribute{
				Description: "The profile settings.",
				Computed:    true,
				Attributes: map[string]dsschema.Attribute{
					"groups": dsschema.ListAttribute{
						Description: "The groups.",
						Computed:    true,
						ElementType: types.StringType,
					},
					"profiles": dsschema.SingleNestedAttribute{
						Description: "Profiles.",
						Computed:    true,
						Attributes: map[string]dsschema.Attribute{
							"url_filtering_profiles": dsschema.ListAttribute{
								Description: "URL filtering profiles.",
								Computed:    true,
								ElementType: types.StringType,
							},
							"data_filtering_profiles": dsschema.ListAttribute{
								Description: "Data filtering profiles.",
								Computed:    true,
								ElementType: types.StringType,
							},
							"file_blocking_profiles": dsschema.ListAttribute{
								Description: "File blocking profiles.",
								Computed:    true,
								ElementType: types.StringType,
							},
							"wildfire_analysis_profiles": dsschema.ListAttribute{
								Description: "Wildfire analysis profiles.",
								Computed:    true,
								ElementType: types.StringType,
							},
							"anti_virus_profiles": dsschema.ListAttribute{
								Description: "Anti-virus profiles.",
								Computed:    true,
								ElementType: types.StringType,
							},
							"anti_spyware_profiles": dsschema.ListAttribute{
								Description: "Anti-spyware profiles.",
								Computed:    true,
								ElementType: types.StringType,
							},
							"vulnerability_profiles": dsschema.ListAttribute{
								Description: "Vulnerability profiles.",
								Computed:    true,
								ElementType: types.StringType,
							},
						},
					},
				},
			},
			"qos": dsschema.SingleNestedAttribute{
				Description: "The QOS settings.",
				Computed:    true,
				Attributes: map[string]dsschema.Attribute{
					"ip_dscp": dsschema.StringAttribute{
						Description: "IP DSCP.",
						Computed:    true,
					},
					"ip_precedence": dsschema.StringAttribute{
						Description: "IP precedence.",
						Computed:    true,
					},
					"follow_client_to_server_flow": dsschema.BoolAttribute{
						Description: "Follow client to server flow.",
						Computed:    true,
					},
				},
			},
			// TODO: targets schema
			"negate_target": dsschema.BoolAttribute{
				Description: "Negate the target.",
				Computed:    true,
			},
			"disable_inspect": dsschema.BoolAttribute{
				Description: "(PAN-OS 10.2+) Disable inspect.",
				Computed:    true,
			},
		},
	}
}

func (d *SecurityRuleDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*pango.XmlApiClient)
}

func (d *SecurityRuleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state SecurityRuleDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Determine the location.
	var loc security.Location
	if state.Location.Shared != nil {
		loc.Shared = &security.SharedLocation{}

		// Rulebase.
		if state.Location.Shared.Rulebase.IsNull() {
			loc.Shared.Rulebase = "pre-rulebase"
		} else {
			loc.Shared.Rulebase = state.Location.Shared.Rulebase.ValueString()
		}
	} else if state.Location.Vsys != nil {
		loc.Vsys = &security.VsysLocation{}

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
		loc.DeviceGroup = &security.DeviceGroupLocation{}

		// PanoramaDevice.
		if state.Location.DeviceGroup.PanoramaDevice.IsNull() {
			loc.DeviceGroup.PanoramaDevice = "localhost.localdomain"
		} else {
			loc.DeviceGroup.PanoramaDevice = state.Location.DeviceGroup.PanoramaDevice.ValueString()
		}

		// Rulebase.
		if state.Location.DeviceGroup.Rulebase.IsNull() {
			loc.DeviceGroup.Rulebase = "pre-rulebase"
		} else {
			loc.DeviceGroup.Rulebase = state.Location.DeviceGroup.Rulebase.ValueString()
		}

		// Name.
		if state.Location.DeviceGroup.Name.IsNull() {
			resp.Diagnostics.AddError("Invalid location", "The device group name must be specified.")
			return
		}
		loc.DeviceGroup.Name = state.Location.DeviceGroup.Name.ValueString()
	} else {
		resp.Diagnostics.AddError("Unknown location", "Location for object is unknown")
		return
	}

	// Determine the rest of the Read params.
	var action string
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
	}

	// Either name or uuid must be given.
	var name, uuid string
	if !state.Name.IsNull() && !state.Uuid.IsNull() {
		resp.Diagnostics.AddError("Invalid input", `Specify either "name" or "uuid" not both`)
		return
	} else if !state.Name.IsNull() {
		name = state.Name.ValueString()
	} else if !state.Uuid.IsNull() {
		uuid = state.Uuid.ValueString()
	} else {
		resp.Diagnostics.AddError("Invalid input", `Specify either "name" or "uuid"`)
		return
	}

	// Create the service.
	svc := security.NewService(d.client)

	var err error
	var ans *security.Entry

	// Perform the operation.
	if d.client.Hostname != "" {
		if name != "" {
			ans, err = svc.Read(ctx, loc, name, action)
		} else {
			ans, err = svc.ReadById(ctx, loc, uuid, action)
		}
	} else {
		if name != "" {
			ans, err = svc.ReadFromConfig(ctx, loc, name)
		} else {
			ans, err = svc.ReadFromConfigById(ctx, loc, uuid)
		}
	}

	if err != nil {
		resp.Diagnostics.AddError("Error in read", err.Error())
		return
	}

	// Save the information to state.
	state.Name = types.StringValue(ans.Name)
	state.Uuid = types.StringValue(ans.Uuid)
	var2, var3 := types.SetValueFrom(ctx, types.StringType, ans.SourceZones)
	state.SourceZones = var2
	resp.Diagnostics.Append(var3.Errors()...)
	var4, var5 := types.SetValueFrom(ctx, types.StringType, ans.DestinationZones)
	state.DestinationZones = var4
	resp.Diagnostics.Append(var5.Errors()...)
	var6, var7 := types.SetValueFrom(ctx, types.StringType, ans.SourceAddresses)
	state.SourceAddresses = var6
	resp.Diagnostics.Append(var7.Errors()...)
	var8, var9 := types.SetValueFrom(ctx, types.StringType, ans.SourceUsers)
	state.SourceUsers = var8
	resp.Diagnostics.Append(var9.Errors()...)
	var10, var11 := types.SetValueFrom(ctx, types.StringType, ans.DestinationAddresses)
	state.DestinationAddresses = var10
	resp.Diagnostics.Append(var11.Errors()...)
	var12, var13 := types.SetValueFrom(ctx, types.StringType, ans.Services)
	state.Services = var12
	resp.Diagnostics.Append(var13.Errors()...)
	var14, var15 := types.SetValueFrom(ctx, types.StringType, ans.Categories)
	state.Categories = var14
	resp.Diagnostics.Append(var15.Errors()...)
	var16, var17 := types.SetValueFrom(ctx, types.StringType, ans.Applications)
	state.Applications = var16
	resp.Diagnostics.Append(var17.Errors()...)
	var18, var19 := types.SetValueFrom(ctx, types.StringType, ans.SourceDevices)
	state.SourceDevices = var18
	resp.Diagnostics.Append(var19.Errors()...)
	var20, var21 := types.SetValueFrom(ctx, types.StringType, ans.DestinationDevices)
	state.DestinationDevices = var20
	resp.Diagnostics.Append(var21.Errors()...)
	state.Schedule = types.StringPointerValue(ans.Schedule)
	var22, var23 := types.ListValueFrom(ctx, types.StringType, ans.Tags)
	state.Tags = var22
	resp.Diagnostics.Append(var23.Errors()...)
	state.NegateSource = types.BoolPointerValue(ans.NegateSource)
	state.NegateDestination = types.BoolPointerValue(ans.NegateDestination)
	state.Disabled = types.BoolPointerValue(ans.Disabled)
	state.Description = types.StringPointerValue(ans.Description)
	state.GroupTag = types.StringPointerValue(ans.GroupTag)
	state.Action = types.StringValue(ans.Action)
	state.IcmpUnreachable = types.BoolPointerValue(ans.IcmpUnreachable)
	state.Type = types.StringPointerValue(ans.Type)
	state.DisableServerResponseInspection = types.BoolPointerValue(ans.DisableServerResponseInspection)
	state.LogSetting = types.StringPointerValue(ans.LogSetting)
	state.LogStart = types.BoolPointerValue(ans.LogStart)
	state.LogEnd = types.BoolPointerValue(ans.LogEnd)
	if ans.ProfileSettings != nil {
		state.ProfileSettings = &SecurityRuleDsProfileSettingsObject{}
		var24, var25 := types.ListValueFrom(ctx, types.StringType, ans.ProfileSettings.Groups)
		state.ProfileSettings.Groups = var24
		resp.Diagnostics.Append(var25.Errors()...)
		if ans.ProfileSettings.Profiles != nil {
			state.ProfileSettings.Profiles = &SecurityRuleDsProfilesObject{}
			var26, var27 := types.ListValueFrom(ctx, types.StringType, ans.ProfileSettings.Profiles.UrlFilteringProfiles)
			state.ProfileSettings.Profiles.UrlFilteringProfiles = var26
			resp.Diagnostics.Append(var27.Errors()...)
			var28, var29 := types.ListValueFrom(ctx, types.StringType, ans.ProfileSettings.Profiles.DataFilteringProfiles)
			state.ProfileSettings.Profiles.DataFilteringProfiles = var28
			resp.Diagnostics.Append(var29.Errors()...)
			var30, var31 := types.ListValueFrom(ctx, types.StringType, ans.ProfileSettings.Profiles.FileBlockingProfiles)
			state.ProfileSettings.Profiles.FileBlockingProfiles = var30
			resp.Diagnostics.Append(var31.Errors()...)
			var32, var33 := types.ListValueFrom(ctx, types.StringType, ans.ProfileSettings.Profiles.WildfireAnalysisProfiles)
			state.ProfileSettings.Profiles.WildfireAnalysisProfiles = var32
			resp.Diagnostics.Append(var33.Errors()...)
			var34, var35 := types.ListValueFrom(ctx, types.StringType, ans.ProfileSettings.Profiles.AntiVirusProfiles)
			state.ProfileSettings.Profiles.AntiVirusProfiles = var34
			resp.Diagnostics.Append(var35.Errors()...)
			var36, var37 := types.ListValueFrom(ctx, types.StringType, ans.ProfileSettings.Profiles.AntiSpywareProfiles)
			state.ProfileSettings.Profiles.AntiSpywareProfiles = var36
			resp.Diagnostics.Append(var37.Errors()...)
			var38, var39 := types.ListValueFrom(ctx, types.StringType, ans.ProfileSettings.Profiles.VulnerabilityProfiles)
			state.ProfileSettings.Profiles.VulnerabilityProfiles = var38
			resp.Diagnostics.Append(var39.Errors()...)
		}
	}
	if ans.Qos != nil {
		state.Qos = &SecurityRuleDsQosObject{}
		state.Qos.IpDscp = types.StringPointerValue(ans.Qos.IpDscp)
		state.Qos.IpPrecedence = types.StringPointerValue(ans.Qos.IpPrecedence)
		state.Qos.FollowClientToServerFlow = types.BoolValue(ans.Qos.FollowClientToServerFlow != nil)
	}
	// TODO: targets
	state.NegateTarget = types.BoolPointerValue(ans.NegateTarget)
	state.DisableInspect = types.BoolPointerValue(ans.DisableInspect)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Resource.
var (
	_ resource.Resource                = &SecurityPolicyRulesResource{}
	_ resource.ResourceWithConfigure   = &SecurityPolicyRulesResource{}
	_ resource.ResourceWithImportState = &SecurityPolicyRulesResource{}
)

func NewSecurityPolicyRulesResource() resource.Resource {
	return &SecurityPolicyRulesResource{}
}

type SecurityPolicyRulesResource struct {
	client *pango.XmlApiClient
}

type SecurityPolicyRulesTfid struct {
	Rules    []RuleInfo        `json:"rules"`
	Position rule.Position     `json:"position"`
	Location security.Location `json:"location"`
}

func (o *SecurityPolicyRulesTfid) IsValid() error {
	var err error
	if len(o.Rules) == 0 {
		return fmt.Errorf("No rules present")
	}

	names := make(map[string]bool)
	uuids := make(map[string]bool)
	for _, rInfo := range o.Rules {
		if names[rInfo.Name] {
			return fmt.Errorf("Multiple rules have name %q", rInfo.Name)
		}
		names[rInfo.Name] = true

		if uuids[rInfo.Uuid] {
			return fmt.Errorf("Multiple rules have uuid %q", rInfo.Uuid)
		}
		uuids[rInfo.Uuid] = true
	}

	if err = o.Location.IsValid(); err != nil {
		return err
	}

	if err = o.Position.IsValid(false); err != nil {
		return err
	}

	return nil
}

type SecurityPolicyRulesResourceModel struct {
	//Timeouts crudTimeouts `tfsdk:"timeouts"`
	Tfid     types.String                        `tfsdk:"tfid"`
	Location SecurityPolicyRulesResourceLocation `tfsdk:"location"`
	Position SecurityPolicyRulesResourcePosition `tfsdk:"position"`
	Rules    []SecurityPolicyRulesResourceEntry  `tfsdk:"rules"`
	Ordered  types.Bool                          `tfsdk:"ordered"`
}

type SecurityPolicyRulesResourceLocation struct {
	Shared      *SecurityPolicyRulesResourceSharedLocation      `tfsdk:"shared"`
	Vsys        *SecurityPolicyRulesResourceVsysLocation        `tfsdk:"vsys"`
	DeviceGroup *SecurityPolicyRulesResourceDeviceGroupLocation `tfsdk:"device_group"`
}

type SecurityPolicyRulesResourceSharedLocation struct {
	Rulebase types.String `tfsdk:"rulebase"`
}

type SecurityPolicyRulesResourceVsysLocation struct {
	NgfwDevice types.String `tfsdk:"ngfw_device"`
	Name       types.String `tfsdk:"name"`
}

type SecurityPolicyRulesResourceDeviceGroupLocation struct {
	PanoramaDevice types.String `tfsdk:"panorama_device"`
	Name           types.String `tfsdk:"name"`
	Rulebase       types.String `tfsdk:"rulebase"`
}

type SecurityPolicyRulesResourcePosition struct {
	Ok              types.Bool   `tfsdk:"ok"`
	First           types.Bool   `tfsdk:"first"`
	Last            types.Bool   `tfsdk:"last"`
	SomewhereBefore types.String `tfsdk:"somewhere_before"`
	DirectlyBefore  types.String `tfsdk:"directly_before"`
	SomewhereAfter  types.String `tfsdk:"somewhere_after"`
	DirectlyAfter   types.String `tfsdk:"directly_after"`
}

type SecurityPolicyRulesResourceEntry struct {
	Name                            types.String                                      `tfsdk:"name"`
	Uuid                            types.String                                      `tfsdk:"uuid"`
	SourceZones                     types.Set                                         `tfsdk:"source_zones"`
	DestinationZones                types.Set                                         `tfsdk:"destination_zones"`
	SourceAddresses                 types.Set                                         `tfsdk:"source_addresses"`
	SourceUsers                     types.Set                                         `tfsdk:"source_users"`
	DestinationAddresses            types.Set                                         `tfsdk:"destination_addresses"`
	Services                        types.Set                                         `tfsdk:"services"`
	Categories                      types.Set                                         `tfsdk:"categories"`
	Applications                    types.Set                                         `tfsdk:"applications"`
	SourceDevices                   types.Set                                         `tfsdk:"source_devices"`
	DestinationDevices              types.Set                                         `tfsdk:"destination_devices"`
	Schedule                        types.String                                      `tfsdk:"schedule"`
	Tags                            types.List                                        `tfsdk:"tags"`
	NegateSource                    types.Bool                                        `tfsdk:"negate_source"`
	NegateDestination               types.Bool                                        `tfsdk:"negate_destination"`
	Disabled                        types.Bool                                        `tfsdk:"disabled"`
	Description                     types.String                                      `tfsdk:"description"`
	GroupTag                        types.String                                      `tfsdk:"group_tag"`
	Action                          types.String                                      `tfsdk:"action"`
	IcmpUnreachable                 types.Bool                                        `tfsdk:"icmp_unreachable"`
	Type                            types.String                                      `tfsdk:"type"`
	DisableServerResponseInspection types.Bool                                        `tfsdk:"disable_server_response_inspection"`
	LogSetting                      types.String                                      `tfsdk:"log_setting"`
	LogStart                        types.Bool                                        `tfsdk:"log_start"`
	LogEnd                          types.Bool                                        `tfsdk:"log_end"`
	ProfileSettings                 *SecurityPolicyRulesResourceProfileSettingsObject `tfsdk:"profile_settings"`
	Qos                             *SecurityPolicyRulesResourceQosObject             `tfsdk:"qos"`
	// TODO: figure out Targets
	NegateTarget   types.Bool `tfsdk:"negate_target"`
	DisableInspect types.Bool `tfsdk:"disable_inspect"`
}

type SecurityPolicyRulesResourceProfileSettingsObject struct {
	Groups   types.List                                 `tfsdk:"groups"`
	Profiles *SecurityPolicyRulesResourceProfilesObject `tfsdk:"profiles"`
}

type SecurityPolicyRulesResourceProfilesObject struct {
	UrlFilteringProfiles     types.List `tfsdk:"url_filtering_profiles"`
	DataFilteringProfiles    types.List `tfsdk:"data_filtering_profiles"`
	FileBlockingProfiles     types.List `tfsdk:"file_blocking_profiles"`
	WildfireAnalysisProfiles types.List `tfsdk:"wildfire_analysis_profiles"`
	AntiVirusProfiles        types.List `tfsdk:"anti_virus_profiles"`
	AntiSpywareProfiles      types.List `tfsdk:"anti_spyware_profiles"`
	VulnerabilityProfiles    types.List `tfsdk:"vulnerability_profiles"`
}

type SecurityPolicyRulesResourceQosObject struct {
	IpDscp                   types.String `tfsdk:"ip_dscp"`
	IpPrecedence             types.String `tfsdk:"ip_precedence"`
	FollowClientToServerFlow types.Bool   `tfsdk:"follow_client_to_server_flow"`
}

func (r *SecurityPolicyRulesResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_security_policy_rules"
}

func (r *SecurityPolicyRulesResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = rsschema.Schema{
		Description: "Manges a group of security rules in the specified order.",

		Attributes: map[string]rsschema.Attribute{
			"tfid": rsschema.StringAttribute{
				Description: "The Terraform ID.",
				Computed:    true,
			},
			"location": rsschema.SingleNestedAttribute{
				Description: "The location of this object. One and only one of the locations should be specified.",
				Required:    true,
				Attributes: map[string]rsschema.Attribute{
					"device_group": rsschema.SingleNestedAttribute{
						Description: "(Panorama) The given device group.",
						Optional:    true,
						Validators: []validator.Object{
							objectvalidator.ExactlyOneOf(
								path.MatchRoot("location").AtName("shared"),
								path.MatchRoot("location").AtName("vsys"),
							),
						},
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
							"rulebase": rsschema.StringAttribute{
								Description: "The rulebase. Valid values are `pre-rulebase` or `post-rulebase`. Default: `pre-rulebase`",
								Optional:    true,
								Computed:    true,
								Default:     stringdefault.StaticString("pre-rulebase"),
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.RequiresReplace(),
								},
								Validators: []validator.String{
									stringvalidator.OneOf("pre-rulebase", "post-rulebase"),
								},
							},
						},
					},
					"shared": rsschema.SingleNestedAttribute{
						Description: "(Panorama) Located in shared.",
						Optional:    true,
						Attributes: map[string]rsschema.Attribute{
							"rulebase": rsschema.StringAttribute{
								Description: "The rulebase. Valid values are `pre-rulebase` or `post-rulebase`. Default: `pre-rulebase`",
								Optional:    true,
								Computed:    true,
								Default:     stringdefault.StaticString("pre-rulebase"),
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.RequiresReplace(),
								},
								Validators: []validator.String{
									stringvalidator.OneOf("pre-rulebase", "post-rulebase"),
								},
							},
						},
					},
					"vsys": rsschema.SingleNestedAttribute{
						Description: "(NGFW) The given vsys.",
						Optional:    true,
						Attributes: map[string]rsschema.Attribute{
							"name": rsschema.StringAttribute{
								Description: "The vsys name. Default: `vsys1`.",
								Optional:    true,
								Computed:    true,
								Default:     stringdefault.StaticString("vsys1"),
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.RequiresReplace(),
								},
							},
							"ngfw_device": rsschema.StringAttribute{
								Description: "The NGFW device. Default: `localhost.localdomain`.",
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
			"ordered": rsschema.BoolAttribute{
				Description: "(Internal use) Validates that the rule ordering is ok. If the rules are unordered, then this will be set to `false` when state is read.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"position": rsschema.SingleNestedAttribute{
				Description: "The location of the rule group. If none of the positions are defined, then the rules will only be sequentially placed, not in any particular location in the rulebase as a whole.",
				Required:    true,
				Attributes: map[string]rsschema.Attribute{
					"ok": rsschema.BoolAttribute{
						Description: "(Internal use) Validates that the position is ok. If the position is not ok, this will be set to `false` when state is read.",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(true),
					},
					"first": rsschema.BoolAttribute{
						Description: "Locate the rule group at the top of the policy.",
						Optional:    true,
						Validators: []validator.Bool{
							boolvalidator.ConflictsWith(
								path.MatchRoot("position").AtName("last"),
								path.MatchRoot("position").AtName("somewhere_before"),
								path.MatchRoot("position").AtName("somewhere_after"),
								path.MatchRoot("position").AtName("directly_before"),
								path.MatchRoot("position").AtName("directly_after"),
							),
						},
					},
					"last": rsschema.BoolAttribute{
						Description: "Locate the rule group at the bottom of the policy.",
						Optional:    true,
					},
					"somewhere_before": rsschema.StringAttribute{
						Description: "Locate the rule group somewhere before the given rule name.",
						Optional:    true,
					},
					"somewhere_after": rsschema.StringAttribute{
						Description: "Locate the rule group somewhere after the given rule name.",
						Optional:    true,
					},
					"directly_before": rsschema.StringAttribute{
						Description: "Locate the rule group directly before the given rule name.",
						Optional:    true,
					},
					"directly_after": rsschema.StringAttribute{
						Description: "Locate the rule group directly after the given rule name.",
						Optional:    true,
					},
				},
			},
			"rules": rsschema.ListNestedAttribute{
				Description: "The list of security rules.",
				Required:    true,
				NestedObject: rsschema.NestedAttributeObject{
					Attributes: map[string]rsschema.Attribute{
						"name": rsschema.StringAttribute{
							Description: "Alphanumeric string [ 0-9a-zA-Z._-].",
							Required:    true,
						},
						"uuid": rsschema.StringAttribute{
							Description: "The UUID.",
							Computed:    true,
						},
						"source_zones": rsschema.SetAttribute{
							Description: "The source zones.",
							Optional:    true,
							ElementType: types.StringType,
						},
						"destination_zones": rsschema.SetAttribute{
							Description: "The destination zones.",
							Optional:    true,
							ElementType: types.StringType,
						},
						"source_addresses": rsschema.SetAttribute{
							Description: "The source addresses.",
							Optional:    true,
							ElementType: types.StringType,
						},
						"source_users": rsschema.SetAttribute{
							Description: "The source users.",
							Optional:    true,
							ElementType: types.StringType,
						},
						"destination_addresses": rsschema.SetAttribute{
							Description: "The destination addresses.",
							Optional:    true,
							ElementType: types.StringType,
						},
						"services": rsschema.SetAttribute{
							Description: "The services.",
							Optional:    true,
							ElementType: types.StringType,
						},
						"categories": rsschema.SetAttribute{
							Description: "The categories.",
							Optional:    true,
							ElementType: types.StringType,
						},
						"applications": rsschema.SetAttribute{
							Description: "The applications.",
							Optional:    true,
							ElementType: types.StringType,
						},
						"source_devices": rsschema.SetAttribute{
							Description: "Source HIP devices.",
							Optional:    true,
							ElementType: types.StringType,
						},
						"destination_devices": rsschema.SetAttribute{
							Description: "Destination HIP devices.",
							Optional:    true,
							ElementType: types.StringType,
						},
						"schedule": rsschema.StringAttribute{
							Description: "Schedule.",
							Optional:    true,
						},
						"tags": rsschema.ListAttribute{
							Description: "Tags for address object.",
							Optional:    true,
							ElementType: types.StringType,
						},
						"negate_source": rsschema.BoolAttribute{
							Description: "Negate the source.",
							Optional:    true,
						},
						"negate_destination": rsschema.BoolAttribute{
							Description: "Negate the destination.",
							Optional:    true,
						},
						"disabled": rsschema.BoolAttribute{
							Description: "If the rule is disabled or not.",
							Optional:    true,
						},
						"description": rsschema.StringAttribute{
							Description: "The description.",
							Optional:    true,
						},
						"group_tag": rsschema.StringAttribute{
							Description: "The group.",
							Optional:    true,
						},
						"action": rsschema.StringAttribute{
							Description: "The rule action.",
							Required:    true,
						},
						"icmp_unreachable": rsschema.BoolAttribute{
							Description: "ICMP unreachable.",
							Optional:    true,
						},
						"type": rsschema.StringAttribute{
							Description: "Rule type.",
							Optional:    true,
						},
						"disable_server_response_inspection": rsschema.BoolAttribute{
							Description: "Disable server response inspection.",
							Optional:    true,
						},
						"log_setting": rsschema.StringAttribute{
							Description: "Log setting.",
							Optional:    true,
						},
						"log_start": rsschema.BoolAttribute{
							Description: "Log at session start.",
							Optional:    true,
						},
						"log_end": rsschema.BoolAttribute{
							Description: "Log at session end.",
							Optional:    true,
							Computed:    true,
							Default:     booldefault.StaticBool(true),
						},
						"profile_settings": rsschema.SingleNestedAttribute{
							Description: "The profile settings.",
							Optional:    true,
							Attributes: map[string]rsschema.Attribute{
								"groups": rsschema.ListAttribute{
									Description: "The groups.",
									Optional:    true,
									ElementType: types.StringType,
								},
								"profiles": rsschema.SingleNestedAttribute{
									Description: "Profiles.",
									Optional:    true,
									Attributes: map[string]rsschema.Attribute{
										"url_filtering_profiles": rsschema.ListAttribute{
											Description: "URL filtering profiles.",
											Optional:    true,
											ElementType: types.StringType,
										},
										"data_filtering_profiles": rsschema.ListAttribute{
											Description: "Data filtering profiles.",
											Optional:    true,
											ElementType: types.StringType,
										},
										"file_blocking_profiles": rsschema.ListAttribute{
											Description: "File blocking profiles.",
											Optional:    true,
											ElementType: types.StringType,
										},
										"wildfire_analysis_profiles": rsschema.ListAttribute{
											Description: "Wildfire analysis profiles.",
											Optional:    true,
											ElementType: types.StringType,
										},
										"anti_virus_profiles": rsschema.ListAttribute{
											Description: "Anti-virus profiles.",
											Optional:    true,
											ElementType: types.StringType,
										},
										"anti_spyware_profiles": rsschema.ListAttribute{
											Description: "Anti-spyware profiles.",
											Optional:    true,
											ElementType: types.StringType,
										},
										"vulnerability_profiles": rsschema.ListAttribute{
											Description: "Vulnerability profiles.",
											Optional:    true,
											ElementType: types.StringType,
										},
									},
								},
							},
						},
						"qos": rsschema.SingleNestedAttribute{
							Description: "The QOS settings.",
							Optional:    true,
							Attributes: map[string]rsschema.Attribute{
								"ip_dscp": rsschema.StringAttribute{
									Description: "IP DSCP.",
									Optional:    true,
								},
								"ip_precedence": rsschema.StringAttribute{
									Description: "IP precedence.",
									Optional:    true,
								},
								"follow_client_to_server_flow": rsschema.BoolAttribute{
									Description: "Follow client to server flow.",
									Optional:    true,
									Computed:    true,
									Default:     booldefault.StaticBool(false),
								},
							},
						},
						// TODO: targets schema
						"negate_target": rsschema.BoolAttribute{
							Description: "Negate the target.",
							Optional:    true,
						},
						"disable_inspect": rsschema.BoolAttribute{
							Description: "(PAN-OS 10.2+) Disable inspect.",
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func (r *SecurityPolicyRulesResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*pango.XmlApiClient)
}

func (r *SecurityPolicyRulesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var err error
	var listing []security.Entry
	var tfid SecurityPolicyRulesTfid
	var state SecurityPolicyRulesResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource create", map[string]any{
		"resource_name": "panos_security_policy_rules",
		"function":      "Create",
	})

	// Verify mode.
	if r.client.Hostname == "" {
		resp.Diagnostics.AddError("Invalid mode error", InspectionModeError)
		return
	}

	// Create the service.
	svc := security.NewService(r.client)

	// Determine the location.
	if state.Location.Shared != nil {
		tfid.Location.Shared = &security.SharedLocation{
			Rulebase: state.Location.Shared.Rulebase.ValueString(),
		}
	}
	if state.Location.Vsys != nil {
		tfid.Location.Vsys = &security.VsysLocation{
			NgfwDevice: state.Location.Vsys.NgfwDevice.ValueString(),
			Name:       state.Location.Vsys.Name.ValueString(),
		}
	}
	if state.Location.DeviceGroup != nil {
		tfid.Location.DeviceGroup = &security.DeviceGroupLocation{
			PanoramaDevice: state.Location.DeviceGroup.PanoramaDevice.ValueString(),
			Rulebase:       state.Location.DeviceGroup.Rulebase.ValueString(),
			Name:           state.Location.DeviceGroup.Name.ValueString(),
		}
	}
	if err = tfid.Location.IsValid(); err != nil {
		resp.Diagnostics.AddError("Invalid path", err.Error())
		return
	}

	// Determine the position.
	tfid.Position = rule.Position{
		First:           state.Position.First.ValueBoolPointer(),
		Last:            state.Position.Last.ValueBoolPointer(),
		SomewhereBefore: state.Position.SomewhereBefore.ValueStringPointer(),
		SomewhereAfter:  state.Position.SomewhereAfter.ValueStringPointer(),
		DirectlyBefore:  state.Position.DirectlyBefore.ValueStringPointer(),
		DirectlyAfter:   state.Position.DirectlyAfter.ValueStringPointer(),
	}
	if err = tfid.Position.IsValid(false); err != nil {
		resp.Diagnostics.AddError("Invalid position", err.Error())
		return
	}

	// Load the desired config.
	isNewName := make(map[string]bool)
	entryNames := make([]string, 0, len(state.Rules))
	entries := make([]security.Entry, 0, len(state.Rules))
	for _, var0 := range state.Rules {
		isNewName[var0.Name.ValueString()] = true
		entryNames = append(entryNames, var0.Name.ValueString())
		var1 := security.Entry{Name: var0.Name.ValueString()}
		resp.Diagnostics.Append(var0.SourceZones.ElementsAs(ctx, &var1.SourceZones, false)...)
		resp.Diagnostics.Append(var0.DestinationZones.ElementsAs(ctx, &var1.DestinationZones, false)...)
		resp.Diagnostics.Append(var0.SourceAddresses.ElementsAs(ctx, &var1.SourceAddresses, false)...)
		resp.Diagnostics.Append(var0.SourceUsers.ElementsAs(ctx, &var1.SourceUsers, false)...)
		resp.Diagnostics.Append(var0.DestinationAddresses.ElementsAs(ctx, &var1.DestinationAddresses, false)...)
		resp.Diagnostics.Append(var0.Services.ElementsAs(ctx, &var1.Services, false)...)
		resp.Diagnostics.Append(var0.Categories.ElementsAs(ctx, &var1.Categories, false)...)
		resp.Diagnostics.Append(var0.Applications.ElementsAs(ctx, &var1.Applications, false)...)
		resp.Diagnostics.Append(var0.SourceDevices.ElementsAs(ctx, &var1.SourceDevices, false)...)
		resp.Diagnostics.Append(var0.DestinationDevices.ElementsAs(ctx, &var1.DestinationDevices, false)...)
		var1.Schedule = var0.Schedule.ValueStringPointer()
		resp.Diagnostics.Append(var0.Tags.ElementsAs(ctx, &var1.Tags, false)...)
		var1.NegateSource = var0.NegateSource.ValueBoolPointer()
		var1.NegateDestination = var0.NegateDestination.ValueBoolPointer()
		var1.Disabled = var0.Disabled.ValueBoolPointer()
		var1.Description = var0.Description.ValueStringPointer()
		var1.GroupTag = var0.GroupTag.ValueStringPointer()
		var1.Action = var0.Action.ValueString()
		var1.IcmpUnreachable = var0.IcmpUnreachable.ValueBoolPointer()
		var1.Type = var0.Type.ValueStringPointer()
		var1.DisableServerResponseInspection = var0.DisableServerResponseInspection.ValueBoolPointer()
		var1.LogSetting = var0.LogSetting.ValueStringPointer()
		var1.LogStart = var0.LogStart.ValueBoolPointer()
		var1.LogEnd = var0.LogEnd.ValueBoolPointer()
		if var0.ProfileSettings != nil {
			var1.ProfileSettings = &security.ProfileSettingsObject{}
			resp.Diagnostics.Append(var0.ProfileSettings.Groups.ElementsAs(ctx, &var1.ProfileSettings.Groups, false)...)
			if var0.ProfileSettings.Profiles != nil {
				var1.ProfileSettings.Profiles = &security.ProfilesObject{}
				resp.Diagnostics.Append(var0.ProfileSettings.Profiles.UrlFilteringProfiles.ElementsAs(ctx, &var1.ProfileSettings.Profiles.UrlFilteringProfiles, false)...)
				resp.Diagnostics.Append(var0.ProfileSettings.Profiles.DataFilteringProfiles.ElementsAs(ctx, &var1.ProfileSettings.Profiles.DataFilteringProfiles, false)...)
				resp.Diagnostics.Append(var0.ProfileSettings.Profiles.FileBlockingProfiles.ElementsAs(ctx, &var1.ProfileSettings.Profiles.FileBlockingProfiles, false)...)
				resp.Diagnostics.Append(var0.ProfileSettings.Profiles.WildfireAnalysisProfiles.ElementsAs(ctx, &var1.ProfileSettings.Profiles.WildfireAnalysisProfiles, false)...)
				resp.Diagnostics.Append(var0.ProfileSettings.Profiles.AntiVirusProfiles.ElementsAs(ctx, &var1.ProfileSettings.Profiles.AntiVirusProfiles, false)...)
				resp.Diagnostics.Append(var0.ProfileSettings.Profiles.AntiSpywareProfiles.ElementsAs(ctx, &var1.ProfileSettings.Profiles.AntiSpywareProfiles, false)...)
				resp.Diagnostics.Append(var0.ProfileSettings.Profiles.VulnerabilityProfiles.ElementsAs(ctx, &var1.ProfileSettings.Profiles.VulnerabilityProfiles, false)...)
			}
		}
		if var0.Qos != nil {
			var1.Qos = &security.QosObject{}
			var1.Qos.IpDscp = var0.Qos.IpDscp.ValueStringPointer()
			var1.Qos.IpPrecedence = var0.Qos.IpPrecedence.ValueStringPointer()
			if var0.Qos.FollowClientToServerFlow.ValueBool() {
				var1.Qos.FollowClientToServerFlow = ""
			}
		}
		// TODO: Targets
		var1.NegateTarget = var0.NegateTarget.ValueBoolPointer()
		var1.DisableInspect = var0.DisableInspect.ValueBoolPointer()

		entries = append(entries, var1)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Timeout handling.
	//ctx, cancel := context.WithTimeout(ctx, GetTimeout(state.Timeouts.Create))
	//defer cancel()

	// Get the current list of security rules.
	listing, err = svc.List(ctx, tfid.Location, "get", "", "")
	if err != nil {
		resp.Diagnostics.AddError("Error during Create's refresh", err.Error())
		return
	}

	// Verify no rules are already present.
	for _, x := range listing {
		if isNewName[x.Name] {
			resp.Diagnostics.AddError("Rule to be created already exists", x.Name)
		}
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Prepare the multi-config.
	vn := r.client.Versioning()
	updates := xmlapi.NewMultiConfig(len(entries))
	specifier, _, err := security.Versioning(vn)
	if err != nil {
		resp.Diagnostics.AddError("Error getting specifier", err.Error())
		return
	}
	for _, entry := range entries {
		path, err := tfid.Location.Xpath(vn, entry.Name, "")
		if err != nil {
			resp.Diagnostics.AddError("Error creating path", err.Error())
			return
		}

		elm, err := specifier(entry)
		if err != nil {
			resp.Diagnostics.AddError("Error specifying item", err.Error())
			return
		}

		updates.Add(&xmlapi.Config{
			Action:  "edit",
			Xpath:   util.AsXpath(path),
			Element: elm,
			Target:  r.client.GetTarget(),
		})
	}

	// Create the rules.
	if _, _, _, err = r.client.MultiConfig(ctx, updates, false, nil); err != nil {
		resp.Diagnostics.AddError("Error in create", err.Error())
		return
	}

	// Position the rules.
	if err = svc.MoveGroup(ctx, tfid.Location, tfid.Position, entries); err != nil {
		_ = svc.Delete(ctx, tfid.Location, entryNames...)
		resp.Diagnostics.AddError("Error positioning rules", err.Error())
		return
	}

	// Retrieve the list of rules again.
	listing, err = svc.List(ctx, tfid.Location, "get", "", "")
	if err != nil {
		_ = svc.Delete(ctx, tfid.Location, entryNames...)
		resp.Diagnostics.AddError("Error during Create's refresh", err.Error())
		return
	}

	// Find the first UUID.
	tfid.Rules = make([]RuleInfo, 0, len(entries))
	ans := make([]security.Entry, 0, len(entries))
	first := -1
	for index, live := range listing {
		if live.Name == entries[0].Name {
			first = index
			break
		}
	}

	// First UUID should be present and we shouldn't go out of bounds.
	if first < 0 || first+len(entries) > len(listing) {
		_ = svc.Delete(ctx, tfid.Location, entryNames...)
		resp.Diagnostics.AddError("All rules not present", fmt.Sprintf("%d/%d rules were present", len(ans), len(entries)))
		return
	}

	// Verify rule group appears sequentially.
	for entryIndex := range entries {
		listingIndex := first + entryIndex
		if listing[listingIndex].Name != entries[entryIndex].Name {
			resp.Diagnostics.AddError("Error in create verification", fmt.Sprintf("Expected %q at group index %d", entries[entryIndex].Name, entryIndex))
			continue
		}
		ans = append(ans, listing[listingIndex])
		tfid.Rules = append(tfid.Rules, RuleInfo{
			Name: listing[listingIndex].Name,
			Uuid: listing[listingIndex].Uuid,
		})
	}

	// Verify no errors so far in post-create validation.
	if resp.Diagnostics.HasError() {
		_ = svc.Delete(ctx, tfid.Location, entryNames...)
		resp.Diagnostics.AddError("All rules not present", fmt.Sprintf("%d/%d rules were present", len(ans), len(entries)))
		return
	}

	// Tfid handling.
	tfidstr, err := EncodeLocation(&tfid)
	if err != nil {
		_ = svc.Delete(ctx, tfid.Location, entryNames...)
		resp.Diagnostics.AddError("Error creating tfid", err.Error())
		return
	}

	// Save the state.
	state.Tfid = types.StringValue(tfidstr)
	state.Position.Ok = types.BoolValue(true)
	rules := make([]SecurityPolicyRulesResourceEntry, 0, len(ans))
	for _, var2 := range ans {
		var3 := SecurityPolicyRulesResourceEntry{}
		var3.Name = types.StringValue(var2.Name)
		var3.Uuid = types.StringValue(var2.Uuid)
		var4, var5 := types.SetValueFrom(ctx, types.StringType, var2.SourceZones)
		var3.SourceZones = var4
		resp.Diagnostics.Append(var5.Errors()...)
		var6, var7 := types.SetValueFrom(ctx, types.StringType, var2.DestinationZones)
		var3.DestinationZones = var6
		resp.Diagnostics.Append(var7.Errors()...)
		var8, var9 := types.SetValueFrom(ctx, types.StringType, var2.SourceAddresses)
		var3.SourceAddresses = var8
		resp.Diagnostics.Append(var9.Errors()...)
		var10, var11 := types.SetValueFrom(ctx, types.StringType, var2.SourceUsers)
		var3.SourceUsers = var10
		resp.Diagnostics.Append(var11.Errors()...)
		var12, var13 := types.SetValueFrom(ctx, types.StringType, var2.DestinationAddresses)
		var3.DestinationAddresses = var12
		resp.Diagnostics.Append(var13.Errors()...)
		var14, var15 := types.SetValueFrom(ctx, types.StringType, var2.Services)
		var3.Services = var14
		resp.Diagnostics.Append(var15.Errors()...)
		var16, var17 := types.SetValueFrom(ctx, types.StringType, var2.Categories)
		var3.Categories = var16
		resp.Diagnostics.Append(var17.Errors()...)
		var18, var19 := types.SetValueFrom(ctx, types.StringType, var2.Applications)
		var3.Applications = var18
		resp.Diagnostics.Append(var19.Errors()...)
		var20, var21 := types.SetValueFrom(ctx, types.StringType, var2.SourceDevices)
		var3.SourceDevices = var20
		resp.Diagnostics.Append(var21.Errors()...)
		var22, var23 := types.SetValueFrom(ctx, types.StringType, var2.DestinationDevices)
		var3.DestinationDevices = var22
		resp.Diagnostics.Append(var23.Errors()...)
		var3.Schedule = types.StringPointerValue(var2.Schedule)
		var24, var25 := types.ListValueFrom(ctx, types.StringType, var2.Tags)
		var3.Tags = var24
		resp.Diagnostics.Append(var25.Errors()...)
		var3.NegateSource = types.BoolPointerValue(var2.NegateSource)
		var3.NegateDestination = types.BoolPointerValue(var2.NegateDestination)
		var3.Disabled = types.BoolPointerValue(var2.Disabled)
		var3.Description = types.StringPointerValue(var2.Description)
		var3.GroupTag = types.StringPointerValue(var2.GroupTag)
		var3.Action = types.StringValue(var2.Action)
		var3.IcmpUnreachable = types.BoolPointerValue(var2.IcmpUnreachable)
		var3.Type = types.StringPointerValue(var2.Type)
		var3.DisableServerResponseInspection = types.BoolPointerValue(var2.DisableServerResponseInspection)
		var3.LogSetting = types.StringPointerValue(var2.LogSetting)
		var3.LogStart = types.BoolPointerValue(var2.LogStart)
		var3.LogEnd = types.BoolPointerValue(var2.LogEnd)
		if var2.ProfileSettings != nil {
			var3.ProfileSettings = &SecurityPolicyRulesResourceProfileSettingsObject{}
			var26, var27 := types.ListValueFrom(ctx, types.StringType, var2.ProfileSettings.Groups)
			var3.ProfileSettings.Groups = var26
			resp.Diagnostics.Append(var27.Errors()...)
			if var2.ProfileSettings.Profiles != nil {
				var3.ProfileSettings.Profiles = &SecurityPolicyRulesResourceProfilesObject{}
				var28, var29 := types.ListValueFrom(ctx, types.StringType, var2.ProfileSettings.Profiles.UrlFilteringProfiles)
				var3.ProfileSettings.Profiles.UrlFilteringProfiles = var28
				resp.Diagnostics.Append(var29.Errors()...)
				var30, var31 := types.ListValueFrom(ctx, types.StringType, var2.ProfileSettings.Profiles.DataFilteringProfiles)
				var3.ProfileSettings.Profiles.DataFilteringProfiles = var30
				resp.Diagnostics.Append(var31.Errors()...)
				var32, var33 := types.ListValueFrom(ctx, types.StringType, var2.ProfileSettings.Profiles.FileBlockingProfiles)
				var3.ProfileSettings.Profiles.FileBlockingProfiles = var32
				resp.Diagnostics.Append(var33.Errors()...)
				var34, var35 := types.ListValueFrom(ctx, types.StringType, var2.ProfileSettings.Profiles.WildfireAnalysisProfiles)
				var3.ProfileSettings.Profiles.WildfireAnalysisProfiles = var34
				resp.Diagnostics.Append(var35.Errors()...)
				var36, var37 := types.ListValueFrom(ctx, types.StringType, var2.ProfileSettings.Profiles.AntiVirusProfiles)
				var3.ProfileSettings.Profiles.AntiVirusProfiles = var36
				resp.Diagnostics.Append(var37.Errors()...)
				var38, var39 := types.ListValueFrom(ctx, types.StringType, var2.ProfileSettings.Profiles.AntiSpywareProfiles)
				var3.ProfileSettings.Profiles.AntiSpywareProfiles = var38
				resp.Diagnostics.Append(var39.Errors()...)
				var40, var41 := types.ListValueFrom(ctx, types.StringType, var2.ProfileSettings.Profiles.VulnerabilityProfiles)
				var3.ProfileSettings.Profiles.VulnerabilityProfiles = var40
				resp.Diagnostics.Append(var41.Errors()...)
			}
		}
		if var2.Qos != nil {
			var3.Qos = &SecurityPolicyRulesResourceQosObject{}
			var3.Qos.IpDscp = types.StringPointerValue(var2.Qos.IpDscp)
			var3.Qos.IpPrecedence = types.StringPointerValue(var2.Qos.IpPrecedence)
			if var2.Qos.FollowClientToServerFlow == nil {
				var3.Qos.FollowClientToServerFlow = types.BoolValue(false)
			} else {
				var3.Qos.FollowClientToServerFlow = types.BoolValue(true)
			}
		}
		// TODO: Targets
		var3.NegateTarget = types.BoolPointerValue(var2.NegateTarget)
		var3.DisableInspect = types.BoolPointerValue(var2.DisableInspect)

		rules = append(rules, var3)
	}
	state.Rules = rules
	state.Ordered = types.BoolValue(true)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *SecurityPolicyRulesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var err error
	var savestate, state SecurityPolicyRulesResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &savestate)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the tfid info.
	var tfid SecurityPolicyRulesTfid
	if err = DecodeLocation(savestate.Tfid.ValueString(), &tfid); err != nil {
		resp.Diagnostics.AddError("Error parsing tfid", err.Error())
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource read", map[string]any{
		"resource_name": "panos_security_policy_rules",
		"function":      "Read",
	})

	// Verify mode.
	if r.client.Hostname == "" {
		resp.Diagnostics.AddError("Invalid mode error", InspectionModeError)
		return
	}

	// Create the service.
	svc := security.NewService(r.client)

	// Timeout handling.
	//ctx, cancel := context.WithTimeout(ctx, GetTimeout(savestate.Timeouts.Read))
	//defer cancel()

	// Perform a list to get all rules.
	listing, err := svc.List(ctx, tfid.Location, "get", "", "")
	if err != nil {
		resp.Diagnostics.AddError("Error in listing", err.Error())
		return
	}

	// If there are no rules, we can remove the state, we need a full redeploy.
	if len(listing) == 0 {
		resp.State.RemoveResource(ctx)
		return
	}

	// Build the UUID map.
	uuidMap := make(map[string]int, len(listing))
	for index, live := range listing {
		uuidMap[live.Uuid] = index
	}

	// Find the rules.
	ordered := true
	prev := -1
	rules := make([]SecurityPolicyRulesResourceEntry, 0, len(tfid.Rules))
	for index, rInfo := range tfid.Rules {
		rid, ok := uuidMap[rInfo.Uuid]
		if !ok {
			ordered = false
			prev = -1
			continue
		}

		if index != 0 && ordered {
			ordered = prev+1 == rid
		}

		prev = rid
		var2 := listing[rid]

		var3 := SecurityPolicyRulesResourceEntry{}
		var3.Name = types.StringValue(var2.Name)
		var3.Uuid = types.StringValue(var2.Uuid)
		var4, var5 := types.SetValueFrom(ctx, types.StringType, var2.SourceZones)
		var3.SourceZones = var4
		resp.Diagnostics.Append(var5.Errors()...)
		var6, var7 := types.SetValueFrom(ctx, types.StringType, var2.DestinationZones)
		var3.DestinationZones = var6
		resp.Diagnostics.Append(var7.Errors()...)
		var8, var9 := types.SetValueFrom(ctx, types.StringType, var2.SourceAddresses)
		var3.SourceAddresses = var8
		resp.Diagnostics.Append(var9.Errors()...)
		var10, var11 := types.SetValueFrom(ctx, types.StringType, var2.SourceUsers)
		var3.SourceUsers = var10
		resp.Diagnostics.Append(var11.Errors()...)
		var12, var13 := types.SetValueFrom(ctx, types.StringType, var2.DestinationAddresses)
		var3.DestinationAddresses = var12
		resp.Diagnostics.Append(var13.Errors()...)
		var14, var15 := types.SetValueFrom(ctx, types.StringType, var2.Services)
		var3.Services = var14
		resp.Diagnostics.Append(var15.Errors()...)
		var16, var17 := types.SetValueFrom(ctx, types.StringType, var2.Categories)
		var3.Categories = var16
		resp.Diagnostics.Append(var17.Errors()...)
		var18, var19 := types.SetValueFrom(ctx, types.StringType, var2.Applications)
		var3.Applications = var18
		resp.Diagnostics.Append(var19.Errors()...)
		var20, var21 := types.SetValueFrom(ctx, types.StringType, var2.SourceDevices)
		var3.SourceDevices = var20
		resp.Diagnostics.Append(var21.Errors()...)
		var22, var23 := types.SetValueFrom(ctx, types.StringType, var2.DestinationDevices)
		var3.DestinationDevices = var22
		resp.Diagnostics.Append(var23.Errors()...)
		var3.Schedule = types.StringPointerValue(var2.Schedule)
		var24, var25 := types.ListValueFrom(ctx, types.StringType, var2.Tags)
		var3.Tags = var24
		resp.Diagnostics.Append(var25.Errors()...)
		var3.NegateSource = types.BoolPointerValue(var2.NegateSource)
		var3.NegateDestination = types.BoolPointerValue(var2.NegateDestination)
		var3.Disabled = types.BoolPointerValue(var2.Disabled)
		var3.Description = types.StringPointerValue(var2.Description)
		var3.GroupTag = types.StringPointerValue(var2.GroupTag)
		var3.Action = types.StringValue(var2.Action)
		var3.IcmpUnreachable = types.BoolPointerValue(var2.IcmpUnreachable)
		var3.Type = types.StringPointerValue(var2.Type)
		var3.DisableServerResponseInspection = types.BoolPointerValue(var2.DisableServerResponseInspection)
		var3.LogSetting = types.StringPointerValue(var2.LogSetting)
		var3.LogStart = types.BoolPointerValue(var2.LogStart)
		var3.LogEnd = types.BoolPointerValue(var2.LogEnd)
		if var2.ProfileSettings != nil {
			var3.ProfileSettings = &SecurityPolicyRulesResourceProfileSettingsObject{}
			var26, var27 := types.ListValueFrom(ctx, types.StringType, var2.ProfileSettings.Groups)
			var3.ProfileSettings.Groups = var26
			resp.Diagnostics.Append(var27.Errors()...)
			if var2.ProfileSettings.Profiles != nil {
				var3.ProfileSettings.Profiles = &SecurityPolicyRulesResourceProfilesObject{}
				var28, var29 := types.ListValueFrom(ctx, types.StringType, var2.ProfileSettings.Profiles.UrlFilteringProfiles)
				var3.ProfileSettings.Profiles.UrlFilteringProfiles = var28
				resp.Diagnostics.Append(var29.Errors()...)
				var30, var31 := types.ListValueFrom(ctx, types.StringType, var2.ProfileSettings.Profiles.DataFilteringProfiles)
				var3.ProfileSettings.Profiles.DataFilteringProfiles = var30
				resp.Diagnostics.Append(var31.Errors()...)
				var32, var33 := types.ListValueFrom(ctx, types.StringType, var2.ProfileSettings.Profiles.FileBlockingProfiles)
				var3.ProfileSettings.Profiles.FileBlockingProfiles = var32
				resp.Diagnostics.Append(var33.Errors()...)
				var34, var35 := types.ListValueFrom(ctx, types.StringType, var2.ProfileSettings.Profiles.WildfireAnalysisProfiles)
				var3.ProfileSettings.Profiles.WildfireAnalysisProfiles = var34
				resp.Diagnostics.Append(var35.Errors()...)
				var36, var37 := types.ListValueFrom(ctx, types.StringType, var2.ProfileSettings.Profiles.AntiVirusProfiles)
				var3.ProfileSettings.Profiles.AntiVirusProfiles = var36
				resp.Diagnostics.Append(var37.Errors()...)
				var38, var39 := types.ListValueFrom(ctx, types.StringType, var2.ProfileSettings.Profiles.AntiSpywareProfiles)
				var3.ProfileSettings.Profiles.AntiSpywareProfiles = var38
				resp.Diagnostics.Append(var39.Errors()...)
				var40, var41 := types.ListValueFrom(ctx, types.StringType, var2.ProfileSettings.Profiles.VulnerabilityProfiles)
				var3.ProfileSettings.Profiles.VulnerabilityProfiles = var40
				resp.Diagnostics.Append(var41.Errors()...)
			}
		}
		if var2.Qos != nil {
			var3.Qos = &SecurityPolicyRulesResourceQosObject{}
			var3.Qos.IpDscp = types.StringPointerValue(var2.Qos.IpDscp)
			var3.Qos.IpPrecedence = types.StringPointerValue(var2.Qos.IpPrecedence)
			if var2.Qos.FollowClientToServerFlow == nil {
				var3.Qos.FollowClientToServerFlow = types.BoolValue(false)
			} else {
				var3.Qos.FollowClientToServerFlow = types.BoolValue(true)
			}
		}
		// TODO: Targets
		var3.NegateTarget = types.BoolPointerValue(var2.NegateTarget)
		var3.DisableInspect = types.BoolPointerValue(var2.DisableInspect)

		rules = append(rules, var3)
	}

	// If there are no rules, we can remove the state, we need a full redeploy.
	if len(rules) == 0 {
		resp.State.RemoveResource(ctx)
		return
	}

	// Save the location.
	if tfid.Location.Shared != nil {
		state.Location.Shared = &SecurityPolicyRulesResourceSharedLocation{
			Rulebase: types.StringValue(tfid.Location.Shared.Rulebase),
		}
	}
	if tfid.Location.Vsys != nil {
		state.Location.Vsys = &SecurityPolicyRulesResourceVsysLocation{
			NgfwDevice: types.StringValue(tfid.Location.Vsys.NgfwDevice),
			Name:       types.StringValue(tfid.Location.Vsys.Name),
		}
	}
	if tfid.Location.DeviceGroup != nil {
		state.Location.DeviceGroup = &SecurityPolicyRulesResourceDeviceGroupLocation{
			PanoramaDevice: types.StringValue(tfid.Location.DeviceGroup.PanoramaDevice),
			Name:           types.StringValue(tfid.Location.DeviceGroup.Name),
			Rulebase:       types.StringValue(tfid.Location.DeviceGroup.Rulebase),
		}
	}

	// Save the position.
	state.Position = SecurityPolicyRulesResourcePosition{
		Ok:              types.BoolValue(false),
		First:           types.BoolPointerValue(tfid.Position.First),
		Last:            types.BoolPointerValue(tfid.Position.Last),
		SomewhereBefore: types.StringPointerValue(tfid.Position.SomewhereBefore),
		SomewhereAfter:  types.StringPointerValue(tfid.Position.SomewhereAfter),
		DirectlyBefore:  types.StringPointerValue(tfid.Position.DirectlyBefore),
		DirectlyAfter:   types.StringPointerValue(tfid.Position.DirectlyAfter),
	}
	switch {
	case tfid.Position.First != nil && *tfid.Position.First:
		if tfid.Rules[0].Uuid == listing[0].Uuid {
			state.Position.Ok = types.BoolValue(true)
		}
	case tfid.Position.Last != nil && *tfid.Position.Last:
		if tfid.Rules[len(tfid.Rules)-1].Uuid == listing[len(listing)-1].Uuid {
			state.Position.Ok = types.BoolValue(true)
		}
	case tfid.Position.SomewhereBefore != nil:
		var found bool
		for index, live := range listing {
			if live.Name == *tfid.Position.SomewhereBefore {
				found = true
				if rid, ok := uuidMap[tfid.Rules[len(tfid.Rules)-1].Uuid]; ok {
					if rid < index {
						state.Position.Ok = types.BoolValue(true)
					}
				}
				break
			}
		}
		if !found {
			resp.Diagnostics.AddError("Error in rule positioning", fmt.Sprintf("Cannot place group somewhere before %q - referenced rule does not exist", *tfid.Position.SomewhereBefore))
			return
		}
	case tfid.Position.SomewhereAfter != nil:
		var found bool
		for index, live := range listing {
			if live.Name == *tfid.Position.SomewhereAfter {
				found = true
				if rid, ok := uuidMap[tfid.Rules[0].Uuid]; ok {
					if rid > index {
						state.Position.Ok = types.BoolValue(true)
					}
				}
				break
			}
		}
		if !found {
			resp.Diagnostics.AddError("Error in rule positioning", fmt.Sprintf("Cannot place group somewhere after %q - referenced rule does not exist", *tfid.Position.SomewhereAfter))
			return
		}
	case tfid.Position.DirectlyBefore != nil:
		var found bool
		for index, live := range listing {
			if live.Name == *tfid.Position.DirectlyBefore {
				found = true
				if rid, ok := uuidMap[tfid.Rules[len(tfid.Rules)-1].Uuid]; ok {
					if rid+1 == index {
						state.Position.Ok = types.BoolValue(true)
					}
				}
				break
			}
		}
		if !found {
			resp.Diagnostics.AddError("Error in rule positioning", fmt.Sprintf("Cannot place group directly before %q - referenced rule does not exist", *tfid.Position.DirectlyBefore))
			return
		}
	case tfid.Position.DirectlyAfter != nil:
		var found bool
		for index, live := range listing {
			if live.Name == *tfid.Position.DirectlyAfter {
				found = true
				if rid, ok := uuidMap[tfid.Rules[0].Uuid]; ok {
					if rid == index+1 {
						state.Position.Ok = types.BoolValue(true)
					}
				}
				break
			}
		}
		if !found {
			resp.Diagnostics.AddError("Error in rule positioning", fmt.Sprintf("Cannot place group directly after %q - referenced rule does not exist", *tfid.Position.DirectlyAfter))
			return
		}
	default:
		state.Position.Ok = types.BoolValue(true)
	}

	// Save state.
	state.Tfid = savestate.Tfid
	state.Rules = rules
	state.Ordered = types.BoolValue(ordered)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *SecurityPolicyRulesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var err error
	var updates *xmlapi.MultiConfig
	var plan, state SecurityPolicyRulesResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the tfid info.
	var tfid, newtfid SecurityPolicyRulesTfid
	if err = DecodeLocation(state.Tfid.ValueString(), &tfid); err != nil {
		resp.Diagnostics.AddError("Error parsing tfid", err.Error())
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource read", map[string]any{
		"resource_name": "panos_security_policy_rules",
		"function":      "Update",
	})

	// Determine the location.
	if plan.Location.Shared != nil {
		newtfid.Location.Shared = &security.SharedLocation{
			Rulebase: plan.Location.Shared.Rulebase.ValueString(),
		}
	}
	if plan.Location.Vsys != nil {
		newtfid.Location.Vsys = &security.VsysLocation{
			NgfwDevice: plan.Location.Vsys.NgfwDevice.ValueString(),
			Name:       plan.Location.Vsys.Name.ValueString(),
		}
	}
	if plan.Location.DeviceGroup != nil {
		newtfid.Location.DeviceGroup = &security.DeviceGroupLocation{
			PanoramaDevice: plan.Location.DeviceGroup.PanoramaDevice.ValueString(),
			Rulebase:       plan.Location.DeviceGroup.Rulebase.ValueString(),
			Name:           plan.Location.DeviceGroup.Name.ValueString(),
		}
	}
	if err = newtfid.Location.IsValid(); err != nil {
		resp.Diagnostics.AddError("Invalid path", err.Error())
		return
	}

	// Determine the position.
	newtfid.Position = rule.Position{
		First:           plan.Position.First.ValueBoolPointer(),
		Last:            plan.Position.Last.ValueBoolPointer(),
		SomewhereBefore: plan.Position.SomewhereBefore.ValueStringPointer(),
		SomewhereAfter:  plan.Position.SomewhereAfter.ValueStringPointer(),
		DirectlyBefore:  plan.Position.DirectlyBefore.ValueStringPointer(),
		DirectlyAfter:   plan.Position.DirectlyAfter.ValueStringPointer(),
	}
	if err = newtfid.Position.IsValid(false); err != nil {
		resp.Diagnostics.AddError("Invalid position", err.Error())
		return
	}

	// Verify mode.
	if r.client.Hostname == "" {
		resp.Diagnostics.AddError("Invalid mode error", InspectionModeError)
		return
	}

	// Create the service.
	svc := security.NewService(r.client)

	// Prepare to handle versioning.
	vn := r.client.Versioning()
	specifier, _, err := security.Versioning(vn)
	if err != nil {
		resp.Diagnostics.AddError("Error getting specifier", err.Error())
		return
	}

	// Load the desired config.
	entries := make([]security.Entry, 0, len(plan.Rules))
	for _, var0 := range plan.Rules {
		var1 := security.Entry{Name: var0.Name.ValueString()}
		resp.Diagnostics.Append(var0.SourceZones.ElementsAs(ctx, &var1.SourceZones, false)...)
		resp.Diagnostics.Append(var0.DestinationZones.ElementsAs(ctx, &var1.DestinationZones, false)...)
		resp.Diagnostics.Append(var0.SourceAddresses.ElementsAs(ctx, &var1.SourceAddresses, false)...)
		resp.Diagnostics.Append(var0.SourceUsers.ElementsAs(ctx, &var1.SourceUsers, false)...)
		resp.Diagnostics.Append(var0.DestinationAddresses.ElementsAs(ctx, &var1.DestinationAddresses, false)...)
		resp.Diagnostics.Append(var0.Services.ElementsAs(ctx, &var1.Services, false)...)
		resp.Diagnostics.Append(var0.Categories.ElementsAs(ctx, &var1.Categories, false)...)
		resp.Diagnostics.Append(var0.Applications.ElementsAs(ctx, &var1.Applications, false)...)
		resp.Diagnostics.Append(var0.SourceDevices.ElementsAs(ctx, &var1.SourceDevices, false)...)
		resp.Diagnostics.Append(var0.DestinationDevices.ElementsAs(ctx, &var1.DestinationDevices, false)...)
		var1.Schedule = var0.Schedule.ValueStringPointer()
		resp.Diagnostics.Append(var0.Tags.ElementsAs(ctx, &var1.Tags, false)...)
		var1.NegateSource = var0.NegateSource.ValueBoolPointer()
		var1.NegateDestination = var0.NegateDestination.ValueBoolPointer()
		var1.Disabled = var0.Disabled.ValueBoolPointer()
		var1.Description = var0.Description.ValueStringPointer()
		var1.GroupTag = var0.GroupTag.ValueStringPointer()
		var1.Action = var0.Action.ValueString()
		var1.IcmpUnreachable = var0.IcmpUnreachable.ValueBoolPointer()
		var1.Type = var0.Type.ValueStringPointer()
		var1.DisableServerResponseInspection = var0.DisableServerResponseInspection.ValueBoolPointer()
		var1.LogSetting = var0.LogSetting.ValueStringPointer()
		var1.LogStart = var0.LogStart.ValueBoolPointer()
		var1.LogEnd = var0.LogEnd.ValueBoolPointer()
		if var0.ProfileSettings != nil {
			var1.ProfileSettings = &security.ProfileSettingsObject{}
			resp.Diagnostics.Append(var0.ProfileSettings.Groups.ElementsAs(ctx, &var1.ProfileSettings.Groups, false)...)
			if var0.ProfileSettings.Profiles != nil {
				var1.ProfileSettings.Profiles = &security.ProfilesObject{}
				resp.Diagnostics.Append(var0.ProfileSettings.Profiles.UrlFilteringProfiles.ElementsAs(ctx, &var1.ProfileSettings.Profiles.UrlFilteringProfiles, false)...)
				resp.Diagnostics.Append(var0.ProfileSettings.Profiles.DataFilteringProfiles.ElementsAs(ctx, &var1.ProfileSettings.Profiles.DataFilteringProfiles, false)...)
				resp.Diagnostics.Append(var0.ProfileSettings.Profiles.FileBlockingProfiles.ElementsAs(ctx, &var1.ProfileSettings.Profiles.FileBlockingProfiles, false)...)
				resp.Diagnostics.Append(var0.ProfileSettings.Profiles.WildfireAnalysisProfiles.ElementsAs(ctx, &var1.ProfileSettings.Profiles.WildfireAnalysisProfiles, false)...)
				resp.Diagnostics.Append(var0.ProfileSettings.Profiles.AntiVirusProfiles.ElementsAs(ctx, &var1.ProfileSettings.Profiles.AntiVirusProfiles, false)...)
				resp.Diagnostics.Append(var0.ProfileSettings.Profiles.AntiSpywareProfiles.ElementsAs(ctx, &var1.ProfileSettings.Profiles.AntiSpywareProfiles, false)...)
				resp.Diagnostics.Append(var0.ProfileSettings.Profiles.VulnerabilityProfiles.ElementsAs(ctx, &var1.ProfileSettings.Profiles.VulnerabilityProfiles, false)...)
			}
		}
		if var0.Qos != nil {
			var1.Qos = &security.QosObject{}
			var1.Qos.IpDscp = var0.Qos.IpDscp.ValueStringPointer()
			var1.Qos.IpPrecedence = var0.Qos.IpPrecedence.ValueStringPointer()
			if var0.Qos.FollowClientToServerFlow.ValueBool() {
				var1.Qos.FollowClientToServerFlow = ""
			}
		}
		// TODO: Targets
		var1.NegateTarget = var0.NegateTarget.ValueBoolPointer()
		var1.DisableInspect = var0.DisableInspect.ValueBoolPointer()

		entries = append(entries, var1)
	}

	// Get the list of all rules.
	listing, err := svc.List(ctx, newtfid.Location, "get", "", "")
	if err != nil {
		resp.Diagnostics.AddError("Error in refresh for update", err.Error())
		return
	}

	managedUuidIsUsed := make(map[string]bool)
	uuidMap := make(map[string]int, len(listing))
	nameMap := make(map[string]int, len(listing))
	for index, live := range listing {
		uuidMap[live.Uuid] = index
		nameMap[live.Name] = index
	}

	updates = xmlapi.NewMultiConfig(len(entries) * 2)

	delayed := make([]security.Entry, 0, len(entries))
	for _, entry := range entries {
		var ruleInfo *RuleInfo
		for index := range tfid.Rules {
			if tfid.Rules[index].Name == entry.Name {
				ruleInfo = &tfid.Rules[index]
				break
			}
		}

		// Delay processing if this is a new rule or a rename.
		if ruleInfo == nil {
			delayed = append(delayed, entry)
			continue
		}

		// Verify the rule name we want isn't taken by another rule we don't
		// manage (someone deleted our rule out-of-band then created a new
		// rule with the same name; this results in a unique UUID).
		rid, ok := nameMap[entry.Name]
		if ok && ruleInfo.Uuid != listing[rid].Uuid {
			resp.Diagnostics.AddError(fmt.Sprintf("Rule was deleted and replaced: %s", ruleInfo.Name), fmt.Sprintf("uuid has changed: old:%s new:%s", ruleInfo.Uuid, listing[rid].Uuid))
			return
		}

		managedUuidIsUsed[ruleInfo.Uuid] = true

		// Check if the UUID we made has been deleted out-of-band.
		listingIndex, ok := uuidMap[ruleInfo.Uuid]
		if !ok {
			delayed = append(delayed, entry)
			continue
		}

		// Check if we need to rename the rule.
		if listing[listingIndex].Name != entry.Name {
			path, err := newtfid.Location.Xpath(vn, listing[listingIndex].Name, "")
			if err != nil {
				resp.Diagnostics.AddError("Error creating rename xpath", fmt.Sprintf("(%s) %s - %s", listing[listingIndex].Uuid, listing[listingIndex].Name, err))
				return
			}

			updates.Add(&xmlapi.Config{
				Action:  "rename",
				Xpath:   util.AsXpath(path),
				NewName: entry.Name,
				Target:  r.client.GetTarget(),
			})

			delete(nameMap, listing[listingIndex].Name)
			nameMap[entry.Name] = listingIndex
		}

		// The rule name and UUID matches our records, verify spec.
		if !security.SpecMatches(&entry, &listing[listingIndex]) {
			path, err := newtfid.Location.Xpath(vn, entry.Name, "")
			if err != nil {
				resp.Diagnostics.AddError("Error creating update xpath", fmt.Sprintf("(%s) %s - %s", ruleInfo.Uuid, ruleInfo.Name, err))
				return
			}

			entry.CopyMiscFrom(&listing[listingIndex])
			elm, err := specifier(entry)
			if err != nil {
				resp.Diagnostics.AddError("Error specifying update", fmt.Sprintf("(%s) %s - %s", ruleInfo.Uuid, ruleInfo.Name, err))
				return
			}

			updates.Add(&xmlapi.Config{
				Action:  "edit",
				Xpath:   util.AsXpath(path),
				Element: elm,
				Target:  r.client.GetTarget(),
			})
		}
	}

	// Now that we've done a first pass, we need to address the delayed entries. Any
	// rules in here are either renamed from a previous rule or brand new. So first we
	// should reuse any UUIDs that we're responsible for, then after that we will just
	// create new rules.
	for _, entry := range delayed {
		if rid, ok := nameMap[entry.Name]; ok {
			resp.Diagnostics.AddError("Rule already exists", listing[rid].Name)
			return
		}

		// See if we can repurpose a previous rule UUID first.
		var ruleInfo *RuleInfo
		if len(managedUuidIsUsed) != len(tfid.Rules) {
			for index := range tfid.Rules {
				if !managedUuidIsUsed[tfid.Rules[index].Uuid] {
					managedUuidIsUsed[tfid.Rules[index].Uuid] = true
					if _, ok := uuidMap[tfid.Rules[index].Uuid]; ok {
						ruleInfo = &tfid.Rules[index]
						break
					}
				}
			}
		}

		if ruleInfo == nil {
			path, err := newtfid.Location.Xpath(vn, entry.Name, "")
			if err != nil {
				resp.Diagnostics.AddError("Error creating missing elm xpath", fmt.Sprintf("name:%q - %s", entry.Name, err))
				return
			}

			elm, err := specifier(entry)
			if err != nil {
				resp.Diagnostics.AddError("Error specifying missing elm Element", fmt.Sprintf("name:%q - %s", entry.Name, err))
				return
			}

			updates.Add(&xmlapi.Config{
				Action:  "edit",
				Xpath:   util.AsXpath(path),
				Element: elm,
				Target:  r.client.GetTarget(),
			})
		} else {
			listingIndex := uuidMap[ruleInfo.Uuid]
			// Rename the old rule to the desired name.
			path, err := newtfid.Location.Xpath(vn, listing[listingIndex].Name, "")
			if err != nil {
				resp.Diagnostics.AddError("Error creating repurpose rename xpath", fmt.Sprintf("(%s) %s - %s", listing[listingIndex].Uuid, listing[listingIndex].Name, err))
				return
			}

			updates.Add(&xmlapi.Config{
				Action:  "rename",
				Xpath:   util.AsXpath(path),
				NewName: entry.Name,
				Target:  r.client.GetTarget(),
			})

			if !security.SpecMatches(&entry, &listing[listingIndex]) {
				path, err := newtfid.Location.Xpath(vn, entry.Name, "")
				if err != nil {
					resp.Diagnostics.AddError("Error creating update xpath", fmt.Sprintf("(%s) %s - %s", ruleInfo.Uuid, ruleInfo.Name, err))
					return
				}

				entry.CopyMiscFrom(&listing[listingIndex])
				elm, err := specifier(entry)
				if err != nil {
					resp.Diagnostics.AddError("Error specifying update", fmt.Sprintf("(%s) %s - %s", ruleInfo.Uuid, ruleInfo.Name, err))
					return
				}

				updates.Add(&xmlapi.Config{
					Action:  "edit",
					Xpath:   util.AsXpath(path),
					Element: elm,
					Target:  r.client.GetTarget(),
				})
			}
		}
	}

	// Finally delete any rules of UUIDs we manage that are no longer used.
	for _, ruleInfo := range tfid.Rules {
		if _, ok := uuidMap[ruleInfo.Uuid]; !managedUuidIsUsed[ruleInfo.Uuid] && ok {
			path, err := newtfid.Location.Xpath(vn, "", ruleInfo.Uuid)
			if err != nil {
				resp.Diagnostics.AddError("Error building unused uuid delete xpath", fmt.Sprintf("(%s) %s", ruleInfo.Uuid, err))
				return
			}

			updates.Add(&xmlapi.Config{
				Action: "delete",
				Xpath:  util.AsXpath(path),
				Target: r.client.GetTarget(),
			})
		}
	}

	// Perform updates if needed.
	if len(updates.Operations) > 0 {
		if _, _, _, err = r.client.MultiConfig(ctx, updates, false, nil); err != nil {
			resp.Diagnostics.AddError("Error performing multi-config", err.Error())
			return
		}
	}

	// Position the rules.
	if err = svc.MoveGroup(ctx, newtfid.Location, newtfid.Position, entries); err != nil {
		resp.Diagnostics.AddError("Error positioning rules", err.Error())
		return
	}

	// Retrieve the list of rules again.
	listing, err = svc.List(ctx, newtfid.Location, "get", "", "")
	if err != nil {
		resp.Diagnostics.AddError("Error during Create's refresh", err.Error())
		return
	}

	// Find the first UUID.
	newtfid.Rules = make([]RuleInfo, 0, len(entries))
	ans := make([]security.Entry, 0, len(entries))
	first := -1
	for index, live := range listing {
		if live.Name == entries[0].Name {
			first = index
			break
		}
	}

	// First UUID should be present and we shouldn't go out of bounds.
	if first < 0 || first+len(entries) > len(listing) {
		resp.Diagnostics.AddError("All rules not present", fmt.Sprintf("%d/%d rules were present", len(ans), len(entries)))
		return
	}

	// Verify rule group appears sequentially.
	for entryIndex := range entries {
		listingIndex := first + entryIndex
		if listing[listingIndex].Name != entries[entryIndex].Name {
			resp.Diagnostics.AddError("Error in create verification", fmt.Sprintf("Expected %q at group index %d", entries[entryIndex].Name, entryIndex))
			continue
		}
		ans = append(ans, listing[listingIndex])
		newtfid.Rules = append(newtfid.Rules, RuleInfo{
			Name: listing[listingIndex].Name,
			Uuid: listing[listingIndex].Uuid,
		})
	}

	// Verify no errors so far in post-create validation.
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("All rules not present", fmt.Sprintf("%d/%d rules were present", len(ans), len(entries)))
		return
	}

	// Tfid handling.
	tfidstr, err := EncodeLocation(&newtfid)
	if err != nil {
		resp.Diagnostics.AddError("Error creating new tfid", err.Error())
		return
	}

	// Save the state.
	state.Tfid = types.StringValue(tfidstr)
	state.Position = SecurityPolicyRulesResourcePosition{
		Ok:              types.BoolValue(true),
		First:           types.BoolPointerValue(newtfid.Position.First),
		Last:            types.BoolPointerValue(newtfid.Position.Last),
		SomewhereBefore: types.StringPointerValue(newtfid.Position.SomewhereBefore),
		SomewhereAfter:  types.StringPointerValue(newtfid.Position.SomewhereAfter),
		DirectlyBefore:  types.StringPointerValue(newtfid.Position.DirectlyBefore),
		DirectlyAfter:   types.StringPointerValue(newtfid.Position.DirectlyAfter),
	}
	rules := make([]SecurityPolicyRulesResourceEntry, 0, len(ans))
	for _, var2 := range ans {
		var3 := SecurityPolicyRulesResourceEntry{}
		var3.Name = types.StringValue(var2.Name)
		var3.Uuid = types.StringValue(var2.Uuid)
		var4, var5 := types.SetValueFrom(ctx, types.StringType, var2.SourceZones)
		var3.SourceZones = var4
		resp.Diagnostics.Append(var5.Errors()...)
		var6, var7 := types.SetValueFrom(ctx, types.StringType, var2.DestinationZones)
		var3.DestinationZones = var6
		resp.Diagnostics.Append(var7.Errors()...)
		var8, var9 := types.SetValueFrom(ctx, types.StringType, var2.SourceAddresses)
		var3.SourceAddresses = var8
		resp.Diagnostics.Append(var9.Errors()...)
		var10, var11 := types.SetValueFrom(ctx, types.StringType, var2.SourceUsers)
		var3.SourceUsers = var10
		resp.Diagnostics.Append(var11.Errors()...)
		var12, var13 := types.SetValueFrom(ctx, types.StringType, var2.DestinationAddresses)
		var3.DestinationAddresses = var12
		resp.Diagnostics.Append(var13.Errors()...)
		var14, var15 := types.SetValueFrom(ctx, types.StringType, var2.Services)
		var3.Services = var14
		resp.Diagnostics.Append(var15.Errors()...)
		var16, var17 := types.SetValueFrom(ctx, types.StringType, var2.Categories)
		var3.Categories = var16
		resp.Diagnostics.Append(var17.Errors()...)
		var18, var19 := types.SetValueFrom(ctx, types.StringType, var2.Applications)
		var3.Applications = var18
		resp.Diagnostics.Append(var19.Errors()...)
		var20, var21 := types.SetValueFrom(ctx, types.StringType, var2.SourceDevices)
		var3.SourceDevices = var20
		resp.Diagnostics.Append(var21.Errors()...)
		var22, var23 := types.SetValueFrom(ctx, types.StringType, var2.DestinationDevices)
		var3.DestinationDevices = var22
		resp.Diagnostics.Append(var23.Errors()...)
		var3.Schedule = types.StringPointerValue(var2.Schedule)
		var24, var25 := types.ListValueFrom(ctx, types.StringType, var2.Tags)
		var3.Tags = var24
		resp.Diagnostics.Append(var25.Errors()...)
		var3.NegateSource = types.BoolPointerValue(var2.NegateSource)
		var3.NegateDestination = types.BoolPointerValue(var2.NegateDestination)
		var3.Disabled = types.BoolPointerValue(var2.Disabled)
		var3.Description = types.StringPointerValue(var2.Description)
		var3.GroupTag = types.StringPointerValue(var2.GroupTag)
		var3.Action = types.StringValue(var2.Action)
		var3.IcmpUnreachable = types.BoolPointerValue(var2.IcmpUnreachable)
		var3.Type = types.StringPointerValue(var2.Type)
		var3.DisableServerResponseInspection = types.BoolPointerValue(var2.DisableServerResponseInspection)
		var3.LogSetting = types.StringPointerValue(var2.LogSetting)
		var3.LogStart = types.BoolPointerValue(var2.LogStart)
		var3.LogEnd = types.BoolPointerValue(var2.LogEnd)
		if var2.ProfileSettings != nil {
			var3.ProfileSettings = &SecurityPolicyRulesResourceProfileSettingsObject{}
			var26, var27 := types.ListValueFrom(ctx, types.StringType, var2.ProfileSettings.Groups)
			var3.ProfileSettings.Groups = var26
			resp.Diagnostics.Append(var27.Errors()...)
			if var2.ProfileSettings.Profiles != nil {
				var3.ProfileSettings.Profiles = &SecurityPolicyRulesResourceProfilesObject{}
				var28, var29 := types.ListValueFrom(ctx, types.StringType, var2.ProfileSettings.Profiles.UrlFilteringProfiles)
				var3.ProfileSettings.Profiles.UrlFilteringProfiles = var28
				resp.Diagnostics.Append(var29.Errors()...)
				var30, var31 := types.ListValueFrom(ctx, types.StringType, var2.ProfileSettings.Profiles.DataFilteringProfiles)
				var3.ProfileSettings.Profiles.DataFilteringProfiles = var30
				resp.Diagnostics.Append(var31.Errors()...)
				var32, var33 := types.ListValueFrom(ctx, types.StringType, var2.ProfileSettings.Profiles.FileBlockingProfiles)
				var3.ProfileSettings.Profiles.FileBlockingProfiles = var32
				resp.Diagnostics.Append(var33.Errors()...)
				var34, var35 := types.ListValueFrom(ctx, types.StringType, var2.ProfileSettings.Profiles.WildfireAnalysisProfiles)
				var3.ProfileSettings.Profiles.WildfireAnalysisProfiles = var34
				resp.Diagnostics.Append(var35.Errors()...)
				var36, var37 := types.ListValueFrom(ctx, types.StringType, var2.ProfileSettings.Profiles.AntiVirusProfiles)
				var3.ProfileSettings.Profiles.AntiVirusProfiles = var36
				resp.Diagnostics.Append(var37.Errors()...)
				var38, var39 := types.ListValueFrom(ctx, types.StringType, var2.ProfileSettings.Profiles.AntiSpywareProfiles)
				var3.ProfileSettings.Profiles.AntiSpywareProfiles = var38
				resp.Diagnostics.Append(var39.Errors()...)
				var40, var41 := types.ListValueFrom(ctx, types.StringType, var2.ProfileSettings.Profiles.VulnerabilityProfiles)
				var3.ProfileSettings.Profiles.VulnerabilityProfiles = var40
				resp.Diagnostics.Append(var41.Errors()...)
			}
		}
		if var2.Qos != nil {
			var3.Qos = &SecurityPolicyRulesResourceQosObject{}
			var3.Qos.IpDscp = types.StringPointerValue(var2.Qos.IpDscp)
			var3.Qos.IpPrecedence = types.StringPointerValue(var2.Qos.IpPrecedence)
			if var2.Qos.FollowClientToServerFlow == nil {
				var3.Qos.FollowClientToServerFlow = types.BoolValue(false)
			} else {
				var3.Qos.FollowClientToServerFlow = types.BoolValue(true)
			}
		}
		// TODO: Targets
		var3.NegateTarget = types.BoolPointerValue(var2.NegateTarget)
		var3.DisableInspect = types.BoolPointerValue(var2.DisableInspect)

		rules = append(rules, var3)
	}
	state.Rules = rules
	state.Ordered = types.BoolValue(true)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *SecurityPolicyRulesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var err error
	var state SecurityPolicyRulesResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the tfid info.
	var tfid SecurityPolicyRulesTfid
	if err = DecodeLocation(state.Tfid.ValueString(), &tfid); err != nil {
		resp.Diagnostics.AddError("Error parsing tfid", err.Error())
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource read", map[string]any{
		"resource_name": "panos_security_policy_rules",
		"function":      "Delete",
	})

	// Verify mode.
	if r.client.Hostname == "" {
		resp.Diagnostics.AddError("Invalid mode error", InspectionModeError)
		return
	}

	uuids := make([]string, 0, len(tfid.Rules))
	for _, rInfo := range tfid.Rules {
		uuids = append(uuids, rInfo.Uuid)
	}

	// Create the service.
	svc := security.NewService(r.client)

	// Timeout handling.
	//ctx, cancel := context.WithTimeout(ctx, GetTimeout(state.Timeouts.Delete))
	//defer cancel()

	// Perform the operation.
	if err := svc.DeleteById(ctx, tfid.Location, uuids...); err != nil {
		resp.Diagnostics.AddError("Error in delete", err.Error())
	}
}

func (r *SecurityPolicyRulesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("tfid"), req, resp)
}
