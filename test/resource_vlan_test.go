package provider_test

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

// TestAccPanosVlan_Basic creates a VLAN with one interface and a virtual_interface
// block and checks every attribute in state.
func TestAccPanosVlan_Basic(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)
	templateName := fmt.Sprintf("%s-tmpl", prefix)
	vlanName := fmt.Sprintf("%s-vlan", prefix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPanosVlanBasic,
				ConfigVariables: map[string]config.Variable{
					"template_name": config.StringVariable(templateName),
					"vlan_name":     config.StringVariable(vlanName),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_vlan.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact(vlanName),
					),
					statecheck.ExpectKnownValue(
						"panos_vlan.test",
						tfjsonpath.New("interfaces").AtSliceIndex(0),
						knownvalue.StringExact("ethernet1/1"),
					),
					statecheck.ExpectKnownValue(
						"panos_vlan.test",
						tfjsonpath.New("virtual_interface").AtMapKey("interface"),
						knownvalue.StringExact("vlan.1"),
					),
				},
			},
		},
	})
}

const testAccPanosVlanBasic = `
variable "template_name" { type = string }
variable "vlan_name"     { type = string }

resource "panos_template" "test" {
  location = { panorama = {} }
  name     = var.template_name
}

resource "panos_ethernet_interface" "eth1" {
  location = {
    template = {
      vsys = "vsys1"
      name = panos_template.test.name
    }
  }

  name   = "ethernet1/1"
  layer2 = {}
}

resource "panos_vlan_interface" "vlan1" {
  location = { template = { name = panos_template.test.name } }
  name     = "vlan.1"
}

resource "panos_vlan" "test" {
  depends_on = [panos_template.test]

  location = {
    template = {
      name = panos_template.test.name
    }
  }

  name       = var.vlan_name
  interfaces = [panos_ethernet_interface.eth1.name]

  virtual_interface = {
    interface = panos_vlan_interface.vlan1.name
  }
}
`

// TestAccPanosVlan_NoOptionals creates a VLAN with only a name — no interfaces
// and no virtual_interface.
func TestAccPanosVlan_NoOptionals(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)
	templateName := fmt.Sprintf("%s-tmpl", prefix)
	vlanName := fmt.Sprintf("%s-vlan", prefix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPanosVlanNoOptionals,
				ConfigVariables: map[string]config.Variable{
					"template_name": config.StringVariable(templateName),
					"vlan_name":     config.StringVariable(vlanName),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_vlan.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact(vlanName),
					),
					statecheck.ExpectKnownValue(
						"panos_vlan.test",
						tfjsonpath.New("interfaces"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"panos_vlan.test",
						tfjsonpath.New("virtual_interface"),
						knownvalue.Null(),
					),
				},
			},
		},
	})
}

const testAccPanosVlanNoOptionals = `
variable "template_name" { type = string }
variable "vlan_name"     { type = string }

resource "panos_template" "test" {
  location = { panorama = {} }
  name     = var.template_name
}

resource "panos_vlan" "test" {
  depends_on = [panos_template.test]

  location = {
    template = {
      name = panos_template.test.name
    }
  }

  name = var.vlan_name
}
`

// TestAccPanosVlan_MultipleInterfaces verifies that a VLAN can hold two
// ethernet interfaces in its member list.
func TestAccPanosVlan_MultipleInterfaces(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)
	templateName := fmt.Sprintf("%s-tmpl", prefix)
	vlanName := fmt.Sprintf("%s-vlan", prefix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPanosVlanMultipleInterfaces,
				ConfigVariables: map[string]config.Variable{
					"template_name": config.StringVariable(templateName),
					"vlan_name":     config.StringVariable(vlanName),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_vlan.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact(vlanName),
					),
					statecheck.ExpectKnownValue(
						"panos_vlan.test",
						tfjsonpath.New("interfaces").AtSliceIndex(0),
						knownvalue.StringExact("ethernet1/1"),
					),
					statecheck.ExpectKnownValue(
						"panos_vlan.test",
						tfjsonpath.New("interfaces").AtSliceIndex(1),
						knownvalue.StringExact("ethernet1/2"),
					),
				},
			},
		},
	})
}

const testAccPanosVlanMultipleInterfaces = `
variable "template_name" { type = string }
variable "vlan_name"     { type = string }

resource "panos_template" "test" {
  location = { panorama = {} }
  name     = var.template_name
}

resource "panos_ethernet_interface" "eth1" {
  location = {
    template = {
      vsys = "vsys1"
      name = panos_template.test.name
    }
  }

  name   = "ethernet1/1"
  layer2 = {}
}

resource "panos_ethernet_interface" "eth2" {
  location = {
    template = {
      vsys = "vsys1"
      name = panos_template.test.name
    }
  }

  name   = "ethernet1/2"
  layer2 = {}
}

resource "panos_vlan" "test" {
  depends_on = [panos_template.test]

  location = {
    template = {
      name = panos_template.test.name
    }
  }

  name = var.vlan_name
  interfaces = [
    panos_ethernet_interface.eth1.name,
    panos_ethernet_interface.eth2.name,
  ]
}
`

// TestAccPanosVlan_Import verifies that a VLAN can be imported by its
// base64-encoded JSON import ID.
func TestAccPanosVlan_Import(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)
	templateName := fmt.Sprintf("%s-tmpl", prefix)
	vlanName := fmt.Sprintf("%s-vlan", prefix)

	importStateGenerateID := func(state *terraform.State) (string, error) {
		importState := map[string]any{
			"location": map[string]any{
				"template": map[string]any{
					"name":            templateName,
					"panorama_device": "localhost.localdomain",
					"ngfw_device":     "localhost.localdomain",
				},
			},
			"name": vlanName,
		}

		marshalled, err := json.Marshal(importState)
		if err != nil {
			return "", fmt.Errorf("failed to marshal import state: %w", err)
		}

		return base64.StdEncoding.EncodeToString(marshalled), nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPanosVlanImportConfig,
				ConfigVariables: map[string]config.Variable{
					"template_name": config.StringVariable(templateName),
					"vlan_name":     config.StringVariable(vlanName),
				},
			},
			{
				Config: testAccPanosVlanImportStep,
				ConfigVariables: map[string]config.Variable{
					"template_name": config.StringVariable(templateName),
					"vlan_name":     config.StringVariable(vlanName),
				},
				ResourceName:      "panos_vlan.imported",
				ImportState:       true,
				ImportStateIdFunc: importStateGenerateID,
			},
		},
	})
}

const testAccPanosVlanImportConfig = `
variable "template_name" { type = string }
variable "vlan_name"     { type = string }

resource "panos_template" "test" {
  location = { panorama = {} }
  name     = var.template_name
}

resource "panos_ethernet_interface" "eth1" {
  location = {
    template = {
      vsys = "vsys1"
      name = panos_template.test.name
    }
  }

  name   = "ethernet1/1"
  layer2 = {}
}

resource "panos_vlan_interface" "vlan1" {
  location = { template = { name = panos_template.test.name } }
  name     = "vlan.1"
}

resource "panos_vlan" "test" {
  depends_on = [panos_template.test]

  location = {
    template = {
      name = panos_template.test.name
    }
  }

  name       = var.vlan_name
  interfaces = [panos_ethernet_interface.eth1.name]

  virtual_interface = {
    interface = panos_vlan_interface.vlan1.name
  }
}
`

const testAccPanosVlanImportStep = `
variable "template_name" { type = string }
variable "vlan_name"     { type = string }

resource "panos_template" "test" {
  location = { panorama = {} }
  name     = var.template_name
}

resource "panos_ethernet_interface" "eth1" {
  location = {
    template = {
      vsys = "vsys1"
      name = panos_template.test.name
    }
  }

  name   = "ethernet1/1"
  layer2 = {}
}

resource "panos_vlan_interface" "vlan1" {
  location = { template = { name = panos_template.test.name } }
  name     = "vlan.1"
}

resource "panos_vlan" "imported" {
  depends_on = [panos_template.test]

  location = {
    template = {
      name = panos_template.test.name
    }
  }

  name       = var.vlan_name
  interfaces = [panos_ethernet_interface.eth1.name]

  virtual_interface = {
    interface = panos_vlan_interface.vlan1.name
  }
}
`
