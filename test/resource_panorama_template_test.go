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

func TestAccPanosTemplate_RequiredInputs(t *testing.T) {
	t.Parallel()

	resourceName := "acc_test_template"
	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	templateName := fmt.Sprintf("%s-%s", resourceName, nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: makePanosTemplateConfig(resourceName),
				ConfigVariables: map[string]config.Variable{
					"template_name": config.StringVariable(templateName),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_template."+resourceName,
						tfjsonpath.New("name"),
						knownvalue.StringExact(templateName),
					),
				},
			},
		},
	})
}

func makePanosTemplateConfig(label string) string {
	configTpl := `
    variable "template_name" { type = string }

    resource "panos_template" "%s" {
        name = var.template_name

        location = {
            panorama = {
                panorama_device = "localhost.localdomain"
            }
        }
    }
    `
	return fmt.Sprintf(configTpl, label)
}
