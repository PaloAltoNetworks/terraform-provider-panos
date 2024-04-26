package provider

import (
	"context"
	"fmt"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/objects/service"
	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/xmlapi"

	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	rsschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Resource.
var (
	_ resource.Resource                = &ServiceObjectsResource{}
	_ resource.ResourceWithConfigure   = &ServiceObjectsResource{}
	_ resource.ResourceWithImportState = &ServiceObjectsResource{}
)

func NewServiceObjectsResource() resource.Resource {
	return &ServiceObjectsResource{}
}

type ServiceObjectsResource struct {
	client *pango.XmlApiClient
}

type ServiceObjectsTfid struct {
	Names    []string         `json:"names"`
	Location service.Location `json:"location"`
}

func (o *ServiceObjectsTfid) IsValid() error {
	var err error
	if len(o.Names) == 0 {
		return fmt.Errorf("No objects present")
	}

	names := make(map[string]bool)
	for _, name := range o.Names {
		if names[name] {
			return fmt.Errorf("Object %q present multiple times", name)
		}
		names[name] = true
	}

	if err = o.Location.IsValid(); err != nil {
		return err
	}

	return nil
}

type ServiceObjectsResourceModel struct {
	//Timeouts crudTimeouts `tfsdk:"timeouts"`
	Tfid     types.String                   `tfsdk:"tfid"`
	Location ServiceObjectsResourceLocation `tfsdk:"location"`
	Objects  []ServiceObjectsEntry          `tfsdk:"objects"`
}

type ServiceObjectsEntry struct {
	Name        types.String                 `tfsdk:"name"`
	Description types.String                 `tfsdk:"description"`
	Tags        types.List                   `tfsdk:"tags"`
	Protocol    ServiceObjectsProtocolObject `tfsdk:"protocol"`
}

type ServiceObjectsProtocolObject struct {
	Tcp  types.Object `tfsdk:"tcp"`
	Udp  types.Object `tfsdk:"udp"`
	Sctp types.Object `tfsdk:"sctp"`
}

func (o ServiceObjectsProtocolObject) Types() map[string]attr.Type {
	return map[string]attr.Type{
		"tcp": types.ObjectType{
			AttrTypes: ServiceObjectsTcpObject{}.Types(),
		},
		"udp": types.ObjectType{
			AttrTypes: ServiceObjectsUdpObject{}.Types(),
		},
		"sctp": types.ObjectType{
			AttrTypes: ServiceObjectsSctpObject{}.Types(),
		},
	}
}

type ServiceObjectsTcpObject struct {
	DestinationPort types.String `tfsdk:"destination_port"`
	SourcePort      types.String `tfsdk:"source_port"`
	Override        types.Object `tfsdk:"override"`
}

func (o ServiceObjectsTcpObject) Types() map[string]attr.Type {
	return map[string]attr.Type{
		"destination_port": types.StringType,
		"source_port":      types.StringType,
		"override": types.ObjectType{
			AttrTypes: ServiceObjectsTcpOverrideObject{}.Types(),
		},
	}
}

type ServiceObjectsTcpOverrideObject struct {
	No  types.Bool   `tfsdk:"no"`
	Yes types.Object `tfsdk:"yes"`
}

func (o ServiceObjectsTcpOverrideObject) Types() map[string]attr.Type {
	return map[string]attr.Type{
		"no": types.BoolType,
		"yes": types.ObjectType{
			AttrTypes: ServiceObjectsYesTcpOverrideObject{}.Types(),
		},
	}
}

type ServiceObjectsYesTcpOverrideObject struct {
	Timeout           types.Int64 `tfsdk:"timeout"`
	HalfClosedTimeout types.Int64 `tfsdk:"half_closed_timeout"`
	TimeWaitTimeout   types.Int64 `tfsdk:"time_wait_timeout"`
}

func (o ServiceObjectsYesTcpOverrideObject) Types() map[string]attr.Type {
	return map[string]attr.Type{
		"timeout":             types.Int64Type,
		"half_closed_timeout": types.Int64Type,
		"time_wait_timeout":   types.Int64Type,
	}
}

type ServiceObjectsUdpObject struct {
	DestinationPort types.String `tfsdk:"destination_port"`
	SourcePort      types.String `tfsdk:"source_port"`
	Override        types.Object `tfsdk:"override"`
}

func (o ServiceObjectsUdpObject) Types() map[string]attr.Type {
	return map[string]attr.Type{
		"destination_port": types.StringType,
		"source_port":      types.StringType,
		"override": types.ObjectType{
			AttrTypes: ServiceObjectsUdpOverrideObject{}.Types(),
		},
	}
}

type ServiceObjectsUdpOverrideObject struct {
	No  types.Bool   `tfsdk:"no"`
	Yes types.Object `tfsdk:"yes"`
}

func (o ServiceObjectsUdpOverrideObject) Types() map[string]attr.Type {
	return map[string]attr.Type{
		"no": types.BoolType,
		"yes": types.ObjectType{
			AttrTypes: ServiceObjectsYesUdpOverrideObject{}.Types(),
		},
	}
}

type ServiceObjectsYesUdpOverrideObject struct {
	Timeout types.Int64 `tfsdk:"timeout"`
}

func (o ServiceObjectsYesUdpOverrideObject) Types() map[string]attr.Type {
	return map[string]attr.Type{
		"timeout": types.Int64Type,
	}
}

type ServiceObjectsSctpObject struct {
	DestinationPort types.String `tfsdk:"destination_port"`
	SourcePort      types.String `tfsdk:"source_port"`
}

func (o ServiceObjectsSctpObject) Types() map[string]attr.Type {
	return map[string]attr.Type{
		"destination_port": types.StringType,
		"source_port":      types.StringType,
	}
}

type ServiceObjectsResourceLocation struct {
	Shared       types.Bool                                 `tfsdk:"shared"`
	Vsys         *ServiceObjectsResourceVsysLocation        `tfsdk:"vsys"`
	DeviceGroup  *ServiceObjectsResourceDeviceGroupLocation `tfsdk:"device_group"`
	FromPanorama types.Bool                                 `tfsdk:"from_panorama"`
}

type ServiceObjectsResourceVsysLocation struct {
	NgfwDevice types.String `tfsdk:"ngfw_device"`
	Name       types.String `tfsdk:"name"`
}

type ServiceObjectsResourceDeviceGroupLocation struct {
	PanoramaDevice types.String `tfsdk:"panorama_device"`
	Name           types.String `tfsdk:"name"`
}

func (r *ServiceObjectsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_objects"
}

func (r *ServiceObjectsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = rsschema.Schema{
		Description: "Manges a group of service objects.",

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
								path.MatchRoot("location").AtName("device_group"),
								path.MatchRoot("location").AtName("shared"),
								path.MatchRoot("location").AtName("from_panorama"),
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
						},
					},
					"shared": rsschema.BoolAttribute{
						Description: "(NGFW and Panorama) Located in shared.",
						Optional:    true,
					},
					"from_panorama": rsschema.BoolAttribute{
						Description: "(NGFW) Pushed from Panorama.",
						Optional:    true,
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
			"objects": rsschema.ListNestedAttribute{
				Description: "The list of service objects.",
				Required:    true,
				NestedObject: rsschema.NestedAttributeObject{
					Attributes: map[string]rsschema.Attribute{
						"name": rsschema.StringAttribute{
							Description: "Alphanumeric string [ 0-9a-zA-Z._-].",
							Required:    true,
						},
						"description": rsschema.StringAttribute{
							Description: "The description.",
							Optional:    true,
						},
						"tags": rsschema.ListAttribute{
							Description: "Tags for address object.",
							Optional:    true,
							ElementType: types.StringType,
						},
						"protocol": rsschema.SingleNestedAttribute{
							Description: "The protocol specification.",
							Required:    true,
							Attributes: map[string]rsschema.Attribute{
								"tcp": rsschema.SingleNestedAttribute{
									Description: "TCP protocol spec.",
									Optional:    true,
									Attributes: map[string]rsschema.Attribute{
										"destination_port": rsschema.StringAttribute{
											Description: "The destination port.",
											Required:    true,
										},
										"source_port": rsschema.StringAttribute{
											Description: "The source port.",
											Optional:    true,
										},
										"override": rsschema.SingleNestedAttribute{
											Description: "Override spec.",
											Optional:    true,
											Attributes: map[string]rsschema.Attribute{
												"no": rsschema.BoolAttribute{
													Description: "No override.",
													Optional:    true,
													Computed:    true,
													Default:     booldefault.StaticBool(false),
												},
												"yes": rsschema.SingleNestedAttribute{
													Description: "Enable TCP override.",
													Optional:    true,
													Attributes: map[string]rsschema.Attribute{
														"timeout": rsschema.Int64Attribute{
															Description: "TCP timeout.",
															Optional:    true,
														},
														"half_closed_timeout": rsschema.Int64Attribute{
															Description: "Half closed timeout.",
															Optional:    true,
														},
														"time_wait_timeout": rsschema.Int64Attribute{
															Description: "Time wait timeout.",
															Optional:    true,
														},
													},
												},
											},
										},
									},
								},
								"udp": rsschema.SingleNestedAttribute{
									Description: "UDP protocol spec.",
									Optional:    true,
									Attributes: map[string]rsschema.Attribute{
										"destination_port": rsschema.StringAttribute{
											Description: "The destination port.",
											Required:    true,
										},
										"source_port": rsschema.StringAttribute{
											Description: "The source port.",
											Optional:    true,
										},
										"override": rsschema.SingleNestedAttribute{
											Description: "Override spec.",
											Optional:    true,
											Attributes: map[string]rsschema.Attribute{
												"no": rsschema.BoolAttribute{
													Description: "No override.",
													Optional:    true,
													Computed:    true,
													Default:     booldefault.StaticBool(false),
												},
												"yes": rsschema.SingleNestedAttribute{
													Description: "Enable TCP override.",
													Optional:    true,
													Attributes: map[string]rsschema.Attribute{
														"timeout": rsschema.Int64Attribute{
															Description: "TCP timeout.",
															Optional:    true,
														},
													},
												},
											},
										},
									},
								},
								"sctp": rsschema.SingleNestedAttribute{
									Description: "SCTP protocol spec.",
									Optional:    true,
									Attributes: map[string]rsschema.Attribute{
										"destination_port": rsschema.StringAttribute{
											Description: "The destination port.",
											Required:    true,
										},
										"source_port": rsschema.StringAttribute{
											Description: "The source port.",
											Optional:    true,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r *ServiceObjectsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*pango.XmlApiClient)
}

func (r *ServiceObjectsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var err error
	var listing []service.Entry
	var tfid ServiceObjectsTfid
	var state ServiceObjectsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource create", map[string]any{
		"resource_name": "panos_service_objects",
		"function":      "Create",
	})

	// Verify mode.
	if r.client.Hostname == "" {
		resp.Diagnostics.AddError("Invalid mode error", InspectionModeError)
		return
	}

	// Create the service.
	svc := service.NewService(r.client)

	// Determine the location.
	if !state.Location.Shared.IsNull() && state.Location.Shared.ValueBool() {
		tfid.Location.Shared = true
	}
	if state.Location.Vsys != nil {
		tfid.Location.Vsys = &service.VsysLocation{
			NgfwDevice: state.Location.Vsys.NgfwDevice.ValueString(),
			Name:       state.Location.Vsys.Name.ValueString(),
		}
	}
	if state.Location.DeviceGroup != nil {
		tfid.Location.DeviceGroup = &service.DeviceGroupLocation{
			PanoramaDevice: state.Location.DeviceGroup.PanoramaDevice.ValueString(),
			Name:           state.Location.DeviceGroup.Name.ValueString(),
		}
	}
	if !state.Location.FromPanorama.IsNull() && state.Location.FromPanorama.ValueBool() {
		tfid.Location.FromPanorama = true
	}
	if err = tfid.Location.IsValid(); err != nil {
		resp.Diagnostics.AddError("Invalid path", err.Error())
		return
	}

	// Load the desired config.
	isNewName := make(map[string]bool)
	tfid.Names = make([]string, 0, len(state.Objects))
	entries := make([]service.Entry, 0, len(state.Objects))
	for _, var0 := range state.Objects {
		isNewName[var0.Name.ValueString()] = true
		tfid.Names = append(tfid.Names, var0.Name.ValueString())
		var1 := service.Entry{Name: var0.Name.ValueString()}
		var1.Description = var0.Description.ValueStringPointer()
		resp.Diagnostics.Append(var0.Tags.ElementsAs(ctx, &var1.Tags, false)...)
		/*
		   if !var0.Protocol.IsNull() {
		       var var2 ServiceObjectsProtocolObject
		       resp.Diagnostics.Append(var0.Protocol.As(ctx, &var2, basetypes.ObjectAsOptions{})...)
		*/
		var2 := var0.Protocol

		if !var2.Tcp.IsNull() {
			var1.Protocol.Tcp = &service.TcpObject{}
			var var3 ServiceObjectsTcpObject
			resp.Diagnostics.Append(var2.Tcp.As(ctx, &var3, basetypes.ObjectAsOptions{})...)
			var1.Protocol.Tcp.DestinationPort = var3.DestinationPort.ValueString()
			var1.Protocol.Tcp.SourcePort = var3.SourcePort.ValueStringPointer()
			if !var3.Override.IsNull() {
				var1.Protocol.Tcp.Override = &service.TcpOverrideObject{}
				var var4 ServiceObjectsTcpOverrideObject
				resp.Diagnostics.Append(var3.Override.As(ctx, &var4, basetypes.ObjectAsOptions{})...)
				var1.Protocol.Tcp.Override.No = var4.No.ValueBool()
				if !var4.Yes.IsNull() {
					var1.Protocol.Tcp.Override.Yes = &service.YesTcpOverrideObject{}
					var var5 ServiceObjectsYesTcpOverrideObject
					resp.Diagnostics.Append(var4.Yes.As(ctx, &var5, basetypes.ObjectAsOptions{})...)

					var1.Protocol.Tcp.Override.Yes.Timeout = var5.Timeout.ValueInt64Pointer()
					var1.Protocol.Tcp.Override.Yes.HalfClosedTimeout = var5.HalfClosedTimeout.ValueInt64Pointer()
					var1.Protocol.Tcp.Override.Yes.TimeWaitTimeout = var5.TimeWaitTimeout.ValueInt64Pointer()
				}
			}
		}
		if !var2.Udp.IsNull() {
			var1.Protocol.Udp = &service.UdpObject{}
			var var3 ServiceObjectsUdpObject
			resp.Diagnostics.Append(var2.Udp.As(ctx, &var3, basetypes.ObjectAsOptions{})...)
			var1.Protocol.Udp.DestinationPort = var3.DestinationPort.ValueString()
			var1.Protocol.Udp.SourcePort = var3.SourcePort.ValueStringPointer()
			if !var3.Override.IsNull() {
				var1.Protocol.Udp.Override = &service.UdpOverrideObject{}
				var var4 ServiceObjectsUdpOverrideObject
				resp.Diagnostics.Append(var3.Override.As(ctx, &var4, basetypes.ObjectAsOptions{})...)
				var1.Protocol.Udp.Override.No = var4.No.ValueBool()
				if !var4.Yes.IsNull() {
					var1.Protocol.Udp.Override.Yes = &service.YesUdpOverrideObject{}
					var var5 ServiceObjectsYesUdpOverrideObject
					resp.Diagnostics.Append(var4.Yes.As(ctx, &var5, basetypes.ObjectAsOptions{})...)

					var1.Protocol.Udp.Override.Yes.Timeout = var5.Timeout.ValueInt64Pointer()
				}
			}
		}
		if !var2.Sctp.IsNull() {
			var1.Protocol.Sctp = &service.SctpObject{}
			var var3 ServiceObjectsSctpObject
			resp.Diagnostics.Append(var2.Sctp.As(ctx, &var3, basetypes.ObjectAsOptions{})...)
			var1.Protocol.Sctp.DestinationPort = var3.DestinationPort.ValueString()
			var1.Protocol.Sctp.SourcePort = var3.SourcePort.ValueStringPointer()
		}
		//}

		entries = append(entries, var1)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Timeout handling.
	//ctx, cancel := context.WithTimeout(ctx, GetTimeout(state.Timeouts.Create))
	//defer cancel()

	// Get the current list of objects.
	listing, err = svc.List(ctx, tfid.Location, "get", "", "")
	if err != nil {
		resp.Diagnostics.AddError("Error during Create's refresh", err.Error())
		return
	}

	// Verify no objects are already present.
	for _, x := range listing {
		if isNewName[x.Name] {
			resp.Diagnostics.AddError("Object to be created already exists", x.Name)
		}
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Prepare the multi-config.
	vn := r.client.Versioning()
	updates := xmlapi.NewMultiConfig(len(entries))
	specifier, _, err := service.Versioning(vn)
	if err != nil {
		resp.Diagnostics.AddError("Error getting specifier", err.Error())
		return
	}
	for _, entry := range entries {
		path, err := tfid.Location.Xpath(vn, entry.Name)
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

	// Create the objects.
	if _, _, _, err = r.client.MultiConfig(ctx, updates, false, nil); err != nil {
		resp.Diagnostics.AddError("Error in create", err.Error())
		return
	}

	// Retrieve the list of objects again.
	listing, err = svc.List(ctx, tfid.Location, "get", "", "")
	if err != nil {
		_ = svc.Delete(ctx, tfid.Location, tfid.Names...)
		resp.Diagnostics.AddError("Error during Create's refresh", err.Error())
		return
	}

	ans := make([]service.Entry, 0, len(entries))
	for _, live := range listing {
		if isNewName[live.Name] {
			ans = append(ans, live)
		}
	}

	// Tfid handling.
	tfidstr, err := EncodeLocation(&tfid)
	if err != nil {
		_ = svc.Delete(ctx, tfid.Location, tfid.Names...)
		resp.Diagnostics.AddError("Error creating tfid", err.Error())
		return
	}

	// Save the state.
	state.Tfid = types.StringValue(tfidstr)
	objs := make([]ServiceObjectsEntry, 0, len(ans))
	for _, var50 := range ans {
		var51 := ServiceObjectsEntry{}
		var51.Name = types.StringValue(var50.Name)
		var51.Description = types.StringPointerValue(var50.Description)

		var52, var53 := types.ListValueFrom(ctx, types.StringType, var50.Tags)
		var51.Tags = var52
		resp.Diagnostics.Append(var53.Errors()...)

		var var54 ServiceObjectsProtocolObject
		if var50.Protocol.Tcp == nil {
			var54.Tcp = types.ObjectNull(ServiceObjectsTcpObject{}.Types())
		} else {
			var var55 ServiceObjectsTcpObject
			var55.DestinationPort = types.StringValue(var50.Protocol.Tcp.DestinationPort)
			var55.SourcePort = types.StringPointerValue(var50.Protocol.Tcp.SourcePort)
			if var50.Protocol.Tcp.Override == nil {
				var55.Override = types.ObjectNull(ServiceObjectsTcpOverrideObject{}.Types())
			} else {
				var var56 ServiceObjectsTcpOverrideObject
				var56.No = types.BoolValue(var50.Protocol.Tcp.Override.No)
				if var50.Protocol.Tcp.Override.Yes == nil {
					var56.Yes = types.ObjectNull(ServiceObjectsYesTcpOverrideObject{}.Types())
				} else {
					var var57 ServiceObjectsYesTcpOverrideObject
					var57.Timeout = types.Int64PointerValue(var50.Protocol.Tcp.Override.Yes.Timeout)
					var57.HalfClosedTimeout = types.Int64PointerValue(var50.Protocol.Tcp.Override.Yes.HalfClosedTimeout)
					var57.TimeWaitTimeout = types.Int64PointerValue(var50.Protocol.Tcp.Override.Yes.TimeWaitTimeout)
					var58, var59 := types.ObjectValueFrom(ctx, var57.Types(), var57)
					var56.Yes = var58
					resp.Diagnostics.Append(var59...)
				}
				var60, var61 := types.ObjectValueFrom(ctx, var56.Types(), var56)
				var55.Override = var60
				resp.Diagnostics.Append(var61...)
			}
			var62, var63 := types.ObjectValueFrom(ctx, var55.Types(), var55)
			var54.Tcp = var62
			resp.Diagnostics.Append(var63...)
		}

		if var50.Protocol.Udp == nil {
			var54.Udp = types.ObjectNull(ServiceObjectsUdpObject{}.Types())
		} else {
			var var55 ServiceObjectsUdpObject
			var55.DestinationPort = types.StringValue(var50.Protocol.Udp.DestinationPort)
			var55.SourcePort = types.StringPointerValue(var50.Protocol.Udp.SourcePort)
			if var50.Protocol.Udp.Override == nil {
				var55.Override = types.ObjectNull(ServiceObjectsUdpOverrideObject{}.Types())
			} else {
				var var56 ServiceObjectsUdpOverrideObject
				var56.No = types.BoolValue(var50.Protocol.Udp.Override.No)
				if var50.Protocol.Udp.Override.Yes == nil {
					var56.Yes = types.ObjectNull(ServiceObjectsYesUdpOverrideObject{}.Types())
				} else {
					var var57 ServiceObjectsYesUdpOverrideObject
					var57.Timeout = types.Int64PointerValue(var50.Protocol.Udp.Override.Yes.Timeout)
					var58, var59 := types.ObjectValueFrom(ctx, var57.Types(), var57)
					var56.Yes = var58
					resp.Diagnostics.Append(var59...)
				}
				var60, var61 := types.ObjectValueFrom(ctx, var56.Types(), var56)
				var55.Override = var60
				resp.Diagnostics.Append(var61...)
			}
			var62, var63 := types.ObjectValueFrom(ctx, var55.Types(), var55)
			var54.Udp = var62
			resp.Diagnostics.Append(var63...)
		}

		if var50.Protocol.Sctp == nil {
			var54.Sctp = types.ObjectNull(ServiceObjectsSctpObject{}.Types())
		} else {
			var var55 ServiceObjectsSctpObject
			var55.DestinationPort = types.StringValue(var50.Protocol.Sctp.DestinationPort)
			var55.SourcePort = types.StringPointerValue(var50.Protocol.Sctp.SourcePort)

			var62, var63 := types.ObjectValueFrom(ctx, var55.Types(), var55)
			var54.Udp = var62
			resp.Diagnostics.Append(var63...)
		}

		var51.Protocol = var54
		objs = append(objs, var51)
	}

	state.Objects = objs

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ServiceObjectsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var err error
	var savestate, state ServiceObjectsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &savestate)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the tfid info.
	var tfid ServiceObjectsTfid
	if err = DecodeLocation(savestate.Tfid.ValueString(), &tfid); err != nil {
		resp.Diagnostics.AddError("Error parsing tfid", err.Error())
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource read", map[string]any{
		"resource_name": "panos_service_objects",
		"function":      "Read",
	})

	// Verify mode.
	if r.client.Hostname == "" {
		resp.Diagnostics.AddError("Invalid mode error", InspectionModeError)
		return
	}

	// Create the service.
	svc := service.NewService(r.client)

	// Timeout handling.
	//ctx, cancel := context.WithTimeout(ctx, GetTimeout(savestate.Timeouts.Read))
	//defer cancel()

	// Perform a list to get all objects.
	listing, err := svc.List(ctx, tfid.Location, "get", "", "")
	if err != nil {
		resp.Diagnostics.AddError("Error in listing", err.Error())
		return
	}

	// If there are no objects, we can remove the state, we need a full redeploy.
	if len(listing) == 0 {
		resp.State.RemoveResource(ctx)
		return
	}

	nameMap := make(map[string]int, len(listing))
	for index, live := range listing {
		nameMap[live.Name] = index
	}

	// Find the objects.
	objs := make([]ServiceObjectsEntry, 0, len(tfid.Names))
	for _, name := range tfid.Names {
		oid, ok := nameMap[name]
		if !ok {
			continue
		}

		var50 := listing[oid]
		var51 := ServiceObjectsEntry{}
		var51.Name = types.StringValue(var50.Name)
		var51.Description = types.StringPointerValue(var50.Description)

		var52, var53 := types.ListValueFrom(ctx, types.StringType, var50.Tags)
		var51.Tags = var52
		resp.Diagnostics.Append(var53.Errors()...)

		var var54 ServiceObjectsProtocolObject
		if var50.Protocol.Tcp == nil {
			var54.Tcp = types.ObjectNull(ServiceObjectsTcpObject{}.Types())
		} else {
			var var55 ServiceObjectsTcpObject
			var55.DestinationPort = types.StringValue(var50.Protocol.Tcp.DestinationPort)
			var55.SourcePort = types.StringPointerValue(var50.Protocol.Tcp.SourcePort)
			if var50.Protocol.Tcp.Override == nil {
				var55.Override = types.ObjectNull(ServiceObjectsTcpOverrideObject{}.Types())
			} else {
				var var56 ServiceObjectsTcpOverrideObject
				var56.No = types.BoolValue(var50.Protocol.Tcp.Override.No)
				if var50.Protocol.Tcp.Override.Yes == nil {
					var56.Yes = types.ObjectNull(ServiceObjectsYesTcpOverrideObject{}.Types())
				} else {
					var var57 ServiceObjectsYesTcpOverrideObject
					var57.Timeout = types.Int64PointerValue(var50.Protocol.Tcp.Override.Yes.Timeout)
					var57.HalfClosedTimeout = types.Int64PointerValue(var50.Protocol.Tcp.Override.Yes.HalfClosedTimeout)
					var57.TimeWaitTimeout = types.Int64PointerValue(var50.Protocol.Tcp.Override.Yes.TimeWaitTimeout)
					var58, var59 := types.ObjectValueFrom(ctx, var57.Types(), var57)
					var56.Yes = var58
					resp.Diagnostics.Append(var59...)
				}
				var60, var61 := types.ObjectValueFrom(ctx, var56.Types(), var56)
				var55.Override = var60
				resp.Diagnostics.Append(var61...)
			}
			var62, var63 := types.ObjectValueFrom(ctx, var55.Types(), var55)
			var54.Tcp = var62
			resp.Diagnostics.Append(var63...)
		}

		if var50.Protocol.Udp == nil {
			var54.Udp = types.ObjectNull(ServiceObjectsUdpObject{}.Types())
		} else {
			var var55 ServiceObjectsUdpObject
			var55.DestinationPort = types.StringValue(var50.Protocol.Udp.DestinationPort)
			var55.SourcePort = types.StringPointerValue(var50.Protocol.Udp.SourcePort)
			if var50.Protocol.Udp.Override == nil {
				var55.Override = types.ObjectNull(ServiceObjectsUdpOverrideObject{}.Types())
			} else {
				var var56 ServiceObjectsUdpOverrideObject
				var56.No = types.BoolValue(var50.Protocol.Udp.Override.No)
				if var50.Protocol.Udp.Override.Yes == nil {
					var56.Yes = types.ObjectNull(ServiceObjectsYesUdpOverrideObject{}.Types())
				} else {
					var var57 ServiceObjectsYesUdpOverrideObject
					var57.Timeout = types.Int64PointerValue(var50.Protocol.Udp.Override.Yes.Timeout)
					var58, var59 := types.ObjectValueFrom(ctx, var57.Types(), var57)
					var56.Yes = var58
					resp.Diagnostics.Append(var59...)
				}
				var60, var61 := types.ObjectValueFrom(ctx, var56.Types(), var56)
				var55.Override = var60
				resp.Diagnostics.Append(var61...)
			}
			var62, var63 := types.ObjectValueFrom(ctx, var55.Types(), var55)
			var54.Udp = var62
			resp.Diagnostics.Append(var63...)
		}

		if var50.Protocol.Sctp == nil {
			var54.Sctp = types.ObjectNull(ServiceObjectsSctpObject{}.Types())
		} else {
			var var55 ServiceObjectsSctpObject
			var55.DestinationPort = types.StringValue(var50.Protocol.Sctp.DestinationPort)
			var55.SourcePort = types.StringPointerValue(var50.Protocol.Sctp.SourcePort)

			var62, var63 := types.ObjectValueFrom(ctx, var55.Types(), var55)
			var54.Udp = var62
			resp.Diagnostics.Append(var63...)
		}

		var51.Protocol = var54
		objs = append(objs, var51)
	}

	// If there are no rules, we can remove the state, we need a full redeploy.
	if len(objs) == 0 {
		resp.State.RemoveResource(ctx)
		return
	}

	// Save the location.
	if tfid.Location.Shared {
		state.Location.Shared = types.BoolValue(true)
	}
	if tfid.Location.Vsys != nil {
		state.Location.Vsys = &ServiceObjectsResourceVsysLocation{
			NgfwDevice: types.StringValue(tfid.Location.Vsys.NgfwDevice),
			Name:       types.StringValue(tfid.Location.Vsys.Name),
		}
	}
	if tfid.Location.DeviceGroup != nil {
		state.Location.DeviceGroup = &ServiceObjectsResourceDeviceGroupLocation{
			PanoramaDevice: types.StringValue(tfid.Location.DeviceGroup.PanoramaDevice),
			Name:           types.StringValue(tfid.Location.DeviceGroup.Name),
		}
	}
	if tfid.Location.FromPanorama {
		state.Location.FromPanorama = types.BoolValue(true)
	}

	// Save state.
	state.Tfid = savestate.Tfid
	state.Objects = objs

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ServiceObjectsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var err error
	var updates *xmlapi.MultiConfig
	var plan, state ServiceObjectsResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the tfid info.
	var tfid, newtfid ServiceObjectsTfid
	if err = DecodeLocation(state.Tfid.ValueString(), &tfid); err != nil {
		resp.Diagnostics.AddError("Error parsing tfid", err.Error())
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource read", map[string]any{
		"resource_name": "panos_service_objects",
		"function":      "Update",
	})

	// Determine the location.
	if !plan.Location.Shared.IsNull() && plan.Location.Shared.ValueBool() {
		newtfid.Location.Shared = true
	}
	if plan.Location.Vsys != nil {
		newtfid.Location.Vsys = &service.VsysLocation{
			NgfwDevice: plan.Location.Vsys.NgfwDevice.ValueString(),
			Name:       plan.Location.Vsys.Name.ValueString(),
		}
	}
	if plan.Location.DeviceGroup != nil {
		newtfid.Location.DeviceGroup = &service.DeviceGroupLocation{
			PanoramaDevice: plan.Location.DeviceGroup.PanoramaDevice.ValueString(),
			Name:           plan.Location.DeviceGroup.Name.ValueString(),
		}
	}
	if !plan.Location.FromPanorama.IsNull() && plan.Location.FromPanorama.ValueBool() {
		newtfid.Location.FromPanorama = true
	}
	if err = newtfid.Location.IsValid(); err != nil {
		resp.Diagnostics.AddError("Invalid path", err.Error())
		return
	}

	// Verify mode.
	if r.client.Hostname == "" {
		resp.Diagnostics.AddError("Invalid mode error", InspectionModeError)
		return
	}

	// Create the service.
	svc := service.NewService(r.client)

	// Prepare to handle versioning.
	vn := r.client.Versioning()
	specifier, _, err := service.Versioning(vn)
	if err != nil {
		resp.Diagnostics.AddError("Error getting specifier", err.Error())
		return
	}

	// Build up the old object names we previously made.
	ownedNames := make(map[string]bool)
	for _, name := range tfid.Names {
		ownedNames[name] = true
	}

	// Load the desired config.
	desiredNames := make(map[string]bool)
	newtfid.Names = make([]string, 0, len(plan.Objects))
	entries := make([]service.Entry, 0, len(plan.Objects))
	for _, var0 := range plan.Objects {
		desiredNames[var0.Name.ValueString()] = true
		newtfid.Names = append(newtfid.Names, var0.Name.ValueString())
		var1 := service.Entry{Name: var0.Name.ValueString()}
		var1.Description = var0.Description.ValueStringPointer()
		resp.Diagnostics.Append(var0.Tags.ElementsAs(ctx, &var1.Tags, false)...)
		/*
		   if !var0.Protocol.IsNull() {
		       var var2 ServiceObjectsProtocolObject
		       resp.Diagnostics.Append(var0.Protocol.As(ctx, &var2, basetypes.ObjectAsOptions{})...)
		*/
		var2 := var0.Protocol

		if !var2.Tcp.IsNull() {
			var1.Protocol.Tcp = &service.TcpObject{}
			var var3 ServiceObjectsTcpObject
			resp.Diagnostics.Append(var2.Tcp.As(ctx, &var3, basetypes.ObjectAsOptions{})...)
			var1.Protocol.Tcp.DestinationPort = var3.DestinationPort.ValueString()
			var1.Protocol.Tcp.SourcePort = var3.SourcePort.ValueStringPointer()
			if !var3.Override.IsNull() {
				var1.Protocol.Tcp.Override = &service.TcpOverrideObject{}
				var var4 ServiceObjectsTcpOverrideObject
				resp.Diagnostics.Append(var3.Override.As(ctx, &var4, basetypes.ObjectAsOptions{})...)
				var1.Protocol.Tcp.Override.No = var4.No.ValueBool()
				if !var4.Yes.IsNull() {
					var1.Protocol.Tcp.Override.Yes = &service.YesTcpOverrideObject{}
					var var5 ServiceObjectsYesTcpOverrideObject
					resp.Diagnostics.Append(var4.Yes.As(ctx, &var5, basetypes.ObjectAsOptions{})...)

					var1.Protocol.Tcp.Override.Yes.Timeout = var5.Timeout.ValueInt64Pointer()
					var1.Protocol.Tcp.Override.Yes.HalfClosedTimeout = var5.HalfClosedTimeout.ValueInt64Pointer()
					var1.Protocol.Tcp.Override.Yes.TimeWaitTimeout = var5.TimeWaitTimeout.ValueInt64Pointer()
				}
			}
		}
		if !var2.Udp.IsNull() {
			var1.Protocol.Udp = &service.UdpObject{}
			var var3 ServiceObjectsUdpObject
			resp.Diagnostics.Append(var2.Udp.As(ctx, &var3, basetypes.ObjectAsOptions{})...)
			var1.Protocol.Udp.DestinationPort = var3.DestinationPort.ValueString()
			var1.Protocol.Udp.SourcePort = var3.SourcePort.ValueStringPointer()
			if !var3.Override.IsNull() {
				var1.Protocol.Udp.Override = &service.UdpOverrideObject{}
				var var4 ServiceObjectsUdpOverrideObject
				resp.Diagnostics.Append(var3.Override.As(ctx, &var4, basetypes.ObjectAsOptions{})...)
				var1.Protocol.Udp.Override.No = var4.No.ValueBool()
				if !var4.Yes.IsNull() {
					var1.Protocol.Udp.Override.Yes = &service.YesUdpOverrideObject{}
					var var5 ServiceObjectsYesUdpOverrideObject
					resp.Diagnostics.Append(var4.Yes.As(ctx, &var5, basetypes.ObjectAsOptions{})...)

					var1.Protocol.Udp.Override.Yes.Timeout = var5.Timeout.ValueInt64Pointer()
				}
			}
		}
		if !var2.Sctp.IsNull() {
			var1.Protocol.Sctp = &service.SctpObject{}
			var var3 ServiceObjectsSctpObject
			resp.Diagnostics.Append(var2.Sctp.As(ctx, &var3, basetypes.ObjectAsOptions{})...)
			var1.Protocol.Sctp.DestinationPort = var3.DestinationPort.ValueString()
			var1.Protocol.Sctp.SourcePort = var3.SourcePort.ValueStringPointer()
		}
		//}

		entries = append(entries, var1)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the list of all rules.
	listing, err := svc.List(ctx, newtfid.Location, "get", "", "")
	if err != nil {
		resp.Diagnostics.AddError("Error in refresh for update", err.Error())
		return
	}

	nameIsUsed := make(map[string]bool)
	nameMap := make(map[string]int, len(listing))
	for index, live := range listing {
		nameMap[live.Name] = index
	}

	updates = xmlapi.NewMultiConfig(len(entries) * 2)

	delayed := make([]service.Entry, 0, len(entries))
	for _, entry := range entries {
		if _, ok := nameMap[entry.Name]; !ownedNames[entry.Name] || !ok {
			delayed = append(delayed, entry)
		}

		nameIsUsed[entry.Name] = true
		listingIndex := nameMap[entry.Name]

		if !service.SpecMatches(&entry, &listing[listingIndex]) {
			path, err := newtfid.Location.Xpath(vn, entry.Name)
			if err != nil {
				resp.Diagnostics.AddError("Error creating update xpath", fmt.Sprintf("%s - %s", entry.Name, err))
				return
			}

			entry.CopyMiscFrom(&listing[listingIndex])
			elm, err := specifier(entry)
			if err != nil {
				resp.Diagnostics.AddError("Error specifying update", fmt.Sprintf("%s - %s", entry.Name, err))
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

	for _, entry := range delayed {
		// Verify the name we want isn't already present.
		if _, ok := nameMap[entry.Name]; ok {
			resp.Diagnostics.AddError("Object already exists", entry.Name)
			return
		}

		var oldName string
		if len(nameIsUsed) != len(tfid.Names) {
			for _, name := range tfid.Names {
				if !nameIsUsed[name] {
					nameIsUsed[name] = true
					if _, ok := nameMap[name]; ok {
						oldName = name
						break
					}
				}
			}
		}

		if oldName == "" {
			path, err := newtfid.Location.Xpath(vn, entry.Name)
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
			// Rename the old object to the desired name.
			path, err := newtfid.Location.Xpath(vn, oldName)
			if err != nil {
				resp.Diagnostics.AddError("Error creating repurpose rename xpath", fmt.Sprintf("%s - %s", oldName, err))
				return
			}

			updates.Add(&xmlapi.Config{
				Action:  "rename",
				Xpath:   util.AsXpath(path),
				NewName: entry.Name,
				Target:  r.client.GetTarget(),
			})

			listingIndex := nameMap[oldName]

			if service.SpecMatches(&entry, &listing[listingIndex]) {
				path, err := newtfid.Location.Xpath(vn, entry.Name)
				if err != nil {
					resp.Diagnostics.AddError("Error creating update xpath", fmt.Sprintf("%s - %s", entry.Name, err))
					return
				}

				entry.CopyMiscFrom(&listing[listingIndex])
				elm, err := specifier(entry)
				if err != nil {
					resp.Diagnostics.AddError("Error specifying update", fmt.Sprintf("%s - %s", entry.Name, err))
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
	for _, name := range tfid.Names {
		if _, ok := nameMap[name]; !nameIsUsed[name] && ok {
			path, err := newtfid.Location.Xpath(vn, name)
			if err != nil {
				resp.Diagnostics.AddError("Error building unused name delete xpath", fmt.Sprintf("(%s) %s", name, err))
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

	// Retrieve the list of objects again.
	listing, err = svc.List(ctx, newtfid.Location, "get", "", "")
	if err != nil {
		resp.Diagnostics.AddError("Error during Create's refresh", err.Error())
		return
	}

	// TODO: pick up from here: build up the config from the listing and save
	// it to state along with newtfid as a string.  This is similar to the ending
	// code from Create().

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ServiceObjectsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var err error
	var state ServiceObjectsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the tfid info.
	var tfid ServiceObjectsTfid
	if err = DecodeLocation(state.Tfid.ValueString(), &tfid); err != nil {
		resp.Diagnostics.AddError("Error parsing tfid", err.Error())
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource read", map[string]any{
		"resource_name": "panos_service_objects",
		"function":      "Delete",
	})

	// Verify mode.
	if r.client.Hostname == "" {
		resp.Diagnostics.AddError("Invalid mode error", InspectionModeError)
		return
	}

	// Create the service.
	svc := service.NewService(r.client)

	// Timeout handling.
	//ctx, cancel := context.WithTimeout(ctx, GetTimeout(state.Timeouts.Delete))
	//defer cancel()

	// Perform the operation.
	if err := svc.Delete(ctx, tfid.Location, tfid.Names...); err != nil {
		resp.Diagnostics.AddError("Error in delete", err.Error())
	}
}

func (r *ServiceObjectsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("tfid"), req, resp)
}
