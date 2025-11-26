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

func TestAccBgpTimerRoutingProfile_Basic(t *testing.T) {
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
				Config: bgpTimerRoutingProfile_Basic_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_bgp_timer_routing_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_bgp_timer_routing_profile.example",
						tfjsonpath.New("hold_time"),
						knownvalue.StringExact("120"),
					),
					statecheck.ExpectKnownValue(
						"panos_bgp_timer_routing_profile.example",
						tfjsonpath.New("keep_alive_interval"),
						knownvalue.StringExact("40"),
					),
					statecheck.ExpectKnownValue(
						"panos_bgp_timer_routing_profile.example",
						tfjsonpath.New("min_route_advertisement_interval"),
						knownvalue.Int64Exact(60),
					),
					statecheck.ExpectKnownValue(
						"panos_bgp_timer_routing_profile.example",
						tfjsonpath.New("open_delay_time"),
						knownvalue.Int64Exact(10),
					),
					statecheck.ExpectKnownValue(
						"panos_bgp_timer_routing_profile.example",
						tfjsonpath.New("reconnect_retry_interval"),
						knownvalue.Int64Exact(30),
					),
				},
			},
		},
	})
}

const bgpTimerRoutingProfile_Basic_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_bgp_timer_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix
  hold_time = "120"
  keep_alive_interval = "40"
  min_route_advertisement_interval = 60
  open_delay_time = 10
  reconnect_retry_interval = 30
}
`
