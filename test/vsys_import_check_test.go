package provider_test

import (
	"context"
	"encoding/xml"
	"fmt"

	"github.com/PaloAltoNetworks/pango/errors"
	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/xmlapi"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
)

type ImportType string

const (
	ImportTypeInterface     ImportType = "interface"
	ImportTypeVirtualRouter ImportType = "virtual-router"
)

type expectVsysImport struct {
	ResourceName     string     // e.g., "panos_virtual_router.test"
	VsysName         string     // e.g., "vsys1"
	ImportType       ImportType // "interface" or "virtual-router"
	ShouldBeImported bool       // true = expect import, false = expect no import
}

func ExpectVsysImportExists(resourceName, vsysName string, importType ImportType) *expectVsysImport {
	return &expectVsysImport{
		ResourceName:     resourceName,
		VsysName:         vsysName,
		ImportType:       importType,
		ShouldBeImported: true,
	}
}

func ExpectVsysImportAbsent(resourceName, vsysName string, importType ImportType) *expectVsysImport {
	return &expectVsysImport{
		ResourceName:     resourceName,
		VsysName:         vsysName,
		ImportType:       importType,
		ShouldBeImported: false,
	}
}

func (o *expectVsysImport) CheckState(ctx context.Context, req statecheck.CheckStateRequest, resp *statecheck.CheckStateResponse) {
	// Find resource in Terraform state
	var resource *tfjson.StateResource
	for _, r := range req.State.Values.RootModule.Resources {
		if r.Address == o.ResourceName {
			resource = r
			break
		}
	}

	if resource == nil {
		resp.Error = fmt.Errorf("resource %s not found in state", o.ResourceName)
		return
	}

	// Get the actual resource name (e.g., "ethernet1/1", "vr-1")
	resourceId, ok := resource.AttributeValues["name"].(string)
	if !ok || resourceId == "" {
		resp.Error = fmt.Errorf("resource %s has no name attribute", o.ResourceName)
		return
	}

	// Extract template name from Terraform state
	// location.template.name is nested: location -> template -> name
	locationMap, ok := resource.AttributeValues["location"].(map[string]interface{})
	if !ok {
		resp.Error = fmt.Errorf("resource %s has no location attribute", o.ResourceName)
		return
	}

	templateMap, ok := locationMap["template"].(map[string]interface{})
	if !ok {
		resp.Error = fmt.Errorf("resource %s has no location.template attribute", o.ResourceName)
		return
	}

	templateName, ok := templateMap["name"].(string)
	if !ok || templateName == "" {
		resp.Error = fmt.Errorf("resource %s has no location.template.name attribute", o.ResourceName)
		return
	}

	// Build XPath for template-based vsys import list
	xpath := []string{
		"config", "devices",
		util.AsEntryXpath("localhost.localdomain"),
		"template", util.AsEntryXpath(templateName),
		"config", "devices",
		util.AsEntryXpath("localhost.localdomain"),
		"vsys", util.AsEntryXpath(o.VsysName),
		"import", "network", string(o.ImportType),
	}

	// Make API call
	cmd := &xmlapi.Config{
		Action: "get",
		Xpath:  util.AsXpath(xpath),
	}

	bytes, _, err := sdkClient.Communicate(ctx, cmd, false, nil)
	if err != nil {
		// ObjectNotFound means import list is empty - this is acceptable for negative checks
		if errors.IsObjectNotFound(err) {
			if o.ShouldBeImported {
				resp.Error = fmt.Errorf("expected %s %q to be imported to vsys %q, but vsys has no %s imports",
					o.ImportType, resourceId, o.VsysName, o.ImportType)
			}
			// For negative check, ObjectNotFound is success (no imports = not imported)
			return
		}
		resp.Error = fmt.Errorf("failed to check vsys import state: %w", err)
		return
	}

	// Parse XML response based on import type
	var members []string
	if o.ImportType == ImportTypeVirtualRouter {
		var vrResponse struct {
			Members []struct {
				Name string `xml:",chardata"`
			} `xml:"result>virtual-router>member"`
		}
		if err := xml.Unmarshal(bytes, &vrResponse); err != nil {
			resp.Error = fmt.Errorf("failed to parse virtual-router import response: %w", err)
			return
		}
		for _, m := range vrResponse.Members {
			members = append(members, m.Name)
		}
	} else {
		var ifResponse struct {
			Members []struct {
				Name string `xml:",chardata"`
			} `xml:"result>interface>member"`
		}
		if err := xml.Unmarshal(bytes, &ifResponse); err != nil {
			resp.Error = fmt.Errorf("failed to parse interface import response: %w", err)
			return
		}
		for _, m := range ifResponse.Members {
			members = append(members, m.Name)
		}
	}

	// Check if resource is in import list
	found := false
	for _, member := range members {
		if member == resourceId {
			found = true
			break
		}
	}

	// Validate based on expectation
	if o.ShouldBeImported && !found {
		resp.Error = fmt.Errorf("expected %s %q to be imported to vsys %q, but it was not found in import list (found: %v)",
			o.ImportType, resourceId, o.VsysName, members)
	} else if !o.ShouldBeImported && found {
		resp.Error = fmt.Errorf("expected %s %q to NOT be imported to vsys %q, but it was found in import list",
			o.ImportType, resourceId, o.VsysName)
	}
}
