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

// TestAccPanosTemplate_DefaultVsys verifies that default_vsys can be set
// during the initial create step. The CRUD hooks strip default_vsys from the
// create call and set it via a post-create update.
func TestAccPanosTemplate_DefaultVsys(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	location := config.ObjectVariable(map[string]config.Variable{
		"panorama": config.ObjectVariable(map[string]config.Variable{}),
	})

	configVars := map[string]config.Variable{
		"prefix":   config.StringVariable(prefix),
		"location": location,
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:          template_DefaultVsys_Tmpl,
				ConfigVariables: configVars,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_template.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-template", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_template.test",
						tfjsonpath.New("default_vsys"),
						knownvalue.StringExact("vsys1"),
					),
				},
			},
		},
	})
}

const template_DefaultVsys_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "test" {
  location = var.location
  name = "${var.prefix}-template"
  default_vsys = "vsys1"
}
`

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
