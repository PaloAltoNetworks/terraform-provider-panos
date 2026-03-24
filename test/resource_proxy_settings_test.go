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

func TestAccProxySettings_Basic(t *testing.T) {
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
				Config: proxySettingsConfig_Basic_step1,
				ConfigVariables: map[string]config.Variable{
					"location": location,
					"prefix":   config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_proxy_settings.settings",
						tfjsonpath.New("lcaas_use_proxy"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_proxy_settings.settings",
						tfjsonpath.New("secure_proxy_server"),
						knownvalue.StringExact("proxy.example.com"),
					),
					statecheck.ExpectKnownValue(
						"panos_proxy_settings.settings",
						tfjsonpath.New("secure_proxy_port"),
						knownvalue.Int64Exact(8080),
					),
					statecheck.ExpectKnownValue(
						"panos_proxy_settings.settings",
						tfjsonpath.New("secure_proxy_user"),
						knownvalue.StringExact("proxy-user"),
					),
					statecheck.ExpectKnownValue(
						"panos_proxy_settings.settings",
						tfjsonpath.New("secure_proxy_password"),
						knownvalue.StringExact("proxy-password"),
					),
				},
			},
			{
				Config: proxySettingsConfig_Basic_step2,
				ConfigVariables: map[string]config.Variable{
					"location": location,
					"prefix":   config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_proxy_settings.settings",
						tfjsonpath.New("lcaas_use_proxy"),
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownValue(
						"panos_proxy_settings.settings",
						tfjsonpath.New("secure_proxy_server"),
						knownvalue.StringExact("proxy2.example.com"),
					),
					statecheck.ExpectKnownValue(
						"panos_proxy_settings.settings",
						tfjsonpath.New("secure_proxy_port"),
						knownvalue.Int64Exact(8443),
					),
					statecheck.ExpectKnownValue(
						"panos_proxy_settings.settings",
						tfjsonpath.New("secure_proxy_user"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"panos_proxy_settings.settings",
						tfjsonpath.New("secure_proxy_password"),
						knownvalue.Null(),
					),
				},
			},
			{
				Config: proxySettingsConfig_Basic_step3,
				ConfigVariables: map[string]config.Variable{
					"location": location,
					"prefix":   config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_dns_settings.settings",
						tfjsonpath.New("dns_settings").AtMapKey("servers").AtMapKey("primary"),
						knownvalue.StringExact("10.0.0.1"),
					),
					statecheck.ExpectKnownValue(
						"panos_dns_settings.settings",
						tfjsonpath.New("dns_settings").AtMapKey("servers").AtMapKey("secondary"),
						knownvalue.StringExact("20.0.0.1"),
					),
				},
			},
		},
	})
}

const proxySettingsConfig_Basic_step1 = `
variable "location" { type = any }
variable "prefix" { type = string }

resource "panos_template" "template" {
  location = { panorama = {} }

  name = var.prefix
}

resource "panos_proxy_settings" "settings" {
  depends_on = [panos_template.template]
  location = var.location

  lcaas_use_proxy     = true
  secure_proxy_server = "proxy.example.com"
  secure_proxy_port   = 8080
  secure_proxy_user   = "proxy-user"
  secure_proxy_password = "proxy-password"
}
`

const proxySettingsConfig_Basic_step2 = `
variable "location" { type = any }
variable "prefix" { type = string }

resource "panos_template" "template" {
  location = { panorama = {} }

  name = var.prefix
}

resource "panos_dns_settings" "settings" {
  depends_on = [panos_template.template]
  location = var.location

  dns_settings = {
    servers = {
      primary = "10.0.0.1"
      secondary = "20.0.0.2"
    }
  }
}

resource "panos_proxy_settings" "settings" {
  depends_on = [panos_template.template, panos_dns_settings.settings]
  location = var.location

  lcaas_use_proxy     = false
  secure_proxy_server = "proxy2.example.com"
  secure_proxy_port   = 8443
}
`

const proxySettingsConfig_Basic_step3 = `
variable "location" { type = any }
variable "prefix" { type = string }

resource "panos_template" "template" {
  location = { panorama = {} }

  name = var.prefix
}

resource "panos_dns_settings" "settings" {
  depends_on = [panos_template.template]
  location = var.location

  dns_settings = {
    servers = {
      primary = "10.0.0.1"
      secondary = "20.0.0.1"
    }
  }
}
`
