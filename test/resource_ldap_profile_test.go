package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	//"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccServerLdapProfile(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	location := config.ObjectVariable(map[string]config.Variable{
		"panorama": config.ObjectVariable(map[string]config.Variable{}),
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: panosServerLdapProfileTmpl1,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_ldap_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_ldap_profile.example",
						tfjsonpath.New("base"),
						knownvalue.StringExact("dc=example,dc=com"),
					),
					statecheck.ExpectKnownValue(
						"panos_ldap_profile.example",
						tfjsonpath.New("bind_dn"),
						knownvalue.StringExact("cn=admin,dc=example,dc=com"),
					),
					statecheck.ExpectKnownValue(
						"panos_ldap_profile.example",
						tfjsonpath.New("bind_password"),
						knownvalue.StringExact("admin_password"),
					),
					statecheck.ExpectKnownValue(
						"panos_ldap_profile.example",
						tfjsonpath.New("bind_timelimit"),
						knownvalue.Int64Exact(30),
					),
					statecheck.ExpectKnownValue(
						"panos_ldap_profile.example",
						tfjsonpath.New("disabled"),
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownValue(
						"panos_ldap_profile.example",
						tfjsonpath.New("ldap_type"),
						knownvalue.StringExact("active-directory"),
					),
					statecheck.ExpectKnownValue(
						"panos_ldap_profile.example",
						tfjsonpath.New("retry_interval"),
						knownvalue.Int64Exact(60),
					),
					statecheck.ExpectKnownValue(
						"panos_ldap_profile.example",
						tfjsonpath.New("ssl"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_ldap_profile.example",
						tfjsonpath.New("timelimit"),
						knownvalue.Int64Exact(30),
					),
					statecheck.ExpectKnownValue(
						"panos_ldap_profile.example",
						tfjsonpath.New("verify_server_certificate"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_ldap_profile.example",
						tfjsonpath.New("servers"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":    knownvalue.StringExact("ADSRV1"),
								"address": knownvalue.StringExact("ldap.example.com"),
								"port":    knownvalue.Int64Exact(389),
							}),
						}),
					),
				},
			},
		},
	})
}

const panosServerLdapProfileTmpl1 = `
variable "location" { type = any }
variable "prefix" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }

  name = format("%s-tmpl", var.prefix)
}

resource "panos_ldap_profile" "example" {
  location = var.location

  name = var.prefix

  base = "dc=example,dc=com"
  bind_dn = "cn=admin,dc=example,dc=com"
  bind_password = "admin_password"
  bind_timelimit = 30
  disabled = false
  ldap_type = "active-directory"
  retry_interval = 60
  ssl = true
  timelimit = 30
  verify_server_certificate = true

  servers = [
    {
      name = "ADSRV1"
      address = "ldap.example.com"
      port = 389
    }
  ]
}
`
