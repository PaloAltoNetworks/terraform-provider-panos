package provider

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/PaloAltoNetworks/pango/commit"
	"github.com/PaloAltoNetworks/pango/xmlapi"
)

func (o *CommitAction) InvokeCustom(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var model CommitActionModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := o.client.RetrieveSystemInfo(ctx); err != nil {
		resp.Diagnostics.AddError("Failed to retrieve system information", err.Error())
		return
	}

	isFirewall, err := o.client.IsFirewall()
	if err != nil {
		resp.Diagnostics.AddError("Failed to determine device type", err.Error())
		return
	}

	var description string
	var admins []string
	var excludeDeviceAndNetwork bool
	var excludeSharedObjects bool
	var force bool
	var excludePolicyAndObjects bool
	var deviceGroups []string
	var templates []string
	var templateStacks []string
	var wildfireAppliances []string
	var wildfireClusters []string
	var logCollectors []string
	var logCollectorGroups []string

	if !model.Description.IsNull() && !model.Description.IsUnknown() {
		description = model.Description.ValueString()
	}

	if !model.Admins.IsNull() && !model.Admins.IsUnknown() {
		resp.Diagnostics.Append(model.Admins.ElementsAs(ctx, &admins, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if !model.ExcludeDeviceAndNetwork.IsNull() && !model.ExcludeDeviceAndNetwork.IsUnknown() {
		excludeDeviceAndNetwork = model.ExcludeDeviceAndNetwork.ValueBool()
	}

	if !model.ExcludeSharedObjects.IsNull() && !model.ExcludeSharedObjects.IsUnknown() {
		excludeSharedObjects = model.ExcludeSharedObjects.ValueBool()
	}

	if !model.Force.IsNull() && !model.Force.IsUnknown() {
		force = model.Force.ValueBool()
	}

	if !model.ExcludePolicyAndObjects.IsNull() && !model.ExcludePolicyAndObjects.IsUnknown() {
		excludePolicyAndObjects = model.ExcludePolicyAndObjects.ValueBool()
	}

	if !model.DeviceGroups.IsNull() && !model.DeviceGroups.IsUnknown() {
		resp.Diagnostics.Append(model.DeviceGroups.ElementsAs(ctx, &deviceGroups, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if !model.Templates.IsNull() && !model.Templates.IsUnknown() {
		resp.Diagnostics.Append(model.Templates.ElementsAs(ctx, &templates, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if !model.TemplateStacks.IsNull() && !model.TemplateStacks.IsUnknown() {
		resp.Diagnostics.Append(model.TemplateStacks.ElementsAs(ctx, &templateStacks, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if !model.WildfireAppliances.IsNull() && !model.WildfireAppliances.IsUnknown() {
		resp.Diagnostics.Append(model.WildfireAppliances.ElementsAs(ctx, &wildfireAppliances, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if !model.WildfireClusters.IsNull() && !model.WildfireClusters.IsUnknown() {
		resp.Diagnostics.Append(model.WildfireClusters.ElementsAs(ctx, &wildfireClusters, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if !model.LogCollectors.IsNull() && !model.LogCollectors.IsUnknown() {
		resp.Diagnostics.Append(model.LogCollectors.ElementsAs(ctx, &logCollectors, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if !model.LogCollectorGroups.IsNull() && !model.LogCollectorGroups.IsUnknown() {
		resp.Diagnostics.Append(model.LogCollectorGroups.ElementsAs(ctx, &logCollectorGroups, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if isFirewall {
		if len(deviceGroups) > 0 || len(templates) > 0 || len(templateStacks) > 0 ||
			len(wildfireAppliances) > 0 || len(wildfireClusters) > 0 ||
			len(logCollectors) > 0 || len(logCollectorGroups) > 0 {
			resp.Diagnostics.AddError(
				"Invalid parameters for firewall",
				"Parameters device_groups, templates, template_stacks, wildfire_appliances, wildfire_clusters, log_collectors, and log_collector_groups are only available on Panorama",
			)
			return
		}
	} else {
		if excludePolicyAndObjects {
			resp.Diagnostics.AddError(
				"Invalid parameter for Panorama",
				"Parameter exclude_policy_and_objects is only available on NGFW",
			)
			return
		}
	}

	var commitCmd xmlapi.CommitAction
	if isFirewall {
		commitCmd = commit.FirewallCommit{
			Description:             description,
			Admins:                  admins,
			ExcludeDeviceAndNetwork: excludeDeviceAndNetwork,
			ExcludeSharedObjects:    excludeSharedObjects,
			ExcludePolicyAndObjects: excludePolicyAndObjects,
			Force:                   force,
		}
	} else {
		commitCmd = commit.PanoramaCommit{
			Description:             description,
			Admins:                  admins,
			DeviceGroups:            deviceGroups,
			Templates:               templates,
			TemplateStacks:          templateStacks,
			WildfireAppliances:      wildfireAppliances,
			WildfireClusters:        wildfireClusters,
			LogCollectors:           logCollectors,
			LogCollectorGroups:      logCollectorGroups,
			ExcludeDeviceAndNetwork: excludeDeviceAndNetwork,
			ExcludeSharedObjects:    excludeSharedObjects,
			Force:                   force,
		}
	}

	cmd := &xmlapi.Commit{
		Command: commitCmd,
		Target:  o.client.GetTarget(),
	}

	var commitResp xmlapi.JobResponse

	_, _, err = o.client.Communicate(ctx, cmd, false, &commitResp)
	if err != nil {
		resp.Diagnostics.AddError("Failed to schedule a commit", err.Error())
		return
	}

	if commitResp.Id > 0 {
		err = o.client.WaitForJob(ctx, commitResp.Id, 2*time.Second, nil)
		if err != nil {
			resp.Diagnostics.AddError("Failed to wait for commit task to finish", err.Error())
			return
		}
	}

	if model.PushConfiguration.IsNull() || model.PushConfiguration.IsUnknown() {
		return
	}

	if isFirewall {
		resp.Diagnostics.AddError(
			"Push configuration not available on firewall",
			"The push_configuration parameter is only available on Panorama devices",
		)
		return
	}

	var pushConfig CommitActionPushConfigurationObject
	resp.Diagnostics.Append(model.PushConfiguration.As(ctx, &pushConfig, basetypes.ObjectAsOptions{})...)
	if resp.Diagnostics.HasError() {
		return
	}

	var pushDescription string
	var devices []string

	if !pushConfig.Description.IsNull() && !pushConfig.Description.IsUnknown() {
		pushDescription = pushConfig.Description.ValueString()
	} else {
		pushDescription = description
	}

	if !pushConfig.Devices.IsNull() && !pushConfig.Devices.IsUnknown() {
		resp.Diagnostics.Append(pushConfig.Devices.ElementsAs(ctx, &devices, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	resp.Diagnostics.Append(pushToDevices(
		ctx,
		o.client,
		pushToDevicesSpec{
			Type:                pushConfig.Type.ValueString(),
			Name:                pushConfig.Name.ValueString(),
			Description:         pushDescription,
			IncludeTemplate:     pushConfig.IncludeTemplate.ValueBool(),
			ForceTemplateValues: pushConfig.ForceTemplateValues.ValueBool(),
			Devices:             devices,
			FailOnError:         false, // Warnings only for post-commit push
		},
	)...)
}
