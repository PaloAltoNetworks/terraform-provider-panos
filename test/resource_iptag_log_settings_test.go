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

func TestAccIptagLogSettings(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	filter1 := "(tag_name neq '')"
	filter2 := "(datasourcename eq test)"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: iptagLogSettingsTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":           config.StringVariable(prefix),
					"description":      config.StringVariable("test description"),
					"filter":           config.StringVariable(filter1),
					"send_to_panorama": config.BoolVariable(true),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_iptag_log_settings.settings",
						tfjsonpath.New("description"),
						knownvalue.StringExact("test description"),
					),
					statecheck.ExpectKnownValue(
						"panos_iptag_log_settings.settings",
						tfjsonpath.New("filter"),
						knownvalue.StringExact(filter1),
					),
					statecheck.ExpectKnownValue(
						"panos_iptag_log_settings.settings",
						tfjsonpath.New("send_to_panorama"),
						knownvalue.Bool(true),
					),
				},
			},
			{
				Config: iptagLogSettingsTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":           config.StringVariable(prefix),
					"description":      config.StringVariable("updated description"),
					"filter":           config.StringVariable(filter2),
					"send_to_panorama": config.BoolVariable(false),
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
						"panos_iptag_log_settings.settings",
						tfjsonpath.New("description"),
						knownvalue.StringExact("updated description"),
					),
					statecheck.ExpectKnownValue(
						"panos_iptag_log_settings.settings",
						tfjsonpath.New("filter"),
						knownvalue.StringExact(filter2),
					),
					statecheck.ExpectKnownValue(
						"panos_iptag_log_settings.settings",
						tfjsonpath.New("send_to_panorama"),
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownValue(
						"panos_iptag_log_settings.settings",
						tfjsonpath.New("actions").AtSliceIndex(0).AtMapKey("name"),
						knownvalue.StringExact("tag-action"),
					),
				},
			},
		},
	})
}

const iptagLogSettingsTmpl = `
variable "prefix" { type = string }
variable "description" { type = string }
variable "filter" { type = string }
variable "send_to_panorama" { type = bool }
variable "actions" {
  type    = any
  default = []
}


resource "panos_template" "tmpl" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_iptag_log_settings" "settings" {
  location = { template = { name = panos_template.tmpl.name } }
  name = var.prefix
  description = var.description
  filter = var.filter
  send_to_panorama = var.send_to_panorama
  actions = var.actions
}
`
