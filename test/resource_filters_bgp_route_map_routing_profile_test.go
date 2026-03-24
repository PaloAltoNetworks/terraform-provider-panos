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

func TestAccFiltersBgpRouteMapRoutingProfile_Basic(t *testing.T) {
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
				Config: filtersBgpRouteMapRoutingProfile_Basic_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_filters_bgp_route_map_routing_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_filters_bgp_route_map_routing_profile.example",
						tfjsonpath.New("description"),
						knownvalue.StringExact("BGP Route Map for testing"),
					),
					statecheck.ExpectKnownValue(
						"panos_filters_bgp_route_map_routing_profile.example",
						tfjsonpath.New("route_map"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":        knownvalue.StringExact("10"),
								"action":      knownvalue.StringExact("deny"),
								"description": knownvalue.StringExact("First route map entry"),
								"match":       knownvalue.Null(),
								"set":         knownvalue.Null(),
							}),
						}),
					),
				},
			},
		},
	})
}

const filtersBgpRouteMapRoutingProfile_Basic_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_filters_bgp_route_map_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix
  description = "BGP Route Map for testing"
  route_map = [
    {
      name = "10"
      action = "deny"
      description = "First route map entry"
    }
  ]
}
`

func TestAccFiltersBgpRouteMapRoutingProfile_RouteMap_ActionPermit(t *testing.T) {
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
				Config: filtersBgpRouteMapRoutingProfile_RouteMap_ActionPermit_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_filters_bgp_route_map_routing_profile.example",
						tfjsonpath.New("route_map").AtSliceIndex(0).AtMapKey("action"),
						knownvalue.StringExact("permit"),
					),
				},
			},
		},
	})
}

const filtersBgpRouteMapRoutingProfile_RouteMap_ActionPermit_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_filters_bgp_route_map_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix
  route_map = [
    {
      name = "10"
      action = "permit"
    }
  ]
}
`

func TestAccFiltersBgpRouteMapRoutingProfile_Match_BasicStrings(t *testing.T) {
	t.Parallel()
	t.Skip("requires AS path access list and community list resources")

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
				Config: filtersBgpRouteMapRoutingProfile_Match_BasicStrings_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_filters_bgp_route_map_routing_profile.example",
						tfjsonpath.New("route_map").AtSliceIndex(0).AtMapKey("match"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"as_path_access_list": knownvalue.StringExact("as-path-list1"),
							"regular_community":   knownvalue.StringExact("regular-comm1"),
							"large_community":     knownvalue.StringExact("large-comm1"),
							"extended_community":  knownvalue.StringExact("extended-comm1"),
							"interface":           knownvalue.StringExact("ethernet1/1"),
							"peer":                knownvalue.StringExact("10.0.0.1"),
							"origin":              knownvalue.Null(),
							"metric":              knownvalue.Null(),
							"tag":                 knownvalue.Null(),
							"local_preference":    knownvalue.Null(),
							"ipv4":                knownvalue.Null(),
							"ipv6":                knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

const filtersBgpRouteMapRoutingProfile_Match_BasicStrings_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_filters_bgp_route_map_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix
  route_map = [
    {
      name = "10"
      match = {
        as_path_access_list = "as-path-list1"
        regular_community = "regular-comm1"
        large_community = "large-comm1"
        extended_community = "extended-comm1"
        interface = "ethernet1/1"
        peer = "10.0.0.1"
      }
    }
  ]
}
`

func TestAccFiltersBgpRouteMapRoutingProfile_Match_Enums(t *testing.T) {
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
				Config: filtersBgpRouteMapRoutingProfile_Match_Enums_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_filters_bgp_route_map_routing_profile.example",
						tfjsonpath.New("route_map").AtSliceIndex(0).AtMapKey("match"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"as_path_access_list": knownvalue.Null(),
							"regular_community":   knownvalue.Null(),
							"large_community":     knownvalue.Null(),
							"extended_community":  knownvalue.Null(),
							"interface":           knownvalue.Null(),
							"peer":                knownvalue.Null(),
							"origin":              knownvalue.StringExact("igp"),
							"metric":              knownvalue.Int64Exact(100),
							"tag":                 knownvalue.Int64Exact(200),
							"local_preference":    knownvalue.Int64Exact(150),
							"ipv4":                knownvalue.Null(),
							"ipv6":                knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

const filtersBgpRouteMapRoutingProfile_Match_Enums_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_filters_bgp_route_map_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix
  route_map = [
    {
      name = "10"
      match = {
        origin = "igp"
        metric = 100
        tag = 200
        local_preference = 150
      }
    }
  ]
}
`

func TestAccFiltersBgpRouteMapRoutingProfile_Set_BasicValues(t *testing.T) {
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
				Config: filtersBgpRouteMapRoutingProfile_Set_BasicValues_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_filters_bgp_route_map_routing_profile.example",
						tfjsonpath.New("route_map").AtSliceIndex(0).AtMapKey("set"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"tag":                              knownvalue.Int64Exact(300),
							"local_preference":                 knownvalue.Int64Exact(250),
							"weight":                           knownvalue.Int64Exact(500),
							"origin":                           knownvalue.StringExact("incomplete"),
							"originator_id":                    knownvalue.StringExact("1.1.1.1"),
							"aggregator":                       knownvalue.Null(),
							"metric":                           knownvalue.Null(),
							"ipv4":                             knownvalue.Null(),
							"ipv6":                             knownvalue.Null(),
							"atomic_aggregate":                 knownvalue.Null(),
							"ipv6_nexthop_prefer_global":       knownvalue.Null(),
							"overwrite_regular_community":      knownvalue.Null(),
							"overwrite_large_community":        knownvalue.Null(),
							"remove_regular_community":         knownvalue.Null(),
							"remove_large_community":           knownvalue.Null(),
							"aspath_prepend":                   knownvalue.Null(),
							"regular_community":                knownvalue.Null(),
							"large_community":                  knownvalue.Null(),
							"aspath_exclude":                   knownvalue.Null(),
							"extended_community":               knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

const filtersBgpRouteMapRoutingProfile_Set_BasicValues_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_filters_bgp_route_map_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix
  route_map = [
    {
      name = "10"
      set = {
        tag = 300
        local_preference = 250
        weight = 500
        origin = "incomplete"
        originator_id = "1.1.1.1"
      }
    }
  ]
}
`

func TestAccFiltersBgpRouteMapRoutingProfile_Set_BooleanFlags(t *testing.T) {
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
				Config: filtersBgpRouteMapRoutingProfile_Set_BooleanFlags_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_filters_bgp_route_map_routing_profile.example",
						tfjsonpath.New("route_map").AtSliceIndex(0).AtMapKey("set").AtMapKey("atomic_aggregate"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_filters_bgp_route_map_routing_profile.example",
						tfjsonpath.New("route_map").AtSliceIndex(0).AtMapKey("set").AtMapKey("ipv6_nexthop_prefer_global"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_filters_bgp_route_map_routing_profile.example",
						tfjsonpath.New("route_map").AtSliceIndex(0).AtMapKey("set").AtMapKey("overwrite_regular_community"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_filters_bgp_route_map_routing_profile.example",
						tfjsonpath.New("route_map").AtSliceIndex(0).AtMapKey("set").AtMapKey("overwrite_large_community"),
						knownvalue.Bool(true),
					),
				},
			},
		},
	})
}

const filtersBgpRouteMapRoutingProfile_Set_BooleanFlags_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_filters_bgp_route_map_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix
  route_map = [
    {
      name = "10"
      set = {
        atomic_aggregate = true
        ipv6_nexthop_prefer_global = true
        overwrite_regular_community = true
        overwrite_large_community = true
      }
    }
  ]
}
`

func TestAccFiltersBgpRouteMapRoutingProfile_Set_Metric(t *testing.T) {
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
				Config: filtersBgpRouteMapRoutingProfile_Set_Metric_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_filters_bgp_route_map_routing_profile.example",
						tfjsonpath.New("route_map").AtSliceIndex(0).AtMapKey("set").AtMapKey("metric"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"value":  knownvalue.Int64Exact(400),
							"action": knownvalue.StringExact("add"),
						}),
					),
				},
			},
		},
	})
}

const filtersBgpRouteMapRoutingProfile_Set_Metric_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_filters_bgp_route_map_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix
  route_map = [
    {
      name = "10"
      set = {
        metric = {
          value = 400
          action = "add"
        }
      }
    }
  ]
}
`

func TestAccFiltersBgpRouteMapRoutingProfile_Set_Aggregator(t *testing.T) {
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
				Config: filtersBgpRouteMapRoutingProfile_Set_Aggregator_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_filters_bgp_route_map_routing_profile.example",
						tfjsonpath.New("route_map").AtSliceIndex(0).AtMapKey("set").AtMapKey("aggregator"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"as":        knownvalue.Int64Exact(65001),
							"router_id": knownvalue.StringExact("2.2.2.2"),
						}),
					),
				},
			},
		},
	})
}

const filtersBgpRouteMapRoutingProfile_Set_Aggregator_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_filters_bgp_route_map_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix
  route_map = [
    {
      name = "10"
      set = {
        aggregator = {
          as = 65001
          router_id = "2.2.2.2"
        }
      }
    }
  ]
}
`

func TestAccFiltersBgpRouteMapRoutingProfile_Set_Ipv4(t *testing.T) {
	t.Parallel()
	t.Skip("source_address validation requires additional routing configuration beyond interface creation")

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
				Config: filtersBgpRouteMapRoutingProfile_Set_Ipv4_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_filters_bgp_route_map_routing_profile.example",
						tfjsonpath.New("route_map").AtSliceIndex(0).AtMapKey("set").AtMapKey("ipv4"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"source_address": knownvalue.StringExact("ethernet1/1"),
							"next_hop":       knownvalue.StringExact("10.2.2.2"),
						}),
					),
				},
			},
		},
	})
}

const filtersBgpRouteMapRoutingProfile_Set_Ipv4_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_ethernet_interface" "example" {
  location = { template = { name = panos_template.example.name } }

  name = "ethernet1/1"
  layer3 = {
    ips = [
      { name = "10.1.1.1/24" }
    ]
  }
}

resource "panos_filters_bgp_route_map_routing_profile" "example" {
  depends_on = [panos_ethernet_interface.example]
  location = var.location

  name = var.prefix
  route_map = [
    {
      name = "10"
      set = {
        ipv4 = {
          source_address = panos_ethernet_interface.example.name
          next_hop = "10.2.2.2"
        }
      }
    }
  ]
}
`

func TestAccFiltersBgpRouteMapRoutingProfile_Set_Ipv6(t *testing.T) {
	t.Parallel()
	t.Skip("source_address validation requires additional routing configuration beyond interface creation")

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
				Config: filtersBgpRouteMapRoutingProfile_Set_Ipv6_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_filters_bgp_route_map_routing_profile.example",
						tfjsonpath.New("route_map").AtSliceIndex(0).AtMapKey("set").AtMapKey("ipv6"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"source_address": knownvalue.StringExact("ethernet1/1"),
							"next_hop":       knownvalue.StringExact("2001:db8:2::1"),
						}),
					),
				},
			},
		},
	})
}

const filtersBgpRouteMapRoutingProfile_Set_Ipv6_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_ethernet_interface" "example" {
  location = { template = { name = panos_template.example.name } }

  name = "ethernet1/1"
  layer3 = {
    ipv6 = {
      addresses = [
        { name = "2001:db8:1::1/64" }
      ]
    }
  }
}

resource "panos_filters_bgp_route_map_routing_profile" "example" {
  depends_on = [panos_ethernet_interface.example]
  location = var.location

  name = var.prefix
  route_map = [
    {
      name = "10"
      set = {
        ipv6 = {
          source_address = panos_ethernet_interface.example.name
          next_hop = "2001:db8:2::1"
        }
      }
    }
  ]
}
`

func TestAccFiltersBgpRouteMapRoutingProfile_Set_AspathPrepend(t *testing.T) {
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
				Config: filtersBgpRouteMapRoutingProfile_Set_AspathPrepend_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_filters_bgp_route_map_routing_profile.example",
						tfjsonpath.New("route_map").AtSliceIndex(0).AtMapKey("set").AtMapKey("aspath_prepend"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.Int64Exact(65001),
							knownvalue.Int64Exact(65002),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_filters_bgp_route_map_routing_profile.example",
						tfjsonpath.New("route_map").AtSliceIndex(0).AtMapKey("set").AtMapKey("aspath_exclude"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.Int64Exact(65003),
						}),
					),
				},
			},
		},
	})
}

const filtersBgpRouteMapRoutingProfile_Set_AspathPrepend_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_filters_bgp_route_map_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix
  route_map = [
    {
      name = "10"
      set = {
        aspath_prepend = [65001, 65002]
        aspath_exclude = [65003]
      }
    }
  ]
}
`

func TestAccFiltersBgpRouteMapRoutingProfile_MultipleRouteMaps(t *testing.T) {
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
				Config: filtersBgpRouteMapRoutingProfile_MultipleRouteMaps_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_filters_bgp_route_map_routing_profile.example",
						tfjsonpath.New("route_map"),
						knownvalue.ListSizeExact(3),
					),
					statecheck.ExpectKnownValue(
						"panos_filters_bgp_route_map_routing_profile.example",
						tfjsonpath.New("route_map").AtSliceIndex(0).AtMapKey("name"),
						knownvalue.StringExact("10"),
					),
					statecheck.ExpectKnownValue(
						"panos_filters_bgp_route_map_routing_profile.example",
						tfjsonpath.New("route_map").AtSliceIndex(1).AtMapKey("name"),
						knownvalue.StringExact("20"),
					),
					statecheck.ExpectKnownValue(
						"panos_filters_bgp_route_map_routing_profile.example",
						tfjsonpath.New("route_map").AtSliceIndex(2).AtMapKey("name"),
						knownvalue.StringExact("30"),
					),
				},
			},
		},
	})
}

const filtersBgpRouteMapRoutingProfile_MultipleRouteMaps_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_filters_bgp_route_map_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix
  route_map = [
    {
      name = "10"
      action = "permit"
      description = "First route map"
    },
    {
      name = "20"
      action = "deny"
      description = "Second route map"
    },
    {
      name = "30"
      action = "permit"
      description = "Third route map"
    }
  ]
}
`

func TestAccFiltersBgpRouteMapRoutingProfile_Set_RegularCommunity(t *testing.T) {
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
				Config: filtersBgpRouteMapRoutingProfile_Set_RegularCommunity_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_filters_bgp_route_map_routing_profile.example",
						tfjsonpath.New("route_map").AtSliceIndex(0).AtMapKey("set").AtMapKey("regular_community"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("65001:100"),
							knownvalue.StringExact("65001:200"),
						}),
					),
				},
			},
		},
	})
}

const filtersBgpRouteMapRoutingProfile_Set_RegularCommunity_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_filters_bgp_route_map_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix
  route_map = [
    {
      name = "10"
      set = {
        regular_community = ["65001:100", "65001:200"]
      }
    }
  ]
}
`

func TestAccFiltersBgpRouteMapRoutingProfile_Set_LargeCommunity(t *testing.T) {
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
				Config: filtersBgpRouteMapRoutingProfile_Set_LargeCommunity_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_filters_bgp_route_map_routing_profile.example",
						tfjsonpath.New("route_map").AtSliceIndex(0).AtMapKey("set").AtMapKey("large_community"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("65001:100:200"),
							knownvalue.StringExact("65001:300:400"),
						}),
					),
				},
			},
		},
	})
}

const filtersBgpRouteMapRoutingProfile_Set_LargeCommunity_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_filters_bgp_route_map_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix
  route_map = [
    {
      name = "10"
      set = {
        large_community = ["65001:100:200", "65001:300:400"]
      }
    }
  ]
}
`

func TestAccFiltersBgpRouteMapRoutingProfile_Set_ExtendedCommunity(t *testing.T) {
	t.Parallel()
	t.Skip("extended_community only available in PAN-OS 11.0.2-11.0.3")

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
				Config: filtersBgpRouteMapRoutingProfile_Set_ExtendedCommunity_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_filters_bgp_route_map_routing_profile.example",
						tfjsonpath.New("route_map").AtSliceIndex(0).AtMapKey("set").AtMapKey("extended_community"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("rt:65001:100"),
							knownvalue.StringExact("soo:65001:200"),
						}),
					),
				},
			},
		},
	})
}

const filtersBgpRouteMapRoutingProfile_Set_ExtendedCommunity_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_filters_bgp_route_map_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix
  route_map = [
    {
      name = "10"
      set = {
        extended_community = ["rt:65001:100", "soo:65001:200"]
      }
    }
  ]
}
`

func TestAccFiltersBgpRouteMapRoutingProfile_Set_RemoveCommunities(t *testing.T) {
	t.Parallel()
	t.Skip("remove community requires community list resources")

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
				Config: filtersBgpRouteMapRoutingProfile_Set_RemoveCommunities_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_filters_bgp_route_map_routing_profile.example",
						tfjsonpath.New("route_map").AtSliceIndex(0).AtMapKey("set").AtMapKey("remove_regular_community"),
						knownvalue.StringExact("remove-regular-comm"),
					),
					statecheck.ExpectKnownValue(
						"panos_filters_bgp_route_map_routing_profile.example",
						tfjsonpath.New("route_map").AtSliceIndex(0).AtMapKey("set").AtMapKey("remove_large_community"),
						knownvalue.StringExact("remove-large-comm"),
					),
				},
			},
		},
	})
}

const filtersBgpRouteMapRoutingProfile_Set_RemoveCommunities_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_filters_bgp_route_map_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix
  route_map = [
    {
      name = "10"
      set = {
        remove_regular_community = "remove-regular-comm"
        remove_large_community = "remove-large-comm"
      }
    }
  ]
}
`

func TestAccFiltersBgpRouteMapRoutingProfile_Set_Metric_ActionSet(t *testing.T) {
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
				Config: filtersBgpRouteMapRoutingProfile_Set_Metric_ActionSet_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_filters_bgp_route_map_routing_profile.example",
						tfjsonpath.New("route_map").AtSliceIndex(0).AtMapKey("set").AtMapKey("metric"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"value":  knownvalue.Int64Exact(500),
							"action": knownvalue.StringExact("set"),
						}),
					),
				},
			},
		},
	})
}

const filtersBgpRouteMapRoutingProfile_Set_Metric_ActionSet_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_filters_bgp_route_map_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix
  route_map = [
    {
      name = "10"
      set = {
        metric = {
          value = 500
          action = "set"
        }
      }
    }
  ]
}
`

func TestAccFiltersBgpRouteMapRoutingProfile_Set_Metric_ActionSubtract(t *testing.T) {
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
				Config: filtersBgpRouteMapRoutingProfile_Set_Metric_ActionSubtract_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_filters_bgp_route_map_routing_profile.example",
						tfjsonpath.New("route_map").AtSliceIndex(0).AtMapKey("set").AtMapKey("metric"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"value":  knownvalue.Int64Exact(100),
							"action": knownvalue.StringExact("subtract"),
						}),
					),
				},
			},
		},
	})
}

const filtersBgpRouteMapRoutingProfile_Set_Metric_ActionSubtract_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_filters_bgp_route_map_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix
  route_map = [
    {
      name = "10"
      set = {
        metric = {
          value = 100
          action = "subtract"
        }
      }
    }
  ]
}
`

func TestAccFiltersBgpRouteMapRoutingProfile_Match_Origin_Variants(t *testing.T) {
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
				Config: filtersBgpRouteMapRoutingProfile_Match_Origin_Variants_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_filters_bgp_route_map_routing_profile.example",
						tfjsonpath.New("route_map").AtSliceIndex(0).AtMapKey("match").AtMapKey("origin"),
						knownvalue.StringExact("egp"),
					),
					statecheck.ExpectKnownValue(
						"panos_filters_bgp_route_map_routing_profile.example",
						tfjsonpath.New("route_map").AtSliceIndex(1).AtMapKey("match").AtMapKey("origin"),
						knownvalue.StringExact("incomplete"),
					),
					statecheck.ExpectKnownValue(
						"panos_filters_bgp_route_map_routing_profile.example",
						tfjsonpath.New("route_map").AtSliceIndex(2).AtMapKey("match").AtMapKey("origin"),
						knownvalue.StringExact("none"),
					),
				},
			},
		},
	})
}

const filtersBgpRouteMapRoutingProfile_Match_Origin_Variants_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_filters_bgp_route_map_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix
  route_map = [
    {
      name = "10"
      match = {
        origin = "egp"
      }
    },
    {
      name = "20"
      match = {
        origin = "incomplete"
      }
    },
    {
      name = "30"
      match = {
        origin = "none"
      }
    }
  ]
}
`

func TestAccFiltersBgpRouteMapRoutingProfile_Set_Origin_Variants(t *testing.T) {
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
				Config: filtersBgpRouteMapRoutingProfile_Set_Origin_Variants_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_filters_bgp_route_map_routing_profile.example",
						tfjsonpath.New("route_map").AtSliceIndex(0).AtMapKey("set").AtMapKey("origin"),
						knownvalue.StringExact("egp"),
					),
					statecheck.ExpectKnownValue(
						"panos_filters_bgp_route_map_routing_profile.example",
						tfjsonpath.New("route_map").AtSliceIndex(1).AtMapKey("set").AtMapKey("origin"),
						knownvalue.StringExact("igp"),
					),
					statecheck.ExpectKnownValue(
						"panos_filters_bgp_route_map_routing_profile.example",
						tfjsonpath.New("route_map").AtSliceIndex(2).AtMapKey("set").AtMapKey("origin"),
						knownvalue.StringExact("none"),
					),
				},
			},
		},
	})
}

const filtersBgpRouteMapRoutingProfile_Set_Origin_Variants_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_filters_bgp_route_map_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix
  route_map = [
    {
      name = "10"
      set = {
        origin = "egp"
      }
    },
    {
      name = "20"
      set = {
        origin = "igp"
      }
    },
    {
      name = "30"
      set = {
        origin = "none"
      }
    }
  ]
}
`

func TestAccFiltersBgpRouteMapRoutingProfile_Match_Ipv4_Address_AccessList(t *testing.T) {
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
				Config: filtersBgpRouteMapRoutingProfile_Match_Ipv4_Address_AccessList_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_filters_bgp_route_map_routing_profile.example",
						tfjsonpath.New("route_map").AtSliceIndex(0).AtMapKey("match").AtMapKey("ipv4").AtMapKey("address").AtMapKey("access_list"),
						knownvalue.StringExact(prefix+"-acl"),
					),
				},
			},
		},
	})
}

const filtersBgpRouteMapRoutingProfile_Match_Ipv4_Address_AccessList_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_filters_access_list_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = "${var.prefix}-acl"
  type = {
    ipv4 = {
      ipv4_entries = [
        {
          name = "10"
          action = "permit"
          source_address = {
            address = "any"
          }
        }
      ]
    }
  }
}

resource "panos_filters_bgp_route_map_routing_profile" "example" {
  depends_on = [panos_filters_access_list_routing_profile.example]
  location = var.location

  name = var.prefix
  route_map = [
    {
      name = "10"
      match = {
        ipv4 = {
          address = {
            access_list = panos_filters_access_list_routing_profile.example.name
          }
        }
      }
    }
  ]
}
`

func TestAccFiltersBgpRouteMapRoutingProfile_Match_Ipv4_Address_PrefixList(t *testing.T) {
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
				Config: filtersBgpRouteMapRoutingProfile_Match_Ipv4_Address_PrefixList_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_filters_bgp_route_map_routing_profile.example",
						tfjsonpath.New("route_map").AtSliceIndex(0).AtMapKey("match").AtMapKey("ipv4").AtMapKey("address").AtMapKey("prefix_list"),
						knownvalue.StringExact(prefix+"-pfl"),
					),
				},
			},
		},
	})
}

const filtersBgpRouteMapRoutingProfile_Match_Ipv4_Address_PrefixList_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_filters_prefix_list_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = "${var.prefix}-pfl"
  type = {
    ipv4 = {
      ipv4_entries = [
        {
          name = "10"
          action = "deny"
          prefix = {
            entry = {
              network = "10.0.0.0/8"
            }
          }
        }
      ]
    }
  }
}

resource "panos_filters_bgp_route_map_routing_profile" "example" {
  depends_on = [panos_filters_prefix_list_routing_profile.example]
  location = var.location

  name = var.prefix
  route_map = [
    {
      name = "10"
      match = {
        ipv4 = {
          address = {
            prefix_list = panos_filters_prefix_list_routing_profile.example.name
          }
        }
      }
    }
  ]
}
`

func TestAccFiltersBgpRouteMapRoutingProfile_Match_Ipv4_NextHop_AccessList(t *testing.T) {
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
				Config: filtersBgpRouteMapRoutingProfile_Match_Ipv4_NextHop_AccessList_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_filters_bgp_route_map_routing_profile.example",
						tfjsonpath.New("route_map").AtSliceIndex(0).AtMapKey("match").AtMapKey("ipv4").AtMapKey("next_hop").AtMapKey("access_list"),
						knownvalue.StringExact(prefix+"-acl"),
					),
				},
			},
		},
	})
}

const filtersBgpRouteMapRoutingProfile_Match_Ipv4_NextHop_AccessList_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_filters_access_list_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = "${var.prefix}-acl"
  type = {
    ipv4 = {
      ipv4_entries = [
        {
          name = "10"
          action = "permit"
          source_address = {
            address = "any"
          }
        }
      ]
    }
  }
}

resource "panos_filters_bgp_route_map_routing_profile" "example" {
  depends_on = [panos_filters_access_list_routing_profile.example]
  location = var.location

  name = var.prefix
  route_map = [
    {
      name = "10"
      match = {
        ipv4 = {
          next_hop = {
            access_list = panos_filters_access_list_routing_profile.example.name
          }
        }
      }
    }
  ]
}
`

func TestAccFiltersBgpRouteMapRoutingProfile_Match_Ipv4_NextHop_PrefixList(t *testing.T) {
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
				Config: filtersBgpRouteMapRoutingProfile_Match_Ipv4_NextHop_PrefixList_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_filters_bgp_route_map_routing_profile.example",
						tfjsonpath.New("route_map").AtSliceIndex(0).AtMapKey("match").AtMapKey("ipv4").AtMapKey("next_hop").AtMapKey("prefix_list"),
						knownvalue.StringExact(prefix+"-pfl"),
					),
				},
			},
		},
	})
}

const filtersBgpRouteMapRoutingProfile_Match_Ipv4_NextHop_PrefixList_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_filters_prefix_list_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = "${var.prefix}-pfl"
  type = {
    ipv4 = {
      ipv4_entries = [
        {
          name = "10"
          action = "deny"
          prefix = {
            entry = {
              network = "10.0.0.0/8"
            }
          }
        }
      ]
    }
  }
}

resource "panos_filters_bgp_route_map_routing_profile" "example" {
  depends_on = [panos_filters_prefix_list_routing_profile.example]
  location = var.location

  name = var.prefix
  route_map = [
    {
      name = "10"
      match = {
        ipv4 = {
          next_hop = {
            prefix_list = panos_filters_prefix_list_routing_profile.example.name
          }
        }
      }
    }
  ]
}
`

func TestAccFiltersBgpRouteMapRoutingProfile_Match_Ipv4_RouteSource_AccessList(t *testing.T) {
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
				Config: filtersBgpRouteMapRoutingProfile_Match_Ipv4_RouteSource_AccessList_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_filters_bgp_route_map_routing_profile.example",
						tfjsonpath.New("route_map").AtSliceIndex(0).AtMapKey("match").AtMapKey("ipv4").AtMapKey("route_source").AtMapKey("access_list"),
						knownvalue.StringExact(prefix+"-acl"),
					),
				},
			},
		},
	})
}

const filtersBgpRouteMapRoutingProfile_Match_Ipv4_RouteSource_AccessList_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_filters_access_list_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = "${var.prefix}-acl"
  type = {
    ipv4 = {
      ipv4_entries = [
        {
          name = "10"
          action = "permit"
          source_address = {
            address = "any"
          }
        }
      ]
    }
  }
}

resource "panos_filters_bgp_route_map_routing_profile" "example" {
  depends_on = [panos_filters_access_list_routing_profile.example]
  location = var.location

  name = var.prefix
  route_map = [
    {
      name = "10"
      match = {
        ipv4 = {
          route_source = {
            access_list = panos_filters_access_list_routing_profile.example.name
          }
        }
      }
    }
  ]
}
`

func TestAccFiltersBgpRouteMapRoutingProfile_Match_Ipv4_RouteSource_PrefixList(t *testing.T) {
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
				Config: filtersBgpRouteMapRoutingProfile_Match_Ipv4_RouteSource_PrefixList_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_filters_bgp_route_map_routing_profile.example",
						tfjsonpath.New("route_map").AtSliceIndex(0).AtMapKey("match").AtMapKey("ipv4").AtMapKey("route_source").AtMapKey("prefix_list"),
						knownvalue.StringExact(prefix+"-pfl"),
					),
				},
			},
		},
	})
}

const filtersBgpRouteMapRoutingProfile_Match_Ipv4_RouteSource_PrefixList_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_filters_prefix_list_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = "${var.prefix}-pfl"
  type = {
    ipv4 = {
      ipv4_entries = [
        {
          name = "10"
          action = "deny"
          prefix = {
            entry = {
              network = "10.0.0.0/8"
            }
          }
        }
      ]
    }
  }
}

resource "panos_filters_bgp_route_map_routing_profile" "example" {
  depends_on = [panos_filters_prefix_list_routing_profile.example]
  location = var.location

  name = var.prefix
  route_map = [
    {
      name = "10"
      match = {
        ipv4 = {
          route_source = {
            prefix_list = panos_filters_prefix_list_routing_profile.example.name
          }
        }
      }
    }
  ]
}
`

func TestAccFiltersBgpRouteMapRoutingProfile_Match_Ipv6_Address_AccessList(t *testing.T) {
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
				Config: filtersBgpRouteMapRoutingProfile_Match_Ipv6_Address_AccessList_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_filters_bgp_route_map_routing_profile.example",
						tfjsonpath.New("route_map").AtSliceIndex(0).AtMapKey("match").AtMapKey("ipv6").AtMapKey("address").AtMapKey("access_list"),
						knownvalue.StringExact(prefix+"-acl"),
					),
				},
			},
		},
	})
}

const filtersBgpRouteMapRoutingProfile_Match_Ipv6_Address_AccessList_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_filters_access_list_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = "${var.prefix}-acl"
  type = {
    ipv6 = {
      ipv6_entries = [
        {
          name = "10"
          action = "permit"
          source_address = {
            address = "any"
          }
        }
      ]
    }
  }
}

resource "panos_filters_bgp_route_map_routing_profile" "example" {
  depends_on = [panos_filters_access_list_routing_profile.example]
  location = var.location

  name = var.prefix
  route_map = [
    {
      name = "10"
      match = {
        ipv6 = {
          address = {
            access_list = panos_filters_access_list_routing_profile.example.name
          }
        }
      }
    }
  ]
}
`

func TestAccFiltersBgpRouteMapRoutingProfile_Match_Ipv6_Address_PrefixList(t *testing.T) {
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
				Config: filtersBgpRouteMapRoutingProfile_Match_Ipv6_Address_PrefixList_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_filters_bgp_route_map_routing_profile.example",
						tfjsonpath.New("route_map").AtSliceIndex(0).AtMapKey("match").AtMapKey("ipv6").AtMapKey("address").AtMapKey("prefix_list"),
						knownvalue.StringExact(prefix+"-pfl"),
					),
				},
			},
		},
	})
}

const filtersBgpRouteMapRoutingProfile_Match_Ipv6_Address_PrefixList_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_filters_prefix_list_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = "${var.prefix}-pfl"
  type = {
    ipv6 = {
      ipv6_entries = [
        {
          name = "10"
          action = "deny"
          prefix = {
            entry = {
              network = "2001:db8::/32"
            }
          }
        }
      ]
    }
  }
}

resource "panos_filters_bgp_route_map_routing_profile" "example" {
  depends_on = [panos_filters_prefix_list_routing_profile.example]
  location = var.location

  name = var.prefix
  route_map = [
    {
      name = "10"
      match = {
        ipv6 = {
          address = {
            prefix_list = panos_filters_prefix_list_routing_profile.example.name
          }
        }
      }
    }
  ]
}
`

func TestAccFiltersBgpRouteMapRoutingProfile_Match_Ipv6_NextHop_AccessList(t *testing.T) {
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
				Config: filtersBgpRouteMapRoutingProfile_Match_Ipv6_NextHop_AccessList_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_filters_bgp_route_map_routing_profile.example",
						tfjsonpath.New("route_map").AtSliceIndex(0).AtMapKey("match").AtMapKey("ipv6").AtMapKey("next_hop").AtMapKey("access_list"),
						knownvalue.StringExact(prefix+"-acl"),
					),
				},
			},
		},
	})
}

const filtersBgpRouteMapRoutingProfile_Match_Ipv6_NextHop_AccessList_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_filters_access_list_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = "${var.prefix}-acl"
  type = {
    ipv6 = {
      ipv6_entries = [
        {
          name = "10"
          action = "permit"
          source_address = {
            address = "any"
          }
        }
      ]
    }
  }
}

resource "panos_filters_bgp_route_map_routing_profile" "example" {
  depends_on = [panos_filters_access_list_routing_profile.example]
  location = var.location

  name = var.prefix
  route_map = [
    {
      name = "10"
      match = {
        ipv6 = {
          next_hop = {
            access_list = panos_filters_access_list_routing_profile.example.name
          }
        }
      }
    }
  ]
}
`

func TestAccFiltersBgpRouteMapRoutingProfile_Match_Ipv6_NextHop_PrefixList(t *testing.T) {
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
				Config: filtersBgpRouteMapRoutingProfile_Match_Ipv6_NextHop_PrefixList_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_filters_bgp_route_map_routing_profile.example",
						tfjsonpath.New("route_map").AtSliceIndex(0).AtMapKey("match").AtMapKey("ipv6").AtMapKey("next_hop").AtMapKey("prefix_list"),
						knownvalue.StringExact(prefix+"-pfl"),
					),
				},
			},
		},
	})
}

const filtersBgpRouteMapRoutingProfile_Match_Ipv6_NextHop_PrefixList_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_filters_prefix_list_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = "${var.prefix}-pfl"
  type = {
    ipv6 = {
      ipv6_entries = [
        {
          name = "10"
          action = "deny"
          prefix = {
            entry = {
              network = "2001:db8::/32"
            }
          }
        }
      ]
    }
  }
}

resource "panos_filters_bgp_route_map_routing_profile" "example" {
  depends_on = [panos_filters_prefix_list_routing_profile.example]
  location = var.location

  name = var.prefix
  route_map = [
    {
      name = "10"
      match = {
        ipv6 = {
          next_hop = {
            prefix_list = panos_filters_prefix_list_routing_profile.example.name
          }
        }
      }
    }
  ]
}
`
