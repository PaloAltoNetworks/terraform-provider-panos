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

func TestAccPanosVirtualRouter_Basic(t *testing.T) {
	t.Parallel()

	resName := "vr_profile"
	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	templateName := "acc-vrouter"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: makePanosVirtualRouterConfig(resName),
				ConfigVariables: map[string]config.Variable{
					"name_suffix":   config.StringVariable(nameSuffix),
					"router_name":   config.StringVariable(resName),
					"template_name": config.StringVariable(templateName),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_virtual_router."+resName,
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf(
							"%s-%s",
							resName,
							nameSuffix,
						)),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router."+resName,
						tfjsonpath.New("administrative_distances").AtMapKey("static"),
						knownvalue.Int64Exact(15),
					),
				},
			},
		},
	})
}

func TestAccPanosVirtualRouter_WithEthernetInterface(t *testing.T) {
	t.Parallel()

	resName := "vr_profile_if"
	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	templateName := "acc-vrouter"
	interfaceName := "ethernet1/41"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: makePanosVirtualRouterConfig(resName),
				ConfigVariables: map[string]config.Variable{
					"name_suffix":      config.StringVariable(nameSuffix),
					"router_name":      config.StringVariable(resName),
					"with_interface":   config.BoolVariable(true),
					"ethernet_if_name": config.StringVariable(interfaceName),
					"template_name":    config.StringVariable(templateName),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_virtual_router."+resName,
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf(
							"%s-%s",
							resName,
							nameSuffix,
						)),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router."+resName,
						tfjsonpath.
							New("interfaces").
							AtSliceIndex(0),
						knownvalue.StringExact("ethernet1/41"),
					),
				},
			},
		},
	})
}

func TestAccPanosVirtualRouter_AdministrativeDistances(t *testing.T) {
	t.Parallel()

	resName := "vr_admin_dist"
	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	templateName := "acc-vrouter-admin-dist"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualRouterAdministrativeDistances,
				ConfigVariables: map[string]config.Variable{
					"name_suffix":   config.StringVariable(nameSuffix),
					"router_name":   config.StringVariable(resName),
					"template_name": config.StringVariable(templateName),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_virtual_router."+resName,
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-%s", resName, nameSuffix)),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router."+resName,
						tfjsonpath.New("administrative_distances").AtMapKey("ebgp"),
						knownvalue.Int64Exact(25),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router."+resName,
						tfjsonpath.New("administrative_distances").AtMapKey("ibgp"),
						knownvalue.Int64Exact(210),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router."+resName,
						tfjsonpath.New("administrative_distances").AtMapKey("ospf_ext"),
						knownvalue.Int64Exact(115),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router."+resName,
						tfjsonpath.New("administrative_distances").AtMapKey("ospf_int"),
						knownvalue.Int64Exact(35),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router."+resName,
						tfjsonpath.New("administrative_distances").AtMapKey("ospfv3_ext"),
						knownvalue.Int64Exact(120),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router."+resName,
						tfjsonpath.New("administrative_distances").AtMapKey("ospfv3_int"),
						knownvalue.Int64Exact(40),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router."+resName,
						tfjsonpath.New("administrative_distances").AtMapKey("rip"),
						knownvalue.Int64Exact(125),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router."+resName,
						tfjsonpath.New("administrative_distances").AtMapKey("static"),
						knownvalue.Int64Exact(15),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router."+resName,
						tfjsonpath.New("administrative_distances").AtMapKey("static_ipv6"),
						knownvalue.Int64Exact(20),
					),
				},
			},
		},
	})
}

func TestAccPanosVirtualRouter_Ecmp_BalancedRoundRobin(t *testing.T) {
	t.Parallel()

	resName := "vr_ecmp_brr"
	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	templateName := "acc-vrouter-ecmp-brr"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualRouterEcmpBalancedRoundRobin,
				ConfigVariables: map[string]config.Variable{
					"name_suffix":   config.StringVariable(nameSuffix),
					"router_name":   config.StringVariable(resName),
					"template_name": config.StringVariable(templateName),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_virtual_router."+resName,
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-%s", resName, nameSuffix)),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router."+resName,
						tfjsonpath.New("ecmp").AtMapKey("enable"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router."+resName,
						tfjsonpath.New("ecmp").AtMapKey("max_paths"),
						knownvalue.Int64Exact(4),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router."+resName,
						tfjsonpath.New("ecmp").AtMapKey("symmetric_return"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router."+resName,
						tfjsonpath.New("ecmp").AtMapKey("strict_source_path"),
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router."+resName,
						tfjsonpath.New("ecmp").AtMapKey("algorithm").AtMapKey("balanced_round_robin"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}

func TestAccPanosVirtualRouter_Ecmp_IpModulo(t *testing.T) {
	t.Parallel()

	resName := "vr_ecmp_ipmodulo"
	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	templateName := "acc-vrouter-ecmp-ipmodulo"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualRouterEcmpIpModulo,
				ConfigVariables: map[string]config.Variable{
					"name_suffix":   config.StringVariable(nameSuffix),
					"router_name":   config.StringVariable(resName),
					"template_name": config.StringVariable(templateName),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_virtual_router."+resName,
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-%s", resName, nameSuffix)),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router."+resName,
						tfjsonpath.New("ecmp").AtMapKey("enable"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router."+resName,
						tfjsonpath.New("ecmp").AtMapKey("max_paths"),
						knownvalue.Int64Exact(3),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router."+resName,
						tfjsonpath.New("ecmp").AtMapKey("algorithm").AtMapKey("ip_modulo"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}

const testAccVirtualRouterEcmpIpModulo = `
variable "name_suffix" { type = string }
variable "router_name" { type = string }
variable "template_name" { type = string }

resource "panos_template" "template" {
  name = "${var.template_name}-${var.name_suffix}"

  location = {
    panorama = {
      panorama_device = "localhost.localdomain"
    }
  }
}

resource "panos_virtual_router" "vr_ecmp_ipmodulo" {
  location = {
    template = {
      name = panos_template.template.name
    }
  }

  name = "${var.router_name}-${var.name_suffix}"

  ecmp = {
    enable = true
    max_paths = 3
    algorithm = {
      ip_modulo = {}
    }
  }
}
`

const testAccVirtualRouterEcmpBalancedRoundRobin = `
variable "name_suffix" { type = string }
variable "router_name" { type = string }
variable "template_name" { type = string }

resource "panos_template" "template" {
  name = "${var.template_name}-${var.name_suffix}"

  location = {
    panorama = {
      panorama_device = "localhost.localdomain"
    }
  }
}

resource "panos_virtual_router" "vr_ecmp_brr" {
  location = {
    template = {
      name = panos_template.template.name
    }
  }

  name = "${var.router_name}-${var.name_suffix}"

  ecmp = {
    enable = true
    max_paths = 4
    symmetric_return = true
    strict_source_path = false
    algorithm = {
      balanced_round_robin = {}
    }
  }
}
`

const testAccVirtualRouterAdministrativeDistances = `
variable "name_suffix" { type = string }
variable "router_name" { type = string }
variable "template_name" { type = string }

resource "panos_template" "template" {
  name = "${var.template_name}-${var.name_suffix}"

  location = {
    panorama = {
      panorama_device = "localhost.localdomain"
    }
  }
}

resource "panos_virtual_router" "vr_admin_dist" {
  location = {
    template = {
      name = panos_template.template.name
    }
  }

  name = "${var.router_name}-${var.name_suffix}"

  administrative_distances = {
    ebgp        = 25
    ibgp        = 210
    ospf_ext    = 115
    ospf_int    = 35
    ospfv3_ext  = 120
    ospfv3_int  = 40
    rip         = 125
    static      = 15
    static_ipv6 = 20
  }
}
`

func TestAccPanosVirtualRouter_Multicast(t *testing.T) {
	t.Parallel()

	resName := "vr_multicast"
	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	templateName := "acc-vrouter-multicast"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualRouterMulticast,
				ConfigVariables: map[string]config.Variable{
					"name_suffix":   config.StringVariable(nameSuffix),
					"router_name":   config.StringVariable(resName),
					"template_name": config.StringVariable(templateName),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_virtual_router."+resName,
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-%s", resName, nameSuffix)),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router."+resName,
						tfjsonpath.New("multicast").AtMapKey("enable"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router."+resName,
						tfjsonpath.New("multicast"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}

const testAccVirtualRouterMulticast = `
variable "name_suffix" { type = string }
variable "router_name" { type = string }
variable "template_name" { type = string }

resource "panos_template" "template" {
  name = "${var.template_name}-${var.name_suffix}"

  location = {
    panorama = {
      panorama_device = "localhost.localdomain"
    }
  }
}

resource "panos_ethernet_interface" "eth3" {
  location = {
    template = {
      vsys = "vsys1"
      name = panos_template.template.name
    }
  }

  name = "ethernet1/3"

  layer3 = {
    mtu  = 1400
    ips  = [{ name = "10.1.3.1/24" }]
  }
}

resource "panos_virtual_router" "vr_multicast" {
  location = {
    template = {
      name = panos_template.template.name
    }
  }

  name = "${var.router_name}-${var.name_suffix}"

  interfaces = [panos_ethernet_interface.eth3.name]

  multicast = {
    enable = true
    interface_group = [
      {
        name = "mcast-grp-1"
        description = "Test multicast group"
        interfaces = ["ethernet1/3"]
        group_permission = {
          any_source_multicast = [
            {
              name = "asm-1"
              group_address = "224.0.0.0/4"
            }
          ]
        }
      }
    ]
  }
}
`

func TestAccPanosVirtualRouter_Bgp_Basic(t *testing.T) {
	t.Parallel()

	resName := "vr_bgp_basic"
	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	templateName := "acc-vrouter-bgp-basic"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualRouterBgpBasic,
				ConfigVariables: map[string]config.Variable{
					"name_suffix":   config.StringVariable(nameSuffix),
					"router_name":   config.StringVariable(resName),
					"template_name": config.StringVariable(templateName),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_virtual_router."+resName,
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-%s", resName, nameSuffix)),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router."+resName,
						tfjsonpath.New("protocol").AtMapKey("bgp").AtMapKey("enable"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router."+resName,
						tfjsonpath.New("protocol").AtMapKey("bgp").AtMapKey("router_id"),
						knownvalue.StringExact("10.0.0.1"),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router."+resName,
						tfjsonpath.New("protocol").AtMapKey("bgp").AtMapKey("local_as"),
						knownvalue.StringExact("65001"),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router."+resName,
						tfjsonpath.New("protocol").AtMapKey("bgp").AtMapKey("install_route"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router."+resName,
						tfjsonpath.New("protocol").AtMapKey("bgp").AtMapKey("reject_default_route"),
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router."+resName,
						tfjsonpath.New("protocol").AtMapKey("bgp").AtMapKey("allow_redist_default_route"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router."+resName,
						tfjsonpath.New("protocol").AtMapKey("bgp").AtMapKey("ecmp_multi_as"),
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router."+resName,
						tfjsonpath.New("protocol").AtMapKey("bgp").AtMapKey("enforce_first_as"),
						knownvalue.Bool(true),
					),
				},
			},
		},
	})
}

const testAccVirtualRouterBgpBasic = `
variable "name_suffix" { type = string }
variable "router_name" { type = string }
variable "template_name" { type = string }

resource "panos_template" "template" {
  name = "${var.template_name}-${var.name_suffix}"

  location = {
    panorama = {
      panorama_device = "localhost.localdomain"
    }
  }
}

resource "panos_virtual_router" "vr_bgp_basic" {
  location = {
    template = {
      name = panos_template.template.name
    }
  }

  name = "${var.router_name}-${var.name_suffix}"

  protocol = {
    bgp = {
      enable = true
      router_id = "10.0.0.1"
      local_as = "65001"
      install_route = true
      reject_default_route = false
      allow_redist_default_route = true
      ecmp_multi_as = false
      enforce_first_as = true
    }
  }
}
`

func TestAccPanosVirtualRouter_Bgp_Dampening(t *testing.T) {
	t.Parallel()

	resName := "vr_bgp_dampening"
	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	templateName := "acc-vrouter-bgp-damp"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualRouterBgpDampening,
				ConfigVariables: map[string]config.Variable{
					"name_suffix":   config.StringVariable(nameSuffix),
					"router_name":   config.StringVariable(resName),
					"template_name": config.StringVariable(templateName),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_virtual_router."+resName,
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-%s", resName, nameSuffix)),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router."+resName,
						tfjsonpath.New("protocol").AtMapKey("bgp").AtMapKey("enable"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router."+resName,
						tfjsonpath.New("protocol").AtMapKey("bgp").AtMapKey("dampening_profile").AtSliceIndex(0).AtMapKey("name"),
						knownvalue.StringExact("damp-profile-1"),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router."+resName,
						tfjsonpath.New("protocol").AtMapKey("bgp").AtMapKey("dampening_profile").AtSliceIndex(0).AtMapKey("enable"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router."+resName,
						tfjsonpath.New("protocol").AtMapKey("bgp").AtMapKey("dampening_profile").AtSliceIndex(0).AtMapKey("cutoff"),
						knownvalue.Float64Exact(2.0),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router."+resName,
						tfjsonpath.New("protocol").AtMapKey("bgp").AtMapKey("dampening_profile").AtSliceIndex(0).AtMapKey("reuse"),
						knownvalue.Float64Exact(0.75),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router."+resName,
						tfjsonpath.New("protocol").AtMapKey("bgp").AtMapKey("dampening_profile").AtSliceIndex(0).AtMapKey("max_hold_time"),
						knownvalue.Int64Exact(1800),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router."+resName,
						tfjsonpath.New("protocol").AtMapKey("bgp").AtMapKey("dampening_profile").AtSliceIndex(0).AtMapKey("decay_half_life_reachable"),
						knownvalue.Int64Exact(600),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router."+resName,
						tfjsonpath.New("protocol").AtMapKey("bgp").AtMapKey("dampening_profile").AtSliceIndex(0).AtMapKey("decay_half_life_unreachable"),
						knownvalue.Int64Exact(1200),
					),
				},
			},
		},
	})
}

const testAccVirtualRouterBgpDampening = `
variable "name_suffix" { type = string }
variable "router_name" { type = string }
variable "template_name" { type = string }

resource "panos_template" "template" {
  name = "${var.template_name}-${var.name_suffix}"

  location = {
    panorama = {
      panorama_device = "localhost.localdomain"
    }
  }
}

resource "panos_virtual_router" "vr_bgp_dampening" {
  location = {
    template = {
      name = panos_template.template.name
    }
  }

  name = "${var.router_name}-${var.name_suffix}"

  protocol = {
    bgp = {
      enable = true
      router_id = "10.0.0.2"
      local_as = "65002"
      dampening_profile = [
        {
          name = "damp-profile-1"
          enable = true
          cutoff = 2.0
          reuse = 0.75
          max_hold_time = 1800
          decay_half_life_reachable = 600
          decay_half_life_unreachable = 1200
        }
      ]
    }
  }
}
`

func TestAccPanosVirtualRouter_Bgp_AuthProfile(t *testing.T) {
	t.Parallel()

	resName := "vr_bgp_auth"
	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	templateName := "acc-vrouter-bgp-auth"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualRouterBgpAuthProfile,
				ConfigVariables: map[string]config.Variable{
					"name_suffix":   config.StringVariable(nameSuffix),
					"router_name":   config.StringVariable(resName),
					"template_name": config.StringVariable(templateName),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_virtual_router."+resName,
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-%s", resName, nameSuffix)),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router."+resName,
						tfjsonpath.New("protocol").AtMapKey("bgp").AtMapKey("auth_profile").AtSliceIndex(0).AtMapKey("name"),
						knownvalue.StringExact("auth-profile-1"),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router."+resName,
						tfjsonpath.New("protocol").AtMapKey("bgp").AtMapKey("auth_profile").AtSliceIndex(0).AtMapKey("secret"),
						knownvalue.StringExact("test-secret-key-123"),
					),
				},
			},
		},
	})
}

const testAccVirtualRouterBgpAuthProfile = `
variable "name_suffix" { type = string }
variable "router_name" { type = string }
variable "template_name" { type = string }

resource "panos_template" "template" {
  name = "${var.template_name}-${var.name_suffix}"

  location = {
    panorama = {
      panorama_device = "localhost.localdomain"
    }
  }
}

resource "panos_virtual_router" "vr_bgp_auth" {
  location = {
    template = {
      name = panos_template.template.name
    }
  }

  name = "${var.router_name}-${var.name_suffix}"

  protocol = {
    bgp = {
      enable = true
      router_id = "10.0.0.3"
      local_as = "65003"
      auth_profile = [
        {
          name = "auth-profile-1"
          secret = "test-secret-key-123"
        }
      ]
    }
  }
}
`

func TestAccPanosVirtualRouter_Bgp_AuthProfile_NoSecret(t *testing.T) {
	t.Parallel()

	resName := "vr_bgp_auth_nosecret"
	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	templateName := "acc-vrouter-bgp-auth-nosecret"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualRouterBgpAuthProfileNoSecret,
				ConfigVariables: map[string]config.Variable{
					"name_suffix":   config.StringVariable(nameSuffix),
					"router_name":   config.StringVariable(resName),
					"template_name": config.StringVariable(templateName),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_virtual_router."+resName,
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-%s", resName, nameSuffix)),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router."+resName,
						tfjsonpath.New("protocol").AtMapKey("bgp").AtMapKey("auth_profile").AtSliceIndex(0).AtMapKey("name"),
						knownvalue.StringExact("auth-profile-nosecret"),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router."+resName,
						tfjsonpath.New("protocol").AtMapKey("bgp").AtMapKey("auth_profile").AtSliceIndex(0).AtMapKey("secret"),
						knownvalue.Null(),
					),
				},
			},
		},
	})
}

const testAccVirtualRouterBgpAuthProfileNoSecret = `
variable "name_suffix" { type = string }
variable "router_name" { type = string }
variable "template_name" { type = string }

resource "panos_template" "template" {
  name = "${var.template_name}-${var.name_suffix}"

  location = {
    panorama = {
      panorama_device = "localhost.localdomain"
    }
  }
}

resource "panos_virtual_router" "vr_bgp_auth_nosecret" {
  location = {
    template = {
      name = panos_template.template.name
    }
  }

  name = "${var.router_name}-${var.name_suffix}"

  protocol = {
    bgp = {
      enable = true
      router_id = "10.0.0.10"
      local_as = "65010"
      auth_profile = [
        {
          name = "auth-profile-nosecret"
        }
      ]
    }
  }
}
`

func TestAccPanosVirtualRouter_Bgp_AuthProfile_MaxLength(t *testing.T) {
	t.Parallel()

	resName := "vr_bgp_auth_maxlen"
	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	templateName := "acc-vrouter-bgp-auth-maxlen"
	// Generate a 63-character secret (max allowed length)
	maxLengthSecret := acctest.RandStringFromCharSet(63, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualRouterBgpAuthProfileMaxLength,
				ConfigVariables: map[string]config.Variable{
					"name_suffix":       config.StringVariable(nameSuffix),
					"router_name":       config.StringVariable(resName),
					"template_name":     config.StringVariable(templateName),
					"max_length_secret": config.StringVariable(maxLengthSecret),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_virtual_router."+resName,
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-%s", resName, nameSuffix)),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router."+resName,
						tfjsonpath.New("protocol").AtMapKey("bgp").AtMapKey("auth_profile").AtSliceIndex(0).AtMapKey("name"),
						knownvalue.StringExact("auth-profile-maxlen"),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router."+resName,
						tfjsonpath.New("protocol").AtMapKey("bgp").AtMapKey("auth_profile").AtSliceIndex(0).AtMapKey("secret"),
						knownvalue.StringExact(maxLengthSecret),
					),
				},
			},
		},
	})
}

const testAccVirtualRouterBgpAuthProfileMaxLength = `
variable "name_suffix" { type = string }
variable "router_name" { type = string }
variable "template_name" { type = string }
variable "max_length_secret" { type = string }

resource "panos_template" "template" {
  name = "${var.template_name}-${var.name_suffix}"

  location = {
    panorama = {
      panorama_device = "localhost.localdomain"
    }
  }
}

resource "panos_virtual_router" "vr_bgp_auth_maxlen" {
  location = {
    template = {
      name = panos_template.template.name
    }
  }

  name = "${var.router_name}-${var.name_suffix}"

  protocol = {
    bgp = {
      enable = true
      router_id = "10.0.0.11"
      local_as = "65011"
      auth_profile = [
        {
          name = "auth-profile-maxlen"
          secret = var.max_length_secret
        }
      ]
    }
  }
}
`

func TestAccPanosVirtualRouter_Bgp_AuthProfile_SecretUpdate(t *testing.T) {
	t.Parallel()

	resName := "vr_bgp_auth_update"
	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	templateName := "acc-vrouter-bgp-auth-update"
	initialSecret := "initial-secret-123"
	updatedSecret := "updated-secret-456"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualRouterBgpAuthProfileSecretUpdate,
				ConfigVariables: map[string]config.Variable{
					"name_suffix":   config.StringVariable(nameSuffix),
					"router_name":   config.StringVariable(resName),
					"template_name": config.StringVariable(templateName),
					"secret_value":  config.StringVariable(initialSecret),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_virtual_router."+resName,
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-%s", resName, nameSuffix)),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router."+resName,
						tfjsonpath.New("protocol").AtMapKey("bgp").AtMapKey("auth_profile").AtSliceIndex(0).AtMapKey("secret"),
						knownvalue.StringExact(initialSecret),
					),
				},
			},
			{
				Config: testAccVirtualRouterBgpAuthProfileSecretUpdate,
				ConfigVariables: map[string]config.Variable{
					"name_suffix":   config.StringVariable(nameSuffix),
					"router_name":   config.StringVariable(resName),
					"template_name": config.StringVariable(templateName),
					"secret_value":  config.StringVariable(updatedSecret),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_virtual_router."+resName,
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-%s", resName, nameSuffix)),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router."+resName,
						tfjsonpath.New("protocol").AtMapKey("bgp").AtMapKey("auth_profile").AtSliceIndex(0).AtMapKey("secret"),
						knownvalue.StringExact(updatedSecret),
					),
				},
			},
		},
	})
}

const testAccVirtualRouterBgpAuthProfileSecretUpdate = `
variable "name_suffix" { type = string }
variable "router_name" { type = string }
variable "template_name" { type = string }
variable "secret_value" { type = string }

resource "panos_template" "template" {
  name = "${var.template_name}-${var.name_suffix}"

  location = {
    panorama = {
      panorama_device = "localhost.localdomain"
    }
  }
}

resource "panos_virtual_router" "vr_bgp_auth_update" {
  location = {
    template = {
      name = panos_template.template.name
    }
  }

  name = "${var.router_name}-${var.name_suffix}"

  protocol = {
    bgp = {
      enable = true
      router_id = "10.0.0.12"
      local_as = "65012"
      auth_profile = [
        {
          name = "auth-profile-update"
          secret = var.secret_value
        }
      ]
    }
  }
}
`

func makePanosVirtualRouterConfig(label string) string {
	configTpl := `
    variable "name_suffix" { type = string }
    variable "router_name" { type = string }
    variable "template_name" { type = string }

    variable "with_interface" {
        type = bool
        default = false
    }

    variable "ethernet_if_name" {
        type = string
        default = "ethernet1/40"
    }

    resource "panos_template" "template" {
        name = "${var.template_name}-${var.name_suffix}"

        location = {
            panorama = {
                panorama_device = "localhost.localdomain"
            }
        }
    }

    resource "panos_ethernet_interface" "ethernet" {
      count = var.with_interface ? 1 : 0
      location = {
        template = {
          vsys = "vsys1"
          name = panos_template.template.name
        }
      }


      name = "${var.ethernet_if_name}"

      layer3 = {
        lldp = {
          enable = true
        }
        mtu  = 1350
        ips  = [{ name = "1.1.1.1/32" }]

        ipv6 = {
          enabled = true
          addresses = [
            {
              advertise = {
                enable         = true
                valid_lifetime = "1000000"
              },
              name                = "::1",
              enable_on_interface = true
            },
          ]
        }
      }
    }

    resource "panos_virtual_router" "%s" {
      location = {
       template = {
          name = panos_template.template.name
        }
      }

      name = "${var.router_name}-${var.name_suffix}"

      interfaces = var.with_interface ? [panos_ethernet_interface.ethernet[0].name] : []

      administrative_distances = {
        static = 15
      }
    }
    `

	return fmt.Sprintf(configTpl, label)
}

// BGP Export AS-Path Variant Tests

func TestAccPanosVirtualRouter_BgpExport_AsPath_None(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPanosVirtualRouterBgpExportAsPathNone,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_virtual_router.test",
						tfjsonpath.New("protocol").
							AtMapKey("bgp").
							AtMapKey("policy").
							AtMapKey("export").
							AtMapKey("rules").
							AtSliceIndex(0).
							AtMapKey("action").
							AtMapKey("allow").
							AtMapKey("update").
							AtMapKey("as_path"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"none":               knownvalue.ObjectExact(map[string]knownvalue.Check{}),
							"remove":             knownvalue.Null(),
							"prepend":            knownvalue.Null(),
							"remove_and_prepend": knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

const testAccPanosVirtualRouterBgpExportAsPathNone = `
variable "prefix" { type = string }

resource "panos_template" "test" {
  location = { panorama = {} }
  name     = var.prefix
}

resource "panos_virtual_router" "test" {
  location = {
    template = {
      name = panos_template.test.name
    }
  }

  name = var.prefix

  protocol = {
    bgp = {
      router_id = "1.1.1.1"
      local_as  = 65001

      peer_group = [
        {
          name = "peer-group-1"
          enable = true
          type = {
            ebgp = {
              export_nexthop = "use-self"
            }
          }
          peer = [
            {
              name       = "peer-1"
              enable     = true
              local_ip   = "10.0.0.1"
              peer_ip    = "10.0.0.2"
              peer_as    = 65002
            }
          ]
        }
      ]

      policy = {
        export = {
          rules = [
            {
              name   = "export-rule-1"
              enable = true
              match = {
                as_path = {
                  regex = ".*"
                }
              }
              action = {
                allow = {
                  update = {
                    as_path = {
                      none = {}
                    }
                  }
                }
              }
            }
          ]
        }
      }
    }
  }
}
`

func TestAccPanosVirtualRouter_BgpExport_AsPath_RemoveAndPrepend(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPanosVirtualRouterBgpExportAsPathRemoveAndPrepend,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_virtual_router.test",
						tfjsonpath.New("protocol").
							AtMapKey("bgp").
							AtMapKey("policy").
							AtMapKey("export").
							AtMapKey("rules").
							AtSliceIndex(0).
							AtMapKey("action").
							AtMapKey("allow").
							AtMapKey("update").
							AtMapKey("as_path"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"none":               knownvalue.Null(),
							"remove":             knownvalue.Null(),
							"prepend":            knownvalue.Null(),
							"remove_and_prepend": knownvalue.Int64Exact(3),
						}),
					),
				},
			},
		},
	})
}

const testAccPanosVirtualRouterBgpExportAsPathRemoveAndPrepend = `
variable "prefix" { type = string }

resource "panos_template" "test" {
  location = { panorama = {} }
  name     = var.prefix
}

resource "panos_virtual_router" "test" {
  location = {
    template = {
      name = panos_template.test.name
    }
  }

  name = var.prefix

  protocol = {
    bgp = {
      router_id = "1.1.1.1"
      local_as  = 65001

      peer_group = [
        {
          name = "peer-group-1"
          enable = true
          type = {
            ebgp = {
              export_nexthop = "use-self"
            }
          }
          peer = [
            {
              name       = "peer-1"
              enable     = true
              local_ip   = "10.0.0.1"
              peer_ip    = "10.0.0.2"
              peer_as    = 65002
            }
          ]
        }
      ]

      policy = {
        export = {
          rules = [
            {
              name   = "export-rule-1"
              enable = true
              match = {
                as_path = {
                  regex = ".*"
                }
              }
              action = {
                allow = {
                  update = {
                    as_path = {
                      remove_and_prepend = 3
                    }
                  }
                }
              }
            }
          ]
        }
      }
    }
  }
}
`

// BGP Aggregation AS-Path Variant Tests

func TestAccPanosVirtualRouter_BgpAggregation_AsPath_None(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPanosVirtualRouterBgpAggregationAsPathNone,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_virtual_router.test",
						tfjsonpath.New("protocol").
							AtMapKey("bgp").
							AtMapKey("policy").
							AtMapKey("aggregation").
							AtMapKey("address").
							AtSliceIndex(0).
							AtMapKey("aggregate_route_attributes").
							AtMapKey("as_path"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"none":    knownvalue.ObjectExact(map[string]knownvalue.Check{}),
							"prepend": knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

const testAccPanosVirtualRouterBgpAggregationAsPathNone = `
variable "prefix" { type = string }

resource "panos_template" "test" {
  location = { panorama = {} }
  name     = var.prefix
}

resource "panos_virtual_router" "test" {
  location = {
    template = {
      name = panos_template.test.name
    }
  }

  name = var.prefix

  protocol = {
    bgp = {
      router_id = "1.1.1.1"
      local_as  = 65001

      policy = {
        aggregation = {
          address = [
            {
              name    = "agg-1"
              prefix  = "192.168.0.0/16"
              enable  = true
              summary = true

              aggregate_route_attributes = {
                as_path = {
                  none = {}
                }
              }
            }
          ]
        }
      }
    }
  }
}
`

func TestAccPanosVirtualRouter_BgpAggregation_AsPath_Prepend(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPanosVirtualRouterBgpAggregationAsPathPrepend,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_virtual_router.test",
						tfjsonpath.New("protocol").
							AtMapKey("bgp").
							AtMapKey("policy").
							AtMapKey("aggregation").
							AtMapKey("address").
							AtSliceIndex(0).
							AtMapKey("aggregate_route_attributes").
							AtMapKey("as_path"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"none":    knownvalue.Null(),
							"prepend": knownvalue.Int64Exact(2),
						}),
					),
				},
			},
		},
	})
}

const testAccPanosVirtualRouterBgpAggregationAsPathPrepend = `
variable "prefix" { type = string }

resource "panos_template" "test" {
  location = { panorama = {} }
  name     = var.prefix
}

resource "panos_virtual_router" "test" {
  location = {
    template = {
      name = panos_template.test.name
    }
  }

  name = var.prefix

  protocol = {
    bgp = {
      router_id = "1.1.1.1"
      local_as  = 65001

      policy = {
        aggregation = {
          address = [
            {
              name    = "agg-1"
              prefix  = "192.168.0.0/16"
              enable  = true
              summary = true

              aggregate_route_attributes = {
                as_path = {
                  prepend = 2
                }
              }
            }
          ]
        }
      }
    }
  }
}
`

// BGP Export Community Variant Tests

func TestAccPanosVirtualRouter_BgpExport_Community_RemoveAll(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPanosVirtualRouterBgpExportCommunityRemoveAll,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_virtual_router.test",
						tfjsonpath.New("protocol").
							AtMapKey("bgp").
							AtMapKey("policy").
							AtMapKey("export").
							AtMapKey("rules").
							AtSliceIndex(0).
							AtMapKey("action").
							AtMapKey("allow").
							AtMapKey("update").
							AtMapKey("community"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"none":         knownvalue.Null(),
							"remove_all":   knownvalue.ObjectExact(map[string]knownvalue.Check{}),
							"remove_regex": knownvalue.Null(),
							"append":       knownvalue.Null(),
							"overwrite":    knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

const testAccPanosVirtualRouterBgpExportCommunityRemoveAll = `
variable "prefix" { type = string }

resource "panos_template" "test" {
  location = { panorama = {} }
  name     = var.prefix
}

resource "panos_virtual_router" "test" {
  location = {
    template = {
      name = panos_template.test.name
    }
  }

  name = var.prefix

  protocol = {
    bgp = {
      router_id = "1.1.1.1"
      local_as  = 65001

      peer_group = [
        {
          name = "peer-group-1"
          enable = true
          type = {
            ebgp = {
              export_nexthop = "use-self"
            }
          }
          peer = [
            {
              name       = "peer-1"
              enable     = true
              local_ip   = "10.0.0.1"
              peer_ip    = "10.0.0.2"
              peer_as    = 65002
            }
          ]
        }
      ]

      policy = {
        export = {
          rules = [
            {
              name   = "export-rule-1"
              enable = true
              match = {
                as_path = {
                  regex = ".*"
                }
              }
              action = {
                allow = {
                  update = {
                    community = {
                      remove_all = {}
                    }
                  }
                }
              }
            }
          ]
        }
      }
    }
  }
}
`

func TestAccPanosVirtualRouter_BgpExport_Community_Append(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPanosVirtualRouterBgpExportCommunityAppend,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_virtual_router.test",
						tfjsonpath.New("protocol").
							AtMapKey("bgp").
							AtMapKey("policy").
							AtMapKey("export").
							AtMapKey("rules").
							AtSliceIndex(0).
							AtMapKey("action").
							AtMapKey("allow").
							AtMapKey("update").
							AtMapKey("community"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"none":         knownvalue.Null(),
							"remove_all":   knownvalue.Null(),
							"remove_regex": knownvalue.Null(),
							"append": knownvalue.ListExact([]knownvalue.Check{
								knownvalue.StringExact("65001:100"),
							}),
							"overwrite": knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

const testAccPanosVirtualRouterBgpExportCommunityAppend = `
variable "prefix" { type = string }

resource "panos_template" "test" {
  location = { panorama = {} }
  name     = var.prefix
}

resource "panos_virtual_router" "test" {
  location = {
    template = {
      name = panos_template.test.name
    }
  }

  name = var.prefix

  protocol = {
    bgp = {
      router_id = "1.1.1.1"
      local_as  = 65001

      peer_group = [
        {
          name = "peer-group-1"
          enable = true
          type = {
            ebgp = {
              export_nexthop = "use-self"
            }
          }
          peer = [
            {
              name       = "peer-1"
              enable     = true
              local_ip   = "10.0.0.1"
              peer_ip    = "10.0.0.2"
              peer_as    = 65002
            }
          ]
        }
      ]

      policy = {
        export = {
          rules = [
            {
              name   = "export-rule-1"
              enable = true
              match = {
                as_path = {
                  regex = ".*"
                }
              }
              action = {
                allow = {
                  update = {
                    community = {
                      append = ["65001:100"]
                    }
                  }
                }
              }
            }
          ]
        }
      }
    }
  }
}
`

// ECMP Algorithm Variant Tests

func TestAccPanosVirtualRouter_Ecmp_Algorithm_IpHash(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPanosVirtualRouterEcmpAlgorithmIpHash,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_virtual_router.test",
						tfjsonpath.New("ecmp").AtMapKey("algorithm").AtMapKey("ip_hash").AtMapKey("hash_seed"),
						knownvalue.Int64Exact(12345),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router.test",
						tfjsonpath.New("ecmp").AtMapKey("algorithm").AtMapKey("ip_hash").AtMapKey("src_only"),
						knownvalue.Bool(true),
					),
				},
			},
		},
	})
}

const testAccPanosVirtualRouterEcmpAlgorithmIpHash = `
variable "prefix" { type = string }

resource "panos_template" "test" {
  location = { panorama = {} }
  name     = var.prefix
}

resource "panos_virtual_router" "test" {
  location = {
    template = {
      name = panos_template.test.name
    }
  }

  name = var.prefix

  ecmp = {
    enable = true
    algorithm = {
      ip_hash = {
        hash_seed = 12345
        src_only  = true
      }
    }
  }
}
`

func TestAccPanosVirtualRouter_Ecmp_Algorithm_WeightedRoundRobin(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPanosVirtualRouterEcmpAlgorithmWeightedRoundRobin,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_virtual_router.test",
						tfjsonpath.New("ecmp").AtMapKey("algorithm"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"balanced_round_robin": knownvalue.Null(),
							"ip_hash":              knownvalue.Null(),
							"ip_modulo":            knownvalue.Null(),
							"weighted_round_robin": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"interface": knownvalue.ListExact([]knownvalue.Check{
									knownvalue.ObjectExact(map[string]knownvalue.Check{
										"name":   knownvalue.StringExact("ethernet1/1"),
										"weight": knownvalue.Int64Exact(150),
									}),
									knownvalue.ObjectExact(map[string]knownvalue.Check{
										"name":   knownvalue.StringExact("ethernet1/2"),
										"weight": knownvalue.Int64Exact(100),
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

const testAccPanosVirtualRouterEcmpAlgorithmWeightedRoundRobin = `
variable "prefix" { type = string }

resource "panos_template" "test" {
  location = { panorama = {} }
  name     = var.prefix
}

resource "panos_ethernet_interface" "test1" {
  location = {
    template = {
      vsys = "vsys1"
      name = panos_template.test.name
    }
  }

  name = "ethernet1/1"

  layer3 = {
    mtu = 1500
    ips = [{ name = "10.1.1.1/24" }]
  }
}

resource "panos_ethernet_interface" "test2" {
  location = {
    template = {
      vsys = "vsys1"
      name = panos_template.test.name
    }
  }

  name = "ethernet1/2"

  layer3 = {
    mtu = 1500
    ips = [{ name = "10.1.2.1/24" }]
  }
}

resource "panos_virtual_router" "test" {
  location = {
    template = {
      name = panos_template.test.name
    }
  }

  name = var.prefix

  interfaces = [
    panos_ethernet_interface.test1.name,
    panos_ethernet_interface.test2.name,
  ]

  ecmp = {
    enable = true
    algorithm = {
      weighted_round_robin = {
        interface = [
          {
            name   = "ethernet1/1"
            weight = 150
          },
          {
            name   = "ethernet1/2"
            weight = 100
          }
        ]
      }
    }
  }
}
`

// OSPF Tests

func TestAccPanosVirtualRouter_Ospf_Area_Normal(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualRouterOspfAreaNormal,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_virtual_router.test",
						tfjsonpath.New("protocol").AtMapKey("ospf").AtMapKey("enable"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router.test",
						tfjsonpath.New("protocol").AtMapKey("ospf").AtMapKey("router_id"),
						knownvalue.StringExact("10.0.1.1"),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router.test",
						tfjsonpath.New("protocol").AtMapKey("ospf").AtMapKey("area").AtSliceIndex(0).AtMapKey("name"),
						knownvalue.StringExact("0.0.0.0"),
					),
				},
			},
		},
	})
}

const testAccVirtualRouterOspfAreaNormal = `
variable "prefix" { type = string }

resource "panos_template" "test" {
  location = { panorama = {} }
  name     = var.prefix
}

resource "panos_ethernet_interface" "test" {
  location = {
    template = {
      vsys = "vsys1"
      name = panos_template.test.name
    }
  }

  name = "ethernet1/5"

  layer3 = {
    mtu = 1500
    ips = [{ name = "10.0.5.1/24" }]
  }
}

resource "panos_virtual_router" "test" {
  location = {
    template = {
      name = panos_template.test.name
    }
  }

  name = var.prefix

  interfaces = [panos_ethernet_interface.test.name]

  protocol = {
    ospf = {
      enable = true
      router_id = "10.0.1.1"
      reject_default_route = false
      allow_redist_default_route = true
      area = [
        {
          name = "0.0.0.0"
          type = {
            normal = {}
          }
          interface = [
            {
              name = "ethernet1/5"
              enable = true
              passive = false
              link_type = {
                broadcast = {
                  priority = 100
                  hello_interval = 10
                  dead_counts = 4
                }
              }
            }
          ]
        }
      ]
    }
  }
}
`

func TestAccPanosVirtualRouter_Ospf_Area_Stub(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualRouterOspfAreaStub,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_virtual_router.test",
						tfjsonpath.New("protocol").AtMapKey("ospf").AtMapKey("area").AtSliceIndex(0).AtMapKey("type").AtMapKey("stub"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}

const testAccVirtualRouterOspfAreaStub = `
variable "prefix" { type = string }

resource "panos_template" "test" {
  location = { panorama = {} }
  name     = var.prefix
}

resource "panos_virtual_router" "test" {
  location = {
    template = {
      name = panos_template.test.name
    }
  }

  name = var.prefix

  protocol = {
    ospf = {
      enable = true
      router_id = "10.0.1.2"
      area = [
        {
          name = "0.0.0.10"
          type = {
            stub = {
              accept_summary = true
              default_route = {
                advertise = {
                  metric = 20
                  type = "ext-2"
                }
              }
            }
          }
        }
      ]
    }
  }
}
`

func TestAccPanosVirtualRouter_Ospf_Area_Nssa(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualRouterOspfAreaNssa,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_virtual_router.test",
						tfjsonpath.New("protocol").AtMapKey("ospf").AtMapKey("area").AtSliceIndex(0).AtMapKey("type").AtMapKey("nssa"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}

const testAccVirtualRouterOspfAreaNssa = `
variable "prefix" { type = string }

resource "panos_template" "test" {
  location = { panorama = {} }
  name     = var.prefix
}

resource "panos_virtual_router" "test" {
  location = {
    template = {
      name = panos_template.test.name
    }
  }

  name = var.prefix

  protocol = {
    ospf = {
      enable = true
      router_id = "10.0.1.3"
      area = [
        {
          name = "0.0.0.20"
          type = {
            nssa = {
              accept_summary = false
              default_route = {
                advertise = {
                  metric = 15
                  type = "ext-1"
                }
              }
            }
          }
        }
      ]
    }
  }
}
`

func TestAccPanosVirtualRouter_Ospf_ExportRules(t *testing.T) {
	t.Skip("OSPF export_rules needs OSPF redistribution routing profile - requires implementation")
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualRouterOspfExportRules,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_virtual_router.test",
						tfjsonpath.New("protocol").AtMapKey("ospf").AtMapKey("export_rules"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}

const testAccVirtualRouterOspfExportRules = `
variable "prefix" { type = string }

resource "panos_template" "test" {
  location = { panorama = {} }
  name     = var.prefix
}

resource "panos_virtual_router" "test" {
  location = {
    template = {
      name = panos_template.test.name
    }
  }

  name = var.prefix

  protocol = {
    ospf = {
      enable = true
      router_id = "10.0.1.4"
      export_rules = [
        {
          name = "ospf-export-1"
          metric = 50
        }
      ]
    }
  }
}
`

func TestAccPanosVirtualRouter_Ospf_AuthProfile(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualRouterOspfAuthProfile,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_virtual_router.test",
						tfjsonpath.New("protocol").AtMapKey("ospf").AtMapKey("auth_profile"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}

const testAccVirtualRouterOspfAuthProfile = `
variable "prefix" { type = string }

resource "panos_template" "test" {
  location = { panorama = {} }
  name     = var.prefix
}

resource "panos_virtual_router" "test" {
  location = {
    template = {
      name = panos_template.test.name
    }
  }

  name = var.prefix

  protocol = {
    ospf = {
      enable = true
      router_id = "10.0.1.5"
      auth_profile = [
        {
          name = "ospf-auth-1"
          md5 = [
            {
              name = "1"
              key = "ospf-secret-key"
              preferred = true
            }
          ]
        }
      ]
    }
  }
}
`

// OSPFv3 Tests

func TestAccPanosVirtualRouter_Ospfv3_Basic(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualRouterOspfv3Basic,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_virtual_router.test",
						tfjsonpath.New("protocol").AtMapKey("ospfv3").AtMapKey("enable"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router.test",
						tfjsonpath.New("protocol").AtMapKey("ospfv3").AtMapKey("router_id"),
						knownvalue.StringExact("10.0.2.1"),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router.test",
						tfjsonpath.New("protocol").AtMapKey("ospfv3").AtMapKey("reject_default_route"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router.test",
						tfjsonpath.New("protocol").AtMapKey("ospfv3").AtMapKey("allow_redist_default_route"),
						knownvalue.Bool(false),
					),
				},
			},
		},
	})
}

const testAccVirtualRouterOspfv3Basic = `
variable "prefix" { type = string }

resource "panos_template" "test" {
  location = { panorama = {} }
  name     = var.prefix
}

resource "panos_ethernet_interface" "test" {
  location = {
    template = {
      vsys = "vsys1"
      name = panos_template.test.name
    }
  }

  name = "ethernet1/6"

  layer3 = {
    mtu = 1500
    ipv6 = {
      enabled = true
      addresses = [
        {
          name = "2001:db8::1/64"
          enable_on_interface = true
        }
      ]
    }
  }
}

resource "panos_virtual_router" "test" {
  location = {
    template = {
      name = panos_template.test.name
    }
  }

  name = var.prefix

  interfaces = [panos_ethernet_interface.test.name]

  protocol = {
    ospfv3 = {
      enable = true
      router_id = "10.0.2.1"
      reject_default_route = true
      allow_redist_default_route = false
      area = [
        {
          name = "0.0.0.0"
          type = {
            normal = {}
          }
          interface = [
            {
              name = "ethernet1/6"
              enable = true
              passive = true
            }
          ]
        }
      ]
    }
  }
}
`

func TestAccPanosVirtualRouter_Ospfv3_Areas(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualRouterOspfv3Areas,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_virtual_router.test",
						tfjsonpath.New("protocol").AtMapKey("ospfv3").AtMapKey("enable"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router.test",
						tfjsonpath.New("protocol").AtMapKey("ospfv3").AtMapKey("area").AtSliceIndex(0).AtMapKey("name"),
						knownvalue.StringExact("0.0.0.10"),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router.test",
						tfjsonpath.New("protocol").AtMapKey("ospfv3").AtMapKey("area").AtSliceIndex(0).AtMapKey("type").AtMapKey("stub"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}

const testAccVirtualRouterOspfv3Areas = `
variable "prefix" { type = string }

resource "panos_template" "test" {
  location = { panorama = {} }
  name     = var.prefix
}

resource "panos_virtual_router" "test" {
  location = {
    template = {
      name = panos_template.test.name
    }
  }

  name = var.prefix

  protocol = {
    ospfv3 = {
      enable = true
      router_id = "10.0.2.2"
      area = [
        {
          name = "0.0.0.10"
          type = {
            stub = {
              accept_summary = true
            }
          }
        }
      ]
    }
  }
}
`

// RIP Tests

func TestAccPanosVirtualRouter_Rip_Basic(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualRouterRipBasic,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_virtual_router.test",
						tfjsonpath.New("protocol").AtMapKey("rip").AtMapKey("enable"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router.test",
						tfjsonpath.New("protocol").AtMapKey("rip").AtMapKey("reject_default_route"),
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router.test",
						tfjsonpath.New("protocol").AtMapKey("rip").AtMapKey("allow_redist_default_route"),
						knownvalue.Bool(true),
					),
				},
			},
		},
	})
}

const testAccVirtualRouterRipBasic = `
variable "prefix" { type = string }

resource "panos_template" "test" {
  location = { panorama = {} }
  name     = var.prefix
}

resource "panos_ethernet_interface" "test" {
  location = {
    template = {
      vsys = "vsys1"
      name = panos_template.test.name
    }
  }

  name = "ethernet1/7"

  layer3 = {
    mtu = 1500
    ips = [{ name = "10.0.7.1/24" }]
  }
}

resource "panos_virtual_router" "test" {
  location = {
    template = {
      name = panos_template.test.name
    }
  }

  name = var.prefix

  interfaces = [panos_ethernet_interface.test.name]

  protocol = {
    rip = {
      enable = true
      reject_default_route = false
      allow_redist_default_route = true
      interface = [
        {
          name = "ethernet1/7"
          enable = true
        }
      ]
    }
  }
}
`

func TestAccPanosVirtualRouter_Rip_AuthProfile_Md5(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualRouterRipAuthProfileMd5,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_virtual_router.test",
						tfjsonpath.New("protocol").AtMapKey("rip").AtMapKey("enable"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router.test",
						tfjsonpath.New("protocol").AtMapKey("rip").AtMapKey("auth_profile").AtSliceIndex(0).AtMapKey("name"),
						knownvalue.StringExact("rip-auth-md5"),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router.test",
						tfjsonpath.New("protocol").AtMapKey("rip").AtMapKey("auth_profile").AtSliceIndex(0).AtMapKey("md5"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}

const testAccVirtualRouterRipAuthProfileMd5 = `
variable "prefix" { type = string }

resource "panos_template" "test" {
  location = { panorama = {} }
  name     = var.prefix
}

resource "panos_virtual_router" "test" {
  location = {
    template = {
      name = panos_template.test.name
    }
  }

  name = var.prefix

  protocol = {
    rip = {
      enable = true
      auth_profile = [
        {
          name = "rip-auth-md5"
          md5 = [
            {
              name = "1"
              key = "rip-secret-md5"
              preferred = true
            }
          ]
        }
      ]
    }
  }
}
`

func TestAccPanosVirtualRouter_Rip_AuthProfile_Text(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualRouterRipAuthProfileText,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_virtual_router.test",
						tfjsonpath.New("protocol").AtMapKey("rip").AtMapKey("enable"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router.test",
						tfjsonpath.New("protocol").AtMapKey("rip").AtMapKey("auth_profile").AtSliceIndex(0).AtMapKey("name"),
						knownvalue.StringExact("rip-auth-text"),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router.test",
						tfjsonpath.New("protocol").AtMapKey("rip").AtMapKey("auth_profile").AtSliceIndex(0).AtMapKey("password"),
						knownvalue.StringExact("rip-secret-text"),
					),
				},
			},
		},
	})
}

const testAccVirtualRouterRipAuthProfileText = `
variable "prefix" { type = string }

resource "panos_template" "test" {
  location = { panorama = {} }
  name     = var.prefix
}

resource "panos_virtual_router" "test" {
  location = {
    template = {
      name = panos_template.test.name
    }
  }

  name = var.prefix

  protocol = {
    rip = {
      enable = true
      auth_profile = [
        {
          name = "rip-auth-text"
          password = "rip-secret-text"
        }
      ]
    }
  }
}
`

func TestAccPanosVirtualRouter_Rip_ExportRules(t *testing.T) {
	t.Skip("RIP export_rules needs RIP redistribution routing profile - requires implementation")
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualRouterRipExportRules,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_virtual_router.test",
						tfjsonpath.New("protocol").AtMapKey("rip").AtMapKey("enable"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router.test",
						tfjsonpath.New("protocol").AtMapKey("rip").AtMapKey("export_rules").AtSliceIndex(0).AtMapKey("name"),
						knownvalue.StringExact("rip-export-1"),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router.test",
						tfjsonpath.New("protocol").AtMapKey("rip").AtMapKey("export_rules").AtSliceIndex(0).AtMapKey("metric"),
						knownvalue.Int64Exact(5),
					),
				},
			},
		},
	})
}

const testAccVirtualRouterRipExportRules = `
variable "prefix" { type = string }

resource "panos_template" "test" {
  location = { panorama = {} }
  name     = var.prefix
}

resource "panos_virtual_router" "test" {
  location = {
    template = {
      name = panos_template.test.name
    }
  }

  name = var.prefix

  protocol = {
    rip = {
      enable = true
      export_rules = [
        {
          name = "rip-export-1"
          metric = 5
        }
      ]
    }
  }
}
`

// Additional BGP Tests

func TestAccPanosVirtualRouter_Bgp_PeerGroup_Basic(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualRouterBgpPeerGroupBasic,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_virtual_router.test",
						tfjsonpath.New("protocol").AtMapKey("bgp").AtMapKey("peer_group").AtSliceIndex(0).AtMapKey("name"),
						knownvalue.StringExact("peer-group-basic"),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router.test",
						tfjsonpath.New("protocol").AtMapKey("bgp").AtMapKey("peer_group").AtSliceIndex(0).AtMapKey("enable"),
						knownvalue.Bool(true),
					),
				},
			},
		},
	})
}

const testAccVirtualRouterBgpPeerGroupBasic = `
variable "prefix" { type = string }

resource "panos_template" "test" {
  location = { panorama = {} }
  name     = var.prefix
}

resource "panos_virtual_router" "test" {
  location = {
    template = {
      name = panos_template.test.name
    }
  }

  name = var.prefix

  protocol = {
    bgp = {
      enable = true
      router_id = "10.1.1.1"
      local_as = "65100"
      peer_group = [
        {
          name = "peer-group-basic"
          enable = true
          type = {
            ebgp = {
              export_nexthop = "use-self"
            }
          }
          peer = [
            {
              name = "peer-1"
              enable = true
              local_ip = "10.1.1.1"
              peer_ip = "10.1.1.2"
              peer_as = "65200"
            }
          ]
        }
      ]
    }
  }
}
`

func TestAccPanosVirtualRouter_Bgp_PeerGroup_Advanced(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualRouterBgpPeerGroupAdvanced,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_virtual_router.test",
						tfjsonpath.New("protocol").AtMapKey("bgp").AtMapKey("peer_group").AtSliceIndex(0).AtMapKey("type").AtMapKey("ibgp"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}

const testAccVirtualRouterBgpPeerGroupAdvanced = `
variable "prefix" { type = string }

resource "panos_template" "test" {
  location = { panorama = {} }
  name     = var.prefix
}

resource "panos_virtual_router" "test" {
  location = {
    template = {
      name = panos_template.test.name
    }
  }

  name = var.prefix

  protocol = {
    bgp = {
      enable = true
      router_id = "10.1.2.1"
      local_as = "65300"
      peer_group = [
        {
          name = "peer-group-advanced"
          enable = true
          type = {
            ibgp = {}
          }
          peer = [
            {
              name = "ibgp-peer-1"
              enable = true
              local_ip = "10.1.2.1"
              peer_ip = "10.1.2.2"
              peer_as = "65300"
            }
          ]
        }
      ]
    }
  }
}
`

func TestAccPanosVirtualRouter_Bgp_Import(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualRouterBgpImport,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_virtual_router.test",
						tfjsonpath.New("protocol").AtMapKey("bgp").AtMapKey("policy").AtMapKey("import").AtMapKey("rules").AtSliceIndex(0).AtMapKey("name"),
						knownvalue.StringExact("import-rule-1"),
					),
				},
			},
		},
	})
}

const testAccVirtualRouterBgpImport = `
variable "prefix" { type = string }

resource "panos_template" "test" {
  location = { panorama = {} }
  name     = var.prefix
}

resource "panos_virtual_router" "test" {
  location = {
    template = {
      name = panos_template.test.name
    }
  }

  name = var.prefix

  protocol = {
    bgp = {
      enable = true
      router_id = "10.1.3.1"
      local_as = "65400"
      policy = {
        import = {
          rules = [
            {
              name = "import-rule-1"
              enable = true
              match = {
                as_path = {
                  regex = ".*65500.*"
                }
              }
              action = {
                allow = {
                  update = {
                    local_preference = 200
                  }
                }
              }
            }
          ]
        }
      }
    }
  }
}
`

func TestAccPanosVirtualRouter_Bgp_Redist(t *testing.T) {
	t.Skip("BGP redist_rules validation issue - profile reference not working as expected")
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualRouterBgpRedist,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_virtual_router.test",
						tfjsonpath.New("protocol").AtMapKey("bgp").AtMapKey("redist_rules").AtSliceIndex(0).AtMapKey("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-redist-profile", prefix)),
					),
				},
			},
		},
	})
}

const testAccVirtualRouterBgpRedist = `
variable "prefix" { type = string }

resource "panos_template" "test" {
  location = { panorama = {} }
  name     = var.prefix
}

resource "panos_bgp_redistribution_routing_profile" "test" {
  location = {
    template = {
      name = panos_template.test.name
    }
  }

  name = "${var.prefix}-redist-profile"
  ipv4 = {
    unicast = {
      static = {
        enable = true
        metric = 100
      }
    }
  }
}

resource "panos_virtual_router" "test" {
  location = {
    template = {
      name = panos_template.test.name
    }
  }

  name = var.prefix

  protocol = {
    bgp = {
      enable = true
      router_id = "10.1.4.1"
      local_as = "65500"
      redist_rules = [
        {
          name = panos_bgp_redistribution_routing_profile.test.name
          enable = true
          address_family_identifier = "ipv4"
          route_table = "unicast"
        }
      ]
    }
  }
}
`

func TestAccPanosVirtualRouter_Bgp_Community_RemoveRegex(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualRouterBgpCommunityRemoveRegex,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_virtual_router.test",
						tfjsonpath.New("protocol").AtMapKey("bgp").AtMapKey("policy").AtMapKey("export").AtMapKey("rules").AtSliceIndex(0).AtMapKey("action").AtMapKey("allow").AtMapKey("update").AtMapKey("community").AtMapKey("remove_regex"),
						knownvalue.StringExact("65001:.*"),
					),
				},
			},
		},
	})
}

const testAccVirtualRouterBgpCommunityRemoveRegex = `
variable "prefix" { type = string }

resource "panos_template" "test" {
  location = { panorama = {} }
  name     = var.prefix
}

resource "panos_virtual_router" "test" {
  location = {
    template = {
      name = panos_template.test.name
    }
  }

  name = var.prefix

  protocol = {
    bgp = {
      enable = true
      router_id = "10.1.5.1"
      local_as = "65600"
      peer_group = [
        {
          name = "pg1"
          enable = true
          type = { ebgp = { export_nexthop = "use-self" } }
          peer = [{
            name = "p1"
            enable = true
            local_ip = "10.1.5.1"
            peer_ip = "10.1.5.2"
            peer_as = "65700"
          }]
        }
      ]
      policy = {
        export = {
          rules = [
            {
              name = "export-remove-regex"
              enable = true
              match = {
                as_path = { regex = ".*" }
              }
              action = {
                allow = {
                  update = {
                    community = {
                      remove_regex = "65001:.*"
                    }
                  }
                }
              }
            }
          ]
        }
      }
    }
  }
}
`

func TestAccPanosVirtualRouter_Bgp_Community_Overwrite(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualRouterBgpCommunityOverwrite,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_virtual_router.test",
						tfjsonpath.New("protocol").AtMapKey("bgp").AtMapKey("policy").AtMapKey("export").AtMapKey("rules").AtSliceIndex(0).AtMapKey("action").AtMapKey("allow").AtMapKey("update").AtMapKey("community").AtMapKey("overwrite"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}

const testAccVirtualRouterBgpCommunityOverwrite = `
variable "prefix" { type = string }

resource "panos_template" "test" {
  location = { panorama = {} }
  name     = var.prefix
}

resource "panos_virtual_router" "test" {
  location = {
    template = {
      name = panos_template.test.name
    }
  }

  name = var.prefix

  protocol = {
    bgp = {
      enable = true
      router_id = "10.1.6.1"
      local_as = "65700"
      peer_group = [
        {
          name = "pg1"
          enable = true
          type = { ebgp = { export_nexthop = "use-self" } }
          peer = [{
            name = "p1"
            enable = true
            local_ip = "10.1.6.1"
            peer_ip = "10.1.6.2"
            peer_as = "65800"
          }]
        }
      ]
      policy = {
        export = {
          rules = [
            {
              name = "export-overwrite"
              enable = true
              match = {
                as_path = { regex = ".*" }
              }
              action = {
                allow = {
                  update = {
                    community = {
                      overwrite = ["65700:200"]
                    }
                  }
                }
              }
            }
          ]
        }
      }
    }
  }
}
`

func TestAccPanosVirtualRouter_Bgp_ExtendedCommunity(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualRouterBgpExtendedCommunity,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_virtual_router.test",
						tfjsonpath.New("protocol").AtMapKey("bgp").AtMapKey("policy").AtMapKey("export").AtMapKey("rules").AtSliceIndex(0).AtMapKey("action").AtMapKey("allow").AtMapKey("update").AtMapKey("extended_community"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}

const testAccVirtualRouterBgpExtendedCommunity = `
variable "prefix" { type = string }

resource "panos_template" "test" {
  location = { panorama = {} }
  name     = var.prefix
}

resource "panos_virtual_router" "test" {
  location = {
    template = {
      name = panos_template.test.name
    }
  }

  name = var.prefix

  protocol = {
    bgp = {
      enable = true
      router_id = "10.1.7.1"
      local_as = "65800"
      peer_group = [
        {
          name = "pg1"
          enable = true
          type = { ebgp = { export_nexthop = "use-self" } }
          peer = [{
            name = "p1"
            enable = true
            local_ip = "10.1.7.1"
            peer_ip = "10.1.7.2"
            peer_as = "65900"
          }]
        }
      ]
      policy = {
        export = {
          rules = [
            {
              name = "export-ext-comm"
              enable = true
              match = {
                as_path = { regex = ".*" }
              }
              action = {
                allow = {
                  update = {
                    extended_community = {
                      remove_all = {}
                    }
                  }
                }
              }
            }
          ]
        }
      }
    }
  }
}
`

func TestAccPanosVirtualRouter_Bgp_AggregationComplete(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualRouterBgpAggregationComplete,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_virtual_router.test",
						tfjsonpath.New("protocol").AtMapKey("bgp").AtMapKey("policy").AtMapKey("aggregation").AtMapKey("address"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}

const testAccVirtualRouterBgpAggregationComplete = `
variable "prefix" { type = string }

resource "panos_template" "test" {
  location = { panorama = {} }
  name     = var.prefix
}

resource "panos_virtual_router" "test" {
  location = {
    template = {
      name = panos_template.test.name
    }
  }

  name = var.prefix

  protocol = {
    bgp = {
      enable = true
      router_id = "10.1.8.1"
      local_as = "65900"
      policy = {
        aggregation = {
          address = [
            {
              name = "agg-complete"
              prefix = "192.168.0.0/16"
              enable = true
              summary = true
              aggregate_route_attributes = {
                origin = "incomplete"
                med = 100
                as_path = {
                  prepend = 3
                }
                community = {
                  append = ["65900:100"]
                }
              }
            }
          ]
        }
      }
    }
  }
}
`

func TestAccPanosVirtualRouter_Bgp_ExportComplete(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualRouterBgpExportComplete,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_virtual_router.test",
						tfjsonpath.New("protocol").AtMapKey("bgp").AtMapKey("policy").AtMapKey("export").AtMapKey("rules"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}

const testAccVirtualRouterBgpExportComplete = `
variable "prefix" { type = string }

resource "panos_template" "test" {
  location = { panorama = {} }
  name     = var.prefix
}

resource "panos_virtual_router" "test" {
  location = {
    template = {
      name = panos_template.test.name
    }
  }

  name = var.prefix

  protocol = {
    bgp = {
      enable = true
      router_id = "10.1.9.1"
      local_as = "66000"
      peer_group = [
        {
          name = "pg1"
          enable = true
          type = { ebgp = { export_nexthop = "use-self" } }
          peer = [{
            name = "p1"
            enable = true
            local_ip = "10.1.9.1"
            peer_ip = "10.1.9.2"
            peer_as = "66100"
          }]
        }
      ]
      policy = {
        export = {
          rules = [
            {
              name = "export-complete"
              enable = true
              match = {
                as_path = { regex = ".*" }
                community = { regex = ".*" }
              }
              action = {
                allow = {
                  update = {
                    origin = "igp"
                    med = 50
                    local_preference = 150
                    as_path = {
                      prepend = 2
                    }
                    community = {
                      append = ["66000:500"]
                    }
                  }
                }
              }
            }
          ]
        }
      }
    }
  }
}
`

func TestAccPanosVirtualRouter_Bgp_Complete(t *testing.T) {
	t.Skip("BGP redist_rules validation issue - profile reference not working as expected")
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualRouterBgpComplete,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_virtual_router.test",
						tfjsonpath.New("protocol").AtMapKey("bgp").AtMapKey("enable"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router.test",
						tfjsonpath.New("protocol").AtMapKey("bgp").AtMapKey("peer_group"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router.test",
						tfjsonpath.New("protocol").AtMapKey("bgp").AtMapKey("policy"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router.test",
						tfjsonpath.New("protocol").AtMapKey("bgp").AtMapKey("redist_rules").AtSliceIndex(0).AtMapKey("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-redist-profile", prefix)),
					),
				},
			},
		},
	})
}

const testAccVirtualRouterBgpComplete = `
variable "prefix" { type = string }

resource "panos_template" "test" {
  location = { panorama = {} }
  name     = var.prefix
}

resource "panos_bgp_redistribution_routing_profile" "test" {
  location = {
    template = {
      name = panos_template.test.name
    }
  }

  name = "${var.prefix}-redist-profile"
  ipv4 = {
    unicast = {
      connected = {
        enable = true
        metric = 50
      }
    }
  }
}

resource "panos_virtual_router" "test" {
  location = {
    template = {
      name = panos_template.test.name
    }
  }

  name = var.prefix

  protocol = {
    bgp = {
      enable = true
      router_id = "10.1.10.1"
      local_as = "66200"
      install_route = true
      reject_default_route = false
      allow_redist_default_route = true

      peer_group = [
        {
          name = "complete-pg"
          enable = true
          type = {
            ebgp = {
              export_nexthop = "use-self"
              import_nexthop = "original"
            }
          }
          peer = [
            {
              name = "complete-peer"
              enable = true
              local_ip = "10.1.10.1"
              peer_ip = "10.1.10.2"
              peer_as = "66300"
            }
          ]
        }
      ]

      policy = {
        import = {
          rules = [
            {
              name = "import-all"
              enable = true
              match = {
                as_path = { regex = ".*" }
              }
              action = {
                allow = {
                  update = {
                    local_preference = 100
                  }
                }
              }
            }
          ]
        }
        export = {
          rules = [
            {
              name = "export-all"
              enable = true
              match = {
                as_path = { regex = ".*" }
              }
              action = {
                allow = {
                  update = {
                    as_path = {
                      none = {}
                    }
                  }
                }
              }
            }
          ]
        }
      }

      redist_rules = [
        {
          name = panos_bgp_redistribution_routing_profile.test.name
          enable = true
          address_family_identifier = "ipv4"
          route_table = "unicast"
        }
      ]
    }
  }
}
`

// Enhanced Tests

func TestAccPanosVirtualRouter_Multicast_Complete(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualRouterMulticastComplete,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_virtual_router.test",
						tfjsonpath.New("multicast").AtMapKey("enable"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router.test",
						tfjsonpath.New("multicast").AtMapKey("interface_group"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}

const testAccVirtualRouterMulticastComplete = `
variable "prefix" { type = string }

resource "panos_template" "test" {
  location = { panorama = {} }
  name     = var.prefix
}

resource "panos_ethernet_interface" "test1" {
  location = {
    template = {
      vsys = "vsys1"
      name = panos_template.test.name
    }
  }

  name = "ethernet1/8"

  layer3 = {
    mtu = 1500
    ips = [{ name = "10.0.8.1/24" }]
  }
}

resource "panos_ethernet_interface" "test2" {
  location = {
    template = {
      vsys = "vsys1"
      name = panos_template.test.name
    }
  }

  name = "ethernet1/9"

  layer3 = {
    mtu = 1500
    ips = [{ name = "10.0.9.1/24" }]
  }
}

resource "panos_virtual_router" "test" {
  location = {
    template = {
      name = panos_template.test.name
    }
  }

  name = var.prefix

  interfaces = [
    panos_ethernet_interface.test1.name,
    panos_ethernet_interface.test2.name
  ]

  multicast = {
    enable = true
    interface_group = [
      {
        name = "mcast-grp-comp"
        description = "Complete multicast configuration"
        interfaces = ["ethernet1/8", "ethernet1/9"]
        group_permission = {
          any_source_multicast = [
            {
              name = "asm-complete-1"
              group_address = "224.0.0.0/4"
            },
            {
              name = "asm-complete-2"
              group_address = "239.0.0.0/8"
            }
          ]
        }
      }
    ]
  }
}
`

func TestAccPanosVirtualRouter_Ecmp_Complete(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualRouterEcmpComplete,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_virtual_router.test",
						tfjsonpath.New("ecmp").AtMapKey("enable"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router.test",
						tfjsonpath.New("ecmp").AtMapKey("max_paths"),
						knownvalue.Int64Exact(4),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router.test",
						tfjsonpath.New("ecmp").AtMapKey("symmetric_return"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router.test",
						tfjsonpath.New("ecmp").AtMapKey("strict_source_path"),
						knownvalue.Bool(true),
					),
				},
			},
		},
	})
}

const testAccVirtualRouterEcmpComplete = `
variable "prefix" { type = string }

resource "panos_template" "test" {
  location = { panorama = {} }
  name     = var.prefix
}

resource "panos_virtual_router" "test" {
  location = {
    template = {
      name = panos_template.test.name
    }
  }

  name = var.prefix

  ecmp = {
    enable = true
    max_paths = 4
    symmetric_return = true
    strict_source_path = true
    algorithm = {
      ip_hash = {
        hash_seed = 54321
        src_only = false
        use_port = true
      }
    }
  }
}
`

func TestAccPanosVirtualRouter_TemplateStack(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualRouterTemplateStack,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_virtual_router.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router.test",
						tfjsonpath.New("location").AtMapKey("template_stack").AtMapKey("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-stack", prefix)),
					),
				},
			},
		},
	})
}

const testAccVirtualRouterTemplateStack = `
variable "prefix" { type = string }

resource "panos_template" "test" {
  location = { panorama = {} }
  name     = "${var.prefix}-tpl"
}

resource "panos_template_stack" "test" {
  location = { panorama = {} }
  name     = "${var.prefix}-stack"
  templates = [panos_template.test.name]
}

resource "panos_virtual_router" "test" {
  location = {
    template_stack = {
      name = panos_template_stack.test.name
    }
  }

  name = var.prefix

  administrative_distances = {
    static = 10
  }
}
`
