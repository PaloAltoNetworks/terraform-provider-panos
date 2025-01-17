package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccDnsSettings(t *testing.T) {
	location := config.ObjectVariable(map[string]config.Variable{
		"system": config.ObjectVariable(map[string]config.Variable{}),
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dnsSettingsConfig1,
				ConfigVariables: map[string]config.Variable{
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_dns_settings.settings",
						tfjsonpath.New("fqdn_refresh_time"),
						knownvalue.Int64Exact(1800),
					),
					statecheck.ExpectKnownValue(
						"panos_dns_settings.settings",
						tfjsonpath.New("dns_settings").AtMapKey("servers").AtMapKey("primary"),
						knownvalue.StringExact("172.16.0.1"),
					),
					statecheck.ExpectKnownValue(
						"panos_dns_settings.settings",
						tfjsonpath.New("dns_settings").AtMapKey("servers").AtMapKey("secondary"),
						knownvalue.StringExact("172.16.0.2"),
					),
				},
			},
			{
				Config: dnsSettingsConfig2,
				ConfigVariables: map[string]config.Variable{
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_dns_settings.settings",
						tfjsonpath.New("fqdn_refresh_time"),
						knownvalue.Int64Exact(3600),
					),
					statecheck.ExpectKnownValue(
						"panos_dns_settings.settings",
						tfjsonpath.New("dns_settings").AtMapKey("servers").AtMapKey("primary"),
						knownvalue.StringExact("172.16.0.3"),
					),
				},
			},
			{
				Config: dnsSettingsConfig3,
				ConfigVariables: map[string]config.Variable{
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_dns_settings.settings",
						tfjsonpath.New("fqdn_refresh_time"),
						knownvalue.Int64Exact(1800),
					),
					statecheck.ExpectKnownValue(
						"panos_dns_settings.settings",
						tfjsonpath.New("dns_settings").AtMapKey("servers").AtMapKey("secondary"),
						knownvalue.StringExact("172.16.0.4"),
					),
				},
			},
		},
	})
}

const dnsSettingsConfig1 = `
variable "location" { type = map }

resource "panos_dns_settings" "settings" {
  location = var.location

  dns_settings = {
    servers = {
      primary = "172.16.0.1"
      secondary = "172.16.0.2"
    }
  }
}
`

const dnsSettingsConfig2 = `
variable "location" { type = map }

resource "panos_dns_settings" "settings" {
  location = var.location

  fqdn_refresh_time = 3600
  dns_settings = {
    servers = {
      primary = "172.16.0.3"
    }
  }
}
`

const dnsSettingsConfig3 = `
variable "location" { type = map }

resource "panos_dns_settings" "settings" {
  location = var.location

  dns_settings = {
    servers = {
      secondary = "172.16.0.4"
    }
  }
}
`
