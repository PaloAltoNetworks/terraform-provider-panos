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

func TestAccBgpAddressFamilyRoutingProfile_Ipv6_Unicast_Basic(t *testing.T) {
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
				Config: bgpAddressFamilyRoutingProfile_Ipv6_Unicast_Basic_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_bgp_address_family_routing_profile.example",
						tfjsonpath.New("ipv6"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"unicast": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"enable":                           knownvalue.Bool(true),
								"as_override":                      knownvalue.Bool(false),
								"default_originate":                knownvalue.Bool(true),
								"default_originate_map":            knownvalue.Null(),
								"route_reflector_client":           knownvalue.Bool(false),
								"soft_reconfig_with_stored_info":   knownvalue.Bool(true),
								"add_path": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"tx_all_paths":       knownvalue.Bool(true),
									"tx_bestpath_per_as": knownvalue.Bool(false),
								}),
								"orf": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"orf_prefix_list": knownvalue.StringExact("receive"),
								}),
								"allowas_in":          knownvalue.Null(),
								"maximum_prefix":      knownvalue.Null(),
								"next_hop":            knownvalue.Null(),
								"remove_private_as":   knownvalue.Null(),
								"send_community":      knownvalue.Null(),
							}),
						}),
					),
				},
			},
		},
	})
}

func TestAccBgpAddressFamilyRoutingProfile_Ipv6_Unicast_AllowasIn_Occurrence(t *testing.T) {
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
				Config: bgpAddressFamilyRoutingProfile_Ipv6_Unicast_AllowasIn_Occurrence_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_bgp_address_family_routing_profile.example",
						tfjsonpath.New("ipv6").AtMapKey("unicast").AtMapKey("allowas_in"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"occurrence": knownvalue.Int64Exact(7),
							"origin":     knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

func TestAccBgpAddressFamilyRoutingProfile_Ipv6_Unicast_AllowasIn_Origin(t *testing.T) {
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
				Config: bgpAddressFamilyRoutingProfile_Ipv6_Unicast_AllowasIn_Origin_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_bgp_address_family_routing_profile.example",
						tfjsonpath.New("ipv6").AtMapKey("unicast").AtMapKey("allowas_in"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"occurrence": knownvalue.Null(),
							"origin":     knownvalue.ObjectExact(map[string]knownvalue.Check{}),
						}),
					),
				},
			},
		},
	})
}

func TestAccBgpAddressFamilyRoutingProfile_Ipv6_Unicast_SendCommunity_Extended(t *testing.T) {
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
				Config: bgpAddressFamilyRoutingProfile_Ipv6_Unicast_SendCommunity_Extended_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_bgp_address_family_routing_profile.example",
						tfjsonpath.New("ipv6").AtMapKey("unicast").AtMapKey("send_community"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"all":      knownvalue.Null(),
							"both":     knownvalue.Null(),
							"extended": knownvalue.ObjectExact(map[string]knownvalue.Check{}),
							"large":    knownvalue.Null(),
							"standard": knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

func TestAccBgpAddressFamilyRoutingProfile_Ipv6_Unicast_SendCommunity_Large(t *testing.T) {
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
				Config: bgpAddressFamilyRoutingProfile_Ipv6_Unicast_SendCommunity_Large_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_bgp_address_family_routing_profile.example",
						tfjsonpath.New("ipv6").AtMapKey("unicast").AtMapKey("send_community"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"all":      knownvalue.Null(),
							"both":     knownvalue.Null(),
							"extended": knownvalue.Null(),
							"large":    knownvalue.ObjectExact(map[string]knownvalue.Check{}),
							"standard": knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

func TestAccBgpAddressFamilyRoutingProfile_Ipv6_Unicast_SendCommunity_Standard(t *testing.T) {
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
				Config: bgpAddressFamilyRoutingProfile_Ipv6_Unicast_SendCommunity_Standard_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_bgp_address_family_routing_profile.example",
						tfjsonpath.New("ipv6").AtMapKey("unicast").AtMapKey("send_community"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"all":      knownvalue.Null(),
							"both":     knownvalue.Null(),
							"extended": knownvalue.Null(),
							"large":    knownvalue.Null(),
							"standard": knownvalue.ObjectExact(map[string]knownvalue.Check{}),
						}),
					),
				},
			},
		},
	})
}

const bgpAddressFamilyRoutingProfile_Ipv6_Unicast_Basic_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_bgp_address_family_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix

  ipv6 = {
    unicast = {
      enable = true
      as_override = false
      default_originate = true
      route_reflector_client = false
      soft_reconfig_with_stored_info = true
      add_path = {
        tx_all_paths = true
        tx_bestpath_per_as = false
      }
      orf = {
        orf_prefix_list = "receive"
      }
    }
  }
}
`

const bgpAddressFamilyRoutingProfile_Ipv6_Unicast_AllowasIn_Occurrence_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_bgp_address_family_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix

  ipv6 = {
    unicast = {
      allowas_in = {
        occurrence = 7
      }
    }
  }
}
`

const bgpAddressFamilyRoutingProfile_Ipv6_Unicast_AllowasIn_Origin_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_bgp_address_family_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix

  ipv6 = {
    unicast = {
      allowas_in = {
        origin = {}
      }
    }
  }
}
`

const bgpAddressFamilyRoutingProfile_Ipv6_Unicast_SendCommunity_Extended_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_bgp_address_family_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix

  ipv6 = {
    unicast = {
      send_community = {
        extended = {}
      }
    }
  }
}
`

const bgpAddressFamilyRoutingProfile_Ipv6_Unicast_SendCommunity_Large_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_bgp_address_family_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix

  ipv6 = {
    unicast = {
      send_community = {
        large = {}
      }
    }
  }
}
`

const bgpAddressFamilyRoutingProfile_Ipv6_Unicast_SendCommunity_Standard_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_bgp_address_family_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix

  ipv6 = {
    unicast = {
      send_community = {
        standard = {}
      }
    }
  }
}
`
