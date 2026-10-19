package provider_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango/xmlapi"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccDeviceGroup(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceGroupResourceTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_device_group.dg",
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-dg", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_device_group.dg",
						tfjsonpath.New("description"),
						knownvalue.StringExact("description"),
					),
					statecheck.ExpectKnownValue(
						"panos_device_group.dg",
						tfjsonpath.New("templates"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact(fmt.Sprintf("%s-tmpl", prefix)),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_device_group.dg",
						tfjsonpath.New("authorization_code"),
						knownvalue.StringExact("code"),
					),
				},
			},
		},
	})
}

const testAccDeviceGroupResourceTmpl = `
variable "prefix" { type = string }

resource "panos_template" "template" {
  location = { panorama = {} }

  name = format("%s-tmpl", var.prefix)
}

resource "panos_device_group" "dg" {
  location = { panorama = {} }

  name        = format("%s-dg", var.prefix)
  description = "description"

  templates = [ resource.panos_template.template.name ]

  authorization_code = "code"
}
`

// TestAccDeviceGroup_DevicesVsys verifies that a device group can be bound to
// specific vsys on a firewall, and that the binding is read back both by the
// resource and the data source.
func TestAccDeviceGroup_DevicesVsys(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	serialNumber := fmt.Sprintf("00%s", acctest.RandStringFromCharSet(13, "0123456789"))
	deviceXpath := fmt.Sprintf("/config/mgt-config/devices/entry[@name='%s']", serialNumber)

	// Panorama only accepts device group members that are registered as
	// managed devices, so register the serial number directly via the SDK.
	registerDevice := func() {
		cmd := &xmlapi.Config{
			Action:  "set",
			Xpath:   "/config/mgt-config/devices",
			Element: fmt.Sprintf("<entry name='%s'/>", serialNumber),
			Target:  sdkClient.GetTarget(),
		}

		if _, _, err := sdkClient.Communicate(context.TODO(), cmd, false, nil); err != nil {
			t.Fatalf("Failed to register managed device: %v", err)
		}
	}

	t.Cleanup(func() {
		cmd := &xmlapi.Config{
			Action: "delete",
			Xpath:  deviceXpath,
			Target: sdkClient.GetTarget(),
		}

		if _, _, err := sdkClient.Communicate(context.TODO(), cmd, false, nil); err != nil {
			t.Logf("Failed to remove managed device %s: %v", serialNumber, err)
		}
	})

	expectedDevices := knownvalue.ListExact([]knownvalue.Check{
		knownvalue.ObjectExact(map[string]knownvalue.Check{
			"name": knownvalue.StringExact(serialNumber),
			"vsys": knownvalue.ListExact([]knownvalue.Check{
				knownvalue.ObjectExact(map[string]knownvalue.Check{
					"name": knownvalue.StringExact("vsys1"),
				}),
				knownvalue.ObjectExact(map[string]knownvalue.Check{
					"name": knownvalue.StringExact("vsys2"),
				}),
			}),
		}),
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: registerDevice,
				Config:    deviceGroup_DevicesVsys_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":        config.StringVariable(prefix),
					"serial_number": config.StringVariable(serialNumber),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_device_group.dg",
						tfjsonpath.New("devices"),
						expectedDevices,
					),
					statecheck.ExpectKnownValue(
						"data.panos_device_group.dg",
						tfjsonpath.New("devices"),
						expectedDevices,
					),
				},
			},
		},
	})
}

const deviceGroup_DevicesVsys_Tmpl = `
variable "prefix" { type = string }
variable "serial_number" { type = string }

resource "panos_device_group" "dg" {
  location = { panorama = {} }

  name = format("%s-dg", var.prefix)

  devices = [{
    name = var.serial_number
    vsys = [{ name = "vsys1" }, { name = "vsys2" }]
  }]
}

data "panos_device_group" "dg" {
  location = { panorama = { panorama_device = "localhost.localdomain" } }

  name = panos_device_group.dg.name
}
`
