package provider

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/commit"
	"github.com/PaloAltoNetworks/pango/xmlapi"
)

type pushToDevicesSpec struct {
	Type                string
	Name                string
	Description         string
	IncludeTemplate     bool
	ForceTemplateValues bool
	Devices             []string
	FailOnError         bool
}

// pushToDevices is a shared helper function that executes a push-to-devices operation.
// It validates the push configuration, creates the appropriate commit-all command,
// and executes it on the device.
//
// Parameters:
//   - ctx: Context for the operation
//   - client: PAN-OS API client
//   - spec: Push configuration specification
//
// Returns:
//   - diag.Diagnostics: Contains any errors or warnings from the operation
func pushToDevices(
	ctx context.Context,
	client *pango.Client,
	spec pushToDevicesSpec,
) diag.Diagnostics {
	var diags diag.Diagnostics

	// Validate type-specific constraints
	switch spec.Type {
	case "device_group":
		// include_template and force_template_values are valid
	case "template", "template_stack":
		// Only force_template_values is valid
		if spec.IncludeTemplate {
			diags.AddWarning(
				"include_template ignored",
				"The include_template parameter is only applicable for device_group type and will be ignored",
			)
			spec.IncludeTemplate = false
		}
	case "log_collector_group", "wildfire_appliance", "wildfire_cluster":
		// Neither include_template nor force_template_values are valid
		if spec.IncludeTemplate {
			diags.AddWarning(
				"include_template ignored",
				"The include_template parameter is not applicable for this push type and will be ignored",
			)
			spec.IncludeTemplate = false
		}
		if spec.ForceTemplateValues {
			diags.AddWarning(
				"force_template_values ignored",
				"The force_template_values parameter is not applicable for this push type and will be ignored",
			)
			spec.ForceTemplateValues = false
		}
	default:
		if spec.FailOnError {
			diags.AddError(
				"Invalid push configuration type",
				"Valid types are: device_group, template, template_stack, log_collector_group, wildfire_appliance, wildfire_cluster",
			)
		} else {
			diags.AddWarning(
				"Invalid push configuration type",
				"Valid types are: device_group, template, template_stack, log_collector_group, wildfire_appliance, wildfire_cluster",
			)
		}
		return diags
	}

	var commitType string
	switch spec.Type {
	case "device_group":
		commitType = commit.TypeDeviceGroup
	case "template":
		commitType = commit.TypeTemplate
	case "template_stack":
		commitType = commit.TypeTemplateStack
	case "log_collector_group":
		commitType = commit.TypeLogCollectorGroup
	case "wildfire_appliance":
		commitType = commit.TypeWildfireAppliance
	case "wildfire_cluster":
		commitType = commit.TypeWildfireCluster
	}

	commitAllCmd := commit.PanoramaCommitAll{
		Type:                commitType,
		Name:                spec.Name,
		Description:         spec.Description,
		IncludeTemplate:     spec.IncludeTemplate,
		ForceTemplateValues: spec.ForceTemplateValues,
		Devices:             spec.Devices,
	}

	pushCmd := &xmlapi.Commit{
		Command: commitAllCmd,
		Target:  client.GetTarget(),
	}

	var pushResp xmlapi.JobResponse

	_, _, err := client.Communicate(ctx, pushCmd, false, &pushResp)
	if err != nil {
		if spec.FailOnError {
			diags.AddError("Failed to schedule push to devices", err.Error())
		} else {
			diags.AddWarning(
				"Failed to schedule push to devices",
				"The commit operation completed successfully, but the push operation failed: "+err.Error(),
			)
		}
		return diags
	}

	err = client.WaitForJob(ctx, pushResp.Id, 2*time.Second, nil)
	if err != nil {
		if spec.FailOnError {
			diags.AddError("Failed to wait for push task to finish", err.Error())
		} else {
			diags.AddWarning(
				"Failed to complete push to devices",
				"The commit operation completed successfully, but the push operation failed: "+err.Error(),
			)
		}
		return diags
	}

	return diags
}

func (o *PushToDevicesAction) InvokeCustom(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var model PushToDevicesActionModel
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

	if isFirewall {
		resp.Diagnostics.AddError(
			"Push to devices not available on firewall",
			"The push_to_devices action is only available on Panorama devices",
		)
		return
	}

	var pushDescription string
	var devices []string

	if !model.Description.IsNull() && !model.Description.IsUnknown() {
		pushDescription = model.Description.ValueString()
	}

	if !model.Devices.IsNull() && !model.Devices.IsUnknown() {
		resp.Diagnostics.Append(model.Devices.ElementsAs(ctx, &devices, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Execute the push operation
	resp.Diagnostics.Append(pushToDevices(
		ctx,
		o.client,
		pushToDevicesSpec{
			Type:                model.Type.ValueString(),
			Name:                model.Name.ValueString(),
			Description:         pushDescription,
			IncludeTemplate:     model.IncludeTemplate.ValueBool(),
			ForceTemplateValues: model.ForceTemplateValues.ValueBool(),
			Devices:             devices,
			FailOnError:         true, // Hard errors for standalone push action
		},
	)...)
}
