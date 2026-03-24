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

func TestAccSyslogProfile(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: syslogProfileTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
					"servers": config.ListVariable(
						config.ObjectVariable(map[string]config.Variable{
							"name":      config.StringVariable("server1"),
							"server":    config.StringVariable("10.0.0.1"),
							"transport": config.StringVariable("UDP"),
							"port":      config.IntegerVariable(514),
							"facility":  config.StringVariable("LOG_USER"),
							"format":    config.StringVariable("IETF"),
						}),
					),
					"format": config.ObjectVariable(map[string]config.Variable{
						"traffic": config.StringVariable("traffic-fmt"),
					}),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_syslog_profile.profile",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_syslog_profile.profile",
						tfjsonpath.New("servers").AtSliceIndex(0).AtMapKey("name"),
						knownvalue.StringExact("server1"),
					),
					statecheck.ExpectKnownValue(
						"panos_syslog_profile.profile",
						tfjsonpath.New("format").AtMapKey("traffic"),
						knownvalue.StringExact("traffic-fmt"),
					),
				},
			},
			{
				Config: syslogProfileTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
					"servers": config.ListVariable(
						config.ObjectVariable(map[string]config.Variable{
							"name":      config.StringVariable("server1-upd"),
							"server":    config.StringVariable("10.0.0.1"),
							"transport": config.StringVariable("TCP"),
							"port":      config.IntegerVariable(514),
							"facility":  config.StringVariable("LOG_LOCAL0"),
							"format":    config.StringVariable("BSD"),
						}),
						config.ObjectVariable(map[string]config.Variable{
							"name":      config.StringVariable("server2"),
							"server":    config.StringVariable("10.0.0.2"),
							"transport": config.StringVariable("SSL"),
							"port":      config.IntegerVariable(6514),
							"facility":  config.StringVariable("LOG_LOCAL1"),
							"format":    config.StringVariable("IETF"),
						}),
					),
					"format": config.ObjectVariable(map[string]config.Variable{
						"traffic": config.StringVariable("traffic-fmt-upd"),
						"system":  config.StringVariable("system-fmt"),
						"escaping": config.ObjectVariable(map[string]config.Variable{
							"escape_character":   config.StringVariable("\\"),
							"escaped_characters": config.StringVariable(`'"`),
						}),
					}),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_syslog_profile.profile",
						tfjsonpath.New("servers").AtSliceIndex(0).AtMapKey("name"),
						knownvalue.StringExact("server1-upd"),
					),
					statecheck.ExpectKnownValue(
						"panos_syslog_profile.profile",
						tfjsonpath.New("servers").AtSliceIndex(0).AtMapKey("transport"),
						knownvalue.StringExact("TCP"),
					),
					statecheck.ExpectKnownValue(
						"panos_syslog_profile.profile",
						tfjsonpath.New("servers").AtSliceIndex(1).AtMapKey("name"),
						knownvalue.StringExact("server2"),
					),
					statecheck.ExpectKnownValue(
						"panos_syslog_profile.profile",
						tfjsonpath.New("format").AtMapKey("traffic"),
						knownvalue.StringExact("traffic-fmt-upd"),
					),
					statecheck.ExpectKnownValue(
						"panos_syslog_profile.profile",
						tfjsonpath.New("format").AtMapKey("system"),
						knownvalue.StringExact("system-fmt"),
					),
					statecheck.ExpectKnownValue(
						"panos_syslog_profile.profile",
						tfjsonpath.New("format").AtMapKey("escaping").AtMapKey("escape_character"),
						knownvalue.StringExact("\\"),
					),
				},
			},
		},
	})
}

const syslogProfileTmpl = `
variable "prefix" { type = string }
variable "servers" { type = any }
variable "format" { type = any }

resource "panos_template" "tmpl" {
  location = { panorama = {} }

  name = var.prefix
}

resource "panos_syslog_profile" "profile" {
  location = { template = { name = panos_template.tmpl.name } }

  name = var.prefix
  servers = var.servers
  format = var.format
}
`
