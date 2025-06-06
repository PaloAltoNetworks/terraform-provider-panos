package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccPanosAddressGroup(t *testing.T) {
	t.Parallel()

	resourceName := "dns_addresses"
	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)

	compareValuesDiffer := statecheck.CompareValue(compare.ValuesDiffer())

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: makeAddressGroupConfig(resourceName),
				ConfigVariables: map[string]config.Variable{
					"address_group_name":  config.StringVariable(resourceName),
					"name_suffix":         config.StringVariable(nameSuffix),
					"address_object_name": config.StringVariable("google-dns"),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_address_group."+resourceName,
						tfjsonpath.
							New("static").
							AtSliceIndex(0),
						knownvalue.StringExact("google-dns-"+nameSuffix),
					),
				},
			},
			{
				Config: makeAddressGroupConfig(resourceName),
				ConfigVariables: map[string]config.Variable{
					"address_group_name": config.StringVariable(resourceName),
					"name_suffix":        config.StringVariable(nameSuffix),
					"from_address_group": config.BoolVariable(true),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_address_group."+resourceName,
						tfjsonpath.
							New("static").
							AtSliceIndex(0),
						knownvalue.StringExact(fmt.Sprintf(
							"%s-base-%s",
							resourceName,
							nameSuffix,
						)),
					),
					compareValuesDiffer.AddStateValue(
						"panos_address_group."+resourceName,
						tfjsonpath.New("static"),
					),
				},
			},
		},
	})
}

func makeAddressGroupConfig(label string) string {
	confiTpl := `
    variable "name_suffix" { type = string }
    variable "address_group_name" { type = string }

    variable "address_object_name" {
        type = string
        default = "acct-google-dns"
    }
    variable "address_ip_netmask" {
        type = string
        default = "8.8.8.8/32"
    }

    variable "from_address_group" {
      type    = bool
      default = false
    }

    resource "panos_addresses" "google_dns_servers" {
      location = {
        shared = {}
      }

      addresses = {
        "${var.address_object_name}-${var.name_suffix}" = {
          ip_netmask = var.address_ip_netmask
        },
      }
    }

    resource "panos_address_group" "%s_base" {
      count = var.from_address_group ? 1 : 0
      location = {
        shared = {}
      }

      name   = "${var.address_group_name}-base-${var.name_suffix}"
      static = [for name, data in resource.panos_addresses.google_dns_servers.addresses : name]
    }

    resource "panos_address_group" "%s" {

      location = {
        shared = {}
      }

      name = "${var.address_group_name}-${var.name_suffix}"
      static = var.from_address_group ? (
        [panos_address_group.%s_base[0].name]
        ) : (
        [for name, data in resource.panos_addresses.google_dns_servers.addresses : name]
      )
    }
    `

	return fmt.Sprintf(confiTpl, label, label, label)
}
