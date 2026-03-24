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

func TestAccBgpDampeningRoutingProfile_Basic(t *testing.T) {
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
				Config: bgpDampeningRoutingProfile_Basic_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_bgp_dampening_routing_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_bgp_dampening_routing_profile.example",
						tfjsonpath.New("description"),
						knownvalue.StringExact("BGP dampening profile for testing"),
					),
					statecheck.ExpectKnownValue(
						"panos_bgp_dampening_routing_profile.example",
						tfjsonpath.New("half_life"),
						knownvalue.Int64Exact(20),
					),
					statecheck.ExpectKnownValue(
						"panos_bgp_dampening_routing_profile.example",
						tfjsonpath.New("max_suppress_limit"),
						knownvalue.Int64Exact(90),
					),
					statecheck.ExpectKnownValue(
						"panos_bgp_dampening_routing_profile.example",
						tfjsonpath.New("reuse_limit"),
						knownvalue.Int64Exact(1000),
					),
					statecheck.ExpectKnownValue(
						"panos_bgp_dampening_routing_profile.example",
						tfjsonpath.New("suppress_limit"),
						knownvalue.Int64Exact(3000),
					),
				},
			},
		},
	})
}

const bgpDampeningRoutingProfile_Basic_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_bgp_dampening_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix
  description = "BGP dampening profile for testing"
  half_life = 20
  max_suppress_limit = 90
  reuse_limit = 1000
  suppress_limit = 3000
}
`
