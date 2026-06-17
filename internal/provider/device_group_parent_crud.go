package provider

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"

	sdkerrors "github.com/PaloAltoNetworks/pango/errors"
	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/xmlapi"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
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

type DeviceGroupParentImportState struct {
	Location    types.Object `json:"location"`
	DeviceGroup types.String `json:"device_group"`
}

func (o DeviceGroupParentImportState) MarshalJSON() ([]byte, error) {
	type shadow struct {
		Location    interface{} `json:"location"`
		DeviceGroup *string     `json:"device_group"`
	}
	location_object, err := TypesObjectToMap(o.Location, DeviceGroupParentLocationSchema())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal location into JSON document: %w", err)
	}

	return json.Marshal(shadow{
		Location:    location_object,
		DeviceGroup: o.DeviceGroup.ValueStringPointer(),
	})
}

func (o *DeviceGroupParentImportState) UnmarshalJSON(data []byte) error {
	var shadow struct {
		Location    interface{} `json:"location"`
		DeviceGroup *string     `json:"device_group"`
	}

	if err := json.Unmarshal(data, &shadow); err != nil {
		return err
	}

	location_map, ok := shadow.Location.(map[string]interface{})
	if !ok {
		return NewDiagnosticsError("Failed to unmarshal JSON document into location: expected map[string]interface{}", nil)
	}
	location_object, err := MapToTypesObject(location_map, DeviceGroupParentLocationSchema())
	if err != nil {
		return fmt.Errorf("failed to unmarshal location from JSON: %w", err)
	}
	o.Location = location_object
	o.DeviceGroup = types.StringPointerValue(shadow.DeviceGroup)

	return nil
}

func DeviceGroupParentImportStateCreator(ctx context.Context, resource types.Object) ([]byte, error) {
	attrs := resource.Attributes()
	if attrs == nil {
		return nil, fmt.Errorf("Object has no attributes")
	}

	locationAttr, ok := attrs["location"]
	if !ok {
		return nil, fmt.Errorf("location attribute missing")
	}
	location, ok := locationAttr.(types.Object)
	if !ok {
		return nil, fmt.Errorf("location attribute expected to be an object")
	}

	deviceGroupAttr, ok := attrs["device_group"]
	if !ok {
		return nil, fmt.Errorf("device_group attribute missing")
	}
	deviceGroup, ok := deviceGroupAttr.(types.String)
	if !ok {
		return nil, fmt.Errorf("device_group attribute expected to be a string")
	}

	return json.Marshal(DeviceGroupParentImportState{
		Location:    location,
		DeviceGroup: deviceGroup,
	})
}

func (o *DeviceGroupParentResource) ImportStateCustom(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	data, err := base64.StdEncoding.DecodeString(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Failed to decode Import ID", err.Error())
		return
	}

	var obj DeviceGroupParentImportState
	if err := json.Unmarshal(data, &obj); err != nil {
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
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("device_group"), obj.DeviceGroup)...)
}
