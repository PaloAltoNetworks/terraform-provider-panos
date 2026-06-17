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

func TestAccBfdNetworkProfile_Basic(t *testing.T) {
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
				Config: bfdNetworkProfile_Basic_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_bfd_network_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_bfd_network_profile.example",
						tfjsonpath.New("detection_multiplier"),
						knownvalue.Int64Exact(5),
					),
					statecheck.ExpectKnownValue(
						"panos_bfd_network_profile.example",
						tfjsonpath.New("hold_time"),
						knownvalue.Int64Exact(2000),
					),
					statecheck.ExpectKnownValue(
						"panos_bfd_network_profile.example",
						tfjsonpath.New("min_rx_interval"),
						knownvalue.Int64Exact(300),
					),
					statecheck.ExpectKnownValue(
						"panos_bfd_network_profile.example",
						tfjsonpath.New("min_tx_interval"),
						knownvalue.Int64Exact(400),
					),
					statecheck.ExpectKnownValue(
						"panos_bfd_network_profile.example",
						tfjsonpath.New("mode"),
						knownvalue.StringExact("active"),
					),
					statecheck.ExpectKnownValue(
						"panos_bfd_network_profile.example",
						tfjsonpath.New("multihop"),
						knownvalue.Null(),
					),
				},
			},
		},
	})
}

const bfdNetworkProfile_Basic_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_bfd_network_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix
  detection_multiplier = 5
  hold_time = 2000
  min_rx_interval = 300
  min_tx_interval = 400
  mode = "active"
}
`

func TestAccBfdNetworkProfile_Mode_Passive(t *testing.T) {
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
				Config: bfdNetworkProfile_Mode_Passive_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_bfd_network_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_bfd_network_profile.example",
						tfjsonpath.New("mode"),
						knownvalue.StringExact("passive"),
					),
				},
			},
		},
	})
}

const bfdNetworkProfile_Mode_Passive_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_bfd_network_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix
  mode = "passive"
}
`

func TestAccBfdNetworkProfile_Multihop_MinReceivedTtl(t *testing.T) {
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
				Config: bfdNetworkProfile_Multihop_MinReceivedTtl_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_bfd_network_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_bfd_network_profile.example",
						tfjsonpath.New("multihop"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"min_received_ttl": knownvalue.Int64Exact(64),
						}),
					),
				},
			},
		},
	})
}

const bfdNetworkProfile_Multihop_MinReceivedTtl_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_bfd_network_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix
  multihop = {
    min_received_ttl = 64
  }
}
`

func TestAccBfdNetworkProfile_Multihop_WithAllParams(t *testing.T) {
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
				Config: bfdNetworkProfile_Multihop_WithAllParams_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_bfd_network_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_bfd_network_profile.example",
						tfjsonpath.New("detection_multiplier"),
						knownvalue.Int64Exact(10),
					),
					statecheck.ExpectKnownValue(
						"panos_bfd_network_profile.example",
						tfjsonpath.New("hold_time"),
						knownvalue.Int64Exact(5000),
					),
					statecheck.ExpectKnownValue(
						"panos_bfd_network_profile.example",
						tfjsonpath.New("min_rx_interval"),
						knownvalue.Int64Exact(500),
					),
					statecheck.ExpectKnownValue(
						"panos_bfd_network_profile.example",
						tfjsonpath.New("min_tx_interval"),
						knownvalue.Int64Exact(500),
					),
					statecheck.ExpectKnownValue(
						"panos_bfd_network_profile.example",
						tfjsonpath.New("mode"),
						knownvalue.StringExact("passive"),
					),
					statecheck.ExpectKnownValue(
						"panos_bfd_network_profile.example",
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

const bfdNetworkProfile_Multihop_WithAllParams_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_bfd_network_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix
  detection_multiplier = 10
  hold_time = 5000
  min_rx_interval = 500
  min_tx_interval = 500
  mode = "passive"
  multihop = {
    min_received_ttl = 128
  }
}
`

func TestAccBfdNetworkProfile_Defaults(t *testing.T) {
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
				Config: bfdNetworkProfile_Defaults_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_bfd_network_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_bfd_network_profile.example",
						tfjsonpath.New("detection_multiplier"),
						knownvalue.Int64Exact(3),
					),
					statecheck.ExpectKnownValue(
						"panos_bfd_network_profile.example",
						tfjsonpath.New("hold_time"),
						knownvalue.Int64Exact(0),
					),
					statecheck.ExpectKnownValue(
						"panos_bfd_network_profile.example",
						tfjsonpath.New("min_rx_interval"),
						knownvalue.Int64Exact(1000),
					),
					statecheck.ExpectKnownValue(
						"panos_bfd_network_profile.example",
						tfjsonpath.New("min_tx_interval"),
						knownvalue.Int64Exact(1000),
					),
					statecheck.ExpectKnownValue(
						"panos_bfd_network_profile.example",
						tfjsonpath.New("mode"),
						knownvalue.StringExact("active"),
					),
					statecheck.ExpectKnownValue(
						"panos_bfd_network_profile.example",
						tfjsonpath.New("multihop"),
						knownvalue.Null(),
					),
				},
			},
		},
	})
}

const bfdNetworkProfile_Defaults_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_bfd_network_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix
}
`

func TestAccBfdNetworkProfile_TemplateStack(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: bfdNetworkProfile_TemplateStack_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_bfd_network_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-profile", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_bfd_network_profile.example",
						tfjsonpath.New("detection_multiplier"),
						knownvalue.Int64Exact(4),
					),
					statecheck.ExpectKnownValue(
						"panos_bfd_network_profile.example",
						tfjsonpath.New("mode"),
						knownvalue.StringExact("active"),
					),
					statecheck.ExpectKnownValue(
						"panos_bfd_network_profile.example",
						tfjsonpath.New("location").AtMapKey("template_stack").AtMapKey("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-stack", prefix)),
					),
				},
			},
		},
	})
}

const bfdNetworkProfile_TemplateStack_Tmpl = `
variable "prefix" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_template_stack" "example" {
  location = { panorama = {} }
  name = "${var.prefix}-stack"
  templates = [panos_template.example.name]
}

resource "panos_bfd_network_profile" "example" {
  location = { template_stack = { name = panos_template_stack.example.name } }

  name = "${var.prefix}-profile"
  detection_multiplier = 4
  mode = "active"
}
`
