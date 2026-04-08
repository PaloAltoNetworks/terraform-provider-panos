package provider

import (
	"context"
	"errors"

	"github.com/PaloAltoNetworks/pango/locking"
	vrouter "github.com/PaloAltoNetworks/pango/network/virtual_router"
	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/xmlapi"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	sdkmanager "github.com/PaloAltoNetworks/terraform-provider-panos/internal/manager"
)

type VirtualRouterInterfaceCustom struct {
	specifier vrouter.Specifier
	manager   *sdkmanager.EntryObjectManager[*vrouter.Entry, vrouter.Location, *vrouter.Service]
}

func NewVirtualRouterInterfaceCustom(provider *ProviderData) (*VirtualRouterInterfaceCustom, error) {
	client := provider.Client

	specifier, _, err := vrouter.Versioning(client.Versioning())
	if err != nil {
		return nil, err
	}

	manager := sdkmanager.NewEntryObjectManager[*vrouter.Entry, vrouter.Location, *vrouter.Service](
		client, vrouter.NewService(client), provider.MultiConfigBatchSize, specifier, vrouter.SpecMatches)

	return &VirtualRouterInterfaceCustom{
		specifier: specifier,
		manager:   manager,
	}, nil
}

func terraformToSdkLocation(ctx context.Context, locationObj types.Object) (*vrouter.Location, diag.Diagnostics) {
	var terraformLocation VirtualRouterInterfaceLocation

	diags := locationObj.As(ctx, &terraformLocation, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}

	location := &vrouter.Location{}

	if !terraformLocation.Ngfw.IsNull() {
		location.Ngfw = &vrouter.NgfwLocation{}
		var innerLocation VirtualRouterInterfaceNgfwLocation
		diags.Append(terraformLocation.Ngfw.As(ctx, &innerLocation, basetypes.ObjectAsOptions{})...)
		if diags.HasError() {
			return nil, diags
		}
		location.Ngfw.NgfwDevice = innerLocation.NgfwDevice.ValueString()
	}

	if !terraformLocation.Template.IsNull() {
		location.Template = &vrouter.TemplateLocation{}
		var innerLocation VirtualRouterInterfaceTemplateLocation
		diags.Append(terraformLocation.Template.As(ctx, &innerLocation, basetypes.ObjectAsOptions{})...)
		if diags.HasError() {
			return nil, diags
		}
		location.Template.PanoramaDevice = innerLocation.PanoramaDevice.ValueString()
		location.Template.Template = innerLocation.Name.ValueString()
		location.Template.NgfwDevice = innerLocation.NgfwDevice.ValueString()
	}

	if !terraformLocation.TemplateStack.IsNull() {
		location.TemplateStack = &vrouter.TemplateStackLocation{}
		var innerLocation VirtualRouterInterfaceTemplateStackLocation
		diags.Append(terraformLocation.TemplateStack.As(ctx, &innerLocation, basetypes.ObjectAsOptions{})...)
		if diags.HasError() {
			return nil, diags
		}
		location.TemplateStack.PanoramaDevice = innerLocation.PanoramaDevice.ValueString()
		location.TemplateStack.TemplateStack = innerLocation.Name.ValueString()
		location.TemplateStack.NgfwDevice = innerLocation.NgfwDevice.ValueString()
	}

	if !terraformLocation.Vsys.IsNull() {
		location.Vsys = &vrouter.VsysLocation{}
		var innerLocation VirtualRouterInterfaceVsysLocation
		diags.Append(terraformLocation.Vsys.As(ctx, &innerLocation, basetypes.ObjectAsOptions{})...)
		if diags.HasError() {
			return nil, diags
		}
		location.Vsys.NgfwDevice = innerLocation.NgfwDevice.ValueString()
		location.Vsys.Vsys = innerLocation.Name.ValueString()
	}

	return location, nil
}

func (o *VirtualRouterInterfaceResource) ReadCustom(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state VirtualRouterInterfaceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	location, diags := terraformToSdkLocation(ctx, state.Location)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	components := []string{}
	object, err := o.custom.manager.Read(ctx, *location, components, state.VirtualRouter.ValueString())
	if err != nil {
		if errors.Is(err, sdkmanager.ErrObjectNotFound) {
			resp.Diagnostics.AddError("Error reading data", err.Error())
		} else {
			resp.Diagnostics.AddError("Error reading entry", err.Error())
		}
		return
	}

	var found bool
	for _, elt := range object.Interface {
		if elt == state.Interface.ValueString() {
			found = true
			break
		}
	}

	if !found {
		resp.State.RemoveResource(ctx)
	}
}

func (o *VirtualRouterInterfaceResource) CreateCustom(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan VirtualRouterInterfaceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	location, diags := terraformToSdkLocation(ctx, plan.Location)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	components := []string{}
	xpath, err := location.XpathWithComponents(
		o.client.Versioning(),
		append(components, util.AsEntryXpath(plan.VirtualRouter.ValueString()))...)
	if err != nil {
		resp.Diagnostics.AddError("Error while generating xpath for parent resource", err.Error())
		return
	}

	mutex := locking.GetMutex(locking.XpathLockCategory, util.AsXpath(xpath))
	mutex.Lock()
	defer mutex.Unlock()

	object, err := o.custom.manager.Read(ctx, *location, components, plan.VirtualRouter.ValueString())

	if err != nil {
		if errors.Is(err, sdkmanager.ErrObjectNotFound) {
			resp.Diagnostics.AddError("Parent resource missing", "Virtual router not found")
		} else {
			resp.Diagnostics.AddError("Error while reading parent resource", err.Error())
		}
		return
	}

	var found bool
	for _, elt := range object.Interface {
		if elt == plan.Interface.ValueString() {
			found = true
			break
		}
	}

	if found {
		resp.Diagnostics.AddError("Error while creating interface entry", "entry with a matching name already exists")
		return
	}

	object.Interface = o.addInterface(object.Interface, plan.Interface.ValueString())
	_, err = o.custom.manager.Update(ctx, *location, components, object, "")
	if err != nil {
		resp.Diagnostics.AddError("Error while creating interface entry", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (o *VirtualRouterInterfaceResource) removeInterface(ifaces []string, needle string) []string {
	var result []string
	for _, elt := range ifaces {
		if elt != needle {
			result = append(result, elt)
		}
	}

	return result
}

func (o *VirtualRouterInterfaceResource) addInterface(ifaces []string, needle string) []string {
	var found bool
	var result []string

	for _, elt := range ifaces {
		if elt == needle {
			found = true
		}
		result = append(result, elt)
	}

	if !found {
		result = append(result, needle)
	}

	return result
}

func (o *VirtualRouterInterfaceResource) UpdateCustom(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan VirtualRouterInterfaceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	location, diags := terraformToSdkLocation(ctx, state.Location)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updates := xmlapi.NewMultiConfig(2)

	components := []string{}
	if state.VirtualRouter.ValueString() != plan.VirtualRouter.ValueString() {
		xpath, err := location.XpathWithComponents(
			o.client.Versioning(),
			append(components, util.AsEntryXpath(state.VirtualRouter.ValueString()))...)
		if err != nil {
			resp.Diagnostics.AddError("Error while creating xpath for parent resource", err.Error())
			return
		}

		mutex := locking.GetMutex(locking.XpathLockCategory, util.AsXpath(xpath))
		mutex.Lock()
		defer mutex.Unlock()

		object, err := o.custom.manager.Read(ctx, *location, components, state.VirtualRouter.ValueString())
		if err != nil {
			if !errors.Is(err, sdkmanager.ErrObjectNotFound) {
				resp.Diagnostics.AddError("Error while reading parent resource", err.Error())
				return
			}
		}

		if object != nil {
			object.Interface = o.removeInterface(object.Interface, state.Interface.ValueString())
			xmlEntry, err := o.custom.specifier(object)
			if err != nil {
				resp.Diagnostics.AddError("Error while creating XML document for parent resource", err.Error())
				return
			}

			updates.Add(&xmlapi.Config{
				Action:  "edit",
				Xpath:   util.AsXpath(xpath),
				Element: xmlEntry,
				Target:  o.client.GetTarget(),
			})
		}
	}

	xpath, err := location.XpathWithComponents(
		o.client.Versioning(),
		append(components, util.AsEntryXpath(plan.VirtualRouter.ValueString()))...)
	if err != nil {
		resp.Diagnostics.AddError("Error while creating xpath for parent resource", err.Error())
		return
	}

	mutex := locking.GetMutex(locking.XpathLockCategory, util.AsXpath(xpath))
	mutex.Lock()
	defer mutex.Unlock()

	object, err := o.custom.manager.Read(ctx, *location, components, plan.VirtualRouter.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error while reading parent resource", err.Error())
		return
	}

	object.Interface = o.removeInterface(object.Interface, state.Interface.ValueString())
	object.Interface = o.addInterface(object.Interface, plan.Interface.ValueString())
	xmlEntry, err := o.custom.specifier(object)
	if err != nil {
		resp.Diagnostics.AddError("Error while creating XML document for parent resource", err.Error())
		return
	}

	updates.Add(&xmlapi.Config{
		Action:  "edit",
		Xpath:   util.AsXpath(xpath),
		Element: xmlEntry,
		Target:  o.client.GetTarget(),
	})

	if len(updates.Operations) > 0 {
		if _, _, _, err := o.client.MultiConfig(ctx, updates, false, nil); err != nil {
			resp.Diagnostics.AddError("Error while updating parent resource", err.Error())
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (o *VirtualRouterInterfaceResource) DeleteCustom(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state VirtualRouterInterfaceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	location, diags := terraformToSdkLocation(ctx, state.Location)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	components := []string{}
	xpath, err := location.XpathWithComponents(
		o.client.Versioning(),
		append(components, util.AsEntryXpath(state.VirtualRouter.ValueString()))...)
	if err != nil {
		resp.Diagnostics.AddError("Error while generating xpath for parent resource", err.Error())
		return
	}

	mutex := locking.GetMutex(locking.XpathLockCategory, util.AsXpath(xpath))
	mutex.Lock()
	defer mutex.Unlock()

	object, err := o.custom.manager.Read(ctx, *location, components, state.VirtualRouter.ValueString())
	if err != nil {
		if !errors.Is(err, sdkmanager.ErrObjectNotFound) {
			resp.Diagnostics.AddError("Error while reading parent resource", err.Error())
		}
		return
	}

	object.Interface = o.removeInterface(object.Interface, state.Interface.ValueString())
	o.custom.manager.Update(ctx, *location, components, object, "")
	if err != nil {
		resp.Diagnostics.AddError("Error while deleting interface", err.Error())
	}

	return
}
