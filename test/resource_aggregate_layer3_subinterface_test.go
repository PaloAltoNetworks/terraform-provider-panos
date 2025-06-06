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

func TestAccAggregateLayer3Subinterface_Basic(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: aggregateLayer3Subinterface_Basic,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_aggregate_layer3_subinterface.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact("ae1.1"),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_layer3_subinterface.example",
						tfjsonpath.New("tag"),
						knownvalue.Int64Exact(1),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_layer3_subinterface.example",
						tfjsonpath.New("comment"),
						knownvalue.StringExact("Test aggregate layer3 subinterface"),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_layer3_subinterface.example",
						tfjsonpath.New("ip"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":          knownvalue.StringExact("192.0.2.1/24"),
								"sdwan_gateway": knownvalue.StringExact("10.0.0.1"),
							}),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_layer3_subinterface.example",
						tfjsonpath.New("dhcp_client"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_layer3_subinterface.example",
						tfjsonpath.New("adjust_tcp_mss"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"enable":              knownvalue.Bool(true),
							"ipv4_mss_adjustment": knownvalue.Int64Exact(40),
							"ipv6_mss_adjustment": knownvalue.Int64Exact(60),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_layer3_subinterface.example",
						tfjsonpath.New("arp"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":       knownvalue.StringExact("192.0.2.1"),
								"hw_address": knownvalue.StringExact("00:1a:2b:3c:4d:5e"),
							}),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_layer3_subinterface.example",
						tfjsonpath.New("bonjour"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"enable":    knownvalue.Bool(true),
							"group_id":  knownvalue.Int64Exact(0),
							"ttl_check": knownvalue.Bool(true),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_layer3_subinterface.example",
						tfjsonpath.New("decrypt_forward"),
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_layer3_subinterface.example",
						tfjsonpath.New("df_ignore"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_layer3_subinterface.example",
						tfjsonpath.New("mtu"),
						knownvalue.Int64Exact(1500),
					),
				},
			},
		},
	})
}

const aggregateLayer3Subinterface_Basic = `
variable "prefix" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = "${var.prefix}-tmpl"
}

resource "panos_aggregate_interface" "parent" {
  location = {
    template = {
      vsys = "vsys1"
      name = panos_template.example.name
    }
  }

  name = "ae1"
  layer3 = {}
}

resource "panos_interface_management_profile" "profile" {
  location = { template = { name = panos_template.example.name } }

  name = var.prefix
}

resource "panos_aggregate_layer3_subinterface" "example" {
  location = { template = { name = panos_template.example.name, vsys = "vsys1" } }

  parent = panos_aggregate_interface.parent.name
  name = "ae1.1"
  tag = 1

  comment = "Test aggregate layer3 subinterface"

  adjust_tcp_mss = {
    enable = true
    ipv4_mss_adjustment = 40
    ipv6_mss_adjustment = 60
  }

  arp = [{
    name = "192.0.2.1"
    hw_address = "00:1a:2b:3c:4d:5e"
  }]

  bonjour = {
    enable = true
    group_id = 0
    ttl_check = true
  }

  decrypt_forward = false
  df_ignore = true

  interface_management_profile = panos_interface_management_profile.profile.name

  ip = [{
    name = "192.0.2.1/24"
    sdwan_gateway = "10.0.0.1"
  }]

  ipv6 = {
    enabled = true
    interface_id = "EUI-64"

    neighbor_discovery = {
      dad_attempts = 1
      enable_dad = true
      enable_ndp_monitor = true
      ns_interval = 1
      reachable_time = 30
      router_advertisement = {
        enable = true
        hop_limit = "64"
        lifetime = 1800
        managed_flag = false
        max_interval = 600
        min_interval = 200
        other_flag = false
        router_preference = "Medium"
      }
    }

    ndp_proxy = {
      enabled = false
    }
  }

  mtu = 1500

  # Outermost variant attributes set to null
  dhcp_client = null
  sdwan_link_settings = null
}
`

func TestAccAggregateLayer3Subinterface_DhcpClient(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: aggregateLayer3Subinterface_DhcpClient,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_aggregate_layer3_subinterface.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact("ae1.1"),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_layer3_subinterface.example",
						tfjsonpath.New("tag"),
						knownvalue.Int64Exact(1),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_layer3_subinterface.example",
						tfjsonpath.New("comment"),
						knownvalue.StringExact("Test aggregate layer3 subinterface with DHCP client"),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_layer3_subinterface.example",
						tfjsonpath.New("ip"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_layer3_subinterface.example",
						tfjsonpath.New("dhcp_client"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"enable":               knownvalue.Bool(true),
							"create_default_route": knownvalue.Bool(true),
							"default_route_metric": knownvalue.Int64Exact(10),
							"send_hostname": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"enable":   knownvalue.Bool(true),
								"hostname": knownvalue.StringExact("system-hostname"),
							}),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_layer3_subinterface.example",
						tfjsonpath.New("interface_management_profile"),
						knownvalue.StringExact(prefix),
					),
				},
			},
		},
	})
}

const aggregateLayer3Subinterface_DhcpClient = `
variable "prefix" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = "${var.prefix}-tmpl"
}

resource "panos_aggregate_interface" "parent" {
  location = {
    template = {
      vsys = "vsys1"
      name = panos_template.example.name
    }
  }

  name = "ae1"
  layer3 = {}
}

resource "panos_interface_management_profile" "profile" {
  location = { template = { name = panos_template.example.name } }

  name = var.prefix
}

resource "panos_aggregate_layer3_subinterface" "example" {
  location = { template = { name = panos_template.example.name, vsys = "vsys1" } }

  parent = panos_aggregate_interface.parent.name
  name = "ae1.1"
  tag = 1

  comment = "Test aggregate layer3 subinterface with DHCP client"

  dhcp_client = {
    enable = true
    create_default_route = true
    default_route_metric = 10
    send_hostname = {
      enable = true
      hostname = "system-hostname"
    }
  }

  ip = null

  interface_management_profile = panos_interface_management_profile.profile.name
}
`

const aggregateLayer3Subinterface_Ipv6_Address = `
variable "prefix" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = "${var.prefix}-tmpl"
}

resource "panos_aggregate_interface" "parent" {
  location = {
    template = {
      vsys = "vsys1"
      name = panos_template.example.name
    }
  }

  name = "ae1"
  layer3 = {}
}

resource "panos_interface_management_profile" "profile" {
  location = { template = { name = panos_template.example.name } }

  name = var.prefix
}

resource "panos_aggregate_layer3_subinterface" "example" {
  location = { template = { name = panos_template.example.name, vsys = "vsys1" } }

  parent = panos_aggregate_interface.parent.name
  name = "ae1.1"
  tag = 1

  comment = "Test aggregate layer3 subinterface with IPv6 address"

  ipv6 = {
    enabled = true
    interface_id = "EUI-64"
    address = [{
      name = "2001:db8::1/64"
      enable_on_interface = true
      prefix = {}
      anycast = {}
      advertise = {
        enable = true
        valid_lifetime = "2592000"
        preferred_lifetime = "604800"
        onlink_flag = true
        auto_config_flag = true
      }
    }]
    dhcp_client = null
    inherited = null
    neighbor_discovery = null
  }

  ip = null
  dhcp_client = null
  mtu = 1500
  ndp_proxy = {
    enabled = true
    address = [{
      name = "2001:db8::/64"
      negate = false
    }]
  }

  interface_management_profile = panos_interface_management_profile.profile.name
}
`

func TestAccAggregateLayer3Subinterface_Ipv6_Address(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: aggregateLayer3Subinterface_Ipv6_Address,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_aggregate_layer3_subinterface.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact("ae1.1"),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_layer3_subinterface.example",
						tfjsonpath.New("tag"),
						knownvalue.Int64Exact(1),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_layer3_subinterface.example",
						tfjsonpath.New("comment"),
						knownvalue.StringExact("Test aggregate layer3 subinterface with IPv6 address"),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_layer3_subinterface.example",
						tfjsonpath.New("ip"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_layer3_subinterface.example",
						tfjsonpath.New("dhcp_client"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_layer3_subinterface.example",
						tfjsonpath.New("ipv6"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"enabled":      knownvalue.Bool(true),
							"interface_id": knownvalue.StringExact("EUI-64"),
							"address": knownvalue.ListExact([]knownvalue.Check{
								knownvalue.ObjectExact(map[string]knownvalue.Check{
									"name":                knownvalue.StringExact("2001:db8::1/64"),
									"enable_on_interface": knownvalue.Bool(true),
									"prefix":              knownvalue.ObjectExact(map[string]knownvalue.Check{}),
									"anycast":             knownvalue.ObjectExact(map[string]knownvalue.Check{}),
									"advertise": knownvalue.ObjectExact(map[string]knownvalue.Check{
										"enable":             knownvalue.Bool(true),
										"valid_lifetime":     knownvalue.StringExact("2592000"),
										"preferred_lifetime": knownvalue.StringExact("604800"),
										"onlink_flag":        knownvalue.Bool(true),
										"auto_config_flag":   knownvalue.Bool(true),
									}),
								}),
							}),
							"dhcp_client":        knownvalue.Null(),
							"inherited":          knownvalue.Null(),
							"neighbor_discovery": knownvalue.Null(),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_layer3_subinterface.example",
						tfjsonpath.New("mtu"),
						knownvalue.Int64Exact(1500),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_layer3_subinterface.example",
						tfjsonpath.New("ndp_proxy"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"enabled": knownvalue.Bool(true),
							"address": knownvalue.ListExact([]knownvalue.Check{
								knownvalue.ObjectExact(map[string]knownvalue.Check{
									"name":   knownvalue.StringExact("2001:db8::/64"),
									"negate": knownvalue.Bool(false),
								}),
							}),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_layer3_subinterface.example",
						tfjsonpath.New("interface_management_profile"),
						knownvalue.StringExact(prefix),
					),
				},
			},
		},
	})
}

const aggregateLayer3Subinterface_Ipv6_Inherited = `
variable "prefix" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = "${var.prefix}-tmpl"
}

resource "panos_aggregate_interface" "parent" {
  location = {
    template = {
      vsys = "vsys1"
      name = panos_template.example.name
    }
  }

  name = "ae1"
  layer3 = {}
}

resource "panos_interface_management_profile" "profile" {
  location = { template = { name = panos_template.example.name } }

  name = var.prefix
}

resource "panos_aggregate_layer3_subinterface" "example" {
  location = { template = { name = panos_template.example.name, vsys = "vsys1" } }

  parent = panos_aggregate_interface.parent.name
  name = "ae1.1"
  tag = 1

  comment = "Test aggregate layer3 subinterface with IPv6 inherited configuration"

  ipv6 = {
    enabled = true
    interface_id = "EUI-64"
    inherited = {
      enable = true
      assign_addr = [
        {
          name = "gua-address-dynamic"
          type = {
            gua = {
              enable_on_interface = true
              # prefix_pool = "ipv6-pool"  # prefix-pool 'ipv6-pool' is not a valid reference
              pool_type = {
                dynamic = {}
                dynamic_id = null
              }
              advertise = {
                enable = true
                onlink_flag = true
                auto_config_flag = true
              }
            }
            ula = null
          }
        },
        {
          name = "gua-address-dynamic-id"
          type = {
            gua = {
              enable_on_interface = true
              # prefix_pool = "ipv6-pool"  # prefix-pool 'ipv6-pool' is not a valid reference
              pool_type = {
                dynamic = null
                dynamic_id = {
                  identifier = 1
                }
              }
              advertise = {
                enable = true
                onlink_flag = true
                auto_config_flag = true
              }
            }
            ula = null
          }
        },
        {
          name = "ula-address"
          type = {
            gua = null
            ula = {
              enable_on_interface = true
              address = "fd00:1234:5678::/48"
              prefix = true
              anycast = false
              advertise = {
                enable = true
                valid_lifetime = "2592000"
                preferred_lifetime = "604800"
                onlink_flag = true
                auto_config_flag = true
              }
            }
          }
        }
      ]
      neighbor_discovery = {
        dad_attempts = 1
        enable_dad = true
        ns_interval = 1000
        reachable_time = 30000
        dns_server = {
          enable = true
          source = {
            dhcpv6 = {
              # prefix_pool = "ipv6-pool"  # prefix-pool 'ipv6-pool' is not a valid reference
            }
            manual = null
          }
        }
        dns_suffix = {
          enable = true
          source = {
            dhcpv6 = {
              # prefix_pool = "ipv6-pool"  # prefix-pool 'ipv6-pool' is not a valid reference
            }
            manual = null
          }
        }
        enable_ndp_monitor = true
        neighbor = [{
          name = "2001:db8::1"
          hw_address = "00:1a:2b:3c:4d:5e"
        }]
        router_advertisement = {
          enable = true
          hop_limit = "64"
          lifetime = 1800
          managed_flag = true
          max_interval = 600
          min_interval = 200
          other_flag = true
          router_preference = "Medium"
          enable_consistency_check = true
          link_mtu = "unspecified"
          reachable_time = "unspecified"
          retransmission_timer = "unspecified"
        }
      }
    }
    address = null
    dhcp_client = null
    neighbor_discovery = null
  }

  ip = null
  dhcp_client = null
  mtu = 1500
  ndp_proxy = {
    enabled = false
    address = [{
      name = "2001:db8::/64"
      negate = false
    }]
  }

  interface_management_profile = panos_interface_management_profile.profile.name
}`

func TestAccAggregateLayer3Subinterface_Ipv6_Inherited(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: aggregateLayer3Subinterface_Ipv6_Inherited,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_aggregate_layer3_subinterface.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact("ae1.1"),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_layer3_subinterface.example",
						tfjsonpath.New("tag"),
						knownvalue.Int64Exact(1),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_layer3_subinterface.example",
						tfjsonpath.New("comment"),
						knownvalue.StringExact("Test aggregate layer3 subinterface with IPv6 inherited configuration"),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_layer3_subinterface.example",
						tfjsonpath.New("ip"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_layer3_subinterface.example",
						tfjsonpath.New("dhcp_client"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_layer3_subinterface.example",
						tfjsonpath.New("ipv6"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"enabled":      knownvalue.Bool(true),
							"interface_id": knownvalue.StringExact("EUI-64"),
							"inherited": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"enable": knownvalue.Bool(true),
								"assign_addr": knownvalue.ListExact([]knownvalue.Check{
									knownvalue.ObjectExact(map[string]knownvalue.Check{
										"name": knownvalue.StringExact("gua-address-dynamic"),
										"type": knownvalue.ObjectExact(map[string]knownvalue.Check{
											"gua": knownvalue.ObjectExact(map[string]knownvalue.Check{
												"enable_on_interface": knownvalue.Bool(true),
												"prefix_pool":         knownvalue.Null(),
												"pool_type": knownvalue.ObjectExact(map[string]knownvalue.Check{
													"dynamic":    knownvalue.ObjectExact(map[string]knownvalue.Check{}),
													"dynamic_id": knownvalue.Null(),
												}),
												"advertise": knownvalue.ObjectExact(map[string]knownvalue.Check{
													"enable":           knownvalue.Bool(true),
													"onlink_flag":      knownvalue.Bool(true),
													"auto_config_flag": knownvalue.Bool(true),
												}),
											}),
											"ula": knownvalue.Null(),
										}),
									}),
									knownvalue.ObjectExact(map[string]knownvalue.Check{
										"name": knownvalue.StringExact("gua-address-dynamic-id"),
										"type": knownvalue.ObjectExact(map[string]knownvalue.Check{
											"gua": knownvalue.ObjectExact(map[string]knownvalue.Check{
												"enable_on_interface": knownvalue.Bool(true),
												"prefix_pool":         knownvalue.Null(),
												"pool_type": knownvalue.ObjectExact(map[string]knownvalue.Check{
													"dynamic": knownvalue.Null(),
													"dynamic_id": knownvalue.ObjectExact(map[string]knownvalue.Check{
														"identifier": knownvalue.Int64Exact(1),
													}),
												}),
												"advertise": knownvalue.ObjectExact(map[string]knownvalue.Check{
													"enable":           knownvalue.Bool(true),
													"onlink_flag":      knownvalue.Bool(true),
													"auto_config_flag": knownvalue.Bool(true),
												}),
											}),
											"ula": knownvalue.Null(),
										}),
									}),
									knownvalue.ObjectExact(map[string]knownvalue.Check{
										"name": knownvalue.StringExact("ula-address"),
										"type": knownvalue.ObjectExact(map[string]knownvalue.Check{
											"gua": knownvalue.Null(),
											"ula": knownvalue.ObjectExact(map[string]knownvalue.Check{
												"enable_on_interface": knownvalue.Bool(true),
												"address":             knownvalue.StringExact("fd00:1234:5678::/48"),
												"prefix":              knownvalue.Bool(true),
												"anycast":             knownvalue.Bool(false),
												"advertise": knownvalue.ObjectExact(map[string]knownvalue.Check{
													"enable":             knownvalue.Bool(true),
													"valid_lifetime":     knownvalue.StringExact("2592000"),
													"preferred_lifetime": knownvalue.StringExact("604800"),
													"onlink_flag":        knownvalue.Bool(true),
													"auto_config_flag":   knownvalue.Bool(true),
												}),
											}),
										}),
									}),
								}),
								"neighbor_discovery": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"dad_attempts":   knownvalue.Int64Exact(1),
									"enable_dad":     knownvalue.Bool(true),
									"ns_interval":    knownvalue.Int64Exact(1000),
									"reachable_time": knownvalue.Int64Exact(30000),
									"dns_server": knownvalue.ObjectExact(map[string]knownvalue.Check{
										"enable": knownvalue.Bool(true),
										"source": knownvalue.ObjectExact(map[string]knownvalue.Check{
											"dhcpv6": knownvalue.ObjectExact(map[string]knownvalue.Check{
												"prefix_pool": knownvalue.Null(),
											}),
											"manual": knownvalue.Null(),
										}),
									}),
									"dns_suffix": knownvalue.ObjectExact(map[string]knownvalue.Check{
										"enable": knownvalue.Bool(true),
										"source": knownvalue.ObjectExact(map[string]knownvalue.Check{
											"dhcpv6": knownvalue.ObjectExact(map[string]knownvalue.Check{
												"prefix_pool": knownvalue.Null(),
											}),
											"manual": knownvalue.Null(),
										}),
									}),
									"enable_ndp_monitor": knownvalue.Bool(true),
									"neighbor": knownvalue.ListExact([]knownvalue.Check{
										knownvalue.ObjectExact(map[string]knownvalue.Check{
											"name":       knownvalue.StringExact("2001:db8::1"),
											"hw_address": knownvalue.StringExact("00:1a:2b:3c:4d:5e"),
										}),
									}),
									"router_advertisement": knownvalue.ObjectExact(map[string]knownvalue.Check{
										"enable":                   knownvalue.Bool(true),
										"hop_limit":                knownvalue.StringExact("64"),
										"lifetime":                 knownvalue.Int64Exact(1800),
										"managed_flag":             knownvalue.Bool(true),
										"max_interval":             knownvalue.Int64Exact(600),
										"min_interval":             knownvalue.Int64Exact(200),
										"other_flag":               knownvalue.Bool(true),
										"router_preference":        knownvalue.StringExact("Medium"),
										"enable_consistency_check": knownvalue.Bool(true),
										"link_mtu":                 knownvalue.StringExact("unspecified"),
										"reachable_time":           knownvalue.StringExact("unspecified"),
										"retransmission_timer":     knownvalue.StringExact("unspecified"),
									}),
								}),
							}),
							"address":            knownvalue.Null(),
							"dhcp_client":        knownvalue.Null(),
							"neighbor_discovery": knownvalue.Null(),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_layer3_subinterface.example",
						tfjsonpath.New("mtu"),
						knownvalue.Int64Exact(1500),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_layer3_subinterface.example",
						tfjsonpath.New("ndp_proxy"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"enabled": knownvalue.Bool(false),
							"address": knownvalue.ListExact([]knownvalue.Check{
								knownvalue.ObjectExact(map[string]knownvalue.Check{
									"name":   knownvalue.StringExact("2001:db8::/64"),
									"negate": knownvalue.Bool(false),
								}),
							}),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_layer3_subinterface.example",
						tfjsonpath.New("interface_management_profile"),
						knownvalue.StringExact(prefix),
					),
				},
			},
		},
	})
}

const aggregateLayer3Subinterface_Ipv6_Inherited_NeighborDiscovery_DnsServer_Manual = `
variable "prefix" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = "${var.prefix}-tmpl"
}

resource "panos_aggregate_interface" "parent" {
  location = {
    template = {
      vsys = "vsys1"
      name = panos_template.example.name
    }
  }

  name = "ae1"
  layer3 = {}
}

resource "panos_interface_management_profile" "profile" {
  location = { template = { name = panos_template.example.name } }

  name = var.prefix
}

resource "panos_aggregate_layer3_subinterface" "example" {
  location = { template = { name = panos_template.example.name, vsys = "vsys1" } }

  parent = panos_aggregate_interface.parent.name
  name = "ae1.1"
  tag = 1

  comment = "Test aggregate layer3 subinterface with IPv6 inherited neighbor discovery manual DNS server"

  ipv6 = {
    enabled = true
    interface_id = "EUI-64"
    inherited = {
      enable = true
      assign_addr = [{
        name = "ula-address"
        type = {
          gua = null
          ula = {
            enable_on_interface = true
            address = "fd00:1234:5678::/48"
            prefix = true
            anycast = false
            advertise = {
              enable = true
              valid_lifetime = "2592000"
              preferred_lifetime = "604800"
              onlink_flag = true
              auto_config_flag = true
            }
          }
        }
      }]
      neighbor_discovery = {
        dns_server = {
          enable = true
          source = {
            dhcpv6 = null
            manual = {
              server = [{
                name = "2001:db8::53"
                lifetime = 1200
              }]
            }
          }
        }
        dns_suffix = {
          enable = true
          source = {
            dhcpv6 = null
            manual = {
              suffix = [{
                name = "example.com"
                lifetime = 1200
              }]
            }
          }
        }
      }
    }
  }

  interface_management_profile = panos_interface_management_profile.profile.name
}
`

func TestAccAggregateLayer3Subinterface_Ipv6_Inherited_NeighborDiscovery_DnsServer_Manual(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: aggregateLayer3Subinterface_Ipv6_Inherited_NeighborDiscovery_DnsServer_Manual,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_aggregate_layer3_subinterface.example",
						tfjsonpath.New("ipv6").AtMapKey("inherited").AtMapKey("assign_addr"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name": knownvalue.StringExact("ula-address"),
								"type": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"gua": knownvalue.Null(),
									"ula": knownvalue.ObjectExact(map[string]knownvalue.Check{
										"enable_on_interface": knownvalue.Bool(true),
										"address":             knownvalue.StringExact("fd00:1234:5678::/48"),
										"prefix":              knownvalue.Bool(true),
										"anycast":             knownvalue.Bool(false),
										"advertise": knownvalue.ObjectExact(map[string]knownvalue.Check{
											"enable":             knownvalue.Bool(true),
											"valid_lifetime":     knownvalue.StringExact("2592000"),
											"preferred_lifetime": knownvalue.StringExact("604800"),
											"onlink_flag":        knownvalue.Bool(true),
											"auto_config_flag":   knownvalue.Bool(true),
										}),
									}),
								}),
							}),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_layer3_subinterface.example",
						tfjsonpath.New("ipv6").AtMapKey("inherited").AtMapKey("neighbor_discovery").AtMapKey("dns_server"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"enable": knownvalue.Bool(true),
							"source": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"dhcpv6": knownvalue.Null(),
								"manual": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"server": knownvalue.ListExact([]knownvalue.Check{
										knownvalue.ObjectExact(map[string]knownvalue.Check{
											"name":     knownvalue.StringExact("2001:db8::53"),
											"lifetime": knownvalue.Int64Exact(1200),
										}),
									}),
								}),
							}),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_layer3_subinterface.example",
						tfjsonpath.New("ipv6").AtMapKey("inherited").AtMapKey("neighbor_discovery").AtMapKey("dns_suffix"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"enable": knownvalue.Bool(true),
							"source": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"dhcpv6": knownvalue.Null(),
								"manual": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"suffix": knownvalue.ListExact([]knownvalue.Check{
										knownvalue.ObjectExact(map[string]knownvalue.Check{
											"name":     knownvalue.StringExact("example.com"),
											"lifetime": knownvalue.Int64Exact(1200),
										}),
									}),
								}),
							}),
						}),
					),
				},
			},
		},
	})
}

const aggregateLayer3Subinterface_Ipv6_NeighborDiscovery = `
variable "prefix" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = "${var.prefix}-tmpl"
}

resource "panos_aggregate_interface" "parent" {
  location = {
    template = {
      vsys = "vsys1"
      name = panos_template.example.name
    }
  }

  name = "ae1"
  layer3 = {}
}

resource "panos_interface_management_profile" "profile" {
  location = { template = { name = panos_template.example.name } }

  name = var.prefix
}

resource "panos_aggregate_layer3_subinterface" "example" {
  location = { template = { name = panos_template.example.name, vsys = "vsys1" } }

  parent = panos_aggregate_interface.parent.name
  name = "ae1.1"
  tag = 1

  comment = "Test aggregate layer3 subinterface with IPv6 neighbor discovery"

  ipv6 = {
    enabled = true
    interface_id = "EUI-64"
    neighbor_discovery = {
      dad_attempts = 1
      enable_dad = true
      ns_interval = 1000
      reachable_time = 30000
      enable_ndp_monitor = true
      neighbor = [{
        name = "2001:db8::1"
        hw_address = "00:1a:2b:3c:4d:5e"
      }]
      router_advertisement = {
        enable = true
        hop_limit = "64"
        lifetime = 1800
        managed_flag = false
        max_interval = 600
        min_interval = 200
        other_flag = false
        router_preference = "Medium"
        enable_consistency_check = true
        link_mtu = "unspecified"
        reachable_time = "unspecified"
        retransmission_timer = "unspecified"
        dns_support = {
          enable = true
          server = [{
            name = "2001:db8::53"
            lifetime = 1200
          }]
          suffix = [{
            name = "example.com"
            lifetime = 1200
          }]
        }
      }
    }
    address = null
    dhcp_client = null
    inherited = null
  }

  ip = null
  dhcp_client = null
  mtu = 1500
  ndp_proxy = {
    enabled = false
  }

  interface_management_profile = panos_interface_management_profile.profile.name
}
`

func TestAccAggregateLayer3Subinterface_Ipv6_NeighborDiscovery(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: aggregateLayer3Subinterface_Ipv6_NeighborDiscovery,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_aggregate_layer3_subinterface.example",
						tfjsonpath.New("ipv6").AtMapKey("neighbor_discovery"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"dad_attempts":       knownvalue.Int64Exact(1),
							"enable_dad":         knownvalue.Bool(true),
							"ns_interval":        knownvalue.Int64Exact(1000),
							"reachable_time":     knownvalue.Int64Exact(30000),
							"enable_ndp_monitor": knownvalue.Bool(true),
							"neighbor": knownvalue.ListExact([]knownvalue.Check{
								knownvalue.ObjectExact(map[string]knownvalue.Check{
									"name":       knownvalue.StringExact("2001:db8::1"),
									"hw_address": knownvalue.StringExact("00:1a:2b:3c:4d:5e"),
								}),
							}),
							"router_advertisement": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"enable":                   knownvalue.Bool(true),
								"hop_limit":                knownvalue.StringExact("64"),
								"lifetime":                 knownvalue.Int64Exact(1800),
								"managed_flag":             knownvalue.Bool(false),
								"max_interval":             knownvalue.Int64Exact(600),
								"min_interval":             knownvalue.Int64Exact(200),
								"other_flag":               knownvalue.Bool(false),
								"router_preference":        knownvalue.StringExact("Medium"),
								"enable_consistency_check": knownvalue.Bool(true),
								"link_mtu":                 knownvalue.StringExact("unspecified"),
								"reachable_time":           knownvalue.StringExact("unspecified"),
								"retransmission_timer":     knownvalue.StringExact("unspecified"),
								"dns_support": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"enable": knownvalue.Bool(true),
									"server": knownvalue.ListExact([]knownvalue.Check{
										knownvalue.ObjectExact(map[string]knownvalue.Check{
											"name":     knownvalue.StringExact("2001:db8::53"),
											"lifetime": knownvalue.Int64Exact(1200),
										}),
									}),
									"suffix": knownvalue.ListExact([]knownvalue.Check{
										knownvalue.ObjectExact(map[string]knownvalue.Check{
											"name":     knownvalue.StringExact("example.com"),
											"lifetime": knownvalue.Int64Exact(1200),
										}),
									}),
								}),
							}),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_layer3_subinterface.example",
						tfjsonpath.New("ipv6").AtMapKey("address"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_layer3_subinterface.example",
						tfjsonpath.New("ipv6").AtMapKey("dhcp_client"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_layer3_subinterface.example",
						tfjsonpath.New("ipv6").AtMapKey("inherited"),
						knownvalue.Null(),
					),
				},
			},
		},
	})
}

const aggregateLayer3Subinterface_SdwanLinkSettings_UpstreamNat_Ddns = `
variable "prefix" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = "${var.prefix}-tmpl"
}

resource "panos_aggregate_interface" "parent" {
  location = {
    template = {
      vsys = "vsys1"
      name = panos_template.example.name
    }
  }

  name = "ae1"
  layer3 = {}
}

resource "panos_interface_management_profile" "profile" {
  location = { template = { name = panos_template.example.name } }

  name = var.prefix
}

resource "panos_aggregate_layer3_subinterface" "example" {
  location = { template = { name = panos_template.example.name, vsys = "vsys1" } }

  parent = panos_aggregate_interface.parent.name
  name = "ae1.1"
  tag = 1

  comment = "Test aggregate layer3 subinterface with SD-WAN link settings and upstream NAT DDNS"

  ddns_config = {
    ddns_enabled = true
    ddns_hostname = "test-hostname"
    ddns_update_interval = 15
    ddns_vendor = "Palo Alto Networks DDNS"
    ddns_cert_profile = "test-cert-profile"
    ddns_ip = ["192.0.2.1", "192.0.2.2"]
    ddns_ipv6 = ["2001:db8::1", "2001:db8::2"]
    ddns_vendor_config = [{
      name = "key1"
      value = "value1"
    }]
  }

  dhcp_client = {
    enable = true
  }

  sdwan_link_settings = {
    enable = true
    sdwan_interface_profile = "test-profile"
    upstream_nat = {
      enable = true
      ddns = {}
    }
  }

  ip = null
  ipv6 = null

  interface_management_profile = panos_interface_management_profile.profile.name
}
`

func TestAccAggregateLayer3Subinterface_SdwanLinkSettings_UpstreamNat_Ddns(t *testing.T) {
	t.Parallel()
	t.Skip("missing required resources: sdwan profile and certificate profile")

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: aggregateLayer3Subinterface_SdwanLinkSettings_UpstreamNat_Ddns,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_aggregate_layer3_subinterface.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact("ae1.1"),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_layer3_subinterface.example",
						tfjsonpath.New("tag"),
						knownvalue.Int64Exact(1),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_layer3_subinterface.example",
						tfjsonpath.New("comment"),
						knownvalue.StringExact("Test aggregate layer3 subinterface with SD-WAN link settings and upstream NAT DDNS"),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_layer3_subinterface.example",
						tfjsonpath.New("ddns_config"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"ddns_enabled":         knownvalue.Bool(true),
							"ddns_hostname":        knownvalue.StringExact("test-hostname"),
							"ddns_update_interval": knownvalue.Int64Exact(15),
							"ddns_vendor":          knownvalue.StringExact("test-vendor"),
							"ddns_cert_profile":    knownvalue.StringExact("test-cert-profile"),
							"ddns_ip": knownvalue.ListExact([]knownvalue.Check{
								knownvalue.StringExact("192.0.2.1"),
								knownvalue.StringExact("192.0.2.2"),
							}),
							"ddns_ipv6": knownvalue.ListExact([]knownvalue.Check{
								knownvalue.StringExact("2001:db8::1"),
								knownvalue.StringExact("2001:db8::2"),
							}),
							"ddns_vendor_config": knownvalue.ListExact([]knownvalue.Check{
								knownvalue.ObjectExact(map[string]knownvalue.Check{
									"name":  knownvalue.StringExact("key1"),
									"value": knownvalue.StringExact("value1"),
								}),
							}),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_layer3_subinterface.example",
						tfjsonpath.New("sdwan_link_settings"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"enable":                  knownvalue.Bool(true),
							"sdwan_interface_profile": knownvalue.StringExact("test-profile"),
							"upstream_nat": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"enable": knownvalue.Bool(true),
								"ddns":   knownvalue.ObjectExact(map[string]knownvalue.Check{}),
							}),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_layer3_subinterface.example",
						tfjsonpath.New("ip"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_layer3_subinterface.example",
						tfjsonpath.New("dhcp_client").AtMapKey("enable"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_layer3_subinterface.example",
						tfjsonpath.New("ipv6"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_layer3_subinterface.example",
						tfjsonpath.New("interface_management_profile"),
						knownvalue.StringExact(prefix),
					),
				},
			},
		},
	})
}

const aggregateLayer3Subinterface_SdwanLinkSettings_UpstreamNat_StaticIp_Fqdn = `
variable "prefix" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = "${var.prefix}-tmpl"
}

resource "panos_aggregate_interface" "parent" {
  location = {
    template = {
      vsys = "vsys1"
      name = panos_template.example.name
    }
  }

  name = "ae1"
  layer3 = {}
}

resource "panos_interface_management_profile" "profile" {
  location = { template = { name = panos_template.example.name } }

  name = var.prefix
}

resource "panos_aggregate_layer3_subinterface" "example" {
  location = { template = { name = panos_template.example.name, vsys = "vsys1" } }

  parent = panos_aggregate_interface.parent.name
  name = "ae1.1"
  tag = 1

  comment = "Test aggregate layer3 subinterface with SD-WAN link settings and upstream NAT static IP FQDN"

  ddns_config = {
    ddns_enabled = true
    ddns_hostname = "test-hostname"
    ddns_update_interval = 15
    ddns_vendor = "test-vendor"
    ddns_cert_profile = "test-cert-profile"
    ddns_ip = ["192.0.2.1", "192.0.2.2"]
    ddns_ipv6 = ["2001:db8::1", "2001:db8::2"]
    ddns_vendor_config = [{
      name = "key1"
      value = "value1"
    }]
  }

  dhcp_client = {
    enable = true
  }

  sdwan_link_settings = {
    enable = true
    sdwan_interface_profile = "test-profile"
    upstream_nat = {
      enable = true
      static_ip = {
        fqdn = "example.com"
      }
    }
  }

  ip = null
  ipv6 = null

  interface_management_profile = panos_interface_management_profile.profile.name
}
`

func TestAccAggregateLayer3Subinterface_SdwanLinkSettings_UpstreamNat_StaticIp_Fqdn(t *testing.T) {
	t.Parallel()
	t.Skip("missing required resource: sdwan profile")

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: aggregateLayer3Subinterface_SdwanLinkSettings_UpstreamNat_StaticIp_Fqdn,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_aggregate_layer3_subinterface.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact("ae1.1"),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_layer3_subinterface.example",
						tfjsonpath.New("tag"),
						knownvalue.Int64Exact(1),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_layer3_subinterface.example",
						tfjsonpath.New("comment"),
						knownvalue.StringExact("Test aggregate layer3 subinterface with SD-WAN link settings and upstream NAT static IP FQDN"),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_layer3_subinterface.example",
						tfjsonpath.New("sdwan_link_settings"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"enable":                  knownvalue.Bool(true),
							"sdwan_interface_profile": knownvalue.StringExact("test-profile"),
							"upstream_nat": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"enable": knownvalue.Bool(true),
								"static_ip": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"fqdn": knownvalue.StringExact("example.com"),
								}),
							}),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_layer3_subinterface.example",
						tfjsonpath.New("ip"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_layer3_subinterface.example",
						tfjsonpath.New("dhcp_client").AtMapKey("enable"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_layer3_subinterface.example",
						tfjsonpath.New("ipv6"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_layer3_subinterface.example",
						tfjsonpath.New("interface_management_profile"),
						knownvalue.StringExact(prefix),
					),
				},
			},
		},
	})
}

const aggregateLayer3Subinterface_SdwanLinkSettings_UpstreamNat_StaticIp_IpAddress = `
variable "prefix" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = "${var.prefix}-tmpl"
}

resource "panos_aggregate_interface" "parent" {
  location = {
    template = {
      vsys = "vsys1"
      name = panos_template.example.name
    }
  }

  name = "ae1"
  layer3 = {}
}

resource "panos_interface_management_profile" "profile" {
  location = { template = { name = panos_template.example.name } }

  name = var.prefix
}

resource "panos_aggregate_layer3_subinterface" "example" {
  location = { template = { name = panos_template.example.name, vsys = "vsys1" } }

  parent = panos_aggregate_interface.parent.name
  name = "ae1.1"
  tag = 1

  comment = "Test aggregate layer3 subinterface with SD-WAN link settings and upstream NAT static IP address"

  dhcp_client = {
    enable = true
  }

  sdwan_link_settings = {
    enable = true
    sdwan_interface_profile = "test-profile"
    upstream_nat = {
      enable = true
      static_ip = {
        ip_address = "203.0.113.1"
      }
    }
  }

  ip = null
  ipv6 = null

  interface_management_profile = panos_interface_management_profile.profile.name
}
`

func TestAccAggregateLayer3Subinterface_SdwanLinkSettings_UpstreamNat_StaticIp_Address(t *testing.T) {
	t.Parallel()
	t.Skip("missing required resource: sdwan profile")

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: aggregateLayer3Subinterface_SdwanLinkSettings_UpstreamNat_StaticIp_IpAddress,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_aggregate_layer3_subinterface.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact("ae1.1"),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_layer3_subinterface.example",
						tfjsonpath.New("tag"),
						knownvalue.Int64Exact(1),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_layer3_subinterface.example",
						tfjsonpath.New("comment"),
						knownvalue.StringExact("Test aggregate layer3 subinterface with SD-WAN link settings and upstream NAT static IP address"),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_layer3_subinterface.example",
						tfjsonpath.New("sdwan_link_settings"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"enable":                  knownvalue.Bool(true),
							"sdwan_interface_profile": knownvalue.StringExact("test-profile"),
							"upstream_nat": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"enable": knownvalue.Bool(true),
								"static_ip": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"ip_address": knownvalue.StringExact("203.0.113.1"),
								}),
							}),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_layer3_subinterface.example",
						tfjsonpath.New("ip"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_layer3_subinterface.example",
						tfjsonpath.New("dhcp_client").AtMapKey("enable"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_layer3_subinterface.example",
						tfjsonpath.New("ipv6"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_layer3_subinterface.example",
						tfjsonpath.New("interface_management_profile"),
						knownvalue.StringExact(prefix),
					),
				},
			},
		},
	})
}
