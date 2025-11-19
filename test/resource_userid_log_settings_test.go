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

func TestAccUserIdLogSettings(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: userIdLogSettingsTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":           config.StringVariable(prefix),
					"description":      config.StringVariable("test description"),
					"filter":           config.StringVariable("(datasourcename eq test)"),
					"send_to_panorama": config.BoolVariable(true),
					"quarantine":       config.BoolVariable(false),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_userid_log_settings.settings",
						tfjsonpath.New("description"),
						knownvalue.StringExact("test description"),
					),
					statecheck.ExpectKnownValue(
						"panos_userid_log_settings.settings",
						tfjsonpath.New("filter"),
						knownvalue.StringExact("(datasourcename eq test)"),
					),
					statecheck.ExpectKnownValue(
						"panos_userid_log_settings.settings",
						tfjsonpath.New("send_to_panorama"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_userid_log_settings.settings",
						tfjsonpath.New("quarantine"),
						knownvalue.Bool(false),
					),
				},
			},
			{
				Config: userIdLogSettingsTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":           config.StringVariable(prefix),
					"description":      config.StringVariable("updated description"),
					"send_to_panorama": config.BoolVariable(false),
					"quarantine":       config.BoolVariable(true),
					"actions": config.ListVariable(config.ObjectVariable(map[string]config.Variable{
						"name": config.StringVariable("tag-action"),
						"type": config.ObjectVariable(map[string]config.Variable{
							"tagging": config.ObjectVariable(map[string]config.Variable{
								"action": config.StringVariable("add-tag"),
								"target": config.StringVariable("source-address"),
								"tags":   config.ListVariable(config.StringVariable("tag1"), config.StringVariable("tag2")),
							}),
						}),
					})),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_userid_log_settings.settings",
						tfjsonpath.New("description"),
						knownvalue.StringExact("updated description"),
					),
					statecheck.ExpectKnownValue(
						"panos_userid_log_settings.settings",
						tfjsonpath.New("filter"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"panos_userid_log_settings.settings",
						tfjsonpath.New("send_to_panorama"),
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownValue(
						"panos_userid_log_settings.settings",
						tfjsonpath.New("quarantine"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_userid_log_settings.settings",
						tfjsonpath.New("actions").AtSliceIndex(0).AtMapKey("name"),
						knownvalue.StringExact("tag-action"),
					),
					statecheck.ExpectKnownValue(
						"panos_userid_log_settings.settings",
						tfjsonpath.New("actions").AtSliceIndex(0).AtMapKey("type").AtMapKey("tagging").AtMapKey("action"),
						knownvalue.StringExact("add-tag"),
					),
					statecheck.ExpectKnownValue(
						"panos_userid_log_settings.settings",
						tfjsonpath.New("actions").AtSliceIndex(0).AtMapKey("type").AtMapKey("tagging").AtMapKey("target"),
						knownvalue.StringExact("source-address"),
					),
					statecheck.ExpectKnownValue(
						"panos_userid_log_settings.settings",
						tfjsonpath.New("actions").AtSliceIndex(0).AtMapKey("type").AtMapKey("tagging").AtMapKey("tags").AtSliceIndex(0),
						knownvalue.StringExact("tag1"),
					),
					statecheck.ExpectKnownValue(
						"panos_userid_log_settings.settings",
						tfjsonpath.New("actions").AtSliceIndex(0).AtMapKey("type").AtMapKey("tagging").AtMapKey("tags").AtSliceIndex(1),
						knownvalue.StringExact("tag2"),
					),
				},
			},
		},
	})
}

const userIdLogSettingsTmpl = `
variable "prefix" { type = string }
variable "description" { type = string }
variable "filter" {
  type = string
  default = null
}
variable "send_to_panorama" { type = bool }
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

resource "panos_userid_log_settings" "settings" {
  location = { template = { name = panos_template.tmpl.name } }
  name = var.prefix
  description = var.description
  filter = var.filter
  send_to_panorama = var.send_to_panorama
  syslog_profiles = [panos_syslog_profile.syslog1.name, panos_syslog_profile.syslog2.name]
  quarantine = var.quarantine
  actions = var.actions
}
`
