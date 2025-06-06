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

func TestAccAggregateInterface_DecryptMirror(t *testing.T) {
	t.Parallel()

	interfaceName := "ae1"
	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: aggregateInterfaceDecryptMirror1,
				ConfigVariables: map[string]config.Variable{
					"prefix":         config.StringVariable(prefix),
					"interface_name": config.StringVariable(interfaceName),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_aggregate_interface.iface",
						tfjsonpath.New("name"),
						knownvalue.StringExact(interfaceName),
					),

					statecheck.ExpectKnownValue(
						"panos_aggregate_interface.iface",
						tfjsonpath.
							New("decrypt_mirror"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{}),
					),
				},
			},
			{
				Config: aggregateInterfaceDecryptMirror1,
				ConfigVariables: map[string]config.Variable{
					"prefix":         config.StringVariable(prefix),
					"interface_name": config.StringVariable(interfaceName),
				},
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

const aggregateInterfaceDecryptMirror1 = `
variable "prefix" { type = string }
variable "interface_name" { type = string }

locals {
  template_name = format("%s-tmpl", var.prefix)
}

resource "panos_template" "template" {
  location = { panorama = {} }
  name = local.template_name
}

resource "panos_aggregate_interface" "iface" {
  location = { template = { name = local.template_name } }

  name = var.interface_name

  comment = "aggregate interface comment"
  decrypt_mirror = {}
}
`

func TestAccAggregateInterface_HA(t *testing.T) {
	t.Parallel()

	interfaceName := "ae1"
	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: aggregateInterfaceHa1,
				ConfigVariables: map[string]config.Variable{
					"prefix":         config.StringVariable(prefix),
					"interface_name": config.StringVariable(interfaceName),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_aggregate_interface.iface",
						tfjsonpath.New("name"),
						knownvalue.StringExact(interfaceName),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_interface.iface",
						tfjsonpath.
							New("ha").
							AtMapKey("lacp").
							AtMapKey("enable"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_interface.iface",
						tfjsonpath.
							New("ha").
							AtMapKey("lacp").
							AtMapKey("fast_failover"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_interface.iface",
						tfjsonpath.
							New("ha").
							AtMapKey("lacp").
							AtMapKey("max_ports"),
						knownvalue.Int64Exact(4),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_interface.iface",
						tfjsonpath.
							New("ha").
							AtMapKey("lacp").
							AtMapKey("mode"),
						knownvalue.StringExact("active"),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_interface.iface",
						tfjsonpath.
							New("ha").
							AtMapKey("lacp").
							AtMapKey("system_priority"),
						knownvalue.Int64Exact(10),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_interface.iface",
						tfjsonpath.
							New("ha").
							AtMapKey("lacp").
							AtMapKey("transmission_rate"),
						knownvalue.StringExact("fast"),
					),
				},
			},
			{
				Config: aggregateInterfaceHa1,
				ConfigVariables: map[string]config.Variable{
					"prefix":         config.StringVariable(prefix),
					"interface_name": config.StringVariable(interfaceName),
				},
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

const aggregateInterfaceHa1 = `
variable "prefix" { type = string }
variable "interface_name" { type = string }

locals {
  template_name = format("%s-tmpl", var.prefix)
}

resource "panos_template" "template" {
  location = { panorama = {} }
  name = local.template_name
}

resource "panos_aggregate_interface" "iface" {
  location = { template = { name = local.template_name } }

  name = var.interface_name

  comment = "aggregate interface comment"
  ha = {
    lacp = {
      enable = true
      fast_failover = true
      max_ports = 4
      mode = "active"
      system_priority = 10
      transmission_rate = "fast"
    }
  }
}
`

func TestAccAggregateInterface_Layer2(t *testing.T) {
	t.Parallel()

	interfaceName := "ae1"
	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: aggregateInterfaceLayer21,
				ConfigVariables: map[string]config.Variable{
					"prefix":         config.StringVariable(prefix),
					"interface_name": config.StringVariable(interfaceName),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_aggregate_interface.iface",
						tfjsonpath.New("name"),
						knownvalue.StringExact(interfaceName),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_interface.iface",
						tfjsonpath.
							New("layer2"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							//"netflow_profile": knownvalue.StringExact("netflow-profile"),
							"netflow_profile": knownvalue.Null(),
							"lacp": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"enable":            knownvalue.Bool(true),
								"fast_failover":     knownvalue.Bool(true),
								"max_ports":         knownvalue.Int64Exact(2),
								"mode":              knownvalue.StringExact("active"),
								"system_priority":   knownvalue.Int64Exact(10),
								"transmission_rate": knownvalue.StringExact("fast"),
								"high_availability": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"passive_pre_negotiation": knownvalue.Bool(true),
								}),
							}),
							"lldp": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"enable":  knownvalue.Bool(true),
								"profile": knownvalue.Null(),
								"high_availability": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"passive_pre_negotiation": knownvalue.Bool(true),
								}),
							}),
						}),
					),
				},
			},
			{
				Config: aggregateInterfaceLayer21,
				ConfigVariables: map[string]config.Variable{
					"prefix":         config.StringVariable(prefix),
					"interface_name": config.StringVariable(interfaceName),
				},
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

const aggregateInterfaceLayer21 = `
variable "prefix" { type = string }
variable "interface_name" { type = string }

locals {
  template_name = format("%s-tmpl", var.prefix)
}

resource "panos_template" "template" {
  location = { panorama = {} }
  name = local.template_name
}

resource "panos_aggregate_interface" "iface" {
  location = { template = { name = local.template_name } }

  name = var.interface_name

  layer2 = {
    # netflow_profile = "netflow-profile",
    lacp = {
      enable = true
      fast_failover = true
      max_ports = 2
      mode = "active"
      system_priority = 10
      transmission_rate = "fast"
      high_availability = {
        passive_pre_negotiation = true
      }
    }
    lldp = {
      enable = true
      #profile = format("%s-lldp-profile", var.profile)
      high_availability = {
        passive_pre_negotiation = true
      }
    }
  }
}
`

func TestAccAggregateInterface_Layer3_1(t *testing.T) {
	t.Parallel()

	interfaceName := "ae1"
	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: aggregateInterfaceLayer31,
				ConfigVariables: map[string]config.Variable{
					"prefix":         config.StringVariable(prefix),
					"interface_name": config.StringVariable(interfaceName),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_aggregate_interface.iface",
						tfjsonpath.New("name"),
						knownvalue.StringExact(interfaceName),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_interface.iface",
						tfjsonpath.
							New("layer3"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"decrypt_forward":              knownvalue.Bool(true),
							"df_ignore":                    knownvalue.Bool(true),
							"interface_management_profile": knownvalue.StringExact(fmt.Sprintf("%s-profile", prefix)),
							"mtu":                          knownvalue.Int64Exact(9216),
							// FIXME: panos_netflow_profile implementation needed
							"netflow_profile":        knownvalue.Null(),
							"untagged_sub_interface": knownvalue.Bool(true),
							"adjust_tcp_mss": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"enable":              knownvalue.Bool(true),
								"ipv4_mss_adjustment": knownvalue.Int64Exact(40),
								"ipv6_mss_adjustment": knownvalue.Int64Exact(60),
							}),
							"arp": knownvalue.ListExact([]knownvalue.Check{
								knownvalue.ObjectExact(map[string]knownvalue.Check{
									"name":       knownvalue.StringExact("172.16.0.1"),
									"hw_address": knownvalue.StringExact("aa:bb:cc:dd:ee:ff"),
								}),
							}),
							"bonjour": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"enable":    knownvalue.Bool(true),
								"group_id":  knownvalue.Int64Exact(1),
								"ttl_check": knownvalue.Bool(true),
							}),
							"ddns_config": knownvalue.Null(),
							"dhcp_client": knownvalue.Null(),
							"ip": knownvalue.ListExact([]knownvalue.Check{
								knownvalue.ObjectExact(map[string]knownvalue.Check{
									"name":          knownvalue.StringExact("172.16.0.1"),
									"sdwan_gateway": knownvalue.Null(), // FIXME: panos_sdwan_gateway
								}),
							}),
							"ipv6": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"enabled":      knownvalue.Bool(true),
								"interface_id": knownvalue.StringExact("10"),
								"address": knownvalue.ListExact([]knownvalue.Check{
									knownvalue.ObjectExact(map[string]knownvalue.Check{
										"name":                knownvalue.StringExact("::1"),
										"enable_on_interface": knownvalue.Bool(true),
										"advertise": knownvalue.ObjectExact(map[string]knownvalue.Check{
											"auto_config_flag":   knownvalue.Bool(true),
											"enable":             knownvalue.Bool(true),
											"onlink_flag":        knownvalue.Bool(true),
											"preferred_lifetime": knownvalue.StringExact("1200"),
											"valid_lifetime":     knownvalue.StringExact("1200"),
										}),
										"anycast": knownvalue.ObjectExact(map[string]knownvalue.Check{}),
										"prefix":  knownvalue.ObjectExact(map[string]knownvalue.Check{}),
									}),
								}),
								"dhcp_client": knownvalue.Null(),
								"inherited":   knownvalue.Null(),
								"neighbor_discovery": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"dad_attempts":       knownvalue.Int64Exact(10),
									"enable_dad":         knownvalue.Bool(true),
									"enable_ndp_monitor": knownvalue.Bool(true),
									"ns_interval":        knownvalue.Int64Exact(100),
									"reachable_time":     knownvalue.Int64Exact(100),
									"neighbor": knownvalue.ListExact([]knownvalue.Check{
										knownvalue.ObjectExact(map[string]knownvalue.Check{
											"name":       knownvalue.StringExact("::2"),
											"hw_address": knownvalue.StringExact("bb:cc:dd:ee:ff:aa"),
										}),
									}),
									"router_advertisement": knownvalue.ObjectExact(map[string]knownvalue.Check{
										"enable":                   knownvalue.Bool(true),
										"enable_consistency_check": knownvalue.Bool(true),
										"hop_limit":                knownvalue.StringExact("10"),
										"lifetime":                 knownvalue.Int64Exact(200),
										"link_mtu":                 knownvalue.StringExact("9216"),
										"managed_flag":             knownvalue.Bool(false),
										"max_interval":             knownvalue.Int64Exact(100),
										"min_interval":             knownvalue.Int64Exact(40),
										"other_flag":               knownvalue.Bool(false),
										"reachable_time":           knownvalue.StringExact("100"),
										"retransmission_timer":     knownvalue.StringExact("200"),
										"router_preference":        knownvalue.StringExact("High"),
										"dns_support": knownvalue.ObjectExact(map[string]knownvalue.Check{
											"enable": knownvalue.Bool(true),
											"server": knownvalue.ListExact([]knownvalue.Check{
												knownvalue.ObjectExact(map[string]knownvalue.Check{
													"name":     knownvalue.StringExact("::2"),
													"lifetime": knownvalue.Int64Exact(100),
												}),
											}),
											"suffix": knownvalue.ListExact([]knownvalue.Check{
												knownvalue.ObjectExact(map[string]knownvalue.Check{
													"name":     knownvalue.StringExact("example.com"),
													"lifetime": knownvalue.Int64Exact(100),
												}),
											}),
										}),
									}),
								}),
							}),
							"lacp": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"enable":            knownvalue.Bool(true),
								"fast_failover":     knownvalue.Bool(true),
								"max_ports":         knownvalue.Int64Exact(2),
								"mode":              knownvalue.StringExact("active"),
								"system_priority":   knownvalue.Int64Exact(10),
								"transmission_rate": knownvalue.StringExact("fast"),
								"high_availability": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"passive_pre_negotiation": knownvalue.Bool(true),
								}),
							}),
							"lldp": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"enable": knownvalue.Bool(true),
								//"profile": knownvalue.StringExact(fmt.Sprintf("%s-profile", prefix)), // FIXME: missing resource
								"profile": knownvalue.Null(),
								"high_availability": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"passive_pre_negotiation": knownvalue.Bool(true),
								}),
							}),
							"ndp_proxy": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"enabled": knownvalue.Bool(true),
								"address": knownvalue.ListExact([]knownvalue.Check{
									knownvalue.ObjectExact(map[string]knownvalue.Check{
										"name":   knownvalue.StringExact("172.16.0.1"),
										"negate": knownvalue.Bool(true),
									}),
								}),
							}),
							// FIXME: missing resource
							// "sdwan_link_settings": knownvalue.ObjectExact(map[string]knownvalue.Check{
							// 	"enable":                  knownvalue.Bool(true),
							// 	"sdwan_interface_profile": knownvalue.StringExact(fmt.Sprintf("%s-profile", prefix)),
							// 	"upstream_nat": knownvalue.ObjectExact(map[string]knownvalue.Check{
							// 		"enable": knownvalue.Bool(true),
							// 		"static_ip": knownvalue.ObjectExact(map[string]knownvalue.Check{
							// 			"fqdn":       knownvalue.StringExact("example.com"),
							// 			"ip_address": knownvalue.StringExact("172.16.0.1"),
							// 		}),
							// 	}),
							// }),
							"sdwan_link_settings": knownvalue.Null(),
						}),
					),
				},
			},
			{
				Config: aggregateInterfaceLayer31,
				ConfigVariables: map[string]config.Variable{
					"prefix":         config.StringVariable(prefix),
					"interface_name": config.StringVariable(interfaceName),
				},
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccAggregateInterface_Layer3_2(t *testing.T) {
	t.Parallel()

	interfaceName := "ae1"
	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: aggregateInterfaceLayer32,
				ConfigVariables: map[string]config.Variable{
					"prefix":         config.StringVariable(prefix),
					"interface_name": config.StringVariable(interfaceName),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_aggregate_interface.iface",
						tfjsonpath.New("name"),
						knownvalue.StringExact(interfaceName),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_interface.iface",
						tfjsonpath.
							New("layer3"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"decrypt_forward":              knownvalue.Null(),
							"df_ignore":                    knownvalue.Null(),
							"interface_management_profile": knownvalue.Null(),
							"mtu":                          knownvalue.Null(),
							// FIXME: panos_netflow_profile implementation needed
							"netflow_profile":        knownvalue.Null(),
							"untagged_sub_interface": knownvalue.Null(),
							"adjust_tcp_mss":         knownvalue.Null(),
							"arp":                    knownvalue.Null(),
							"bonjour":                knownvalue.Null(),
							"ddns_config":            knownvalue.Null(),
							"dhcp_client": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"create_default_route": knownvalue.Bool(true),
								"default_route_metric": knownvalue.Int64Exact(100),
								"enable":               knownvalue.Bool(true),
								"send_hostname": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"enable":   knownvalue.Bool(true),
									"hostname": knownvalue.StringExact("example.com"),
								}),
							}),
							"ip": knownvalue.Null(),
							"ipv6": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"address": knownvalue.Null(),
								"dhcp_client": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"accept_ra_route":      knownvalue.Bool(true),
									"default_route_metric": knownvalue.Int64Exact(100),
									"enable":               knownvalue.Bool(true),
									"preference":           knownvalue.StringExact("high"),
									"neighbor_discovery": knownvalue.ObjectExact(map[string]knownvalue.Check{
										"dad_attempts":       knownvalue.Int64Exact(10),
										"enable_dad":         knownvalue.Bool(true),
										"enable_ndp_monitor": knownvalue.Bool(true),
										"ns_interval":        knownvalue.Int64Exact(100),
										"reachable_time":     knownvalue.Int64Exact(100),
										"dns_server": knownvalue.ObjectExact(map[string]knownvalue.Check{
											"enable": knownvalue.Bool(true),
											"source": knownvalue.ObjectExact(map[string]knownvalue.Check{
												"dhcpv6": knownvalue.Null(),
												"manual": knownvalue.ObjectExact(map[string]knownvalue.Check{
													"server": knownvalue.ListExact([]knownvalue.Check{
														knownvalue.ObjectExact(map[string]knownvalue.Check{
															"name":     knownvalue.StringExact("::1"),
															"lifetime": knownvalue.Int64Exact(4),
														}),
													}),
												}),
											}),
										}),
										"dns_suffix": knownvalue.ObjectExact(map[string]knownvalue.Check{
											"enable": knownvalue.Bool(true),
											"source": knownvalue.ObjectExact(map[string]knownvalue.Check{
												"dhcpv6": knownvalue.Null(),
												"manual": knownvalue.ObjectExact(map[string]knownvalue.Check{
													"suffix": knownvalue.ListExact([]knownvalue.Check{
														knownvalue.ObjectExact(map[string]knownvalue.Check{
															"name":     knownvalue.StringExact("example.com"),
															"lifetime": knownvalue.Int64Exact(4),
														}),
													}),
												}),
											}),
										}),
										"neighbor": knownvalue.ListExact([]knownvalue.Check{
											knownvalue.ObjectExact(map[string]knownvalue.Check{
												"name":       knownvalue.StringExact("::2"),
												"hw_address": knownvalue.StringExact("aa:bb:cc:dd:ee:ff"),
											}),
										}),
									}),
									"prefix_delegation": knownvalue.ObjectExact(map[string]knownvalue.Check{
										"enable": knownvalue.ObjectExact(map[string]knownvalue.Check{
											"no": knownvalue.Null(),
											"yes": knownvalue.ObjectExact(map[string]knownvalue.Check{
												"pfx_pool_name":   knownvalue.StringExact("pfx-pool"),
												"prefix_len":      knownvalue.Int64Exact(8),
												"prefix_len_hint": knownvalue.Bool(true),
											}),
										}),
									}),
									"v6_options": knownvalue.ObjectExact(map[string]knownvalue.Check{
										"duid_type":             knownvalue.StringExact("duid-type-llt"),
										"rapid_commit":          knownvalue.Bool(true),
										"support_srvr_reconfig": knownvalue.Bool(true),
										"enable": knownvalue.ObjectExact(map[string]knownvalue.Check{
											"no": knownvalue.Null(),
											"yes": knownvalue.ObjectExact(map[string]knownvalue.Check{
												"non_temp_addr": knownvalue.Bool(true),
												"temp_addr":     knownvalue.Bool(true),
											}),
										}),
									}),
								}),
								"enabled":            knownvalue.Null(),
								"inherited":          knownvalue.Null(),
								"interface_id":       knownvalue.StringExact("EUI-64"),
								"neighbor_discovery": knownvalue.Null(),
							}),
							"lacp":                knownvalue.Null(),
							"lldp":                knownvalue.Null(),
							"ndp_proxy":           knownvalue.Null(),
							"sdwan_link_settings": knownvalue.Null(),
						}),
					),
				},
			},
			{
				Config: aggregateInterfaceLayer32,
				ConfigVariables: map[string]config.Variable{
					"prefix":         config.StringVariable(prefix),
					"interface_name": config.StringVariable(interfaceName),
				},
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccAggregateInterface_Layer3_3(t *testing.T) {
	t.Parallel()

	interfaceName := "ae1"
	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: aggregateInterfaceLayer33,
				ConfigVariables: map[string]config.Variable{
					"prefix":         config.StringVariable(prefix),
					"interface_name": config.StringVariable(interfaceName),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_aggregate_interface.iface",
						tfjsonpath.New("name"),
						knownvalue.StringExact(interfaceName),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_interface.iface",
						tfjsonpath.
							New("layer3"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"decrypt_forward":              knownvalue.Null(),
							"df_ignore":                    knownvalue.Null(),
							"interface_management_profile": knownvalue.Null(),
							"mtu":                          knownvalue.Null(),
							// FIXME: panos_netflow_profile implementation needed
							"netflow_profile":        knownvalue.Null(),
							"untagged_sub_interface": knownvalue.Null(),
							"adjust_tcp_mss":         knownvalue.Null(),
							"arp":                    knownvalue.Null(),
							"bonjour":                knownvalue.Null(),
							"ddns_config":            knownvalue.Null(),
							"dhcp_client":            knownvalue.Null(),
							"ip":                     knownvalue.Null(),
							"ipv6": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"address":     knownvalue.Null(),
								"dhcp_client": knownvalue.Null(),
								"enabled":     knownvalue.Null(),
								"inherited": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"enable": knownvalue.Bool(true),
									"assign_addr": knownvalue.ListExact([]knownvalue.Check{
										knownvalue.ObjectExact(map[string]knownvalue.Check{
											"name": knownvalue.StringExact("172.16.0.1"),
											"type": knownvalue.ObjectExact(map[string]knownvalue.Check{
												"ula": knownvalue.Null(),
												"gua": knownvalue.ObjectExact(map[string]knownvalue.Check{
													"enable_on_interface": knownvalue.Bool(true),
													// "prefix_pool": knownvalue.StringExact(fmt.Sprintf("%s-pool", prefix)) FIXME: missing resource
													"prefix_pool": knownvalue.Null(),
													"advertise": knownvalue.ObjectExact(map[string]knownvalue.Check{
														"auto_config_flag": knownvalue.Bool(true),
														"enable":           knownvalue.Bool(true),
														"onlink_flag":      knownvalue.Bool(true),
													}),
													"pool_type": knownvalue.ObjectExact(map[string]knownvalue.Check{
														"dynamic":    knownvalue.ObjectExact(nil),
														"dynamic_id": knownvalue.Null(),
													}),
												}),
											}),
										}),
										knownvalue.ObjectExact(map[string]knownvalue.Check{
											"name": knownvalue.StringExact("172.16.0.2"),
											"type": knownvalue.ObjectExact(map[string]knownvalue.Check{
												"ula": knownvalue.Null(),
												"gua": knownvalue.ObjectExact(map[string]knownvalue.Check{
													"enable_on_interface": knownvalue.Null(),
													"prefix_pool":         knownvalue.Null(),
													"advertise":           knownvalue.Null(),
													"pool_type": knownvalue.ObjectExact(map[string]knownvalue.Check{
														"dynamic": knownvalue.Null(),
														"dynamic_id": knownvalue.ObjectExact(map[string]knownvalue.Check{
															"identifier": knownvalue.Int64Exact(100),
														}),
													}),
												}),
											}),
										}),
										knownvalue.ObjectExact(map[string]knownvalue.Check{
											"name": knownvalue.StringExact("172.16.0.3"),
											"type": knownvalue.ObjectExact(map[string]knownvalue.Check{
												"ula": knownvalue.ObjectExact(map[string]knownvalue.Check{
													"address":             knownvalue.StringExact("fd12:3456:789a:1::1"),
													"anycast":             knownvalue.Bool(true),
													"enable_on_interface": knownvalue.Bool(true),
													"prefix":              knownvalue.Bool(true),
													"advertise": knownvalue.ObjectExact(map[string]knownvalue.Check{
														"auto_config_flag":   knownvalue.Bool(true),
														"enable":             knownvalue.Bool(true),
														"onlink_flag":        knownvalue.Bool(true),
														"preferred_lifetime": knownvalue.StringExact("100"),
														"valid_lifetime":     knownvalue.StringExact("200"),
													}),
												}),
												"gua": knownvalue.Null(),
											}),
										}),
									}),
									"neighbor_discovery": knownvalue.ObjectExact(map[string]knownvalue.Check{
										"dad_attempts":       knownvalue.Int64Exact(10),
										"enable_dad":         knownvalue.Bool(true),
										"enable_ndp_monitor": knownvalue.Bool(true),
										"ns_interval":        knownvalue.Int64Exact(100),
										"reachable_time":     knownvalue.Int64Exact(100),
										"dns_server": knownvalue.ObjectExact(map[string]knownvalue.Check{
											"enable": knownvalue.Bool(true),
											"source": knownvalue.ObjectExact(map[string]knownvalue.Check{
												"dhcpv6": knownvalue.Null(),
												"manual": knownvalue.ObjectExact(map[string]knownvalue.Check{
													"server": knownvalue.ListExact([]knownvalue.Check{
														knownvalue.ObjectExact(map[string]knownvalue.Check{
															"name":     knownvalue.StringExact("::2"),
															"lifetime": knownvalue.Int64Exact(4),
														}),
													}),
												}),
											}),
										}),
										"dns_suffix": knownvalue.ObjectExact(map[string]knownvalue.Check{
											"enable": knownvalue.Bool(true),
											"source": knownvalue.ObjectExact(map[string]knownvalue.Check{
												"dhcpv6": knownvalue.Null(),
												"manual": knownvalue.ObjectExact(map[string]knownvalue.Check{
													"suffix": knownvalue.ListExact([]knownvalue.Check{
														knownvalue.ObjectExact(map[string]knownvalue.Check{
															"name":     knownvalue.StringExact("example.com"),
															"lifetime": knownvalue.Int64Exact(4),
														}),
													}),
												}),
											}),
										}),
										"neighbor": knownvalue.ListExact([]knownvalue.Check{
											knownvalue.ObjectExact(map[string]knownvalue.Check{
												"name":       knownvalue.StringExact("::2"),
												"hw_address": knownvalue.StringExact("aa:bb:cc:dd:ee:ff"),
											}),
										}),
										"router_advertisement": knownvalue.ObjectExact(map[string]knownvalue.Check{
											"enable":                   knownvalue.Bool(true),
											"enable_consistency_check": knownvalue.Bool(true),
											"hop_limit":                knownvalue.StringExact("10"),
											"lifetime":                 knownvalue.Int64Exact(200),
											"link_mtu":                 knownvalue.StringExact("9216"),
											"managed_flag":             knownvalue.Bool(true),
											"max_interval":             knownvalue.Int64Exact(100),
											"min_interval":             knownvalue.Int64Exact(40),
											"other_flag":               knownvalue.Bool(true),
											"reachable_time":           knownvalue.StringExact("100"),
											"retransmission_timer":     knownvalue.StringExact("1000"),
											"router_preference":        knownvalue.StringExact("High"),
										}),
									}),
								}),
								"interface_id":       knownvalue.StringExact("EUI-64"),
								"neighbor_discovery": knownvalue.Null(),
							}),
							"lacp":                knownvalue.Null(),
							"lldp":                knownvalue.Null(),
							"ndp_proxy":           knownvalue.Null(),
							"sdwan_link_settings": knownvalue.Null(),
						}),
					),
				},
			},
			{
				Config: aggregateInterfaceLayer33,
				ConfigVariables: map[string]config.Variable{
					"prefix":         config.StringVariable(prefix),
					"interface_name": config.StringVariable(interfaceName),
				},
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

const aggregateInterfaceLayer31 = `
variable "prefix" { type = string }
variable "interface_name" { type = string }

locals {
  template_name = format("%s-tmpl", var.prefix)
}

resource "panos_template" "template" {
  location = { panorama = {} }
  name = local.template_name
}

resource "panos_interface_management_profile" "profile" {
  location = { template = { name = panos_template.template.name } }

  name = format("%s-profile", var.prefix)
}


resource "panos_aggregate_interface" "iface" {
  location = { template = { name = panos_template.template.name } }

  name = var.interface_name

  layer3 = {
    decrypt_forward = true
    df_ignore = true
    interface_management_profile = panos_interface_management_profile.profile.name
    mtu = 9216
    #netflow_profile = "netflow-profile"
    untagged_sub_interface = true
    adjust_tcp_mss = {
      enable = true
      ipv4_mss_adjustment = 40
      ipv6_mss_adjustment = 60
    }
    arp = [{
      name = "172.16.0.1"
      hw_address = "aa:bb:cc:dd:ee:ff"
    }]
    bonjour = {
      enable = true
      group_id = 1
      ttl_check = true
    }
    #ddns_config = {
    #  ddns_cert_profile = format("%s-ddns-cert-profile", var.prefix)
    #  ddns_enabled = true
    #  ddns_hostname = "example.com"
    #  ddns_ip = ["172.16.0.1"]
    #  ddns_ipv6 = ["::1"]
    #  ddns_update_interval = 100
    #  ddns_vendor = "Some Vendor"
    #  ddns_vendor_config = [{
    #    name = "some-config"
    #    value = "some-value"
    #  }]
    #}
    ip = [{
      name = "172.16.0.1",
      #sdwan_gateway = "sdwan-gateway"
    }]
    ipv6 = {
      enabled = true
      interface_id = "10"
      address = [{
        name = "::1"
        enable_on_interface = true
        advertise = {
          auto_config_flag = true
          enable = true
          onlink_flag = true
          preferred_lifetime = "1200"
          valid_lifetime = "1200"
        }
        anycast = {}
        prefix = {}
       }]
       neighbor_discovery = {
         dad_attempts = 10
         enable_dad = true
         enable_ndp_monitor = true
         ns_interval = 100
         reachable_time = 100
         neighbor = [{
           name = "::2"
           hw_address = "bb:cc:dd:ee:ff:aa"
         }]
         router_advertisement = {
           enable = true
           enable_consistency_check = true
           hop_limit = "10"
           lifetime = 200
           link_mtu = "9216"
           managed_flag = false
           max_interval = 100
           min_interval = 40
           other_flag = false
           reachable_time = "100"
           retransmission_timer = "200"
           router_preference = "High"
           dns_support = {
             enable = true
             server = [{ name = "::2", lifetime = 100 }]
             suffix = [{ name = "example.com", lifetime = 100 }]
           }
         }
       }
    }
    lacp = {
      enable = true
      fast_failover = true
      max_ports = 2
      mode = "active"
      system_priority = 10
      transmission_rate = "fast"
      high_availability = {
        passive_pre_negotiation = true
      }
    }
    lldp = {
      enable = true
      #profile = format("%s-profile", var.prefix)
      high_availability = {
        passive_pre_negotiation = true
      }
    }
    ndp_proxy = {
      enabled = true
      address = [{
        name = "172.16.0.1"
        negate = true
      }]
    }
    #sdwan_link_settings = {
    #  enable = true
    #  sdwan_interface_profile = format("%s-profile", var.prefix)
    #  upstream_nat = {
    #    enable = true
    #    static_ip = {
    #      fqdn = "example.com"
    #    }
    #  }
    #}
  }
}
`

const aggregateInterfaceLayer32 = `
variable "prefix" { type = string }
variable "interface_name" { type = string }

locals {
  template_name = format("%s-tmpl", var.prefix)
}

resource "panos_template" "template" {
  location = { panorama = {} }
  name = local.template_name
}

resource "panos_aggregate_interface" "iface" {
  location = { template = { name = panos_template.template.name } }

  name = var.interface_name

  layer3 = {
    dhcp_client = {
      create_default_route = true
      default_route_metric = 100
      enable = true
      send_hostname = {
        enable = true
        hostname = "example.com"
      }
    }
    ipv6 = {
      dhcp_client = {
        accept_ra_route = true
        default_route_metric = 100
        enable = true
        preference = "high"
        neighbor_discovery = {
          dad_attempts = 10
          enable_dad = true
          enable_ndp_monitor = true
          ns_interval = 100
          reachable_time = 100
          dns_server = { enable = true, source = { manual = { server = [{ name = "::1", lifetime = 4 }]}}}
          dns_suffix = { enable = true, source = { manual = { suffix = [{ name = "example.com", lifetime = 4 }]}}}
          neighbor = [{ name = "::2", hw_address = "aa:bb:cc:dd:ee:ff" }]
        }
        prefix_delegation = {
          enable = {
            yes = {
              pfx_pool_name = "pfx-pool"
              prefix_len = 8
              prefix_len_hint = true
            }
          }
        }
        v6_options = {
          duid_type = "duid-type-llt"
          rapid_commit = true
          support_srvr_reconfig = true
          enable = {
            yes = {
              non_temp_addr = true
              temp_addr = true
            }
          }
        }
      }
    }
    #sdwan_link_settings = {
    #  enable = true
    #  sdwan_interface_profile = format("%s-profile", var.prefix)
    #  upstream_nat = {
    #    enable = true
    #    static_ip = {
    #      ip_address = "172.16.0.1"
    #    }
    #  }
    #}
  }
}
`

const aggregateInterfaceLayer33 = `
variable "prefix" { type = string }
variable "interface_name" { type = string }

locals {
  template_name = format("%s-tmpl", var.prefix)
}

resource "panos_template" "template" {
  location = { panorama = {} }
  name = local.template_name
}

resource "panos_aggregate_interface" "iface" {
  location = { template = { name = panos_template.template.name } }

  name = var.interface_name

  layer3 = {
    ipv6 = {
       inherited = {
         enable = true
         assign_addr = [
           {
             name = "172.16.0.1"
             type = {
               gua = {
                 enable_on_interface = true
                 #prefix_pool = "pool"
                 advertise = {
                   auto_config_flag = true
                   enable = true
                   onlink_flag = true
                 }
                 pool_type = {
                   dynamic = {}
                 }
               }
             }
           },
           {
             name = "172.16.0.2"
             type = {
               gua = {
                 #prefix_pool = "pool"
                 pool_type = {
                   dynamic_id = { identifier = 100 }
                 }
               }
             }
           },
           {
             name = "172.16.0.3"
             type = {
               ula = {
                 address = "fd12:3456:789a:1::1"
                 anycast = true
                 enable_on_interface = true
                 prefix = true
                 advertise = {
                   auto_config_flag = true
                   enable = true
                   onlink_flag = true
                   preferred_lifetime = "100"
                   valid_lifetime = "200"
                 }
               }
             }
           },
         ]
         neighbor_discovery = {
           dad_attempts = 10
           enable_dad = true
           enable_ndp_monitor = true
           ns_interval = 100
           reachable_time = 100
           dns_server = { enable = true, source = { manual = { server = [{ name = "::2", lifetime = 4 }]}}}
           dns_suffix = { enable = true, source = { manual = { suffix = [{ name = "example.com", lifetime = 4 }]}}}
           neighbor = [{ name = "::2", hw_address = "aa:bb:cc:dd:ee:ff" }]
           router_advertisement = {
             enable = true
             enable_consistency_check = true
             hop_limit = "10"
             lifetime = 200
             link_mtu = "9216"
             managed_flag = true
             max_interval = 100
             min_interval = 40
             other_flag = true
             reachable_time = "100"
             retransmission_timer = "1000"
             router_preference = "High"
           }
         }
       }
    }
  }
}
`

func TestAccAggregateInterface_VirtualWire(t *testing.T) {
	t.Parallel()

	interfaceName := "ae1"
	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: aggregateInterfaceVirtualWire1,
				ConfigVariables: map[string]config.Variable{
					"prefix":         config.StringVariable(prefix),
					"interface_name": config.StringVariable(interfaceName),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_aggregate_interface.iface",
						tfjsonpath.New("name"),
						knownvalue.StringExact(interfaceName),
					),
					statecheck.ExpectKnownValue(
						"panos_aggregate_interface.iface",
						tfjsonpath.
							New("virtual_wire"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							// "netflow_profile": knownvalue.StringExact(fmt.Sprintf("%s-profile", prefix)), FIXME: missing resource
							"netflow_profile": knownvalue.Null(),
							"lldp": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"enable": knownvalue.Bool(true),
								// "profile": knownvalue.StringExact(fmt.Sprintf("%s-profile", prefix)), FIXME: missing resource
								"profile": knownvalue.Null(),
								"high_availability": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"passive_pre_negotiation": knownvalue.Bool(true),
								}),
							}),
						}),
					),
				},
			},
		},
	})
}

const aggregateInterfaceVirtualWire1 = `
variable "prefix" { type = string }
variable "interface_name" { type = string }

locals {
  template_name = format("%s-tmpl", var.prefix)
}

resource "panos_template" "template" {
  location = { panorama = {} }
  name = local.template_name
}

resource "panos_aggregate_interface" "iface" {
  location = { template = { name = local.template_name } }

  name = var.interface_name

  virtual_wire = {
    # netflow_profile = "netflow-profile"
    lldp = {
      enable = true
      # profile = format("%s-profile", var.prefix)
      high_availability = {
        passive_pre_negotiation = true
      }
    }
  }
}
`
