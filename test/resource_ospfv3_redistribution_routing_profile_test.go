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

func TestAccOspfv3RedistributionRoutingProfile_Basic(t *testing.T) {
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
				Config: ospfv3RedistributionRoutingProfile_Basic_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_ospfv3_redistribution_routing_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_ospfv3_redistribution_routing_profile.example",
						tfjsonpath.New("bgp"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"enable":      knownvalue.Bool(true),
							"metric":      knownvalue.Int64Exact(100),
							"metric_type": knownvalue.StringExact("type-1"),
							"route_map":   knownvalue.Null(),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_ospfv3_redistribution_routing_profile.example",
						tfjsonpath.New("connected"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"panos_ospfv3_redistribution_routing_profile.example",
						tfjsonpath.New("default_route"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"panos_ospfv3_redistribution_routing_profile.example",
						tfjsonpath.New("static"),
						knownvalue.Null(),
					),
				},
			},
		},
	})
}

const ospfv3RedistributionRoutingProfile_Basic_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_ospfv3_redistribution_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix
  bgp = {
    enable = true
    metric = 100
    metric_type = "type-1"
  }
}
`

func TestAccOspfv3RedistributionRoutingProfile_DefaultRoute(t *testing.T) {
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
				Config: ospfv3RedistributionRoutingProfile_DefaultRoute_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_ospfv3_redistribution_routing_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_ospfv3_redistribution_routing_profile.example",
						tfjsonpath.New("default_route"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"always":      knownvalue.Bool(true),
							"enable":      knownvalue.Bool(true),
							"metric":      knownvalue.Int64Exact(10),
							"metric_type": knownvalue.StringExact("type-1"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_ospfv3_redistribution_routing_profile.example",
						tfjsonpath.New("bgp"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"panos_ospfv3_redistribution_routing_profile.example",
						tfjsonpath.New("connected"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"panos_ospfv3_redistribution_routing_profile.example",
						tfjsonpath.New("static"),
						knownvalue.Null(),
					),
				},
			},
		},
	})
}

const ospfv3RedistributionRoutingProfile_DefaultRoute_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_ospfv3_redistribution_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix
  default_route = {
    always = true
    enable = true
    metric = 10
    metric_type = "type-1"
  }
}
`

func TestAccOspfv3RedistributionRoutingProfile_StaticAndConnected(t *testing.T) {
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
				Config: ospfv3RedistributionRoutingProfile_StaticAndConnected_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_ospfv3_redistribution_routing_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_ospfv3_redistribution_routing_profile.example",
						tfjsonpath.New("static"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"enable":      knownvalue.Bool(true),
							"metric":      knownvalue.Int64Exact(200),
							"metric_type": knownvalue.StringExact("type-2"),
							"route_map":   knownvalue.Null(),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_ospfv3_redistribution_routing_profile.example",
						tfjsonpath.New("connected"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"enable":      knownvalue.Bool(true),
							"metric":      knownvalue.Int64Exact(50),
							"metric_type": knownvalue.StringExact("type-1"),
							"route_map":   knownvalue.Null(),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_ospfv3_redistribution_routing_profile.example",
						tfjsonpath.New("bgp"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"panos_ospfv3_redistribution_routing_profile.example",
						tfjsonpath.New("default_route"),
						knownvalue.Null(),
					),
				},
			},
		},
	})
}

const ospfv3RedistributionRoutingProfile_StaticAndConnected_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_ospfv3_redistribution_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix
  static = {
    enable = true
    metric = 200
    metric_type = "type-2"
  }
  connected = {
    enable = true
    metric = 50
    metric_type = "type-1"
  }
}
`
