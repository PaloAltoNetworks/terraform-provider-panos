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

// TestAccPanosVlanEntry_Basic creates a single static-MAC entry under a parent
// VLAN and verifies all three attributes (vlan, name, interface) in state.
func TestAccPanosVlanEntry_Basic(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)
	templateName := fmt.Sprintf("%s-tmpl", prefix)
	vlanName := fmt.Sprintf("%s-vlan", prefix)
	macAddr := "00:30:48:55:66:77"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPanosVlanEntryBasic,
				ConfigVariables: map[string]config.Variable{
					"template_name": config.StringVariable(templateName),
					"vlan_name":     config.StringVariable(vlanName),
					"mac_addr":      config.StringVariable(macAddr),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_vlan_entry.test",
						tfjsonpath.New("vlan"),
						knownvalue.StringExact(vlanName),
					),
					statecheck.ExpectKnownValue(
						"panos_vlan_entry.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact(macAddr),
					),
					statecheck.ExpectKnownValue(
						"panos_vlan_entry.test",
						tfjsonpath.New("interface"),
						knownvalue.StringExact("ethernet1/1"),
					),
				},
			},
		},
	})
}

const testAccPanosVlanEntryBasic = `
variable "template_name" { type = string }
variable "vlan_name"     { type = string }
variable "mac_addr"      { type = string }

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

resource "panos_vlan" "test" {
  depends_on = [panos_template.test]

  location = {
    template = {
      name = panos_template.test.name
    }
  }

  name       = var.vlan_name
  interfaces = [panos_ethernet_interface.eth1.name]
}

resource "panos_vlan_entry" "test" {
  location = {
    template = {
      name = panos_template.test.name
    }
  }

  vlan      = panos_vlan.test.name
  name      = var.mac_addr
  interface = panos_ethernet_interface.eth1.name
}
`

// TestAccPanosVlanEntry_NoInterface creates a static-MAC entry without setting
// the optional interface attribute.
func TestAccPanosVlanEntry_NoInterface(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)
	templateName := fmt.Sprintf("%s-tmpl", prefix)
	vlanName := fmt.Sprintf("%s-vlan", prefix)
	macAddr := "00:30:48:aa:00:01"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPanosVlanEntryNoInterface,
				ConfigVariables: map[string]config.Variable{
					"template_name": config.StringVariable(templateName),
					"vlan_name":     config.StringVariable(vlanName),
					"mac_addr":      config.StringVariable(macAddr),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_vlan_entry.test",
						tfjsonpath.New("vlan"),
						knownvalue.StringExact(vlanName),
					),
					statecheck.ExpectKnownValue(
						"panos_vlan_entry.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact(macAddr),
					),
					statecheck.ExpectKnownValue(
						"panos_vlan_entry.test",
						tfjsonpath.New("interface"),
						knownvalue.Null(),
					),
				},
			},
		},
	})
}

const testAccPanosVlanEntryNoInterface = `
variable "template_name" { type = string }
variable "vlan_name"     { type = string }
variable "mac_addr"      { type = string }

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

resource "panos_vlan_entry" "test" {
  location = {
    template = {
      name = panos_template.test.name
    }
  }

  vlan = panos_vlan.test.name
  name = var.mac_addr
}
`

// TestAccPanosVlanEntry_MultipleEntries creates two static-MAC entries sharing
// one parent VLAN.  Each entry targets a different ethernet interface.
func TestAccPanosVlanEntry_MultipleEntries(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)
	templateName := fmt.Sprintf("%s-tmpl", prefix)
	vlanName := fmt.Sprintf("%s-vlan", prefix)
	mac1 := "00:30:48:10:20:30"
	mac2 := "00:30:48:40:50:60"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPanosVlanEntryMultipleEntries,
				ConfigVariables: map[string]config.Variable{
					"template_name": config.StringVariable(templateName),
					"vlan_name":     config.StringVariable(vlanName),
					"mac1":          config.StringVariable(mac1),
					"mac2":          config.StringVariable(mac2),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					// First entry
					statecheck.ExpectKnownValue(
						"panos_vlan_entry.mac1",
						tfjsonpath.New("vlan"),
						knownvalue.StringExact(vlanName),
					),
					statecheck.ExpectKnownValue(
						"panos_vlan_entry.mac1",
						tfjsonpath.New("name"),
						knownvalue.StringExact(mac1),
					),
					statecheck.ExpectKnownValue(
						"panos_vlan_entry.mac1",
						tfjsonpath.New("interface"),
						knownvalue.StringExact("ethernet1/1"),
					),
					// Second entry
					statecheck.ExpectKnownValue(
						"panos_vlan_entry.mac2",
						tfjsonpath.New("vlan"),
						knownvalue.StringExact(vlanName),
					),
					statecheck.ExpectKnownValue(
						"panos_vlan_entry.mac2",
						tfjsonpath.New("name"),
						knownvalue.StringExact(mac2),
					),
					statecheck.ExpectKnownValue(
						"panos_vlan_entry.mac2",
						tfjsonpath.New("interface"),
						knownvalue.StringExact("ethernet1/2"),
					),
				},
			},
		},
	})
}

const testAccPanosVlanEntryMultipleEntries = `
variable "template_name" { type = string }
variable "vlan_name"     { type = string }
variable "mac1"          { type = string }
variable "mac2"          { type = string }

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

resource "panos_vlan_entry" "mac1" {
  location = {
    template = {
      name = panos_template.test.name
    }
  }

  vlan      = panos_vlan.test.name
  name      = var.mac1
  interface = panos_ethernet_interface.eth1.name
}

resource "panos_vlan_entry" "mac2" {
  location = {
    template = {
      name = panos_template.test.name
    }
  }

  vlan      = panos_vlan.test.name
  name      = var.mac2
  interface = panos_ethernet_interface.eth2.name
}
`

// TestAccPanosVlanEntry_VlanForceNew verifies that changing the `vlan`
// xpath_variable causes Terraform to destroy-and-recreate the entry rather
// than performing an in-place update.  Two VLANs are declared up-front and a
// ConfigVariables-driven `vlan_name` variable selects which one the entry
// belongs to.
func TestAccPanosVlanEntry_VlanForceNew(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)
	templateName := fmt.Sprintf("%s-tmpl", prefix)
	vlan1Name := fmt.Sprintf("%s-vlan1", prefix)
	vlan2Name := fmt.Sprintf("%s-vlan2", prefix)
	macAddr := "00:30:48:aa:bb:cc"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			// Step 1 – entry lives under vlan1
			{
				Config: testAccPanosVlanEntryVlanForceNew,
				ConfigVariables: map[string]config.Variable{
					"template_name": config.StringVariable(templateName),
					"vlan1_name":    config.StringVariable(vlan1Name),
					"vlan2_name":    config.StringVariable(vlan2Name),
					"vlan_name":     config.StringVariable(vlan1Name),
					"mac_addr":      config.StringVariable(macAddr),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_vlan_entry.test",
						tfjsonpath.New("vlan"),
						knownvalue.StringExact(vlan1Name),
					),
					statecheck.ExpectKnownValue(
						"panos_vlan_entry.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact(macAddr),
					),
				},
			},
			// Step 2 – move entry to vlan2; should replace the resource
			{
				Config: testAccPanosVlanEntryVlanForceNew,
				ConfigVariables: map[string]config.Variable{
					"template_name": config.StringVariable(templateName),
					"vlan1_name":    config.StringVariable(vlan1Name),
					"vlan2_name":    config.StringVariable(vlan2Name),
					"vlan_name":     config.StringVariable(vlan2Name),
					"mac_addr":      config.StringVariable(macAddr),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_vlan_entry.test",
						tfjsonpath.New("vlan"),
						knownvalue.StringExact(vlan2Name),
					),
					statecheck.ExpectKnownValue(
						"panos_vlan_entry.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact(macAddr),
					),
				},
			},
		},
	})
}

const testAccPanosVlanEntryVlanForceNew = `
variable "template_name" { type = string }
variable "vlan1_name"    { type = string }
variable "vlan2_name"    { type = string }
variable "vlan_name"     { type = string }
variable "mac_addr"      { type = string }

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

# vlan1 and vlan2 use distinct member interfaces because PAN-OS enforces
# uniquein on <interface>/<member> — a single interface cannot belong to
# two VLANs simultaneously.
resource "panos_vlan" "vlan1" {
  depends_on = [panos_template.test]

  location = {
    template = {
      name = panos_template.test.name
    }
  }

  name       = var.vlan1_name
  interfaces = [panos_ethernet_interface.eth1.name]
}

resource "panos_vlan" "vlan2" {
  depends_on = [panos_template.test]

  location = {
    template = {
      name = panos_template.test.name
    }
  }

  name       = var.vlan2_name
  interfaces = [panos_ethernet_interface.eth2.name]
}

# Explicit depends_on because vlan is a raw variable (var.vlan_name) rather
# than an attribute reference, so Terraform has no implicit dependency on
# the parent VLANs. Without this, the entry's create may run before either
# vlan exists; PAN-OS would then auto-create the parent vlan xpath as a
# side-effect of the deep set, and the subsequent panos_vlan create fails
# with "entry already exists".
resource "panos_vlan_entry" "test" {
  depends_on = [panos_vlan.vlan1, panos_vlan.vlan2]

  location = {
    template = {
      name = panos_template.test.name
    }
  }

  vlan = var.vlan_name
  name = var.mac_addr
}
`

// TestAccPanosVlanEntry_Import verifies that a VLAN entry can be imported
// using its base64-encoded JSON import ID that carries the location, the
// parent vlan name, and the MAC-address entry name.
func TestAccPanosVlanEntry_Import(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)
	templateName := fmt.Sprintf("%s-tmpl", prefix)
	vlanName := fmt.Sprintf("%s-vlan", prefix)
	macAddr := "00:30:48:cc:dd:ee"

	importStateGenerateID := func(state *terraform.State) (string, error) {
		importState := map[string]any{
			"location": map[string]any{
				"template": map[string]any{
					"name":            templateName,
					"panorama_device": "localhost.localdomain",
					"ngfw_device":     "localhost.localdomain",
				},
			},
			"vlan": vlanName,
			"name": macAddr,
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
				Config: testAccPanosVlanEntryImportConfig,
				ConfigVariables: map[string]config.Variable{
					"template_name": config.StringVariable(templateName),
					"vlan_name":     config.StringVariable(vlanName),
					"mac_addr":      config.StringVariable(macAddr),
				},
			},
			{
				Config: testAccPanosVlanEntryImportStep,
				ConfigVariables: map[string]config.Variable{
					"template_name": config.StringVariable(templateName),
					"vlan_name":     config.StringVariable(vlanName),
					"mac_addr":      config.StringVariable(macAddr),
				},
				ResourceName:      "panos_vlan_entry.imported",
				ImportState:       true,
				ImportStateIdFunc: importStateGenerateID,
			},
		},
	})
}

const testAccPanosVlanEntryImportConfig = `
variable "template_name" { type = string }
variable "vlan_name"     { type = string }
variable "mac_addr"      { type = string }

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

resource "panos_vlan" "test" {
  depends_on = [panos_template.test]

  location = {
    template = {
      name = panos_template.test.name
    }
  }

  name       = var.vlan_name
  interfaces = [panos_ethernet_interface.eth1.name]
}

resource "panos_vlan_entry" "test" {
  location = {
    template = {
      name = panos_template.test.name
    }
  }

  vlan      = panos_vlan.test.name
  name      = var.mac_addr
  interface = panos_ethernet_interface.eth1.name
}
`

const testAccPanosVlanEntryImportStep = `
variable "template_name" { type = string }
variable "vlan_name"     { type = string }
variable "mac_addr"      { type = string }

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

resource "panos_vlan" "test" {
  depends_on = [panos_template.test]

  location = {
    template = {
      name = panos_template.test.name
    }
  }

  name       = var.vlan_name
  interfaces = [panos_ethernet_interface.eth1.name]
}

resource "panos_vlan_entry" "imported" {
  location = {
    template = {
      name = panos_template.test.name
    }
  }

  vlan      = panos_vlan.test.name
  name      = var.mac_addr
  interface = panos_ethernet_interface.eth1.name
}
`
