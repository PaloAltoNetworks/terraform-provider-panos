
package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccPanosVirtualWire(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: virtualWireResourceTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":       config.StringVariable(prefix),
					"name":         config.StringVariable("vw-1"),
					"interface1":   config.StringVariable("ethernet1/1"),
					"interface2":   config.StringVariable("ethernet1/2"),
					"lsp_enable":   config.BoolVariable(false),
					"mf_enable":    config.BoolVariable(false),
					"tag_allowed":  config.StringVariable(""),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_virtual_wire.vw",
						tfjsonpath.New("name"),
						knownvalue.StringExact("vw-1"),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_wire.vw",
						tfjsonpath.New("interface1"),
						knownvalue.StringExact("ethernet1/1"),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_wire.vw",
						tfjsonpath.New("interface2"),
						knownvalue.StringExact("ethernet1/2"),
					),
				},
			},
		},
	})
}

func TestAccPanosVirtualWireWithOptions(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: virtualWireResourceTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":       config.StringVariable(prefix),
					"name":         config.StringVariable("vw-2"),
					"interface1":   config.StringVariable("ethernet1/3"),
					"interface2":   config.StringVariable("ethernet1/4"),
					"lsp_enable":   config.BoolVariable(true),
					"mf_enable":    config.BoolVariable(true),
					"tag_allowed":  config.StringVariable("100-200"),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_virtual_wire.vw",
						tfjsonpath.New("name"),
						knownvalue.StringExact("vw-2"),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_wire.vw",
						tfjsonpath.New("link_state_pass_through").AtMapKey("enable"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_wire.vw",
						tfjsonpath.New("multicast_firewalling").AtMapKey("enable"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_wire.vw",
						tfjsonpath.New("tag_allowed"),
						knownvalue.StringExact("100-200"),
					),
				},
			},
		},
	})
}

const virtualWireResourceTmpl = `
variable "prefix" { type = string }
variable "name" { type = string }
variable "interface1" { type = string }
variable "interface2" { type = string }
variable "lsp_enable" { type = bool }
variable "mf_enable" { type = bool }
variable "tag_allowed" { type = string }

resource "panos_template" "template" {
  location = { panorama = {} }
  name     = format("%s-tmpl", var.prefix)
}

resource "panos_ethernet_interface" "iface1" {
  location = { template = { name = resource.panos_template.template.name, vsys = "vsys1" } }
  name     = var.interface1
  virtual_wire = {}
}

resource "panos_ethernet_interface" "iface2" {
  location = { template = { name = resource.panos_template.template.name, vsys = "vsys1" } }
  name     = var.interface2
  virtual_wire = {}
}

resource "panos_virtual_wire" "vw" {
  depends_on = [
    resource.panos_ethernet_interface.iface1, resource.panos_ethernet_interface.iface2
  ]
  location = { template = { name = resource.panos_template.template.name } }

  name                      = var.name
  interface1                = var.interface1
  interface2                = var.interface2
  link_state_pass_through = var.lsp_enable ? {
    enable = true
  } : null
  multicast_firewalling = var.mf_enable ? {
    enable = true
  } : null
  tag_allowed               = var.tag_allowed != "" ? var.tag_allowed : null
}
`
