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
    "github.com/hashicorp/terraform-plugin-framework/path"
    "github.com/hashicorp/terraform-plugin-framework/resource"
    rsschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/schema/validator"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Resource.
var (
    _ resource.Resource = &nestedAddressObjectResource{}
    _ resource.ResourceWithConfigure = &nestedAddressObjectResource{}
    _ resource.ResourceWithImportState = &nestedAddressObjectResource{}
)

func NewNestedAddressObjectResource() resource.Resource {
    return &nestedAddressObjectResource{}
}

type nestedAddressObjectResource struct {
    client *pango.XmlApiClient
}

type nestedAddressObjectLocation struct {
    Name string `json:"name"`
    Location address.Location `json:"location"`
}

func (o *nestedAddressObjectLocation) IsValid() error {
    if o.Name == "" {
        return fmt.Errorf("name is unspecified")
    }

    return o.Location.IsValid()
}

type nestedEntryModel struct {
    Tfid types.String `tfsdk:"tfid"`

    // Input.
    Location nestedLocationModel `tfsdk:"location"`

    Name types.String `tfsdk:"name"`
    Description types.String `tfsdk:"description"`
    Tags types.List `tfsdk:"tags"`
    IpNetmask types.String `tfsdk:"ip_netmask"`
    IpRange types.String `tfsdk:"ip_range"`
    Fqdn types.String `tfsdk:"fqdn"`
    IpWildcard types.String `tfsdk:"ip_wildcard"`
}

type nestedLocationModel struct {
    Shared types.Bool `tfsdk:"shared"`
    FromPanorama types.Bool `tfsdk:"from_panorama"`
    Vsys *nestedVsysLocation `tfsdk:"vsys"`
    DeviceGroup *nestedDeviceGroupLocation `tfsdk:"device_group"`
}

type nestedVsysLocation struct {
    Name types.String `tfsdk:"name"`
    NgfwDevice types.String `tfsdk:"ngfw_device"`
}

type nestedDeviceGroupLocation struct {
    Name types.String `tfsdk:"name"`
    PanoramaDevice types.String `tfsdk:"panorama_device"`
}

func (r *nestedAddressObjectResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_nested_address_object"
}

func (r *nestedAddressObjectResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = rsschema.Schema{
        Description: "Manages an address object.  This is the \"nested\" style where the location is a struct.",

        Attributes: map[string] rsschema.Attribute{
            "description": rsschema.StringAttribute{
                Description: "The description.",
                Optional: true,
                Validators: []validator.String{
                    stringvalidator.LengthAtMost(1023),
                },
            },
            "fqdn": rsschema.StringAttribute{
                Description: "The Fqdn param. String length must be between 1 and 255 characters. String validation regex: `^[a-zA-Z0-9_]([a-zA-Z0-9._-])+[a-zA-Z0-9]$`. Ensure that only one of the following is specified: `fqdn`, `ip_netmask`, `ip_range`, `ip_wildcard`",
                Optional:    true,
                Validators: []validator.String{
                    stringvalidator.LengthBetween(1, 255),
                    stringvalidator.RegexMatches(regexp.MustCompile("^[a-zA-Z0-9_]([a-zA-Z0-9._-])+[a-zA-Z0-9]$"), ""),
                    stringvalidator.ExactlyOneOf(
                        path.MatchRelative(),
                        path.MatchRelative().AtParent().AtName("ip_netmask"),
                        path.MatchRelative().AtParent().AtName("ip_range"),
                        path.MatchRelative().AtParent().AtName("ip_wildcard"),
                    ),
                },
            },
            "ip_netmask": rsschema.StringAttribute{
                Description: "The IpNetmask param. Ensure that only one of the following is specified: `fqdn`, `ip_netmask`, `ip_range`, `ip_wildcard`",
                Optional:    true,
            },
            "ip_range": rsschema.StringAttribute{
                Description: "The IpRange param. Ensure that only one of the following is specified: `fqdn`, `ip_netmask`, `ip_range`, `ip_wildcard`",
                Optional:    true,
            },
            "ip_wildcard": rsschema.StringAttribute{
                Description: "The IpWildcard param. Ensure that only one of the following is specified: `fqdn`, `ip_netmask`, `ip_range`, `ip_wildcard`",
                Optional:    true,
            },
            "location": rsschema.SingleNestedAttribute{
                Description: "The location of this object.",
                Required: true,
                Attributes: map[string] rsschema.Attribute{
                    "device_group": rsschema.SingleNestedAttribute{
                        Description: "(Panorama) In the given device group.",
                        Optional: true,
                        Attributes: map[string] rsschema.Attribute{
                            "name": rsschema.StringAttribute{
                                Description: "The device group name.",
                                Optional: true,
                                Computed: true,
                                Default: stringdefault.StaticString("vsys1"),
                                PlanModifiers: []planmodifier.String{
                                    stringplanmodifier.RequiresReplace(),
                                },
                            },
                            "panorama_device": rsschema.StringAttribute{
                                Description: "The Panorama device.",
                                Optional: true,
                                Computed: true,
                                Default: stringdefault.StaticString("localhost.localdomain"),
                                PlanModifiers: []planmodifier.String{
                                    stringplanmodifier.RequiresReplace(),
                                },
                            },
                        },
                    },
                    "from_panorama": rsschema.BoolAttribute{
                        Description: "(NGFW) Pushed from Panorama. This is a read-only location and only suitable for data sources.",
                        Optional: true,
                        Validators: []validator.Bool{
                            boolvalidator.ExactlyOneOf(
                                path.MatchRelative(),
                                path.MatchRelative().AtParent().AtName("device_group"),
                                path.MatchRelative().AtParent().AtName("shared"),
                                path.MatchRelative().AtParent().AtName("vsys"),
                            ),
                        },
                        PlanModifiers: []planmodifier.Bool{
                            boolplanmodifier.RequiresReplace(),
                        },
                    },
                    "shared": rsschema.BoolAttribute{
                        Description: "(NGFW and Panorama) Located in shared.",
                        Optional: true,
                        PlanModifiers: []planmodifier.Bool{
                            boolplanmodifier.RequiresReplace(),
                        },
                    },
                    "vsys": rsschema.SingleNestedAttribute{
                        Description: "(NGFW) In the given vsys.",
                        Optional: true,
                        Attributes: map[string] rsschema.Attribute{
                            "name": rsschema.StringAttribute{
                                Description: "The vsys name.",
                                Optional: true,
                                Computed: true,
                                Default: stringdefault.StaticString("vsys1"),
                                PlanModifiers: []planmodifier.String{
                                    stringplanmodifier.RequiresReplace(),
                                },
                            },
                            "ngfw_device": rsschema.StringAttribute{
                                Description: "The NGFW device.",
                                Optional: true,
                                Computed: true,
                                Default: stringdefault.StaticString("localhost.localdomain"),
                                PlanModifiers: []planmodifier.String{
                                    stringplanmodifier.RequiresReplace(),
                                },
                            },
                        },
                    },
                },
            },
            "name": rsschema.StringAttribute{
                Description: "Alphanumeric string [ 0-9a-zA-Z._-]. String length must not exceed 63 characters.",
                Required:    true,
                Validators: []validator.String{
                    stringvalidator.LengthAtMost(63),
                },
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
            "tfid": rsschema.StringAttribute{
                Description: "The Terraform ID.",
                Computed:    true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
        },
    }
}

func (r *nestedAddressObjectResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
    if req.ProviderData == nil {
        return
    }

    r.client = req.ProviderData.(*pango.XmlApiClient)
}

func (r *nestedAddressObjectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var state nestedEntryModel
    resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Basic logging.
    tflog.Info(ctx, "performing resource create", map[string] any{
        "resource_name": "panos_nested_address_object",
        "function": "Create",
        "name": state.Name.ValueString(),
    })

    // Create the service.
    svc := address.NewService(r.client)

    // Determine the location.
    loc := nestedAddressObjectLocation{Name: state.Name.ValueString()}
    if state.Location.Shared.ValueBool() {
        loc.Location.Shared = true
    } else if state.Location.FromPanorama.ValueBool() {
        loc.Location.FromPanorama = true
    } else if state.Location.Vsys != nil {
        loc.Location.Vsys = &address.VsysLocation{}
        loc.Location.Vsys.Name = state.Location.Vsys.Name.ValueString()
        loc.Location.Vsys.NgfwDevice = state.Location.Vsys.NgfwDevice.ValueString()
    } else if state.Location.DeviceGroup != nil {
        loc.Location.DeviceGroup = &address.DeviceGroupLocation{}
        loc.Location.DeviceGroup.Name = state.Location.DeviceGroup.Name.ValueString()
        loc.Location.DeviceGroup.PanoramaDevice = state.Location.DeviceGroup.PanoramaDevice.ValueString()
    } else {
        resp.Diagnostics.AddError("Unknown location", "Location for object is unknown")
        return
    }

    // Load the desired config.
    obj := address.Entry{Name: state.Name.ValueString()}
    obj.Description = state.Description.ValueStringPointer()
    obj.IpNetmask = state.IpNetmask.ValueStringPointer()
    obj.IpRange = state.IpRange.ValueStringPointer()
    obj.Fqdn = state.Fqdn.ValueStringPointer()
    obj.IpWildcard = state.IpWildcard.ValueStringPointer()
    resp.Diagnostics.Append(state.Tags.ElementsAs(ctx, &obj.Tags, false)...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Perform the operation.
    ans, err := svc.Create(ctx, loc.Location, obj)
    if err != nil {
        resp.Diagnostics.AddError("Error in create", err.Error())
        return
    }

    // Save the tfid.
    tfidstr, err := EncodeLocation(&loc)
    if err != nil {
        resp.Diagnostics.AddError("error creating tfid", err.Error())
        return
    }
    state.Tfid = types.StringValue(tfidstr)

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

// Read performs Read for the struct.
func (r *nestedAddressObjectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var savestate, state nestedEntryModel
    resp.Diagnostics.Append(req.State.Get(ctx, &savestate)...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Parse the location from tfid.
    var loc nestedAddressObjectLocation
    if err := DecodeLocation(savestate.Tfid.ValueString(), &loc); err != nil {
        resp.Diagnostics.AddError("error parsing tfid", err.Error())
        return
    }

    // Basic logging.
    tflog.Info(ctx, "performing resource read", map[string] any{
        "resource_name": "panos_nested_address_object",
        "function": "Read",
        "name": loc.Name,
    })

    // Create the service.
    svc := address.NewService(r.client)

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
    } else if loc.Location.FromPanorama {
        state.Location.FromPanorama = types.BoolValue(true)
    } else if loc.Location.Vsys != nil {
        state.Location.Vsys = &nestedVsysLocation{}
        state.Location.Vsys.Name = types.StringValue(loc.Location.Vsys.Name)
        state.Location.Vsys.NgfwDevice = types.StringValue(loc.Location.Vsys.NgfwDevice)
    } else if loc.Location.DeviceGroup != nil {
        state.Location.DeviceGroup = &nestedDeviceGroupLocation{}
        state.Location.DeviceGroup.Name = types.StringValue(loc.Location.DeviceGroup.Name)
        state.Location.DeviceGroup.PanoramaDevice = types.StringValue(loc.Location.DeviceGroup.PanoramaDevice)
    }

    // Save the answer to state.
    state.Tfid = savestate.Tfid
    state.Name = types.StringValue(loc.Name)
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

func (r *nestedAddressObjectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var plan, state nestedEntryModel
    resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
    resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
    if resp.Diagnostics.HasError() {
        return
    }

    var loc nestedAddressObjectLocation
    if err := DecodeLocation(state.Tfid.ValueString(), &loc); err != nil {
        resp.Diagnostics.AddError("error parsing tfid", err.Error())
        return
    }

    // Basic logging.
    tflog.Info(ctx, "performing resource update", map[string] any{
        "resource_name": "panos_nested_address_object",
        "function": "Update",
        "tfid": state.Tfid.ValueString(),
    })

    // Create the service.
    svc := address.NewService(r.client)

    // Load the desired config.
    obj := address.Entry{Name: plan.Name.ValueString()}
    obj.Description = plan.Description.ValueStringPointer()
    obj.IpNetmask = plan.IpNetmask.ValueStringPointer()
    obj.IpRange = plan.IpRange.ValueStringPointer()
    obj.Fqdn = plan.Fqdn.ValueStringPointer()
    obj.IpWildcard = plan.IpWildcard.ValueStringPointer()
    resp.Diagnostics.Append(plan.Tags.ElementsAs(ctx, &obj.Tags, false)...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Perform the operation.
    ans, err := svc.Update(ctx, loc.Location, obj, loc.Name)
    if err != nil {
        resp.Diagnostics.AddError("Error in update", err.Error())
        return
    }

    // Save the tfid.
    loc.Name = obj.Name
    tfidstr, err := EncodeLocation(&loc)
    if err != nil {
        resp.Diagnostics.AddError("error creating tfid", err.Error())
        return
    }
    state.Tfid = types.StringValue(tfidstr)

    // Save the state.
    state.Location = plan.Location
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

func (r *nestedAddressObjectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var idType types.String
    resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("tfid"), &idType)...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Parse the location from tfid.
    var loc nestedAddressObjectLocation
    if err := DecodeLocation(idType.ValueString(), &loc); err != nil {
        resp.Diagnostics.AddError("error parsing tfid", err.Error())
        return
    }

    // Basic logging.
    tflog.Info(ctx, "performing resource delete", map[string] any{
        "resource_name": "panos_nested_address_object",
        "function": "Delete",
        "name": loc.Name,
    })

    // Create the service.
    svc := address.NewService(r.client)

    // Perform the operation.
    if err := svc.Delete(ctx, loc.Location, loc.Name); err != nil && !IsObjectNotFound(err) {
        resp.Diagnostics.AddError("Error in delete", err.Error())
    }
}

func (r *nestedAddressObjectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    resource.ImportStatePassthroughID(ctx, path.Root("tfid"), req, resp)
}
