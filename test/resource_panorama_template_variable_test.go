package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccPanosPanoramaTemplateVariable_Basic(t *testing.T) {
	t.Parallel()

	testData := []struct {
		variableType string
		value        string
	}{
		{variableType: "qos_profile", value: "value"},
		{variableType: "egress_max", value: "12.12"},
		{variableType: "link_tag", value: "value"},
		{variableType: "device_priority", value: "12"},
		{variableType: "interface", value: "ethernet1/20"},
		{variableType: "fqdn", value: "www.paloaltonetworks.com"},
		{variableType: "group_id", value: "12"},
		{variableType: "device_id", value: "1"},
		{variableType: "as_number", value: "12"},
		{variableType: "ip_netmask", value: "8.8.8.8/32"},
		{variableType: "ip_range", value: "127.0.0.1-127.0.0.255"},
	}

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	compareValuesDiffer := statecheck.CompareValue(compare.ValuesDiffer())
	templateName := "acc_codegen"

	templateVariableTypeEntries := make([]resource.TestStep, 0, len(testData))
	for _, testEntry := range testData {
		templateVariableTypeEntries = append(
			templateVariableTypeEntries,
			resource.TestStep{
				Config: templateVariable_Basic_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"name_suffix":     config.StringVariable(nameSuffix),
					"templ_var_type":  config.StringVariable(testEntry.variableType),
					"templ_var_value": config.StringVariable(testEntry.value),
					"templ_name":      config.StringVariable(templateName),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_template_variable.example",
						tfjsonpath.New("type").AtMapKey(testEntry.variableType),
						knownvalue.StringExact(testEntry.value),
					),
					statecheck.ExpectKnownValue(
						"panos_template_variable.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact("$tempvar-"+nameSuffix),
					),
					compareValuesDiffer.AddStateValue(
						"panos_template_variable.example",
						tfjsonpath.New("type"),
					),
				},
			},
		)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps:                    templateVariableTypeEntries,
	})
}

func TestAccPanosPanoramaTemplateVariable_Override(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: templateVariable_Stack_Override_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_template_variable.tmpl-var",
						tfjsonpath.New("type").AtMapKey("interface"),
						knownvalue.StringExact("None"),
					),
					statecheck.ExpectKnownValue(
						"panos_template_variable.stack-var",
						tfjsonpath.New("type").AtMapKey("interface"),
						knownvalue.StringExact("ethernet1/1"),
					),
				},
			},
		},
	})
}

const templateVariable_Basic_Tmpl = `
variable "name_suffix" { type = string }
variable "templ_var_type" { type = string }
variable "templ_var_value" { type = string }
variable "templ_name" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }

  name = "${var.templ_name}_${var.name_suffix}"
}

resource "panos_template_variable" "example" {
  location = {
    template = {
      name = panos_template.example.name
    }
  }

  name        = "$tempvar-${var.name_suffix}"
  description = "Temp variable description"
  type        = {
    "${var.templ_var_type}": var.templ_var_value
  }
}
`

const templateVariable_Stack_Override_Tmpl = `
variable "prefix" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }

  name = var.prefix
}

resource "panos_template_stack" "example" {
  location = { panorama = {} }

  templates = [panos_template.example.name]

  name = "${var.prefix}-stack"
}

resource "panos_template_variable" "tmpl-var" {
  location = { template = { name = panos_template.example.name } }

  name = format("$%s", var.prefix)
  type = { interface = "None" }
}

resource "panos_template_variable" "stack-var" {
  location = { template_stack = { name = panos_template_stack.example.name } }

  name = format("$%s", var.prefix)
  type = { interface = "ethernet1/1" }
}
`
