package provider

import (
	"context"
	//"fmt"
	//"regexp"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/policies/rules/security"

	//"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	//"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	//"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	//"github.com/hashicorp/terraform-plugin-framework/path"
	//"github.com/hashicorp/terraform-plugin-framework/resource"
	//rsschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	//"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	//"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	//"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	//"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	//"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	//"github.com/hashicorp/terraform-plugin-log/tflog"
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
