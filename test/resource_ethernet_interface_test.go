package provider_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

const ethernetInterface_Basic = `
variable "location" { type = any }
variable "create_template" { type = bool }
variable "prefix" { type = string }

resource "panos_template" "example" {
  count = var.create_template ? 1 : 0
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_ethernet_interface" "example" {
  depends_on = [ panos_template.example ]
  location = var.location

  name = "ethernet1/1"
  layer3 = {}
}
`

func TestAccEthernetInterface_Basic(t *testing.T) {
	t.Parallel()
	if os.Getenv("TF_ACC") != "1" {
		t.Skip("environment setup not complete")
	}

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	var location config.Variable
	var createTemplate config.Variable

	err := sdkClient.RetrieveSystemInfo(context.Background())
	if err != nil {
		panic(err)
	}
	firewall, err := sdkClient.IsFirewall()
	if err != nil {
		panic(err)
	}

	if firewall {
		location = config.ObjectVariable(map[string]config.Variable{
			"ngfw": config.ObjectVariable(map[string]config.Variable{}),
		})
		createTemplate = config.BoolVariable(false)
	} else {
		location = config.ObjectVariable(map[string]config.Variable{
			"template": config.ObjectVariable(map[string]config.Variable{
				"name": config.StringVariable(prefix),
			}),
		})
		createTemplate = config.BoolVariable(true)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: ethernetInterface_Basic,
				ConfigVariables: map[string]config.Variable{
					"location":        location,
					"create_template": createTemplate,
					"prefix":          config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_ethernet_interface.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact("ethernet1/1"),
					),
				},
			},
		},
	})
}

func TestAccEthernetInterface_Concurrent(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	location := config.ObjectVariable(map[string]config.Variable{
		"template": config.ObjectVariable(map[string]config.Variable{
			"name": config.StringVariable(prefix),
		}),
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: ethernetInterface_Concurrent_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
			},
		},
	})
}

const ethernetInterface_Concurrent_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }

  name = var.prefix
}

resource "panos_ethernet_interface" "example1" {
  depends_on = [panos_template.example]
  location = var.location

  name = "ethernet1/1"
  layer3 = {}
}

resource "panos_ethernet_interface" "example2" {
  depends_on = [panos_template.example]
  location = var.location

  name = "ethernet1/2"
  layer3 = {}
}

resource "panos_ethernet_interface" "example3" {
  depends_on = [panos_template.example]
  location = var.location

  name = "ethernet1/3"
  layer3 = {}
}

resource "panos_ethernet_interface" "example4" {
  depends_on = [panos_template.example]
  location = var.location

  name = "ethernet1/4"
  layer3 = {}
}

resource "panos_ethernet_interface" "example5" {
  depends_on = [panos_template.example]
  location = var.location

  name = "ethernet1/5"
  layer3 = {}
}

resource "panos_ethernet_interface" "example6" {
  depends_on = [panos_template.example]
  location = var.location

  name = "ethernet1/6"
  layer3 = {}
}
`

func TestAccPanosEthernetInterface_Layer3(t *testing.T) {
	t.Parallel()

	resName := "ethernet"
	templateName := "acc-codegen"
	interfaceName := "ethernet1/23"
	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: makePanosEthernetInterface_Layer3(resName),
				ConfigVariables: map[string]config.Variable{
					"name_suffix":     config.StringVariable(nameSuffix),
					"interface_name":  config.StringVariable(interfaceName),
					"template_name":   config.StringVariable(templateName),
					"ip_addr_netmask": config.StringVariable("1.1.1.1/32"),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_ethernet_interface."+resName,
						tfjsonpath.New("name"),
						knownvalue.StringExact("ethernet1/23"),
					),

					statecheck.ExpectKnownValue(
						"panos_ethernet_interface."+resName,
						tfjsonpath.
							New("layer3").
							AtMapKey("ips").
							AtSliceIndex(0).
							AtMapKey("name"),
						knownvalue.StringExact("1.1.1.1/32"),
					),
				},
			},
		},
	})
}

func makePanosEthernetInterface_Layer3(label string) string {
	configTpl := `
    variable "name_suffix" { type = string }
    variable "interface_name" { type = string }
    variable "ip_addr_netmask" { type = string }
    variable "template_name" { type = string }

    resource "panos_template" "acc_codegen_template" {
        name = "${var.template_name}-${var.name_suffix}"

        location = {
            panorama = {
                panorama_device = "localhost.localdomain"
            }
        }
    }

    resource "panos_ethernet_interface" "%s" {
      location = {
        template = {
          vsys = "vsys1"
          name = panos_template.acc_codegen_template.name
        }
      }


      name = var.interface_name

      layer3 = {
        lldp = {
          enable = true
        }

        mtu  = 1350
        ips  = [{ name = var.ip_addr_netmask }]

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
    `

	return fmt.Sprintf(configTpl, label)
}
