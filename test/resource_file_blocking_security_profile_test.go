
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

func TestAccFileBlockingSecurityProfile_Basic(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	location := config.ObjectVariable(map[string]config.Variable{
		"device_group": config.ObjectVariable(map[string]config.Variable{
			"name": config.StringVariable(prefix),
		}),
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fileBlockingSecurityProfile_Basic_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_file_blocking_security_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_file_blocking_security_profile.example",
						tfjsonpath.New("description"),
						knownvalue.StringExact("test description"),
					),
					statecheck.ExpectKnownValue(
						"panos_file_blocking_security_profile.example",
						tfjsonpath.New("disable_override"),
						knownvalue.StringExact("yes"),
					),
					statecheck.ExpectKnownValue(
						"panos_file_blocking_security_profile.example",
						tfjsonpath.New("rules"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name": knownvalue.StringExact("rule1"),
								"applications": knownvalue.ListExact([]knownvalue.Check{
									knownvalue.StringExact("paloalto-prisma-sdwan"),
								}),
								"file_types": knownvalue.ListExact([]knownvalue.Check{
									knownvalue.StringExact("pdf"),
								}),
								"direction": knownvalue.StringExact("both"),
								"action":    knownvalue.StringExact("block"),
							}),
						}),
					),
				},
			},
		},
	})
}

const fileBlockingSecurityProfile_Basic_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_device_group" "example" {
	location = { panorama = {} }
	name = var.prefix
}

resource "panos_file_blocking_security_profile" "example" {
	depends_on = [panos_device_group.example]
	location = var.location
	name = var.prefix
	description = "test description"
	disable_override = "yes"
	rules = [
		{
			name = "rule1"
			applications = ["paloalto-prisma-sdwan"]
			file_types = ["pdf"]
			direction = "both"
			action = "block"
		}
	]
}
`
