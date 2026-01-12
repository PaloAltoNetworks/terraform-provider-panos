package provider

import (
	"context"
	"encoding/xml"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/PaloAltoNetworks/pango/xmlapi"
)

type commitAllReq struct {
	XMLName      xml.Name                   `xml:"commit-all"`
	SharedPolicy *commitAllSharedPolicy     `xml:"shared-policy,omitempty"`
	DeviceGroup  *commitAllDeviceGroupEntry `xml:"device-group,omitempty"`
}

type commitAllSharedPolicy struct {
	DeviceGroup *commitAllDeviceGroupEntry `xml:"device-group,omitempty"`
}

type commitAllDeviceGroupEntry struct {
	Entry []commitAllEntry `xml:"entry"`
}

type commitAllEntry struct {
	Name string `xml:"name,attr"`
}

func (o commitAllReq) Action() string {
	return "all"
}

func (o commitAllReq) Element() any {
	return o
}

func (o *CommitAllAction) InvokeCustom(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	// Note: This action pushes committed configuration to device groups.
	// Ensure you run panos_commit first to commit changes to Panorama,
	// then run this action to push to managed devices.

	// Parse the action input
	var config CommitAllActionModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var deviceGroups []string

	// Check if device_groups parameter is specified
	if !config.DeviceGroups.IsNull() && !config.DeviceGroups.IsUnknown() {
		// User specified device groups
		resp.Diagnostics.Append(config.DeviceGroups.ElementsAs(ctx, &deviceGroups, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		tflog.Info(ctx, "Pushing to specified device groups", map[string]interface{}{
			"device_groups": deviceGroups,
		})
	} else {
		// No device groups specified - fetch all from Panorama
		tflog.Info(ctx, "No device groups specified, fetching all from Panorama")
		var err error
		deviceGroups, err = o.fetchAllDeviceGroups(ctx)
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to fetch device groups",
				fmt.Sprintf("Could not retrieve device groups from Panorama: %s\n\n"+
					"Possible causes:\n"+
					"1. No device groups are configured in Panorama\n"+
					"2. Insufficient permissions to view device groups\n"+
					"3. Connected to a firewall instead of Panorama\n\n"+
					"To push to specific device groups, use the 'device_groups' parameter.",
					err.Error()),
			)
			return
		}

		if len(deviceGroups) == 0 {
			resp.Diagnostics.AddError(
				"No device groups found",
				"Panorama has no configured device groups. The commit-all operation "+
					"pushes policies to device groups.\n\n"+
					"Solutions:\n"+
					"1. Create at least one device group in Panorama\n"+
					"2. Specify device groups using the 'device_groups' parameter\n"+
					"3. Use 'panos_commit' action for local commits only",
			)
			return
		}

		tflog.Info(ctx, "Found device groups", map[string]interface{}{
			"count":         len(deviceGroups),
			"device_groups": deviceGroups,
		})
	}

	// Build entry list for all device groups
	entries := make([]commitAllEntry, len(deviceGroups))
	for i, dg := range deviceGroups {
		entries[i] = commitAllEntry{Name: dg}
	}

	// Create commit-all request structure
	commitReq := &commitAllReq{
		SharedPolicy: &commitAllSharedPolicy{
			DeviceGroup: &commitAllDeviceGroupEntry{
				Entry: entries,
			},
		},
	}

	cmd := &xmlapi.Commit{
		Command: commitReq,
		Target:  "", // Don't use target here, it's in the cmd structure
	}

	tflog.Info(ctx, "Submitting commit-all request to Panorama")

	var commitResp xmlapi.JobResponse

	_, _, err := o.client.Communicate(ctx, cmd, false, &commitResp)
	if err != nil {
		resp.Diagnostics.AddError("Failed to schedule a commit-all (push to devices)", err.Error())
		return
	}

	tflog.Info(ctx, "Commit-all job started", map[string]interface{}{
		"job_id": commitResp.Id,
	})

	err = o.client.WaitForJob(ctx, commitResp.Id, 2*time.Second, nil)
	if err != nil {
		resp.Diagnostics.AddError("Failed to wait for commit-all task to finish", err.Error())
		return
	}

	tflog.Info(ctx, "Commit-all completed successfully")
}

func (o *CommitAllAction) fetchAllDeviceGroups(ctx context.Context) ([]string, error) {
	// Use operational command to get device groups
	// This is more reliable than config xpath queries across PAN-OS versions
	opCmd := &xmlapi.Op{
		Command: "<show><dg-hierarchy></dg-hierarchy></show>",
	}

	rawResp, _, err := o.client.Communicate(ctx, opCmd, false, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to query device groups: %w", err)
	}

	// Try multiple XML structures to support different PAN-OS versions

	// Format 1: Nested dg-info structure
	var xmlResp1 struct {
		Result struct {
			DgInfo struct {
				Dg []struct {
					Name string `xml:"dg-name"`
				} `xml:"dg"`
			} `xml:"dg-info"`
		} `xml:"result"`
	}

	if err = xml.Unmarshal(rawResp, &xmlResp1); err == nil && len(xmlResp1.Result.DgInfo.Dg) > 0 {
		deviceGroups := make([]string, len(xmlResp1.Result.DgInfo.Dg))
		for i, dg := range xmlResp1.Result.DgInfo.Dg {
			deviceGroups[i] = dg.Name
		}
		tflog.Debug(ctx, "Successfully fetched device groups", map[string]interface{}{
			"count":  len(deviceGroups),
			"groups": deviceGroups,
		})
		return deviceGroups, nil
	}

	// Format 2: Direct dg entries with name attribute
	var xmlResp2 struct {
		Result struct {
			Dg []struct {
				Name string `xml:"name,attr"`
			} `xml:"dg"`
		} `xml:"result"`
	}

	if err = xml.Unmarshal(rawResp, &xmlResp2); err == nil && len(xmlResp2.Result.Dg) > 0 {
		deviceGroups := make([]string, len(xmlResp2.Result.Dg))
		for i, dg := range xmlResp2.Result.Dg {
			deviceGroups[i] = dg.Name
		}
		tflog.Debug(ctx, "Successfully fetched device groups", map[string]interface{}{
			"count":  len(deviceGroups),
			"groups": deviceGroups,
		})
		return deviceGroups, nil
	}

	tflog.Debug(ctx, "No device groups found via operational command, trying config API")

	// Fallback to config API
	return o.fetchDeviceGroupsViaConfig(ctx)
}

func (o *CommitAllAction) fetchDeviceGroupsViaConfig(ctx context.Context) ([]string, error) {
	// Fallback method using config API when operational command fails
	configCmd := &xmlapi.Config{
		Action: "get",
		Xpath:  "/config/devices/entry[@name='localhost.localdomain']/device-group/entry",
	}

	rawResp, _, err := o.client.Communicate(ctx, configCmd, false, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to query device groups via config API: %w", err)
	}

	var xmlResp struct {
		Result struct {
			Entry []struct {
				Name string `xml:"name,attr"`
			} `xml:"entry"`
		} `xml:"result"`
	}

	if err = xml.Unmarshal(rawResp, &xmlResp); err != nil {
		return nil, fmt.Errorf("failed to parse device group response: %w", err)
	}

	if len(xmlResp.Result.Entry) == 0 {
		return nil, fmt.Errorf("no device groups found in Panorama")
	}

	deviceGroups := make([]string, len(xmlResp.Result.Entry))
	for i, entry := range xmlResp.Result.Entry {
		deviceGroups[i] = entry.Name
	}

	tflog.Debug(ctx, "Successfully fetched device groups via config API", map[string]interface{}{
		"count":  len(deviceGroups),
		"groups": deviceGroups,
	})

	return deviceGroups, nil
}
