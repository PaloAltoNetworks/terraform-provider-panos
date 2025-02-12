package provider_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	sdkErrors "github.com/PaloAltoNetworks/pango/errors"
	"github.com/PaloAltoNetworks/pango/objects/application"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccApplication(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		CheckDestroy:             testAccApplicationDestroy(prefix),
		Steps: []resource.TestStep{
			{
				Config: testApplicationTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_application.app1",
						tfjsonpath.New("able_to_transfer_file"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_application.app1",
						tfjsonpath.New("category"),
						knownvalue.StringExact("general-internet"),
					),
					statecheck.ExpectKnownValue(
						"panos_application.app1",
						tfjsonpath.New("consume_big_bandwidth"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_application.app1",
						tfjsonpath.New("data_ident"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_application.app1",
						tfjsonpath.New("description"),
						knownvalue.StringExact("description"),
					),
					statecheck.ExpectKnownValue(
						"panos_application.app1",
						tfjsonpath.New("disable_override"),
						knownvalue.StringExact("no"),
					),
					statecheck.ExpectKnownValue(
						"panos_application.app1",
						tfjsonpath.New("evasive_behavior"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_application.app1",
						tfjsonpath.New("file_type_ident"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_application.app1",
						tfjsonpath.New("has_known_vulnerability"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_application.app1",
						tfjsonpath.New("no_appid_caching"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_application.app1",
						tfjsonpath.New("parent_app"),
						knownvalue.StringExact("8x8"),
					),
					statecheck.ExpectKnownValue(
						"panos_application.app1",
						tfjsonpath.New("pervasive_use"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_application.app1",
						tfjsonpath.New("prone_to_misuse"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_application.app1",
						tfjsonpath.New("risk"),
						knownvalue.Int64Exact(1),
					),
					statecheck.ExpectKnownValue(
						"panos_application.app1",
						tfjsonpath.New("subcategory"),
						knownvalue.StringExact("internet-utility"),
					),
					statecheck.ExpectKnownValue(
						"panos_application.app1",
						tfjsonpath.New("tcp_half_closed_timeout"),
						knownvalue.Int64Exact(60),
					),
					statecheck.ExpectKnownValue(
						"panos_application.app1",
						tfjsonpath.New("tcp_time_wait_timeout"),
						knownvalue.Int64Exact(120),
					),
					statecheck.ExpectKnownValue(
						"panos_application.app1",
						tfjsonpath.New("tcp_timeout"),
						knownvalue.Int64Exact(180),
					),
					statecheck.ExpectKnownValue(
						"panos_application.app1",
						tfjsonpath.New("technology"),
						knownvalue.StringExact("browser-based"),
					),
					statecheck.ExpectKnownValue(
						"panos_application.app1",
						tfjsonpath.New("timeout"),
						knownvalue.Int64Exact(240),
					),
					statecheck.ExpectKnownValue(
						"panos_application.app1",
						tfjsonpath.New("tunnel_applications"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_application.app1",
						tfjsonpath.New("tunnel_other_application"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_application.app1",
						tfjsonpath.New("udp_timeout"),
						knownvalue.Int64Exact(120),
					),
					statecheck.ExpectKnownValue(
						"panos_application.app1",
						tfjsonpath.New("used_by_malware"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_application.app1",
						tfjsonpath.New("virus_ident"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_application.app1",
						tfjsonpath.New("default"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"ident_by_icmp_type": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"code": knownvalue.StringExact("1"),
								"type": knownvalue.StringExact("1"),
							}),
							"ident_by_icmp6_type":  knownvalue.Null(),
							"ident_by_ip_protocol": knownvalue.Null(),
							"port":                 knownvalue.Null(),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_application.app1",
						tfjsonpath.New("signature"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":          knownvalue.StringExact("signature"),
								"comment":       knownvalue.StringExact("comment"),
								"scope":         knownvalue.StringExact("protocol-data-unit"),
								"order_free":    knownvalue.Bool(true),
								"and_condition": knownvalue.Null(),
							}),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_application.app2",
						tfjsonpath.New("default"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"ident_by_icmp6_type": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"code": knownvalue.StringExact("1"),
								"type": knownvalue.StringExact("1"),
							}),
							"ident_by_icmp_type":   knownvalue.Null(),
							"ident_by_ip_protocol": knownvalue.Null(),
							"port":                 knownvalue.Null(),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_application.app3",
						tfjsonpath.New("default"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"ident_by_icmp_type":   knownvalue.Null(),
							"ident_by_icmp6_type":  knownvalue.Null(),
							"ident_by_ip_protocol": knownvalue.StringExact("1"),
							"port":                 knownvalue.Null(),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_application.app4",
						tfjsonpath.New("default"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"ident_by_icmp_type":   knownvalue.Null(),
							"ident_by_icmp6_type":  knownvalue.Null(),
							"ident_by_ip_protocol": knownvalue.Null(),
							"port": knownvalue.ListExact([]knownvalue.Check{
								knownvalue.StringExact("tcp/80"),
							}),
						}),
					),
				},
			},
		},
	})
}

const testApplicationTmpl = `
variable prefix { type = string }

resource "panos_device_group" "dg" {
  location = { panorama = {} }

  name = format("%s-dg", var.prefix)
}

resource "panos_application" "app1" {
  location = { device_group = { name = panos_device_group.dg.name } }

  name = format("%s-app1", var.prefix)

  able_to_transfer_file = true
  #alg_disable_capability = "yes"
  category = "general-internet"
  consume_big_bandwidth = true
  data_ident = true
  description = "description"
  disable_override = "no"
  evasive_behavior = true
  file_type_ident = true
  has_known_vulnerability = true
  no_appid_caching = true
  parent_app = "8x8"
  pervasive_use = true
  prone_to_misuse = true
  risk = 1
  subcategory = "internet-utility"
  tcp_half_closed_timeout = 60
  tcp_time_wait_timeout = 120
  tcp_timeout = 180
  technology = "browser-based"
  timeout = 240
  tunnel_applications = true
  tunnel_other_application = true
  udp_timeout = 120
  used_by_malware = true
  virus_ident = true
  default = {
    ident_by_icmp_type = {
      code = "1"
      type = "1"
    }
  }
  signature = [{
    name = "signature"
    comment = "comment"
    scope = "protocol-data-unit"
    order_free = true
  }]
}

resource "panos_application" "app2" {
  location = { device_group = { name = panos_device_group.dg.name } }

  name = format("%s-app2", var.prefix)

  default = {
    ident_by_icmp6_type = {
      code = "1"
      type = "1"
    }
  }
}

resource "panos_application" "app3" {
  location = { device_group = { name = panos_device_group.dg.name } }

  name = format("%s-app3", var.prefix)

  default = {
    ident_by_ip_protocol = "1"
  }
}

resource "panos_application" "app4" {
  location = { device_group = { name = panos_device_group.dg.name } }

  name = format("%s-app4", var.prefix)

  default = {
    port = ["tcp/80"]
  }
}

`

func testAccApplicationDestroy(prefix string) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		api := application.NewService(sdkClient)
		ctx := context.TODO()

		location := application.NewSharedLocation()

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
