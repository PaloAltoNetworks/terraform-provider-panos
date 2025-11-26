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

func TestAccInterfaceManagementProfile_Basic(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: interfaceManagementProfile_BasicTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_interface_management_profile.profile",
						tfjsonpath.New("http"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_interface_management_profile.profile",
						tfjsonpath.New("https"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_interface_management_profile.profile",
						tfjsonpath.New("ping"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_interface_management_profile.profile",
						tfjsonpath.New("response_pages"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_interface_management_profile.profile",
						tfjsonpath.New("userid_service"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_interface_management_profile.profile",
						tfjsonpath.New("userid_syslog_listener_ssl"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_interface_management_profile.profile",
						tfjsonpath.New("userid_syslog_listener_udp"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_interface_management_profile.profile",
						tfjsonpath.New("ssh"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_interface_management_profile.profile",
						tfjsonpath.New("telnet"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_interface_management_profile.profile",
						tfjsonpath.New("snmp"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_interface_management_profile.profile",
						tfjsonpath.New("http_ocsp"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_interface_management_profile.profile",
						tfjsonpath.New("permitted_ips"),
						knownvalue.Null(),
					),
				},
			},
		},
	})
}

func TestAccInterfaceManagementProfile_PermittedIps(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: interfaceManagementProfile_PermittedIpsTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
					"permitted_ips": config.ListVariable(
						config.ObjectVariable(map[string]config.Variable{
							"name": config.StringVariable("172.16.0.1"),
						}),
						config.ObjectVariable(map[string]config.Variable{
							"name": config.StringVariable("172.16.0.2"),
						}),
					),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_interface_management_profile.profile",
						tfjsonpath.New("permitted_ips"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name": knownvalue.StringExact("172.16.0.1"),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name": knownvalue.StringExact("172.16.0.2"),
							}),
						}),
					),
				},
			},
		},
	})
}

const interfaceManagementProfile_BasicTmpl = `
variable "prefix" { type = string }

resource "panos_template" "template" {
  location = { panorama = {} }
  name = "${var.prefix}-tmpl"
}

resource "panos_interface_management_profile" "profile" {
  depends_on = [ resource.panos_template.template ]
  location = {
    template = {
      name = panos_template.template.name
    }
  }
  name = var.prefix

  http  = true
  https = true
  ping  = true

  response_pages = true

  userid_service             = true
  userid_syslog_listener_ssl = true
  userid_syslog_listener_udp = true

  ssh    = true
  telnet = true
  snmp   = true

  http_ocsp = true
}
`

const interfaceManagementProfile_PermittedIpsTmpl = `
variable "prefix" { type = string }
variable "permitted_ips" {
  type = list(map(string))
  default = []
}

resource "panos_template" "template" {
  location = { panorama = {} }
  name = "${var.prefix}-tmpl"
}

resource "panos_interface_management_profile" "profile" {
  depends_on = [ resource.panos_template.template ]
  location = {
    template = {
      name = panos_template.template.name
    }
  }
  name = var.prefix

  permitted_ips = var.permitted_ips
}
`
