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

func TestAccPasswordProfile_Basic(t *testing.T) {
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
				Config: passwordProfile_Basic_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_password_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_password_profile.example",
						tfjsonpath.New("password_change"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"expiration_period":                 knownvalue.Int64Exact(90),
							"expiration_warning_period":         knownvalue.Int64Exact(7),
							"post_expiration_admin_login_count": knownvalue.Int64Exact(3),
							"post_expiration_grace_period":      knownvalue.Int64Exact(5),
						}),
					),
				},
			},
		},
	})
}

const passwordProfile_Basic_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_password_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix

  password_change = {
    expiration_period = 90
    expiration_warning_period = 7
    post_expiration_admin_login_count = 3
    post_expiration_grace_period = 5
  }
}
`
