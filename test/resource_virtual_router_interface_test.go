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

func TestAccVirtualRouterInterface_Basic(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)
	vrName := prefix
	ethName := "ethernet1/1"

	location := config.ObjectVariable(map[string]config.Variable{
		"template": config.ObjectVariable(map[string]config.Variable{
			"name": config.StringVariable(prefix),
		}),
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: virtualRouterInterface_Basic_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
					"vr_name":  config.StringVariable(vrName),
					"eth_name": config.StringVariable(ethName),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_virtual_router_interface.example",
						tfjsonpath.New("virtual_router"),
						knownvalue.StringExact(vrName),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router_interface.example",
						tfjsonpath.New("interface"),
						knownvalue.StringExact(ethName),
					),
				},
			},
		},
	})
}

func TestAccVirtualRouterInterface_UpdateVirtualRouter(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)
	vrName1 := prefix + "-1"
	vrName2 := prefix + "-2"
	ethName := "ethernet1/1"

	location := config.ObjectVariable(map[string]config.Variable{
		"template": config.ObjectVariable(map[string]config.Variable{
			"name": config.StringVariable(prefix),
		}),
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: virtualRouterInterface_UpdateVirtualRouter_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":    config.StringVariable(prefix),
					"location":  location,
					"vr_name_1": config.StringVariable(vrName1),
					"vr_name_2": config.StringVariable(vrName2),
					"eth_name":  config.StringVariable(ethName),
					"vr_name":   config.StringVariable(vrName1),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_virtual_router_interface.example",
						tfjsonpath.New("virtual_router"),
						knownvalue.StringExact(vrName1),
					),
				},
			},
			{
				Config: virtualRouterInterface_UpdateVirtualRouter_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":    config.StringVariable(prefix),
					"location":  location,
					"vr_name_1": config.StringVariable(vrName1),
					"vr_name_2": config.StringVariable(vrName2),
					"eth_name":  config.StringVariable(ethName),
					"vr_name":   config.StringVariable(vrName2),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_virtual_router_interface.example",
						tfjsonpath.New("virtual_router"),
						knownvalue.StringExact(vrName2),
					),
				},
			},
		},
	})
}

func TestAccVirtualRouterInterface_UpdateInterface(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)
	vrName := prefix

	location := config.ObjectVariable(map[string]config.Variable{
		"template": config.ObjectVariable(map[string]config.Variable{
			"name": config.StringVariable(prefix),
		}),
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: virtualRouterInterface_UpdateInterface_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
					"vr_name":  config.StringVariable(vrName),
					"eth_name": config.StringVariable("ethernet1/1"),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_virtual_router_interface.example",
						tfjsonpath.New("interface"),
						knownvalue.StringExact("ethernet1/1"),
					),
				},
			},
			{
				Config: virtualRouterInterface_UpdateInterface_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
					"vr_name":  config.StringVariable(vrName),
					"eth_name": config.StringVariable("ethernet1/2"),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_virtual_router_interface.example",
						tfjsonpath.New("interface"),
						knownvalue.StringExact("ethernet1/2"),
					),
				},
			},
		},
	})
}

func TestAccVirtualRouterInterface_MultipleInterfaces(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)
	vrName := prefix

	location := config.ObjectVariable(map[string]config.Variable{
		"template": config.ObjectVariable(map[string]config.Variable{
			"name": config.StringVariable(prefix),
		}),
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: virtualRouterInterface_MultipleInterfaces_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":           config.StringVariable(prefix),
					"location":         location,
					"vr_name":          config.StringVariable(vrName),
					"interfaces_count": config.IntegerVariable(6),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_virtual_router_interface.example[0]",
						tfjsonpath.New("interface"),
						knownvalue.StringExact("ethernet1/1"),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router_interface.example[1]",
						tfjsonpath.New("interface"),
						knownvalue.StringExact("ethernet1/2"),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router_interface.example[2]",
						tfjsonpath.New("interface"),
						knownvalue.StringExact("ethernet1/3"),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router_interface.example[3]",
						tfjsonpath.New("interface"),
						knownvalue.StringExact("ethernet1/4"),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router_interface.example[4]",
						tfjsonpath.New("interface"),
						knownvalue.StringExact("ethernet1/5"),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router_interface.example[5]",
						tfjsonpath.New("interface"),
						knownvalue.StringExact("ethernet1/6"),
					),
				},
			},
			{
				Config: virtualRouterInterface_MultipleInterfaces_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":           config.StringVariable(prefix),
					"location":         location,
					"vr_name":          config.StringVariable(vrName),
					"interfaces_count": config.IntegerVariable(2),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_virtual_router_interface.example[0]",
						tfjsonpath.New("interface"),
						knownvalue.StringExact("ethernet1/1"),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router_interface.example[1]",
						tfjsonpath.New("interface"),
						knownvalue.StringExact("ethernet1/2"),
					),
				},
			},
		},
	})
}

const virtualRouterInterface_UpdateInterface_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }
variable "vr_name" { type = string }
variable "eth_name" { type = string }

resource "panos_template" "example" {
	location = { panorama = {} }
	name     = var.prefix
}

resource "panos_virtual_router" "example" {
	location = var.location
	name     = var.vr_name
	lifecycle {
		ignore_changes = [interfaces]
	}
}

resource "panos_ethernet_interface" "example_1" {
	location = var.location
	name     = "ethernet1/1"
	layer3 = {
		ips = [{ name = "10.1.1.1/24" }]
	}
}

resource "panos_ethernet_interface" "example_2" {
	location = var.location
	name     = "ethernet1/2"
	layer3 = {
		ips = [{ name = "10.1.1.2/24" }]
	}
}

resource "panos_virtual_router_interface" "example" {
	depends_on = [panos_ethernet_interface.example_1, panos_ethernet_interface.example_2]
	location       = var.location
	virtual_router = panos_virtual_router.example.name
	interface      = var.eth_name
}
`
const virtualRouterInterface_UpdateVirtualRouter_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }
variable "vr_name_1" { type = string }
variable "vr_name_2" { type = string }
variable "eth_name" { type = string }
variable "vr_name" { type = string }

resource "panos_template" "example" {
	location = { panorama = {} }
	name     = var.prefix
}

resource "panos_virtual_router" "example_1" {
	location = var.location
	name     = var.vr_name_1
	lifecycle {
		ignore_changes = [interfaces]
	}
}

resource "panos_virtual_router" "example_2" {
	location = var.location
	name     = var.vr_name_2
	lifecycle {
		ignore_changes = [interfaces]
	}
}

resource "panos_ethernet_interface" "example" {
	location = var.location
	name     = var.eth_name
	layer3 = {
		ips = [{ name = "10.1.1.1/24" }]
	}
}

resource "panos_virtual_router_interface" "example" {
	location       = var.location
	virtual_router = var.vr_name
	interface      = panos_ethernet_interface.example.name
}
`
const virtualRouterInterface_Basic_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }
variable "vr_name" { type = string }
variable "eth_name" { type = string }

resource "panos_template" "example" {
	location = { panorama = {} }
	name     = var.prefix
}

resource "panos_virtual_router" "example" {
	location = var.location
	name     = var.vr_name

	lifecycle {
		ignore_changes = [interfaces]
	}
}

resource "panos_ethernet_interface" "example" {
	location = var.location
	name     = var.eth_name
	layer3 = {
		ips = [{ name = "10.1.1.1/24" }]
	}
}

resource "panos_virtual_router_interface" "example" {
	location       = var.location
	virtual_router = panos_virtual_router.example.name
	interface      = panos_ethernet_interface.example.name
}
`

const virtualRouterInterface_MultipleInterfaces_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }
variable "vr_name" { type = string }
variable "interfaces_count" { type = number }

resource "panos_template" "example" {
	location = { panorama = {} }
	name     = var.prefix
}

resource "panos_virtual_router" "example" {
	location = var.location
	name     = var.vr_name

	lifecycle {
		ignore_changes = [interfaces]
	}
}

resource "panos_ethernet_interface" "example" {
	count    = var.interfaces_count
	location = var.location
	name     = "ethernet1/${count.index + 1}"
	layer3 = {
		ips = [{ name = "10.1.1.${count.index + 1}/24"}]
	}
}

resource "panos_virtual_router_interface" "example" {
	count          = var.interfaces_count
	location       = var.location
	virtual_router = panos_virtual_router.example.name
	interface      = panos_ethernet_interface.example[count.index].name
}
`

const virtualRouterInterface_ServerSideDrift_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
	location = { panorama = {} }

	name = var.prefix
}

resource "panos_virtual_router" "example" {
	location = var.location
	name     = var.vr_name

	lifecycle {
		ignore_changes = [interfaces]
	}
}

resource "panos_ethernet_interface" "example" {
	count    = var.interfaces_count
	location = var.location
	name     = "ethernet1/${count.index + 1}"
	layer3 = {
		ips = [{ name = "10.1.1.${count.index + 1}/24"}]
	}
}
resource "panos_virtual_router_interface" "example" {
	location = var.location

	virtual_router = panos_virtual_router.example.name
}
`
