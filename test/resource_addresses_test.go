package provider_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	sdkErrors "github.com/PaloAltoNetworks/pango/errors"
	"github.com/PaloAltoNetworks/pango/objects/address"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccAddresses(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		CheckDestroy:             testAccAddressesDestroy(prefix),
		Steps: []resource.TestStep{
			{
				Config: testAccAddressesResourceTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
					"addresses": config.MapVariable(map[string]config.Variable{
						fmt.Sprintf("%s-ip-netmask", prefix): config.ObjectVariable(map[string]config.Variable{
							"tags":       config.ListVariable(config.StringVariable(fmt.Sprintf("%s-tag", prefix))),
							"ip_netmask": config.StringVariable("172.16.0.1/32"),
						}),
						fmt.Sprintf("%s-ip-range", prefix): config.ObjectVariable(map[string]config.Variable{
							"description": config.StringVariable("description"),
							"ip_range":    config.StringVariable("172.16.0.1-172.16.0.255"),
						}),
						fmt.Sprintf("%s-ip-wildcard", prefix): config.ObjectVariable(map[string]config.Variable{
							"ip_wildcard": config.StringVariable("172.16.0.0/0.0.0.255"),
						}),
						fmt.Sprintf("%s-fqdn", prefix): config.ObjectVariable(map[string]config.Variable{
							"fqdn": config.StringVariable("example.com"),
						}),
					}),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_addresses.addresses",
						tfjsonpath.
							New("addresses").
							AtMapKey(fmt.Sprintf("%s-ip-netmask", prefix)).
							AtMapKey("ip_netmask"),
						knownvalue.StringExact("172.16.0.1/32"),
					),
					statecheck.ExpectKnownValue(
						"panos_addresses.addresses",
						tfjsonpath.
							New("addresses").
							AtMapKey(fmt.Sprintf("%s-ip-netmask", prefix)).
							AtMapKey("tags"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact(fmt.Sprintf("%s-tag", prefix)),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_addresses.addresses",
						tfjsonpath.
							New("addresses").
							AtMapKey(fmt.Sprintf("%s-ip-range", prefix)).
							AtMapKey("description"),
						knownvalue.StringExact("description"),
					),
					statecheck.ExpectKnownValue(
						"panos_addresses.addresses",
						tfjsonpath.
							New("addresses").
							AtMapKey(fmt.Sprintf("%s-ip-range", prefix)).
							AtMapKey("ip_range"),
						knownvalue.StringExact("172.16.0.1-172.16.0.255"),
					),
					statecheck.ExpectKnownValue(
						"panos_addresses.addresses",
						tfjsonpath.
							New("addresses").
							AtMapKey(fmt.Sprintf("%s-ip-wildcard", prefix)).
							AtMapKey("ip_wildcard"),
						knownvalue.StringExact("172.16.0.0/0.0.0.255"),
					),
					statecheck.ExpectKnownValue(
						"panos_addresses.addresses",
						tfjsonpath.
							New("addresses").
							AtMapKey(fmt.Sprintf("%s-fqdn", prefix)).
							AtMapKey("fqdn"),
						knownvalue.StringExact("example.com"),
					),
				},
			},
		},
	})
}

const testAccAddressesResourceTmpl = `
variable prefix { type = string }
variable "addresses" {
  type = map(object({
    disable_override = optional(bool),
    description = optional(string),
    ip_netmask = optional(string),
    ip_range = optional(string),
    ip_wildcard = optional(string),
    fqdn = optional(string),
    tags = optional(list(string)),
  }))
}

resource "panos_administrative_tag" "tag" {
  location = { shared = true }

  name = format("%s-tag", var.prefix)
}

resource "panos_addresses" "addresses" {
  depends_on = [ resource.panos_administrative_tag.tag ]

  location = { shared = true }

  addresses = { for name, value in var.addresses : name => {
    disable_override = value.disable_override,
    description = value.description,
    ip_netmask = value.ip_netmask,
    ip_range = value.ip_range,
    fqdn = value.fqdn,
    ip_wildcard = value.ip_wildcard,
    tags = value.tags
  }}
}
`

func testAccAddressesDestroy(prefix string) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		api := address.NewService(sdkClient)
		ctx := context.TODO()

		location := address.NewSharedLocation()

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
