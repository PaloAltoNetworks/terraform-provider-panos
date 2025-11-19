package provider_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/PaloAltoNetworks/pango/objects/address"
)

func TestAccPanosAddress(t *testing.T) {
	t.Parallel()

	prefix := fmt.Sprintf("test-acc-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: panosAddress_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_address.netmask",
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-netmask", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_address.netmask",
						tfjsonpath.New("ip_netmask"),
						knownvalue.StringExact("192.168.80.151/32"),
					),
					statecheck.ExpectKnownValue(
						"panos_address.range",
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-range", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_address.range",
						tfjsonpath.New("ip_range"),
						knownvalue.StringExact("192.168.80.151-192.168.80.155"),
					),
					statecheck.ExpectKnownValue(
						"panos_address.fqdn",
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-fqdn", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_address.fqdn",
						tfjsonpath.New("fqdn"),
						knownvalue.StringExact("example.com"),
					),
					statecheck.ExpectKnownValue(
						"panos_address.wildcard",
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-wildcard", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_address.wildcard",
						tfjsonpath.New("ip_wildcard"),
						knownvalue.StringExact("192.168.0.0/0.0.255.255"),
					),
				},
			},
		},
	})
}

func TestAccPanosAddress_Rename_Basic(t *testing.T) {
	t.Parallel()

	prefix := fmt.Sprintf("test-acc-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: panosAddress_Rename_Initial_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_address.update",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_address.update",
						tfjsonpath.New("ip_netmask"),
						knownvalue.StringExact("10.0.0.1/32"),
					),
				},
			},
			{
				Config: panosAddress_Rename_Updated_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_address.update",
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-renamed", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_address.update",
						tfjsonpath.New("ip_netmask"),
						knownvalue.StringExact("10.0.0.2/32"),
					),
				},
			},
		},
	})
}

func TestAccPanosAddress_Rename_Missing_Source(t *testing.T) {
	t.Parallel()

	prefix := fmt.Sprintf("test-acc-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: panosAddress_Rename_Initial_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_address.update",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_address.update",
						tfjsonpath.New("ip_netmask"),
						knownvalue.StringExact("10.0.0.1/32"),
					),
				},
			},
			{
				PreConfig: func() {
					panosAddressDeleteObject(prefix, prefix)
				},
				Config: panosAddress_Rename_Updated_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_address.update",
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-renamed", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_address.update",
						tfjsonpath.New("ip_netmask"),
						knownvalue.StringExact("10.0.0.2/32"),
					),
				},
			},
		},
	})
}

func TestAccPanosAddress_Rename_Existing_Target(t *testing.T) {
	t.Parallel()

	prefix := fmt.Sprintf("test-acc-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: panosAddress_Rename_Initial_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_address.update",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_address.update",
						tfjsonpath.New("ip_netmask"),
						knownvalue.StringExact("10.0.0.1/32"),
					),
				},
			},
			{
				PreConfig: func() {
					panosAddressCreateObject(prefix, fmt.Sprintf("%s-renamed", prefix))
				},
				Config: panosAddress_Rename_Updated_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ExpectError: regexp.MustCompile(fmt.Sprintf("entry '%s-renamed' already exists", prefix)),
			},
		},
	})
}

const panosAddress_Tmpl = `
variable "prefix" { type = string }

resource "panos_device_group" "example" {
	location = { panorama = {} }
	name = var.prefix
}

resource "panos_administrative_tag" "a" {
	location = { device_group = { name = panos_device_group.example.name } }
	name = "a"
}

resource "panos_address" "netmask" {
    name = "${var.prefix}-netmask"
    location = { device_group = { name = panos_device_group.example.name } }
    ip_netmask = "192.168.80.151/32"
    description = "made by terraform"
    tags = [panos_administrative_tag.a.name]
}

resource "panos_address" "range" {
    name = "${var.prefix}-range"
    location = { device_group = { name = panos_device_group.example.name } }
    ip_range = "192.168.80.151-192.168.80.155"
    description = "made by terraform"
    tags = [panos_administrative_tag.a.name]
}

resource "panos_address" "fqdn" {
    name = "${var.prefix}-fqdn"
    location = { device_group = { name = panos_device_group.example.name } }
    fqdn = "example.com"
    description = "made by terraform"
    tags = [panos_administrative_tag.a.name]
}

resource "panos_address" "wildcard" {
    name = "${var.prefix}-wildcard"
    location = { device_group = { name = panos_device_group.example.name } }
    ip_wildcard = "192.168.0.0/0.0.255.255"
    description = "made by terraform"
}
`

const panosAddress_Rename_Initial_Tmpl = `
variable "prefix" { type = string }

resource "panos_device_group" "example" {
	location = { panorama = {} }
	name = var.prefix
}

resource "panos_address" "update" {
    name = "${var.prefix}"
    location = { device_group = { name = panos_device_group.example.name } }
    ip_netmask = "10.0.0.1/32"
    description = "made by terraform"
}
`

const panosAddress_Rename_Updated_Tmpl = `
variable "prefix" { type = string }

resource "panos_device_group" "example" {
	location = { panorama = {} }
	name = var.prefix
}

resource "panos_address" "update" {
    name = "${var.prefix}-renamed"
    location = { device_group = { name = panos_device_group.example.name } }
    ip_netmask = "10.0.0.2/32"
    description = "made by terraform"
}
`

func panosAddressDeleteObject(prefix string, name string) {
	svc := address.NewService(sdkClient)

	location := address.NewDeviceGroupLocation()
	location.DeviceGroup.DeviceGroup = prefix

	err := svc.Delete(context.TODO(), *location, name)
	if err != nil {
		panic(fmt.Sprintf("Failed to delete object from the device: %v", err))
	}
}

func panosAddressCreateObject(prefix string, name string) {
	svc := address.NewService(sdkClient)

	location := address.NewDeviceGroupLocation()
	location.DeviceGroup.DeviceGroup = prefix

	entry := &address.Entry{}
	entry.Name = name
	_, err := svc.Create(context.TODO(), *location, entry)
	if err != nil {
		panic(fmt.Sprintf("Failed to delete object from the device: %v", err))
	}
}
