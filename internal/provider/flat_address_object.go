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
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/schema/validator"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Resource.
var (
    _ resource.Resource = &flatAddressObjectResource{}
    _ resource.ResourceWithConfigure = &flatAddressObjectResource{}
    _ resource.ResourceWithImportState = &flatAddressObjectResource{}
)

func NewFlatAddressObjectResource() resource.Resource {
    return &flatAddressObjectResource{}
}

type flatAddressObjectResource struct {
    client *pango.XmlApiClient
}

type flatAddressObjectLocation struct {
    Name string `json:"name"`
    Location address.Location `json:"location"`
}

func (o *flatAddressObjectLocation) IsValid() error {
    if o.Name == "" {
        return fmt.Errorf("name is unspecified")
    }

    return o.Location.IsValid()
}

type flatEntryModel struct {
    Tfid types.String `tfsdk:"tfid"`

    // Location.
    Shared types.Bool `tfsdk:"shared"`
    FromPanorama types.Bool `tfsdk:"from_panorama"`
    Vsys types.String `tfsdk:"vsys"`
    DeviceGroup types.String `tfsdk:"device_group"`

    // Input.

    Name types.String `tfsdk:"name"`
    Description types.String `tfsdk:"description"`
    Tags types.List `tfsdk:"tags"`
    IpNetmask types.String `tfsdk:"ip_netmask"`
    IpRange types.String `tfsdk:"ip_range"`
    Fqdn types.String `tfsdk:"fqdn"`
    IpWildcard types.String `tfsdk:"ip_wildcard"`
}

func (r *flatAddressObjectResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_flat_address_object"
}

func (r *flatAddressObjectResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = rsschema.Schema{
        Description: "Manages an address object.  This is the \"flat\" style where the location is mixed in with the params for the object itself.",

        Attributes: map[string] rsschema.Attribute{
            // Location params.
            "from_panorama": rsschema.BoolAttribute{
                Description: "(Location param; NGFW only) Pushed from Panorama. This is a read-only location and only suitable for data sources.",
                Optional: true,
                Validators: []validator.Bool{
                    boolvalidator.ExactlyOneOf(
                        path.MatchRoot("from_panorama"),
                        path.MatchRoot("device_group"),
                        path.MatchRoot("shared"),
                        path.MatchRoot("vsys"),
                    ),
                },
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.RequiresReplace(),
                },
            },
            "device_group": rsschema.StringAttribute{
                Description: "(Location param; Panorama only) The device group name.",
                Optional: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.RequiresReplace(),
                },
            },
            "shared": rsschema.BoolAttribute{
                Description: "(Location param; NGFW and Panorama) Located in shared.",
                Optional: true,
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.RequiresReplace(),
                },
            },
            "vsys": rsschema.StringAttribute{
                Description: "(Location param; NGFW only) The vsys name.",
                Optional: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.RequiresReplace(),
                },
            },

            // Params.
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
                        path.MatchRoot("ip_netmask"),
                        path.MatchRoot("ip_range"),
                        path.MatchRoot("ip_wildcard"),
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

func (r *flatAddressObjectResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
    if req.ProviderData == nil {
        return
    }

    r.client = req.ProviderData.(*pango.XmlApiClient)
}

func (r *flatAddressObjectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var state flatEntryModel
    resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Basic logging.
    tflog.Info(ctx, "performing resource create", map[string] any{
        "resource_name": "panos_flat_address_object",
        "function": "Create",
        "name": state.Name.ValueString(),
    })

    // Create the service.
    svc := address.NewService(r.client)

    // Determine the location.
    loc := flatAddressObjectLocation{Name: state.Name.ValueString()}
    if !state.Shared.IsNull() && state.Shared.ValueBool() {
        loc.Location.Shared = true
    } else if !state.FromPanorama.IsNull() && state.FromPanorama.ValueBool() {
        loc.Location.FromPanorama = true
    } else if !state.Vsys.IsNull() && state.Vsys.ValueString() != "" {
        loc.Location.Vsys = &address.VsysLocation{}
        loc.Location.Vsys.Name = state.Vsys.ValueString()
        loc.Location.Vsys.NgfwDevice = "localhost.localdomain"
    } else if !state.DeviceGroup.IsNull() && state.DeviceGroup.ValueString() != "" {
        loc.Location.DeviceGroup = &address.DeviceGroupLocation{}
        loc.Location.DeviceGroup.Name = state.DeviceGroup.ValueString()
        loc.Location.DeviceGroup.PanoramaDevice = "localhost.localdomain"
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
func (r *flatAddressObjectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var savestate, state flatEntryModel
    resp.Diagnostics.Append(req.State.Get(ctx, &savestate)...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Parse the location from tfid.
    var loc flatAddressObjectLocation
    if err := DecodeLocation(savestate.Tfid.ValueString(), &loc); err != nil {
        resp.Diagnostics.AddError("error parsing tfid", err.Error())
        return
    }

    // Basic logging.
    tflog.Info(ctx, "performing resource read", map[string] any{
        "resource_name": "panos_flat_address_object",
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
        state.Shared = types.BoolValue(true)
    } else if loc.Location.FromPanorama {
        state.FromPanorama = types.BoolValue(true)
    } else if loc.Location.Vsys != nil {
        state.Vsys = types.StringValue(loc.Location.Vsys.Name)
    } else if loc.Location.DeviceGroup != nil {
        state.DeviceGroup = types.StringValue(loc.Location.DeviceGroup.Name)
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

func (r *flatAddressObjectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var plan, state flatEntryModel
    resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
    resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
    if resp.Diagnostics.HasError() {
        return
    }

    var loc flatAddressObjectLocation
    if err := DecodeLocation(state.Tfid.ValueString(), &loc); err != nil {
        resp.Diagnostics.AddError("error parsing tfid", err.Error())
        return
    }

    // Basic logging.
    tflog.Info(ctx, "performing resource update", map[string] any{
        "resource_name": "panos_flat_address_object",
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

    // Save the location.
    state.Shared = plan.Shared
    state.FromPanorama = plan.FromPanorama
    state.Vsys = plan.Vsys
    state.DeviceGroup = plan.DeviceGroup

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

func (r *flatAddressObjectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var idType types.String
    resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("tfid"), &idType)...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Parse the location from tfid.
    var loc flatAddressObjectLocation
    if err := DecodeLocation(idType.ValueString(), &loc); err != nil {
        resp.Diagnostics.AddError("error parsing tfid", err.Error())
        return
    }

    // Basic logging.
    tflog.Info(ctx, "performing resource delete", map[string] any{
        "resource_name": "panos_flat_address_object",
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

func (r *flatAddressObjectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    resource.ImportStatePassthroughID(ctx, path.Root("tfid"), req, resp)
}
