// Given YAML definition of a resource (logical-router.yaml) and a resource under test (panosLogicalRouterTmpl1) add all attributes from logical-router.yaml that are not within resource under test to knownvalue.ObjectExact calls as knownvalue.Null().
package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	//"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccPanosLogicalRouter_1(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: panosLogicalRouterTmpl1,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_logical_router.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-router", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_logical_router.example",
						tfjsonpath.New("vrf"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":       knownvalue.StringExact("default"),
								"rib_filter": knownvalue.Null(),
								"administrative_distances": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"static":       knownvalue.Int64Exact(10),
									"static_ipv6":  knownvalue.Int64Exact(10),
									"ospf_inter":   knownvalue.Int64Exact(110),
									"ospf_intra":   knownvalue.Int64Exact(110),
									"ospf_ext":     knownvalue.Int64Exact(110),
									"ospfv3_inter": knownvalue.Int64Exact(110),
									"ospfv3_intra": knownvalue.Int64Exact(110),
									"ospfv3_ext":   knownvalue.Int64Exact(110),
									"bgp_internal": knownvalue.Int64Exact(200),
									"bgp_external": knownvalue.Int64Exact(20),
									"bgp_local":    knownvalue.Int64Exact(20),
									"rip":          knownvalue.Int64Exact(120),
								}),
								"bgp": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"enable":                         knownvalue.Bool(true),
									"router_id":                      knownvalue.StringExact("10.0.0.1"),
									"local_as":                       knownvalue.StringExact("65000"),
									"install_route":                  knownvalue.Bool(false),
									"enforce_first_as":               knownvalue.Bool(false),
									"fast_external_failover":         knownvalue.Bool(false),
									"ecmp_multi_as":                  knownvalue.Bool(false),
									"default_local_preference":       knownvalue.Int64Exact(100),
									"graceful_shutdown":              knownvalue.Bool(false),
									"always_advertise_network_route": knownvalue.Bool(false),
									"med": knownvalue.ObjectExact(map[string]knownvalue.Check{
										"always_compare_med":           knownvalue.Bool(false),
										"deterministic_med_comparison": knownvalue.Bool(false),
									}),
									"graceful_restart": knownvalue.ObjectExact(map[string]knownvalue.Check{
										"enable":                knownvalue.Bool(false),
										"stale_route_time":      knownvalue.Int64Exact(120),
										"max_peer_restart_time": knownvalue.Int64Exact(120),
										"local_restart_time":    knownvalue.Int64Exact(120),
									}),
									"global_bfd": knownvalue.ObjectExact(map[string]knownvalue.Check{
										"profile": knownvalue.StringExact("None"),
									}),
									"redistribution_profile": knownvalue.Null(),
									"advertise_network":      knownvalue.Null(),
									"peer_group":             knownvalue.Null(),
									"aggregate_routes":       knownvalue.Null(),
								}),
								"ecmp": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"enable":             knownvalue.Bool(false),
									"max_paths":          knownvalue.Int64Exact(2),
									"symmetric_return":   knownvalue.Bool(false),
									"strict_source_path": knownvalue.Bool(false),
									"algorithm": knownvalue.ObjectExact(map[string]knownvalue.Check{
										"ip_hash": knownvalue.ObjectExact(map[string]knownvalue.Check{
											"hash_seed": knownvalue.Int64Exact(100),
											"src_only":  knownvalue.Bool(true),
											"use_port":  knownvalue.Bool(true),
										}),
										"balanced_round_robin": knownvalue.Null(),
										"ip_modulo":            knownvalue.Null(),
										"weighted_round_robin": knownvalue.Null(),
									}),
								}),
								"interface": knownvalue.ListExact([]knownvalue.Check{}),
								"multicast": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"enable": knownvalue.Bool(true),
									"pim": knownvalue.ObjectExact(map[string]knownvalue.Check{
										"enable":            knownvalue.Bool(false),
										"rpf_lookup_mode":   knownvalue.StringExact("mrib-then-urib"),
										"route_ageout_time": knownvalue.Int64Exact(210),
										"ssm_address_space": knownvalue.ObjectExact(map[string]knownvalue.Check{
											"group_list": knownvalue.StringExact("None"),
										}),
										"if_timer_global":  knownvalue.Null(),
										"group_permission": knownvalue.Null(),
										"rp":               knownvalue.Null(),
										"spt_threshold":    knownvalue.Null(),
										"interface":        knownvalue.Null(),
									}),
									"igmp": knownvalue.ObjectExact(map[string]knownvalue.Check{
										"enable":  knownvalue.Bool(true),
										"dynamic": knownvalue.Null(),
										"static":  knownvalue.Null(),
									}),
									"msdp": knownvalue.ObjectExact(map[string]knownvalue.Check{
										"enable":                knownvalue.Bool(false),
										"global_timer":          knownvalue.StringExact("default"),
										"originator_id":         knownvalue.Null(),
										"peer":                  knownvalue.Null(),
										"global_authentication": knownvalue.Null(),
									}),
									"static_route": knownvalue.Null(),
								}),
								"ospf": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"router_id": knownvalue.StringExact("10.0.0.1"),
									"enable":    knownvalue.Bool(true),
									"rfc1583":   knownvalue.Bool(false),
									"global_bfd": knownvalue.ObjectExact(map[string]knownvalue.Check{
										"profile": knownvalue.StringExact("None"),
									}),
									"graceful_restart": knownvalue.ObjectExact(map[string]knownvalue.Check{
										"enable":                    knownvalue.Bool(false),
										"grace_period":              knownvalue.Int64Exact(120),
										"helper_enable":             knownvalue.Bool(false),
										"strict_lsa_checking":       knownvalue.Bool(false),
										"max_neighbor_restart_time": knownvalue.Int64Exact(140),
									}),
									"spf_timer":              knownvalue.Null(),
									"global_if_timer":        knownvalue.Null(),
									"redistribution_profile": knownvalue.Null(),
									"area":                   knownvalue.Null(),
								}),
								"ospfv3": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"enable":                  knownvalue.Bool(true),
									"router_id":               knownvalue.StringExact("10.0.0.1"),
									"disable_transit_traffic": knownvalue.Bool(false),
									"global_bfd": knownvalue.ObjectExact(map[string]knownvalue.Check{
										"profile": knownvalue.StringExact("None"),
									}),
									"graceful_restart": knownvalue.ObjectExact(map[string]knownvalue.Check{
										"enable":                    knownvalue.Bool(false),
										"grace_period":              knownvalue.Int64Exact(120),
										"helper_enable":             knownvalue.Bool(false),
										"strict_lsa_checking":       knownvalue.Bool(false),
										"max_neighbor_restart_time": knownvalue.Int64Exact(140),
									}),
									"spf_timer":              knownvalue.Null(),
									"global_if_timer":        knownvalue.Null(),
									"redistribution_profile": knownvalue.Null(),
									"area":                   knownvalue.Null(),
								}),
								"rip": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"enable":                        knownvalue.Bool(true),
									"default_information_originate": knownvalue.Bool(false),
									"global_bfd": knownvalue.ObjectExact(map[string]knownvalue.Check{
										"profile": knownvalue.StringExact("None"),
									}),
									"global_timer":                    knownvalue.Null(),
									"auth_profile":                    knownvalue.Null(),
									"redistribution_profile":          knownvalue.Null(),
									"global_inbound_distribute_list":  knownvalue.Null(),
									"global_outbound_distribute_list": knownvalue.Null(),
									"interfaces":                      knownvalue.Null(),
								}),
								"routing_table": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"ip": knownvalue.ObjectExact(map[string]knownvalue.Check{
										"static_route": knownvalue.ListExact([]knownvalue.Check{
											knownvalue.ObjectExact(map[string]knownvalue.Check{
												"name":        knownvalue.StringExact("default-route"),
												"destination": knownvalue.StringExact("0.0.0.0/0"),
												"nexthop": knownvalue.ObjectExact(map[string]knownvalue.Check{
													"ip_address": knownvalue.StringExact("10.0.0.1"),
													"discard":    knownvalue.Null(),
													"next_lr":    knownvalue.Null(),
													"fqdn":       knownvalue.Null(),
												}),
												"interface":               knownvalue.Null(),
												"metric":                  knownvalue.Int64Exact(10),
												"administrative_distance": knownvalue.Null(),
												"bfd":                     knownvalue.Null(),
												"path_monitor":            knownvalue.Null(),
											}),
										}),
									}),
									"ipv6": knownvalue.Null(),
								}),
							}),
						}),
					),
				},
			},
		},
	})
}

const panosLogicalRouterTmpl1 = `
variable "prefix" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }

  name = format("%s-tmpl", var.prefix)
}

resource "panos_ethernet_interface" "iface" {
  location = { template = { name = panos_template.example.name, vsys = "vsys1" } }

  name = "ethernet1/1"

  layer3 = {
    ips = [{ name = "10.0.0.1/32" }]
  }
}

resource "panos_logical_router" "example" {
  location = { template = { name = panos_template.example.name } }

  name = format("%s-router", var.prefix)

  vrf = [{
    name = "default"

    administrative_distances = {
      static = 10
      static_ipv6 = 10
      ospf_inter = 110
      ospf_intra = 110
      ospf_ext = 110
      ospfv3_inter = 110
      ospfv3_intra = 110
      ospfv3_ext = 110
      bgp_internal = 200
      bgp_external = 20
      bgp_local = 20
      rip = 120
    }
    bgp = {
      enable = true
      router_id = "10.0.0.1"
      local_as = "65000"
      install_route = false
      enforce_first_as = false
      fast_external_failover = false
      ecmp_multi_as = false
      default_local_preference = 100
      graceful_shutdown = false
      always_advertise_network_route = false
      med = {
        always_compare_med = false
        deterministic_med_comparison = false
      }
      graceful_restart = {
        enable = false
        stale_route_time = 120
        max_peer_restart_time = 120
        local_restart_time = 120
      }
      global_bfd = {
        profile = "None"
      }
      #advertise_network = {
      #  ipv4 = {
      #    network = []
      #  }
      #  ipv6 = {
      #    network = []
      #  }
      #}
      #peer_group = []
      #aggregate_routes = []
    }
    ecmp = {
      enable = false
      max_paths = 2
      symmetric_return = false
      strict_source_path = false
      algorithm = {
        ip_hash = {
          hash_seed = 100
          src_only = true
          use_port = true
        }
      }
    }
    interface = []
    multicast = {
      enable = true
      #static_route = []
      pim = {
        enable = false
        rpf_lookup_mode = "mrib-then-urib"
        route_ageout_time = 210
        #if_timer_global = ""
        #group_permission = ""
        ssm_address_space = {
          group_list = "None"
        }
        #rp = {
        #  local_rp = {
        #    candidate_rp = {
        #      interface = panos_ethernet_interface.iface.name
        #      address = "10.0.0.1/32"
        #      priority = 200
        #      advertisement_interval = 300
        #      group_list = "group-list"
        #    }
        #  }
        #  external_rp = []
        #}
        #spt_threshold = []
        #interface = []
      }
      igmp = {
        enable = true
        #dynamic = {
        #  interface = []
        #}
        #static = []
      }
      msdp = {
        enable = false
        global_timer = "default"
        #originator_id = {
        #  interface = panos_ethernet_interface.iface.name
        #  ip = "10.0.0.1"
        #}
        #peer = []
      }
    }
    ospf = {
      router_id = "10.0.0.1"
      enable = true
      rfc1583 = false
      #spf_timer = ""
      #global_if_timer = ""
      #redistribution_profile = ""
      global_bfd = {
        profile = "None"
      }
      graceful_restart = {
        enable = false
        grace_period = 120
        helper_enable = false
        strict_lsa_checking = false
        max_neighbor_restart_time = 140
      }
    }
    ospfv3 = {
      enable = true
      router_id = "10.0.0.1"
      disable_transit_traffic = false
      #spf_timer = ""
      #global_if_timer = ""
      #redistribution_profile = ""
      global_bfd = {
        profile = "None"
      }
      graceful_restart = {
        enable = false
        grace_period = 120
        helper_enable = false
        strict_lsa_checking = false
        max_neighbor_restart_time = 140
      }
    }
    #rib_filter = {
    #  ipv4 = {
    #    static = {
    #      route_map = ""
    #    }
    #    bgp = {
    #      route_map = ""
    #    }
    #    ospf = {
    #      route_map = ""
    #    }
    #    rip = {
    #      route_map = ""
    #    }
    #  }
    #  ipv6 = {
    #    static = {
    #      route_map = ""
    #    }
    #    bgp = {
    #      route_map = ""
    #    }
    #    ospfv3 = {
    #      route_map = ""
    #    }
    #  }
    #}
    rip = {
      enable = true
      default_information_originate = false
      #global_timer = ""
      #auth_profile = ""
      #redistribution_profile = ""
      global_bfd = {
        profile = "None"
      }
      #global_inbound_distribute_list = {
      #  access_list = ""
      #}
      #global_outbound_distribute_list = {
      #  access_list = ""
      #}
      #interfaces = []
    }
    routing_table = {
      ip = {
        static_route = [{
          name = "default-route"
          destination = "0.0.0.0/0"
          #interface = panos_ethernet_interface.iface.name
          preference = 100
          nexthop = { ip_address = "10.0.0.1" }
        }]
      }
      #ipv6 = {
      #  static_route = []
      #}
    }
  }]
}
`
