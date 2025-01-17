package provider_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	sdkErrors "github.com/PaloAltoNetworks/pango/errors"
	"github.com/PaloAltoNetworks/pango/network/interface/vlan"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccVlanInterface_1(t *testing.T) {
	t.Parallel()

	interfaceName := "vlan.1"
	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		CheckDestroy: testAccCheckPanosVlanInterfaceDestroy(
			prefix, interfaceName,
		),
		Steps: []resource.TestStep{
			{
				Config: vlanInterfaceResource1,
				ConfigVariables: map[string]config.Variable{
					"prefix":         config.StringVariable(prefix),
					"interface_name": config.StringVariable(interfaceName),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_vlan_interface.iface",
						tfjsonpath.New("name"),
						knownvalue.StringExact("vlan.1"),
					),
					statecheck.ExpectKnownValue(
						"panos_vlan_interface.iface",
						tfjsonpath.New("comment"),
						knownvalue.StringExact("vlan interface comment"),
					),
					statecheck.ExpectKnownValue(
						"panos_vlan_interface.iface",
						tfjsonpath.New("interface_management_profile"),
						knownvalue.StringExact(fmt.Sprintf("%s-profile", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_vlan_interface.iface",
						tfjsonpath.New("df_ignore"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_vlan_interface.iface",
						tfjsonpath.New("mtu"),
						knownvalue.Int64Exact(9216),
					),
					statecheck.ExpectKnownValue(
						"panos_vlan_interface.iface",
						tfjsonpath.New("adjust_tcp_mss"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"enable":              knownvalue.Bool(true),
							"ipv4_mss_adjustment": knownvalue.Int64Exact(100),
							"ipv6_mss_adjustment": knownvalue.Int64Exact(200),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_vlan_interface.iface",
						tfjsonpath.New("arp"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":       knownvalue.StringExact("192.168.0.1"),
								"hw_address": knownvalue.StringExact("aa:bb:cc:dd:ee:ff"),
								"interface":  knownvalue.StringExact("ethernet1/1"),
							}),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_vlan_interface.iface",
						tfjsonpath.New("bonjour"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"enable":    knownvalue.Bool(true),
							"group_id":  knownvalue.Int64Exact(10),
							"ttl_check": knownvalue.Bool(true),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_vlan_interface.iface",
						tfjsonpath.New("dhcp_client"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"create_default_route": knownvalue.Bool(true),
							"default_route_metric": knownvalue.Int64Exact(10),
							"enable":               knownvalue.Bool(true),
							"send_hostname": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"enable":   knownvalue.Bool(true),
								"hostname": knownvalue.StringExact("example.com"),
							}),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_vlan_interface.iface",
						tfjsonpath.New("ipv6"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"enabled":      knownvalue.Bool(true),
							"interface_id": knownvalue.StringExact("10"),
							"address":      knownvalue.Null(),
							"dhcp_client": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"accept_ra_route":      knownvalue.Bool(true),
								"default_route_metric": knownvalue.Int64Exact(10),
								"enable":               knownvalue.Bool(true),
								"preference":           knownvalue.StringExact("high"),
								"neighbor_discovery": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"dad_attempts":       knownvalue.Int64Exact(10),
									"enable_dad":         knownvalue.Bool(true),
									"enable_ndp_monitor": knownvalue.Bool(true),
									"ns_interval":        knownvalue.Int64Exact(10),
									"reachable_time":     knownvalue.Int64Exact(10),
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
											"name":       knownvalue.StringExact("::3"),
											"hw_address": knownvalue.StringExact("aa:bb:cc:dd:ee:ff"),
										}),
									}),
								}),
								"prefix_delegation": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"enable": knownvalue.ObjectExact(map[string]knownvalue.Check{
										"no": knownvalue.Null(),
										"yes": knownvalue.ObjectExact(map[string]knownvalue.Check{
											"pfx_pool_name":   knownvalue.StringExact("pool"),
											"prefix_len":      knownvalue.Int64Exact(12),
											"prefix_len_hint": knownvalue.Bool(true),
										}),
									}),
								}),
								"v6_options": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"duid_type":             knownvalue.StringExact("duid-type-ll"),
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
							"inherited":          knownvalue.Null(),
							"neighbor_discovery": knownvalue.Null(),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_vlan_interface.iface",
						tfjsonpath.New("ndp_proxy"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"enabled": knownvalue.Bool(true),
							"address": knownvalue.ListExact([]knownvalue.Check{
								knownvalue.ObjectExact(map[string]knownvalue.Check{
									"name":   knownvalue.StringExact("172.16.0.1"),
									"negate": knownvalue.Bool(true),
								}),
							}),
						}),
					),
				},
			},
		},
	})
}

const vlanInterfaceResource1 = `
variable "prefix" { type = string }
variable "interface_name" { type = string }

locals {
  template_name = format("%s-tmpl", var.prefix)
}

resource "panos_template" "template" {
  location = { panorama = {} }
  name     = local.template_name
}

resource "panos_interface_management_profile" "profile" {
  location = { template = { name = panos_template.template.name } }

  name = format("%s-profile", var.prefix)
}

resource "panos_ethernet_interface" "interface" {
  location = { template = { vsys = "vsys1", name = panos_template.template.name } }

  name   = "ethernet1/1"
  layer2 = {}
}

resource "panos_vlan_interface" "iface" {
  location = { template = { name = panos_template.template.name } }

  name    = var.interface_name
  comment = "vlan interface comment"

  interface_management_profile = panos_interface_management_profile.profile.name
  df_ignore                    = true
  mtu                          = 9216
  #netflow_profile = format("%s-profile", var.prefix)
  adjust_tcp_mss = {
    enable              = true
    ipv4_mss_adjustment = 100
    ipv6_mss_adjustment = 200
  }
  arp = [{
    name       = "192.168.0.1"
    hw_address = "aa:bb:cc:dd:ee:ff"
    interface  = panos_ethernet_interface.interface.name
  }]
  bonjour = {
    enable    = true
    group_id  = 10
    ttl_check = true
  }
  # ddns_config = {
  #   ddns_cert_profile    = format("%s-cert-profile", var.prefix)
  #   ddns_enabled         = true
  #   ddns_hostname        = "example.com"
  #   ddns_ip              = ["172.16.0.1", "172.16.0.2"]
  #   ddns_ipv6            = ["::1", "::2"]
  #   ddns_update_interval = 100
  #   ddns_vendor          = "Vendor"
  #   ddns_vendor_config   = [{ name = "name", value = "value" }]
  # }
  dhcp_client = {
    create_default_route = true
    default_route_metric = 10
    enable               = true
    send_hostname        = { enable = true, hostname = "example.com" }
  }
  ipv6 = {
    enabled      = true
    interface_id = "10"
    dhcp_client = {
      accept_ra_route      = true
      default_route_metric = 10
      enable               = true
      preference           = "high"
      neighbor_discovery = {
        dad_attempts       = 10
        enable_dad         = true
        enable_ndp_monitor = true
        ns_interval        = 10
        reachable_time     = 10
        dns_server = {
          enable = true
          source = {
            manual = {
              server = [{
                name     = "::2"
                lifetime = 4
              }]
            }
          }
        }
        dns_suffix = {
          enable = true
          source = {
            manual = {
              suffix = [{
                name     = "example.com"
                lifetime = 4
              }]
            }
          }
        }
        neighbor = [{
          name       = "::3"
          hw_address = "aa:bb:cc:dd:ee:ff"
        }]
      }
      prefix_delegation = {
        enable = {
          yes = {
            pfx_pool_name   = "pool"
            prefix_len      = 12
            prefix_len_hint = true
          }
        }
      }
      v6_options = {
        duid_type             = "duid-type-ll"
        rapid_commit          = true
        support_srvr_reconfig = true
        enable = {
          yes = {
            non_temp_addr = true
            temp_addr     = true
          }
        }
      }
    }
  }
  ndp_proxy = {
    enabled = true
    address = [{
      name   = "172.16.0.1"
      negate = true
    }]
  }
}
`

func TestAccVlanInterface_2(t *testing.T) {
	t.Parallel()

	interfaceName := "vlan.1"
	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		CheckDestroy: testAccCheckPanosVlanInterfaceDestroy(
			prefix, interfaceName,
		),
		Steps: []resource.TestStep{
			{
				Config: vlanInterfaceResource2,
				ConfigVariables: map[string]config.Variable{
					"prefix":         config.StringVariable(prefix),
					"interface_name": config.StringVariable(interfaceName),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_vlan_interface.iface",
						tfjsonpath.New("name"),
						knownvalue.StringExact("vlan.1"),
					),
					statecheck.ExpectKnownValue(
						"panos_vlan_interface.iface",
						tfjsonpath.New("ip"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name": knownvalue.StringExact("172.16.0.1"),
							}),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_vlan_interface.iface",
						tfjsonpath.New("ipv6"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
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
										"preferred_lifetime": knownvalue.StringExact("100"),
										"valid_lifetime":     knownvalue.StringExact("200"),
									}),
									"anycast": knownvalue.ObjectExact(nil),
									"prefix":  knownvalue.ObjectExact(nil),
								}),
							}),
							"dhcp_client": knownvalue.Null(),
							"inherited":   knownvalue.Null(),
							"neighbor_discovery": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"dad_attempts":       knownvalue.Int64Exact(10),
								"enable_dad":         knownvalue.Bool(true),
								"enable_ndp_monitor": knownvalue.Bool(true),
								"ns_interval":        knownvalue.Int64Exact(30),
								"reachable_time":     knownvalue.Int64Exact(50),
								"neighbor": knownvalue.ListExact([]knownvalue.Check{
									knownvalue.ObjectExact(map[string]knownvalue.Check{
										"name":       knownvalue.StringExact("::6"),
										"hw_address": knownvalue.StringExact("aa:bb:cc:dd:ee:ff"),
									}),
								}),
								"router_advertisement": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"enable":                   knownvalue.Bool(true),
									"enable_consistency_check": knownvalue.Bool(true),
									"hop_limit":                knownvalue.StringExact("5"),
									"lifetime":                 knownvalue.Int64Exact(400),
									"link_mtu":                 knownvalue.StringExact("9216"),
									"managed_flag":             knownvalue.Bool(false),
									"max_interval":             knownvalue.Int64Exact(200),
									"min_interval":             knownvalue.Int64Exact(150),
									"other_flag":               knownvalue.Bool(false),
									"reachable_time":           knownvalue.StringExact("1500"),
									"retransmission_timer":     knownvalue.StringExact("3000"),
									"router_preference":        knownvalue.StringExact("High"),
									"dns_support": knownvalue.ObjectExact(map[string]knownvalue.Check{
										"enable": knownvalue.Bool(true),
										"server": knownvalue.ListExact([]knownvalue.Check{
											knownvalue.ObjectExact(map[string]knownvalue.Check{
												"name":     knownvalue.StringExact("::7"),
												"lifetime": knownvalue.Int64Exact(400),
											}),
										}),
										"suffix": knownvalue.ListExact([]knownvalue.Check{
											knownvalue.ObjectExact(map[string]knownvalue.Check{
												"name":     knownvalue.StringExact("example.com"),
												"lifetime": knownvalue.Int64Exact(400),
											}),
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

const vlanInterfaceResource2 = `
variable "prefix" { type = string }
variable "interface_name" { type = string }

locals {
  template_name = format("%s-tmpl", var.prefix)
}

resource "panos_template" "template" {
  location = { panorama = {} }
  name     = local.template_name
}

resource "panos_vlan_interface" "iface" {
  location = { template = { name = panos_template.template.name } }

  name    = var.interface_name

  ip = [{
    name = "172.16.0.1"
  }]
  ipv6 = {
    enabled      = true
    interface_id = "10"
    address = [{
      name                = "::1"
      enable_on_interface = true
      advertise = {
        auto_config_flag   = true
        enable             = true
        onlink_flag        = true
        preferred_lifetime = "100"
        valid_lifetime     = "200"
      }
      anycast = {}
      prefix  = {}
    }]
    neighbor_discovery = {
      dad_attempts       = 10
      enable_dad         = true
      enable_ndp_monitor = true
      ns_interval        = 30
      reachable_time     = 50
      neighbor = [{
        name       = "::6"
        hw_address = "aa:bb:cc:dd:ee:ff"
      }]
      router_advertisement = {
        enable                   = true
        enable_consistency_check = true
        hop_limit                = "5"
        lifetime                 = 400
        link_mtu                 = "9216"
        managed_flag             = false
        max_interval             = 200
        min_interval             = 150
        other_flag               = false
        reachable_time           = "1500"
        retransmission_timer     = "3000"
        router_preference        = "High"
        dns_support = {
          enable = true
          server = [{ name = "::7", lifetime = 400 }]
          suffix = [{ name = "example.com", lifetime = 400 }]
        }
      }
    }
  }
}
`

func TestAccVlanInterface_3(t *testing.T) {
	t.Parallel()

	interfaceName := "vlan.1"
	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		CheckDestroy: testAccCheckPanosVlanInterfaceDestroy(
			prefix, interfaceName,
		),
		Steps: []resource.TestStep{
			{
				Config: vlanInterfaceResource3,
				ConfigVariables: map[string]config.Variable{
					"prefix":         config.StringVariable(prefix),
					"interface_name": config.StringVariable(interfaceName),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_vlan_interface.iface",
						tfjsonpath.New("name"),
						knownvalue.StringExact("vlan.1"),
					),
					statecheck.ExpectKnownValue(
						"panos_vlan_interface.iface",
						tfjsonpath.New("ipv6"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
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
														"identifier": knownvalue.Int64Exact(4095),
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
													"preferred_lifetime": knownvalue.StringExact("200"),
													"valid_lifetime":     knownvalue.StringExact("300"),
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
									"ns_interval":        knownvalue.Int64Exact(30),
									"reachable_time":     knownvalue.Int64Exact(60),
									"dns_server": knownvalue.ObjectExact(map[string]knownvalue.Check{
										"enable": knownvalue.Bool(true),
										"source": knownvalue.ObjectExact(map[string]knownvalue.Check{
											"dhcpv6": knownvalue.ObjectExact(map[string]knownvalue.Check{
												// "prefix_pool": knownvalue.StringExact(fmt.Sprintf("%s-pool", prefix))
												"prefix_pool": knownvalue.Null(),
											}),
											"manual": knownvalue.Null(),
										}),
									}),
									"dns_suffix": knownvalue.ObjectExact(map[string]knownvalue.Check{
										"enable": knownvalue.Bool(true),
										"source": knownvalue.ObjectExact(map[string]knownvalue.Check{
											"dhcpv6": knownvalue.Null(),
											"manual": knownvalue.ObjectExact(map[string]knownvalue.Check{
												"suffix": knownvalue.ListExact([]knownvalue.Check{
													knownvalue.ObjectExact(map[string]knownvalue.Check{
														"name":     knownvalue.StringExact("::6"),
														"lifetime": knownvalue.Int64Exact(100),
													}),
												}),
											}),
										}),
									}),
									"neighbor": knownvalue.ListExact([]knownvalue.Check{
										knownvalue.ObjectExact(map[string]knownvalue.Check{
											"name":       knownvalue.StringExact("::4"),
											"hw_address": knownvalue.StringExact("aa:bb:cc:dd:ee:ff"),
										}),
									}),
									"router_advertisement": knownvalue.ObjectExact(map[string]knownvalue.Check{
										"enable":                   knownvalue.Bool(true),
										"enable_consistency_check": knownvalue.Bool(true),
										"hop_limit":                knownvalue.StringExact("10"),
										"lifetime":                 knownvalue.Int64Exact(400),
										"link_mtu":                 knownvalue.StringExact("9216"),
										"managed_flag":             knownvalue.Bool(true),
										"max_interval":             knownvalue.Int64Exact(200),
										"min_interval":             knownvalue.Int64Exact(100),
										"other_flag":               knownvalue.Bool(true),
										"reachable_time":           knownvalue.StringExact("2000"),
										"retransmission_timer":     knownvalue.StringExact("5000"),
										"router_preference":        knownvalue.StringExact("High"),
									}),
								}),
							}),
							"interface_id":       knownvalue.StringExact("EUI-64"),
							"neighbor_discovery": knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

const vlanInterfaceResource3 = `
variable "prefix" { type = string }
variable "interface_name" { type = string }

locals {
  template_name = format("%s-tmpl", var.prefix)
}

resource "panos_template" "template" {
  location = { panorama = {} }
  name     = local.template_name
}

resource "panos_vlan_interface" "iface" {
  location = { template = { name = panos_template.template.name } }

  name    = var.interface_name

  ipv6 = {
    inherited = {
      enable = true
      assign_addr = [
        {
          name = "172.16.0.1"
          type = {
            gua = {
              enable_on_interface = true
              #prefix_pool         = format("%s-pool", var.prefix)
              advertise = {
                auto_config_flag = true
                enable           = true
                onlink_flag      = true
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
              pool_type = {
                dynamic_id = { identifier = 4095 }
              }
            }
          }
        },
        {
          name = "172.16.0.3"
          type = {
            ula = {
              address             = "fd12:3456:789a:1::1"
              anycast             = true
              enable_on_interface = true
              prefix              = true
              advertise = {
                auto_config_flag   = true
                enable             = true
                onlink_flag        = true
                preferred_lifetime = "200"
                valid_lifetime     = "300"
              }
            }
          }
        },
      ]
      neighbor_discovery = {
        dad_attempts       = 10
        enable_dad         = true
        enable_ndp_monitor = true
        ns_interval        = 30
        reachable_time     = 60
        dns_server = {
          enable = true
          source = {
            dhcpv6 = {
              #prefix_pool = format("%s-pool", var.prefix)
            }
          }
        }
        dns_suffix = {
          enable = true
          source = {
            manual = {
              suffix = [{ name = "::6", lifetime = 100 }]
            }
          }
        }
        neighbor = [{
          name       = "::4"
          hw_address = "aa:bb:cc:dd:ee:ff"
        }]
        router_advertisement = {
          enable                   = true
          enable_consistency_check = true
          hop_limit                = "10"
          lifetime                 = 400
          link_mtu                 = "9216"
          managed_flag             = true
          max_interval             = 200
          min_interval             = 100
          other_flag               = true
          reachable_time           = "2000"
          retransmission_timer     = "5000"
          router_preference        = "High"
        }
      }
    }
  }
}
`

func testAccCheckPanosVlanInterfaceDestroy(prefix string, entry string) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		api := vlan.NewService(sdkClient)
		ctx := context.TODO()

		location := vlan.NewTemplateLocation()
		location.Template.Template = fmt.Sprintf("%s-tmpl", prefix)

		reply, err := api.Read(ctx, *location, entry, "show")
		if err != nil && !sdkErrors.IsObjectNotFound(err) {
			return fmt.Errorf("reading ethernet entry via sdk: %v", err)
		}

		if reply != nil {
			err := fmt.Errorf("terraform didn't delete the server entry properly")
			delErr := api.Delete(ctx, *location, entry)
			if delErr != nil {
				return errors.Join(err, delErr)
			}
			return err
		}

		return nil
	}
}
