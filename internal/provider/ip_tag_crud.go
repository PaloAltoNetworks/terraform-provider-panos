package provider

import (
	"context"
	"encoding/xml"
	"fmt"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/xmlapi"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// ipTagUidMessage is the type=user-id payload that registers or unregisters
// tags against an IP address. Exactly one of Register/Unregister is set.
type ipTagUidMessage struct {
	XMLName    xml.Name      `xml:"uid-message"`
	Version    string        `xml:"version"`
	Type       string        `xml:"type"`
	Register   *ipTagPayload `xml:"payload>register,omitempty"`
	Unregister *ipTagPayload `xml:"payload>unregister,omitempty"`
}

type ipTagPayload struct {
	Entries []ipTagEntry `xml:"entry"`
}

type ipTagEntry struct {
	Ip   string   `xml:"ip,attr"`
	Tags []string `xml:"tag>member"`
}

// registerIpTagCommand builds a user-id message that registers tags on an IP.
func registerIpTagCommand(ip string, tags []string) ipTagUidMessage {
	return ipTagUidMessage{
		Version:  "1.0",
		Type:     "update",
		Register: &ipTagPayload{Entries: []ipTagEntry{{Ip: ip, Tags: tags}}},
	}
}

// unregisterIpTagCommand builds a user-id message that removes tags from an IP.
func unregisterIpTagCommand(ip string, tags []string) ipTagUidMessage {
	return ipTagUidMessage{
		Version:    "1.0",
		Type:       "update",
		Unregister: &ipTagPayload{Entries: []ipTagEntry{{Ip: ip, Tags: tags}}},
	}
}

// registeredIpPageLimit is the maximum number of registered-ip entries PAN-OS
// returns per page; the read loop pages through results using start-point.
const registeredIpPageLimit = 500

// ipTagShowRequest is the type=op command that lists registered IP/tag mappings
// (PAN-OS 8.0+ paginated form).
type ipTagShowRequest struct {
	XMLName xml.Name        `xml:"show"`
	Filter  ipTagShowFilter `xml:"object>registered-ip"`
}

type ipTagShowFilter struct {
	Tag   *ipTagShowTagFilter `xml:"tag,omitempty"`
	Ip    string              `xml:"ip,omitempty"`
	Limit int                 `xml:"limit"`
	Start int                 `xml:"start-point"`
}

type ipTagShowTagFilter struct {
	Entry ipTagShowTagName `xml:"entry"`
}

type ipTagShowTagName struct {
	Name string `xml:"name,attr"`
}

// registeredIpRequest builds a single page of a registered-ip query. Both
// ipFilter and tagFilter are optional server-side filters; startPoint is the
// 1-based index of the first entry to return.
func registeredIpRequest(ipFilter, tagFilter string, startPoint int) ipTagShowRequest {
	req := ipTagShowRequest{}
	req.Filter.Ip = ipFilter
	req.Filter.Limit = registeredIpPageLimit
	req.Filter.Start = startPoint
	if tagFilter != "" {
		req.Filter.Tag = &ipTagShowTagFilter{Entry: ipTagShowTagName{Name: tagFilter}}
	}
	return req
}

// ipTagShowResponse decodes a registered-ip op response. The outfile path is
// only populated by pre-8.0 PAN-OS, which dumps results to a file instead of
// returning them inline; we treat that as an unsupported response shape.
type ipTagShowResponse struct {
	Entries []ipTagRespEntry `xml:"result>entry"`
	Outfile string           `xml:"result>msg>line>outfile"`
}

type ipTagRespEntry struct {
	Ip   string   `xml:"ip,attr"`
	Tags []string `xml:"tag>member"`
}

// toMap converts a decoded response into an IP→tags map, guarding against the
// pre-8.0 outfile response shape.
func (r ipTagShowResponse) toMap() (map[string][]string, error) {
	if r.Outfile != "" {
		return nil, fmt.Errorf("PAN-OS returned an outfile (%q) instead of registered-ip entries; PAN-OS 10.1+ is required", r.Outfile)
	}

	ans := make(map[string][]string, len(r.Entries))
	for _, entry := range r.Entries {
		ans[entry.Ip] = entry.Tags
	}
	return ans, nil
}

// parseRegisteredIpResponse decodes raw registered-ip op response bytes into an
// IP→tags map.
func parseRegisteredIpResponse(body []byte) (map[string][]string, error) {
	var resp ipTagShowResponse
	if err := xml.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return resp.toMap()
}

// pageFetcher returns a single page of registered-ip entries starting at the
// given 1-based start point. The IO layer backs this with a Communicate call.
type pageFetcher func(startPoint int) ([]ipTagRespEntry, error)

// collectRegisteredIps pages through registered-ip entries until a page shorter
// than pageLimit is returned, assembling them into a single IP→tags map.
func collectRegisteredIps(fetch pageFetcher, pageLimit int) (map[string][]string, error) {
	ans := make(map[string][]string)

	start := 1
	for {
		entries, err := fetch(start)
		if err != nil {
			return nil, err
		}

		for _, entry := range entries {
			ans[entry.Ip] = entry.Tags
		}

		if len(entries) < pageLimit {
			break
		}
		start += len(entries)
	}

	return ans, nil
}

type IpTagCustom struct{}

// tagSetDiff compares the current set of tags registered against an IP with the
// desired set, returning the tags that must be registered (in desired but not
// current) and unregistered (in current but not desired).
func tagSetDiff(current, desired []string) (toAdd, toRemove []string) {
	currentSet := make(map[string]struct{}, len(current))
	for _, tag := range current {
		currentSet[tag] = struct{}{}
	}

	desiredSet := make(map[string]struct{}, len(desired))
	for _, tag := range desired {
		desiredSet[tag] = struct{}{}
	}

	for tag := range desiredSet {
		if _, ok := currentSet[tag]; !ok {
			toAdd = append(toAdd, tag)
		}
	}

	for tag := range currentSet {
		if _, ok := desiredSet[tag]; !ok {
			toRemove = append(toRemove, tag)
		}
	}

	return toAdd, toRemove
}

func NewIpTagCustom(provider *ProviderData) (*IpTagCustom, error) {
	return &IpTagCustom{}, nil
}

// ReadCustom looks up every tag currently registered against the configured IP
// at the configured location. Unlike the resource, the data source has no notion
// of "managed" tags: it reports all tags present on the IP.
func (o *IpTagDataSource) ReadCustom(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data IpTagDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	vsys, target, diags := resolveIpTagLocation(ctx, data.Location)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ip := data.Ip.ValueString()
	if ip == "" {
		resp.Diagnostics.AddError("Invalid configuration", "'ip' must be set to the IP address to look up")
		return
	}

	present, err := registeredTagsForIp(ctx, o.client, vsys, target, ip)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read registered IP tags", err.Error())
		return
	}

	// Report every tag registered on the IP. An IP with no tags yields an empty
	// set rather than an error, so callers can detect "nothing registered".
	tagsValue, diags := types.SetValueFrom(ctx, types.StringType, present)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Tags = tagsValue

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// reconcileManagedTags returns the subset of managed tags that are still
// present on the firewall. Tags present on the IP but not managed by this
// resource are ignored. An empty result means the resource no longer owns any
// tags on this IP and should be removed from state.
func reconcileManagedTags(managed, present []string) []string {
	presentSet := make(map[string]struct{}, len(present))
	for _, tag := range present {
		presentSet[tag] = struct{}{}
	}

	var result []string
	for _, tag := range managed {
		if _, ok := presentSet[tag]; ok {
			result = append(result, tag)
		}
	}

	return result
}

// resolveIpTagLocation maps an ip_tag location into the vsys and target
// (Panorama-managed firewall serial) used by the user-id/op API calls. The
// location xpath is vestigial for this runtime resource/data source; only these
// two values matter. It is shared by both the resource and the data source.
func resolveIpTagLocation(ctx context.Context, locationObj types.Object) (vsys string, target string, diags diag.Diagnostics) {
	var loc IpTagLocation
	diags.Append(locationObj.As(ctx, &loc, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return "", "", diags
	}

	switch {
	case !loc.Vsys.IsNull() && !loc.Vsys.IsUnknown():
		var vsysLoc IpTagVsysLocation
		diags.Append(loc.Vsys.As(ctx, &vsysLoc, basetypes.ObjectAsOptions{})...)
		// The vsys name attribute is generated as "name" (provider-wide
		// convention), not "vsys".
		vsys = vsysLoc.Name.ValueString()
	case !loc.Panorama.IsNull() && !loc.Panorama.IsUnknown():
		// Registered directly on Panorama: no vsys, no target firewall.
	case !loc.TargetDevice.IsNull() && !loc.TargetDevice.IsUnknown():
		var targetLoc IpTagTargetDeviceLocation
		diags.Append(loc.TargetDevice.As(ctx, &targetLoc, basetypes.ObjectAsOptions{})...)
		vsys = targetLoc.Vsys.ValueString()
		target = targetLoc.Serial.ValueString()
		if target == "" {
			diags.AddError("Invalid location", "'target_device.serial' must be set to the serial number of the managed firewall")
		}
	default:
		diags.AddError("Invalid location", "Exactly one of 'vsys' or 'target_device' must be set")
	}

	return vsys, target, diags
}

// tagsFromSet reads a terraform string set into a plain slice.
func tagsFromSet(ctx context.Context, set types.Set) ([]string, diag.Diagnostics) {
	var tags []string
	diags := set.ElementsAs(ctx, &tags, false)
	return tags, diags
}

// sendUserId sends a register/unregister user-id message to PAN-OS.
func (o *IpTagResource) sendUserId(ctx context.Context, vsys, target string, msg ipTagUidMessage) error {
	cmd := &xmlapi.UserId{
		Command: msg,
		Vsys:    vsys,
		Target:  target,
	}
	_, _, err := o.client.Communicate(ctx, cmd, false, nil)
	return err
}

// registeredTagsForIp returns the tags currently registered against ip,
// paging through the registered-ip table as needed. It is shared by both the
// resource and the data source.
func registeredTagsForIp(ctx context.Context, client *pango.Client, vsys, target, ip string) ([]string, error) {
	fetch := func(startPoint int) ([]ipTagRespEntry, error) {
		cmd := &xmlapi.Op{
			Command: registeredIpRequest(ip, "", startPoint),
			Vsys:    vsys,
			Target:  target,
		}

		var resp ipTagShowResponse
		if _, _, err := client.Communicate(ctx, cmd, false, &resp); err != nil {
			return nil, err
		}
		if resp.Outfile != "" {
			return nil, fmt.Errorf("PAN-OS returned an outfile (%q) instead of registered-ip entries; PAN-OS 10.1+ is required", resp.Outfile)
		}
		return resp.Entries, nil
	}

	all, err := collectRegisteredIps(fetch, registeredIpPageLimit)
	if err != nil {
		return nil, err
	}
	return all[ip], nil
}

func (o *IpTagResource) ImportStateCustom(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.AddError("Import not supported", "The panos_ip_tag resource does not support terraform import.")
}

func (o *IpTagResource) CreateCustom(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data IpTagResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	vsys, target, diags := resolveIpTagLocation(ctx, data.Location)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tags, diags := tagsFromSet(ctx, data.Tags)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if len(tags) == 0 {
		resp.Diagnostics.AddError("Invalid configuration", "At least one tag must be specified")
		return
	}

	// Registering is additive and idempotent: re-registering an existing tag is
	// a no-op, and tags already on the IP that we don't manage are untouched.
	ip := data.Ip.ValueString()
	if err := o.sendUserId(ctx, vsys, target, registerIpTagCommand(ip, tags)); err != nil {
		resp.Diagnostics.AddError("Failed to register IP tags", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (o *IpTagResource) ReadCustom(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data IpTagResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	vsys, target, diags := resolveIpTagLocation(ctx, data.Location)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	managed, diags := tagsFromSet(ctx, data.Tags)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ip := data.Ip.ValueString()
	present, err := registeredTagsForIp(ctx, o.client, vsys, target, ip)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read registered IP tags", err.Error())
		return
	}

	// State reflects only the managed tags still present on the firewall; tags
	// we don't manage are ignored. If none remain, the resource is gone.
	reconciled := reconcileManagedTags(managed, present)
	if len(reconciled) == 0 {
		resp.State.RemoveResource(ctx)
		return
	}

	tagsValue, diags := types.SetValueFrom(ctx, types.StringType, reconciled)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Tags = tagsValue

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (o *IpTagResource) UpdateCustom(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan IpTagResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	var state IpTagResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	vsys, target, diags := resolveIpTagLocation(ctx, plan.Location)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	newTags, diags := tagsFromSet(ctx, plan.Tags)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if len(newTags) == 0 {
		resp.Diagnostics.AddError("Invalid configuration", "At least one tag must be specified")
		return
	}

	oldTags, diags := tagsFromSet(ctx, state.Tags)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	newIp := plan.Ip.ValueString()
	oldIp := state.Ip.ValueString()

	if oldIp != newIp {
		// The IP cannot be marked RequiresReplace via the spec, so handle a
		// change here: move all managed tags from the old IP to the new one.
		if len(oldTags) > 0 {
			if err := o.sendUserId(ctx, vsys, target, unregisterIpTagCommand(oldIp, oldTags)); err != nil {
				resp.Diagnostics.AddError("Failed to unregister IP tags from previous IP", err.Error())
				return
			}
		}
		if err := o.sendUserId(ctx, vsys, target, registerIpTagCommand(newIp, newTags)); err != nil {
			resp.Diagnostics.AddError("Failed to register IP tags", err.Error())
			return
		}
	} else {
		toAdd, toRemove := tagSetDiff(oldTags, newTags)
		if len(toAdd) > 0 {
			if err := o.sendUserId(ctx, vsys, target, registerIpTagCommand(newIp, toAdd)); err != nil {
				resp.Diagnostics.AddError("Failed to register IP tags", err.Error())
				return
			}
		}
		if len(toRemove) > 0 {
			if err := o.sendUserId(ctx, vsys, target, unregisterIpTagCommand(newIp, toRemove)); err != nil {
				resp.Diagnostics.AddError("Failed to unregister IP tags", err.Error())
				return
			}
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (o *IpTagResource) DeleteCustom(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data IpTagResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	vsys, target, diags := resolveIpTagLocation(ctx, data.Location)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	managed, diags := tagsFromSet(ctx, data.Tags)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ip := data.Ip.ValueString()
	present, err := registeredTagsForIp(ctx, o.client, vsys, target, ip)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read registered IP tags", err.Error())
		return
	}

	// Only unregister the managed tags that are still present, leaving any
	// tags we don't own (including overlapping ones from other resources) be.
	stillOurs := reconcileManagedTags(managed, present)
	if len(stillOurs) > 0 {
		if err := o.sendUserId(ctx, vsys, target, unregisterIpTagCommand(ip, stillOurs)); err != nil {
			resp.Diagnostics.AddError("Failed to unregister IP tags", err.Error())
			return
		}
	}
}
