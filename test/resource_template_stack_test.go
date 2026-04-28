package provider_test

import (
	"context"
	"encoding/xml"
	"fmt"
	"regexp"
	"testing"

	"github.com/PaloAltoNetworks/pango/generic"
	"github.com/PaloAltoNetworks/pango/panorama/template_stack"
	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/xmlapi"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccTemplateStack_Basic(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	location := config.ObjectVariable(map[string]config.Variable{
		"panorama": config.ObjectVariable(map[string]config.Variable{}),
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: templateStack_Basic_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_template_stack.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_template_stack.example",
						tfjsonpath.New("description"),
						knownvalue.StringExact("Template stack description"),
					),
					statecheck.ExpectKnownValue(
						"panos_template_stack.example",
						tfjsonpath.New("templates"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact(fmt.Sprintf("%s-template1", prefix)),
							knownvalue.StringExact(fmt.Sprintf("%s-template2", prefix)),
						}),
					),
				},
			},
		},
	})
}

const templateStack_Basic_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "template1" {
  location = var.location
  name = "${var.prefix}-template1"
}

resource "panos_template" "template2" {
  location = var.location
  name = "${var.prefix}-template2"
}

resource "panos_template_stack" "example" {
  location = var.location
  name = var.prefix
  description = "Template stack description"
  templates = [
    panos_template.template1.name,
    panos_template.template2.name
  ]
}
`

func TestAccTemplateStack_Devices(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	// Generate serial numbers for firewall devices
	suffix1 := acctest.RandStringFromCharSet(13, "0123456789")
	serialNumber1 := fmt.Sprintf("00%s", suffix1)
	suffix2 := acctest.RandStringFromCharSet(13, "0123456789")
	serialNumber2 := fmt.Sprintf("00%s", suffix2)

	location := config.ObjectVariable(map[string]config.Variable{
		"panorama": config.ObjectVariable(map[string]config.Variable{}),
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: templateStack_Devices_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":          config.StringVariable(prefix),
					"location":        location,
					"serial_number_1": config.StringVariable(serialNumber1),
					"serial_number_2": config.StringVariable(serialNumber2),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_template_stack.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_template_stack.example",
						tfjsonpath.New("devices"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name": knownvalue.StringExact(serialNumber1),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name": knownvalue.StringExact(serialNumber2),
							}),
						}),
					),
				},
			},
		},
	})
}

const templateStack_Devices_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }
variable "serial_number_1" { type = string }
variable "serial_number_2" { type = string }

resource "panos_firewall_device" "device1" {
  location = var.location
  name = var.serial_number_1
  hostname = "fw1.example.com"
  ip = "192.0.2.1"
}

resource "panos_firewall_device" "device2" {
  location = var.location
  name = var.serial_number_2
  hostname = "fw2.example.com"
  ip = "192.0.2.2"
}

resource "panos_template_stack" "example" {
  location = var.location
  name = var.prefix
  devices = [
    { name = panos_firewall_device.device1.name },
    { name = panos_firewall_device.device2.name }
  ]
}
`

// TestAccTemplateStack_DefaultVsys verifies that default_vsys can be set
// during the initial create step. The template is created with default_vsys
// first (which triggers its hooks to create the vsys), then the template-stack
// references that template and sets its own default_vsys.
func TestAccTemplateStack_DefaultVsys(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	location := config.ObjectVariable(map[string]config.Variable{
		"panorama": config.ObjectVariable(map[string]config.Variable{}),
	})

	configVars := map[string]config.Variable{
		"prefix":   config.StringVariable(prefix),
		"location": location,
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:          templateStack_DefaultVsys_Tmpl,
				ConfigVariables: configVars,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_template_stack.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_template_stack.example",
						tfjsonpath.New("default_vsys"),
						knownvalue.StringExact("vsys1"),
					),
				},
			},
		},
	})
}

const templateStack_DefaultVsys_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "test" {
  location = var.location
  name = "${var.prefix}-template"
  default_vsys = "vsys1"
}

resource "panos_template_stack" "example" {
  location = var.location
  name = var.prefix
  default_vsys = "vsys1"
  templates = [panos_template.test.name]
}
`

// TestAccTemplateStack_DefaultVsysNoTemplateVsys verifies that setting
// default_vsys on a template-stack fails when the referenced template does not
// have the vsys created (i.e. template is created without default_vsys).
func TestAccTemplateStack_DefaultVsysNoTemplateVsys(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	location := config.ObjectVariable(map[string]config.Variable{
		"panorama": config.ObjectVariable(map[string]config.Variable{}),
	})

	configVars := map[string]config.Variable{
		"prefix":   config.StringVariable(prefix),
		"location": location,
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:          templateStack_DefaultVsysNoTemplateVsys_Tmpl,
				ConfigVariables: configVars,
				ExpectError: regexp.MustCompile(`Error in create`),
			},
		},
	})
}

const templateStack_DefaultVsysNoTemplateVsys_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "test" {
  location = var.location
  name = "${var.prefix}-template"
}

resource "panos_template_stack" "example" {
  location = var.location
  name = var.prefix
  default_vsys = "vsys1"
  templates = [panos_template.test.name]
}
`

// TestAccTemplateStack_Complete creates a template stack with all fields
// including default_vsys set directly during the initial create step.
func TestAccTemplateStack_Complete(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	// Generate serial numbers for firewall devices
	suffix1 := acctest.RandStringFromCharSet(13, "0123456789")
	serialNumber1 := fmt.Sprintf("00%s", suffix1)
	suffix2 := acctest.RandStringFromCharSet(13, "0123456789")
	serialNumber2 := fmt.Sprintf("00%s", suffix2)

	location := config.ObjectVariable(map[string]config.Variable{
		"panorama": config.ObjectVariable(map[string]config.Variable{}),
	})

	configVars := map[string]config.Variable{
		"prefix":          config.StringVariable(prefix),
		"location":        location,
		"serial_number_1": config.StringVariable(serialNumber1),
		"serial_number_2": config.StringVariable(serialNumber2),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:          templateStack_Complete_Tmpl,
				ConfigVariables: configVars,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_template_stack.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_template_stack.example",
						tfjsonpath.New("description"),
						knownvalue.StringExact("Complete template stack"),
					),
					statecheck.ExpectKnownValue(
						"panos_template_stack.example",
						tfjsonpath.New("devices"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name": knownvalue.StringExact(serialNumber1),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name": knownvalue.StringExact(serialNumber2),
							}),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_template_stack.example",
						tfjsonpath.New("default_vsys"),
						knownvalue.StringExact("vsys1"),
					),
				},
			},
		},
	})
}

const templateStack_Complete_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }
variable "serial_number_1" { type = string }
variable "serial_number_2" { type = string }

resource "panos_template" "test" {
  location = var.location
  name = "${var.prefix}-template"
  default_vsys = "vsys1"
}

resource "panos_firewall_device" "device1" {
  location = var.location
  name = var.serial_number_1
  hostname = "fw1.example.com"
  ip = "192.0.2.1"
}

resource "panos_firewall_device" "device2" {
  location = var.location
  name = var.serial_number_2
  hostname = "fw2.example.com"
  ip = "192.0.2.2"
}

resource "panos_template_stack" "example" {
  location = var.location
  name = var.prefix
  description = "Complete template stack"
  templates = [panos_template.test.name]
  devices = [
    { name = panos_firewall_device.device1.name },
    { name = panos_firewall_device.device2.name }
  ]
  default_vsys = "vsys1"
}
`

// TestAccTemplateStack_DeviceVariablePreservation verifies that updating a
// template stack (e.g. changing description) does not delete per-device
// variable overrides stored as sub-elements of device entries.
//
// This reproduces a reported bug where the PUT request sends device entries
// without their child elements, causing PAN-OS to delete per-device template
// variable overrides.
func TestAccTemplateStack_DeviceVariablePreservation(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	suffix := acctest.RandStringFromCharSet(13, "0123456789")
	serialNumber := fmt.Sprintf("00%s", suffix)

	stackName := fmt.Sprintf("%s-stack", prefix)

	location := config.ObjectVariable(map[string]config.Variable{
		"panorama": config.ObjectVariable(map[string]config.Variable{}),
	})

	configVars := map[string]config.Variable{
		"prefix":        config.StringVariable(prefix),
		"location":      location,
		"serial_number": config.StringVariable(serialNumber),
	}

	// Build the xpath for direct SDK access since Read/Update have a bug
	// with name formatting. Use ReadWithXpath/UpdateWithXpath instead.
	loc := template_stack.Location{
		Panorama: &template_stack.PanoramaLocation{
			PanoramaDevice: "localhost.localdomain",
		},
	}
	stackXpathParts, err := loc.XpathWithComponents(sdkClient.Versioning(), util.AsEntryXpath(stackName))
	if err != nil {
		t.Fatalf("Failed to build template stack xpath: %v", err)
	}
	stackXpath := util.AsXpath(stackXpathParts)

	// injectDeviceVariable uses a direct API call to add a per-device variable
	// override to the device entry, simulating a user setting local values
	// via the Panorama UI.
	//
	// We use sdkClient.Communicate directly because the SDK's UpdateWithXpath
	// uses SpecMatches which doesn't compare Misc fields, so it would skip
	// the update.
	injectDeviceVariable := func() {
		// Build the xpath to the specific device entry within the template stack.
		deviceXpath := fmt.Sprintf("%s/devices/entry[@name='%s']", stackXpath, serialNumber)

		// Override the stack-level template variable with a per-device value.
		// The variable name must match the existing template variable.
		varName := fmt.Sprintf("$%s-var", prefix)
		variableXml := generic.Xml{
			XMLName: xml.Name{Local: "variable"},
			Nodes: []generic.Xml{
				{
					XMLName: xml.Name{Local: "entry"},
					Name:    &varName,
					Nodes: []generic.Xml{
						{
							XMLName: xml.Name{Local: "type"},
							Nodes: []generic.Xml{
								{
									XMLName: xml.Name{Local: "ip-netmask"},
									Text:    []byte("10.0.0.1/24"),
								},
							},
						},
					},
				},
			},
		}

		cmd := &xmlapi.Config{
			Action:  "set",
			Xpath:   deviceXpath,
			Element: variableXml,
			Target:  sdkClient.GetTarget(),
		}

		if _, _, err := sdkClient.Communicate(context.TODO(), cmd, false, nil); err != nil {
			t.Fatalf("Failed to inject per-device variable: %v", err)
		}
	}

	// checkDeviceVariableExists verifies that the per-device variable data
	// survived the Terraform update by reading the template stack via the
	// pango SDK and checking the device entry's Misc field.
	checkDeviceVariableExists := func(s *terraform.State) error {
		svc := template_stack.NewService(sdkClient)

		entry, err := svc.ReadWithXpath(context.TODO(), stackXpath, "get")
		if err != nil {
			return fmt.Errorf("failed to read template stack: %v", err)
		}

		for _, device := range entry.Devices {
			if device.Name == serialNumber {
				if len(device.Misc) == 0 {
					return fmt.Errorf(
						"device %s lost its Misc data: per-device variable overrides were deleted during template stack update",
						serialNumber,
					)
				}
				return nil
			}
		}
		return fmt.Errorf("device %s not found in template stack after update", serialNumber)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:          templateStack_DeviceVarPreservation_Step1_Tmpl,
				ConfigVariables: configVars,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_template_stack.test",
						tfjsonpath.New("description"),
						knownvalue.StringExact("Original description"),
					),
					statecheck.ExpectKnownValue(
						"panos_template_stack.test",
						tfjsonpath.New("devices"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name": knownvalue.StringExact(serialNumber),
							}),
						}),
					),
				},
			},
			{
				PreConfig:       injectDeviceVariable,
				Config:          templateStack_DeviceVarPreservation_Step2_Tmpl,
				ConfigVariables: configVars,
				Check:           checkDeviceVariableExists,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_template_stack.test",
						tfjsonpath.New("description"),
						knownvalue.StringExact("Updated description"),
					),
					statecheck.ExpectKnownValue(
						"panos_template_stack.test",
						tfjsonpath.New("devices"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name": knownvalue.StringExact(serialNumber),
							}),
						}),
					),
				},
			},
		},
	})
}

const templateStack_DeviceVarPreservation_Step1_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }
variable "serial_number" { type = string }

resource "panos_template" "test" {
  location = var.location
  name = "${var.prefix}-template"
}

resource "panos_firewall_device" "test" {
  location = var.location
  name = var.serial_number
  hostname = "fw-devvar.example.com"
  ip = "192.0.2.10"
}

resource "panos_template_stack" "test" {
  location = var.location
  name = "${var.prefix}-stack"
  description = "Original description"
  templates = [panos_template.test.name]
  devices = [{ name = panos_firewall_device.test.name }]
}

resource "panos_template_variable" "test" {
  location = { template_stack = { name = panos_template_stack.test.name } }
  name = format("$%s-var", var.prefix)
  type = { ip_netmask = "10.0.0.0/24" }
}
`

const templateStack_DeviceVarPreservation_Step2_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }
variable "serial_number" { type = string }

resource "panos_template" "test" {
  location = var.location
  name = "${var.prefix}-template"
}

resource "panos_firewall_device" "test" {
  location = var.location
  name = var.serial_number
  hostname = "fw-devvar.example.com"
  ip = "192.0.2.10"
}

resource "panos_template_stack" "test" {
  location = var.location
  name = "${var.prefix}-stack"
  description = "Updated description"
  templates = [panos_template.test.name]
  devices = [{ name = panos_firewall_device.test.name }]
}

resource "panos_template_variable" "test" {
  location = { template_stack = { name = panos_template_stack.test.name } }
  name = format("$%s-var", var.prefix)
  type = { ip_netmask = "10.0.0.0/24" }
}
`
