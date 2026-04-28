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

func TestAccOspfInterfaceTimerRoutingProfile_Basic(t *testing.T) {
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
				Config: ospfInterfaceTimerRoutingProfile_Basic_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_ospf_interface_timer_routing_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_ospf_interface_timer_routing_profile.example",
						tfjsonpath.New("dead_counts"),
						knownvalue.Int64Exact(10),
					),
					statecheck.ExpectKnownValue(
						"panos_ospf_interface_timer_routing_profile.example",
						tfjsonpath.New("gr_delay"),
						knownvalue.Int64Exact(5),
					),
					statecheck.ExpectKnownValue(
						"panos_ospf_interface_timer_routing_profile.example",
						tfjsonpath.New("hello_interval"),
						knownvalue.Int64Exact(30),
					),
					statecheck.ExpectKnownValue(
						"panos_ospf_interface_timer_routing_profile.example",
						tfjsonpath.New("retransmit_interval"),
						knownvalue.Int64Exact(10),
					),
					statecheck.ExpectKnownValue(
						"panos_ospf_interface_timer_routing_profile.example",
						tfjsonpath.New("transit_delay"),
						knownvalue.Int64Exact(5),
					),
				},
			},
		},
	})
}

const ospfInterfaceTimerRoutingProfile_Basic_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_ospf_interface_timer_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix
  dead_counts = 10
  gr_delay = 5
  hello_interval = 30
  retransmit_interval = 10
  transit_delay = 5
}
`
