package provider

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// vendorConfigs defines the required configuration keys for each MFA vendor type.
// All keys listed are mandatory when the vendor type is selected.
var vendorConfigs = map[string][]string{
	"duo-security-v2": {
		"duo-api-host",
		"duo-integration-key",
		"duo-secret-key",
		"duo-timeout",
		"duo-baseuri",
	},
	"okta-adaptive-v1": {
		"okta-api-host",
		"okta-baseuri",
		"okta-token",
		"okta-org",
		"okta-timeout",
	},
	"ping-identity-v1": {
		"ping-api-host",
		"ping-baseuri",
		"ping-token",
		"ping-org-alias",
		"ping-timeout",
	},
	"rsa-securid-access-v1": {
		"rsa-api-host",
		"rsa-baseuri",
		"rsa-accesskey",
		"rsa-accessid",
		"rsa-assurancepolicyid",
		"rsa-timeout",
	},
}

func (o *MfaServerProfileResource) ValidateConfigCustom(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data MfaServerProfileResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Skip validation if values are unknown (computed at plan time)
	if data.MfaVendorType.IsUnknown() || data.MfaConfig.IsUnknown() {
		return
	}

	// Check if mfa_config is provided without mfa_vendor_type
	vendorTypeEmpty := data.MfaVendorType.IsNull() || data.MfaVendorType.ValueString() == ""
	configProvided := !data.MfaConfig.IsNull() && len(data.MfaConfig.Elements()) > 0

	if vendorTypeEmpty && configProvided {
		resp.Diagnostics.AddAttributeError(
			path.Root("mfa_config"),
			"Configuration Requires Vendor Type",
			"The 'mfa_config' attribute cannot be set without specifying 'mfa_vendor_type'. "+
				"MFA configuration keys are vendor-specific and require a vendor type to be defined.",
		)
		return
	}

	// Skip validation if vendor type is not set (allows certificate-only profiles)
	if vendorTypeEmpty {
		return
	}

	vendorType := data.MfaVendorType.ValueString()

	// Extract config items from the list
	var configItems []MfaServerProfileResourceMfaConfigObject
	diags := data.MfaConfig.ElementsAs(ctx, &configItems, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build set of provided configuration keys
	providedKeys := make(map[string]bool)
	for _, item := range configItems {
		if !item.Name.IsNull() && !item.Name.IsUnknown() {
			providedKeys[item.Name.ValueString()] = true
		}
	}

	// Get required keys for this vendor type
	requiredKeys, vendorExists := vendorConfigs[vendorType]
	if !vendorExists {
		// Build list of valid vendor types for error message
		validVendors := make([]string, 0, len(vendorConfigs))
		for vendor := range vendorConfigs {
			validVendors = append(validVendors, vendor)
		}
		sort.Strings(validVendors)

		resp.Diagnostics.AddAttributeError(
			path.Root("mfa_vendor_type"),
			"Invalid MFA Vendor Type",
			fmt.Sprintf(
				"MFA vendor type '%s' is not supported. Valid vendor types are: %s.",
				vendorType,
				strings.Join(validVendors, ", "),
			),
		)
		return
	}

	// Build required keys set for O(1) lookup
	requiredKeySet := make(map[string]bool)
	for _, key := range requiredKeys {
		requiredKeySet[key] = true
	}

	// Check for missing required keys
	var missingKeys []string
	for _, requiredKey := range requiredKeys {
		if !providedKeys[requiredKey] {
			missingKeys = append(missingKeys, requiredKey)
		}
	}

	if len(missingKeys) > 0 {
		sort.Strings(missingKeys)

		resp.Diagnostics.AddAttributeError(
			path.Root("mfa_config"),
			"Missing Required Configuration Keys",
			fmt.Sprintf(
				"MFA vendor type '%s' requires the following configuration keys that are missing: %s. "+
					"All required keys must be provided for the vendor configuration to be valid.",
				vendorType,
				strings.Join(missingKeys, ", "),
			),
		)
	}

	// Check for invalid keys (keys that don't belong to this vendor)
	var invalidKeys []string
	for providedKey := range providedKeys {
		if !requiredKeySet[providedKey] {
			invalidKeys = append(invalidKeys, providedKey)
		}
	}

	if len(invalidKeys) > 0 {
		sort.Strings(invalidKeys)

		resp.Diagnostics.AddAttributeError(
			path.Root("mfa_config"),
			"Invalid Configuration Keys",
			fmt.Sprintf(
				"MFA vendor type '%s' does not support the following configuration keys: %s. "+
					"Only these keys are valid for this vendor: %s.",
				vendorType,
				strings.Join(invalidKeys, ", "),
				strings.Join(requiredKeys, ", "),
			),
		)
	}
}
