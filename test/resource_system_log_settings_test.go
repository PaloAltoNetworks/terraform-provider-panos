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

func TestAccSystemLogSettings(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: systemLogSettingsTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":           config.StringVariable(prefix),
					"description":      config.StringVariable("test description"),
					"filter":           config.StringVariable("(severity eq high)"),
					"send_to_panorama": config.BoolVariable(true),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_system_log_settings.settings",
						tfjsonpath.New("description"),
						knownvalue.StringExact("test description"),
					),
					statecheck.ExpectKnownValue(
						"panos_system_log_settings.settings",
						tfjsonpath.New("filter"),
						knownvalue.StringExact("(severity eq high)"),
					),
					statecheck.ExpectKnownValue(
						"panos_system_log_settings.settings",
						tfjsonpath.New("send_to_panorama"),
						knownvalue.Bool(true),
					),
				},
			},
			{
				Config: systemLogSettingsTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":           config.StringVariable(prefix),
					"description":      config.StringVariable("updated description"),
					"filter":           config.StringVariable("(severity eq critical)"),
					"send_to_panorama": config.BoolVariable(false),
					"actions": config.ListVariable(config.ObjectVariable(map[string]config.Variable{
						"name": config.StringVariable("azure-action"),
						"type": config.ObjectVariable(map[string]config.Variable{
							"integration": config.ObjectVariable(map[string]config.Variable{
								"action": config.StringVariable("Azure-Security-Center-Integration"),
							}),
						}),
					})),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_system_log_settings.settings",
						tfjsonpath.New("description"),
						knownvalue.StringExact("updated description"),
					),
					statecheck.ExpectKnownValue(
						"panos_system_log_settings.settings",
						tfjsonpath.New("filter"),
						knownvalue.StringExact("(severity eq critical)"),
					),
					statecheck.ExpectKnownValue(
						"panos_system_log_settings.settings",
						tfjsonpath.New("send_to_panorama"),
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownValue(
						"panos_system_log_settings.settings",
						tfjsonpath.New("actions").AtSliceIndex(0).AtMapKey("name"),
						knownvalue.StringExact("azure-action"),
					),
					statecheck.ExpectKnownValue(
						"panos_system_log_settings.settings",
						tfjsonpath.New("actions").AtSliceIndex(0).AtMapKey("type").AtMapKey("integration").AtMapKey("action"),
						knownvalue.StringExact("Azure-Security-Center-Integration"),
					),
				},
			},
		},
	})
}

const systemLogSettingsTmpl = `
variable "prefix" { type = string }
variable "description" { type = string }
variable "filter" { type = string }
variable "send_to_panorama" { type = bool }
variable "actions" {
  type = list(object({
    name = string
    type = object({
      integration = object({
        action = string
      })
    })
  }))
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

resource "panos_system_log_settings" "settings" {
  location = { template = { name = panos_template.tmpl.name } }
  name = var.prefix
  description = var.description
  filter = var.filter
  send_to_panorama = var.send_to_panorama
  syslog_profiles = [panos_syslog_profile.syslog1.name, panos_syslog_profile.syslog2.name]
  actions = var.actions
}
`
