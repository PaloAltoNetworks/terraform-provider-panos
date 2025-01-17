package provider_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	sdkErrors "github.com/PaloAltoNetworks/pango/errors"
	"github.com/PaloAltoNetworks/pango/panorama/devicegroup"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccDeviceGroup(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		CheckDestroy:             testAccDeviceGroupDestroy(prefix),
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

func testAccDeviceGroupDestroy(prefix string) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		api := devicegroup.NewService(sdkClient)
		ctx := context.TODO()

		location := devicegroup.NewPanoramaLocation()

		entries, err := api.List(ctx, *location, "get", "", "")
		if err != nil && !sdkErrors.IsObjectNotFound(err) {
			return fmt.Errorf("listing interface management entries via sdk: %v", err)
		}

		var leftEntries []string
		for _, elt := range entries {
			if strings.HasPrefix(elt.Name, prefix) {
				leftEntries = append(leftEntries, elt.Name)
			}
		}

		if len(leftEntries) > 0 {
			err := fmt.Errorf("terraform failed to remove entries from the server")
			delErr := api.Delete(ctx, *location, leftEntries...)
			if delErr != nil {
				return errors.Join(err, delErr)
			}
			return err
		}

		return nil
	}
}
