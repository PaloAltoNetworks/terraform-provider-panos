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

func TestAccBgpFilteringRoutingProfile_Ipv4_Multicast_Basic(t *testing.T) {
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
				Config: bgpFilteringRoutingProfile_Ipv4_Multicast_Basic_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_bgp_filtering_routing_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_bgp_filtering_routing_profile.example",
						tfjsonpath.New("description"),
						knownvalue.StringExact("Test BGP filtering profile"),
					),
					statecheck.ExpectKnownValue(
						"panos_bgp_filtering_routing_profile.example",
						tfjsonpath.New("ipv4"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"multicast": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"inherit":                     knownvalue.Bool(true),
								"conditional_advertisement":   knownvalue.Null(),
								"filter_list":                 knownvalue.Null(),
								"inbound_network_filters":     knownvalue.Null(),
								"outbound_network_filters":    knownvalue.Null(),
								"route_maps":                  knownvalue.Null(),
								"unsuppress_map":              knownvalue.Null(),
							}),
							"unicast": knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

func TestAccBgpFilteringRoutingProfile_Ipv4_Multicast_ConditionalAdvertisement_Exist(t *testing.T) {
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
				Config: bgpFilteringRoutingProfile_Ipv4_Multicast_ConditionalAdvertisement_Exist_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_bgp_filtering_routing_profile.example",
						tfjsonpath.New("ipv4").AtMapKey("multicast").AtMapKey("conditional_advertisement"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"exist": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"advertise_map": knownvalue.StringExact(prefix + "-advertise-map"),
								"exist_map":     knownvalue.StringExact(prefix + "-exist-map"),
							}),
							"non_exist": knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

func TestAccBgpFilteringRoutingProfile_Ipv4_Multicast_ConditionalAdvertisement_NonExist(t *testing.T) {
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
				Config: bgpFilteringRoutingProfile_Ipv4_Multicast_ConditionalAdvertisement_NonExist_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_bgp_filtering_routing_profile.example",
						tfjsonpath.New("ipv4").AtMapKey("multicast").AtMapKey("conditional_advertisement"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"exist": knownvalue.Null(),
							"non_exist": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"advertise_map":  knownvalue.StringExact(prefix + "-advertise-map"),
								"non_exist_map": knownvalue.StringExact(prefix + "-non-exist-map"),
							}),
						}),
					),
				},
			},
		},
	})
}

func TestAccBgpFilteringRoutingProfile_Ipv4_Multicast_FilterList(t *testing.T) {
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
				Config: bgpFilteringRoutingProfile_Ipv4_Multicast_FilterList_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_bgp_filtering_routing_profile.example",
						tfjsonpath.New("ipv4").AtMapKey("multicast").AtMapKey("filter_list"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"inbound":  knownvalue.StringExact(prefix + "-inbound"),
							"outbound": knownvalue.StringExact(prefix + "-outbound"),
						}),
					),
				},
			},
		},
	})
}

func TestAccBgpFilteringRoutingProfile_Ipv4_Multicast_InboundNetworkFilters_DistributeList(t *testing.T) {
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
				Config: bgpFilteringRoutingProfile_Ipv4_Multicast_InboundNetworkFilters_DistributeList_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_bgp_filtering_routing_profile.example",
						tfjsonpath.New("ipv4").AtMapKey("multicast").AtMapKey("inbound_network_filters"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"distribute_list": knownvalue.StringExact(prefix + "-inbound-acl"),
							"prefix_list":     knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

func TestAccBgpFilteringRoutingProfile_Ipv4_Multicast_InboundNetworkFilters_PrefixList(t *testing.T) {
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
				Config: bgpFilteringRoutingProfile_Ipv4_Multicast_InboundNetworkFilters_PrefixList_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_bgp_filtering_routing_profile.example",
						tfjsonpath.New("ipv4").AtMapKey("multicast").AtMapKey("inbound_network_filters"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"distribute_list": knownvalue.Null(),
							"prefix_list":     knownvalue.StringExact(prefix + "-inbound-prefix"),
						}),
					),
				},
			},
		},
	})
}

func TestAccBgpFilteringRoutingProfile_Ipv4_Multicast_OutboundNetworkFilters_DistributeList(t *testing.T) {
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
				Config: bgpFilteringRoutingProfile_Ipv4_Multicast_OutboundNetworkFilters_DistributeList_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_bgp_filtering_routing_profile.example",
						tfjsonpath.New("ipv4").AtMapKey("multicast").AtMapKey("outbound_network_filters"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"distribute_list": knownvalue.StringExact(prefix + "-outbound-acl"),
							"prefix_list":     knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

func TestAccBgpFilteringRoutingProfile_Ipv4_Multicast_OutboundNetworkFilters_PrefixList(t *testing.T) {
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
				Config: bgpFilteringRoutingProfile_Ipv4_Multicast_OutboundNetworkFilters_PrefixList_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_bgp_filtering_routing_profile.example",
						tfjsonpath.New("ipv4").AtMapKey("multicast").AtMapKey("outbound_network_filters"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"distribute_list": knownvalue.Null(),
							"prefix_list":     knownvalue.StringExact(prefix + "-outbound-prefix"),
						}),
					),
				},
			},
		},
	})
}

func TestAccBgpFilteringRoutingProfile_Ipv4_Multicast_RouteMaps(t *testing.T) {
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
				Config: bgpFilteringRoutingProfile_Ipv4_Multicast_RouteMaps_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_bgp_filtering_routing_profile.example",
						tfjsonpath.New("ipv4").AtMapKey("multicast").AtMapKey("route_maps"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"inbound":  knownvalue.StringExact(prefix + "-route-map-in"),
							"outbound": knownvalue.StringExact(prefix + "-route-map-out"),
						}),
					),
				},
			},
		},
	})
}

func TestAccBgpFilteringRoutingProfile_Ipv4_Multicast_UnsuppressMap(t *testing.T) {
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
				Config: bgpFilteringRoutingProfile_Ipv4_Multicast_UnsuppressMap_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_bgp_filtering_routing_profile.example",
						tfjsonpath.New("ipv4").AtMapKey("multicast").AtMapKey("unsuppress_map"),
						knownvalue.StringExact(prefix + "-unsuppress"),
					),
				},
			},
		},
	})
}


func TestAccBgpFilteringRoutingProfile_Ipv4_Unicast_Basic(t *testing.T) {
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
				Config: bgpFilteringRoutingProfile_Ipv4_Unicast_Basic_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_bgp_filtering_routing_profile.example",
						tfjsonpath.New("ipv4"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"multicast": knownvalue.Null(),
							"unicast": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"conditional_advertisement":   knownvalue.Null(),
								"filter_list":                 knownvalue.Null(),
								"inbound_network_filters":     knownvalue.Null(),
								"outbound_network_filters":    knownvalue.Null(),
								"route_maps":                  knownvalue.Null(),
								"unsuppress_map":              knownvalue.Null(),
							}),
						}),
					),
				},
			},
		},
	})
}

func TestAccBgpFilteringRoutingProfile_Ipv4_Unicast_ConditionalAdvertisement_Exist(t *testing.T) {
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
				Config: bgpFilteringRoutingProfile_Ipv4_Unicast_ConditionalAdvertisement_Exist_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_bgp_filtering_routing_profile.example",
						tfjsonpath.New("ipv4").AtMapKey("unicast").AtMapKey("conditional_advertisement"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"exist": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"advertise_map": knownvalue.StringExact(prefix + "-advertise-map"),
								"exist_map":     knownvalue.StringExact(prefix + "-exist-map"),
							}),
							"non_exist": knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

func TestAccBgpFilteringRoutingProfile_Ipv4_Unicast_ConditionalAdvertisement_NonExist(t *testing.T) {
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
				Config: bgpFilteringRoutingProfile_Ipv4_Unicast_ConditionalAdvertisement_NonExist_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_bgp_filtering_routing_profile.example",
						tfjsonpath.New("ipv4").AtMapKey("unicast").AtMapKey("conditional_advertisement"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"exist": knownvalue.Null(),
							"non_exist": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"advertise_map":  knownvalue.StringExact(prefix + "-advertise-map"),
								"non_exist_map": knownvalue.StringExact(prefix + "-non-exist-map"),
							}),
						}),
					),
				},
			},
		},
	})
}

func TestAccBgpFilteringRoutingProfile_Ipv4_Unicast_FilterList(t *testing.T) {
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
				Config: bgpFilteringRoutingProfile_Ipv4_Unicast_FilterList_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_bgp_filtering_routing_profile.example",
						tfjsonpath.New("ipv4").AtMapKey("unicast").AtMapKey("filter_list"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"inbound":  knownvalue.StringExact(prefix + "-inbound"),
							"outbound": knownvalue.StringExact(prefix + "-outbound"),
						}),
					),
				},
			},
		},
	})
}

func TestAccBgpFilteringRoutingProfile_Ipv4_Unicast_InboundNetworkFilters_DistributeList(t *testing.T) {
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
				Config: bgpFilteringRoutingProfile_Ipv4_Unicast_InboundNetworkFilters_DistributeList_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_bgp_filtering_routing_profile.example",
						tfjsonpath.New("ipv4").AtMapKey("unicast").AtMapKey("inbound_network_filters"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"distribute_list": knownvalue.StringExact(prefix + "-inbound-acl"),
							"prefix_list":     knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

func TestAccBgpFilteringRoutingProfile_Ipv4_Unicast_InboundNetworkFilters_PrefixList(t *testing.T) {
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
				Config: bgpFilteringRoutingProfile_Ipv4_Unicast_InboundNetworkFilters_PrefixList_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_bgp_filtering_routing_profile.example",
						tfjsonpath.New("ipv4").AtMapKey("unicast").AtMapKey("inbound_network_filters"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"distribute_list": knownvalue.Null(),
							"prefix_list":     knownvalue.StringExact(prefix + "-inbound-prefix"),
						}),
					),
				},
			},
		},
	})
}

func TestAccBgpFilteringRoutingProfile_Ipv4_Unicast_OutboundNetworkFilters_DistributeList(t *testing.T) {
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
				Config: bgpFilteringRoutingProfile_Ipv4_Unicast_OutboundNetworkFilters_DistributeList_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_bgp_filtering_routing_profile.example",
						tfjsonpath.New("ipv4").AtMapKey("unicast").AtMapKey("outbound_network_filters"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"distribute_list": knownvalue.StringExact(prefix + "-outbound-acl"),
							"prefix_list":     knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

func TestAccBgpFilteringRoutingProfile_Ipv4_Unicast_OutboundNetworkFilters_PrefixList(t *testing.T) {
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
				Config: bgpFilteringRoutingProfile_Ipv4_Unicast_OutboundNetworkFilters_PrefixList_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_bgp_filtering_routing_profile.example",
						tfjsonpath.New("ipv4").AtMapKey("unicast").AtMapKey("outbound_network_filters"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"distribute_list": knownvalue.Null(),
							"prefix_list":     knownvalue.StringExact(prefix + "-outbound-prefix"),
						}),
					),
				},
			},
		},
	})
}

func TestAccBgpFilteringRoutingProfile_Ipv4_Unicast_RouteMaps(t *testing.T) {
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
				Config: bgpFilteringRoutingProfile_Ipv4_Unicast_RouteMaps_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_bgp_filtering_routing_profile.example",
						tfjsonpath.New("ipv4").AtMapKey("unicast").AtMapKey("route_maps"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"inbound":  knownvalue.StringExact(prefix + "-route-map-in"),
							"outbound": knownvalue.StringExact(prefix + "-route-map-out"),
						}),
					),
				},
			},
		},
	})
}

func TestAccBgpFilteringRoutingProfile_Ipv4_Unicast_UnsuppressMap(t *testing.T) {
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
				Config: bgpFilteringRoutingProfile_Ipv4_Unicast_UnsuppressMap_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_bgp_filtering_routing_profile.example",
						tfjsonpath.New("ipv4").AtMapKey("unicast").AtMapKey("unsuppress_map"),
						knownvalue.StringExact(prefix + "-unsuppress"),
					),
				},
			},
		},
	})
}

const bgpFilteringRoutingProfile_Ipv4_Multicast_Basic_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_bgp_filtering_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix
  description = "Test BGP filtering profile"

  ipv4 = {
    multicast = {
      inherit = true
    }
  }
}
`

const bgpFilteringRoutingProfile_Ipv4_Multicast_ConditionalAdvertisement_Exist_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_filters_bgp_route_map_routing_profile" "advertise_map" {
  depends_on = [panos_template.example]
  location = var.location

  name = "${var.prefix}-advertise-map"
  route_map = [
    {
      name = "10"
      action = "permit"
    }
  ]
}

resource "panos_filters_bgp_route_map_routing_profile" "exist_map" {
  depends_on = [panos_template.example]
  location = var.location

  name = "${var.prefix}-exist-map"
  route_map = [
    {
      name = "10"
      action = "permit"
    }
  ]
}

resource "panos_bgp_filtering_routing_profile" "example" {
  depends_on = [panos_filters_bgp_route_map_routing_profile.advertise_map, panos_filters_bgp_route_map_routing_profile.exist_map]
  location = var.location

  name = var.prefix

  ipv4 = {
    multicast = {
      conditional_advertisement = {
        exist = {
          advertise_map = panos_filters_bgp_route_map_routing_profile.advertise_map.name
          exist_map = panos_filters_bgp_route_map_routing_profile.exist_map.name
        }
      }
    }
  }
}
`

const bgpFilteringRoutingProfile_Ipv4_Multicast_ConditionalAdvertisement_NonExist_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_filters_bgp_route_map_routing_profile" "advertise_map" {
  depends_on = [panos_template.example]
  location = var.location

  name = "${var.prefix}-advertise-map"
  route_map = [
    {
      name = "10"
      action = "permit"
    }
  ]
}

resource "panos_filters_bgp_route_map_routing_profile" "non_exist_map" {
  depends_on = [panos_template.example]
  location = var.location

  name = "${var.prefix}-non-exist-map"
  route_map = [
    {
      name = "10"
      action = "permit"
    }
  ]
}

resource "panos_bgp_filtering_routing_profile" "example" {
  depends_on = [panos_filters_bgp_route_map_routing_profile.advertise_map, panos_filters_bgp_route_map_routing_profile.non_exist_map]
  location = var.location

  name = var.prefix

  ipv4 = {
    multicast = {
      conditional_advertisement = {
        non_exist = {
          advertise_map = panos_filters_bgp_route_map_routing_profile.advertise_map.name
          non_exist_map = panos_filters_bgp_route_map_routing_profile.non_exist_map.name
        }
      }
    }
  }
}
`

const bgpFilteringRoutingProfile_Ipv4_Multicast_FilterList_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_filters_as_path_access_list_routing_profile" "inbound" {
  depends_on = [panos_template.example]
  location = var.location

  name = "${var.prefix}-inbound"
  aspath_entries = [
    {
      name = "1"
      action = "permit"
      aspath_regex = "^65001_"
    }
  ]
}

resource "panos_filters_as_path_access_list_routing_profile" "outbound" {
  depends_on = [panos_template.example]
  location = var.location

  name = "${var.prefix}-outbound"
  aspath_entries = [
    {
      name = "1"
      action = "permit"
      aspath_regex = "^65002_"
    }
  ]
}

resource "panos_bgp_filtering_routing_profile" "example" {
  depends_on = [panos_filters_as_path_access_list_routing_profile.inbound, panos_filters_as_path_access_list_routing_profile.outbound]
  location = var.location

  name = var.prefix

  ipv4 = {
    multicast = {
      filter_list = {
        inbound = panos_filters_as_path_access_list_routing_profile.inbound.name
        outbound = panos_filters_as_path_access_list_routing_profile.outbound.name
      }
    }
  }
}
`

const bgpFilteringRoutingProfile_Ipv4_Multicast_InboundNetworkFilters_DistributeList_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_filters_access_list_routing_profile" "inbound_acl" {
  depends_on = [panos_template.example]
  location = var.location

  name = "${var.prefix}-inbound-acl"
  type = {
    ipv4 = {
      ipv4_entries = [
        {
          name = "10"
          action = "permit"
        }
      ]
    }
  }
}

resource "panos_bgp_filtering_routing_profile" "example" {
  depends_on = [panos_filters_access_list_routing_profile.inbound_acl]
  location = var.location

  name = var.prefix

  ipv4 = {
    multicast = {
      inbound_network_filters = {
        distribute_list = panos_filters_access_list_routing_profile.inbound_acl.name
      }
    }
  }
}
`

const bgpFilteringRoutingProfile_Ipv4_Multicast_InboundNetworkFilters_PrefixList_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_filters_prefix_list_routing_profile" "inbound_prefix" {
  depends_on = [panos_template.example]
  location = var.location

  name = "${var.prefix}-inbound-prefix"
  type = {
    ipv4 = {
      ipv4_entries = [
        {
          name = "10"
          action = "permit"
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

resource "panos_bgp_filtering_routing_profile" "example" {
  depends_on = [panos_filters_prefix_list_routing_profile.inbound_prefix]
  location = var.location

  name = var.prefix

  ipv4 = {
    multicast = {
      inbound_network_filters = {
        prefix_list = panos_filters_prefix_list_routing_profile.inbound_prefix.name
      }
    }
  }
}
`

const bgpFilteringRoutingProfile_Ipv4_Multicast_OutboundNetworkFilters_DistributeList_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_filters_access_list_routing_profile" "outbound_acl" {
  depends_on = [panos_template.example]
  location = var.location

  name = "${var.prefix}-outbound-acl"
  type = {
    ipv4 = {
      ipv4_entries = [
        {
          name = "10"
          action = "permit"
        }
      ]
    }
  }
}

resource "panos_bgp_filtering_routing_profile" "example" {
  depends_on = [panos_filters_access_list_routing_profile.outbound_acl]
  location = var.location

  name = var.prefix

  ipv4 = {
    multicast = {
      outbound_network_filters = {
        distribute_list = panos_filters_access_list_routing_profile.outbound_acl.name
      }
    }
  }
}
`

const bgpFilteringRoutingProfile_Ipv4_Multicast_OutboundNetworkFilters_PrefixList_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_filters_prefix_list_routing_profile" "outbound_prefix" {
  depends_on = [panos_template.example]
  location = var.location

  name = "${var.prefix}-outbound-prefix"
  type = {
    ipv4 = {
      ipv4_entries = [
        {
          name = "10"
          action = "permit"
          prefix = {
            entry = {
              network = "192.168.0.0/16"
            }
          }
        }
      ]
    }
  }
}

resource "panos_bgp_filtering_routing_profile" "example" {
  depends_on = [panos_filters_prefix_list_routing_profile.outbound_prefix]
  location = var.location

  name = var.prefix

  ipv4 = {
    multicast = {
      outbound_network_filters = {
        prefix_list = panos_filters_prefix_list_routing_profile.outbound_prefix.name
      }
    }
  }
}
`

const bgpFilteringRoutingProfile_Ipv4_Multicast_RouteMaps_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_filters_bgp_route_map_routing_profile" "route_map_in" {
  depends_on = [panos_template.example]
  location = var.location

  name = "${var.prefix}-route-map-in"
  route_map = [
    {
      name = "10"
      action = "permit"
    }
  ]
}

resource "panos_filters_bgp_route_map_routing_profile" "route_map_out" {
  depends_on = [panos_template.example]
  location = var.location

  name = "${var.prefix}-route-map-out"
  route_map = [
    {
      name = "10"
      action = "permit"
    }
  ]
}

resource "panos_bgp_filtering_routing_profile" "example" {
  depends_on = [panos_filters_bgp_route_map_routing_profile.route_map_in, panos_filters_bgp_route_map_routing_profile.route_map_out]
  location = var.location

  name = var.prefix

  ipv4 = {
    multicast = {
      route_maps = {
        inbound = panos_filters_bgp_route_map_routing_profile.route_map_in.name
        outbound = panos_filters_bgp_route_map_routing_profile.route_map_out.name
      }
    }
  }
}
`

const bgpFilteringRoutingProfile_Ipv4_Multicast_UnsuppressMap_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_filters_bgp_route_map_routing_profile" "unsuppress" {
  depends_on = [panos_template.example]
  location = var.location

  name = "${var.prefix}-unsuppress"
  route_map = [
    {
      name = "10"
      action = "permit"
    }
  ]
}

resource "panos_bgp_filtering_routing_profile" "example" {
  depends_on = [panos_filters_bgp_route_map_routing_profile.unsuppress]
  location = var.location

  name = var.prefix

  ipv4 = {
    multicast = {
      unsuppress_map = panos_filters_bgp_route_map_routing_profile.unsuppress.name
    }
  }
}
`

const bgpFilteringRoutingProfile_Ipv4_Unicast_Basic_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_bgp_filtering_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix

  ipv4 = {
    unicast = {}
  }
}
`

const bgpFilteringRoutingProfile_Ipv4_Unicast_ConditionalAdvertisement_Exist_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_filters_bgp_route_map_routing_profile" "advertise_map" {
  depends_on = [panos_template.example]
  location = var.location

  name = "${var.prefix}-advertise-map"
  route_map = [
    {
      name = "10"
      action = "permit"
    }
  ]
}

resource "panos_filters_bgp_route_map_routing_profile" "exist_map" {
  depends_on = [panos_template.example]
  location = var.location

  name = "${var.prefix}-exist-map"
  route_map = [
    {
      name = "10"
      action = "permit"
    }
  ]
}

resource "panos_bgp_filtering_routing_profile" "example" {
  depends_on = [panos_filters_bgp_route_map_routing_profile.advertise_map, panos_filters_bgp_route_map_routing_profile.exist_map]
  location = var.location

  name = var.prefix

  ipv4 = {
    unicast = {
      conditional_advertisement = {
        exist = {
          advertise_map = panos_filters_bgp_route_map_routing_profile.advertise_map.name
          exist_map = panos_filters_bgp_route_map_routing_profile.exist_map.name
        }
      }
    }
  }
}
`

const bgpFilteringRoutingProfile_Ipv4_Unicast_ConditionalAdvertisement_NonExist_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_filters_bgp_route_map_routing_profile" "advertise_map" {
  depends_on = [panos_template.example]
  location = var.location

  name = "${var.prefix}-advertise-map"
  route_map = [
    {
      name = "10"
      action = "permit"
    }
  ]
}

resource "panos_filters_bgp_route_map_routing_profile" "non_exist_map" {
  depends_on = [panos_template.example]
  location = var.location

  name = "${var.prefix}-non-exist-map"
  route_map = [
    {
      name = "10"
      action = "permit"
    }
  ]
}

resource "panos_bgp_filtering_routing_profile" "example" {
  depends_on = [panos_filters_bgp_route_map_routing_profile.advertise_map, panos_filters_bgp_route_map_routing_profile.non_exist_map]
  location = var.location

  name = var.prefix

  ipv4 = {
    unicast = {
      conditional_advertisement = {
        non_exist = {
          advertise_map = panos_filters_bgp_route_map_routing_profile.advertise_map.name
          non_exist_map = panos_filters_bgp_route_map_routing_profile.non_exist_map.name
        }
      }
    }
  }
}
`

const bgpFilteringRoutingProfile_Ipv4_Unicast_FilterList_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_filters_as_path_access_list_routing_profile" "inbound" {
  depends_on = [panos_template.example]
  location = var.location

  name = "${var.prefix}-inbound"
  aspath_entries = [
    {
      name = "1"
      action = "permit"
      aspath_regex = "^65001_"
    }
  ]
}

resource "panos_filters_as_path_access_list_routing_profile" "outbound" {
  depends_on = [panos_template.example]
  location = var.location

  name = "${var.prefix}-outbound"
  aspath_entries = [
    {
      name = "1"
      action = "permit"
      aspath_regex = "^65002_"
    }
  ]
}

resource "panos_bgp_filtering_routing_profile" "example" {
  depends_on = [panos_filters_as_path_access_list_routing_profile.inbound, panos_filters_as_path_access_list_routing_profile.outbound]
  location = var.location

  name = var.prefix

  ipv4 = {
    unicast = {
      filter_list = {
        inbound = panos_filters_as_path_access_list_routing_profile.inbound.name
        outbound = panos_filters_as_path_access_list_routing_profile.outbound.name
      }
    }
  }
}
`

const bgpFilteringRoutingProfile_Ipv4_Unicast_InboundNetworkFilters_DistributeList_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_filters_access_list_routing_profile" "inbound_acl" {
  depends_on = [panos_template.example]
  location = var.location

  name = "${var.prefix}-inbound-acl"
  type = {
    ipv4 = {
      ipv4_entries = [
        {
          name = "10"
          action = "permit"
        }
      ]
    }
  }
}

resource "panos_bgp_filtering_routing_profile" "example" {
  depends_on = [panos_filters_access_list_routing_profile.inbound_acl]
  location = var.location

  name = var.prefix

  ipv4 = {
    unicast = {
      inbound_network_filters = {
        distribute_list = panos_filters_access_list_routing_profile.inbound_acl.name
      }
    }
  }
}
`

const bgpFilteringRoutingProfile_Ipv4_Unicast_InboundNetworkFilters_PrefixList_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_filters_prefix_list_routing_profile" "inbound_prefix" {
  depends_on = [panos_template.example]
  location = var.location

  name = "${var.prefix}-inbound-prefix"
  type = {
    ipv4 = {
      ipv4_entries = [
        {
          name = "10"
          action = "permit"
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

resource "panos_bgp_filtering_routing_profile" "example" {
  depends_on = [panos_filters_prefix_list_routing_profile.inbound_prefix]
  location = var.location

  name = var.prefix

  ipv4 = {
    unicast = {
      inbound_network_filters = {
        prefix_list = panos_filters_prefix_list_routing_profile.inbound_prefix.name
      }
    }
  }
}
`

const bgpFilteringRoutingProfile_Ipv4_Unicast_OutboundNetworkFilters_DistributeList_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_filters_access_list_routing_profile" "outbound_acl" {
  depends_on = [panos_template.example]
  location = var.location

  name = "${var.prefix}-outbound-acl"
  type = {
    ipv4 = {
      ipv4_entries = [
        {
          name = "10"
          action = "permit"
        }
      ]
    }
  }
}

resource "panos_bgp_filtering_routing_profile" "example" {
  depends_on = [panos_filters_access_list_routing_profile.outbound_acl]
  location = var.location

  name = var.prefix

  ipv4 = {
    unicast = {
      outbound_network_filters = {
        distribute_list = panos_filters_access_list_routing_profile.outbound_acl.name
      }
    }
  }
}
`

const bgpFilteringRoutingProfile_Ipv4_Unicast_OutboundNetworkFilters_PrefixList_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_filters_prefix_list_routing_profile" "outbound_prefix" {
  depends_on = [panos_template.example]
  location = var.location

  name = "${var.prefix}-outbound-prefix"
  type = {
    ipv4 = {
      ipv4_entries = [
        {
          name = "10"
          action = "permit"
          prefix = {
            entry = {
              network = "192.168.0.0/16"
            }
          }
        }
      ]
    }
  }
}

resource "panos_bgp_filtering_routing_profile" "example" {
  depends_on = [panos_filters_prefix_list_routing_profile.outbound_prefix]
  location = var.location

  name = var.prefix

  ipv4 = {
    unicast = {
      outbound_network_filters = {
        prefix_list = panos_filters_prefix_list_routing_profile.outbound_prefix.name
      }
    }
  }
}
`

const bgpFilteringRoutingProfile_Ipv4_Unicast_RouteMaps_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_filters_bgp_route_map_routing_profile" "route_map_in" {
  depends_on = [panos_template.example]
  location = var.location

  name = "${var.prefix}-route-map-in"
  route_map = [
    {
      name = "10"
      action = "permit"
    }
  ]
}

resource "panos_filters_bgp_route_map_routing_profile" "route_map_out" {
  depends_on = [panos_template.example]
  location = var.location

  name = "${var.prefix}-route-map-out"
  route_map = [
    {
      name = "10"
      action = "permit"
    }
  ]
}

resource "panos_bgp_filtering_routing_profile" "example" {
  depends_on = [panos_filters_bgp_route_map_routing_profile.route_map_in, panos_filters_bgp_route_map_routing_profile.route_map_out]
  location = var.location

  name = var.prefix

  ipv4 = {
    unicast = {
      route_maps = {
        inbound = panos_filters_bgp_route_map_routing_profile.route_map_in.name
        outbound = panos_filters_bgp_route_map_routing_profile.route_map_out.name
      }
    }
  }
}
`

const bgpFilteringRoutingProfile_Ipv4_Unicast_UnsuppressMap_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_filters_bgp_route_map_routing_profile" "unsuppress" {
  depends_on = [panos_template.example]
  location = var.location

  name = "${var.prefix}-unsuppress"
  route_map = [
    {
      name = "10"
      action = "permit"
    }
  ]
}

resource "panos_bgp_filtering_routing_profile" "example" {
  depends_on = [panos_filters_bgp_route_map_routing_profile.unsuppress]
  location = var.location

  name = var.prefix

  ipv4 = {
    unicast = {
      unsuppress_map = panos_filters_bgp_route_map_routing_profile.unsuppress.name
    }
  }
}
`
