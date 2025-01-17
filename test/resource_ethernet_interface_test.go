package provider_test

import (
	"context"
	"fmt"
	"testing"

	sdkErrors "github.com/PaloAltoNetworks/pango/errors"
	"github.com/PaloAltoNetworks/pango/network/interface/ethernet"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccPanosEthernetInterface_Layer3(t *testing.T) {
	t.Parallel()

	resName := "ethernet"
	templateName := "acc-codegen"
	interfaceName := "ethernet1/23"
	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		CheckDestroy: testAccCheckPanosEthernetInterfaceDestroy(
			interfaceName,
			fmt.Sprintf("%s-%s", templateName, nameSuffix),
		),
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

func testAccCheckPanosEthernetInterfaceDestroy(entryName, templateName string) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		api := ethernet.NewService(sdkClient)
		ctx := context.TODO()

		location := ethernet.NewTemplateLocation()
		location.Template.Template = templateName

		reply, err := api.Read(ctx, *location, entryName, "show")
		if err != nil && !sdkErrors.IsObjectNotFound(err) {
			return fmt.Errorf("reading ethernet entry via sdk: %v", err)
		}

		if reply != nil {
			if reply.EntryName() == entryName {
				return fmt.Errorf("ethernet object still exists: %s", entryName)
			}
		}

		return nil
	}
}
