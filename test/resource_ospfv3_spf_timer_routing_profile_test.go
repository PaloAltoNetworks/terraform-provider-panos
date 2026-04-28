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

func TestAccOspfv3SpfTimerRoutingProfile_Basic(t *testing.T) {
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
				Config: ospfv3SpfTimerRoutingProfile_Basic_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_ospfv3_spf_timer_routing_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_ospfv3_spf_timer_routing_profile.example",
						tfjsonpath.New("initial_hold_time"),
						knownvalue.Int64Exact(100),
					),
					statecheck.ExpectKnownValue(
						"panos_ospfv3_spf_timer_routing_profile.example",
						tfjsonpath.New("lsa_interval"),
						knownvalue.Int64Exact(8),
					),
					statecheck.ExpectKnownValue(
						"panos_ospfv3_spf_timer_routing_profile.example",
						tfjsonpath.New("max_hold_time"),
						knownvalue.Int64Exact(300),
					),
					statecheck.ExpectKnownValue(
						"panos_ospfv3_spf_timer_routing_profile.example",
						tfjsonpath.New("spf_calculation_delay"),
						knownvalue.Int64Exact(200),
					),
				},
			},
		},
	})
}

const ospfv3SpfTimerRoutingProfile_Basic_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_ospfv3_spf_timer_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix
  initial_hold_time = 100
  lsa_interval = 8
  max_hold_time = 300
  spf_calculation_delay = 200
}
`
