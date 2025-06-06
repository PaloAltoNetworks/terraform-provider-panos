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

func TestAccZone(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	suffix := "tmpl1"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: zoneResourceTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":          config.StringVariable(prefix),
					"template_suffix": config.StringVariable(suffix),
					"interface_type":  config.StringVariable("layer3"),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_zone.zone",
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-zone", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_zone.zone",
						tfjsonpath.New("enable_device_identification"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_zone.zone",
						tfjsonpath.New("enable_user_identification"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_zone.zone",
						tfjsonpath.New("network").
							AtMapKey("enable_packet_buffer_protection"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_zone.zone",
						tfjsonpath.New("network").
							AtMapKey("layer3"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("ethernet1/1"),
							knownvalue.StringExact("ethernet1/2"),
						}),
					),
				},
			},
		},
	})

	suffix = "tmpl2"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: zoneResourceTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":          config.StringVariable(prefix),
					"template_suffix": config.StringVariable(suffix),
					"interface_type":  config.StringVariable("layer2"),
				},
				ConfigStateChecks: []statecheck.StateCheck{},
			},
		},
	})

	// suffix = "tmpl3"
	// resource.Test(t, resource.TestCase{
	// 	PreCheck:                 func() { testAccPreCheck(t) },
	// 	ProtoV6ProviderFactories: testAccProviders,
	// 	Steps: []resource.TestStep{
	// 		{
	// 			Config: zoneResourceTmpl,
	// 			ConfigVariables: map[string]config.Variable{
	// 				"prefix":          config.StringVariable(prefix),
	// 				"template_suffix": config.StringVariable(suffix),
	// 				"interface_type":  config.StringVariable("external"),
	// 			},
	// 			ConfigStateChecks: []statecheck.StateCheck{},
	// 		},
	// 	},
	// })

	suffix = "tmpl4"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: zoneResourceTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":          config.StringVariable(prefix),
					"template_suffix": config.StringVariable(suffix),
					"interface_type":  config.StringVariable("tap"),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_zone.zone",
						tfjsonpath.New("network").
							AtMapKey("tap"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("ethernet1/1"),
							knownvalue.StringExact("ethernet1/2"),
						}),
					),
				},
			},
		},
	})

	// suffix = "tmpl5"
	// resource.Test(t, resource.TestCase{
	// 	PreCheck:                 func() { testAccPreCheck(t) },
	// 	ProtoV6ProviderFactories: testAccProviders,
	// 	Steps: []resource.TestStep{
	// 		{
	// 			Config: zoneResourceTmpl,
	// 			ConfigVariables: map[string]config.Variable{
	// 				"prefix":          config.StringVariable(prefix),
	// 				"template_suffix": config.StringVariable(suffix),
	// 				"interface_type":  config.StringVariable("tunnel"),
	// 			},
	// 			ConfigStateChecks: []statecheck.StateCheck{
	// 				statecheck.ExpectKnownValue(
	// 					"panos_zone.zone",
	// 					tfjsonpath.New("network").
	// 						AtMapKey("tunnel"),
	// 					knownvalue.ObjectExact(map[string]knownvalue.Check{}),
	// 			},
	// 		},
	// 	},
	// })
}

const zoneResourceTmpl = `
variable "prefix" { type = string }
variable "template_suffix" { type = string }
variable "interface_type" { type = string }

locals {
  interfaces = {
    layer2 = var.interface_type == "layer2" ? ["ethernet1/1", "ethernet1/2"] : null
    layer3 = var.interface_type == "layer3" ? ["ethernet1/1", "ethernet1/2"] : null
    external = var.interface_type == "external" ? ["ethernet1/1", "ethernet1/2"] : null
    tap = var.interface_type == "tap" ? ["ethernet1/1", "ethernet1/2"] : null
    tunnel = var.interface_type == "tunnel" ? {} : null
  }

  create_iface = contains(["layer2", "layer3", "tap"], var.interface_type) ? true : false

  template_name = format("%s-%s", var.prefix, var.template_suffix)
  network_common = {
      enable_packet_buffer_protection = true
      # log_setting                      = ["log-setting"]
      # zone_protection_profile          = "zone-protection-profile"
  }
  network = merge(local.network_common, local.interfaces)
}

resource "panos_template" "template" {
  location = { panorama = {} }

  name = local.template_name
}


resource "panos_ethernet_interface" "iface1" {
  count = local.create_iface == true ? 1 : 0
  location = { template = { name = resource.panos_template.template.name, vsys = "vsys1" }}

  name = "ethernet1/1"

  tap = var.interface_type == "tap" ? {} : null
  layer2 = var.interface_type == "layer2" ? {} : null
  layer3 = var.interface_type == "layer3" ? {} : null
}

resource "panos_ethernet_interface" "iface2" {
  count = local.create_iface == true ? 1 : 0

  location = { template = { name = resource.panos_template.template.name, vsys = "vsys1" }}

  name = "ethernet1/2"

  tap = var.interface_type == "tap" ? {} : null
  layer2 = var.interface_type == "layer2" ? {} : null
  layer3 = var.interface_type == "layer3" ? {} : null
}

resource "panos_zone" "zone" {
  depends_on = [
    resource.panos_ethernet_interface.iface1, resource.panos_ethernet_interface.iface2
  ]

  location = { template = { name = resource.panos_template.template.name }}

  name = format("%s-zone", var.prefix)

  device_acl = {
    # exclude_list = ["device-1"]
    # include_list = ["device-2"]
  }

  enable_device_identification = true
  enable_user_identification   = true

  network = local.network
}
`
