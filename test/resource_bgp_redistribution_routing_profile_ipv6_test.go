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

func TestAccBgpRedistributionRoutingProfile_Ipv6_Unicast_Connected(t *testing.T) {
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
				Config: bgpRedistributionRoutingProfile_Ipv6_Unicast_Connected_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_bgp_redistribution_routing_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_bgp_redistribution_routing_profile.example",
						tfjsonpath.New("ipv6"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"unicast": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"connected": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"enable":    knownvalue.Bool(true),
									"metric":    knownvalue.Int64Exact(100),
									"route_map": knownvalue.Null(),
								}),
								"ospfv3": knownvalue.Null(),
								"static": knownvalue.Null(),
							}),
						}),
					),
				},
			},
		},
	})
}

const bgpRedistributionRoutingProfile_Ipv6_Unicast_Connected_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_bgp_redistribution_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix

  ipv6 = {
    unicast = {
      connected = {
        enable = true
        metric = 100
      }
    }
  }
}
`

func TestAccBgpRedistributionRoutingProfile_Ipv6_Unicast_Ospfv3(t *testing.T) {
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
				Config: bgpRedistributionRoutingProfile_Ipv6_Unicast_Ospfv3_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_bgp_redistribution_routing_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_bgp_redistribution_routing_profile.example",
						tfjsonpath.New("ipv6"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"unicast": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"connected": knownvalue.Null(),
								"ospfv3": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"enable":    knownvalue.Bool(true),
									"metric":    knownvalue.Int64Exact(200),
									"route_map": knownvalue.Null(),
								}),
								"static": knownvalue.Null(),
							}),
						}),
					),
				},
			},
		},
	})
}

const bgpRedistributionRoutingProfile_Ipv6_Unicast_Ospfv3_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_bgp_redistribution_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix

  ipv6 = {
    unicast = {
      ospfv3 = {
        enable = true
        metric = 200
      }
    }
  }
}
`

func TestAccBgpRedistributionRoutingProfile_Ipv6_Unicast_Static(t *testing.T) {
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
				Config: bgpRedistributionRoutingProfile_Ipv6_Unicast_Static_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_bgp_redistribution_routing_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_bgp_redistribution_routing_profile.example",
						tfjsonpath.New("ipv6"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"unicast": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"connected": knownvalue.Null(),
								"ospfv3":    knownvalue.Null(),
								"static": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"enable":    knownvalue.Bool(true),
									"metric":    knownvalue.Int64Exact(250),
									"route_map": knownvalue.Null(),
								}),
							}),
						}),
					),
				},
			},
		},
	})
}

const bgpRedistributionRoutingProfile_Ipv6_Unicast_Static_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_bgp_redistribution_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix

  ipv6 = {
    unicast = {
      static = {
        enable = true
        metric = 250
      }
    }
  }
}
`

func TestAccBgpRedistributionRoutingProfile_Ipv6_Unicast_Multiple(t *testing.T) {
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
				Config: bgpRedistributionRoutingProfile_Ipv6_Unicast_Multiple_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_bgp_redistribution_routing_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_bgp_redistribution_routing_profile.example",
						tfjsonpath.New("ipv6"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"unicast": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"connected": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"enable":    knownvalue.Bool(true),
									"metric":    knownvalue.Int64Exact(100),
									"route_map": knownvalue.Null(),
								}),
								"ospfv3": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"enable":    knownvalue.Bool(true),
									"metric":    knownvalue.Int64Exact(200),
									"route_map": knownvalue.Null(),
								}),
								"static": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"enable":    knownvalue.Bool(true),
									"metric":    knownvalue.Int64Exact(250),
									"route_map": knownvalue.Null(),
								}),
							}),
						}),
					),
				},
			},
		},
	})
}

const bgpRedistributionRoutingProfile_Ipv6_Unicast_Multiple_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_bgp_redistribution_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix

  ipv6 = {
    unicast = {
      connected = {
        enable = true
        metric = 100
      }
      ospfv3 = {
        enable = true
        metric = 200
      }
      static = {
        enable = true
        metric = 250
      }
    }
  }
}
`
