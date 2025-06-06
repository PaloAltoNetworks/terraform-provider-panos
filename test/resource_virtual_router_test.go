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

func TestAccPanosVirtualRouter_RequiredInputs(t *testing.T) {
	t.Parallel()

	resName := "vr_profile"
	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	templateName := "acc-vrouter"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: makePanosVirtualRouterConfig(resName),
				ConfigVariables: map[string]config.Variable{
					"name_suffix":   config.StringVariable(nameSuffix),
					"router_name":   config.StringVariable(resName),
					"template_name": config.StringVariable(templateName),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_virtual_router."+resName,
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf(
							"%s-%s",
							resName,
							nameSuffix,
						)),
					),
				},
			},
		},
	})
}

func TestAccPanosVirtualRouter_WithEthernetInterface(t *testing.T) {
	t.Parallel()

	resName := "vr_profile_if"
	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	templateName := "acc-vrouter"
	interfaceName := "ethernet1/41"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: makePanosVirtualRouterConfig(resName),
				ConfigVariables: map[string]config.Variable{
					"name_suffix":      config.StringVariable(nameSuffix),
					"router_name":      config.StringVariable(resName),
					"with_interface":   config.BoolVariable(true),
					"ethernet_if_name": config.StringVariable(interfaceName),
					"template_name":    config.StringVariable(templateName),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_virtual_router."+resName,
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf(
							"%s-%s",
							resName,
							nameSuffix,
						)),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router."+resName,
						tfjsonpath.
							New("interfaces").
							AtSliceIndex(0),
						knownvalue.StringExact("ethernet1/41"),
					),
				},
			},
		},
	})
}

func makePanosVirtualRouterConfig(label string) string {
	configTpl := `
    variable "name_suffix" { type = string }
    variable "router_name" { type = string }
    variable "template_name" { type = string }

    variable "with_interface" {
        type = bool
        default = false
    }

    variable "ethernet_if_name" {
        type = string
        default = "ethernet1/40"
    }

    resource "panos_template" "template" {
        name = "${var.template_name}-${var.name_suffix}"

        location = {
            panorama = {
                panorama_device = "localhost.localdomain"
            }
        }
    }

    resource "panos_ethernet_interface" "ethernet" {
      count = var.with_interface ? 1 : 0
      location = {
        template = {
          vsys = "vsys1"
          name = panos_template.template.name
        }
      }


      name = "${var.ethernet_if_name}"

      layer3 = {
        lldp = {
          enable = true
        }
        mtu  = 1350
        ips  = [{ name = "1.1.1.1/32" }]

        ipv6 = {
          enabled = true
          addresses = [
            {
              advertise = {
                enable         = true
                valid_lifetime = "1000000"
              },
              name                = "::1",
              enable_on_interface = true
            },
          ]
        }
      }
    }

    resource "panos_virtual_router" "%s" {
      location = {
       template = {
          name = panos_template.template.name
        }
      }

      name = "${var.router_name}-${var.name_suffix}"

      interfaces = var.with_interface ? [panos_ethernet_interface.ethernet[0].name] : []

      administrative_distances = {
        static = 15
      }
    }
    `

	return fmt.Sprintf(configTpl, label)
}
