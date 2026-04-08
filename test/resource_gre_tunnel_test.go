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

func TestAccGreTunnel_Basic(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: greTunnelBasic,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_gre_tunnel.tunnel1",
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-tunnel1", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_gre_tunnel.tunnel1",
						tfjsonpath.New("copy_tos"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_gre_tunnel.tunnel1",
						tfjsonpath.New("disabled"),
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownValue(
						"panos_gre_tunnel.tunnel1",
						tfjsonpath.New("tunnel_interface"),
						knownvalue.StringExact("tunnel.1"),
					),
					statecheck.ExpectKnownValue(
						"panos_gre_tunnel.tunnel1",
						tfjsonpath.New("ttl"),
						knownvalue.Int64Exact(128),
					),
				},
			},
		},
	})
}

func TestAccGreTunnel_KeepAlive(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: greTunnelKeepAlive,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_gre_tunnel.tunnel1",
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-tunnel1", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_gre_tunnel.tunnel1",
						tfjsonpath.New("keep_alive"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"enable":     knownvalue.Bool(true),
							"interval":   knownvalue.Int64Exact(20),
							"retry":      knownvalue.Int64Exact(5),
							"hold_timer": knownvalue.Int64Exact(10),
						}),
					),
				},
			},
		},
	})
}

func TestAccGreTunnel_LocalAddressIp(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: greTunnelLocalAddressIp,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_gre_tunnel.tunnel1",
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-tunnel1", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_gre_tunnel.tunnel1",
						tfjsonpath.New("local_address"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"interface":   knownvalue.StringExact("ethernet1/1"),
							"ip":          knownvalue.StringExact("1.1.1.1/24"),
							"floating_ip": knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

func TestAccGreTunnel_PeerAddress(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: greTunnelPeerAddress,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_gre_tunnel.tunnel1",
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-tunnel1", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_gre_tunnel.tunnel1",
						tfjsonpath.New("peer_address"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"ip": knownvalue.StringExact("3.3.3.3"),
						}),
					),
				},
			},
		},
	})
}

func TestAccGreTunnel_Full(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: greTunnelFull,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_gre_tunnel.tunnel1",
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-tunnel1", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_gre_tunnel.tunnel1",
						tfjsonpath.New("copy_tos"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_gre_tunnel.tunnel1",
						tfjsonpath.New("disabled"),
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownValue(
						"panos_gre_tunnel.tunnel1",
						tfjsonpath.New("tunnel_interface"),
						knownvalue.StringExact("tunnel.2"),
					),
					statecheck.ExpectKnownValue(
						"panos_gre_tunnel.tunnel1",
						tfjsonpath.New("ttl"),
						knownvalue.Int64Exact(200),
					),
					statecheck.ExpectKnownValue(
						"panos_gre_tunnel.tunnel1",
						tfjsonpath.New("keep_alive"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"enable":     knownvalue.Bool(true),
							"interval":   knownvalue.Int64Exact(30),
							"retry":      knownvalue.Int64Exact(4),
							"hold_timer": knownvalue.Int64Exact(15),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_gre_tunnel.tunnel1",
						tfjsonpath.New("local_address"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"interface":   knownvalue.StringExact("ethernet1/2"),
							"ip":          knownvalue.StringExact("4.4.4.4/24"),
							"floating_ip": knownvalue.Null(),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_gre_tunnel.tunnel1",
						tfjsonpath.New("peer_address"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"ip": knownvalue.StringExact("5.5.5.5"),
						}),
					),
				},
			},
		},
	})
}

const greTunnelBasic = `
variable "prefix" { type = string }

locals {
  template_name = format("%s-tmpl", var.prefix)
}

resource "panos_template" "template" {
  location = { panorama = {} }
  name     = local.template_name
}

resource "panos_tunnel_interface" "tunnelif" {
  location = { template = { name = panos_template.template.name } }
  name = "tunnel.1"
}

resource "panos_gre_tunnel" "tunnel1" {
  location = { template = { name = panos_template.template.name } }

  name    = format("%s-tunnel1", var.prefix)
  copy_tos = true
  disabled = false
  tunnel_interface = panos_tunnel_interface.tunnelif.name
  ttl = 128
}
`

const greTunnelKeepAlive = `
variable "prefix" { type = string }

locals {
  template_name = format("%s-tmpl", var.prefix)
}

resource "panos_template" "template" {
  location = { panorama = {} }
  name     = local.template_name
}

resource "panos_gre_tunnel" "tunnel1" {
  location = { template = { name = panos_template.template.name } }

  name    = format("%s-tunnel1", var.prefix)
  keep_alive = {
    enable = true
    interval = 20
    retry = 5
    hold_timer = 10
  }
}
`

const greTunnelLocalAddressIp = `
variable "prefix" { type = string }

locals {
  template_name = format("%s-tmpl", var.prefix)
}

resource "panos_template" "template" {
  location = { panorama = {} }
  name     = local.template_name
}

resource "panos_ethernet_interface" "ethif" {
  location = { template = { name = panos_template.template.name } }
  name = "ethernet1/1"
  layer3 = {
    ips = [{
      name = "1.1.1.1/24"
    }]
  }
}

resource "panos_gre_tunnel" "tunnel1" {
  location = { template = { name = panos_template.template.name } }

  name    = format("%s-tunnel1", var.prefix)
  local_address = {
    interface = panos_ethernet_interface.ethif.name
    ip = panos_ethernet_interface.ethif.layer3.ips.0.name
  }
}
`

const greTunnelPeerAddress = `
variable "prefix" { type = string }

locals {
  template_name = format("%s-tmpl", var.prefix)
}

resource "panos_template" "template" {
  location = { panorama = {} }
  name     = local.template_name
}

resource "panos_gre_tunnel" "tunnel1" {
  location = { template = { name = panos_template.template.name } }

  name    = format("%s-tunnel1", var.prefix)
  peer_address = {
    ip = "3.3.3.3"
  }
}
`

const greTunnelFull = `
variable "prefix" { type = string }

locals {
  template_name = format("%s-tmpl", var.prefix)
}

resource "panos_template" "template" {
  location = { panorama = {} }
  name     = local.template_name
}

resource "panos_ethernet_interface" "ethif" {
  location = { template = { name = panos_template.template.name } }
  name = "ethernet1/2"
  layer3 = {
    ips = [{
      name = "4.4.4.4/24"
    }]
  }
}

resource "panos_tunnel_interface" "tunnelif" {
  location = { template = { name = panos_template.template.name } }
  name = "tunnel.2"
}

resource "panos_gre_tunnel" "tunnel1" {
  location = { template = { name = panos_template.template.name } }

  name    = format("%s-tunnel1", var.prefix)
  copy_tos = true
  disabled = false
  tunnel_interface = panos_tunnel_interface.tunnelif.name
  ttl = 200
  keep_alive = {
    enable = true
    interval = 30
    retry = 4
    hold_timer = 15
  }
  local_address = {
    interface = panos_ethernet_interface.ethif.name
    ip = "4.4.4.4/24"
  }
  peer_address = {
    ip = "5.5.5.5"
  }
}
`
