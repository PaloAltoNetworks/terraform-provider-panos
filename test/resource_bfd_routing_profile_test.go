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

func TestAccBfdRoutingProfile_Basic(t *testing.T) {
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
				Config: bfdRoutingProfile_Basic_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_bfd_routing_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_bfd_routing_profile.example",
						tfjsonpath.New("detection_multiplier"),
						knownvalue.Int64Exact(5),
					),
					statecheck.ExpectKnownValue(
						"panos_bfd_routing_profile.example",
						tfjsonpath.New("hold_time"),
						knownvalue.Int64Exact(1000),
					),
					statecheck.ExpectKnownValue(
						"panos_bfd_routing_profile.example",
						tfjsonpath.New("min_rx_interval"),
						knownvalue.Int64Exact(500),
					),
					statecheck.ExpectKnownValue(
						"panos_bfd_routing_profile.example",
						tfjsonpath.New("min_tx_interval"),
						knownvalue.Int64Exact(500),
					),
					statecheck.ExpectKnownValue(
						"panos_bfd_routing_profile.example",
						tfjsonpath.New("mode"),
						knownvalue.StringExact("active"),
					),
				},
			},
		},
	})
}

const bfdRoutingProfile_Basic_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_bfd_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix
  detection_multiplier = 5
  hold_time = 1000
  min_rx_interval = 500
  min_tx_interval = 500
  mode = "active"
}
`

func TestAccBfdRoutingProfile_Mode_Passive(t *testing.T) {
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
				Config: bfdRoutingProfile_Mode_Passive_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_bfd_routing_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_bfd_routing_profile.example",
						tfjsonpath.New("mode"),
						knownvalue.StringExact("passive"),
					),
				},
			},
		},
	})
}

const bfdRoutingProfile_Mode_Passive_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_bfd_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix
  mode = "passive"
}
`

func TestAccBfdRoutingProfile_Multihop(t *testing.T) {
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
				Config: bfdRoutingProfile_Multihop_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_bfd_routing_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_bfd_routing_profile.example",
						tfjsonpath.New("multihop"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"min_received_ttl": knownvalue.Int64Exact(128),
						}),
					),
				},
			},
		},
	})
}

const bfdRoutingProfile_Multihop_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_bfd_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix
  multihop = {
    min_received_ttl = 128
  }
}
`
