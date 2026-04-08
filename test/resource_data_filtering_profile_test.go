
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

func TestAccDataFilteringProfile_Basic(t *testing.T) {
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
				Config: dataFilteringProfile_Basic_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_data_filtering_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_data_filtering_profile.example",
						tfjsonpath.New("data_capture"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_data_filtering_profile.example",
						tfjsonpath.New("description"),
						knownvalue.StringExact("test description"),
					),
					statecheck.ExpectKnownValue(
						"panos_data_filtering_profile.example",
						tfjsonpath.New("disable_override"),
						knownvalue.StringExact("yes"),
					),
					statecheck.ExpectKnownValue(
						"panos_data_filtering_profile.example",
						tfjsonpath.New("rules"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":          knownvalue.StringExact("rule1"),
								"data_object": knownvalue.StringExact(prefix),
								"direction":     knownvalue.StringExact("both"),
								"alert_threshold": knownvalue.Int64Exact(10),
								"block_threshold": knownvalue.Int64Exact(20),
								"log_severity":  knownvalue.StringExact("high"),
								"application": knownvalue.ListExact([]knownvalue.Check{
									knownvalue.StringExact("any"),
								}),
								"file_type": knownvalue.ListExact([]knownvalue.Check{
									knownvalue.StringExact("any"),
								}),
							}),
						}),
					),
				},
			},
		},
	})
}

const dataFilteringProfile_Basic_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_device_group" "example" {
	location = { panorama = {} }
	name = var.prefix
}

resource "panos_custom_data_object" "example" {
	location = var.location
	name = var.prefix
	description = "test custom data object"
	pattern_type = {
		regex = {
			pattern = [
				{
					name = "pattern1"
					regex = "test-regex"
				}
			]
		}
	}
}

resource "panos_data_filtering_profile" "example" {
	location = var.location
	name = var.prefix
	data_capture = true
	description = "test description"
	disable_override = "yes"
	rules = [
		{
			name = "rule1"
			data_object = panos_custom_data_object.example.name
			direction = "both"
			alert_threshold = 10
			block_threshold = 20
			log_severity = "high"
			application = ["any"]
			file_type = ["any"]
		}
	]
}
`
