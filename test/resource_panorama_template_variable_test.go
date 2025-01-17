package provider_test

import (
	"context"
	"fmt"
	"testing"

	sdkErrors "github.com/PaloAltoNetworks/pango/errors"
	"github.com/PaloAltoNetworks/pango/panorama/template_variable"
	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccPanosPanoramaTemplateVariable(t *testing.T) {
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

	resourceName := "acc_test_template"
	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	compareValuesDiffer := statecheck.CompareValue(compare.ValuesDiffer())
	templateName := "acc_codegen"

	templateVariableTypeEntries := make([]resource.TestStep, 0, len(testData))
	for _, testEntry := range testData {
		templateVariableTypeEntries = append(
			templateVariableTypeEntries,
			resource.TestStep{
				Config: makePanoramaTemplateVariableConfig(resourceName),
				ConfigVariables: map[string]config.Variable{
					"name_suffix":     config.StringVariable(nameSuffix),
					"templ_var_type":  config.StringVariable(testEntry.variableType),
					"templ_var_value": config.StringVariable(testEntry.value),
					"templ_name":      config.StringVariable(templateName),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_template_variable."+resourceName,
						tfjsonpath.New("type").AtMapKey(testEntry.variableType),
						knownvalue.StringExact(testEntry.value),
					),
					statecheck.ExpectKnownValue(
						"panos_template_variable."+resourceName,
						tfjsonpath.New("name"),
						knownvalue.StringExact("$tempvar-"+nameSuffix),
					),
					compareValuesDiffer.AddStateValue(
						"panos_template_variable."+resourceName,
						tfjsonpath.New("type"),
					),
				},
			},
		)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		CheckDestroy: testAccPanosPanoramaTemplateVariableDestroy(
			"$tempvar-"+nameSuffix,
			fmt.Sprintf("%s_%s", templateName, nameSuffix),
		),
		Steps: templateVariableTypeEntries,
	})
}

func makePanoramaTemplateVariableConfig(label string) string {
	configTpl := `
    variable "name_suffix" { type = string }
    variable "templ_var_type" { type = string }
    variable "templ_var_value" { type = string }
    variable "templ_name" { type = string }

    resource "panos_template" "%s" {
      name = "${var.templ_name}_${var.name_suffix}"

      location = {
        panorama = {
          panorama_device = "localhost.localdomain"
        }
      }
    }

    resource "panos_template_variable" "%s" {
      location = {
        template = {
          name = panos_template.%s.name
        }
      }

      name        = "$tempvar-${var.name_suffix}"
      description = "Temp variable description"
      type        = {
        "${var.templ_var_type}": var.templ_var_value
      }
    }
    `

	return fmt.Sprintf(configTpl, label, label, label)
}

func testAccPanosPanoramaTemplateVariableDestroy(entryName, templateName string) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		api := template_variable.NewService(sdkClient)
		ctx := context.TODO()

		location := template_variable.NewTemplateLocation()
		location.Template.Template = templateName

		reply, err := api.Read(ctx, *location, entryName, "show")
		if err != nil && !sdkErrors.IsObjectNotFound(err) {
			return fmt.Errorf("reading template variable entry via sdk: %v", err)
		}

		if reply != nil {
			if reply.EntryName() == entryName {
				return fmt.Errorf("template object still exists: %s", entryName)
			}
		}

		return nil
	}
}
