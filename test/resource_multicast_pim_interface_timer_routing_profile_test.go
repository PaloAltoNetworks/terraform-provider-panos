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

func TestAccMulticastPimInterfaceTimerRoutingProfile_Basic(t *testing.T) {
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
				Config: multicastPimInterfaceTimerRoutingProfile_Basic_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_multicast_pim_interface_timer_routing_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_multicast_pim_interface_timer_routing_profile.example",
						tfjsonpath.New("assert_interval"),
						knownvalue.Int64Exact(200),
					),
					statecheck.ExpectKnownValue(
						"panos_multicast_pim_interface_timer_routing_profile.example",
						tfjsonpath.New("hello_interval"),
						knownvalue.Int64Exact(60),
					),
					statecheck.ExpectKnownValue(
						"panos_multicast_pim_interface_timer_routing_profile.example",
						tfjsonpath.New("join_prune_interval"),
						knownvalue.Int64Exact(120),
					),
				},
			},
		},
	})
}

const multicastPimInterfaceTimerRoutingProfile_Basic_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_multicast_pim_interface_timer_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix
  assert_interval = 200
  hello_interval = 60
  join_prune_interval = 120
}
`
