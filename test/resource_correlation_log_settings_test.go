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

func TestAccCorrelationLogSettings(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: correlationLogSettingsTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":      config.StringVariable(prefix),
					"description": config.StringVariable("test description"),
					"filter":      config.StringVariable("(severity eq high)"),
					"quarantine":  config.BoolVariable(false),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_correlation_log_settings.settings",
						tfjsonpath.New("description"),
						knownvalue.StringExact("test description"),
					),
					statecheck.ExpectKnownValue(
						"panos_correlation_log_settings.settings",
						tfjsonpath.New("filter"),
						knownvalue.StringExact("(severity eq high)"),
					),
					statecheck.ExpectKnownValue(
						"panos_correlation_log_settings.settings",
						tfjsonpath.New("quarantine"),
						knownvalue.Bool(false),
					),
				},
			},
			{
				Config: correlationLogSettingsTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":      config.StringVariable(prefix),
					"description": config.StringVariable("updated description"),
					"filter":      config.StringVariable("(severity eq critical)"),
					"quarantine":  config.BoolVariable(true),
					"actions": config.ListVariable(config.ObjectVariable(map[string]config.Variable{
						"name": config.StringVariable("integration-action"),
						"type": config.ObjectVariable(map[string]config.Variable{
							"integration": config.ObjectVariable(map[string]config.Variable{
								"action": config.StringVariable("Azure-Security-Center-Integration"),
							}),
						}),
					})),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_correlation_log_settings.settings",
						tfjsonpath.New("description"),
						knownvalue.StringExact("updated description"),
					),
					statecheck.ExpectKnownValue(
						"panos_correlation_log_settings.settings",
						tfjsonpath.New("filter"),
						knownvalue.StringExact("(severity eq critical)"),
					),
					statecheck.ExpectKnownValue(
						"panos_correlation_log_settings.settings",
						tfjsonpath.New("quarantine"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_correlation_log_settings.settings",
						tfjsonpath.New("actions").AtSliceIndex(0).AtMapKey("name"),
						knownvalue.StringExact("integration-action"),
					),
					statecheck.ExpectKnownValue(
						"panos_correlation_log_settings.settings",
						tfjsonpath.New("actions").AtSliceIndex(0).AtMapKey("type").AtMapKey("integration").AtMapKey("action"),
						knownvalue.StringExact("Azure-Security-Center-Integration"),
					),
				},
			},
		},
	})
}

const correlationLogSettingsTmpl = `
variable "prefix" { type = string }
variable "description" { type = string }
variable "filter" { type = string }
variable "quarantine" { type = bool }
variable "actions" {
  type    = any
  default = []
}


resource "panos_template" "tmpl" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_syslog_profile" "syslog1" {
  location = { template = { name = panos_template.tmpl.name } }

  name = "${var.prefix}1"

  servers = [{
    name = "server2"
    server = "10.0.0.2"
  }]
}

resource "panos_syslog_profile" "syslog2" {
  location = { template = { name = panos_template.tmpl.name } }

  name = "${var.prefix}2"

  servers = [{
    name = "server2"
    server = "10.0.0.2"
  }]
}

resource "panos_correlation_log_settings" "settings" {
  location = { template = { name = panos_template.tmpl.name } }
  name = var.prefix
  description = var.description
  filter = var.filter
  quarantine = var.quarantine
  syslog_profiles = [panos_syslog_profile.syslog1.name, panos_syslog_profile.syslog2.name]
  actions = var.actions
}
`
