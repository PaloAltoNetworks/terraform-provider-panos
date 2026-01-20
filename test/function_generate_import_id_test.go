package provider_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// TestAccFunction_GenerateImportId_Address_DeviceGroup tests the generate_import_id function
// with an address resource in a device-group location.
func TestAccFunction_GenerateImportId_Address_DeviceGroup(t *testing.T) {
	t.Parallel()

	prefix := fmt.Sprintf("test-acc-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0),
		},
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFunctionGenerateImportIdAddressDeviceGroup_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
			},
		},
	})
}

const testAccFunctionGenerateImportIdAddressDeviceGroup_Tmpl = `
variable "prefix" { type = string }

resource "panos_device_group" "test" {
	location = { panorama = {} }
	name = var.prefix
}

resource "panos_address" "test" {
	name = "${var.prefix}-addr"
	location = { device_group = { name = panos_device_group.test.name } }
	ip_netmask = "10.0.0.1/32"
	description = "test address"
}

output "import_id" {
	value = provider::panos::generate_import_id("panos_address", panos_address.test)
}
`

// TestAccFunction_GenerateImportId_EthernetInterface_Template tests the generate_import_id function
// with an ethernet interface resource in a template location.
func TestAccFunction_GenerateImportId_EthernetInterface_Template(t *testing.T) {
	t.Parallel()

	prefix := fmt.Sprintf("test-acc-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0),
		},
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFunctionGenerateImportIdEthernetTemplate_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
			},
		},
	})
}

const testAccFunctionGenerateImportIdEthernetTemplate_Tmpl = `
variable "prefix" { type = string }

resource "panos_template" "test" {
	location = { panorama = {} }
	name = var.prefix
}

resource "panos_ethernet_interface" "test" {
	name = "ethernet1/1"
	location = { template = { name = panos_template.test.name } }
	layer3 = {
		ips = [{name = "192.168.1.1/32"}]
	}
}

output "import_id" {
	value = provider::panos::generate_import_id("panos_ethernet_interface", panos_ethernet_interface.test)
}
`

// TestAccFunction_GenerateImportId_EthernetInterface_Shared tests the generate_import_id function
// with an ethernet interface resource in a shared location using layer3 variant.
// Note: This test is skipped because shared location is not supported for ethernet interfaces on Panorama.
func TestAccFunction_GenerateImportId_EthernetInterface_Shared(t *testing.T) {
	t.Skip("Skipping shared location test - not supported for ethernet interfaces on Panorama")
}

const testAccFunctionGenerateImportIdEthernetShared_Tmpl = `
variable "prefix" { type = string }

resource "panos_ethernet_interface" "test" {
	name = "ethernet1/2"
	location = { shared = {} }
	layer3 = {
		ips = [{name = "192.168.2.1/32"}]
	}
}

output "import_id" {
	value = provider::panos::generate_import_id("panos_ethernet_interface", panos_ethernet_interface.test)
}
`

// TestAccFunction_GenerateImportId_EthernetInterface_TemplateStack tests the generate_import_id function
// with an ethernet interface resource in a template-stack location using layer3 variant.
func TestAccFunction_GenerateImportId_EthernetInterface_TemplateStack(t *testing.T) {
	t.Parallel()

	prefix := fmt.Sprintf("test-acc-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0),
		},
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFunctionGenerateImportIdEthernetTemplateStack_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
			},
		},
	})
}

const testAccFunctionGenerateImportIdEthernetTemplateStack_Tmpl = `
variable "prefix" { type = string }

resource "panos_template" "test" {
	location = { panorama = {} }
	name = "${var.prefix}-tmpl"
}

resource "panos_template_stack" "test" {
	location = { panorama = {} }
	name = var.prefix
	templates = [panos_template.test.name]
}

resource "panos_ethernet_interface" "test" {
	name = "ethernet1/3"
	location = { template_stack = { name = panos_template_stack.test.name } }
	layer3 = {
		ips = [{name = "192.168.3.1/32"}]
	}
}

output "import_id" {
	value = provider::panos::generate_import_id("panos_ethernet_interface", panos_ethernet_interface.test)
}
`

// TestAccFunction_GenerateImportId_EthernetInterface_VirtualWire_Template tests the generate_import_id function
// with an ethernet interface resource using virtual-wire variant in template location.
func TestAccFunction_GenerateImportId_EthernetInterface_VirtualWire_Template(t *testing.T) {
	t.Parallel()

	prefix := fmt.Sprintf("test-acc-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0),
		},
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFunctionGenerateImportIdEthernetVirtualWireTemplate_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
			},
		},
	})
}

const testAccFunctionGenerateImportIdEthernetVirtualWireTemplate_Tmpl = `
variable "prefix" { type = string }

resource "panos_template" "test" {
	location = { panorama = {} }
	name = var.prefix
}

resource "panos_ethernet_interface" "test" {
	name = "ethernet1/4"
	location = { template = { name = panos_template.test.name } }
	virtual_wire = {}
}

output "import_id" {
	value = provider::panos::generate_import_id("panos_ethernet_interface", panos_ethernet_interface.test)
}
`

// TestAccFunction_GenerateImportId_Address_Vsys tests the generate_import_id function
// with an address resource in a vsys location.
// Note: This test is skipped because vsys is only available on NGFW, not Panorama.
// TODO: Enable this test when running against NGFW test environment.
func TestAccFunction_GenerateImportId_Address_Vsys(t *testing.T) {
	t.Skip("Skipping vsys test - only applicable on NGFW, not Panorama")
}

const testAccFunctionGenerateImportIdAddressVsys_Tmpl = `
variable "prefix" { type = string }

resource "panos_address" "test" {
	name = "${var.prefix}-addr"
	location = { vsys = { name = "vsys1" } }
	ip_netmask = "10.0.0.1/32"
	description = "test address"
}

output "import_id" {
	value = provider::panos::generate_import_id("panos_address", panos_address.test)
}
`

// TestAccFunction_GenerateImportId_Address_Shared tests the generate_import_id function
// with an address resource in a shared location.
func TestAccFunction_GenerateImportId_Address_Shared(t *testing.T) {
	t.Parallel()

	prefix := fmt.Sprintf("test-acc-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0),
		},
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFunctionGenerateImportIdAddressShared_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
			},
		},
	})
}

const testAccFunctionGenerateImportIdAddressShared_Tmpl = `
variable "prefix" { type = string }

resource "panos_address" "test" {
	name = "${var.prefix}-addr"
	location = { shared = {} }
	ip_netmask = "10.0.0.1/32"
	description = "test address"
}

output "import_id" {
	value = provider::panos::generate_import_id("panos_address", panos_address.test)
}
`

// TestAccFunction_GenerateImportId_UnsupportedResource tests the generate_import_id function
// with an unsupported resource type.
func TestAccFunction_GenerateImportId_UnsupportedResource(t *testing.T) {
	t.Parallel()

	prefix := fmt.Sprintf("test-acc-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0),
		},
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFunctionGenerateImportIdUnsupportedResource_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ExpectError: regexp.MustCompile("Unsupported resource type"),
			},
		},
	})
}

const testAccFunctionGenerateImportIdUnsupportedResource_Tmpl = `
variable "prefix" { type = string }

resource "panos_address" "test" {
	name = "${var.prefix}-addr"
	location = { shared = {} }
	ip_netmask = "10.0.0.1/32"
	description = "test address"
}

output "import_id" {
	value = provider::panos::generate_import_id("panos_nonexistent_resource", panos_address.test)
}
`

// TestAccFunction_GenerateImportId_MissingName tests the generate_import_id function
// with a resource missing the name attribute.
func TestAccFunction_GenerateImportId_MissingName(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0),
		},
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccFunctionGenerateImportIdMissingName_Tmpl,
				ExpectError: regexp.MustCompile("(?s)name.*attribute.*missing"),
			},
		},
	})
}

const testAccFunctionGenerateImportIdMissingName_Tmpl = `
output "import_id" {
	value = provider::panos::generate_import_id("panos_address", {
		location = { shared = {} }
	})
}
`

// TestAccFunction_GenerateImportId_MissingLocation tests the generate_import_id function
// with a resource missing the location attribute.
func TestAccFunction_GenerateImportId_MissingLocation(t *testing.T) {
	t.Parallel()

	prefix := fmt.Sprintf("test-acc-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0),
		},
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFunctionGenerateImportIdMissingLocation_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ExpectError: regexp.MustCompile("(?s)location.*attribute.*missing"),
			},
		},
	})
}

const testAccFunctionGenerateImportIdMissingLocation_Tmpl = `
variable "prefix" { type = string }

output "import_id" {
	value = provider::panos::generate_import_id("panos_address", {
		name = "${var.prefix}-addr"
	})
}
`

// TestAccFunction_GenerateImportId_MultipleResourceTypes tests the generate_import_id function
// with multiple different resource types in one configuration.
func TestAccFunction_GenerateImportId_MultipleResourceTypes(t *testing.T) {
	t.Parallel()

	prefix := fmt.Sprintf("test-acc-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0),
		},
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFunctionGenerateImportIdMultipleTypes_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
			},
		},
	})
}

const testAccFunctionGenerateImportIdMultipleTypes_Tmpl = `
variable "prefix" { type = string }

resource "panos_device_group" "test" {
	location = { panorama = {} }
	name = var.prefix
}

resource "panos_address" "test" {
	name = "${var.prefix}-addr"
	location = { device_group = { name = panos_device_group.test.name } }
	ip_netmask = "10.0.0.1/32"
	description = "test address"
}

resource "panos_service" "test" {
	name = "${var.prefix}-svc"
	location = { device_group = { name = panos_device_group.test.name } }
	protocol = { tcp = { port = "8080" } }
}

resource "panos_address_group" "test" {
	name = "${var.prefix}-grp"
	location = { device_group = { name = panos_device_group.test.name } }
	static = [panos_address.test.name]
}

output "address_import_id" {
	value = provider::panos::generate_import_id("panos_address", panos_address.test)
}

output "service_import_id" {
	value = provider::panos::generate_import_id("panos_service", panos_service.test)
}

output "address_group_import_id" {
	value = provider::panos::generate_import_id("panos_address_group", panos_address_group.test)
}
`

// TestAccFunction_GenerateImportId_Inline_Address_Shared tests the generate_import_id function
// with an address resource in a shared location using inline object literals.
func TestAccFunction_GenerateImportId_Inline_Address_Shared(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0),
		},
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFunctionGenerateImportIdInlineAddressShared_Tmpl,
			},
		},
	})
}

const testAccFunctionGenerateImportIdInlineAddressShared_Tmpl = `
output "import_id" {
	value = provider::panos::generate_import_id("panos_address", {
		name = "test-inline-addr"
		location = { shared = {} }
	})
}
`

// TestAccFunction_GenerateImportId_Inline_Address_DeviceGroup tests the generate_import_id function
// with an address resource in a device-group location using inline object literals.
func TestAccFunction_GenerateImportId_Inline_Address_DeviceGroup(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0),
		},
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFunctionGenerateImportIdInlineAddressDeviceGroup_Tmpl,
			},
		},
	})
}

const testAccFunctionGenerateImportIdInlineAddressDeviceGroup_Tmpl = `
output "import_id" {
	value = provider::panos::generate_import_id("panos_address", {
		name = "test-inline-addr-dg"
		location = { device_group = { name = "TestDeviceGroup" } }
	})
}
`

// TestAccFunction_GenerateImportId_Inline_EthernetInterface_Template tests the generate_import_id function
// with an ethernet interface resource in a template location using inline object literals.
func TestAccFunction_GenerateImportId_Inline_EthernetInterface_Template(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0),
		},
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFunctionGenerateImportIdInlineEthernetTemplate_Tmpl,
			},
		},
	})
}

const testAccFunctionGenerateImportIdInlineEthernetTemplate_Tmpl = `
output "import_id" {
	value = provider::panos::generate_import_id("panos_ethernet_interface", {
		name = "ethernet1/5"
		location = { template = { name = "TestTemplate" } }
	})
}
`

// TestAccFunction_GenerateImportId_Inline_Service_DeviceGroup tests the generate_import_id function
// with a service resource in a device-group location using inline object literals.
func TestAccFunction_GenerateImportId_Inline_Service_DeviceGroup(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0),
		},
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFunctionGenerateImportIdInlineServiceDeviceGroup_Tmpl,
			},
		},
	})
}

const testAccFunctionGenerateImportIdInlineServiceDeviceGroup_Tmpl = `
output "import_id" {
	value = provider::panos::generate_import_id("panos_service", {
		name = "test-inline-svc"
		location = { device_group = { name = "TestDeviceGroup" } }
	})
}
`

// TestAccFunction_GenerateImportId_Inline_AddressGroup_DeviceGroup tests the generate_import_id function
// with an address group resource in a device-group location using inline object literals.
func TestAccFunction_GenerateImportId_Inline_AddressGroup_DeviceGroup(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0),
		},
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFunctionGenerateImportIdInlineAddressGroupDeviceGroup_Tmpl,
			},
		},
	})
}

const testAccFunctionGenerateImportIdInlineAddressGroupDeviceGroup_Tmpl = `
output "import_id" {
	value = provider::panos::generate_import_id("panos_address_group", {
		name = "test-inline-grp"
		location = { device_group = { name = "TestDeviceGroup" } }
	})
}
`

// TestAccFunction_GenerateImportId_Inline_Zone_Template tests the generate_import_id function
// with a zone resource in a complex template location using inline object literals.
func TestAccFunction_GenerateImportId_Inline_Zone_Template(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0),
		},
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFunctionGenerateImportIdInlineZoneTemplate_Tmpl,
			},
		},
	})
}

const testAccFunctionGenerateImportIdInlineZoneTemplate_Tmpl = `
output "import_id" {
	value = provider::panos::generate_import_id("panos_zone", {
		name = "test-inline-zone"
		location = {
			template = {
				name = "TestTemplate"
				ngfw_device = "TestNGFW"
				vsys = "vsys1"
			}
		}
	})
}
`

// TestAccFunction_GenerateImportId_Inline_MultipleResources tests the generate_import_id function
// with multiple inline objects in a single configuration.
func TestAccFunction_GenerateImportId_Inline_MultipleResources(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0),
		},
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFunctionGenerateImportIdInlineMultipleResources_Tmpl,
			},
		},
	})
}

const testAccFunctionGenerateImportIdInlineMultipleResources_Tmpl = `
output "address_import_id" {
	value = provider::panos::generate_import_id("panos_address", {
		name = "inline-addr-1"
		location = { shared = {} }
	})
}

output "service_import_id" {
	value = provider::panos::generate_import_id("panos_service", {
		name = "inline-svc-1"
		location = { device_group = { name = "DG1" } }
	})
}

output "address_group_import_id" {
	value = provider::panos::generate_import_id("panos_address_group", {
		name = "inline-grp-1"
		location = { device_group = { name = "DG1" } }
	})
}
`

// TestAccFunction_GenerateImportId_Inline_InvalidLocation tests the generate_import_id function
// with an invalid location type using inline object literals.
func TestAccFunction_GenerateImportId_Inline_InvalidLocation(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0),
		},
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccFunctionGenerateImportIdInlineInvalidLocation_Tmpl,
				ExpectError: regexp.MustCompile("unknown field.*invalid_location_type"),
			},
		},
	})
}

const testAccFunctionGenerateImportIdInlineInvalidLocation_Tmpl = `
output "import_id" {
	value = provider::panos::generate_import_id("panos_address", {
		name = "test-invalid-location"
		location = { invalid_location_type = {} }
	})
}
`

// TestAccFunction_GenerateImportId_Inline_Address_Vsys tests the generate_import_id function
// with an address resource in a vsys location using inline object literals.
// Note: This test is skipped because vsys is only available on NGFW, not Panorama.
func TestAccFunction_GenerateImportId_Inline_Address_Vsys(t *testing.T) {
	t.Skip("Skipping vsys test - only applicable on NGFW, not Panorama")
}

const testAccFunctionGenerateImportIdInlineAddressVsys_Tmpl = `
output "import_id" {
	value = provider::panos::generate_import_id("panos_address", {
		name = "test-inline-addr-vsys"
		location = { vsys = { name = "vsys1" } }
	})
}
`
