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

func TestAccDnsProxy(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)
	name := "dns-proxy-test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dnsProxyResourceTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":  config.StringVariable(prefix),
					"name":    config.StringVariable(name),
					"enabled": config.BoolVariable(true),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_dns_proxy.proxy",
						tfjsonpath.New("name"),
						knownvalue.StringExact(name),
					),
					statecheck.ExpectKnownValue(
						"panos_dns_proxy.proxy",
						tfjsonpath.New("enabled"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_dns_proxy.proxy",
						tfjsonpath.New("default").AtMapKey("primary"),
						knownvalue.StringExact("8.8.8.8"),
					),
					statecheck.ExpectKnownValue(
						"panos_dns_proxy.proxy",
						tfjsonpath.New("default").AtMapKey("secondary"),
						knownvalue.StringExact("8.8.4.4"),
					),
					statecheck.ExpectKnownValue(
						"panos_dns_proxy.proxy",
						tfjsonpath.New("interface"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("ethernet1/1"),
							knownvalue.StringExact("ethernet1/2"),
						}),
					),
				},
			},
			{
				Config: dnsProxyResourceTmplFull,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
					"name":   config.StringVariable(name),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_dns_proxy.proxy",
						tfjsonpath.New("name"),
						knownvalue.StringExact(name),
					),
					statecheck.ExpectKnownValue(
						"panos_dns_proxy.proxy",
						tfjsonpath.New("enabled"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_dns_proxy.proxy",
						tfjsonpath.New("cache").AtMapKey("enabled"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_dns_proxy.proxy",
						tfjsonpath.New("cache").AtMapKey("cache_edns"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_dns_proxy.proxy",
						tfjsonpath.New("cache").AtMapKey("max_ttl").AtMapKey("enabled"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_dns_proxy.proxy",
						tfjsonpath.New("cache").AtMapKey("max_ttl").AtMapKey("time_to_live"),
						knownvalue.Int64Exact(3600),
					),
					statecheck.ExpectKnownValue(
						"panos_dns_proxy.proxy",
						tfjsonpath.New("domain_servers").AtSliceIndex(0).AtMapKey("name"),
						knownvalue.StringExact("internal-servers"),
					),
					statecheck.ExpectKnownValue(
						"panos_dns_proxy.proxy",
						tfjsonpath.New("domain_servers").AtSliceIndex(0).AtMapKey("domain_name"),
						knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("internal.example.com")}),
					),
					statecheck.ExpectKnownValue(
						"panos_dns_proxy.proxy",
						tfjsonpath.New("static_entries").AtSliceIndex(0).AtMapKey("name"),
						knownvalue.StringExact("static-entry-1"),
					),
					statecheck.ExpectKnownValue(
						"panos_dns_proxy.proxy",
						tfjsonpath.New("static_entries").AtSliceIndex(0).AtMapKey("domain"),
						knownvalue.StringExact("static.example.com"),
					),
					statecheck.ExpectKnownValue(
						"panos_dns_proxy.proxy",
						tfjsonpath.New("tcp_queries").AtMapKey("enabled"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_dns_proxy.proxy",
						tfjsonpath.New("udp_queries").AtMapKey("retries").AtMapKey("attempts"),
						knownvalue.Int64Exact(10),
					),
				},
			},
		},
	})
}

const dnsProxyResourceTmpl = `
variable "prefix" { type = string }
variable "name" { type = string }
variable "enabled" { type = bool }

resource "panos_template" "tmpl" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_ethernet_interface" "iface1" {
  location = { template = { name = panos_template.tmpl.name } }
  name = "ethernet1/1"
  layer3 = {}
}

resource "panos_ethernet_interface" "iface2" {
  location = { template = { name = panos_template.tmpl.name } }
  name = "ethernet1/2"
  layer3 = {}
}

resource "panos_dns_proxy" "proxy" {
  location = { template = { name = panos_template.tmpl.name } }
  name = var.name
  enabled = var.enabled
  default = {
    primary = "8.8.8.8"
    secondary = "8.8.4.4"
  }
  interface = [panos_ethernet_interface.iface1.name, panos_ethernet_interface.iface2.name]
}
`

const dnsProxyResourceTmplFull = `
variable "prefix" { type = string }
variable "name" { type = string }

resource "panos_template" "tmpl" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_ethernet_interface" "iface1" {
  location = { template = { name = panos_template.tmpl.name } }
  name = "ethernet1/1"
  layer3 = {}
}

resource "panos_ethernet_interface" "iface2" {
  location = { template = { name = panos_template.tmpl.name } }
  name = "ethernet1/2"
  layer3 = {}
}

resource "panos_dns_proxy" "proxy" {
  location = { template = { name = panos_template.tmpl.name } }
  name = var.name
  enabled = true
  interface = [panos_ethernet_interface.iface1.name, panos_ethernet_interface.iface2.name]

  cache = {
    enabled = true
    cache_edns = true
    max_ttl = {
      enabled = true
      time_to_live = 3600
    }
  }

  default = {
    primary = "8.8.8.8"
    secondary = "8.8.4.4"
  }

  domain_servers = [{
    name = "internal-servers"
    domain_name = ["internal.example.com"]
    primary = "10.0.0.1"
    secondary = "10.0.0.2"
    cacheable = true
  }]

  static_entries = [{
    name = "static-entry-1"
    domain = "static.example.com"
    address = ["192.168.1.100"]
  }]

  tcp_queries = {
    enabled = true
    max_pending_requests = 128
  }

  udp_queries = {
    retries = {
      attempts = 10
      interval = 5
    }
  }
}
`
