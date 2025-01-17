package provider_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	sdkErrors "github.com/PaloAltoNetworks/pango/errors"
	"github.com/PaloAltoNetworks/pango/network/interface/loopback"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccLoopbackInterface(t *testing.T) {
	t.Parallel()

	interfaceName := "loopback.1"
	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		CheckDestroy: testAccCheckPanosLoopbackInterfaceDestroy(
			prefix, interfaceName,
		),
		Steps: []resource.TestStep{
			{
				Config: loopbackInterfaceResource1,
				ConfigVariables: map[string]config.Variable{
					"prefix":         config.StringVariable(prefix),
					"interface_name": config.StringVariable(interfaceName),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_loopback_interface.iface",
						tfjsonpath.New("name"),
						knownvalue.StringExact("loopback.1"),
					),
					statecheck.ExpectKnownValue(
						"panos_loopback_interface.iface",
						tfjsonpath.New("comment"),
						knownvalue.StringExact("loopback interface comment"),
					),
					statecheck.ExpectKnownValue(
						"panos_loopback_interface.iface",
						tfjsonpath.New("interface_management_profile"),
						knownvalue.StringExact(fmt.Sprintf("%s-profile", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_loopback_interface.iface",
						tfjsonpath.New("adjust_tcp_mss"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"enable":              knownvalue.Bool(true),
							"ipv4_mss_adjustment": knownvalue.Int64Exact(100),
							"ipv6_mss_adjustment": knownvalue.Int64Exact(200),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_loopback_interface.iface",
						tfjsonpath.New("ip"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name": knownvalue.StringExact("127.0.0.1"),
							}),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_loopback_interface.iface",
						tfjsonpath.New("ipv6"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"enabled":      knownvalue.Bool(true),
							"interface_id": knownvalue.StringExact("100"),
							"address": knownvalue.ListExact([]knownvalue.Check{
								knownvalue.ObjectExact(map[string]knownvalue.Check{
									"name":                knownvalue.StringExact("::1"),
									"enable_on_interface": knownvalue.Bool(true),
									"anycast":             knownvalue.ObjectExact(nil),
									"prefix":              knownvalue.ObjectExact(nil),
								}),
							}),
						}),
					),
				},
			},
		},
	})
}

const loopbackInterfaceResource1 = `
variable "prefix" { type = string }
variable "interface_name" { type = string }

locals {
  template_name = format("%s-tmpl", var.prefix)
}

resource "panos_template" "template" {
  location = { panorama = {} }
  name = local.template_name
}

resource "panos_interface_management_profile" "profile" {
  location = { template = { name = panos_template.template.name }}

  name = format("%s-profile", var.prefix)
}

resource "panos_loopback_interface" "iface" {
  location = { template = { name = panos_template.template.name } }

  name = var.interface_name
  comment = "loopback interface comment"

  interface_management_profile = panos_interface_management_profile.profile.name
  mtu = "9126"
  #netflow_profile = format("%s-profile", var.prefix)
  adjust_tcp_mss = {
    enable = true
    ipv4_mss_adjustment = 100
    ipv6_mss_adjustment = 200
  }
  ip = [{
    name = "127.0.0.1"
  }]
  ipv6 = {
    enabled = true
    interface_id = "100"
    address = [{
      name = "::1"
      enable_on_interface = true
      anycast = {}
      prefix = {}
    }]

  }
}
`

func testAccCheckPanosLoopbackInterfaceDestroy(prefix string, entry string) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		api := loopback.NewService(sdkClient)
		ctx := context.TODO()

		location := loopback.NewTemplateLocation()
		location.Template.Template = fmt.Sprintf("%s-tmpl", prefix)

		reply, err := api.Read(ctx, *location, entry, "show")
		if err != nil && !sdkErrors.IsObjectNotFound(err) {
			return fmt.Errorf("reading ethernet entry via sdk: %v", err)
		}

		if reply != nil {
			err := fmt.Errorf("terraform didn't delete the server entry properly")
			delErr := api.Delete(ctx, *location, entry)
			if delErr != nil {
				return errors.Join(err, delErr)
			}
			return err
		}

		return nil
	}
}
