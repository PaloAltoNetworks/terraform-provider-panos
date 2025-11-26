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

func TestAccBgpAddressFamilyRoutingProfile_Ipv4_Multicast_Basic(t *testing.T) {
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
				Config: bgpAddressFamilyRoutingProfile_Ipv4_Multicast_Basic_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_bgp_address_family_routing_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_bgp_address_family_routing_profile.example",
						tfjsonpath.New("ipv4"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"multicast": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"enable":                           knownvalue.Bool(true),
								"as_override":                      knownvalue.Bool(true),
								"default_originate":                knownvalue.Bool(true),
								"default_originate_map":            knownvalue.Null(),
								"route_reflector_client":           knownvalue.Bool(true),
								"soft_reconfig_with_stored_info":   knownvalue.Bool(true),
								"add_path": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"tx_all_paths":       knownvalue.Bool(true),
									"tx_bestpath_per_as": knownvalue.Bool(false),
								}),
								"orf": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"orf_prefix_list": knownvalue.StringExact("both"),
								}),
								"allowas_in":          knownvalue.Null(),
								"maximum_prefix":      knownvalue.Null(),
								"next_hop":            knownvalue.Null(),
								"remove_private_as":   knownvalue.Null(),
								"send_community":      knownvalue.Null(),
							}),
							"unicast": knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

func TestAccBgpAddressFamilyRoutingProfile_Ipv4_Multicast_AllowasIn_Occurrence(t *testing.T) {
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
				Config: bgpAddressFamilyRoutingProfile_Ipv4_Multicast_AllowasIn_Occurrence_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_bgp_address_family_routing_profile.example",
						tfjsonpath.New("ipv4").AtMapKey("multicast"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"enable":                         knownvalue.Null(),
							"as_override":                    knownvalue.Null(),
							"default_originate":              knownvalue.Null(),
							"default_originate_map":          knownvalue.Null(),
							"route_reflector_client":         knownvalue.Null(),
							"soft_reconfig_with_stored_info": knownvalue.Null(),
							"add_path":                       knownvalue.Null(),
							"orf":                            knownvalue.Null(),
							"allowas_in": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"occurrence": knownvalue.Int64Exact(5),
								"origin":     knownvalue.Null(),
							}),
							"maximum_prefix": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"max_prefixes": knownvalue.Int64Exact(5000),
								"threshold":    knownvalue.Int64Exact(80),
								"action": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"restart": knownvalue.ObjectExact(map[string]knownvalue.Check{
										"interval": knownvalue.Int64Exact(10),
									}),
									"warning_only": knownvalue.Null(),
								}),
							}),
							"next_hop": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"self":       knownvalue.ObjectExact(map[string]knownvalue.Check{}),
								"self_force": knownvalue.Null(),
							}),
							"remove_private_as": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"all":        knownvalue.ObjectExact(map[string]knownvalue.Check{}),
								"replace_as": knownvalue.Null(),
							}),
							"send_community": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"all":      knownvalue.ObjectExact(map[string]knownvalue.Check{}),
								"both":     knownvalue.Null(),
								"extended": knownvalue.Null(),
								"large":    knownvalue.Null(),
								"standard": knownvalue.Null(),
							}),
						}),
					),
				},
			},
		},
	})
}

func TestAccBgpAddressFamilyRoutingProfile_Ipv4_Multicast_AllowasIn_Origin(t *testing.T) {
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
				Config: bgpAddressFamilyRoutingProfile_Ipv4_Multicast_AllowasIn_Origin_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_bgp_address_family_routing_profile.example",
						tfjsonpath.New("ipv4").AtMapKey("multicast"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"enable":                         knownvalue.Null(),
							"as_override":                    knownvalue.Null(),
							"default_originate":              knownvalue.Null(),
							"default_originate_map":          knownvalue.Null(),
							"route_reflector_client":         knownvalue.Null(),
							"soft_reconfig_with_stored_info": knownvalue.Null(),
							"add_path":                       knownvalue.Null(),
							"orf":                            knownvalue.Null(),
							"allowas_in": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"occurrence": knownvalue.Null(),
								"origin":     knownvalue.ObjectExact(map[string]knownvalue.Check{}),
							}),
							"maximum_prefix": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"max_prefixes": knownvalue.Int64Exact(10000),
								"threshold":    knownvalue.Int64Exact(90),
								"action": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"restart":      knownvalue.Null(),
									"warning_only": knownvalue.ObjectExact(map[string]knownvalue.Check{}),
								}),
							}),
							"next_hop": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"self":       knownvalue.Null(),
								"self_force": knownvalue.ObjectExact(map[string]knownvalue.Check{}),
							}),
							"remove_private_as": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"all":        knownvalue.Null(),
								"replace_as": knownvalue.ObjectExact(map[string]knownvalue.Check{}),
							}),
							"send_community": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"all":      knownvalue.Null(),
								"both":     knownvalue.ObjectExact(map[string]knownvalue.Check{}),
								"extended": knownvalue.Null(),
								"large":    knownvalue.Null(),
								"standard": knownvalue.Null(),
							}),
						}),
					),
				},
			},
		},
	})
}

func TestAccBgpAddressFamilyRoutingProfile_Ipv4_Multicast_SendCommunity_Extended(t *testing.T) {
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
				Config: bgpAddressFamilyRoutingProfile_Ipv4_Multicast_SendCommunity_Extended_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_bgp_address_family_routing_profile.example",
						tfjsonpath.New("ipv4").AtMapKey("multicast").AtMapKey("send_community"),
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

func TestAccBgpAddressFamilyRoutingProfile_Ipv4_Multicast_SendCommunity_Large(t *testing.T) {
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
				Config: bgpAddressFamilyRoutingProfile_Ipv4_Multicast_SendCommunity_Large_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_bgp_address_family_routing_profile.example",
						tfjsonpath.New("ipv4").AtMapKey("multicast").AtMapKey("send_community"),
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

func TestAccBgpAddressFamilyRoutingProfile_Ipv4_Multicast_SendCommunity_Standard(t *testing.T) {
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
				Config: bgpAddressFamilyRoutingProfile_Ipv4_Multicast_SendCommunity_Standard_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_bgp_address_family_routing_profile.example",
						tfjsonpath.New("ipv4").AtMapKey("multicast").AtMapKey("send_community"),
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

const bgpAddressFamilyRoutingProfile_Ipv4_Multicast_SendCommunity_Extended_Tmpl = `
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

  ipv4 = {
    multicast = {
      send_community = {
        extended = {}
      }
    }
  }
}
`

const bgpAddressFamilyRoutingProfile_Ipv4_Multicast_SendCommunity_Large_Tmpl = `
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

  ipv4 = {
    multicast = {
      send_community = {
        large = {}
      }
    }
  }
}
`

const bgpAddressFamilyRoutingProfile_Ipv4_Multicast_SendCommunity_Standard_Tmpl = `
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

  ipv4 = {
    multicast = {
      send_community = {
        standard = {}
      }
    }
  }
}
`

const bgpAddressFamilyRoutingProfile_Ipv4_Multicast_AllowasIn_Origin_Tmpl = `
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

  ipv4 = {
    multicast = {
      allowas_in = {
        origin = {}
      }
      maximum_prefix = {
        max_prefixes = 10000
        threshold = 90
        action = {
          warning_only = {}
        }
      }
      next_hop = {
        self_force = {}
      }
      remove_private_as = {
        replace_as = {}
      }
      send_community = {
        both = {}
      }
    }
  }
}
`

const bgpAddressFamilyRoutingProfile_Ipv4_Multicast_AllowasIn_Occurrence_Tmpl = `
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

  ipv4 = {
    multicast = {
      allowas_in = {
        occurrence = 5
      }
      maximum_prefix = {
        max_prefixes = 5000
        threshold = 80
        action = {
          restart = {
            interval = 10
          }
        }
      }
      next_hop = {
        self = {}
      }
      remove_private_as = {
        all = {}
      }
      send_community = {
        all = {}
      }
    }
  }
}
`

func TestAccBgpAddressFamilyRoutingProfile_Ipv4_Unicast_Basic(t *testing.T) {
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
				Config: bgpAddressFamilyRoutingProfile_Ipv4_Unicast_Basic_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_bgp_address_family_routing_profile.example",
						tfjsonpath.New("ipv4"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"multicast": knownvalue.Null(),
							"unicast": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"enable":                           knownvalue.Bool(true),
								"as_override":                      knownvalue.Bool(true),
								"default_originate":                knownvalue.Bool(true),
								"default_originate_map":            knownvalue.Null(),
								"route_reflector_client":           knownvalue.Bool(true),
								"soft_reconfig_with_stored_info":   knownvalue.Bool(true),
								"add_path": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"tx_all_paths":       knownvalue.Bool(false),
									"tx_bestpath_per_as": knownvalue.Bool(true),
								}),
								"orf": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"orf_prefix_list": knownvalue.StringExact("send"),
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

func TestAccBgpAddressFamilyRoutingProfile_Ipv4_Unicast_AllowasIn_Occurrence(t *testing.T) {
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
				Config: bgpAddressFamilyRoutingProfile_Ipv4_Unicast_AllowasIn_Occurrence_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_bgp_address_family_routing_profile.example",
						tfjsonpath.New("ipv4").AtMapKey("unicast").AtMapKey("allowas_in"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"occurrence": knownvalue.Int64Exact(3),
							"origin":     knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

func TestAccBgpAddressFamilyRoutingProfile_Ipv4_Unicast_AllowasIn_Origin(t *testing.T) {
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
				Config: bgpAddressFamilyRoutingProfile_Ipv4_Unicast_AllowasIn_Origin_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_bgp_address_family_routing_profile.example",
						tfjsonpath.New("ipv4").AtMapKey("unicast").AtMapKey("allowas_in"),
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

func TestAccBgpAddressFamilyRoutingProfile_Ipv4_Unicast_SendCommunity_Extended(t *testing.T) {
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
				Config: bgpAddressFamilyRoutingProfile_Ipv4_Unicast_SendCommunity_Extended_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_bgp_address_family_routing_profile.example",
						tfjsonpath.New("ipv4").AtMapKey("unicast").AtMapKey("send_community"),
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

func TestAccBgpAddressFamilyRoutingProfile_Ipv4_Unicast_SendCommunity_Large(t *testing.T) {
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
				Config: bgpAddressFamilyRoutingProfile_Ipv4_Unicast_SendCommunity_Large_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_bgp_address_family_routing_profile.example",
						tfjsonpath.New("ipv4").AtMapKey("unicast").AtMapKey("send_community"),
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

func TestAccBgpAddressFamilyRoutingProfile_Ipv4_Unicast_SendCommunity_Standard(t *testing.T) {
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
				Config: bgpAddressFamilyRoutingProfile_Ipv4_Unicast_SendCommunity_Standard_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_bgp_address_family_routing_profile.example",
						tfjsonpath.New("ipv4").AtMapKey("unicast").AtMapKey("send_community"),
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

const bgpAddressFamilyRoutingProfile_Ipv4_Unicast_Basic_Tmpl = `
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

  ipv4 = {
    unicast = {
      enable = true
      as_override = true
      default_originate = true
      route_reflector_client = true
      soft_reconfig_with_stored_info = true
      add_path = {
        tx_all_paths = false
        tx_bestpath_per_as = true
      }
      orf = {
        orf_prefix_list = "send"
      }
    }
  }
}
`

const bgpAddressFamilyRoutingProfile_Ipv4_Unicast_AllowasIn_Occurrence_Tmpl = `
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

  ipv4 = {
    unicast = {
      allowas_in = {
        occurrence = 3
      }
    }
  }
}
`

const bgpAddressFamilyRoutingProfile_Ipv4_Unicast_AllowasIn_Origin_Tmpl = `
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

  ipv4 = {
    unicast = {
      allowas_in = {
        origin = {}
      }
    }
  }
}
`

const bgpAddressFamilyRoutingProfile_Ipv4_Unicast_SendCommunity_Extended_Tmpl = `
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

  ipv4 = {
    unicast = {
      send_community = {
        extended = {}
      }
    }
  }
}
`

const bgpAddressFamilyRoutingProfile_Ipv4_Unicast_SendCommunity_Large_Tmpl = `
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

  ipv4 = {
    unicast = {
      send_community = {
        large = {}
      }
    }
  }
}
`

const bgpAddressFamilyRoutingProfile_Ipv4_Unicast_SendCommunity_Standard_Tmpl = `
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

  ipv4 = {
    unicast = {
      send_community = {
        standard = {}
      }
    }
  }
}
`

const bgpAddressFamilyRoutingProfile_Ipv4_Multicast_Basic_Tmpl = `
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

  ipv4 = {
    multicast = {
      enable = true
      as_override = true
      default_originate = true
      route_reflector_client = true
      soft_reconfig_with_stored_info = true
      add_path = {
        tx_all_paths = true
        tx_bestpath_per_as = false
      }
      orf = {
        orf_prefix_list = "both"
      }
    }
  }
}
`
