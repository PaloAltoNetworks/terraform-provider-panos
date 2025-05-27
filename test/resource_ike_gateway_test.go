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

func TestAccIkeGateway(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: ikeGatewayConfigTmpl1,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_ike_gateway.example1",
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-gw1", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_ike_gateway.example1",
						tfjsonpath.New("comment"),
						knownvalue.StringExact("ike gateway comment"),
					),
					statecheck.ExpectKnownValue(
						"panos_ike_gateway.example1",
						tfjsonpath.New("disabled"),
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownValue(
						"panos_ike_gateway.example1",
						tfjsonpath.New("ipv6"),
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownValue(
						"panos_ike_gateway.example1",
						tfjsonpath.New("authentication"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"certificate": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"certificate_profile":          knownvalue.Null(),
								"local_certificate":            knownvalue.Null(),
								"allow_id_payload_mismatch":    knownvalue.Bool(true),
								"strict_validation_revocation": knownvalue.Bool(true),
								"use_management_as_source":     knownvalue.Null(),
							}),
							"pre_shared_key": knownvalue.Null(),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_ike_gateway.example1",
						tfjsonpath.New("local_address"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"interface":   knownvalue.StringExact("ethernet1/1"),
							"ip":          knownvalue.StringExact("10.0.0.1/32"),
							"floating_ip": knownvalue.Null(),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_ike_gateway.example1",
						tfjsonpath.New("peer_address"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"ip":      knownvalue.StringExact("10.10.0.1/32"),
							"dynamic": knownvalue.Null(),
							"fqdn":    knownvalue.Null(),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_ike_gateway.example1",
						tfjsonpath.New("peer_id"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"type":     knownvalue.StringExact("fqdn"),
							"id":       knownvalue.StringExact("example.com"),
							"matching": knownvalue.StringExact("exact"),
						}),
					),
					// example2 - local_address without ip/floating_ip
					statecheck.ExpectKnownValue(
						"panos_ike_gateway.example2",
						tfjsonpath.New("local_address"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"interface":   knownvalue.StringExact("ethernet1/1"),
							"ip":          knownvalue.Null(),
							"floating_ip": knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

const ikeGatewayConfigTmpl1 = `
variable "prefix" { type = string }

resource "panos_ike_gateway" "example1" {
  location = { template = { name = panos_template.example.name } }

  name    = format("%s-gw1", var.prefix)

  comment = "ike gateway comment"
  disabled = false
  ipv6 = false

  authentication = {
    certificate = {
      allow_id_payload_mismatch = true
      strict_validation_revocation = true
    }
  }

  local_address = {
    interface = panos_ethernet_interface.example.name
    ip = panos_ethernet_interface.example.layer3.ips.0.name
  }

  peer_address = {
    ip = "10.10.0.1/32"
  }

  peer_id = {
    type     = "fqdn"
    id       = "example.com"
    matching = "exact"
  }
}

resource "panos_ike_gateway" "example2" {
  location = { template = { name = panos_template.example.name } }

  name    = format("%s-gw2", var.prefix)

  local_address = {
    interface = panos_ethernet_interface.example.name
  }
}

resource "panos_ethernet_interface" "example" {
  location = { template = { name = panos_template.example.name, vsys = "vsys1" } }

  name = "ethernet1/1"

  layer3 = {
    ips = [{ name = "10.0.0.1/32" }]
  }
}

resource "panos_device_group" "example" {
  location = { panorama = {} }

  name = format("%s-dg", var.prefix)

  templates = [panos_template.example.name]
}

resource "panos_template" "example" {
   location = { panorama = {} }
   name     = format("%s-tmpl", var.prefix)
}
`
