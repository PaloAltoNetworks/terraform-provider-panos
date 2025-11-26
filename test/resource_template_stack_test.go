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

func TestAccTemplateStack_Basic(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	location := config.ObjectVariable(map[string]config.Variable{
		"panorama": config.ObjectVariable(map[string]config.Variable{}),
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: templateStack_Basic_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_template_stack.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_template_stack.example",
						tfjsonpath.New("description"),
						knownvalue.StringExact("Template stack description"),
					),
					statecheck.ExpectKnownValue(
						"panos_template_stack.example",
						tfjsonpath.New("templates"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact(fmt.Sprintf("%s-template1", prefix)),
							knownvalue.StringExact(fmt.Sprintf("%s-template2", prefix)),
						}),
					),
				},
			},
		},
	})
}

const templateStack_Basic_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "template1" {
  location = var.location
  name = "${var.prefix}-template1"
}

resource "panos_template" "template2" {
  location = var.location
  name = "${var.prefix}-template2"
}

resource "panos_template_stack" "example" {
  location = var.location
  name = var.prefix
  description = "Template stack description"
  templates = [
    panos_template.template1.name,
    panos_template.template2.name
  ]
}
`

func TestAccTemplateStack_Devices(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	// Generate serial numbers for firewall devices
	suffix1 := acctest.RandStringFromCharSet(13, "0123456789")
	serialNumber1 := fmt.Sprintf("00%s", suffix1)
	suffix2 := acctest.RandStringFromCharSet(13, "0123456789")
	serialNumber2 := fmt.Sprintf("00%s", suffix2)

	location := config.ObjectVariable(map[string]config.Variable{
		"panorama": config.ObjectVariable(map[string]config.Variable{}),
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: templateStack_Devices_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":          config.StringVariable(prefix),
					"location":        location,
					"serial_number_1": config.StringVariable(serialNumber1),
					"serial_number_2": config.StringVariable(serialNumber2),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_template_stack.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_template_stack.example",
						tfjsonpath.New("devices"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name": knownvalue.StringExact(serialNumber1),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name": knownvalue.StringExact(serialNumber2),
							}),
						}),
					),
				},
			},
		},
	})
}

const templateStack_Devices_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }
variable "serial_number_1" { type = string }
variable "serial_number_2" { type = string }

resource "panos_firewall_device" "device1" {
  location = var.location
  name = var.serial_number_1
  hostname = "fw1.example.com"
  ip = "192.0.2.1"
}

resource "panos_firewall_device" "device2" {
  location = var.location
  name = var.serial_number_2
  hostname = "fw2.example.com"
  ip = "192.0.2.2"
}

resource "panos_template_stack" "example" {
  location = var.location
  name = var.prefix
  devices = [
    { name = panos_firewall_device.device1.name },
    { name = panos_firewall_device.device2.name }
  ]
}
`

func TestAccTemplateStack_DefaultVsys(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	location := config.ObjectVariable(map[string]config.Variable{
		"panorama": config.ObjectVariable(map[string]config.Variable{}),
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: templateStack_DefaultVsys_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_template_stack.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_template_stack.example",
						tfjsonpath.New("default_vsys"),
						knownvalue.StringExact("vsys1"),
					),
				},
			},
		},
	})
}

const templateStack_DefaultVsys_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

data "panos_template" "existing" {
  location = {
    panorama = {
      panorama_device = "localhost.localdomain"
    }
  }
  name = "test-acc-tmpl"
}

resource "panos_template_stack" "example" {
  location = var.location
  name = var.prefix
  default_vsys = "vsys1"
  templates = [data.panos_template.existing.name]
}
`

func TestAccTemplateStack_Complete(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	// Generate serial numbers for firewall devices
	suffix1 := acctest.RandStringFromCharSet(13, "0123456789")
	serialNumber1 := fmt.Sprintf("00%s", suffix1)
	suffix2 := acctest.RandStringFromCharSet(13, "0123456789")
	serialNumber2 := fmt.Sprintf("00%s", suffix2)

	location := config.ObjectVariable(map[string]config.Variable{
		"panorama": config.ObjectVariable(map[string]config.Variable{}),
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: templateStack_Complete_Step1_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":          config.StringVariable(prefix),
					"location":        location,
					"serial_number_1": config.StringVariable(serialNumber1),
					"serial_number_2": config.StringVariable(serialNumber2),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_template_stack.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_template_stack.example",
						tfjsonpath.New("description"),
						knownvalue.StringExact("Complete template stack"),
					),
					statecheck.ExpectKnownValue(
						"panos_template_stack.example",
						tfjsonpath.New("devices"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name": knownvalue.StringExact(serialNumber1),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name": knownvalue.StringExact(serialNumber2),
							}),
						}),
					),
				},
			},
			{
				Config: templateStack_Complete_Step2_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":          config.StringVariable(prefix),
					"location":        location,
					"serial_number_1": config.StringVariable(serialNumber1),
					"serial_number_2": config.StringVariable(serialNumber2),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_template_stack.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_template_stack.example",
						tfjsonpath.New("description"),
						knownvalue.StringExact("Complete template stack"),
					),
					statecheck.ExpectKnownValue(
						"panos_template_stack.example",
						tfjsonpath.New("devices"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name": knownvalue.StringExact(serialNumber1),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name": knownvalue.StringExact(serialNumber2),
							}),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_template_stack.example",
						tfjsonpath.New("default_vsys"),
						knownvalue.StringExact("vsys1"),
					),
				},
			},
		},
	})
}

const templateStack_Complete_Step1_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }
variable "serial_number_1" { type = string }
variable "serial_number_2" { type = string }

data "panos_template" "existing" {
  location = {
    panorama = {
      panorama_device = "localhost.localdomain"
    }
  }
  name = "test-acc-tmpl"
}

resource "panos_firewall_device" "device1" {
  location = var.location
  name = var.serial_number_1
  hostname = "fw1.example.com"
  ip = "192.0.2.1"
}

resource "panos_firewall_device" "device2" {
  location = var.location
  name = var.serial_number_2
  hostname = "fw2.example.com"
  ip = "192.0.2.2"
}

resource "panos_template_stack" "example" {
  location = var.location
  name = var.prefix
  description = "Complete template stack"
  templates = [data.panos_template.existing.name]
  devices = [
    { name = panos_firewall_device.device1.name },
    { name = panos_firewall_device.device2.name }
  ]
}
`

const templateStack_Complete_Step2_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }
variable "serial_number_1" { type = string }
variable "serial_number_2" { type = string }

data "panos_template" "existing" {
  location = {
    panorama = {
      panorama_device = "localhost.localdomain"
    }
  }
  name = "test-acc-tmpl"
}

resource "panos_firewall_device" "device1" {
  location = var.location
  name = var.serial_number_1
  hostname = "fw1.example.com"
  ip = "192.0.2.1"
}

resource "panos_firewall_device" "device2" {
  location = var.location
  name = var.serial_number_2
  hostname = "fw2.example.com"
  ip = "192.0.2.2"
}

resource "panos_template_stack" "example" {
  location = var.location
  name = var.prefix
  description = "Complete template stack"
  templates = [data.panos_template.existing.name]
  devices = [
    { name = panos_firewall_device.device1.name },
    { name = panos_firewall_device.device2.name }
  ]
  default_vsys = "vsys1"
}
`
