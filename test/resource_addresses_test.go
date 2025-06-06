package provider_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	sdkerrors "github.com/PaloAltoNetworks/pango/errors"
	"github.com/PaloAltoNetworks/pango/objects/address"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

type expectServerAddressObjects struct {
	Location  address.Location
	Prefix    string
	Addresses []string
}

func ExpectServerAddressObjects(location address.Location, prefix string, addresses []string) *expectServerAddressObjects {
	return &expectServerAddressObjects{
		Location:  location,
		Prefix:    prefix,
		Addresses: addresses,
	}
}

func (o *expectServerAddressObjects) CheckState(ctx context.Context, req statecheck.CheckStateRequest, resp *statecheck.CheckStateResponse) {
	service := address.NewService(sdkClient)

	objects, err := service.List(ctx, o.Location, "get", "", "")
	if err != nil && !sdkerrors.IsObjectNotFound(err) {
		resp.Error = err
		return
	}

	objectsByName := make(map[string]int)
	for _, elt := range o.Addresses {
		objectsByName[fmt.Sprintf("%s-%s", o.Prefix, elt)] = 0
	}

	for _, elt := range objects {
		_, found := objectsByName[elt.Name]
		if !found {
			objectsByName[elt.Name] = -1
		} else {
			objectsByName[elt.Name] = 1
		}
	}

	var errors []string
	for name, state := range objectsByName {
		switch state {
		case -1:
			errors = append(errors, fmt.Sprintf("%s: unexpected", name))
		case 0:
			errors = append(errors, fmt.Sprintf("%s: missing", name))
		}
	}

	if errors != nil {
		resp.Error = fmt.Errorf("Unexpected server state: %s", strings.Join(errors, ", "))
	}
}

func TestAccAddresses(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	location := address.NewDeviceGroupLocation()
	location.DeviceGroup.DeviceGroup = prefix

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAddressesResourceTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":    config.StringVariable(prefix),
					"addresses": config.MapVariable(map[string]config.Variable{})},
			},
			{
				Config: testAccAddressesResourceTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":    config.StringVariable(prefix),
					"addresses": config.MapVariable(map[string]config.Variable{})},
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
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
							"ip_wildcard":      config.StringVariable("172.16.0.0/0.0.0.255"),
							"disable_override": config.StringVariable("no"),
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
							New("addresses"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							fmt.Sprintf("%s-ip-netmask", prefix): knownvalue.ObjectExact(map[string]knownvalue.Check{
								"tags": knownvalue.ListExact([]knownvalue.Check{
									knownvalue.StringExact(fmt.Sprintf("%s-tag", prefix)),
								}),
								"description":      knownvalue.Null(),
								"disable_override": knownvalue.Null(),
								"ip_netmask":       knownvalue.StringExact("172.16.0.1/32"),
								"ip_range":         knownvalue.Null(),
								"ip_wildcard":      knownvalue.Null(),
								"fqdn":             knownvalue.Null(),
							}),
							fmt.Sprintf("%s-ip-range", prefix): knownvalue.ObjectExact(map[string]knownvalue.Check{
								"tags":             knownvalue.Null(),
								"description":      knownvalue.StringExact("description"),
								"disable_override": knownvalue.Null(),
								"ip_netmask":       knownvalue.Null(),
								"ip_range":         knownvalue.StringExact("172.16.0.1-172.16.0.255"),
								"ip_wildcard":      knownvalue.Null(),
								"fqdn":             knownvalue.Null(),
							}),
							fmt.Sprintf("%s-ip-wildcard", prefix): knownvalue.ObjectExact(map[string]knownvalue.Check{
								"tags":             knownvalue.Null(),
								"description":      knownvalue.Null(),
								"disable_override": knownvalue.StringExact("no"),
								"ip_netmask":       knownvalue.Null(),
								"ip_range":         knownvalue.Null(),
								"ip_wildcard":      knownvalue.StringExact("172.16.0.0/0.0.0.255"),
								"fqdn":             knownvalue.Null(),
							}),
							fmt.Sprintf("%s-fqdn", prefix): knownvalue.ObjectExact(map[string]knownvalue.Check{
								"tags":             knownvalue.Null(),
								"description":      knownvalue.Null(),
								"disable_override": knownvalue.Null(),
								"ip_netmask":       knownvalue.Null(),
								"ip_range":         knownvalue.Null(),
								"ip_wildcard":      knownvalue.Null(),
								"fqdn":             knownvalue.StringExact("example.com"),
							}),
						}),
					),
					ExpectServerAddressObjects(*location, prefix, []string{"ip-netmask", "ip-range", "ip-wildcard", "fqdn"}),
				},
			},
			{
				Config: testAccAddressesResourceTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
					"addresses": config.MapVariable(map[string]config.Variable{
						fmt.Sprintf("%s-ip-range", prefix): config.ObjectVariable(map[string]config.Variable{
							"description": config.StringVariable("description"),
							"ip_range":    config.StringVariable("172.16.0.1-172.16.0.255"),
						}),
						fmt.Sprintf("%s-fqdn2", prefix): config.ObjectVariable(map[string]config.Variable{
							"fqdn": config.StringVariable("example.com"),
						}),
					}),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_addresses.addresses",
						tfjsonpath.
							New("addresses"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							fmt.Sprintf("%s-ip-range", prefix): knownvalue.ObjectExact(map[string]knownvalue.Check{
								"tags":             knownvalue.Null(),
								"description":      knownvalue.StringExact("description"),
								"disable_override": knownvalue.Null(),
								"ip_netmask":       knownvalue.Null(),
								"ip_range":         knownvalue.StringExact("172.16.0.1-172.16.0.255"),
								"ip_wildcard":      knownvalue.Null(),
								"fqdn":             knownvalue.Null(),
							}),
							fmt.Sprintf("%s-fqdn2", prefix): knownvalue.ObjectExact(map[string]knownvalue.Check{
								"tags":             knownvalue.Null(),
								"description":      knownvalue.Null(),
								"disable_override": knownvalue.Null(),
								"ip_netmask":       knownvalue.Null(),
								"ip_range":         knownvalue.Null(),
								"ip_wildcard":      knownvalue.Null(),
								"fqdn":             knownvalue.StringExact("example.com"),
							}),
						}),
					),
					ExpectServerAddressObjects(*location, prefix, []string{"ip-range", "fqdn2"}),
				},
			},
			{
				Config: testAccAddressesResourceTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":    config.StringVariable(prefix),
					"addresses": config.MapVariable(map[string]config.Variable{})},
				ConfigStateChecks: []statecheck.StateCheck{
					ExpectServerAddressObjects(*location, prefix, []string{}),
				},
			},
			{
				Config: testAccAddressesResourceTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":    config.StringVariable(prefix),
					"addresses": config.MapVariable(map[string]config.Variable{})},
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

const testAccAddressesResourceTmpl = `
variable prefix { type = string }
variable "addresses" {
  type = map(object({
    disable_override = optional(string),
    description = optional(string),
    ip_netmask = optional(string),
    ip_range = optional(string),
    ip_wildcard = optional(string),
    fqdn = optional(string),
    tags = optional(list(string)),
  }))
}

resource "panos_device_group" "example" {
  location = { panorama = {} }

  name = var.prefix
}

resource "panos_administrative_tag" "tag" {
  location = { device_group = { name = panos_device_group.example.name } }

  name = format("%s-tag", var.prefix)
}

resource "panos_addresses" "addresses" {
  depends_on = [ resource.panos_administrative_tag.tag ]

  location = { device_group = { name = panos_device_group.example.name } }

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
