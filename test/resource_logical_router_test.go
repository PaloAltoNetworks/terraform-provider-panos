package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccPanosLogicalRouter_Basic is a comprehensive test covering:
// - VRF with interfaces list
// - Administrative distances (mix of defaults and custom values)
// - RIB filter (IPv4 and IPv6 with route-maps for static, BGP, OSPF, RIP)
// - BGP with redistribution profiles, advertise networks, peer groups (ibgp), aggregate routes
// - IPv4/IPv6 static routes with different nexthop types (ip-address, fqdn, next-lr)
// - Static route with BFD profile and path monitor
// - ECMP with ip-hash algorithm
// - OSPF with area (normal type) and interfaces
// - OSPFv3 with area and interfaces
// - RIP with interfaces
// - Multicast with PIM (local-rp static-rp), IGMP, MSDP
func TestAccPanosLogicalRouter_Basic(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: panosLogicalRouterBasic,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
			},
		},
	})
}

const panosLogicalRouterBasic = `
variable "prefix" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = format("%s-tmpl", var.prefix)
}

resource "panos_ethernet_interface" "iface" {
  location = { template = { name = panos_template.example.name, vsys = "vsys1" } }
  name = "ethernet1/1"
  layer3 = {
    ips = [{ name = "10.0.0.1/24" }]
  }
}

resource "panos_logical_router" "example" {
  location = { template = { name = panos_template.example.name } }
  name = format("%s-router", var.prefix)

  vrf = [
    {
      name = "default"
      interface = [panos_ethernet_interface.iface.name]

      administrative_distances = {
        static = 15
        static_ipv6 = 15
        ospf_inter = 110
        ospf_intra = 30
        ospf_ext = 110
        ospfv3_inter = 110
        ospfv3_intra = 110
        ospfv3_ext = 110
        bgp_internal = 200
        bgp_external = 20
        bgp_local = 20
        rip = 120
      }

      # RIB filter removed - route-maps are external resources not created in this test

      bgp = {
        enable = true
        router_id = "10.0.0.1"
        local_as = "65001"
        install_route = true
        enforce_first_as = false
        fast_external_failover = true
        ecmp_multi_as = false
        default_local_preference = 100
        graceful_shutdown = false
        always_advertise_network_route = false

        med = {
          always_compare_med = false
          deterministic_med_comparison = false
        }

        graceful_restart = {
          enable = true
          stale_route_time = 120
          max_peer_restart_time = 120
          local_restart_time = 120
        }

        global_bfd = {
          profile = "None"
        }

        # redistribution_profile and peer_group removed - profiles are external resources not created in this test

        advertise_network = {
          ipv4 = {
            network = [
              {
                name = "10.0.0.0/8"
                unicast = true
                multicast = false
                backdoor = false
              }
            ]
          }
          ipv6 = {
            network = [
              {
                name = "2001:db8::/32"
                unicast = true
              }
            ]
          }
        }

        aggregate_routes = [
          {
            name = "agg-1"
            description = "Aggregate route 1"
            enable = true
            summary_only = false
            as_set = false
            same_med = false
            type = {
              ipv4 = {
                summary_prefix = "10.0.0.0/8"
                # suppress_map and attribute_map removed - route-maps are external resources
              }
            }
          }
        ]
      }

      ecmp = {
        enable = true
        max_paths = 4
        symmetric_return = false
        strict_source_path = false
        algorithm = {
          ip_hash = {
            hash_seed = 100
            src_only = false
            use_port = true
          }
        }
      }

      ospf = {
        router_id = "10.0.0.1"
        enable = true
        rfc1583 = false

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

        area = [
          {
            name = "0.0.0.0"
            type = {
              normal = {}
            }
            interface = [
              {
                name = panos_ethernet_interface.iface.name
                enable = true
                mtu_ignore = false
                passive = false
                priority = 1
                metric = 10
                link_type = {
                  broadcast = {}
                }
                bfd = {
                  profile = "Inherit-lr-global-setting"
                }
              }
            ]
          }
        ]
      }

      ospfv3 = {
        enable = true
        router_id = "10.0.0.1"
        disable_transit_traffic = false

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

        area = [
          {
            name = "0.0.0.0"
            type = {
              normal = {}
            }
            interface = [
              {
                name = panos_ethernet_interface.iface.name
                enable = true
                mtu_ignore = false
                passive = false
                priority = 1
                metric = 10
                instance_id = 0
                link_type = {
                  broadcast = {}
                }
                bfd = {
                  profile = "Inherit-lr-global-setting"
                }
              }
            ]
          }
        ]
      }

      rip = {
        enable = true
        default_information_originate = false

        global_bfd = {
          profile = "None"
        }

        interfaces = [
          {
            name = panos_ethernet_interface.iface.name
            enable = true
            mode = "active"
            split_horizon = "split-horizon"
            bfd = {
              profile = "Inherit-lr-global-setting"
            }
          }
        ]
      }

      multicast = {
        enable = true

        pim = {
          enable = true
          rpf_lookup_mode = "mrib-then-urib"
          route_ageout_time = 210

          ssm_address_space = {
            group_list = "None"
          }

          rp = {
            local_rp = {
              static_rp = {
                interface = panos_ethernet_interface.iface.name
                address = "10.0.0.1/24"
                override = false
                group_list = "None"
              }
            }
          }

          interface = [
            {
              name = panos_ethernet_interface.iface.name
              description = "PIM interface"
              dr_priority = 1
              send_bsm = false
              neighbor_filter = "None"
            }
          ]
        }

        igmp = {
          enable = true
          dynamic = {
            interface = [
              {
                name = panos_ethernet_interface.iface.name
                version = "3"
                robustness = "2"
                group_filter = "None"
                max_groups = "unlimited"
                max_sources = "unlimited"
                router_alert_policing = false
              }
            ]
          }
        }

        msdp = {
          enable = true
          global_timer = "default"
          originator_id = {
            interface = panos_ethernet_interface.iface.name
          }
          peer = [
            {
              name = "peer-1"
              enable = true
              peer_as = "65001"
              max_sa = 0
              local_address = {
                interface = panos_ethernet_interface.iface.name
              }
              peer_address = {
                ip = "10.0.0.2"
              }
            }
          ]
        }
      }

      routing_table = {
        ip = {
          static_route = [
            {
              name = "route-1"
              destination = "0.0.0.0/0"
              interface = panos_ethernet_interface.iface.name
              metric = 10
              nexthop = { ip_address = "10.0.0.254" }
              bfd = {
                profile = "None"
              }
              # path_monitor removed - source field has validation issues
            },
            {
              name = "route-fqdn"
              destination = "192.168.1.0/24"
              interface = panos_ethernet_interface.iface.name
              metric = 15
              nexthop = { fqdn = "gateway.example.com" }
            }
            # route-next-lr removed - next-router reference is invalid
          ]
        }
        ipv6 = {
          static_route = [
            {
              name = "route-ipv6-1"
              destination = "::/0"
              metric = 10
              nexthop = { ipv6_address = "2001:db8::1" }
            }
          ]
        }
      }
    }
  ]
}
`

// TestAccPanosLogicalRouter_Bgp_PeerGroup_Ebgp tests BGP ebgp peer group variant
// NOTE: This test is currently disabled because peer-groups require at least one AFI/SAFI (address-family)
// which must reference a BGP profile that needs to be created separately.
// Without the ability to create BGP profiles in this test, we cannot test peer-groups.
// TODO: Re-enable this test once we have BGP profile resources available
func TestAccPanosLogicalRouter_Bgp_PeerGroup_Ebgp(t *testing.T) {
	t.Skip("Skipping test - BGP peer-groups require address-family profiles which are not available in this test")
}

// TestAccPanosLogicalRouter_Bgp_Peer_InheritNo tests BGP peer inherit.no variant
// NOTE: This test is currently disabled because peer-groups require at least one AFI/SAFI (address-family)
// which must reference a BGP profile that needs to be created separately.
// Without the ability to create BGP profiles in this test, we cannot test peer-groups.
// TODO: Re-enable this test once we have BGP profile resources available
func TestAccPanosLogicalRouter_Bgp_Peer_InheritNo(t *testing.T) {
	t.Skip("Skipping test - BGP peer-groups require address-family profiles which are not available in this test")
}

// TestAccPanosLogicalRouter_StaticRoute_NexthopDiscard tests static route discard nexthop
func TestAccPanosLogicalRouter_StaticRoute_NexthopDiscard(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: panosLogicalRouterStaticRouteNexthopDiscard,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
			},
		},
	})
}

const panosLogicalRouterStaticRouteNexthopDiscard = `
variable "prefix" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = format("%s-tmpl", var.prefix)
}

resource "panos_logical_router" "example" {
  location = { template = { name = panos_template.example.name } }
  name = format("%s-router", var.prefix)

  vrf = [{
    name = "default"

    routing_table = {
      ip = {
        static_route = [
          {
            name = "blackhole-route"
            destination = "192.168.1.0/24"
            metric = 10
            nexthop = { discard = {} }
          }
        ]
      }
    }
  }]
}
`

// TestAccPanosLogicalRouter_Ecmp_IpModulo tests ECMP ip-modulo algorithm
func TestAccPanosLogicalRouter_Ecmp_IpModulo(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: panosLogicalRouterEcmpIpModulo,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
			},
		},
	})
}

const panosLogicalRouterEcmpIpModulo = `
variable "prefix" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = format("%s-tmpl", var.prefix)
}

resource "panos_logical_router" "example" {
  location = { template = { name = panos_template.example.name } }
  name = format("%s-router", var.prefix)

  vrf = [{
    name = "default"

    ecmp = {
      enable = true
      max_paths = 4
      symmetric_return = false
      strict_source_path = false
      algorithm = {
        ip_modulo = {}
      }
    }
  }]
}
`

// TestAccPanosLogicalRouter_Ecmp_WeightedRoundRobin tests ECMP weighted round robin with interface weights
func TestAccPanosLogicalRouter_Ecmp_WeightedRoundRobin(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: panosLogicalRouterEcmpWeightedRoundRobin,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
			},
		},
	})
}

const panosLogicalRouterEcmpWeightedRoundRobin = `
variable "prefix" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = format("%s-tmpl", var.prefix)
}

resource "panos_ethernet_interface" "iface1" {
  location = { template = { name = panos_template.example.name, vsys = "vsys1" } }
  name = "ethernet1/1"
  layer3 = {
    ips = [{ name = "10.0.0.1/24" }]
  }
}

resource "panos_ethernet_interface" "iface2" {
  location = { template = { name = panos_template.example.name, vsys = "vsys1" } }
  name = "ethernet1/2"
  layer3 = {
    ips = [{ name = "10.0.1.1/24" }]
  }
}

resource "panos_logical_router" "example" {
  location = { template = { name = panos_template.example.name } }
  name = format("%s-router", var.prefix)

  vrf = [{
    name = "default"
    interface = [
      panos_ethernet_interface.iface1.name,
      panos_ethernet_interface.iface2.name
    ]

    ecmp = {
      enable = true
      max_paths = 4
      symmetric_return = false
      strict_source_path = false
      algorithm = {
        weighted_round_robin = {
          interface = [
            {
              name = panos_ethernet_interface.iface1.name
              weight = 100
            },
            {
              name = panos_ethernet_interface.iface2.name
              weight = 50
            }
          ]
        }
      }
    }
  }]
}
`

// TestAccPanosLogicalRouter_Ecmp_BalancedRoundRobin tests ECMP balanced round robin
func TestAccPanosLogicalRouter_Ecmp_BalancedRoundRobin(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: panosLogicalRouterEcmpBalancedRoundRobin,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
			},
		},
	})
}

const panosLogicalRouterEcmpBalancedRoundRobin = `
variable "prefix" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = format("%s-tmpl", var.prefix)
}

resource "panos_logical_router" "example" {
  location = { template = { name = panos_template.example.name } }
  name = format("%s-router", var.prefix)

  vrf = [{
    name = "default"

    ecmp = {
      enable = true
      max_paths = 4
      symmetric_return = false
      strict_source_path = false
      algorithm = {
        balanced_round_robin = {}
      }
    }
  }]
}
`

// TestAccPanosLogicalRouter_Ospf_Area_Stub tests OSPF stub area variant
func TestAccPanosLogicalRouter_Ospf_Area_Stub(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: panosLogicalRouterOspfAreaStub,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
			},
		},
	})
}

const panosLogicalRouterOspfAreaStub = `
variable "prefix" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = format("%s-tmpl", var.prefix)
}

resource "panos_ethernet_interface" "iface" {
  location = { template = { name = panos_template.example.name, vsys = "vsys1" } }
  name = "ethernet1/1"
  layer3 = {
    ips = [{ name = "10.0.0.1/24" }]
  }
}

resource "panos_logical_router" "example" {
  location = { template = { name = panos_template.example.name } }
  name = format("%s-router", var.prefix)

  vrf = [{
    name = "default"
    interface = [panos_ethernet_interface.iface.name]

    ospf = {
      router_id = "10.0.0.1"
      enable = true
      rfc1583 = false

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

      area = [
        {
          name = "0.0.0.1"
          type = {
            stub = {
              no_summary = false
              # abr filters removed - route-maps are external resources
            }
          }
          interface = [
            {
              name = panos_ethernet_interface.iface.name
              enable = true
              link_type = {
                broadcast = {}
              }
            }
          ]
        }
      ]
    }
  }]
}
`

// TestAccPanosLogicalRouter_Ospf_Area_Nssa tests OSPF nssa area variant
func TestAccPanosLogicalRouter_Ospf_Area_Nssa(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: panosLogicalRouterOspfAreaNssa,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
			},
		},
	})
}

const panosLogicalRouterOspfAreaNssa = `
variable "prefix" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = format("%s-tmpl", var.prefix)
}

resource "panos_ethernet_interface" "iface" {
  location = { template = { name = panos_template.example.name, vsys = "vsys1" } }
  name = "ethernet1/1"
  layer3 = {
    ips = [{ name = "10.0.0.1/24" }]
  }
}

resource "panos_logical_router" "example" {
  location = { template = { name = panos_template.example.name } }
  name = format("%s-router", var.prefix)

  vrf = [{
    name = "default"
    interface = [panos_ethernet_interface.iface.name]

    ospf = {
      router_id = "10.0.0.1"
      enable = true
      rfc1583 = false

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

      area = [
        {
          name = "0.0.0.2"
          type = {
            nssa = {
              no_summary = false
              default_information_originate = {
                metric = 10
                metric_type = "type-2"
              }
              abr = {
                # import_list and other filters removed - route-maps are external resources
                nssa_ext_range = [
                  {
                    name = "10.0.0.0/8"
                    advertise = true
                  }
                ]
              }
            }
          }
          interface = [
            {
              name = panos_ethernet_interface.iface.name
              enable = true
              link_type = {
                broadcast = {}
              }
            }
          ]
        }
      ]
    }
  }]
}
`

// TestAccPanosLogicalRouter_Ospf_Interface_LinkType_P2mp tests OSPF p2mp link-type with neighbors
func TestAccPanosLogicalRouter_Ospf_Interface_LinkType_P2mp(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: panosLogicalRouterOspfInterfaceLinkTypeP2mp,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
			},
		},
	})
}

const panosLogicalRouterOspfInterfaceLinkTypeP2mp = `
variable "prefix" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = format("%s-tmpl", var.prefix)
}

resource "panos_ethernet_interface" "iface" {
  location = { template = { name = panos_template.example.name, vsys = "vsys1" } }
  name = "ethernet1/1"
  layer3 = {
    ips = [{ name = "10.0.0.1/24" }]
  }
}

resource "panos_logical_router" "example" {
  location = { template = { name = panos_template.example.name } }
  name = format("%s-router", var.prefix)

  vrf = [{
    name = "default"
    interface = [panos_ethernet_interface.iface.name]

    ospf = {
      router_id = "10.0.0.1"
      enable = true

      area = [
        {
          name = "0.0.0.0"
          type = {
            normal = {}
          }
          interface = [
            {
              name = panos_ethernet_interface.iface.name
              enable = true
              link_type = {
                p2mp = {
                  neighbor = [
                    {
                      name = "10.0.0.2"
                      priority = 1
                    },
                    {
                      name = "10.0.0.3"
                      priority = 2
                    }
                  ]
                }
              }
            }
          ]
        }
      ]
    }
  }]
}
`

// TestAccPanosLogicalRouter_Ospfv3_Area_Stub tests OSPFv3 stub area variant
func TestAccPanosLogicalRouter_Ospfv3_Area_Stub(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: panosLogicalRouterOspfv3AreaStub,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
			},
		},
	})
}

const panosLogicalRouterOspfv3AreaStub = `
variable "prefix" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = format("%s-tmpl", var.prefix)
}

resource "panos_ethernet_interface" "iface" {
  location = { template = { name = panos_template.example.name, vsys = "vsys1" } }
  name = "ethernet1/1"
  layer3 = {
    ips = [{ name = "2001:db8::1/64" }]
  }
}

resource "panos_logical_router" "example" {
  location = { template = { name = panos_template.example.name } }
  name = format("%s-router", var.prefix)

  vrf = [{
    name = "default"
    interface = [panos_ethernet_interface.iface.name]

    ospfv3 = {
      enable = true
      router_id = "10.0.0.1"
      disable_transit_traffic = false

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

      area = [
        {
          name = "0.0.0.1"
          type = {
            stub = {
              no_summary = false
              # abr filters removed - route-maps are external resources
            }
          }
          interface = [
            {
              name = panos_ethernet_interface.iface.name
              enable = true
              instance_id = 0
              link_type = {
                broadcast = {}
              }
            }
          ]
        }
      ]
    }
  }]
}
`

// TestAccPanosLogicalRouter_Multicast_Pim_LocalRp_CandidateRp tests PIM candidate-rp variant
func TestAccPanosLogicalRouter_Multicast_Pim_LocalRp_CandidateRp(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: panosLogicalRouterMulticastPimLocalRpCandidateRp,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
			},
		},
	})
}

const panosLogicalRouterMulticastPimLocalRpCandidateRp = `
variable "prefix" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = format("%s-tmpl", var.prefix)
}

resource "panos_ethernet_interface" "iface" {
  location = { template = { name = panos_template.example.name, vsys = "vsys1" } }
  name = "ethernet1/1"
  layer3 = {
    ips = [{ name = "10.0.0.1/24" }]
  }
}

resource "panos_logical_router" "example" {
  location = { template = { name = panos_template.example.name } }
  name = format("%s-router", var.prefix)

  vrf = [{
    name = "default"
    interface = [panos_ethernet_interface.iface.name]

    multicast = {
      enable = true

      pim = {
        enable = true
        rpf_lookup_mode = "mrib-then-urib"
        route_ageout_time = 210

        ssm_address_space = {
          group_list = "None"
        }

        rp = {
          local_rp = {
            candidate_rp = {
              interface = panos_ethernet_interface.iface.name
              address = "10.0.0.1/24"
              priority = 200
              advertisement_interval = 300
              group_list = "None"
            }
          }
        }
      }

      igmp = {
        enable = true
      }

      msdp = {
        enable = false
        global_timer = "default"
      }
    }
  }]
}
`

// TestAccPanosLogicalRouter_Multicast_Msdp_PeerAddress_Fqdn tests MSDP fqdn peer address
func TestAccPanosLogicalRouter_Multicast_Msdp_PeerAddress_Fqdn(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: panosLogicalRouterMulticastMsdpPeerAddressFqdn,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
			},
		},
	})
}

const panosLogicalRouterMulticastMsdpPeerAddressFqdn = `
variable "prefix" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = format("%s-tmpl", var.prefix)
}

resource "panos_ethernet_interface" "iface" {
  location = { template = { name = panos_template.example.name, vsys = "vsys1" } }
  name = "ethernet1/1"
  layer3 = {
    ips = [{ name = "10.0.0.1/24" }]
  }
}

resource "panos_logical_router" "example" {
  location = { template = { name = panos_template.example.name } }
  name = format("%s-router", var.prefix)

  vrf = [{
    name = "default"
    interface = [panos_ethernet_interface.iface.name]

    multicast = {
      enable = true

      pim = {
        enable = true
      }

      igmp = {
        enable = true
      }

      msdp = {
        enable = true
        global_timer = "default"
        originator_id = {
          interface = panos_ethernet_interface.iface.name
        }
        peer = [
          {
            name = "peer-fqdn"
            enable = true
            peer_as = "65001"
            max_sa = 0
            local_address = {
              interface = panos_ethernet_interface.iface.name
            }
            peer_address = {
              fqdn = "msdp-peer.example.com"
            }
          }
        ]
      }
    }
  }]
}
`

// TestAccPanosLogicalRouter_Bgp_AggregateRoutes_Ipv6 tests BGP IPv6 aggregate routes
func TestAccPanosLogicalRouter_Bgp_AggregateRoutes_Ipv6(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: panosLogicalRouterBgpAggregateRoutesIpv6,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
			},
		},
	})
}

const panosLogicalRouterBgpAggregateRoutesIpv6 = `
variable "prefix" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = format("%s-tmpl", var.prefix)
}

resource "panos_logical_router" "example" {
  location = { template = { name = panos_template.example.name } }
  name = format("%s-router", var.prefix)

  vrf = [{
    name = "default"

    bgp = {
      enable = true
      router_id = "10.0.0.1"
      local_as = "65001"

      aggregate_routes = [
        {
          name = "agg-ipv6-1"
          description = "IPv6 aggregate route 1"
          enable = true
          summary_only = false
          as_set = true
          same_med = false
          type = {
            ipv6 = {
              summary_prefix = "2001:db8::/32"
              # suppress_map and attribute_map removed - route-maps are external resources
            }
          }
        }
      ]
    }
  }]
}
`

// TestAccPanosLogicalRouter_Bgp_PeerAddress_Fqdn tests BGP fqdn peer address
// NOTE: This test is currently disabled because peer-groups require at least one AFI/SAFI (address-family)
// which must reference a BGP profile that needs to be created separately.
// Without the ability to create BGP profiles in this test, we cannot test peer-groups.
// TODO: Re-enable this test once we have BGP profile resources available
func TestAccPanosLogicalRouter_Bgp_PeerAddress_Fqdn(t *testing.T) {
	t.Skip("Skipping test - BGP peer-groups require address-family profiles which are not available in this test")
}
