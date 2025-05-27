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

func TestAccTunnelInterface(t *testing.T) {
	t.Parallel()

	interfaceName := "tunnel.1"
	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: tunnelInterfaceResource1,
				ConfigVariables: map[string]config.Variable{
					"prefix":         config.StringVariable(prefix),
					"interface_name": config.StringVariable(interfaceName),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_tunnel_interface.iface",
						tfjsonpath.New("name"),
						knownvalue.StringExact("tunnel.1"),
					),
					statecheck.ExpectKnownValue(
						"panos_tunnel_interface.iface",
						tfjsonpath.New("comment"),
						knownvalue.StringExact("tunnel interface comment"),
					),
					statecheck.ExpectKnownValue(
						"panos_tunnel_interface.iface",
						tfjsonpath.New("interface_management_profile"),
						knownvalue.StringExact(fmt.Sprintf("%s-profile", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_tunnel_interface.iface",
						tfjsonpath.New("bonjour"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"enable":    knownvalue.Bool(true),
							"group_id":  knownvalue.Int64Exact(10),
							"ttl_check": knownvalue.Bool(true),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_tunnel_interface.iface",
						tfjsonpath.New("ip"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name": knownvalue.StringExact("127.0.0.1"),
							}),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_tunnel_interface.iface",
						tfjsonpath.New("ipv6"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"enabled":      knownvalue.Bool(true),
							"interface_id": knownvalue.StringExact("100"),
							"address": knownvalue.ListExact([]knownvalue.Check{
								knownvalue.ObjectExact(map[string]knownvalue.Check{
									"name":                knownvalue.StringExact("::1"),
									"enable_on_interface": knownvalue.Bool(true),
									"anycast":             knownvalue.ObjectExact(nil),
									"prefix":              knownvalue.ObjectExact(nil),
								}),
							}),
						}),
					),
				},
			},
			{
				Config: tunnelInterfaceCleanupTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
			},
		},
	})
}

const tunnelInterfaceResource1 = `
variable "prefix" { type = string }
variable "interface_name" { type = string }

locals {
  template_name = format("%s-tmpl", var.prefix)
}

resource "panos_template" "template" {
  location = { panorama = {} }
  name = local.template_name
}

resource "panos_interface_management_profile" "profile" {
  location = { template = { name = panos_template.template.name }}

  name = format("%s-profile", var.prefix)
}

resource "panos_tunnel_interface" "iface" {
  location = { template = { name = panos_template.template.name } }

  name = var.interface_name
  comment = "tunnel interface comment"

  df_ignore = true
  interface_management_profile = panos_interface_management_profile.profile.name
  #link_tag = "tag-1"
  mtu = "9126"
  #netflow_profile = format("%s-profile", var.prefix)
  bonjour = {
    enable = true
    group_id = 10
    ttl_check = true
  }
  ip = [{
    name = "127.0.0.1"
  }]
  ipv6 = {
    enabled = true
    interface_id = "100"
    address = [{
      name = "::1"
      enable_on_interface = true
      anycast = {}
      prefix = {}
    }]

  }
}
`

const tunnelInterfaceCleanupTmpl = `
variable "prefix" { type = string }

locals {
  template_name = format("%s-tmpl", var.prefix)
}

resource "panos_template" "template" {
  location = { panorama = {} }
  name = local.template_name
}
`
