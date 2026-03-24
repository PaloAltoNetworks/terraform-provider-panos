package provider

import (
	"context"
	"encoding/xml"
	"fmt"

	sdkerrors "github.com/PaloAltoNetworks/pango/errors"
	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/xmlapi"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DeviceGroupParentCustom struct{}

func NewDeviceGroupParentCustom(provider *ProviderData) (*DeviceGroupParentCustom, error) {
	return &DeviceGroupParentCustom{}, nil
}

type dgpReq struct {
	XMLName xml.Name `xml:"show"`
	Cmd     string   `xml:"dg-hierarchy"`
}

type dgpResp struct {
	Result *dgHierarchy `xml:"result>dg-hierarchy"`
}

func (o *dgpResp) results() map[string]string {
	ans := make(map[string]string)

	if o.Result != nil {
		for _, v := range o.Result.Info {
			ans[v.Name] = ""
			v.results(ans)
		}
	}

	return ans
}

type dgHierarchy struct {
	Info []dghInfo `xml:"dg"`
}

type dghInfo struct {
	Name     string    `xml:"name,attr"`
	Children []dghInfo `xml:"dg"`
}

func (o *dghInfo) results(ans map[string]string) {
	for _, v := range o.Children {
		ans[v.Name] = o.Name
		v.results(ans)
	}
}

type apReq struct {
	XMLName xml.Name `xml:"request"`
	Info    apInfo   `xml:"move-dg>entry"`
}

type apInfo struct {
	Child  string `xml:"name,attr"`
	Parent string `xml:"new-parent-dg,omitempty"`
}

func getParents(ctx context.Context, client util.PangoClient, deviceGroup string) (map[string]string, error) {
	cmd := &xmlapi.Op{
		Command: dgpReq{},
	}

	var ans dgpResp
	if _, _, err := client.Communicate(ctx, cmd, false, &ans); err != nil {
		return nil, err
	}

	return ans.results(), nil
}

func assignParent(ctx context.Context, client util.PangoClient, deviceGroup string, parent string) error {
	cmd := &xmlapi.Op{
		Command: apReq{
			Info: apInfo{
				Child:  deviceGroup,
				Parent: parent,
			},
		},
	}

	ans := util.JobResponse{}
	if _, _, err := client.Communicate(ctx, cmd, false, &ans); err != nil {
		return err
	}
	if err := client.WaitForJob(ctx, ans.Id, 0, nil); err != nil {
		return err
	}

	return nil
}

func (o *DeviceGroupParentDataSource) ReadCustom(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var state DeviceGroupParentResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := state.DeviceGroup.ValueString()
	hierarchy, err := getParents(ctx, o.client, name)
	if err != nil {
		if sdkerrors.IsObjectNotFound(err) {
			resp.State.RemoveResource(ctx)
		} else {
			resp.Diagnostics.AddError("Failed to query for the device group parent", err.Error())
		}
		return
	}

	parent, ok := hierarchy[name]
	if !ok {
		resp.Diagnostics.AddError("Failed to query for the device group parent", fmt.Sprintf("Device Group '%s' doesn't exist", name))
		return
	}
	state.Parent = types.StringValue(parent)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (o *DeviceGroupParentResource) CreateCustom(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var state DeviceGroupParentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deviceGroup := state.DeviceGroup.ValueString()
	parent := state.Parent.ValueString()
	if err := assignParent(ctx, o.client, deviceGroup, parent); err != nil {
		resp.Diagnostics.AddError("Failed to assign parent to the device group", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

}
func (o *DeviceGroupParentResource) ReadCustom(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var state DeviceGroupParentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := state.DeviceGroup.ValueString()
	hierarchy, err := getParents(ctx, o.client, name)
	if err != nil {
		if sdkerrors.IsObjectNotFound(err) {
			resp.State.RemoveResource(ctx)
		} else {
			resp.Diagnostics.AddError("Failed to query for the device group parent", err.Error())
		}
		return
	}

	parent, ok := hierarchy[name]
	if !ok {
		resp.State.RemoveResource(ctx)
		return
	}
	state.Parent = types.StringValue(parent)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

}
func (o *DeviceGroupParentResource) UpdateCustom(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var state DeviceGroupParentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deviceGroup := state.DeviceGroup.ValueString()
	parent := state.Parent.ValueString()
	if err := assignParent(ctx, o.client, deviceGroup, parent); err != nil {
		resp.Diagnostics.AddError("Failed to assign parent to the device group", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

}
func (o *DeviceGroupParentResource) DeleteCustom(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var state DeviceGroupParentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := state.DeviceGroup.ValueString()
	hierarchy, err := getParents(ctx, o.client, name)
	if err != nil {
		resp.Diagnostics.AddError("Failed to query for the device group parent", err.Error())
		return
	}

	parent, ok := hierarchy[name]
	if !ok {
		resp.Diagnostics.AddError("Failed to query for the device group parent", fmt.Sprintf("Device Group '%s' doesn't exist", name))
		return
	}

	if parent != "" {
		deviceGroup := state.DeviceGroup.ValueString()
		if err := assignParent(ctx, o.client, deviceGroup, ""); err != nil {
			resp.Diagnostics.AddError("Failed to assign parent to the device group", err.Error())
			return
		}
	}

}
