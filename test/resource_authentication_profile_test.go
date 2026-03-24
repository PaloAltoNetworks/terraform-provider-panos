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

func TestAccAuthenticationProfile_Basic(t *testing.T) {
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
				Config: panosAuthenticationProfile_Basic_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_authentication_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_authentication_profile.example",
						tfjsonpath.New("allow_list"),
						knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("all")}),
					),
					statecheck.ExpectKnownValue(
						"panos_authentication_profile.example",
						tfjsonpath.New("lockout"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"failed_attempts": knownvalue.Int64Exact(5),
							"lockout_time":    knownvalue.Int64Exact(30),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_authentication_profile.example",
						tfjsonpath.New("method"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"ldap": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"login_attribute": knownvalue.StringExact("sAMAccountName"),
								"passwd_exp_days": knownvalue.Int64Exact(14),
								"server_profile":  knownvalue.StringExact(prefix),
							}),
							"cloud":          knownvalue.Null(),
							"kerberos":       knownvalue.Null(),
							"local_database": knownvalue.Null(),
							"none":           knownvalue.Null(),
							"radius":         knownvalue.Null(),
							"saml_idp":       knownvalue.Null(),
							"tacplus":        knownvalue.Null(),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_authentication_profile.example",
						tfjsonpath.New("single_sign_on"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"kerberos_keytab":   knownvalue.Null(),
							"realm":             knownvalue.StringExact("EXAMPLE.COM"),
							"service_principal": knownvalue.StringExact("HTTP/firewall.example.com"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_authentication_profile.example",
						tfjsonpath.New("user_domain"),
						knownvalue.StringExact("example.com"),
					),
					statecheck.ExpectKnownValue(
						"panos_authentication_profile.example",
						tfjsonpath.New("username_modifier"),
						knownvalue.StringExact("%USERINPUT%@example.com"),
					),
				},
			},
		},
	})
}

const panosAuthenticationProfile_Basic_Tmpl = `
variable "location" { type = any }
variable "prefix" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }

  name = format("%s-tmpl", var.prefix)
}

resource "panos_ldap_profile" "example" {
  location = var.location

  name = var.prefix
}

resource "panos_authentication_profile" "example" {
  location = var.location

  name = var.prefix

  allow_list = ["all"]

  lockout = {
    failed_attempts = 5
    lockout_time    = 30
  }

  method = {
    ldap = {
      login_attribute = "sAMAccountName"
      passwd_exp_days = 14
      server_profile  = panos_ldap_profile.example.name
    }
  }

  single_sign_on = {
    #kerberos_keytab    = "BASE64_ENCODED_KEYTAB"
    realm              = "EXAMPLE.COM"
    service_principal  = "HTTP/firewall.example.com"
  }

  user_domain        = "example.com"
  username_modifier  = "%USERINPUT%@example.com"
}
`
