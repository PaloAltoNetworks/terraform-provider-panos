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

func TestAccDeviceGroupParent(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceGroupResourceParentTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_device_group_parent.relationship",
						tfjsonpath.New("device_group"),
						knownvalue.StringExact(fmt.Sprintf("%s-dg-child", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_device_group_parent.relationship",
						tfjsonpath.New("parent"),
						knownvalue.StringExact(fmt.Sprintf("%s-dg-parent", prefix)),
					),
				},
			},
		},
	})
}

const testAccDeviceGroupResourceParentTmpl = `
variable "prefix" { type = string }

resource "panos_device_group" "parent" {
  location = { panorama = {} }

  name        = format("%s-dg-parent", var.prefix)
}

resource "panos_device_group" "child" {
  location = { panorama = {} }

  name        = format("%s-dg-child", var.prefix)
}

resource "panos_device_group_parent" "relationship" {
  location = { panorama = {} }

  device_group = resource.panos_device_group.child.name
  parent       = resource.panos_device_group.parent.name
}
`
