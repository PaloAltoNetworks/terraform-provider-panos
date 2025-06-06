package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccEthernetLayer3Subinterface_Basic(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: ethernetLayer3Subinterface_1,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_ethernet_layer3_subinterface.subinterface",
						tfjsonpath.New("name"),
						knownvalue.StringExact("ethernet1/1.1"),
					),
					statecheck.ExpectKnownValue(
						"panos_ethernet_layer3_subinterface.subinterface",
						tfjsonpath.New("tag"),
						knownvalue.Int64Exact(1),
					),
					statecheck.ExpectKnownValue(
						"panos_ethernet_layer3_subinterface.subinterface",
						tfjsonpath.New("parent"),
						knownvalue.StringExact("ethernet1/1"),
					),
					statecheck.ExpectKnownValue(
						"panos_ethernet_layer3_subinterface.subinterface",
						tfjsonpath.New("comment"),
						knownvalue.StringExact("Test subinterface with all top-level parameters"),
					),
					// statecheck.ExpectKnownValue(
					// 	"panos_ethernet_layer3_subinterface.subinterface",
					// 	tfjsonpath.New("netflow_profile"),
					// 	knownvalue.StringExact("NetflowProfile1"),
					// ),
					statecheck.ExpectKnownValue(
						"panos_ethernet_layer3_subinterface.subinterface",
						tfjsonpath.New("mtu"),
						knownvalue.Int64Exact(1500),
					),
					statecheck.ExpectKnownValue(
						"panos_ethernet_layer3_subinterface.subinterface",
						tfjsonpath.New("adjust_tcp_mss"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"enable":              knownvalue.Bool(true),
							"ipv4_mss_adjustment": knownvalue.Int64Exact(100),
							"ipv6_mss_adjustment": knownvalue.Int64Exact(150),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_ethernet_layer3_subinterface.subinterface",
						tfjsonpath.New("arp"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":       knownvalue.StringExact("192.168.0.1"),
								"hw_address": knownvalue.StringExact("00:1a:2b:3c:4d:5e"),
							}),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_ethernet_layer3_subinterface.subinterface",
						tfjsonpath.New("bonjour"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"enable":    knownvalue.Bool(true),
							"group_id":  knownvalue.Int64Exact(5),
							"ttl_check": knownvalue.Bool(true),
						}),
					),
					// statecheck.ExpectKnownValue(
					// 	"panos_ethernet_layer3_subinterface.subinterface",
					// 	tfjsonpath.New("ddns_config"),
					// 	knownvalue.ObjectExact(map[string]knownvalue.Check{
					// 		"ddns_cert_profile":    knownvalue.StringExact("cert-profile-1"),
					// 		"ddns_enabled":         knownvalue.Bool(true),
					// 		"ddns_hostname":        knownvalue.StringExact("test-hostname"),
					// 		"ddns_ip":              knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("10.0.0.1")}),
					// 		"ddns_ipv6":            knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("2001:db8::1")}),
					// 		"ddns_update_interval": knownvalue.Int64Exact(24),
					// 		"ddns_vendor":          knownvalue.StringExact("dyndns"),
					// 		"ddns_vendor_config": knownvalue.ListExact([]knownvalue.Check{
					// 			knownvalue.ObjectExact(map[string]knownvalue.Check{
					// 				"name":  knownvalue.StringExact("vendor_config_1"),
					// 				"value": knownvalue.StringExact("config-value-1"),
					// 			}),
					// 		}),
					// 	}),
					// ),
					statecheck.ExpectKnownValue(
						"panos_ethernet_layer3_subinterface.subinterface",
						tfjsonpath.New("df_ignore"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_ethernet_layer3_subinterface.subinterface",
						tfjsonpath.New("interface_management_profile"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_ethernet_layer3_subinterface.subinterface",
						tfjsonpath.New("ip"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":          knownvalue.StringExact("192.168.1.1"),
								"sdwan_gateway": knownvalue.StringExact("192.168.1.1"),
							}),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_ethernet_layer3_subinterface.subinterface",
						tfjsonpath.New("ndp_proxy"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"enabled": knownvalue.Bool(true),
							"address": knownvalue.ListExact([]knownvalue.Check{
								knownvalue.ObjectExact(map[string]knownvalue.Check{
									"name":   knownvalue.StringExact("ndp_proxy_1"),
									"negate": knownvalue.Bool(false),
								}),
							}),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_ethernet_layer3_subinterface.subinterface",
						tfjsonpath.New("ipv6"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"enabled": knownvalue.Bool(true),
							"address": knownvalue.ListExact([]knownvalue.Check{
								knownvalue.ObjectExact(map[string]knownvalue.Check{
									"name":                knownvalue.StringExact("2001:db8:85a3::8a2e:370:7334"),
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
							"interface_id":       knownvalue.StringExact("EUI-64"),
							"dhcp_client":        knownvalue.Null(),
							"inherited":          knownvalue.Null(),
							"neighbor_discovery": knownvalue.Null(),
						}),
					),
					// statecheck.ExpectKnownValue(
					// 	"panos_ethernet_layer3_subinterface.subinterface",
					// 	tfjsonpath.New("sdwan_link_settings"),
					// 	knownvalue.ObjectExact(map[string]knownvalue.Check{
					// 		"enable":                  knownvalue.Bool(true),
					// 		"sdwan_interface_profile": knownvalue.StringExact("SdwanProfile1"),
					// 		"upstream_nat": knownvalue.ObjectExact(map[string]knownvalue.Check{
					// 			"enable": knownvalue.Bool(true),
					// 			"ddns":   knownvalue.ObjectExact(map[string]knownvalue.Check{}),
					// 		}),
					// 	}),
					// ),
				},
			},
			{
				Config: ethernetLayer3Subinterface_1,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAccEthernetLayer3Subinterface_DHCP_Client(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: ethernetLayer3Subinterface_DHCP_Client,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_ethernet_layer3_subinterface.subinterface",
						tfjsonpath.New("name"),
						knownvalue.StringExact("ethernet1/1.1"),
					),
					statecheck.ExpectKnownValue(
						"panos_ethernet_layer3_subinterface.subinterface",
						tfjsonpath.New("tag"),
						knownvalue.Int64Exact(1),
					),
					statecheck.ExpectKnownValue(
						"panos_ethernet_layer3_subinterface.subinterface",
						tfjsonpath.New("parent"),
						knownvalue.StringExact("ethernet1/1"),
					),
					statecheck.ExpectKnownValue(
						"panos_ethernet_layer3_subinterface.subinterface",
						tfjsonpath.New("comment"),
						knownvalue.StringExact("Test subinterface with DHCP client configuration"),
					),
					statecheck.ExpectKnownValue(
						"panos_ethernet_layer3_subinterface.subinterface",
						tfjsonpath.New("dhcp_client"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"create_default_route": knownvalue.Bool(true),
							"default_route_metric": knownvalue.Int64Exact(10),
							"enable":               knownvalue.Bool(true),
							"send_hostname": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"enable":   knownvalue.Bool(true),
								"hostname": knownvalue.StringExact("dhcp-client-hostname"),
							}),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_ethernet_layer3_subinterface.subinterface",
						tfjsonpath.New("ipv6").AtMapKey("enabled"),
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownValue(
						"panos_ethernet_layer3_subinterface.subinterface",
						tfjsonpath.New("sdwan_link_settings").AtMapKey("enable"),
						knownvalue.Bool(false),
					),
				},
			},
			{
				Config: ethernetLayer3Subinterface_DHCP_Client,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAccEthernetLayer3Subinterface_PPPoE(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: ethernetLayer3Subinterface_PPPoE,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_ethernet_layer3_subinterface.subinterface",
						tfjsonpath.New("pppoe"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"access_concentrator":  knownvalue.StringExact("ac-1"),
							"authentication":       knownvalue.StringExact("auto"),
							"create_default_route": knownvalue.Bool(true),
							"default_route_metric": knownvalue.Int64Exact(10),
							"enable":               knownvalue.Bool(true),
							"passive": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"enable": knownvalue.Bool(true),
							}),
							"password": knownvalue.StringExact("pppoe-password"),
							"service":  knownvalue.StringExact("pppoe-service"),
							"static_address": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"ip": knownvalue.StringExact("192.168.2.1"),
							}),
							"username": knownvalue.StringExact("pppoe-user"),
						}),
					),
					// Additional checks for other relevant attributes
					statecheck.ExpectKnownValue(
						"panos_ethernet_layer3_subinterface.subinterface",
						tfjsonpath.New("ipv6").AtMapKey("enabled"),
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownValue(
						"panos_ethernet_layer3_subinterface.subinterface",
						tfjsonpath.New("sdwan_link_settings").AtMapKey("enable"),
						knownvalue.Bool(false),
					),
				},
			},
			{
				Config: ethernetLayer3Subinterface_PPPoE,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAccEthernetLayer3Subinterface_IPv6_DHCP_Client(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: ethernetLayer3Subinterface_IPv6_DHCP_Client,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_ethernet_layer3_subinterface.subinterface",
						tfjsonpath.New("name"),
						knownvalue.StringExact("ethernet1/1.1"),
					),
					statecheck.ExpectKnownValue(
						"panos_ethernet_layer3_subinterface.subinterface",
						tfjsonpath.New("tag"),
						knownvalue.Int64Exact(1),
					),
					statecheck.ExpectKnownValue(
						"panos_ethernet_layer3_subinterface.subinterface",
						tfjsonpath.New("comment"),
						knownvalue.StringExact("Test subinterface with IPv6 DHCP client configuration"),
					),
					statecheck.ExpectKnownValue(
						"panos_ethernet_layer3_subinterface.subinterface",
						tfjsonpath.New("ipv6"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"enabled": knownvalue.Bool(true),
							"dhcp_client": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"accept_ra_route":      knownvalue.Bool(true),
								"default_route_metric": knownvalue.Int64Exact(10),
								"enable":               knownvalue.Bool(true),
								"neighbor_discovery": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"dad_attempts":       knownvalue.Int64Exact(1),
									"enable_dad":         knownvalue.Bool(true),
									"enable_ndp_monitor": knownvalue.Bool(true),
									"ns_interval":        knownvalue.Int64Exact(1000),
									"reachable_time":     knownvalue.Int64Exact(30000),
									"dns_server":         knownvalue.Null(),
									"dns_suffix":         knownvalue.Null(),
									"neighbor":           knownvalue.Null(),
								}),
								"preference": knownvalue.StringExact("high"),
								"prefix_delegation": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"enable": knownvalue.ObjectExact(map[string]knownvalue.Check{
										"yes": knownvalue.ObjectExact(map[string]knownvalue.Check{
											"pfx_pool_name":   knownvalue.StringExact("prefix-pool-1"),
											"prefix_len":      knownvalue.Int64Exact(64),
											"prefix_len_hint": knownvalue.Bool(true),
										}),
										"no": knownvalue.Null(),
									}),
								}),
								"v6_options": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"duid_type": knownvalue.StringExact("duid-type-llt"),
									"enable": knownvalue.ObjectExact(map[string]knownvalue.Check{
										"yes": knownvalue.ObjectExact(map[string]knownvalue.Check{
											"non_temp_addr": knownvalue.Bool(true),
											"temp_addr":     knownvalue.Bool(false),
										}),
										"no": knownvalue.Null(),
									}),
									"rapid_commit":          knownvalue.Bool(true),
									"support_srvr_reconfig": knownvalue.Bool(true),
								}),
							}),
							"address":            knownvalue.Null(),
							"inherited":          knownvalue.Null(),
							"interface_id":       knownvalue.StringExact("EUI-64"),
							"neighbor_discovery": knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

func TestAccEthernetLayer3Subinterface_IPv6_Neighbor_Discovery(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: ethernetLayer3Subinterface_IPv6_Neighbor_Discovery,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_ethernet_layer3_subinterface.subinterface",
						tfjsonpath.New("ipv6"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"enabled": knownvalue.Bool(true),
							"neighbor_discovery": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"dad_attempts":       knownvalue.Int64Exact(1),
								"enable_dad":         knownvalue.Bool(true),
								"enable_ndp_monitor": knownvalue.Bool(true),
								"ns_interval":        knownvalue.Int64Exact(1000),
								"reachable_time":     knownvalue.Int64Exact(30000),
								"router_advertisement": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"dns_support": knownvalue.ObjectExact(map[string]knownvalue.Check{
										"enable": knownvalue.Bool(true),
										"server": knownvalue.ListExact([]knownvalue.Check{
											knownvalue.ObjectExact(map[string]knownvalue.Check{
												"name":     knownvalue.StringExact("2001:DB8::1/128"),
												"lifetime": knownvalue.Int64Exact(1200),
											}),
										}),
										"suffix": knownvalue.ListExact([]knownvalue.Check{
											knownvalue.ObjectExact(map[string]knownvalue.Check{
												"name":     knownvalue.StringExact("suffix1"),
												"lifetime": knownvalue.Int64Exact(1200),
											}),
										}),
									}),
									"enable":                   knownvalue.Bool(true),
									"enable_consistency_check": knownvalue.Bool(true),
									"hop_limit":                knownvalue.StringExact("64"),
									"lifetime":                 knownvalue.Int64Exact(1800),
									"link_mtu":                 knownvalue.StringExact("1500"),
									"managed_flag":             knownvalue.Bool(false),
									"max_interval":             knownvalue.Int64Exact(600),
									"min_interval":             knownvalue.Int64Exact(200),
									"other_flag":               knownvalue.Bool(false),
									"reachable_time":           knownvalue.StringExact("0"),
									"retransmission_timer":     knownvalue.StringExact("0"),
									"router_preference":        knownvalue.StringExact("Medium"),
								}),
								"neighbor": knownvalue.ListExact([]knownvalue.Check{
									knownvalue.ObjectExact(map[string]knownvalue.Check{
										"name":       knownvalue.StringExact("2001:DB8::1/128"),
										"hw_address": knownvalue.StringExact("00:1a:2b:3c:4d:5e"),
									}),
								}),
							}),
							"interface_id": knownvalue.StringExact("EUI-64"),
							"address":      knownvalue.Null(),
							"dhcp_client":  knownvalue.Null(),
							"inherited":    knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

func TestAccEthernetLayer3Subinterface_IPv6_GUA(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: ethernetLayer3Subinterface_IPv6_GUA,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_ethernet_layer3_subinterface.subinterface",
						tfjsonpath.New("ipv6").AtMapKey("inherited").AtMapKey("assign_addr").AtSliceIndex(0),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"name": knownvalue.StringExact("gua_config"),
							"type": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"gua": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"enable_on_interface": knownvalue.Bool(true),
									"advertise":           knownvalue.Null(),
									"prefix_pool":         knownvalue.Null(),
									"pool_type":           knownvalue.Null(),
								}),
								"ula": knownvalue.Null(),
							}),
						}),
					),
				},
			},
		},
	})
}

func TestAccEthernetLayer3Subinterface_IPv6_ULA(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: ethernetLayer3Subinterface_IPv6_ULA,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_ethernet_layer3_subinterface.subinterface",
						tfjsonpath.New("ipv6").AtMapKey("inherited").AtMapKey("assign_addr").AtSliceIndex(0),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"name": knownvalue.StringExact("ula_config"),
							"type": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"ula": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"enable_on_interface": knownvalue.Bool(true),
									"address":             knownvalue.StringExact("fd00:1234:5678::/48"),
									"advertise":           knownvalue.Null(),
									"anycast":             knownvalue.Null(),
									"prefix":              knownvalue.Null(),
								}),
								"gua": knownvalue.Null(),
							}),
						}),
					),
				},
			},
		},
	})
}

func TestAccEthernetLayer3Subinterface_SDWAN_DDNS(t *testing.T) {
	t.Parallel()
	t.Skip("Missing resource implementation")

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: ethernetLayer3Subinterface_SDWAN_DDNS,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_ethernet_layer3_subinterface.subinterface",
						tfjsonpath.New("sdwan_link_settings"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"enable": knownvalue.Bool(true),
							"upstream_nat": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"enable":    knownvalue.Bool(true),
								"ddns":      knownvalue.ObjectExact(map[string]knownvalue.Check{}),
								"static_ip": knownvalue.Null(),
							}),
						}),
					),
				},
			},
		},
	})
}

func TestAccEthernetLayer3Subinterface_SDWAN_StaticIP_FQDN(t *testing.T) {
	t.Parallel()
	t.Skip("Missing resource implementation")

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: ethernetLayer3Subinterface_SDWAN_StaticIP_FQDN,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_ethernet_layer3_subinterface.subinterface",
						tfjsonpath.New("sdwan_link_settings").AtMapKey("upstream_nat").AtMapKey("static_ip"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"fqdn":       knownvalue.StringExact("example.com"),
							"ip_address": knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

func TestAccEthernetLayer3Subinterface_SDWAN_StaticIP_IPAddress(t *testing.T) {
	t.Parallel()
	t.Skip("Missing resource implementation")

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: ethernetLayer3Subinterface_SDWAN_StaticIP_IPAddress,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_ethernet_layer3_subinterface.subinterface",
						tfjsonpath.New("sdwan_link_settings").AtMapKey("upstream_nat").AtMapKey("static_ip"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"ip_address": knownvalue.StringExact("203.0.113.1"),
							"fqdn":       knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

const ethernetLayer3Subinterface_1 = `
variable "prefix" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = "${var.prefix}-tmpl"
}

resource "panos_ethernet_interface" "parent" {
  location = {
    template = {
      vsys = "vsys1"
      name = panos_template.example.name
    }
  }

  name = "ethernet1/1"
  layer3 = {}
}

resource "panos_interface_management_profile" "profile" {
  location = { template = { name = panos_template.example.name } }

  name = var.prefix
}

resource "panos_ethernet_layer3_subinterface" "subinterface" {
  location = {
    template = {
      vsys = "vsys1"
      name = panos_template.example.name
    }
  }

  parent = panos_ethernet_interface.parent.name
  name = "ethernet1/1.1"
  tag = 1
  comment = "Test subinterface with all top-level parameters"
  #netflow_profile = "NetflowProfile1"
  mtu = 1500

  adjust_tcp_mss = {
    enable = true
    ipv4_mss_adjustment = 100
    ipv6_mss_adjustment = 150
  }

  arp = [
    {
      name = "192.168.0.1"
      hw_address = "00:1a:2b:3c:4d:5e"
    }
  ]

  bonjour = {
    enable = true
    group_id = 5
    ttl_check = true
  }

  #ddns_config = {
  #  ddns_cert_profile = "cert-profile-1"
  #  ddns_enabled = true
  #  ddns_hostname = "test-hostname"
  #  ddns_ip = ["10.0.0.1"]
  #  ddns_ipv6 = ["2001:db8::1"]
  #  ddns_update_interval = 24
  #  ddns_vendor = "dyndns"
  #  ddns_vendor_config = [
  #    {
  #      name = "vendor_config_1"
  #      value = "config-value-1"
  #    }
  #  ]
  #}

  #decrypt_forward = true

  df_ignore = true

  interface_management_profile = panos_interface_management_profile.profile.name

  ip = [
    {
      name = "192.168.1.1"
      sdwan_gateway = "192.168.1.1"
    }
  ]

  ndp_proxy = {
    enabled = true
    address = [
      {
        name = "ndp_proxy_1"
        negate = false
      }
    ]
  }

  ipv6 = {
    enabled = true
    address = [
      {
        name = "2001:db8:85a3::8a2e:370:7334"
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
      }
    ]

    interface_id = "EUI-64"
  }

  #sdwan_link_settings = {
  #  enable = true
  #  sdwan_interface_profile = "SdwanProfile1"
  #  upstream_nat = {
  #    enable = true
  #    ddns = {}
  #  }
  #}
}
`

const ethernetLayer3Subinterface_PPPoE = `
variable "prefix" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = "${var.prefix}-tmpl"
}

resource "panos_ethernet_interface" "parent" {
  location = {
    template = {
      vsys = "vsys1"
      name = panos_template.example.name
    }
  }

  name = "ethernet1/1"
  layer3 = {}
}

resource "panos_ethernet_layer3_subinterface" "subinterface" {
  location = {
    template = {
      vsys = "vsys1"
      name = panos_template.example.name
    }
  }

  parent = panos_ethernet_interface.parent.name
  name = "ethernet1/1.1"
  tag = 1
  comment = "Test subinterface with PPPoE configuration"
  mtu = 1500

  pppoe = {
    access_concentrator = "ac-1"
    authentication = "auto"
    create_default_route = true
    default_route_metric = 10
    enable = true
    passive = {
      enable = true
    }
    password = "pppoe-password"
    service = "pppoe-service"
    static_address = {
      ip = "192.168.2.1"
    }
    username = "pppoe-user"
  }

  // Disable other configurations to focus on PPPoE
  ipv6 = {
    enabled = false
  }

  sdwan_link_settings = {
    enable = false
  }
}
`

const ethernetLayer3Subinterface_DHCP_Client = `
variable "prefix" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = "${var.prefix}-tmpl"
}

resource "panos_ethernet_interface" "parent" {
  location = {
    template = {
      vsys = "vsys1"
      name = panos_template.example.name
    }
  }

  name = "ethernet1/1"
  layer3 = {}
}

resource "panos_ethernet_layer3_subinterface" "subinterface" {
  location = {
    template = {
      vsys = "vsys1"
      name = panos_template.example.name
    }
  }

  parent = panos_ethernet_interface.parent.name
  name = "ethernet1/1.1"
  tag = 1
  comment = "Test subinterface with DHCP client configuration"

  dhcp_client = {
    create_default_route = true
    default_route_metric = 10
    enable = true
    send_hostname = {
      enable = true
      hostname = "dhcp-client-hostname"
    }
  }

  // Explicitly disable IPv6 and SDWAN to focus on DHCP client
  ipv6 = {
    enabled = false
  }

  sdwan_link_settings = {
    enable = false
  }
}
`

const ethernetLayer3Subinterface_IPv6_GUA = `
variable "prefix" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = "${var.prefix}-tmpl"
}

resource "panos_ethernet_interface" "parent" {
  location = {
    template = {
      vsys = "vsys1"
      name = panos_template.example.name
    }
  }

  name = "ethernet1/1"
  layer3 = {}
}

resource "panos_ethernet_layer3_subinterface" "subinterface" {
  location = {
    template = {
      vsys = "vsys1"
      name = panos_template.example.name
    }
  }

  parent = panos_ethernet_interface.parent.name
  name = "ethernet1/1.1"
  tag = 1

  ipv6 = {
    inherited = {
      assign_addr = [{
	name = "gua_config"
        type = {
          gua = {
            enable_on_interface = true
            #prefix_pool = "my-gua-pool"
          }
        }
      }]
    }
  }
}
`

const ethernetLayer3Subinterface_IPv6_DHCP_Client = `
variable "prefix" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = "${var.prefix}-tmpl"
}

resource "panos_ethernet_interface" "parent" {
  location = {
    template = {
      vsys = "vsys1"
      name = panos_template.example.name
    }
  }

  name = "ethernet1/1"
  layer3 = {}
}

resource "panos_ethernet_layer3_subinterface" "subinterface" {
  location = {
    template = {
      vsys = "vsys1"
      name = panos_template.example.name
    }
  }

  parent = panos_ethernet_interface.parent.name
  name = "ethernet1/1.1"
  tag = 1
  comment = "Test subinterface with IPv6 DHCP client configuration"

  ipv6 = {
    enabled = true
    dhcp_client = {
      accept_ra_route = true
      default_route_metric = 10
      enable = true
      neighbor_discovery = {
        dad_attempts = 1
        enable_dad = true
        enable_ndp_monitor = true
        ns_interval = 1000
        reachable_time = 30000
      }
      preference = "high"
      prefix_delegation = {
        enable = {
          yes = {
            pfx_pool_name = "prefix-pool-1"
            prefix_len = 64
            prefix_len_hint = true
          }
        }
      }
      v6_options = {
        duid_type = "duid-type-llt"
        enable = {
          yes = {
            non_temp_addr = true
            temp_addr = false
          }
        }
        rapid_commit = true
        support_srvr_reconfig = true
      }
    }
  }
}
`

const ethernetLayer3Subinterface_IPv6_Neighbor_Discovery = `
variable "prefix" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = "${var.prefix}-tmpl"
}

resource "panos_ethernet_interface" "parent" {
  location = {
    template = {
      vsys = "vsys1"
      name = panos_template.example.name
    }
  }

  name = "ethernet1/1"
  layer3 = {}
}

resource "panos_ethernet_layer3_subinterface" "subinterface" {
  location = {
    template = {
      vsys = "vsys1"
      name = panos_template.example.name
    }
  }

  parent = panos_ethernet_interface.parent.name
  name = "ethernet1/1.1"
  tag = 1
  comment = "Test subinterface with IPv6 Neighbor Discovery configuration"
  mtu = 1500

  ipv6 = {
    enabled = true
    neighbor_discovery = {
      dad_attempts = 1
      enable_dad = true
      enable_ndp_monitor = true
      neighbor = [{
        name = "2001:DB8::1/128"
        hw_address = "00:1a:2b:3c:4d:5e"
      }]
      ns_interval = 1000
      reachable_time = 30000
      router_advertisement = {
        dns_support = {
          enable = true
          server = [
            {
              name = "2001:DB8::1/128"
              lifetime = 1200
            }
          ]
          suffix = [
            {
              name = "suffix1"
              lifetime = 1200
            }
          ]
        }
        enable = true
        enable_consistency_check = true
        hop_limit = "64"
        lifetime = 1800
        link_mtu = "1500"
        managed_flag = false
        max_interval = 600
        min_interval = 200
        other_flag = false
        reachable_time = "0"
        retransmission_timer = "0"
        router_preference = "Medium"
      }
    }
  }
}
`

const ethernetLayer3Subinterface_IPv6_ULA = `
variable "prefix" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = "${var.prefix}-tmpl"
}

resource "panos_ethernet_interface" "parent" {
  location = {
    template = {
      vsys = "vsys1"
      name = panos_template.example.name
    }
  }

  name = "ethernet1/1"
  layer3 = {}
}

resource "panos_ethernet_layer3_subinterface" "subinterface" {
  location = {
    template = {
      vsys = "vsys1"
      name = panos_template.example.name
    }
  }

  parent = panos_ethernet_interface.parent.name
  name = "ethernet1/1.1"
  tag = 1

  ipv6 = {
    inherited = {
      assign_addr = [
        {
          name = "ula_config"
          type = {
            ula = {
              enable_on_interface = true
              address = "fd00:1234:5678::/48"
            }
          }
        }
      ]
    }
  }
}
`

const ethernetLayer3Subinterface_SDWAN_DDNS = `
variable "prefix" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = "${var.prefix}-tmpl"
}

resource "panos_ethernet_interface" "parent" {
  location = {
    template = {
      vsys = "vsys1"
      name = panos_template.example.name
    }
  }

  name = "ethernet1/1"
  layer3 = {}
}

resource "panos_ethernet_layer3_subinterface" "subinterface" {
  location = {
    template = {
      vsys = "vsys1"
      name = panos_template.example.name
    }
  }

  parent = panos_ethernet_interface.parent.name
  name = "ethernet1/1.1"
  tag = 1

  sdwan_link_settings = {
    enable = true
    upstream_nat = {
      enable = true
      ddns = {}
    }
  }
}
`

const ethernetLayer3Subinterface_SDWAN_StaticIP_FQDN = `
variable "prefix" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = "${var.prefix}-tmpl"
}

resource "panos_ethernet_interface" "parent" {
  location = {
    template = {
      vsys = "vsys1"
      name = panos_template.example.name
    }
  }

  name = "ethernet1/1"
  layer3 = {}
}

resource "panos_ethernet_layer3_subinterface" "subinterface" {
  location = {
    template = {
      vsys = "vsys1"
      name = panos_template.example.name
    }
  }

  parent = panos_ethernet_interface.parent.name
  name = "ethernet1/1.1"
  tag = 1

  sdwan_link_settings = {
    enable = true
    upstream_nat = {
      enable = true
      static_ip = {
        fqdn = "example.com"
      }
    }
  }
}
`

const ethernetLayer3Subinterface_SDWAN_StaticIP_IPAddress = `
variable "prefix" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = "${var.prefix}-tmpl"
}

resource "panos_ethernet_interface" "parent" {
  location = {
    template = {
      vsys = "vsys1"
      name = panos_template.example.name
    }
  }

  name = "ethernet1/1"
  layer3 = {}
}

resource "panos_ethernet_layer3_subinterface" "subinterface" {
  location = {
    template = {
      vsys = "vsys1"
      name = panos_template.example.name
    }
  }

  parent = panos_ethernet_interface.parent.name
  name = "ethernet1/1.1"
  tag = 1

  sdwan_link_settings = {
    enable = true
    upstream_nat = {
      enable = true
      static_ip = {
        ip_address = "203.0.113.1"
      }
    }
  }
}
`
