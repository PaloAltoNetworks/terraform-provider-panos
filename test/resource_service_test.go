package provider_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	sdkErrors "github.com/PaloAltoNetworks/pango/errors"
	"github.com/PaloAltoNetworks/pango/objects/service"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccService(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		CheckDestroy:             testAccServiceDestroy(prefix),
		Steps: []resource.TestStep{
			{
				Config: testAccServiceResourceTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
					"services": config.MapVariable(map[string]config.Variable{
						"svc1": config.ObjectVariable(map[string]config.Variable{
							"description": config.StringVariable("description"),
							"tags":        config.ListVariable(config.StringVariable(fmt.Sprintf("%s-tag", prefix))),
							"protocol": config.ObjectVariable(map[string]config.Variable{
								"tcp": config.ObjectVariable(map[string]config.Variable{
									"destination_port": config.StringVariable("8080-9000"),
									"source_port":      config.StringVariable("40000"),
								}),
							}),
						}),
						"svc2": config.ObjectVariable(map[string]config.Variable{
							"protocol": config.ObjectVariable(map[string]config.Variable{
								"tcp": config.ObjectVariable(map[string]config.Variable{
									"destination_port": config.StringVariable("443"),
									"source_port":      config.StringVariable("20000-30000"),
									"override": config.ObjectVariable(map[string]config.Variable{
										"timeout":           config.IntegerVariable(600),
										"halfclose_timeout": config.IntegerVariable(300),
										"timewait_timeout":  config.IntegerVariable(60),
									}),
								}),
							}),
						}),
						"svc3": config.ObjectVariable(map[string]config.Variable{
							"name": config.StringVariable(fmt.Sprintf("%s-svc3", prefix)),
							"protocol": config.ObjectVariable(map[string]config.Variable{
								"udp": config.ObjectVariable(map[string]config.Variable{
									"destination_port": config.StringVariable("443"),
									"source_port":      config.StringVariable("20000"),
									"override": config.ObjectVariable(map[string]config.Variable{
										"timeout": config.IntegerVariable(600),
									}),
								}),
							}),
						}),
						"svc4": config.ObjectVariable(map[string]config.Variable{
							"name": config.StringVariable(fmt.Sprintf("%s-svc4", prefix)),
							"protocol": config.ObjectVariable(map[string]config.Variable{
								"udp": config.ObjectVariable(map[string]config.Variable{
									"destination_port": config.StringVariable("443"),
									"source_port":      config.StringVariable("20000"),
								}),
							}),
						}),
					}),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_service.svc1",
						tfjsonpath.New("description"),
						knownvalue.StringExact("description"),
					),
					statecheck.ExpectKnownValue(
						"panos_service.svc1",
						tfjsonpath.New("tags"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact(fmt.Sprintf("%s-tag", prefix)),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_service.svc1",
						tfjsonpath.New("protocol").
							AtMapKey("tcp").
							AtMapKey("source_port"),
						knownvalue.StringExact("40000"),
					),
					statecheck.ExpectKnownValue(
						"panos_service.svc1",
						tfjsonpath.New("protocol").
							AtMapKey("tcp").
							AtMapKey("destination_port"),
						knownvalue.StringExact("8080-9000"),
					),
					statecheck.ExpectKnownValue(
						"panos_service.svc2",
						tfjsonpath.New("protocol").
							AtMapKey("tcp").
							AtMapKey("source_port"),
						knownvalue.StringExact("20000-30000"),
					),
					statecheck.ExpectKnownValue(
						"panos_service.svc2",
						tfjsonpath.New("protocol").
							AtMapKey("tcp").
							AtMapKey("destination_port"),
						knownvalue.StringExact("443"),
					),
					statecheck.ExpectKnownValue(
						"panos_service.svc2",
						tfjsonpath.New("protocol").
							AtMapKey("tcp").
							AtMapKey("override").
							AtMapKey("timeout"),
						knownvalue.Int64Exact(600),
					),
					statecheck.ExpectKnownValue(
						"panos_service.svc2",
						tfjsonpath.New("protocol").
							AtMapKey("tcp").
							AtMapKey("override").
							AtMapKey("halfclose_timeout"),
						knownvalue.Int64Exact(300),
					),
					statecheck.ExpectKnownValue(
						"panos_service.svc2",
						tfjsonpath.New("protocol").
							AtMapKey("tcp").
							AtMapKey("override").
							AtMapKey("timewait_timeout"),
						knownvalue.Int64Exact(60),
					),
					statecheck.ExpectKnownValue(
						"panos_service.svc3",
						tfjsonpath.New("protocol").
							AtMapKey("udp").
							AtMapKey("source_port"),
						knownvalue.StringExact("20000"),
					),
					statecheck.ExpectKnownValue(
						"panos_service.svc3",
						tfjsonpath.New("protocol").
							AtMapKey("udp").
							AtMapKey("destination_port"),
						knownvalue.StringExact("443"),
					),
					statecheck.ExpectKnownValue(
						"panos_service.svc3",
						tfjsonpath.New("protocol").
							AtMapKey("udp").
							AtMapKey("override").
							AtMapKey("timeout"),
						knownvalue.Int64Exact(600),
					),
					statecheck.ExpectKnownValue(
						"panos_service.svc4",
						tfjsonpath.New("protocol").
							AtMapKey("udp").
							AtMapKey("source_port"),
						knownvalue.StringExact("20000"),
					),
					statecheck.ExpectKnownValue(
						"panos_service.svc4",
						tfjsonpath.New("protocol").
							AtMapKey("udp").
							AtMapKey("destination_port"),
						knownvalue.StringExact("443"),
					),
					statecheck.ExpectKnownValue(
						"panos_service.svc4",
						tfjsonpath.New("protocol").
							AtMapKey("udp").
							AtMapKey("override"),
						knownvalue.Null(),
					),
				},
			},
		},
	})
}

const testAccServiceResourceTmpl = `
variable prefix { type = string }
variable "services" {
  type = map(object({
    description = optional(string),
    tags = optional(list(string)),
    protocol = object({
      tcp = optional(object({
        destination_port = optional(string),
        source_port = optional(string),
        override = optional(object({
          timeout = optional(number),
          halfclose_timeout = optional(number),
          timewait_timeout = optional(number),
        })),
      })),
      udp = optional(object({
        destination_port = optional(string),
        source_port = optional(string),
        override = optional(object({
          timeout = optional(number),
        })),
      })),
    }),
  }))
}

resource "panos_administrative_tag" "tag" {
  location = { shared = true }

  name = format("%s-tag", var.prefix)
}

resource "panos_service" "svc1" {
  depends_on = [ resource.panos_administrative_tag.tag ]

  location = { shared = true }

  name        = format("%s-svc1", var.prefix)
  description = var.services["svc1"].description
  tags        = var.services["svc1"].tags

  protocol = var.services["svc1"].protocol
}

resource "panos_service" "svc2" {
  depends_on = [ resource.panos_administrative_tag.tag ]

  location = { shared = true }

  name        = format("%s-svc2", var.prefix)
  description = var.services["svc2"].description
  tags        = var.services["svc2"].tags

  protocol = var.services["svc2"].protocol
}

resource "panos_service" "svc3" {
  depends_on = [ resource.panos_administrative_tag.tag ]

  location = { shared = true }

  name        = format("%s-svc3", var.prefix)
  description = var.services["svc3"].description
  tags        = var.services["svc3"].tags

  protocol = var.services["svc3"].protocol
}

resource "panos_service" "svc4" {
  depends_on = [ resource.panos_administrative_tag.tag ]

  location = { shared = true }

  name        = format("%s-svc4", var.prefix)
  description = var.services["svc4"].description
  tags        = var.services["svc4"].tags

  protocol = var.services["svc4"].protocol
}
`

func testAccServiceDestroy(prefix string) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		api := service.NewService(sdkClient)
		ctx := context.TODO()

		location := service.NewSharedLocation()

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
		}

		return nil
	}
}
