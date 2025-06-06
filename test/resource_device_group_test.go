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

func TestAccDeviceGroup(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceGroupResourceTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_device_group.dg",
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-dg", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_device_group.dg",
						tfjsonpath.New("description"),
						knownvalue.StringExact("description"),
					),
					statecheck.ExpectKnownValue(
						"panos_device_group.dg",
						tfjsonpath.New("templates"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact(fmt.Sprintf("%s-tmpl", prefix)),
						}),
					),
					// statecheck.ExpectKnownValue(
					// 	"panos_device_group.dg",
					// 	tfjsonpath.New("devices"),
					// 	knownvalue.ListExact([]knownvalue.Check{
					// 		knownvalue.MapExact(map[string]knownvalue.Check{
					// 			"name": knownvalue.StringExact("device-1"),
					// 			"vsys": knownvalue.StringExact("vsys1"),
					// 		}),
					// 	}),
					// ),
					statecheck.ExpectKnownValue(
						"panos_device_group.dg",
						tfjsonpath.New("authorization_code"),
						knownvalue.StringExact("code"),
					),
				},
			},
		},
	})
}

const testAccDeviceGroupResourceTmpl = `
variable "prefix" { type = string }

resource "panos_template" "template" {
  location = { panorama = {} }

  name = format("%s-tmpl", var.prefix)
}

resource "panos_device_group" "dg" {
  location = { panorama = {} }

  name        = format("%s-dg", var.prefix)
  description = "description"

  templates = [ resource.panos_template.template.name ]
  # devices   = [{ name = "device-1", vsys = ["vsys1"] }]

  authorization_code = "code"
}
`
